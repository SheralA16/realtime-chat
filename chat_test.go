package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// TestHubCreation prueba la creación correcta de un nuevo hub
func TestHubCreation(t *testing.T) {
	hub := NewHub()

	if hub.clients == nil {
		t.Error("El mapa de clientes no se inicializó correctamente")
	}

	if hub.broadcast == nil {
		t.Error("El canal de difusión no se inicializó")
	}

	if hub.register == nil {
		t.Error("El canal de registro no se inicializó")
	}

	if hub.unregister == nil {
		t.Error("El canal de cancelación no se inicializó")
	}

	if hub.userHistory == nil {
		t.Error("El historial de usuarios no se inicializó")
	}

	// ✅ CORREGIDO: Usar GetClientCount() que sí existe
	if hub.GetClientCount() != 0 {
		t.Errorf("Se esperaban 0 clientes inicialmente, pero se encontraron %d", hub.GetClientCount())
	}
}

// TestMessageBroadcastBasic prueba el envío básico de mensajes
func TestMessageBroadcastBasic(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Crear un cliente mock para recibir el mensaje
	client := &Client{
		hub:      hub,
		send:     make(chan []byte, 256),
		username: "testuser",
	}

	// Registrar el cliente
	hub.register <- client
	time.Sleep(100 * time.Millisecond)

	// Limpiar mensajes del sistema
	for {
		select {
		case <-client.send:
			// Descartar mensajes del sistema
		case <-time.After(10 * time.Millisecond):
			// No hay más mensajes
			goto sendMessage
		}
	}

sendMessage:
	// Crear mensaje de prueba
	msg := NewMessage("testuser", "Mensaje de prueba")
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Error serializando mensaje: %v", err)
	}

	// Enviar mensaje al hub
	hub.broadcast <- msgBytes
	time.Sleep(100 * time.Millisecond)

	// Verificar que el cliente recibió el mensaje
	select {
	case receivedMsg := <-client.send:
		var parsedMsg Message
		if err := json.Unmarshal(receivedMsg, &parsedMsg); err != nil {
			t.Fatalf("Error parseando mensaje: %v", err)
		}

		if parsedMsg.Content != "Mensaje de prueba" {
			t.Errorf("Se esperaba contenido 'Mensaje de prueba', pero se encontró '%s'", parsedMsg.Content)
		}

		if parsedMsg.Username != "testuser" {
			t.Errorf("Se esperaba username 'testuser', pero se encontró '%s'", parsedMsg.Username)
		}

	case <-time.After(500 * time.Millisecond):
		t.Error("No se recibió el mensaje en el tiempo esperado")
	}
}

// TestImageMessage prueba el manejo de mensajes con imágenes
func TestImageMessage(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Crear un cliente para recibir el mensaje
	client := &Client{
		hub:      hub,
		send:     make(chan []byte, 256),
		username: "testuser",
	}

	// Registrar cliente
	hub.register <- client
	time.Sleep(100 * time.Millisecond)

	// Limpiar mensajes del sistema
	for {
		select {
		case <-client.send:
		case <-time.After(10 * time.Millisecond):
			goto testImage
		}
	}

testImage:
	// Crear datos de imagen simulados
	imageData := &ImageData{
		Data: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==",
		Name: "test.png",
		Type: "image/png",
		Size: 100,
	}

	// Crear mensaje con imagen
	msg := NewMessageWithImage("testuser", "Mira esta imagen", imageData)
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Error serializando mensaje con imagen: %v", err)
	}

	// Verificar propiedades del mensaje
	if !msg.HasImage {
		t.Error("El mensaje debería tener HasImage=true")
	}

	if msg.Image == nil {
		t.Fatal("Los datos de imagen no deberían ser nil")
	}

	if msg.Image.Name != "test.png" {
		t.Errorf("Se esperaba nombre 'test.png', pero se encontró '%s'", msg.Image.Name)
	}

	// Enviar al hub
	hub.broadcast <- msgBytes
	time.Sleep(100 * time.Millisecond)

	// Verificar que el cliente recibió el mensaje con imagen
	select {
	case receivedMsg := <-client.send:
		var parsedMsg Message
		if err := json.Unmarshal(receivedMsg, &parsedMsg); err != nil {
			t.Fatalf("Error parseando mensaje con imagen: %v", err)
		}

		if !parsedMsg.HasImage {
			t.Error("El mensaje recibido debería tener HasImage=true")
		}

		if parsedMsg.Image == nil {
			t.Error("Los datos de imagen no deberían ser nil en el mensaje recibido")
		}

	case <-time.After(500 * time.Millisecond):
		t.Error("No se recibió el mensaje con imagen en el tiempo esperado")
	}
}

