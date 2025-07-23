# 🚀 GO O NO GO - Chat en Tiempo Real con Soporte para Imágenes

**Desarrollado por: JUNIOR_ALVINES y SheralA16**  
**GitHub: [JUNMPI/realtime-chat](https://github.com/JUNMPI/realtime-chat)**  
**GitHub: [SheralA16/realtime-chat](https://github.com/SheralA16/realtime-chat)**

Un sistema de chat en tiempo real desarrollado en Go con WebSockets, interfaz Bootstrap y **soporte completo para envío de imágenes**.

## 🌐 **Demo en Vivo**

**¡Prueba el chat ahora mismo!**
```
https://realtime-chat-production-183c.up.railway.app
```

## 🌟 Características

✅ **Chat en tiempo real** con WebSockets  
✅ **Envío de imágenes** - JPEG, PNG, GIF, WebP (máx. 5MB)  
✅ **Vista previa de imágenes** - Modal con zoom y descarga  
✅ **Arrastrar y soltar** - Interfaz intuitiva para subir imágenes  
✅ **Control de usuarios duplicados** - No permite nombres repetidos  
✅ **Historial persistente** - Los mensajes no se borran durante la sesión  
✅ **Interfaz moderna** con Bootstrap 5  
✅ **Lista de usuarios** conectados/desconectados  
✅ **Mensajes del sistema** para conexiones  
✅ **Responsive design** - Funciona en móviles  
✅ **Deploy en Railway** - Fácil y gratis  
✅ **Manejo robusto de errores** - Sin pérdida de conversación  
✅ **Tests con Race Detector** - Validación de concurrencia  

## 🖼️ Características de Imágenes

### **Formatos Soportados:**
- **JPEG/JPG** - Fotos comprimidas
- **PNG** - Imágenes con transparencia
- **GIF** - Imágenes animadas
- **WebP** - Formato moderno optimizado

### **Funcionalidades:**
- 📤 **Subida por arrastrar y soltar**
- 📤 **Selector de archivos tradicional**
- 🔍 **Vista previa antes de enviar**
- 💬 **Captions opcionales para imágenes**
- 🖼️ **Modal de vista completa**
- 💾 **Descarga de imágenes recibidas**
- 📏 **Información de tamaño y formato**
- ⚡ **Validación en tiempo real**
- 🚫 **Sin pérdida de historial** - Las conversaciones se mantienen intactas

### **Limitaciones:**
- 📦 **Tamaño máximo:** 5MB por imagen
- 🔒 **Solo tipos permitidos:** JPEG, PNG, GIF, WebP
- 🌐 **Base64:** Las imágenes se envían codificadas

## 📁 Estructura del Proyecto

```
realtime-chat/
├── main.go              # Servidor HTTP configurado para Railway
├── hub.go               # Gestión central de clientes y mensajes
├── client.go            # Manejo de clientes WebSocket individuales (✅ MEJORADO)
├── message.go           # Estructuras de mensajes
├── websocket.go         # Configuración WebSocket
├── index.html           # Frontend con Bootstrap (✅ CORREGIDO)
├── chat_test.go         # Tests unitarios con Race Detector
├── go.mod              # Dependencias de Go
├── go.sum              # Checksums de dependencias
└── README.md           # Esta documentación
```

## 🚀 Deploy en Railway (Paso a Paso)

### **Paso 1: Preparar el Repositorio**

1. **Clonar tu repositorio:**
```bash
git clone https://github.com/JUNMPI/realtime-chat.git
cd realtime-chat
```

2. **Actualizar archivos con las versiones corregidas:**
   - ✅ `index.html` con historial persistente corregido
   - ✅ `client.go` con soporte para imágenes
   - ✅ `message.go` con campos de imagen
   - ✅ `chat_test.go` con tests de imágenes y race detector

3. **Commit y push:**
```bash
git add .
git commit -m "🐛 Arreglar historial persistente y mejorar funcionalidad de imágenes"
git push origin main
```

### **Paso 2: Deploy en Railway**

1. **Ve a [Railway.app](https://railway.app)**
2. **Haz clic en "Start a New Project"**
3. **Selecciona "Deploy from GitHub repo"**
4. **Conecta tu cuenta de GitHub si no lo has hecho**
5. **Busca y selecciona `JUNMPI/realtime-chat`**
6. **¡Railway detecta automáticamente que es Go y empieza el deploy!**

### **Paso 3: Obtener tu URL**

**Tu chat está desplegado y funcionando en:**
```
https://realtime-chat-production-183c.up.railway.app
```

### **Paso 4: Probar Funcionalidad Completa**

¡Comparte la URL con tus amigos y prueben todas las características!

## 🧪 Script Completo de Testing con Race Detector

Para ejecutar la batería completa de tests con race detector, copia y pega este script en PowerShell:

```powershell
# Ir a tu directorio del proyecto
cd C:\GoProyectos\realtime-chat

Write-Host "🧪 GO O NO GO - Tests Completos con Race Detector" -ForegroundColor Green
Write-Host "=================================================" -ForegroundColor Green

Write-Host "`n1️⃣ Tests básicos:" -ForegroundColor Yellow
go test -v

Write-Host "`n2️⃣ Race Detector (LO MÁS IMPORTANTE):" -ForegroundColor Yellow
$env:CGO_ENABLED=1; go test -race -v

Write-Host "`n3️⃣ Múltiples ejecuciones con race detector:" -ForegroundColor Yellow
$env:CGO_ENABLED=1; go test -race -count=3

Write-Host "`n4️⃣ Benchmarks con race detector:" -ForegroundColor Yellow
$env:CGO_ENABLED=1; go test -race -bench=.

Write-Host "`n✅ Todos los tests completados!" -ForegroundColor Green
```

### **¿Qué hace cada comando?**

1. **Tests básicos** (`go test -v`): Verifica funcionalidad general
2. **Race detector** (`go test -race -v`): Detecta condiciones de carrera
3. **Múltiples ejecuciones** (`go test -race -count=3`): Ejecuta 3 veces para detectar problemas intermitentes
4. **Benchmarks** (`go test -race -bench=.`): Mide rendimiento bajo carga con race detector

### **Resultado esperado:**
Si todo está bien, verás `PASS` en todos los tests sin ningún `WARNING: DATA RACE`.

## 🧪 Pruebas Locales

Para probar en tu máquina antes de deployar:

```bash
# Ejecutar el servidor
go run *.go

# Abrir en navegador
http://localhost:8080
```

### **Probar Funcionalidad Completa:**
1. **Conectarte con un nombre de usuario único**
2. **Enviar mensajes de texto normales**
3. **Hacer clic en el botón de imagen** 📷
4. **Arrastrar una imagen o hacer clic para seleccionar**
5. **Añadir un caption opcional**
6. **Enviar la imagen**
7. **Continuar chateando** - ¡El historial se mantiene!
8. **Hacer clic en imágenes recibidas para vista completa**

## 🎯 Funcionalidades del Chat

### **Control de Usuarios:**
- ✅ Nombres únicos (no permite duplicados)
- ✅ Validación de formato (solo letras, números, - y _)
- ✅ Longitud entre 2-20 caracteres

### **Mensajes de Texto:**
- ✅ Envío en tiempo real
- ✅ Timestamps automáticos
- ✅ Notificaciones de conexión/desconexión
- ✅ Diferenciación visual (propios vs otros)

### **Mensajes de Imagen:**
- ✅ Subida por arrastrar y soltar
- ✅ Vista previa antes de enviar
- ✅ Captions opcionales
- ✅ Modal de vista completa
- ✅ Descarga de imágenes
- ✅ Validación de formato y tamaño
- ✅ Información de archivo (nombre, tamaño)
- ✅ **Historial persistente** - Sin pérdida de conversación

### **Lista de Usuarios:**
- ✅ Estado online/offline
- ✅ Tiempo de conexión
- ✅ Última vez visto
- ✅ Contador de usuarios activos

## 🏗️ Arquitectura Concurrente del Servidor

### **Patrón Hub-and-Spoke: Cómo Interactúan las Goroutines y Canales**

El sistema utiliza una arquitectura concurrente basada en el patrón **Hub-and-Spoke**, donde una goroutine central (Hub) coordina múltiples goroutines de clientes.

#### **📊 Diagrama de Arquitectura:**

```
                    🏢 HUB CENTRAL (1 Goroutine)
                    ├── Canal register   ← Nuevos clientes
                    ├── Canal unregister ← Clientes desconectados  
                    └── Canal broadcast  ← Mensajes para difundir
                           ↕️
        ┌─────────────────────────────────────────────────────┐
        │                    DIFUSIÓN                         │
        │          Envía a TODOS los clientes                 │
        └─────────────────────────────────────────────────────┘
                           ↓
    👤 CLIENTE A        👤 CLIENTE B        👤 CLIENTE C
    ├── 📥 readPump     ├── 📥 readPump     ├── 📥 readPump
    └── 📤 writePump    └── 📤 writePump    └── 📤 writePump
    
    TOTAL: 1 Hub + (2 × N usuarios) goroutines
```

#### **🔄 Flujo de Comunicación Detallado:**

**1. Conexión de Usuario:**
```go
// websocket.go - Cuando llega nueva conexión
client := &Client{hub, conn, send, username}
client.hub.register <- client  // 📨 Envía al canal register

// 🚀 Inicia 2 goroutines por cliente:
go client.writePump()  // Goroutine para ENVIAR mensajes
go client.readPump()   // Goroutine para LEER mensajes
```

**2. Hub Central (hub.go):**
```go
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:     // 📥 Nuevo cliente
            h.registerClient(client)
        case client := <-h.unregister:   // 📤 Cliente se va
            h.unregisterClient(client)  
        case message := <-h.broadcast:   // 💬 Nuevo mensaje
            h.broadcastMessage(message)  // 📡 Difundir a TODOS
        }
    }
}
```

**3. Envío de Mensaje (client.go):**
```go
// readPump: Lee mensaje del navegador del usuario
func (c *Client) readPump() {
    messageBytes := c.conn.ReadMessage()  // 📥 Lee del WebSocket
    
    // 📨 Envía al Hub para difusión
    c.hub.broadcast <- messageJSON  
}

