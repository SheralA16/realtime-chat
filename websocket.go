package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// upgrader especifica parámetros para actualizar una conexión HTTP a WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin devuelve true si el request debe ser permitido a proceder
	CheckOrigin: func(r *http.Request) bool {
		// Permitir conexiones desde cualquier origen
		// En producción, deberías validar el origen apropiadamente
		return true
	},
}

// serveWS maneja las solicitudes WebSocket del cliente
func serveWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error actualizando conexión a WebSocket: %v", err)
		return
	}

	// Obtener nombre de usuario de los parámetros de consulta
	username := r.URL.Query().Get("username")
	if username == "" {
		// Generar nombre de usuario por defecto si no se proporciona
		username = "Usuario" + time.Now().Format("150405")
	}

	// Crear nuevo cliente
	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256), // Canal con buffer para evitar bloqueos
		username: username,
	}

	// Registrar cliente en el hub
	client.hub.register <- client

	// Permitir recolección de memoria de referencia al cliente
	// ejecutando todas las goroutines
	go client.writePump()
	go client.readPump()
}
