package main

import (
    "time"
    "fmt"
    "flag"
    "strconv"
    "p2p/peer"
)

const (
    default_peer_port = 8180
)

func main() {
	peer_port := flag.Int("p", default_peer_port, "Port where peer will listen")
	flag.Parse()

    p := peer.New(strconv.Itoa(*peer_port))
    fmt.Println("Created peer")

    defer p.Close()

    go p.Listen()
    fmt.Println("Listening for Connections")

    go p.Poison()
    fmt.Println("Poisson Loop Initiated")

    go p.ProcessQueue()
    fmt.Println("Starting processing the queue")

    time.Sleep(time.Duration(3 * time.Minute))
    fmt.Println("Program finished time limit, exiting")
    return 
}
