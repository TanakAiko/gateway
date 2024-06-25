package main

import (
	"fmt"
	conf "gateway/config"
	hd "gateway/internals/handlers"
	ws "gateway/internals/handlers/ws"
	mw "gateway/internals/middleware"
	"net/http"
)

func main() {
	http.HandleFunc("/authorized", mw.CorsMiddleware(hd.AuthorizedHandler))
	http.HandleFunc("/register", mw.CorsMiddleware(hd.RegisterHandler))
	http.HandleFunc("/login", mw.CorsMiddleware(hd.LoginHandler))
	http.HandleFunc("/ws", ws.HandlerWS)

	go ws.HandleMessages()

	fmt.Printf("Gateway server starting on port http://localhost:%v\n", conf.Port)
	if err := http.ListenAndServe(":"+conf.Port, nil); err != nil {
		fmt.Println(err)
	}
}
