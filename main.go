package main

import (
	"fmt"
	hd "gateway/internals/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", hd.MainHandler)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("WebSocket server starting on port http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
