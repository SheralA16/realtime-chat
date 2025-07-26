package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

// UserStatus representa el estado de un usuario
type UserStatus struct {
	Username    string    `json:"username"`
	Connected   bool      `json:"connected"`
	LastSeen    time.Time `json:"lastSeen"`
	ConnectedAt time.Time `json:"connectedAt"`
}

// Hub mantiene el conjunto de clientes activos y difunde mensajes a los clientes
type Hub struct {
	// Clientes registrados - mapa protegido por mutex
	clients map[*Client]bool

	// Historial de todos los usuarios que se han conectado
	userHistory map[string]*UserStatus

	//  AGREGAR ESTAS DOS LÍNEAS:
	messageHistory []*Message
	maxHistorySize int

	// Mensajes entrantes de los clientes para difundir
	broadcast chan []byte

	// Solicitudes de registro de nuevos clientes
	register chan *Client

	// Solicitudes de cancelación de registro de clientes
	unregister chan *Client

	// Mutex para proteger acceso concurrente al mapa de clientes y historial
	mu sync.RWMutex
}

// NewHub crea una nueva instancia del hub de chat
func NewHub() *Hub {
	return &Hub{
		broadcast:   make(chan []byte, 1000), // Buffer para evitar bloqueos
		register:    make(chan *Client, 100),
		unregister:  make(chan *Client, 100),
		clients:     make(map[*Client]bool),
		userHistory: make(map[string]*UserStatus),
		// ⭐ AGREGAR ESTAS LÍNEAS:
		messageHistory: make([]*Message, 0),
		maxHistorySize: 50,
	}
}

// Run inicia el loop principal del hub
func (h *Hub) Run() {
	log.Println("🚀 Hub iniciado, esperando conexiones...")

	for {
		select {
		case client := <-h.register:
			h.registerClient(client) //nuevo usuario

		case client := <-h.unregister:
			h.unregisterClient(client) //usuario desconectado

		case message := <-h.broadcast:
			h.broadcastMessage(message) //mensaje entrante
		}
	}
}

// isUsernameAvailable verifica si un nombre de usuario está disponible (método privado)
func (h *Hub) isUsernameAvailable(username string) bool {
	h.mu.RLock()         // Bloquear lectura del mapa de clientes
	defer h.mu.RUnlock() // Desbloquear al final de la función

	// Verificar si hay algún cliente conectado con ese nombre EXACTO
	for client := range h.clients {
		if client.username == username {
			return false
		}
	}

	// También verificar en el historial si está conectado actualmente
	if userStatus, exists := h.userHistory[username]; exists && userStatus.Connected {
		return false
	}

	return true
}

// registerClient registra un nuevo cliente en el hub
func (h *Hub) registerClient(client *Client) {
	// ⭐ VALIDACIÓN: Verificar si el nombre de usuario ya está en uso
	if !h.isUsernameAvailable(client.username) {
		log.Printf("❌ Intento de conexión con nombre duplicado: '%s'", client.username)

		// Enviar mensaje de error al cliente
		errorMsg := map[string]interface{}{
			"type":    "error",
			"message": "El nombre de usuario '" + client.username + "' ya está en uso. Por favor, elige otro nombre.",
			"code":    "USERNAME_TAKEN",
		}

		if msgBytes, err := json.Marshal(errorMsg); err == nil {
			select {
			case client.send <- msgBytes:
				log.Printf("📤 Mensaje de error enviado a cliente con nombre duplicado")
			default:
				log.Printf("❌ No se pudo enviar mensaje de error al cliente")
			}
		}

		// Cerrar la conexión después de un breve delay para que el mensaje llegue
		go func() {
			time.Sleep(100 * time.Millisecond)
			client.conn.Close()
		}()

		return // Salir si el nombre ya está en uso
	}

	// Si llegamos aquí, el nombre está disponible
	h.mu.Lock()              // Bloquear escritura del mapa de clientes
	h.clients[client] = true // Registrar cliente en el mapa

	// Actualizar o crear estado del usuario
	now := time.Now()
	if userStatus, exists := h.userHistory[client.username]; exists {
		userStatus.Connected = true
		userStatus.ConnectedAt = now
		userStatus.LastSeen = now
	} else {
		h.userHistory[client.username] = &UserStatus{
			Username:    client.username,
			Connected:   true,
			ConnectedAt: now,
			LastSeen:    now,
		}
	}

	clientCount := len(h.clients)
	h.mu.Unlock()

	log.Printf("✅ Cliente '%s' conectado exitosamente. Total de clientes: %d", client.username, clientCount)

	successMsg := map[string]interface{}{ // Mensaje de éxito
		"type":     "connectionSuccess",
		"message":  "Conectado exitosamente como " + client.username,
		"username": client.username,
	}

	if msgBytes, err := json.Marshal(successMsg); err == nil {
		select {
		case client.send <- msgBytes:
		default:
		}
	}

	h.broadcastUserList() // Enviar lista de usuarios actualizada

	joinMsg := NewSystemMessage(client.username + " se ha unido al chat") // mensaje de sistema
	joinMsg.Type = MessageTypeJoin

	if msgBytes, err := json.Marshal(joinMsg); err == nil {
		h.broadcastMessage(msgBytes)
	} else {
		log.Printf("Error serializando mensaje de conexión: %v", err)
	}
}

