package main

import (
	"fmt"
	hd "gateway/internals/handlers"
	mw "gateway/internals/middleware"
	"net/http"
)

func main() {
	http.HandleFunc("/register", mw.CorsMiddleware(hd.RegisterHandler))
	http.HandleFunc("/login", mw.CorsMiddleware(hd.LoginHandler))
	http.HandleFunc("/ws", hd.HandlerWS)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("WebSocket server starting on port http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
