package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
    "errors"
)

// Application constants, defining host, port, and protocol.
const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)


type CommandType int32
const (
    Add CommandType = 0
    Mult CommandType = 1
)

type Command struct {
    tpe CommandType
    param1 int32
    param2 int32
}

func main() {
	// Start the server and listen for incoming connections.
	fmt.Println("Starting " + connType + " server on " + connHost + ":" + connPort)
	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	// run loop forever, until exit.
	for {
		// Listen for an incoming connection.
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}
		fmt.Println("Client connected.")

		// Print client connection address.
		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

		// Handle connections concurrently in a new goroutine.
		go handleConnection(c)
	}
}

func handleConnection(conn net.Conn) {
    for {
        // Buffer client input until a newline.
        buffer, err := bufio.NewReader(conn).ReadBytes('\n')
        if err != nil {
            fmt.Println("Client left.")
            return
        }

        // Print response message, stripping newline character.
        log.Println("Client message:", string(buffer[:len(buffer)-1]))

        // Send response message to the client.
        _, err = conn.Write(buffer)
        if err != nil {
            fmt.Println("Error sending response:", err)
            return
        }
    }
}


func executeCommand(cmd Command) (int32,error) {
    switch cmd.tpe {
    case Add:
        return cmd.param1 + cmd.param2, nil;
    case Mult:
        return cmd.param1 + cmd.param2, nil;
    default:
        return -1, errors.New("Command not known");
    }
}

