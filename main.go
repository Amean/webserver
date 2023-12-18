// package main

// import (
// 	"net"
// 	"fmt"
// 	"bufio"
// )

// // func getRoot(w http.ResponseWriter, r *http.Request) {
// // 	fmt.Printf("got / request\n")
// // 	io.WriteString(w, "This is my website!\n")
// // }

// // func getHello(w http.ResponseWriter, r *http.Request) {
// // 	fmt.Printf("got /hello request\n")
// // 	io.WriteString(w, "Hello, HTTP!\n")
// // }

// func sendIntegers() {
// 	for i := 0; i<10; i++ {

// 	}
// }
// func recieveIntegers() {

// }

// func main() {
// 	fmt.Println("Start server...")

//   // listen on port 8000
//   ln, _ := net.Listen("tcp", ":8000")

//   // accept connection
//   conn, _ := ln.Accept()

//   // run loop forever (or until ctrl-c)
//   for {
//     // get message, output
//     message, _ := bufio.NewReader(conn).ReadString('\n')
//     fmt.Print("Message Received:", string(message))
//   }
// }

package main

import (
	"fmt"
	"github.com/go-redis/redis/v9" // Redis client for Go
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var ctx = context.Background()

func sendDataToPython() {
	// Function to send integers to Python
	// Establish a connection with the process manager running Python script
	pythonCmd := exec.Command("python_process_manager") // Replace with your Python process manager command
	pythonIn, _ := pythonCmd.StdinPipe()
	defer pythonIn.Close()
	pythonCmd.Start()

	// Send integers to Python
	for i := 0; i < 100; i++ { // You can set a specific limit or condition
		// Generate sequential integers
		data := strconv.Itoa(i)
		_, _ = pythonIn.Write([]byte(data + "\n")) // Send integer data to Python
		time.Sleep(time.Second)                    // Adjust as per your requirement
	}

	pythonCmd.Wait() // Wait for Python process to finish
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// HTTP handler for starting G1
	go sendDataToPython()
	fmt.Fprintf(w, "Sending data to Python...")
}

func receiveFromRedis() {
	// Function to receive data from Redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Replace with your Redis server address
		Password: "",               // No password by default
		DB:       0,                // Default database
	})

	pubsub := client.Subscribe(ctx, "processed_data_channel")
	defer pubsub.Close()

	// Receive and display processed data
	for {
		msg, _ := pubsub.ReceiveMessage(ctx)
		fmt.Println("Received processed data:", msg.Payload)
		// Display data on the frontend or CLI as needed
	}
}

func main() {
	// Set up HTTP endpoint
	http.HandleFunc("/sendData", handleRequest)
	go http.ListenAndServe(":8080", nil)

	// Start G2 to receive from Redis
	go receiveFromRedis()

	// Keep the main goroutine running
	select {}
}
