package main

import (
    "time"
    "fmt"
    "flag"
    "strconv"
    "p2p/peer"
    "os"
)

const (
    default_peer_port = 8180
)

func main() {
	peer_port := flag.Int("p", default_peer_port, "Port where peer will listen")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Error: A required parameter is missing.")
		fmt.Println("Usage: " + os.Args[0] + " <addr of this machine>")
		os.Exit(1)
	}

    p := peer.New(args[0], strconv.Itoa(*peer_port))
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
