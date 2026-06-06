// Example: Basic search by name and city
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/i5uc/brixhub-go/pkg/brixhub"
)

func main() {
	apiKey := os.Getenv("BRIXHUB_API_KEY")
	if apiKey == "" {
		log.Fatal("Set BRIXHUB_API_KEY environment variable")
	}

	client, err := brixhub.NewClient(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Search by name and city with flexible matching
	req := &brixhub.SearchRequest{
		NomFamille: "Dupont",
		Prenom:     "Jean",
		Ville:      "Paris",
		Flexible:   true,
		PerPage:    10,
	}

	results, _, err := client.Search(ctx, req)
	if err != nil {
		if apiErr, ok := err.(*brixhub.APIError); ok {
			log.Fatalf("API error: %s", apiErr.Error())
		}
		log.Fatal(err)
	}

	fmt.Printf("Found %d results\n\n", len(results.Results))
	
	for _, profile := range results.Results {
		fmt.Printf("Name: %s %s\n", profile.Prenom, profile.NomFamille)
		fmt.Printf("Email: %s\n", profile.Email)
		fmt.Printf("Phone: %s\n", profile.Telephone)
		fmt.Printf("City: %s\n", profile.Ville)
		fmt.Printf("Confidence: %d%%\n", profile.Confidence)
		fmt.Printf("Sources: %v\n", profile.Sources)
		fmt.Println("---")
	}

	// Check remaining quota
	if remaining := client.GetRemainingQuota(); remaining >= 0 {
		fmt.Printf("\nRemaining quota: %d requests today\n", remaining)
	}
}