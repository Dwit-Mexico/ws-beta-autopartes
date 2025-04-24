package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/RomanshkVolkov/ws-beta-autopartes/internal/adapters/repository"
	"github.com/RomanshkVolkov/ws-beta-autopartes/internal/core/domain"
	"github.com/RomanshkVolkov/ws-beta-autopartes/internal/core/lg"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	clients      = make(map[*domain.WSClient]bool) // Clientes WebSocket conectados.
	lastData     = make(map[int]time.Time)         // Mapa para almacenar la última fecha de actualización por SucursalID.
	lastDataLock sync.Mutex                        // Mutex para proteger el acceso a lastData.
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Permite todas las conexiones (ajusta según sea necesario).
	},
}

// HandlerWebSocket maneja las conexiones WebSocket.
func HandlerWebSocket(c *gin.Context) {
	// Actualiza la conexión HTTP a WebSocket.
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error al actualizar la conexión: %v", err)
		return
	}
	defer conn.Close()

	// Crea un nuevo cliente.
	client := &domain.WSClient{Conn: conn, Send: make(chan []byte)}
	clients[client] = true
	defer func() {
		delete(clients, client)
		close(client.Send)
	}()

	// Escucha mensajes del cliente (opcional).
	// go handleMessages(client)

	// Envía mensajes desde el canal Send al cliente.
	for {
		select {
		case msg := <-client.Send:
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Error al escribir mensaje: %v", err)
				return
			}
		}
	}
}

// PollWarehouses realiza polling a la base de datos y envía notificaciones a los clientes.
func PollWarehouses() {
	ticker := time.NewTicker(5 * time.Second) // Intervalo de polling.
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println(clients)
		if len(clients) == 0 {
			return
		}
		data, err := repository.GetWebSocketWarehouses()
		if err != nil {
			log.Printf("Error al obtener los almacenes: %v", err)
			continue
		}

		var updated bool
		lastDataLock.Lock()
		for _, warehouse := range data {
			if existing, exists := lastData[warehouse.SucursalID]; !exists || warehouse.LastUpdatedAt.After(existing) {
				lastData[warehouse.SucursalID] = warehouse.LastUpdatedAt
				updated = true
				lg.Info(fmt.Sprintf("Almacén actualizado: %s", strconv.Itoa(warehouse.SucursalID)))
			}
		}
		lastDataLock.Unlock()

		// Si hay cambios, notifica a los clientes.
		if updated {
			notifyClients()
		}
	}
}

// notifyClients envía una notificación a todos los clientes conectados.
func notifyClients() {
	message := []byte("update") // Mensaje simple para notificar cambios.
	for client := range clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(clients, client)
		}
	}
}

// handleMessages maneja mensajes recibidos de los clientes (opcional).
func handleMessages(client *domain.WSClient) {
	for {
		_, _, err := client.Conn.ReadMessage()
		if err != nil {
			log.Printf("Error al leer mensaje del cliente: %v", err)
			delete(clients, client)
			break
		}
	}
}
