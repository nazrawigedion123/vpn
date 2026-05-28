package main

import "C" // Required for C-shared libraries
import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// The //export comment is strictly required! It tells the Go compiler
// to expose this function to C/Rust.
//export StartGoVPNClient
func StartGoVPNClient(serverURL *C.char) {
	// Convert the C string passed from Rust into a Go string
	address := C.GoString(serverURL)
	fmt.Printf("[Go VPN] Starting background client...\n")
	fmt.Printf("[Go VPN] Connecting to Hub: %s\n", address)

	u, err := url.Parse(address)
	if err != nil {
		log.Printf("[Go VPN Error] Invalid URL: %v", err)
		return
	}

	// Connect to the WebSocket Hub via HTTPS/443 (bypassing firewalls)
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("[Go VPN Error] Failed to dial cPanel hub: %v", err)
		return
	}
	defer conn.Close()

	fmt.Println("[Go VPN] Successfully connected to the tunnel!")

	// Keep the tunnel alive and listen for incoming data from the hub.
	// In a full implementation, you would pipe this `message` data 
	// into a local TCP socket or a TUN interface for libp2p to read.
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("[Go VPN Error] Connection dropped:", err)
			return
		}
		log.Printf("[Go VPN] Received packet from hub: %s", message)
	}
}

// main is required by the compiler for c-shared libraries, but left empty.
func main() {}