package peer

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "time"
    "sync"
    "p2p/poisson"
    "math/rand"
)

const (
    connType = "tcp"
    time_limit = 10 * time.Second
)

func CheckToken(s string) bool {
    return s == "123123"
}

type Peer struct {
    mu          sync.Mutex
    req         []string
    poisson     *poisson.PoissonProcess
    rng         *rand.Rand

    self_port   string

    server_con  net.Conn

    prev_l      net.Listener
    prev_conn   net.Conn

    next_addr   string
    next_conn   net.Conn  

    finishing   bool
}

func (p *Peer) Close() {
    fmt.Println("exiting")

    p.finishing = true
    time.Sleep(time.Duration(10 * time.Second))

    p.next_conn.Close()
    p.prev_conn.Close()
    p.server_con.Close()
}

// Create a new peer to peer connection
// the 'port' parameter is fo r
// waits for a peer to be connected to it
// waits to be connected to a peer
func New(port string, next_addr string, server_addr string) Peer {
    fmt.Println("Starting " + connType + " peer on port " + port)

    //set up prev connection
    l, err := net.Listen(connType, "0.0.0.0:"+port)

    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }

    //setup next connection
    var next_c net.Conn
    for {
        var err error
        next_c,err = net.Dial(connType, server_addr)
        if err == nil {
            break
        }
        time.Sleep(2 * time.Second)
        fmt.Println("Error connecting to server:", err.Error())
    }

    //setup poisson
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	lambda := 4.0

	poissonProcess, err := poisson.NewPoissonProcess(lambda, rng)

    return Peer {
        poisson: poissonProcess,
        rng: rng,
        req:[]string{},
        self_port: port, 
        server_con: next_c, 
        prev_l: l, 
        prev_conn: nil,
        next_addr: next_addr, 
        next_conn:nil,
	finishing: false,
    }
}

// wait for someone to connect to us
func (p *Peer) ConnectNext() {
    fmt.Println("Trying to connect to next peer("+p.next_addr+")")

    var next_c net.Conn
    ts := time.Now().Add(time_limit)
    for {
        var err error
        next_c,err = net.Dial(connType, p.next_addr)
        if err == nil { break }


        if time.Now().After(ts) { 
            fmt.Println("Connection Timed out, exiting program")
            os.Exit(0) 
        }
    }
    p.next_conn = next_c
}

//connect to someone else
func (p *Peer) ConnectPrev() {
    fmt.Println("Waiting For previous peer to connect")

    var c net.Conn
    var err error
    deadline := time.Now().Add(time_limit)
    p.prev_l.(*net.TCPListener).SetDeadline(deadline)
    for {
        c, err = p.prev_l.Accept()
        if err == nil {
            break 
        }
		if opErr, ok := err.(net.Error); ok && opErr.Timeout() {
            fmt.Println("Waited too long for a connected, exiting...")
            os.Exit(0)	
		}

    }
    p.prev_conn = c
    fmt.Println("Client " + p.prev_conn.RemoteAddr().String() + " connected.")
}

// Infinitly receives message from someone connected to us
func (p *Peer) Loop()  {
    for {
	if finishing{
		break
	}

        if p.prev_conn == nil {
            p.ConnectPrev()
        }
        if p.next_conn == nil {
            p.ConnectNext()
        }

        buffer, err := bufio.NewReader(p.prev_conn).ReadBytes('\n')
        if err != nil {
            fmt.Println(err)
            p.prev_conn.Close()
            p.prev_conn = nil
            continue
        }

        message := string(buffer[:len(buffer)-1])
        fmt.Println(message)
        if !CheckToken(message) {
            continue;
        }

        p.mu.Lock()

        fmt.Println("Token Received")

        for _, value := range p.req {
            fmt.Println("Sending Req")
            p.server_con.Write(append([]byte(value), byte('\n')))
            buffer, err := bufio.NewReader(p.server_con).ReadBytes('\n')
            if err != nil {
                fmt.Println("Server Error")
                continue
            }


            response := string(buffer[:len(buffer)-1])
            fmt.Println("res:", response)
        }
        p.req = []string{}

        fmt.Println("Requests Sent")

        p.mu.Unlock()

        time.Sleep(500 * time.Millisecond) 


        _, err = p.next_conn.Write(append(buffer, byte('\n')))
        if err != nil {
            fmt.Println("Token not sent")
        }

    }
}

func (p *Peer) Poison() {
    for {
	if finishing {
		break
	}
        cmd := p.GenCommand()
        p.mu.Lock()
        p.req = append(p.req, cmd)
        p.mu.Unlock()


        waitTime := p.poisson.TimeForNextEvent()
		time.Sleep(time.Duration(waitTime * float64(time.Minute)))
    }
}

func (p *Peer) GenCommand() string {
        cmdn := p.rng.Int31n(3)
        cmd := ""
        switch cmdn {
        case 0:
            cmd = "add"
        case 1:
            cmd = "sub"
        case 2:
            cmd = "mult"
        default:
            cmd = "add"

        }
        param1 := p.rng.NormFloat64()
        param2 := p.rng.NormFloat64()
        
	    command := fmt.Sprintf("%s %f %f\n",cmd, param1, param2)
        return command
}

