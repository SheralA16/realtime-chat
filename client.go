package main

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Tiempo máximo de espera para escribir un mensaje al cliente
	writeWait = 10 * time.Second

	// Tiempo máximo de espera para leer el siguiente pong del cliente
	pongWait = 60 * time.Second

	// Enviar pings al cliente con este período. Debe ser menor que pongWait
	pingPeriod = (pongWait * 9) / 10

	// Tamaño máximo del mensaje permitido del cliente
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client representa un cliente WebSocket activo
type Client struct {
	// El hub de chat al que pertenece este cliente
	hub *Hub

	// La conexión WebSocket
	conn *websocket.Conn

	// Canal con buffer para mensajes salientes
	send chan []byte

	// Nombre de usuario del cliente
	username string
}

// readPump bombea mensajes desde la conexión WebSocket al hub
//
// La aplicación ejecuta readPump en una goroutine per-conexión. La aplicación
// asegura que hay como máximo un lector en una conexión ejecutando todos los
// reads desde esta goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)

	// ✅ Manejar error de SetReadDeadline
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("Error estableciendo deadline de lectura para '%s': %v", c.username, err)
		return
	}

	c.conn.SetPongHandler(func(string) error {
		// ✅ Manejar error de SetReadDeadline en pong handler
		if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			log.Printf("Error estableciendo deadline en pong handler para '%s': %v", c.username, err)
		}
		return nil
	})

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error inesperado de WebSocket para '%s': %v", c.username, err)
			} else {
				log.Printf("Cliente '%s' cerró conexión: %v", c.username, err)
			}
			break
		}

		messageBytes = bytes.TrimSpace(bytes.Replace(messageBytes, newline, space, -1))

		// Intentar parsear el mensaje como JSON
		var incomingMsg struct {
			Content string `json:"content"`
		}

		if err := json.Unmarshal(messageBytes, &incomingMsg); err != nil {
			log.Printf("Error parseando mensaje JSON de cliente '%s': %v", c.username, err)
			continue
		}

		// Crear mensaje completo con metadata
		msg := NewMessage(c.username, incomingMsg.Content)

		// Serializar mensaje completo
		messageJSON, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error serializando mensaje de '%s': %v", c.username, err)
			continue
		}

		// Enviar al hub para difusión
		select {
		case c.hub.broadcast <- messageJSON:
			log.Printf("💬 Mensaje de '%s' enviado al hub", c.username)
		default:
			log.Printf("⚠️ Hub ocupado, mensaje de '%s' descartado", c.username)
		}
	}
}

// writePump bombea mensajes desde el hub a la conexión WebSocket
//
// Una goroutine que ejecuta writePump se inicia para cada conexión. La aplicación
// asegura que hay como máximo un writer en una conexión ejecutando todos los
// writes desde esta goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			// ✅ Manejar error de SetWriteDeadline
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("Error estableciendo deadline de escritura para '%s': %v", c.username, err)
				return
			}

			if !ok {
				// El hub cerró el canal
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Printf("Error enviando mensaje de cierre para '%s': %v", c.username, err)
				}
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("Error obteniendo writer para '%s': %v", c.username, err)
				return
			}

			// ✅ Manejar error de Write
			if _, err := w.Write(message); err != nil {
				log.Printf("Error escribiendo mensaje para '%s': %v", c.username, err)
				w.Close()
				return
			}

			// Agregar mensajes de chat en cola al mensaje actual
			n := len(c.send)
			for i := 0; i < n; i++ {
				if _, err := w.Write(newline); err != nil {
					log.Printf("Error escribiendo newline para '%s': %v", c.username, err)
					w.Close()
					return
				}

				queuedMessage := <-c.send
				if _, err := w.Write(queuedMessage); err != nil {
					log.Printf("Error escribiendo mensaje en cola para '%s': %v", c.username, err)
					w.Close()
					return
				}
			}

			// ✅ Manejar error de Close
			if err := w.Close(); err != nil {
				log.Printf("Error cerrando writer para '%s': %v", c.username, err)
				return
			}

		case <-ticker.C:
			// ✅ Manejar error de SetWriteDeadline para ping
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("Error estableciendo deadline para ping para '%s': %v", c.username, err)
				return
			}

			// ✅ Manejar error de WriteMessage para ping
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error enviando ping para '%s': %v", c.username, err)
				return
			}
		}
	}
}
