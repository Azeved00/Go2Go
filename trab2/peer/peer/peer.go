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
)

type Peer struct {
    poisson     *poisson.PoissonProcess
    rng         *rand.Rand

    self_port   string
    listener    net.Listener

    conns       map[net.Addr]net.Conn
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
    l, err := net.Listen(connType, "localhost:"+port)

    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }

    //setup poisson
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	lambda := 4.0

	poissonProcess, err := poisson.NewPoissonProcess(lambda, rng)
    
    return Peer {
        poisson: poissonProcess,
        rng: rng,
        self_port: port, 
        listener: l, 
        peer_map: pmap.NewPeerMap(),
        conns: map[net.Addr]net.Conn{},
    }
}

// wait for someone to connect to us
func (p *Peer) ConnectTo(next_addr string) {
    fmt.Println("Trying to connect to next peer("+next_addr+")")

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
    p.peer_map.UpdatePeer(next_c)
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
    p.conns[conn.RemoteAddr()] = conn
    p.peer_map.UpdatePeer(conn)

	fmt.Println("New client connected:", conn.RemoteAddr())

	// Read messages from the client
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		fmt.Println("Received from client:", msg)

		// Echo the message back to the client
		_, err := conn.Write([]byte("Server: " + msg + "\n"))
		if err != nil {
			fmt.Println("Error writing to client:", err)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Connection error:", err)
	}

	fmt.Println("Client disconnected:", conn.RemoteAddr())
    delete(p.conns, conn.RemoteAddr())
    p.conns[conn.RemoteAddr()] = conn
    conn.Close()
}


func (p *Peer) Poison() {
    for {
        //fmt.Printf("%s\n",p.peer_map.PrettyPrint())
        message, err := p.peer_map.Serialize()
        if err != nil {
            fmt.Printf("%s\n",err)
            continue
        }

        fmt.Printf("%s\n",message)
        for addr, conn := range p.conns {
            fmt.Printf("Sending message to %s...\n", addr)

            _, err := conn.Write(message)
            if err != nil {
                fmt.Printf("Failed to send message %v\n", err)
                continue
            }
            break
        }

        waitTime := p.poisson.TimeForNextEvent()
		time.Sleep(time.Duration(waitTime * float64(time.Minute)))
    }
}