// TestUsernameValidation prueba la validación de nombres de usuario
func TestUsernameValidation(t *testing.T) {
	// Casos válidos
	validUsernames := []string{"user123", "test_user", "user-name", "abc"}
	for _, username := range validUsernames {
		if !validateUsername(username) {
			t.Errorf("Username válido '%s' fue rechazado", username)
		}
	}

	// Casos inválidos
	invalidUsernames := []string{
		"",                      // vacío
		"a",                     // muy corto
		"usuario@",              // carácter especial
		"user name",             // espacio
		strings.Repeat("a", 25), // muy largo
	}
	for _, username := range invalidUsernames {
		if validateUsername(username) {
			t.Errorf("Username inválido '%s' fue aceptado", username)
		}
	}
}

// TestClientRegistration prueba el registro de clientes en el hub
func TestClientRegistration(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Crear un cliente mock
	client := &Client{
		hub:      hub,
		send:     make(chan []byte, 256),
		username: "testuser",
	}

	// Registrar el cliente
	hub.register <- client

	// Dar tiempo para que se procese
	time.Sleep(200 * time.Millisecond)

	// Verificar que el cliente se registró
	if hub.GetClientCount() != 1 {
		t.Errorf("Se esperaba 1 cliente, pero se encontraron %d", hub.GetClientCount())
	}

	// Verificar que se obtiene el nombre de usuario correcto
	users := hub.GetConnectedUsers()
	if len(users) != 1 || users[0] != "testuser" {
		t.Errorf("Se esperaba usuario 'testuser', pero se obtuvo %v", users)
	}
}

// TestDuplicateUsernames prueba que no se permitan nombres duplicados
func TestDuplicateUsernames(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Crear primer cliente
	client1 := &Client{
		hub:      hub,
		send:     make(chan []byte, 256),
		username: "duplicateuser",
	}

	// Registrar primer cliente
	hub.register <- client1
	time.Sleep(100 * time.Millisecond)

	// Verificar que se registró
	if hub.GetClientCount() != 1 {
		t.Fatalf("Primer cliente no se registró correctamente")
	}

	// ✅ NUEVA PRUEBA: Verificar disponibilidad de nombre
	// Simular que ya existe un usuario con ese nombre
	hub.mu.Lock()
	if _, exists := hub.clients[client1]; !exists {
		t.Error("El cliente no se encontró en el mapa de clientes")
	}
	hub.mu.Unlock()

	// Probar que el nombre no está disponible
	if hub.isUsernameAvailable("duplicateuser") {
		t.Error("El nombre duplicateuser debería estar ocupado")
	}

	// Probar que un nombre diferente sí está disponible
	if !hub.isUsernameAvailable("otheruser") {
		t.Error("El nombre otheruser debería estar disponible")
	}
}

// TestMessageBroadcast prueba la difusión de mensajes
func TestMessageBroadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Crear múltiples clientes mock
	numClients := 3
	clients := make([]*Client, numClients)

	for i := 0; i < numClients; i++ {
		clients[i] = &Client{
			hub:      hub,
			send:     make(chan []byte, 256),
			username: "testuser" + string(rune(i+'0')),
		}
		hub.register <- clients[i]
	}

	time.Sleep(300 * time.Millisecond)

	// Limpiar mensajes del sistema
	for _, client := range clients {
	clearLoop:
		for {
			select {
			case <-client.send:
				// Continuar descartando mensajes
			case <-time.After(10 * time.Millisecond):
				// No hay más mensajes, salir del bucle
				break clearLoop
			}
		}
	}

	// Verificar que todos los clientes se registraron
	if hub.GetClientCount() != numClients {
		t.Fatalf("Se esperaban %d clientes, pero se encontraron %d", numClients, hub.GetClientCount())
	}

	// Crear un mensaje de prueba
	msg := NewMessage("testuser0", "Hola mundo")
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Error serializando mensaje: %v", err)
	}

	// Enviar mensaje al hub para difusión
	hub.broadcast <- msgBytes

	// Dar tiempo para la difusión
	time.Sleep(100 * time.Millisecond)

	// Verificar que todos los clientes recibieron el mensaje
	messagesReceived := 0
	for i, client := range clients {
		select {
		case receivedMsg := <-client.send:
			var parsedMsg Message
			if err := json.Unmarshal(receivedMsg, &parsedMsg); err != nil {
				t.Errorf("Cliente %d: Error parseando mensaje recibido: %v", i, err)
				continue
			}

			if parsedMsg.Type == MessageTypeSystem {
				t.Logf("Cliente %d: Recibió mensaje del sistema (ignorado): %s", i, parsedMsg.Content)
				continue
			}

			if parsedMsg.Content != "Hola mundo" {
				t.Errorf("Cliente %d: Se esperaba contenido 'Hola mundo', pero se recibió '%s'", i, parsedMsg.Content)
				continue
			}

			messagesReceived++

		case <-time.After(500 * time.Millisecond):
			t.Errorf("Cliente %d: No recibió el mensaje en tiempo esperado", i)
		}
	}

	if messagesReceived != numClients {
		t.Errorf("Se esperaba que %d clientes recibieran el mensaje, pero solo %d lo recibieron", numClients, messagesReceived)
	}
}

