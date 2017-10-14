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

/*

type SocketLoop struct {
	writeCh chan []byte
	readCh  chan []byte
	close   chan bool
}

type LoopCallback func(read <-chan []byte, send chan<- byte, closeLoop chan<- bool)

func LoopRunner(callback LoopCallback) SocketHandler {
	return func(req *http.Request, socket *websocket.Conn) {
		read := ReadSocket(socket, 3)
		send := make(chan []byte, 3)
		closeLoop := make(chan bool)
		go callback(read, send, closeLoop)
		for {
			select {
			case msg := <-send:
				socket.WriteMessage(1, msg)
			case <-closeLoop:
				socket.Close()
				return
			}
		}
	}
}
*/
