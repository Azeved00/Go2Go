package peer

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "time"
    "p2p/poisson"
    "p2p/pmap"
    "math/rand"
    "sync"
)

const (
    connType = "tcp"
    time_limit = 10 * time.Second
    lambda = 5.0
)

type Peer struct {
    poisson     *poisson.PoissonProcess
    rng         *rand.Rand

    self_port   string
    listener    net.Listener

    conns       map[string]net.Conn
    peer_map    pmap.PMap
	peer_map_mutex sync.Mutex
}

func (p *Peer) Close() {
    fmt.Println("exiting")
	for _, conn := range p.conns {
		conn.Close()
	}
}

// Create a new peer to peer connection
// the 'port' parameter is fo r
// waits for a peer to be connected to it
// waits to be connected to a peer
func New(port string) Peer {
    fmt.Println("Starting " + connType + " peer on port " + port)

    //set up prev connection
    l, err := net.Listen(connType, "0.0.0.0:"+port)

    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }

    //setup poisson
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	poissonProcess, err := poisson.NewPoissonProcess(lambda, rng)
    
    return Peer {
        poisson: poissonProcess,
        rng: rng,
        self_port: port, 
        listener: l, 
        peer_map: pmap.NewPeerMap("localhost:"+port),
        conns: map[string]net.Conn{},
    }
}

// wait for someone to connect to us
func (p *Peer) ConnectTo(next_addr string) {
    if p.conns[next_addr] != nil {
        p.peer_map.UpdatePeer(next_addr)
        return
    }

    fmt.Println("",p.conns)
    fmt.Println("Trying to connect to next peer: "+next_addr)

    var next_c net.Conn
    ts := time.Now().Add(time_limit)
    for {
        var err error
        next_c,err = net.Dial(connType, next_addr)
        
        if err == nil { break }


        if time.Now().After(ts) { 
            fmt.Println("Connection Timed out, exiting program")
            os.Exit(0) 
        }
    }

    p.conns[next_addr] = next_c
    p.peer_map.UpdatePeer(next_addr)
}

func (p *Peer) Listen() {
	for {
		// Accept a new connection
		conn, err := p.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a separate goroutine
		go p.handleCon(conn)
	}
}

func (p *Peer) handleCon(conn net.Conn) {
	fmt.Println("New client connected:", conn.RemoteAddr())

	// Read messages from the client
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Bytes()

        pmap,err := pmap.Deserialize(msg)
        if err != nil {
            fmt.Println("Failed to serialize")
            continue
        }
        p.ConnectTo(pmap.Addr)
        fmt.Println("merging serialization")
        p.peer_map.Merge(pmap)
        fmt.Println("", p.peer_map.PrettyPrint())

	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Connection error:", err)
	}

	fmt.Println("Client disconnected:", conn.RemoteAddr())
    conn.Close()
}


func (p *Peer) Poison() {
    for {
        message, err := p.peer_map.Serialize()
        if err != nil {
            fmt.Printf("%s\n",err)
            continue
        }

        keys := make([]string, 0, len(p.conns))
        fmt.Println("Conns to send:", keys)
        for key := range p.conns {
            keys = append(keys, key)
        }
        if len(keys) <= 0 {
            continue
        }

        randomIndex := p.rng.Intn(len(keys))
        addr := keys[randomIndex]

        conn := p.conns[addr]

        fmt.Printf("Sending message to %s...\n", addr)

        _, err = conn.Write(message)
        if err != nil {
            fmt.Printf("Failed to send message %v\n", err)
            continue
        }

        waitTime := p.poisson.TimeForNextEvent()
		time.Sleep(time.Duration(waitTime * float64(time.Minute)))
    }
}

