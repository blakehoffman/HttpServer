package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	server, error := net.Listen("tcp", "localhost:9988")

	if error != nil{
		fmt.Println("Error listening:", error.Error())
    	os.Exit(1)
	}

	defer server.Close()

	for {
		connection, error := server.Accept()

		if error != nil{
			fmt.Println("Error accepting client: ", error.Error())
			os.Exit(1)
		}
		
		processClient(connection)
	}
}

func processClient(connection net.Conn){
	defer connection.Close()

	buffer := make([]byte, 1024)
	length, error := connection.Read(buffer)

	if error != nil{
		fmt.Println("Error reading buffer", error.Error())
	}

	fmt.Println("Received x bytes: ", length)

	_, error = connection.Write([]byte("Hello world"))

	if error != nil{
		fmt.Println("Error writing to client")
	}
}