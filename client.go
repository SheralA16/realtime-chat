package main

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
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

	// Tamaño máximo del mensaje permitido del cliente (aumentado para imágenes)
	maxMessageSize = 10 * 1024 * 1024 // 10MB para soportar imágenes
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

// IncomingMessage representa un mensaje entrante del cliente
type IncomingMessage struct {
	Content  string     `json:"content"`
	HasImage bool       `json:"hasImage"`
	Image    *ImageData `json:"image,omitempty"`
}

// readPump bombea mensajes desde la conexión WebSocket al hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)

	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("Error estableciendo deadline de lectura para '%s': %v", c.username, err)
		return
	}

	c.conn.SetPongHandler(func(string) error {
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
		var incomingMsg IncomingMessage
		if err := json.Unmarshal(messageBytes, &incomingMsg); err != nil {
			log.Printf("Error parseando mensaje JSON de cliente '%s': %v", c.username, err)
			continue
		}

		// ⭐ VALIDACIONES DE SEGURIDAD PARA IMÁGENES
		if incomingMsg.HasImage && incomingMsg.Image != nil {
			// Validar que sea una imagen válida
			if !c.isValidImage(incomingMsg.Image) {

				c.sendErrorMessage("Imagen inválida. Solo se permiten imágenes de hasta 5MB.")
				continue
			}
			log.Printf("Es valida recibida de '%s': %s (%d bytes)",
				c.username, incomingMsg.Image.Name, incomingMsg.Image.Size)
		}

		// Validar contenido de texto si no hay imagen
		if !incomingMsg.HasImage && strings.TrimSpace(incomingMsg.Content) == "" {
			log.Printf("Esta vacío recibido de '%s'", c.username)
			continue
		}

		// Crear mensaje completo con metadata
		var msg *Message
		if incomingMsg.HasImage && incomingMsg.Image != nil {
			msg = NewMessageWithImage(c.username, incomingMsg.Content, incomingMsg.Image)
			log.Printf(" Con imagen de '%s': texto='%s', imagen='%s'",
				c.username, incomingMsg.Content, incomingMsg.Image.Name)
		} else {
			msg = NewMessage(c.username, incomingMsg.Content)

		}

		// Serializar mensaje completo
		messageJSON, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error serializando mensaje de '%s': %v", c.username, err)
			continue
		}

		// Enviar al hub para difusión
		select {
		case c.hub.broadcast <- messageJSON:
			log.Printf("Mensaje de '%s' enviado al hub para difusión", c.username)
		default:
			log.Printf("Hub ocupado, mensaje de '%s' descartado", c.username)
		}
	}
}

// isValidImage valida que los datos de imagen sean seguros
func (c *Client) isValidImage(image *ImageData) bool {
	// Validar tamaño máximo (5MB)
	const maxImageSize = 5 * 1024 * 1024
	if image.Size > maxImageSize {
		return false
	}

	// Validar que sea un tipo MIME de imagen válido
	validTypes := []string{
		"image/jpeg", "image/jpg", "image/png", "image/gif",
		"image/webp", "image/bmp", "image/svg+xml",
	}

	isValidType := false
	for _, validType := range validTypes {
		if image.Type == validType {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return false
	}

	// Validar que los datos estén en formato data URL válido
	if !strings.HasPrefix(image.Data, "data:") {
		return false
	}

	// Validar nombre de archivo (longitud y caracteres básicos)
	if len(image.Name) == 0 || len(image.Name) > 255 {
		return false
	}

	return true
}

// sendErrorMessage envía un mensaje de error al cliente
func (c *Client) sendErrorMessage(errorText string) {
	errorMsg := map[string]interface{}{
		"type":    "error",
		"message": errorText,
		"code":    "INVALID_IMAGE",
	}

	if msgBytes, err := json.Marshal(errorMsg); err == nil {
		select {
		case c.send <- msgBytes:
			log.Printf(" F bb, mensaje de error enviado a '%s'", c.username)
		default:
			log.Printf("No se pudo enviar mi neto, el mensaje de error a '%s'", c.username)
		}
	}
}

// writePump bombea mensajes desde el hub hacia la conexión WebSocket - enviando cada mensaje como un frame separado
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("El error estableciendo deadline de escritura para '%s': %v", c.username, err)
				return
			}

			if !ok {
				// El hub cerró el canal
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Printf("Enviando error, mensaje de cierre para '%s': %v", c.username, err)
				}
				return
			}

			// ⭐ ENVÍO OPTIMIZADO: Un mensaje por WebSocket frame
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WARNING XD, escribiendo mensaje para '%s': %v", c.username, err)
				return
			}

			// ⭐ PROCESAR MENSAJES ADICIONALES EN BUFFER SIN CONCATENAR
		additionalMessages:
			for {
				select {
				case nextMessage := <-c.send:
					// Enviar cada mensaje adicional como frame separado
					if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
						log.Printf("ERORRRRR, estableciendo deadline para mensaje adicional: %v", err)
						return
					}
					if err := c.conn.WriteMessage(websocket.TextMessage, nextMessage); err != nil {
						log.Printf("FFFFF, enviando mensaje adicional para '%s': %v", c.username, err)
						return
					}
				default:
					// No hay más mensajes en buffer
					break additionalMessages
				}
			}

		case <-ticker.C:
			// Enviar ping
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("FF, deadline para ping para '%s': %v", c.username, err)
				return
			}

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("WARNINGGG enviando ping para '%s': %v", c.username, err)
				return
			}
		}
	}
}
