package peer

import (
    "fmt"
    "net"
    "os"
    "time"
    "p2p/poisson"
    "p2p/lamport"
    "math/rand"
	"encoding/gob"
)

const (
    connType = "tcp"
    time_limit = 10 * time.Second
    shutdown_limit = 15 * time.Second
    lambda = 0.2
)
var peer_addresses = [6]string{
    "localhost:8081", 
    "localhost:8082", 
    "localhost:8083", 
    "localhost:8084", 
    "localhost:8085", 
    "localhost:8086", 
}
var words = []string{
	"apple", "luminous", "gravel", "serenity", "ocean", "thunder",
	"intricate", "blossom", "melody", "vibrant", "shadow", "echo",
	"spark", "horizon", "evergreen",
}


type Peer struct {
    poisson     *poisson.PoissonProcess
    rng         *rand.Rand

    self_port   string
    listener    net.Listener

    conns       ConnMap
    clock       lamport.LampClock

    p_heap      ProcessHeap
}

func (p *Peer) Close() {
    fmt.Println("exiting")
}

// Create a new peer to peer connection
// the 'port' parameter is fo r
// waits for a peer to be connected to it
// waits to be connected to a peer
func New(port string) *Peer {
    fmt.Println("Starting " + connType + " peer on port " + port)

    //set up prev connection
    l, err := net.Listen(connType, "localhost:"+port)

    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }

    //setup poisson
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	poissonProcess, err := poisson.NewPoissonProcess(lambda, rng)
    peer := Peer {
        poisson: poissonProcess,
        rng: rng,
        self_port: port, 
        listener: l, 
        clock: lamport.NewLamportClock("localhost:"+port),
        conns: NewConnMap(),
    }
    for _, addr := range peer_addresses {
        go peer.ConnectTo(addr, false)
    }
    
    return &peer
}

// wait for someone to connect to us
func (p *Peer) ConnectTo(next_addr string, reset bool) {
    if reset {
        p.conns.Set(next_addr, nil)
    }
    if  p.conns.Exists(next_addr) {
        return
    }
    //fmt.Println("",p.conns)
    fmt.Println("Trying to connect to next peer: "+next_addr)

    var enc *gob.Encoder 
    ts := time.Now().Add(time_limit)
    for {
        conn, err := net.Dial(connType, next_addr)
        enc = gob.NewEncoder(conn)
        
        if err == nil { break }


        if time.Now().After(ts) { 
            fmt.Println("Connection Timed out, exiting program")
            return
        }
    }

    p.conns.Set(next_addr, enc)
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


type Message struct {
    Word        string
    Clock_ser   []byte
}
func (p *Peer) handleCon(conn net.Conn) {
	fmt.Println("New client connected:", conn.RemoteAddr())

	decoder := gob.NewDecoder(conn)
	for {
        var message Message
		err := decoder.Decode(&message)
        fmt.Println("Received a message")
	    if err != nil {
			fmt.Printf("Connection closed or error: %v\n", err)
			break
		}

        clock, err := lamport.Deserialize(message.Clock_ser)
        if err != nil {
            fmt.Println("Failed to serialize lamport clock", err)
            continue
        }
        c := clock.Get()

        p.p_heap.Push(HeapObj {
            Counter: c,
            Word: message.Word,
            Addr: clock.Addr,
        })
	}

	fmt.Println("Client disconnected:", conn.RemoteAddr())
    conn.Close()
}

func (p *Peer) Poison() {
	time.Sleep(time.Duration(4 * float64(time.Second)))
    for {
        p.clock.Increment()
        ser_clock, err := p.clock.Serialize()
        if err != nil {
            fmt.Printf("%s\n",err)
            continue
        }
        word := p.SelectRandomWord()
        message := Message{
		    Word:       word,
		    Clock_ser: ser_clock,
	    }

        addrs := p.conns.GetKeys()
        fmt.Println(addrs)
        for _, addr := range addrs {
            if p.conns.Send(addr, message) {
                continue
            } else {
                go p.ConnectTo(addr, true)
            }
        }

        waitTime := p.poisson.TimeForNextEvent()
		time.Sleep(time.Duration(waitTime * float64(time.Second)))
    }
}

func (p *Peer) SelectRandomWord() string {
	randomIndex := p.rng.Intn(len(words))
	return words[randomIndex]
}

func (p *Peer) ProcessQueue() {
    for {
        if p.p_heap.Len() == 0 {
            continue
	    } else {
            item := p.p_heap.Pop()

            clock := lamport.LampClock {
                Counter: item.Counter,
                Addr: item.Addr,
            }
            p.clock.Merge(&clock)


            fmt.Println(clock.PrettyPrint())
            fmt.Println(p.clock.PrettyPrint())
            fmt.Println(item.Word)
        }
    }
}

func (p *Peer) DeadLine() {
    flag := false
    for {
        l :=  p.conns.Len()
        if l == 0 && flag {
            os.Exit(0)
        }
        flag = l == 0

		time.Sleep(time.Duration(shutdown_limit))
    }
}
