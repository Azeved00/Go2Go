package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "time"
)

const (
    connHost = "localhost"
    connType = "tcp"
    token = "123123\n"
    server = "8080"
)

type Peer struct {
    self_port string

    prev_l    net.Listener
    prev_conn net.Conn

    next_addr string
    next_conn net.Conn  
}

func (p *Peer) Close() {
    p.next_conn.Close()
    p.prev_conn.Close()
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

    return Peer {port, l, nil, next_addr, nil}
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
        time.Sleep(10 * time.Second)
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
        time.Sleep(10 * time.Second)
    }
    p.prev_conn = c
    fmt.Println("Client " + p.prev_conn.RemoteAddr().String() + " connected.")
}

// Infinitly receives message from someone connected to us
func (p Peer) Loop()  {
    for {
        time.Sleep(5 * time.Second) 

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
	    p.prev_conn.Write([]byte("ok"))


        message := string(buffer[:len(buffer)-1])
        fmt.Println(message)
        if message == token{
            fmt.Println("Token Received")
        }

        if p.next_conn == nil {
            p.ConnectNext()
        }
        p.next_conn.Write(buffer)
    }
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: " + os.Args[0] + " <machine addr> " + " <addr to connect to>")
        return
    }
    p := New(os.Args[1], os.Args[2])
    defer p.Close()

    p.Loop()
}
