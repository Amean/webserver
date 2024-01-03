package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

func G1(redisClient *redis.Client) {
	for i := 1; i <= 10; i++ {

		strInt := strconv.Itoa(i)

		cmd := exec.Command("python3", "square.py", strInt)
		err := cmd.Run()

		if err != nil {
			log.Println("Error executing Python script:", err)
			continue
		}
	}
}

func G2(redisClient *redis.Client) {
	ctx := context.Background()
	pubsub := redisClient.Subscribe(ctx, "square")
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		square, _ := strconv.Atoi(msg.Payload)
		fmt.Printf("Received square from Redis: %d\n", square)
	}
}

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	go G2(redisClient)

	http.HandleFunc("/startG1", func(w http.ResponseWriter, r *http.Request) {
		go G1(redisClient)
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