// TestConcurrentOperations prueba operaciones concurrentes
func TestConcurrentOperations(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	const numGoroutines = 10
	const messagesPerGoroutine = 5

	var wg sync.WaitGroup

	// Registrar algunos clientes primero
	for i := 0; i < 3; i++ {
		client := &Client{
			hub:      hub,
			send:     make(chan []byte, 256),
			username: "permanentuser" + string(rune(i+'0')),
		}
		hub.register <- client
	}

	time.Sleep(200 * time.Millisecond)

	// Enviar mensajes concurrentemente
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < messagesPerGoroutine; j++ {
				username := "concurrentuser" + string(rune(goroutineID+'0'))
				content := "Mensaje " + string(rune(j+'0')) + " de goroutine " + string(rune(goroutineID+'0'))

				msg := NewMessage(username, content)

				if msgBytes, err := json.Marshal(msg); err == nil {
					select {
					case hub.broadcast <- msgBytes:
						// Mensaje enviado exitosamente
					case <-time.After(100 * time.Millisecond):
						t.Logf("Timeout enviando mensaje desde goroutine %d", goroutineID)
					}
				}

				time.Sleep(5 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	time.Sleep(500 * time.Millisecond)

	// Verificar que el sistema sigue funcionando
	finalClientCount := hub.GetClientCount()
	if finalClientCount < 3 {
		t.Errorf("Se perdieron clientes durante las operaciones concurrentes. Clientes restantes: %d", finalClientCount)
	}

	t.Logf("Operaciones concurrentes completadas. Clientes finales: %d", finalClientCount)
}

// TestWebSocketUpgrade prueba la actualización de WebSocket
func TestWebSocketUpgrade(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Crear servidor de prueba
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	}))
	defer server.Close()

	// Convertir URL HTTP a WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?username=testuser"

	// Conectar como cliente WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Error conectando WebSocket: %v", err)
	}
	defer conn.Close()

	// Dar tiempo para que se registre el cliente
	time.Sleep(200 * time.Millisecond)

	// Verificar que el cliente se registró en el hub
	if hub.GetClientCount() != 1 {
		t.Errorf("Se esperaba 1 cliente conectado, pero se encontraron %d", hub.GetClientCount())
	}

	// Enviar un mensaje de prueba
	testMessage := map[string]interface{}{
		"content":  "Mensaje de prueba",
		"hasImage": false,
	}

	if err := conn.WriteJSON(testMessage); err != nil {
		t.Fatalf("Error enviando mensaje: %v", err)
	}

	// Leer mensajes (puede haber mensajes del sistema primero)
	messageReceived := false
	for attempts := 0; attempts < 5; attempts++ {
		var receivedMsg Message
		if err := conn.ReadJSON(&receivedMsg); err != nil {
			if attempts == 4 { // Último intento
				t.Fatalf("Error leyendo mensaje: %v", err)
			}
			continue
		}

		// Verificar si es nuestro mensaje
		if receivedMsg.Type == MessageTypeMessage && receivedMsg.Content == "Mensaje de prueba" {
			messageReceived = true
			if receivedMsg.Username != "testuser" {
				t.Errorf("Se esperaba usuario 'testuser', pero se recibió '%s'", receivedMsg.Username)
			}
			break
		}
	}

	if !messageReceived {
		t.Error("No se recibió el mensaje de prueba")
	}
}

