// Example: Reverse lookup by email, phone, or IBAN
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/i5uc/brixhub-go/pkg/brixhub"
)

func main() {
	_ = godotenv.Load()
	apiKey := os.Getenv("BRIXHUB_API_KEY")
	if apiKey == "" {
		log.Fatal("Set BRIXHUB_API_KEY environment variable")
	}

	client, err := brixhub.NewClient(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Lookup by email
	fmt.Println("=== Lookup by Email ===")
	profile, err := client.LookupEmail(ctx, "jean.dupont@example.com")
	if err != nil {
		log.Printf("Email lookup error: %v", err)
	} else if profile != nil {
		fmt.Printf("Found: %s %s\n", profile.Prenom, profile.NomFamille)
	}

	// Lookup by phone
	fmt.Println("\n=== Lookup by Phone ===")
	profile, err = client.LookupPhone(ctx, "0612345678")
	if err != nil {
		log.Printf("Phone lookup error: %v", err)
	} else if profile != nil {
		fmt.Printf("Found: %s %s\n", profile.Prenom, profile.NomFamille)
	}

	// Lookup by IBAN
	fmt.Println("\n=== Lookup by IBAN ===")
	profile, err = client.LookupIBAN(ctx, "FR7630006000011234567890189")
	if err != nil {
		log.Printf("IBAN lookup error: %v", err)
	} else if profile != nil {
		fmt.Printf("Found: %s %s\n", profile.Prenom, profile.NomFamille)
		fmt.Printf("BIC: %s\n", profile.BIC)
	}
}