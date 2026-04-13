package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/checkout", func(w http.ResponseWriter, r *http.Request) {
		// Simulate a critical dependency (e.g., a Bank API)
		log.Println("Processing payment...")
		
		// Simulate latency or failure
		time.Sleep(100 * time.Millisecond) 
		
		fmt.Fprint(w, "Order #12345: Payment Successful via GopherShield Protection!")
	})

	log.Println("E-commerce Order Service running on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}