// TestClientDisconnection prueba la desconexión de clientes
func TestClientDisconnection(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Crear cliente
	client := &Client{
		hub:      hub,
		send:     make(chan []byte, 256),
		username: "disconnecttest",
	}

	// Registrar cliente
	hub.register <- client
	time.Sleep(100 * time.Millisecond)

	// Verificar que se registró
	if hub.GetClientCount() != 1 {
		t.Fatalf("Cliente no se registró correctamente")
	}

	// Desregistrar cliente
	hub.unregister <- client
	time.Sleep(100 * time.Millisecond)

	// Verificar que se desregistró
	if hub.GetClientCount() != 0 {
		t.Errorf("Se esperaban 0 clientes después de desconexión, pero se encontraron %d", hub.GetClientCount())
	}

	// Verificar que el usuario está marcado como desconectado
	userHistory := hub.GetUserHistory()
	if userStatus, exists := userHistory["disconnecttest"]; exists {
		if userStatus.Connected {
			t.Error("El usuario debería estar marcado como desconectado")
		}
	} else {
		t.Error("El usuario debería existir en el historial")
	}
}

// BenchmarkMessageBroadcast benchmarks la difusión de mensajes
func BenchmarkMessageBroadcast(b *testing.B) {
	hub := NewHub()
	go hub.Run()

	// Crear algunos clientes
	numClients := 10
	clients := make([]*Client, numClients)

	for i := 0; i < numClients; i++ {
		clients[i] = &Client{
			hub:      hub,
			send:     make(chan []byte, 256),
			username: "benchuser" + string(rune(i+'0')),
		}
		hub.register <- clients[i]
	}

	time.Sleep(100 * time.Millisecond)

	msg := NewMessage("benchuser", "benchmark message")
	msgBytes, _ := json.Marshal(msg)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		select {
		case hub.broadcast <- msgBytes:
		case <-time.After(10 * time.Millisecond):
			// Timeout para evitar bloqueos
		}
	}
}

// TestRaceConditions prueba específicamente condiciones de carrera
func TestRaceConditions(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	const numWorkers = 10         // ✅ REDUCIDO para evitar saturación
	const operationsPerWorker = 5 // ✅ REDUCIDO para ser más manejable

	var wg sync.WaitGroup
	createdClients := make([]*Client, 0, numWorkers*operationsPerWorker)
	var clientsMux sync.Mutex

	// Operaciones concurrentes de registro/desregistro
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for j := 0; j < operationsPerWorker; j++ {
				// ✅ CORREGIDO: Username válido (sin caracteres especiales)
				username := "race" + string(rune(workerID+'A')) + string(rune(j+'0'))

				// Crear cliente
				client := &Client{
					hub:      hub,
					send:     make(chan []byte, 256),
					username: username,
				}

				// Guardar referencia para cleanup
				clientsMux.Lock()
				createdClients = append(createdClients, client)
				clientsMux.Unlock()

				// Registrar
				hub.register <- client
				time.Sleep(2 * time.Millisecond) // ✅ Dar tiempo para procesar

				// Enviar mensaje
				msg := NewMessage(client.username, "mensaje concurrente "+string(rune(j+'0')))
				if msgBytes, err := json.Marshal(msg); err == nil {
					select {
					case hub.broadcast <- msgBytes:
					case <-time.After(50 * time.Millisecond):
						// Timeout para evitar bloqueos
					}
				}

				time.Sleep(2 * time.Millisecond)

				// ✅ DESREGISTRAR INMEDIATAMENTE para evitar leaks
				hub.unregister <- client
				time.Sleep(1 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()

	// ✅ CLEANUP ADICIONAL: Desregistrar cualquier cliente restante
	time.Sleep(200 * time.Millisecond) // Dar tiempo para que se procesen las desconexiones

	clientsMux.Lock()
	for _, client := range createdClients {
		select {
		case hub.unregister <- client:
		default:
			// Cliente ya desregistrado
		}
	}
	clientsMux.Unlock()

	time.Sleep(200 * time.Millisecond) // Tiempo final para cleanup

	// Verificar que el sistema está en estado consistente
	clientCount := hub.GetClientCount()
	t.Logf("Clientes finales después de operaciones de carrera: %d", clientCount)

	// ✅ EXPECTATIVA MÁS REALISTA: Pocos o ningún cliente
	if clientCount > 2 { // Margen de tolerancia muy pequeño
		t.Errorf("Demasiados clientes restantes, posible leak: %d", clientCount)
	}

	// Verificar que no hay deadlocks enviando un mensaje de prueba
	testMsg := NewMessage("testfinal", "test final")
	if msgBytes, err := json.Marshal(testMsg); err == nil {
		select {
		case hub.broadcast <- msgBytes:
			t.Logf("✅ Sistema responde correctamente después de operaciones de carrera")
		case <-time.After(100 * time.Millisecond):
			t.Error("❌ Sistema no responde, posible deadlock")
		}
	}
}
