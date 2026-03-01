package main

import (
	"fmt"
	"net/http"

)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from backend 8083!")
	})
	fmt.Println("Backend running on :8083")
	http.ListenAndServe(":8083", nil)
}
