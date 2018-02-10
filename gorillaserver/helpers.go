package gorillaserver

import (
	"github.com/gorilla/websocket"
	"net/http"
)

// ReadSocket will listen on a websocket and transmit any data it receives over
// a channel
func ReadSocket(socket *websocket.Conn, bufferLength int) <-chan []byte {
	msgChan := make(chan []byte, bufferLength)
	go func(sendMsgChan chan<- []byte) {
		for {
			_, msg, err := socket.ReadMessage()
			if err != nil {
				//lost connection
				close(sendMsgChan)
				break
			}
			sendMsgChan <- msg
		}
	}(msgChan)
	return msgChan
}

func RedirectFunc(path string, code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, code)
	}
}
