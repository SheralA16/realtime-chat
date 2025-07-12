# 🚀 GO O NO GO - Chat en Tiempo Real

Un sistema de chat en tiempo real desarrollado en Go utilizando WebSockets, con interfaz Bootstrap y control avanzado de usuarios duplicados.

## 📋 Características Principales

### ✅ **Control de Usuarios Duplicados**
- **Validación en tiempo real**: No permite conexiones con nombres de usuario ya ocupados
- **Mensajes de error claros**: Notifica al usuario si el nombre está en uso
- **Validación de formato**: Solo permite letras, números, guiones y guiones bajos
- **Límites de longitud**: Entre 2 y 20 caracteres

### 🎨 **Interfaz Moderna con Bootstrap**
- **Diseño responsive**: Funciona en desktop, tablet y móvil
- **Componentes Bootstrap 5**: Cards, alerts, toasts, badges
- **Iconos Bootstrap**: Interfaz visual intuitiva
- **Tema personalizado**: Gradientes verdes

### 💬 **Funcionalidades del Chat**
- **Mensajes en tiempo real**: Difusión instantánea a todos los usuarios
- **Lista de usuarios activos**: Muestra quién está conectado/desconectado
- **Mensajes del sistema**: Notificaciones de conexión/desconexión
- **Historial de usuarios**: Mantiene registro de usuarios pasados
- **Timestamps**: Hora de cada mensaje

### 🔧 **Arquitectura Técnica**
- **Backend en Go**: Gorilla WebSocket para conexiones concurrentes
- **Concurrencia segura**: Mutex para operaciones thread-safe
- **Canales Go**: Comunicación entre goroutines
- **Gestión de memoria**: Cleanup automático de recursos

## 📁 Estructura del Proyecto

```
realtime-chat/
├── main.go              # Servidor HTTP principal
├── hub.go               # Gestión central de clientes y mensajes
├── client.go            # Manejo de clientes WebSocket individuales
├── message.go           # Estructuras de mensajes
├── websocket.go         # Configuración y upgrade de WebSocket
├── index.html           # Frontend con Bootstrap
├── chat_test.go         # Tests unitarios
├── go.mod              # Dependencias de Go
├── go.sum              # Checksums de dependencias
└── README.md           # Esta documentación
```

## 🚀 Instalación y Ejecución

### **Prerrequisitos**
- Go 1.19 o superior
- Navegador web moderno con soporte WebSocket

### **Pasos de Instalación**

1. **Clonar el repositorio**
```bash
git clone <repository-url>
cd realtime-chat
```

2. **Instalar dependencias**
```bash
go mod tidy
```

3. **Ejecutar el servidor**
```bash
go run *.go
```

4. **Abrir en el navegador**
```
http://localhost:8080
```

### **Comandos Útiles**

```bash
# Ejecutar tests
go test -v

# Ejecutar tests con detección de race conditions
go test -race -v

# Benchmark de rendimiento
go test -bench=.

# Ejecutar con logs detallados
go run *.go -v
```

## 🧪 Testing

El proyecto incluye tests completos para validar:

- **Creación del Hub**: Inicialización correcta
- **Registro de clientes**: Conexión y validación
- **Desregistro de clientes**: Desconexión limpia
- **Difusión de mensajes**: Broadcast a todos los usuarios
- **Operaciones concurrentes**: Múltiples usuarios simultáneos
- **Condiciones de carrera**: Seguridad thread-safe
- **Integración WebSocket**: Tests end-to-end

```bash
# Ejecutar todos los tests
go test -v

# Test específico
go test -run TestClientRegistration -v

# Tests con race detection
go test -race -v
```

## 🔒 Validación de Usuarios

### **Reglas de Nombres de Usuario**
- **Longitud**: 2-20 caracteres
- **Caracteres permitidos**: `a-z`, `A-Z`, `0-9`, `-`, `_`
- **Unicidad**: No se permiten nombres duplicados
- **Case sensitive**: "Usuario" y "usuario" son diferentes

### **Flujo de Validación**
1. **Frontend**: Validación inicial en JavaScript
2. **Backend**: Verificación de disponibilidad en el hub
3. **Respuesta**: Error específico si el nombre está en uso
4. **Cleanup**: Limpieza automática de recursos en caso de error

## 🌐 Arquitectura de WebSocket

### **Flujo de Conexión**
```
Cliente → HTTP Upgrade → WebSocket → Validación → Registro en Hub → Broadcast
```