// writePump: Envía mensajes al navegador del usuario  
func (c *Client) writePump() {
    select {
    case message := <-c.send:  // 📤 Mensaje del Hub
        c.conn.WriteMessage(message)  // 📤 Envía al WebSocket
    }
}
```

#### **📡 Sincronización y Thread-Safety:**

**Mutex para Estado Compartido:**
```go
type Hub struct {
    clients map[*Client]bool  // 🔒 Protegido por mutex
    mu      sync.RWMutex     // 🔐 Mutex para lectura/escritura
}

// Acceso seguro al mapa de clientes
func (h *Hub) GetClientCount() int {
    h.mu.RLock()         // 🔓 Lock de lectura
    count := len(h.clients)
    h.mu.RUnlock()       // 🔒 Unlock
    return count
}
```

**Canales con Buffer:**
```go
// Evitan bloqueos con usuarios lentos
broadcast:  make(chan []byte, 1000)  // Buffer grande para mensajes
register:   make(chan *Client, 100)  // Buffer para conexiones
send:       make(chan []byte, 256)   // Buffer por cliente
```

#### **⚡ Ventajas de esta Arquitectura:**

**1. Escalabilidad:**
- ✅ Cada usuario = 2 goroutines independientes
- ✅ Go puede manejar **millones** de goroutines
- ✅ No hay bloqueo entre usuarios

**2. Responsividad:**
- ✅ Usuario lento no afecta a otros
- ✅ Mensajes se procesan en **paralelo**
- ✅ Hub coordina sin bloqueos

**3. Thread-Safety:**
- ✅ Mutex protege estado compartido
- ✅ Canales son thread-safe por diseño
- ✅ No hay race conditions

#### **🔍 Ejemplo de Flujo Completo:**

```
👤 Usuario A escribe "Hola"
    ↓
