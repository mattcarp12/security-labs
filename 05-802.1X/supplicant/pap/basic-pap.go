package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

func main() {
	fmt.Println("--- 802.1X Go Supplicant Initialized ---")

	// 1. Define the Shared Secret
	// This proves to the RADIUS server that we are a trusted network switch
	secret := []byte("GoRadiusSecret1234")

	// 2. Craft the Access-Request Packet (Code 1)
	packet := radius.New(radius.CodeAccessRequest, secret)

	// 3. Append the Attributes (Type-Length-Value)
	// We are passing in the exact credentials we stored in the FreeRADIUS users file
	username := "alice"
	password := "SecurityPlus8021x!"

	rfc2865.UserName_SetString(packet, username)
	rfc2865.UserPassword_SetString(packet, password)

	// 4. Configure the Network Connection
	// Pointing to our local Docker container on the standard RADIUS UDP port
	client := radius.Client{
		Retry: 500 * time.Millisecond, // Retry up to 3 times if no response is received
	}
	radAddr := "127.0.0.1:1812"

	fmt.Printf("Sending Access-Request for user '%s' to %s...\n", username, radAddr)

	// 5. Fire the packet over the network and wait for the Boss to reply
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := client.Exchange(ctx, packet, radAddr)
	if err != nil {
		log.Fatalf("Network Error: %v", err)
	}

	// 6. Evaluate the Decision
	fmt.Println("----------------------------------------")
	if response.Code == radius.CodeAccessAccept {
		fmt.Println("[SUCCESS] RADIUS Code 2: Access-Accept")
		fmt.Println("The credentials are valid. The Switch port is now UNLOCKED.")
	} else if response.Code == radius.CodeAccessReject {
		fmt.Println("[REJECTED] RADIUS Code 3: Access-Reject")
		fmt.Println("Invalid credentials. The Switch port remains LOCKED (Zero Trust).")
	} else {
		fmt.Printf("[UNKNOWN] Received unexpected code: %v\n", response.Code)
	}
}
