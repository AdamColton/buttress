package gorillaserver

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type Server interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	Handler(r *http.Request) (http.Handler, string)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	HandleWebsocket(pattern string, handler SocketHandler)
	ListenAndServe(addr string)
	ServeFile(pattern, path string)
	ServeDirectory(pattern, path string)
	Upgrade(http.ResponseWriter, *http.Request, http.Header) (*websocket.Conn, error)
}

type SocketHandler func(*http.Request, *websocket.Conn)
