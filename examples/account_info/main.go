// Example: Check account info and usage
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

	// Get account info
	info, err := client.GetAccountInfo(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Plan: %s\n", info.Plan)
	fmt.Printf("Daily Quota: %d\n", info.DailyQuota)
	fmt.Printf("Daily Used: %d\n", info.DailyUsed)
	fmt.Printf("Daily Remaining: %d\n", info.DailyRemaining)
	fmt.Printf("Total Requests: %d\n", info.TotalRequests)
	fmt.Printf("Results per Query: %d\n", info.ResultsPerQuery)
	fmt.Printf("Pagination Enabled: %v\n", info.PaginationEnabled)

	// Get usage history (Pro+ only)
	usage, err := client.GetUsage(ctx, 10, 0)
	if err != nil {
		if apiErr, ok := err.(*brixhub.APIError); ok && apiErr.IsPlanLimited() {
			fmt.Println("\nUsage history requires Pro+ plan")
			return
		}
		log.Fatal(err)
	}

	fmt.Printf("\nRecent Usage:\n")
	for _, log := range usage.Logs {
		fmt.Printf("  %s: %s (%d results, %dms)\n", 
			log.Timestamp.Format("2006-01-02 15:04"),
			log.Endpoint,
			log.ResultsCount,
			log.TookMs)
	}
}