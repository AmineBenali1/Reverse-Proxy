package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		fmt.Fprintln(w, "Hello from backend 8084!")
	})
	fmt.Println("Backend running on :8084")
	http.ListenAndServe(":8084", nil)
}
