package gorillaserver

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// GorillaServer uses the Gorilla mux Router and websocket Upgrader
type GorillaServer struct {
	*mux.Router
	websocket.Upgrader
	CSRF struct {
		Auth []byte
		Opts []csrf.Option
	}
}

// New creates a new socket server
func New() *GorillaServer {
	return &GorillaServer{
		Router: mux.NewRouter(),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

// ListenAndServe causes the server to start, the method will not return as long
// as the server is running
func (server *GorillaServer) ListenAndServe(addr string) error {
	var r http.Handler = server.Router
	if server.CSRF.Auth != nil {
		r = csrf.Protect(server.CSRF.Auth, server.CSRF.Opts...)(r)
	}
	return http.ListenAndServe(addr, r)
}

// HandleWebsocket will register a SocketHandler to a url pattern
func (server *GorillaServer) HandleWebsocket(pattern string, handler SocketHandler) {
	server.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		socket, err := server.Upgrade(w, r, nil)
		if err != nil {
			log.Print("while_upgrading_socket", err)
			return
		}
		handler(r, socket)
	})
}

// ServeFile will serve a specific file or directory when the pattern is matched
// it will not serve files with in a directory
func ServeFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}

// ServeDirectory will serve a directory, including all sub files and
// directories. It will not allow parent directories to be accessed. Pattern
// should end with "/"
func (server *GorillaServer) ServeDirectory(pattern, path string) {
	server.PathPrefix(pattern).Handler(http.StripPrefix(pattern, http.FileServer(http.Dir(path))))
}

// ServeDirectory will serve a directory, including all sub files and
// directories. It will not allow parent directories to be accessed. Pattern
// should end with "/". This is useful is the files being served are not off
// the root of the site.
func (server *GorillaServer) ServeFile(pattern, path string) {
	server.HandleFunc(pattern, ServeFile(path))
}
