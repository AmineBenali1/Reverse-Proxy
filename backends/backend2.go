package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(9 * time.Second)
		fmt.Fprintln(w, "Hello from backend 8082!")
	})
	fmt.Println("Backend running on :8082")
	http.ListenAndServe(":8082", nil)
}
