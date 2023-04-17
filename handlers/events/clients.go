package events

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type client struct {
	Id          string
	sendMessage chan EventMessage
}

func newClient(id string) *client {
	return &client{Id: id, sendMessage: make(chan EventMessage)}
}

func (c *client) onLine(ctx context.Context, w io.Writer, flusher http.Flusher) {
	for {
		select {
		case m := <-c.sendMessage:
			data, err := json.Marshal(m.Data)
			if err != nil {
				log.Println(err)
			}
			const format = "event:%s\ndata:%s\n\n"
			fmt.Fprintf(w, format, m.EventName, string(data))
			flusher.Flush()
			return

		case <-ctx.Done():
			return
		}
	}
}