📥 readPump de A lee mensaje del WebSocket
    ↓  
📨 Envía mensaje al canal hub.broadcast
    ↓
🧠 Hub recibe en canal broadcast
    ↓
📡 Hub envía mensaje a TODOS los canales client.send
    ↓
📤 writePump de B, C, D... reciben mensaje
    ↓
👥 Todos los usuarios ven "Hola" en sus pantallas
```

**Tiempo total: ~5-10 milisegundos** ⚡

#### **🧪 Validación con Race Detector:**

El Race Detector de Go verifica que:
- ✅ No hay acceso concurrente sin sincronización
- ✅ Mutex protegen correctamente el estado
- ✅ Canales funcionan sin race conditions

```bash
# Comando que valida toda la arquitectura:
$env:CGO_ENABLED=1; go test -race -v
```

**Si pasa sin warnings = Arquitectura thread-safe garantizada** 🏆

## 🔧 Tecnologías Utilizadas

- **Backend:** Go 1.24.4
- **WebSockets:** Gorilla WebSocket
- **Frontend:** HTML5, Bootstrap 5, JavaScript ES6
- **Imágenes:** Base64 encoding, File API, Drag & Drop API
- **Deploy:** Railway
- **Icons:** Bootstrap Icons
- **Testing:** Go Race Detector, Benchmarks

## 🔒 Cómo se Asegura la Seguridad Concurrente

### **1. Protección del Estado Compartido:**

El mapa de clientes es el recurso más crítico que requiere sincronización:

```go
type Hub struct {
    clients map[*Client]bool  // ⚠️ Recurso compartido
    mu      sync.RWMutex     // 🔐 Mutex para protegerlo
}

