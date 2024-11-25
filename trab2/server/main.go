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
    Sub CommandType = 2
)

type Command struct {
    tpe CommandType
    param1 float64
    param2 float64
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

        log.Println(" & ", string(buffer[:len(buffer)-1]))

        cmd,err := parseCommand(buffer)
        if err!= nil {
            fmt.Println("Error parsing command", err)
            _, err = conn.Write([]byte("error\n"))
            return
        }

        res,err := executeCommand(cmd)
        if err!= nil {
            fmt.Println("Error executing command:", err)
            _, err = conn.Write([]byte("error\n"))
            return
        }

        log.Println("-> ", res)

        // Send response message to the client.
        s := fmt.Sprintf("%f\n", res) 
        _, err = conn.Write([]byte(s))
        if err != nil {
            fmt.Println("Error sending response:", err)
            return
        }
    }
}

func parseCommand(buffer []byte) (Command,error) {
    var command string
    var num1 float64
    var num2 float64
    var cmd CommandType

	input := string(buffer)
	_, err := fmt.Sscanf(input, "%s %f %f", &command, &num1, &num2)
    if err != nil {
        return Command{tpe: Add, param1:0.0, param2:0.0}, err
    }

    switch command {
    case "add":
        cmd = Add
    case "mult":
        cmd = Mult
    case "sub":
        cmd = Sub
    default:
        return Command{tpe: Add, param1:0.0, param2:0.0}, 
            errors.New("Command not known");
    }

    return Command{
        tpe: cmd,
        param1: num1,
        param2: num2,
    }, nil

}

func executeCommand(cmd Command) (float64,error) {
    switch cmd.tpe {
    case Add:
        return cmd.param1 + cmd.param2, nil;
    case Mult:
        return cmd.param1 * cmd.param2, nil;
    default:
        return -1.0, errors.New("Command not known");
    }
}

