package internal

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	UsersCount = make(map[string]uint32)              // Store user counts by domain
	clients    = make(map[http.ResponseWriter]string) // Store clients and their associated domain
	mu         sync.Mutex                             // Mutex to protect concurrent access
)

func GetLiveUsersCount(w http.ResponseWriter, r *http.Request) {
	// Extract the domain from the query parameters
	domain := r.URL.Query().Get("d")
	if domain == "" {
		http.Error(w, "Missing 'd' query parameter", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	clientGone := r.Context().Done()

	// Increment the user count for the specified domain
	mu.Lock()
	UsersCount[domain] += 1
	currentCount := UsersCount[domain]
	clients[w] = domain // Associate this client with the domain
	mu.Unlock()

	// Broadcast the updated count for this domain
	broadcastCount(domain, currentCount)

	// Wait for client disconnect
	<-clientGone

	// When client disconnects, decrement the count for the domain
	fmt.Println("Client disconnected for domain:", domain)

	mu.Lock()
	UsersCount[domain] -= 1
	updatedCount := UsersCount[domain]
	delete(clients, w) // Remove the client
	mu.Unlock()

	// Broadcast the updated count for this domain after disconnection
	broadcastCount(domain, updatedCount)
}

// broadcastCount sends the current user count for a specific domain to all connected clients associated with that domain
func broadcastCount(domain string, count uint32) {
	mu.Lock()
	defer mu.Unlock()

	for client, clientDomain := range clients {
		if clientDomain == domain {
			_, err := fmt.Fprintf(client, "data: %d\n\n", count)
			if err != nil {
				// If writing fails (e.g., client disconnected), remove the client
				delete(clients, client)
				continue
			}

			// Flush to make sure the message is sent immediately
			err = http.NewResponseController(client).Flush()
			if err != nil {
				// Handle flushing error and remove client if necessary
				delete(clients, client)
			}
		}
	}
}