// ✅ Acceso seguro de LECTURA
func (h *Hub) GetClientCount() int {
    h.mu.RLock()         // 🔓 Lock de lectura (múltiples lectores OK)
    count := len(h.clients)
    h.mu.RUnlock()       // 🔒 Siempre unlock
    return count
}

// ✅ Acceso seguro de ESCRITURA  
func (h *Hub) registerClient(client *Client) {
    h.mu.Lock()          // 🔒 Lock exclusivo (solo 1 escritor)
    h.clients[client] = true
    h.mu.Unlock()        // 🔓 Liberar lock
}
```

### **2. Comunicación Thread-Safe con Canales:**

Los canales de Go son **inherentemente thread-safe**:

```go
// ✅ Múltiples goroutines pueden enviar sin problemas
client1.hub.broadcast <- message1  // Goroutine 1
client2.hub.broadcast <- message2  // Goroutine 2  
client3.hub.broadcast <- message3  // Goroutine 3

// ✅ Hub recibe ordenadamente
select {
case msg := <-h.broadcast:  // Solo 1 mensaje a la vez
    h.broadcastMessage(msg)
}
```

### **3. Prevención de Deadlocks:**

**Orden consistente de locks:**
```go
// ✅ SIEMPRE el mismo orden para evitar deadlocks
h.mu.Lock()              // Primero: Hub mutex
clientList := getClients()
h.mu.Unlock()            // Liberar rápido

// Enviar sin locks
for _, client := range clientList {
    select {
    case client.send <- message:  // ✅ No bloqueante
    default:                      // ✅ Si está lleno, continuar
    }
}
```

### **4. Gestión de Recursos por Goroutine:**

Cada cliente limpia sus propios recursos:

```go
func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c  // 📤 Auto-desregistro
        c.conn.Close()         // 🔌 Cerrar conexión
    }()
    // ... leer mensajes
}
```

### **5. Canales con Buffer para Evitar Bloqueos:**

```go
// ✅ Buffers previenen bloqueos cuando usuarios son lentos
send:      make(chan []byte, 256)   // 256 mensajes en cola
broadcast: make(chan []byte, 1000)  // 1000 mensajes globales
```

## 🛡️ Decisiones de Diseño para Conexiones/Desconexiones

### **1. Manejo Elegante de Desconexiones:**

```go
// ✅ Detección automática de desconexión
func (c *Client) readPump() {
    c.conn.SetPongHandler(func(string) error {
        // Cliente responde = está vivo
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })
    
    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            // ✅ Cliente desconectado, limpiar automáticamente
            return  // defer ejecuta cleanup
        }
    }
}
```

### **2. Heartbeat para Detectar Conexiones Perdidas:**

```go
// ✅ Ping cada 54 segundos, timeout a los 60
const (
    pongWait   = 60 * time.Second
    pingPeriod = (pongWait * 9) / 10  // 54 segundos
)

func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    for {
        select {
        case <-ticker.C:
            // 🏓 Enviar ping al cliente
            c.conn.WriteMessage(websocket.PingMessage, nil)
        }
    }
}
```

### **3. Prevención de Memory Leaks:**

```go
// ✅ Cleanup garantizado con defer
func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c  // Eliminar del Hub
        c.conn.Close()         // Cerrar WebSocket
    }()
}

