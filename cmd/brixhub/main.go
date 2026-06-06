// BrixHub CLI - A command-line tool for the BrixHub API
//
// Usage:
//   brixhub search --nom "Dupont" --prenom "Jean"
//   brixhub lookup email jean.dupont@gmail.com
//   brixhub lookup phone 0612345678
//   brixhub account
//   brixhub health
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/yourusername/brixhub-go/pkg/brixhub"
	"github.com/joho/godotenv"
)

var (
	apiKey    string
	userAgent string
	client    *brixhub.Client
)

func main() {
	_ = godotenv.Load()
	rootCmd := &cobra.Command{
		Use:   "brixhub",
		Short: "BrixHub CLI - Search 11B+ documents",
		Long: `BrixHub CLI is a command-line interface for the BrixHub API.
		
Search by name, email, phone, Discord ID, and more across 11+ billion documents.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == "health" {
				return nil // Health doesn't require API key
			}
			
			if apiKey == "" {
				apiKey = os.Getenv("BRIXHUB_API_KEY")
			}
			
			if apiKey == "" && cmd.Name() != "health" {
				return fmt.Errorf("API key required. Set --api-key flag or BRIXHUB_API_KEY environment variable")
			}
			
			var err error
			opts := []brixhub.ClientOption{}
			if userAgent != "" {
				opts = append(opts, brixhub.WithUserAgent(userAgent))
			}
			
			client, err = brixhub.NewClient(apiKey, opts...)
			return err
		},
	}

	rootCmd.PersistentFlags().StringVarP(&apiKey, "api-key", "k", "", "BrixHub API key (or set BRIXHUB_API_KEY)")
	rootCmd.PersistentFlags().StringVarP(&userAgent, "user-agent", "u", "", "Custom User-Agent string")

	// Add commands
	rootCmd.AddCommand(searchCmd())
	rootCmd.AddCommand(lookupCmd())
	rootCmd.AddCommand(accountCmd())
	rootCmd.AddCommand(healthCmd())
	rootCmd.AddCommand(usageCmd())

	if err := rootCmd.Execute(); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
}

func searchCmd() *cobra.Command {
	var req brixhub.SearchRequest
	
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search by multiple criteria",
		Long:  `Perform a multi-criteria search across the BrixHub database.`,
		Example: `  brixhub search --nom "Dupont" --prenom "Jean"
  brixhub search --email "jean.dupont@gmail.com"
  brixhub search --discord-id "123456789"
  brixhub search --nom "Martin" --ville "Paris" --flexible`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			
			color.Cyan("🔍 Searching...")
			
			results, meta, err := client.Search(ctx, &req)
			if err != nil {
				if apiErr, ok := err.(*brixhub.APIError); ok {
					return fmt.Errorf("API error: %s", apiErr.Error())
				}
				return err
			}
			
			// Display results
			if len(results.Results) == 0 {
				color.Yellow("No results found")
				return nil
			}
			
			color.Green("\n✓ Found %d result(s) (page %d of %d)", 
				len(results.Results), meta.Page, meta.Pages)
			
			for i, profile := range results.Results {
				displayProfile(i+1, profile)
			}
			
			// Show quota info
			if remaining := client.GetRemainingQuota(); remaining >= 0 {
				color.Cyan("\n📊 Remaining quota: %d requests today", remaining)
			}
			
			return nil
		},
	}
	
	// Identity flags
	cmd.Flags().StringVarP(&req.NomFamille, "nom", "n", "", "Last name")
	cmd.Flags().StringVarP(&req.Prenom, "prenom", "p", "", "First name")
	cmd.Flags().StringVar(&req.NomNaissance, "nom-naissance", "", "Birth name")
	cmd.Flags().StringVar(&req.NomAffichage, "nom-affichage", "", "Display name")
	cmd.Flags().StringVar(&req.NomUtilisateur, "nom-utilisateur", "", "Username")
	cmd.Flags().StringVar(&req.DateNaissance, "date-naissance", "", "Birth date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&req.Genre, "genre", "", "Gender (M/F)")
	
	// Contact flags
	cmd.Flags().StringVarP(&req.Email, "email", "e", "", "Email address")
	cmd.Flags().StringVarP(&req.Telephone, "telephone", "t", "", "Phone number")
	cmd.Flags().StringVar(&req.Mobile, "mobile", "", "Mobile number")
	cmd.Flags().StringVar(&req.AdresseIP, "ip", "", "IP address")
	
	// Address flags
	cmd.Flags().StringVar(&req.Adresse, "adresse", "", "Street address")
	cmd.Flags().StringVar(&req.CodePostal, "cp", "", "Postal code")
	cmd.Flags().StringVar(&req.Ville, "ville", "", "City")
	cmd.Flags().StringVar(&req.Pays, "pays", "", "Country")
	
	// Gaming flags
	cmd.Flags().StringVar(&req.SteamID, "steam-id", "", "Steam ID")
	cmd.Flags().StringVar(&req.DiscordID, "discord-id", "", "Discord ID")
	cmd.Flags().StringVar(&req.FiveMLicense, "fivem-license", "", "FiveM license")
	
	// Vehicle flags
	cmd.Flags().StringVar(&req.Immatriculation, "immat", "", "License plate")
	cmd.Flags().StringVar(&req.VINPlaque, "vin", "", "VIN number")
	
	// Professional flags
	cmd.Flags().StringVar(&req.SIRET, "siret", "", "SIRET number")
	cmd.Flags().StringVar(&req.SIREN, "siren", "", "SIREN number")
	cmd.Flags().StringVar(&req.Societe, "societe", "", "Company name")
	
	// Options
	cmd.Flags().IntVar(&req.Page, "page", 1, "Page number")
	cmd.Flags().IntVar(&req.PerPage, "per-page", 10, "Results per page")
	cmd.Flags().BoolVar(&req.Flexible, "flexible", false, "Flexible search (contains vs exact)")
	
	return cmd
}

func lookupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lookup [email|phone|iban] [value]",
		Short: "Reverse lookup by single identifier",
		Long:  `Perform a reverse lookup using email, phone number, or IBAN.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			lookupType := strings.ToLower(args[0])
			value := args[1]
			
			var profile *brixhub.Profile
			var err error
			
			switch lookupType {
			case "email":
				color.Cyan("🔍 Looking up email: %s", value)
				profile, err = client.LookupEmail(ctx, value)
			case "phone", "tel":
				color.Cyan("🔍 Looking up phone: %s", value)
				profile, err = client.LookupPhone(ctx, value)
			case "iban":
				color.Cyan("🔍 Looking up IBAN: %s", value)
				profile, err = client.LookupIBAN(ctx, value)
			default:
				return fmt.Errorf("unknown lookup type: %s (use email, phone, or iban)", lookupType)
			}
			
			if err != nil {
				return err
			}
			
			if profile == nil {
				color.Yellow("No results found")
				return nil
			}
			
			displayProfile(1, *profile)
			
			if remaining := client.GetRemainingQuota(); remaining >= 0 {
				color.Cyan("\n📊 Remaining quota: %d requests today", remaining)
			}
			
			return nil
		},
	}
	
	return cmd
}

func accountCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "account",
		Short: "Show account information",
		Long:  `Display your plan, quota, and usage statistics.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			
			info, err := client.GetAccountInfo(ctx)
			if err != nil {
				return err
			}
			
			color.Cyan("👤 Account Information\n")
			color.White("Plan:              %s", info.Plan)
			color.White("Daily Quota:       %d requests", info.DailyQuota)
			color.White("Daily Used:        %d requests", info.DailyUsed)
			color.White("Daily Remaining:   %d requests", info.DailyRemaining)
			color.White("Total Requests:    %d", info.TotalRequests)
			color.White("Results per Query: %d", info.ResultsPerQuery)
			color.White("Pagination:        %v", info.PaginationEnabled)
			
			return nil
		},
	}
}

func usageCmd() *cobra.Command {
	var limit, offset int
	
	cmd := &cobra.Command{
		Use:   "usage",
		Short: "Show usage history (Pro+ plans)",
		Long:  `Display detailed API usage history. Requires Pro plan or higher.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			
			usage, err := client.GetUsage(ctx, limit, offset)
			if err != nil {
				if apiErr, ok := err.(*brixhub.APIError); ok && apiErr.IsPlanLimited() {
					return fmt.Errorf("this feature requires a Pro or higher plan")
				}
				return err
			}
			
			color.Cyan("📊 Usage History\n")
			
			for _, log := range usage.Logs {
				color.White("\n[%s] %s", log.Timestamp.Format("2006-01-02 15:04"), log.Endpoint)
				color.White("  Query: %s", log.Query)
				color.White("  Results: %d | Time: %dms", log.ResultsCount, log.TookMs)
			}
			
			return nil
		},
	}
	
	cmd.Flags().IntVarP(&limit, "limit", "l", 50, "Number of logs to show")
	cmd.Flags().IntVarP(&offset, "offset", "o", 0, "Offset for pagination")
	
	return cmd
}

func healthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check API health status",
		Long:  `Check if the BrixHub API is operational.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create a temporary client without API key
			tempClient, _ := brixhub.NewClient("dummy")
			
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			
			health, err := tempClient.Health(ctx)
			if err != nil {
				color.Red("❌ API is unreachable: %v", err)
				return nil
			}
			
			if health.Status == "operational" {
				color.Green("✅ API is operational")
			} else {
				color.Yellow("⚠️  API status: %s", health.Status)
			}
			
			return nil
		},
	}
}

func displayProfile(index int, p brixhub.Profile) {
	cyan := color.New(color.FgCyan, color.Bold)
	white := color.New(color.FgWhite)
	yellow := color.New(color.FgYellow)
	
	cyan.Printf("\n[%d] ", index)
	
	name := strings.TrimSpace(fmt.Sprintf("%s %s", p.Prenom, p.NomFamille))
	if name == "" && p.NomAffichage != "" {
		name = p.NomAffichage
	}
	if name == "" {
		name = "Unknown"
	}
	white.Printf("%s\n", name)
	
	if p.Email != "" {
		white.Printf("   📧 %s\n", p.Email)
	}
	if p.Telephone != "" {
		white.Printf("   📞 %s\n", p.Telephone)
	}
	if p.Mobile != "" && p.Mobile != p.Telephone {
		white.Printf("   📱 %s\n", p.Mobile)
	}
	if p.Ville != "" || p.CodePostal != "" {
		white.Printf("   📍 %s %s\n", p.CodePostal, p.Ville)
	}
	if p.DateNaissance != "" {
		white.Printf("   🎂 %s\n", p.DateNaissance)
	}
	if p.DiscordID != "" {
		yellow.Printf("   🎮 Discord: %s\n", p.DiscordID)
	}
	if p.SteamID != "" {
		yellow.Printf("   🎮 Steam: %s\n", p.SteamID)
	}
	if len(p.Sources) > 0 {
		yellow.Printf("   📋 Sources: %s\n", strings.Join(p.Sources, ", "))
	}
	if p.Confidence > 0 {
		yellow.Printf("   ✅ Confidence: %d%%\n", p.Confidence)
	}
}