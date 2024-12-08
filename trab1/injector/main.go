package main

import (
    "fmt"
    "net"
    "os"
)

const (
    connType = "tcp"
    token = "123123\n"
)

func main() {
    var port string
    port = os.Args[1]

    // Start the client and connect to the server.
	fmt.Println("Connecting to", connType, "server", port)
	conn, err := net.Dial(connType, port)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
    
    fmt.Println("Sending Token")
    // Send to socket connection.
    conn.Write([]byte(token))

    conn.Close()
}