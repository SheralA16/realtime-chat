package main

import (
	"log"
	"net/http"
)

func main() {
	// Crear el hub de chat
	hub := NewHub()

	// Iniciar el hub en una goroutine separada
	go hub.Run()

	// Configurar rutas HTTP
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	})

	// Servir archivos estáticos desde el directorio ./static/
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Información de inicio
	log.Println("🚀 Servidor de chat iniciado")
	log.Println("📡 Puerto: 8080")
	log.Println("🌐 URL: http://localhost:8080")
	log.Println("💬 WebSocket endpoint: ws://localhost:8080/ws")
	log.Println("📁 Archivos estáticos servidos desde: ./static/")
	log.Println("✅ Servidor listo para recibir conexiones...")

	// Iniciar servidor HTTP
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("❌ Error iniciando servidor HTTP:", err)
	}
}

// serveHome sirve la página principal del chat
func serveHome(w http.ResponseWriter, r *http.Request) {
	// Verificar que sea la ruta raíz
	if r.URL.Path != "/" {
		http.Error(w, "Página no encontrada", http.StatusNotFound)
		return
	}

	// Solo permitir método GET
	if r.Method != "GET" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Servir el archivo index.html
	http.ServeFile(w, r, "index.html")
}