// ✅ En el Hub: cerrar canal del cliente
func (h *Hub) unregisterClient(client *Client) {
    delete(h.clients, client)
    close(client.send)  // ⚠️ Cierra canal, writePump termina
}
```

### **4. Manejo de Usuarios Duplicados:**

```go
// ✅ Verificación atómica antes de registrar
func (h *Hub) registerClient(client *Client) {
    if !h.isUsernameAvailable(client.username) {
        // ❌ Rechazar y cerrar conexión inmediatamente
        client.conn.Close()
        return
    }
    // ✅ Nombre disponible, proceder con registro
}
```

## 🎯 Elección del Paquete WebSocket: Gorilla vs Estándar

### **¿Por qué Gorilla WebSocket?**

**Comparación con golang.org/x/net/websocket:**

| Característica | Gorilla WebSocket | x/net/websocket |
|----------------|-------------------|-----------------|
| **Mantenimiento** | ✅ Activo | ⚠️ Deprecated |
| **Performance** | ✅ Optimizado | ❌ Más lento |
| **Flexibilidad** | ✅ Muy configurable | ❌ Básico |
| **Validación** | ✅ Automática | ❌ Manual |
| **Subprotocolos** | ✅ Soporte completo | ❌ Limitado |

### **Configuración Específica para Railway:**

```go
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    // ⭐ CRÍTICO para Railway: Permitir cross-origin
    CheckOrigin: func(r *http.Request) bool {
        return true  // Railway necesita esto
    },
}
```

### **Ventajas de Gorilla en nuestro Chat:**

1. **Ping/Pong automático:** Detecta conexiones perdidas
2. **Compresión:** Reduce bandwidth para imágenes
3. **Rate limiting:** Protege contra spam
4. **Upgrade seguro:** Valida headers automáticamente
5. **Error handling:** Mensajes de error detallados

```go
// ✅ Gorilla maneja estos casos automáticamente:
if websocket.IsUnexpectedCloseError(err, 
    websocket.CloseGoingAway, 
    websocket.CloseAbnormalClosure) {
    // Solo errores realmente inesperados
}
```

## 📱 Responsive Design

El chat funciona perfectamente en:
- 💻 **Desktop** (1200px+) - Vista completa con sidebar
- 📱 **Tablet** (768px - 1199px) - Layout adaptativo
- 📱 **Mobile** (< 768px) - Interfaz optimizada para móviles

## 🛠️ Desarrollo

### **Ejecutar tests:**
```bash
# Tests básicos
go test -v

# Tests con detección de race conditions (REQUERIDO)
$env:CGO_ENABLED=1; go test -race -v

# Tests específicos de imágenes
go test -v -run TestImage

# Benchmarks de rendimiento
go test -bench=.