// unregisterClient cancela el registro de un cliente del hub
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	if _, ok := h.clients[client]; ok {
		// Eliminar cliente del mapa y cerrar su canal de envío
		delete(h.clients, client)
		close(client.send)

		// Actualizar estado del usuario a desconectado
		if userStatus, exists := h.userHistory[client.username]; exists {
			userStatus.Connected = false
			userStatus.LastSeen = time.Now()
		}

		clientCount := len(h.clients)
		h.mu.Unlock()

		log.Printf("🔌 Cliente '%s' desconectado. Total de clientes: %d", client.username, clientCount)

		// Enviar lista de usuarios actualizada
		h.broadcastUserList()

		// Enviar mensaje de sistema
		leaveMsg := NewSystemMessage(client.username + " ha salido del chat")
		leaveMsg.Type = MessageTypeLeave

		if msgBytes, err := json.Marshal(leaveMsg); err == nil {
			h.broadcastMessage(msgBytes)
		} else {
			log.Printf("Error serializando mensaje de desconexión: %v", err)
		}
	} else {
		h.mu.Unlock()
	}
}

// broadcastMessage envía un mensaje a todos los clientes conectados
func (h *Hub) broadcastMessage(message []byte) {
	h.mu.RLock()
	clients := make([]*Client, 0, len(h.clients))
	for client := range h.clients {
		clients = append(clients, client)
	}
	h.mu.RUnlock()

	// Debug: mostrar qué se está enviando
	log.Printf("Enviando mensaje a %d clientes", len(clients))

	// Enviar mensaje a cada cliente
	for _, client := range clients {
		select {
		case client.send <- message:
			// Mensaje enviado exitosamente
		default:
			// El canal del cliente está lleno o cerrado
			h.mu.Lock()
			delete(h.clients, client)
			h.mu.Unlock()
			close(client.send)
			log.Printf("Cliente '%s' eliminado por canal bloqueado", client.username)
		}
	}
}

// broadcastUserList envía la lista actualizada de usuarios a todos los clientes
func (h *Hub) broadcastUserList() {
	h.mu.RLock()
	users := make([]*UserStatus, 0, len(h.userHistory))
	for _, userStatus := range h.userHistory {
		// Crear copia para evitar problemas de concurrencia
		userCopy := &UserStatus{
			Username:    userStatus.Username,
			Connected:   userStatus.Connected,
			LastSeen:    userStatus.LastSeen,
			ConnectedAt: userStatus.ConnectedAt,
		}
		users = append(users, userCopy)
	}
	h.mu.RUnlock()

	// Crear mensaje con la lista de usuarios
	userListMsg := map[string]interface{}{
		"type":  "userList",
		"users": users,
	}

	log.Printf("Enviando lista de %d usuarios a todos los clientes", len(users))

	if msgBytes, err := json.Marshal(userListMsg); err == nil {
		h.broadcastMessage(msgBytes)
	} else {
		log.Printf("Error serializando lista de usuarios: %v", err)
	}
}

// GetClientCount devuelve el número actual de clientes conectados de forma thread-safe
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	count := len(h.clients)
	h.mu.RUnlock()
	return count
}

// GetConnectedUsers devuelve una lista de nombres de usuarios conectados
func (h *Hub) GetConnectedUsers() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]string, 0, len(h.clients))
	for client := range h.clients {
		users = append(users, client.username)
	}
	return users
}

// GetUserHistory devuelve el historial de todos los usuarios
func (h *Hub) GetUserHistory() map[string]*UserStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Crear copia del mapa para evitar problemas de concurrencia
	history := make(map[string]*UserStatus)
	for username, status := range h.userHistory {
		statusCopy := *status // Copia el valor
		history[username] = &statusCopy
	}
	return history
}

// GetMessageHistory devuelve el historial de mensajes (para debugging y tests)
func (h *Hub) GetMessageHistory() []*Message {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Crear copia del slice para evitar problemas de concurrencia
	history := make([]*Message, len(h.messageHistory))
	copy(history, h.messageHistory)
	return history
}
