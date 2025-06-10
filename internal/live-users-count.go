package internal

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	UsersCount = make(map[string]uint32)              // Store user counts by domain
	clients    = make(map[http.ResponseWriter]string) // Store clients and their associated domain
	mu         sync.Mutex                             // Mutex to protect concurrent access
)

func GetLiveUsersCount(w http.ResponseWriter, r *http.Request) {
	isPeeking := r.URL.Query().Has("peek")

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
	pingInterval := time.NewTicker(15 * time.Second)
	defer pingInterval.Stop()

	done := make(chan struct{})

	// Increment the user count for the specified domain
	mu.Lock()
	if !isPeeking {
		UsersCount[domain] += 1
	}
	clients[w] = domain // Associate this client with the domain
	currentCount := UsersCount[domain]
	mu.Unlock()

	// Broadcast the updated count for this domain
	broadcastCount(domain, currentCount)

	// Start ping loop in a goroutine
	go func() {
		for {
			select {
			case <-pingInterval.C:
				mu.Lock()
				_, err := fmt.Fprintf(w, ": ping\n\n") // SSE comment as heartbeat
				if err != nil {
					mu.Unlock()
					done <- struct{}{}
					return
				}
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
				mu.Unlock()
			case <-clientGone:
				done <- struct{}{}
				return
			}
		}
	}()

	// wait for the client to disconnect
	<-done

	// When client disconnects, decrement the count for the domain
	fmt.Println("Client disconnected for domain:", domain)

	mu.Lock()
	if !isPeeking && UsersCount[domain] > 0 {
		UsersCount[domain] -= 1
	}
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
		if clientDomain != domain {
			continue
		}
		_, err := fmt.Fprintf(client, "data: %d\n\n", count)
		if err != nil {
			// If writing fails (e.g., client disconnected), remove the client
			delete(clients, client)
		}

		// Flush to make sure the message is sent immediately
		err = http.NewResponseController(client).Flush()
		if err != nil {
			// Handle flushing error and remove client if necessary
			delete(clients, client)
		}
	}
}

func StartClientCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			mu.Lock()
			for w, domain := range clients {
				// Try sending an SSE ping
				_, err := fmt.Fprintf(w, ": ping\n\n")
				if err != nil {
					fmt.Println("Cleaning up dead client for domain:", domain)
					delete(clients, w)
					if UsersCount[domain] > 0 {
						UsersCount[domain]--
					}
					continue
				}

				// Try flushing
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				} else {
					// If not flushable, remove anyway
					fmt.Println("Client not flushable, removing:", domain)
					delete(clients, w)
					if UsersCount[domain] > 0 {
						UsersCount[domain]--
					}
				}
			}
			mu.Unlock()
		}
	}()
}
