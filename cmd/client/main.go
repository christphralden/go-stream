package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/christpheralden/go-stream/pkg/client"
)


func main() {
	cl := client.NewClient(
		client.WithProtocol(client.TCP),
	)

	if err := cl.Dial(); err != nil {
		log.Println("Error connecting to server:", err)
		return
	}

	fmt.Println("Input your messages:")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()

		if message == "quit" {
			fmt.Println("Disconnecting from server")
			cl.Stop()
			break
		}

		if err := cl.SendMessage(message); err != nil {
			fmt.Println("Error sending message:", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from input:", err)
	}
}