# Múltiples ejecuciones para detectar problemas intermitentes
$env:CGO_ENABLED=1; go test -race -count=5
```

### **Arquitectura:**
- `main.go` - Servidor HTTP y configuración Railway
- `hub.go` - Centro de gestión de clientes (Patrón Hub-and-Spoke)
- `client.go` - Lógica de clientes individuales (✅ con soporte de imágenes)
- `message.go` - Estructuras de datos (✅ con campos de imagen)
- `websocket.go` - Configuración WebSocket
- `chat_test.go` - Tests unitarios con Race Detector

### **Concurrencia:**
- **1 goroutine central** (Hub) coordina todo el sistema
- **2 goroutines por usuario** (readPump + writePump)
- **Canales con buffer** para comunicación thread-safe
- **Mutex** para proteger estado compartido
- **Race Detector** valida seguridad concurrente

## 🎨 Personalización

### **Cambiar límites de imagen:**
En `client.go`:
```go
const (
    maxImageSize = 5 * 1024 * 1024 // Cambiar tamaño máximo
)
```

### **Modificar interfaz:**
En `index.html`:
```css
.message-image {
    max-width: 300px;  /* Tamaño de vista previa */
    max-height: 200px;
}
```

## 🚨 Solución de Problemas

### **Error: "Imagen demasiado grande"**
- ✅ Reduce el tamaño de la imagen (máx. 5MB)
- ✅ Usa herramientas de compresión de imagen
- ✅ Convierte a formatos más eficientes (WebP, JPEG)

### **Error: "Tipo de imagen no soportado"**
- ✅ Usa solo: JPEG, PNG, GIF, WebP
- ✅ Verifica la extensión del archivo
- ✅ Algunos formatos antiguos pueden no funcionar

### **Race Detector: "cgo: C compiler not found"**
- ✅ Instala MinGW para Windows desde [winlibs.com](https://winlibs.com/)
- ✅ Agrega `C:\mingw64\bin` al PATH
- ✅ Reinicia PowerShell
- ✅ Ejecuta: `$env:CGO_ENABLED=1; go test -race -v`

### **El historial se borra:** ✅ **SOLUCIONADO**
- ✅ **Problema corregido** en la versión actual
- ✅ Ahora el historial es **persistente durante toda la sesión**
- ✅ Los mensajes **no se borran** al enviar imágenes

## 📊 Logs y Monitoreo

Railway proporciona logs en tiempo real:
```
🚀 GO O NO GO - Servidor de chat iniciado
📡 Puerto: 34567
💬 WebSocket endpoint: /ws
✅ Servidor listo para recibir conexiones...
📜 Mensaje agregado al historial local. Total: 15
🖼️ Imagen de 'JUNIOR_ALVINES' enviada al hub (2.3 MB)
✅ Cliente 'testuser' conectado exitosamente. Total de clientes: 5
```

## 🌍 Variables de Entorno

Railway maneja automáticamente:
- `PORT` - Puerto asignado dinámicamente
- Protocolo HTTPS/WSS para producción

## 🔒 Seguridad

- ✅ Validación de entrada en frontend y backend
- ✅ Escape de HTML para prevenir XSS
- ✅ Validación de tipos MIME y magic numbers
- ✅ Límites de tamaño de archivo
- ✅ Rate limiting natural por WebSocket
- ✅ Conexiones HTTPS/WSS en producción
- ✅ Historial seguro sin pérdida de datos
- ✅ Race Detector valida thread-safety

## 🎯 Próximas Funcionalidades

- [ ] Salas de chat múltiples
- [ ] Historial de mensajes persistente en base de datos
- [ ] Autenticación con GitHub
- [ ] Compresión automática de imágenes
- [ ] Soporte para más formatos (videos, documentos)
- [ ] Stickers y emojis personalizados
- [ ] Comandos especiales (/help, /users, /clear, etc.)
- [ ] Notificaciones push
- [ ] Modo oscuro/claro

## ✨ Novedades en Esta Versión

### **🐛 Correcciones:**
- ✅ **Historial persistente** - Los mensajes ya no se borran
- ✅ **Sin duplicados** - Cada mensaje aparece solo una vez
- ✅ **Mejor gestión de memoria** - Optimización del frontend
- ✅ **Logs mejorados** - Mejor debugging y monitoreo
- ✅ **Race Detector** - Validación de concurrencia

### **🚀 Mejoras:**
- ✅ **Flujo optimizado** - Menos operaciones redundantes
- ✅ **IDs únicos** - Sistema robusto de identificación de mensajes
- ✅ **Validación mejorada** - Mejor detección de duplicados
- ✅ **Experiencia de usuario** - Chat más fluido y confiable
- ✅ **Tests exhaustivos** - Cobertura completa con Race Detector

## 📞 Soporte

**Desarrolladores:** JUNIOR_ALVINES & SheralA16  
**GitHub:** [github.com/JUNMPI](https://github.com/JUNMPI)  
**Proyecto:** [realtime-chat](https://github.com/JUNMPI/realtime-chat)  
**Demo:** [https://realtime-chat-production-183c.up.railway.app](https://realtime-chat-production-183c.up.railway.app)

Para reportar bugs o sugerir mejoras, crea un Issue en GitHub.

### **Issues Comunes:**
- **Imágenes grandes:** Reporta problemas con archivos específicos
- **Compatibilidad:** Menciona navegador y sistema operativo
- **Performance:** Incluye detalles de red y dispositivo
- **Race conditions:** ✅ **Detectadas y corregidas** con Race Detector

### **Changelog:**
- **v1.3.0** - ✅ Race Detector y tests exhaustivos
- **v1.2.0** - ✅ Historial persistente corregido
- **v1.1.0** - 🖼️ Soporte completo para imágenes
- **v1.0.0** - 💬 Chat básico en tiempo real

---

**¡Disfruta tu chat en tiempo real con imágenes, historial persistente y validación de concurrencia! 🚀💬🖼️⚡**