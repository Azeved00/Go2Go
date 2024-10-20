package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
)

const (
    connHost = "localhost"
    connPort = "8080"
    connType = "tcp"
)

func main() {
    fmt.Println("Starting " + connType + " server on " + connHost + ":" + connPort)
    l, err := net.Listen(connType, connHost+":"+connPort)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    defer l.Close()

    for {
        c, err := l.Accept()
        if err != nil {
            fmt.Println("Error connecting:", err.Error())
            return
        }
        fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

        go handleConnection(c)
    }
}

func handleConnection(conn net.Conn) {
    buffer, err := bufio.NewReader(conn).ReadBytes('\n')

    if err != nil {
        fmt.Println("Client left.")
        conn.Close()
        return
    }

    log.Println("Client message:", string(buffer[:len(buffer)-1]))

    conn.Write(buffer)

    handleConnection(conn)
}

func main_client() {
	// Start the client and connect to the server.
	fmt.Println("Connecting to", connType, "server", connHost+":"+connPort)
	conn, err := net.Dial(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}

	// Create new reader from Stdin.
	reader := bufio.NewReader(os.Stdin)

	// run loop forever, until exit.
	for {
		// Prompting message.
		fmt.Print("Text to send: ")

		// Read in input until newline, Enter key.
		input, _ := reader.ReadString('\n')

		// Send to socket connection.
		conn.Write([]byte(input))

		// Listen for relay.
		message, _ := bufio.NewReader(conn).ReadString('\n')

		// Print server relay.
		log.Print("Server relay: " + message)
	}
}

