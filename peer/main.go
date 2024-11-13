package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "time"
	"sync"
)

const (
    connHost = "localhost"
    connType = "tcp"
    server_addr = "localhost:8080"
)

func CheckToken(s string) bool {
    return s == "123123"
}

type Peer struct {
    mu sync.Mutex
    req []string

    self_port string

    server_con net.Conn

    prev_l    net.Listener
    prev_conn net.Conn

    next_addr string
    next_conn net.Conn  
}

func (p *Peer) Close() {
    p.next_conn.Close()
    p.prev_conn.Close()
    p.server_con.Close()
}

// Create a new peer to peer connection
// the 'port' parameter is fo r
// waits for a peer to be connected to it
// waits to be connected to a peer
func New(port string, next_addr string) Peer {
    fmt.Println("Starting " + connType + " peer on " + port)

    l, err := net.Listen(connType, connHost+":"+port)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }

    var next_c net.Conn
    for {
        var err error
        next_c,err = net.Dial(connType, server_addr)
        if err == nil {
            break
        }
        time.Sleep(2 * time.Second)
        fmt.Println("Error connecting:", err.Error())
    }

    return Peer {
        req:[]string{},
        self_port: port, 
        server_con: next_c, 
        prev_l: l, 
        prev_conn: nil,
        next_addr: next_addr, 
        next_conn:nil,
    }
}

// wait for someone to connect to us
func (p *Peer) ConnectNext() {
    var next_c net.Conn
    for {
        var err error
        next_c,err = net.Dial(connType, connHost+":"+p.next_addr)
        if err == nil {
            break
        }
        fmt.Println("Error connecting:", err.Error())
        time.Sleep(2 * time.Second)
    }
    p.next_conn = next_c
}

//connect to someone else
func (p *Peer) ConnectPrev() {
    fmt.Println("Looking for a new connection")

    var c net.Conn
    var err error
    for {
        c, err = p.prev_l.Accept()
        if err == nil {
            break 
        }
        fmt.Println("Error connecting:", err.Error())
        time.Sleep(2 * time.Second)
    }
    p.prev_conn = c
    fmt.Println("Client " + p.prev_conn.RemoteAddr().String() + " connected.")
}

// Infinitly receives message from someone connected to us
func (p *Peer) Loop()  {
    for {

        if p.prev_conn == nil {
            p.ConnectPrev()
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
        }
        p.req = []string{}

        fmt.Println("Requests Sent")

        p.mu.Unlock()

        time.Sleep(2 * time.Second) 


        if p.next_conn == nil {
            p.ConnectNext()
        }

        _, err = p.next_conn.Write(append(buffer, byte('\n')))
        if err != nil {
            fmt.Println("Token not sent")
        }

    }
}

func (p *Peer) Poison() {
    for {
        p.mu.Lock()
        p.req = append(p.req, "abcd abcd")
        p.mu.Unlock()
        time.Sleep(2 * time.Second)
    }

}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: " + os.Args[0] + " <machine addr> " + " <addr to connect to>")
        return
    }
    p := New(os.Args[1], os.Args[2])
    defer p.Close()
    fmt.Println("Created peer")

    go p.Poison()
    fmt.Println("Poisson Loop Initiated")

    go p.Loop()
    fmt.Println("Token Loop Initiated")

    for {}
}
