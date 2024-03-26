package webservice

import (
	"context"
	"log"
	"net/http"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func (o *Webservice) handlerWebSocket(w http.ResponseWriter, r *http.Request) {
	// Die URL wird gelesen

	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("WebSocket Accept Error:", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "Der interne Serverfehler ist aufgetreten")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		var v interface{}
		err := wsjson.Read(ctx, c, &v)
		if err != nil {
			log.Println("Lesen fehlgeschlagen:", err)
			break
		}

		log.Printf("Empfangen: %v", v)

		if err := wsjson.Write(ctx, c, v); err != nil {
			log.Println("Senden fehlgeschlagen:", err)
			break
		}
	}
	c.Close(websocket.StatusNormalClosure, "")
}
