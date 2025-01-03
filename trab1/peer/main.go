package main

import (
    "fmt"
    "flag"
    "os"
    "strconv"
    "p2p/peer"
    "time"
)

const (
    default_server_addr = "localhost:8080"
    default_peer_port = 8180
)

func main() {
	peer_port := flag.Int("p", default_peer_port, "Port where peer will listen")
    	next_peer_port := flag.Int("n", default_peer_port, "Port wich this peer will connect to")
    	server_addr := flag.String("s", default_server_addr, "Address of the server")
	flag.Parse()

	args := flag.Args()

    if len(args) < 1 {
		fmt.Println("Error: A required parameters is missing.")
        	fmt.Println("Usage: " + os.Args[0] + " <addr to connect to>")
		os.Exit(1)
    }

   
    p := peer.New(
            strconv.Itoa(*peer_port), 
            args[0]+":"+strconv.Itoa(*next_peer_port),
            *server_addr)



    defer p.Close()
    fmt.Println("Created peer")

    go p.Poison()
    fmt.Println("Poisson Loop Initiated")

    go p.Loop()
    fmt.Println("Token Loop Initiated")

    time.Sleep(time.Duration(2 * time.Minute))
    fmt.Println("Program finished time limit, exiting")
    return 
}
