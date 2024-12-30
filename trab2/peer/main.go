package main

import (
    "fmt"
    "flag"
    "strconv"
    "p2p/peer"
    "time"
)

const (
    default_peer_port = 8180
)

func main() {
	peer_port := flag.Int("p", default_peer_port, "Port where peer will listen")
	address := flag.String("n", "", "Address to connect to (e.g., localhost:8080)")
	flag.Parse()

    p := peer.New(strconv.Itoa(*peer_port))
    fmt.Println("Created peer")

    defer p.Close()

    go p.Listen()
    fmt.Println("Listening for Connections")

	if *address != "" {
        p.ConnectTo(*address)
	}


    go p.Poison()
    fmt.Println("Poisson Loop Initiated")

    time.Sleep(time.Duration(3 * time.Minute))
    fmt.Println("Program finished time limit, exiting")
    return 
}
