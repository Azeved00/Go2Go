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
    for {
        conn, err := net.Dial(connType, port)
        if err == nil {
            fmt.Println("Sending Token")
            // Send to socket connection.
            conn.Write([]byte(token))

            conn.Close()
        }
        fmt.Println("Couldn't Connect, trying again")
    }
    

}
