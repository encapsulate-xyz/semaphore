package sockets

import (
	"fmt"
	"github.com/semaphoreui/semaphore/db"
	"github.com/semaphoreui/semaphore/pkg/tz"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/websocket"
	"github.com/semaphoreui/semaphore/util"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 2 * 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 2 * 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type connection struct {
	ws     *websocket.Conn
	send   chan []byte
	userID int
}

func (c *connection) logError(err error, msg string) {
	log.WithError(err).WithFields(log.Fields{
		"context": "websocket",
		"user_id": c.userID,
	}).Error(msg)
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		_ = c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)

	if err := c.ws.SetReadDeadline(tz.Now().Add(pongWait)); err != nil {
		c.logError(err, "Cannot set read deadline")
	}

	c.ws.SetPongHandler(func(string) error {
		err2 := c.ws.SetReadDeadline(tz.Now().Add(pongWait))
		util.LogErrorF(err2, log.Fields{"error": "Cannot set read deadline"})
		return nil
	})

	for {
		_, message, err := c.ws.ReadMessage()
		fmt.Println(string(message))

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				c.logError(err, "Cannot read message from websocket")
			}
			break
		}
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {

	if err := c.ws.SetWriteDeadline(tz.Now().Add(writeWait)); err != nil {
		c.logError(err, "Cannot set write deadline")
	}

	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		_ = c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				if err := c.write(websocket.CloseMessage, []byte{}); err != nil {
					c.logError(err, "Cannot send close message")
				}
				return
			}

			if err := c.write(websocket.TextMessage, message); err != nil {
				c.logError(err, "Cannot send message")
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				c.logError(err, "Cannot send ping message")
				return
			}
		}
	}
}

// Handler is used by the router to handle the /ws endpoint
func Handler(w http.ResponseWriter, r *http.Request) {
	usr := context.Get(r, "user")
	if usr == nil {
		return
	}

	user := usr.(*db.User)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c := &connection{
		send:   make(chan []byte, 256),
		ws:     ws,
		userID: user.ID,
	}

	h.register <- c

	go c.writePump()
	c.readPump()
}

// Message allows a message to be sent to the websockets, called in API task logging
func Message(userID int, message []byte) {
	h.broadcast <- &sendRequest{
		userID: userID,
		msg:    message,
	}
}
