package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2869" // We need this RFC for EAP attributes
)

func main() {
	fmt.Println("--- 802.1X Go Supplicant (EAP Mode) ---")

	secret := []byte("GoRadiusSecret123")
	username := "alice"

	// 1. Craft the RADIUS Access-Request
	packet := radius.New(radius.CodeAccessRequest, secret)

	// 2. Cryptographic Packet Signing (Message-Authenticator)
	// BEST PRACTICE: Always make the authenticator the FIRST attribute in the packet 
	// to prevent linear parsing bugs in strict RADIUS servers.
	rfc2869.MessageAuthenticator_Set(packet, make([]byte, 16))

	// 3. Add standard attributes
	rfc2865.UserName_SetString(packet, username)

	// 4. Craft the RAW binary EAP-Identity Envelope
	eapLength := 5 + len(username)
	eapIdentityPacket := []byte{
		2,                 // Code 2 = EAP Response
		1,                 // Identifier (Sequence Number)
		0, byte(eapLength),// Length (High byte, Low byte)
		1,                 // Type 1 = Identity
	}
	eapIdentityPacket = append(eapIdentityPacket, []byte(username)...)

	// 5. Stuff the EAP envelope inside the RADIUS packet
	rfc2869.EAPMessage_Set(packet, eapIdentityPacket)

	// 3.5. Cryptographic Packet Signing (Message-Authenticator)
	// We allocate 16 empty bytes. The RADIUS library will automatically 
	// compute the HMAC-MD5 hash of the packet right before it leaves the NIC.
	// rfc2869.MessageAuthenticator_Set(packet, make([]byte, 16))

	// 4. Configure the Network Client
	client := radius.Client{
		Net:        "udp",
		Retry:      500 * time.Millisecond,
	}
	radAddr := "127.0.0.1:1812"


	fmt.Printf("Sending EAP-Identity for '%s' to %s...\n", username, radAddr)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := client.Exchange(ctx, packet, radAddr)
	if err != nil {
		log.Fatalf("Network Error: %v", err)
	}

	// 5. Evaluate the Boss's Decision
	fmt.Println("----------------------------------------")
	if response.Code == radius.CodeAccessAccept {
		fmt.Println("[SUCCESS] RADIUS Code 2: Access-Accept")
	} else if response.Code == radius.CodeAccessReject {
		fmt.Println("[REJECTED] RADIUS Code 3: Access-Reject")
	} else if response.Code == radius.CodeAccessChallenge {
		// THIS IS WHAT WE WANT!
		fmt.Println("[CHALLENGE] RADIUS Code 11: Access-Challenge")
		fmt.Println("The server received our EAP-Identity and is demanding cryptographic proof.")
		
		// Let's peek inside the server's reply to see what EAP type it wants
		eapReply := rfc2869.EAPMessage_Get(response)
		if len(eapReply) > 4 {
			eapType := eapReply[4]
			fmt.Printf("Server is requesting EAP Type: %d\n", eapType)
		}
	} else {
		fmt.Printf("[UNKNOWN] Code: %v\n", response.Code)
	}
}