### **Manejo de Mensajes**
```go
type Message struct {
    Username  string    `json:"username"`
    Content   string    `json:"content"`
    Timestamp time.Time `json:"timestamp"`
    Type      string    `json:"type"` // "message", "system", "join", "leave"
}
```

### **Tipos de Mensajes**
- **message**: Mensajes regulares de chat
- **system**: Mensajes del sistema
- **join**: Usuario se conecta
- **leave**: Usuario se desconecta
- **error**: Errores de validación
- **connectionSuccess**: Confirmación de conexión
- **userList**: Lista actualizada de usuarios

## 📊 Gestión de Concurrencia

### **Primitivas Utilizadas**
- **sync.RWMutex**: Protección de mapas compartidos
- **Channels**: Comunicación entre goroutines
- **Goroutines**: Manejo concurrente de clientes

### **Patrón de Diseño**
```go
// Hub centralizado con canales
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}
```

## 🔧 Configuración Avanzada

### **Timeouts de WebSocket**
```go
const (
    writeWait      = 10 * time.Second  // Timeout de escritura
    pongWait       = 60 * time.Second  // Timeout de pong
    pingPeriod     = 54 * time.Second  // Periodo de ping
    maxMessageSize = 512               // Tamaño máximo de mensaje
)
```

### **Buffers de Canales**
```go
broadcast:  make(chan []byte, 1000)  // Buffer grande para mensajes
register:   make(chan *Client, 100)  // Buffer para nuevos clientes
unregister: make(chan *Client, 100)  // Buffer para desconexiones
send:       make(chan []byte, 256)   // Buffer por cliente
```

## 🚨 Manejo de Errores

### **Tipos de Errores Manejados**
- **Nombres duplicados**: USERNAME_TAKEN
- **Conexión WebSocket**: Upgrade failures
- **Validación de entrada**: Formato inválido
- **Timeouts**: Ping/Pong failures
- **Recursos**: Memory leaks y cleanup

### **Estrategias de Recovery**
- **Graceful shutdown**: Cierre ordenado de conexiones
- **Resource cleanup**: Liberación automática de memoria
- **Error propagation**: Mensajes claros al usuario
- **Logging**: Registro detallado para debugging

## 📈 Rendimiento

### **Métricas Objetivo**
- **Usuarios concurrentes**: 1000+ conexiones simultáneas
- **Latencia de mensajes**: < 50ms
- **Throughput**: 10,000+ mensajes/segundo
- **Memoria por usuario**: < 1MB

### **Optimizaciones Implementadas**
- **Buffered channels**: Evita bloqueos
- **Connection pooling**: Reutilización eficiente
- **Goroutine per connection**: Escalabilidad
- **Memory-efficient structures**: Structs optimizados

## 🔐 Seguridad

### **Medidas Implementadas**
- **Input validation**: Sanitización de nombres y mensajes
- **Rate limiting**: Control de spam (futuro)
- **CORS policy**: Configuración de orígenes permitidos
- **XSS prevention**: Escape de HTML en mensajes

### **Consideraciones de Producción**
- **HTTPS/WSS**: Encriptación en producción
- **Authentication**: Sistema de autenticación (futuro)
- **Authorization**: Permisos por sala (futuro)
- **Monitoring**: Métricas y alertas

## 🤝 Contribución

### **Cómo Contribuir**
1. Fork del repositorio
2. Crear rama feature: `git checkout -b feature/nueva-funcionalidad`
3. Commit cambios: `git commit -am 'Agregar nueva funcionalidad'`
4. Push a la rama: `git push origin feature/nueva-funcionalidad`
5. Crear Pull Request

### **Estándares de Código**
- **Go fmt**: Formato estándar de Go
- **Go vet**: Análisis estático
- **Tests**: Coverage mínimo del 80%
- **Documentación**: Comentarios en funciones públicas

## 📝 Licencia

Este proyecto está bajo la Licencia MIT. Ver `LICENSE` para más detalles.

## 👥 Autores

- **Desarrollador Principal** - Implementación inicial y arquitectura

## 🙏 Agradecimientos

- **Gorilla WebSocket**: Excelente librería para WebSockets en Go
- **Bootstrap**: Framework CSS para interfaz moderna
- **Comunidad Go**: Documentación y mejores prácticas

---

## 📞 Soporte

Para soporte técnico o preguntas:
- **Issues**: Crear issue en el repositorio
- **Documentación**: Revisar este README
- **Tests**: Ejecutar `go test -v` para validar configuraci