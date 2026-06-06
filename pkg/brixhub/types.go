// Package brixhub provides a Go client for the BrixHub API v1.
// BrixHub is a database search platform with over 11 billion documents.
package brixhub

import "time"

const (
	// BaseURL is the production API endpoint
	BaseURL = "https://brixhub.net/api/v1"
	
	// DefaultUserAgent is the default User-Agent header value
	DefaultUserAgent = "brixhub-go/1.0.0"
)

// Response is the standard API response wrapper
type Response struct {
	Status    int             `json:"status"`
	Message   string          `json:"message"`
	Data      json.RawMessage `json:"data"`
	Meta      *Meta           `json:"meta,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

// Meta contains pagination and timing information
type Meta struct {
	Total         int  `json:"total"`
	TotalIsCapped bool `json:"total_is_capped"`
	Page          int  `json:"page"`
	PerPage       int  `json:"per_page"`
	Pages         int  `json:"pages"`
	TookMs        int  `json:"took_ms"`
}

// RateLimit contains rate limiting information from response headers
type RateLimit struct {
	LimitDay      int
	RemainingDay  int
	LimitMin      int
}

// Profile represents a merged profile result from search
type Profile struct {
	// Identity
	NomFamille    string `json:"nom_famille,omitempty"`
	Prenom        string `json:"prenom,omitempty"`
	NomNaissance  string `json:"nom_naissance,omitempty"`
	NomAffichage  string `json:"nom_affichage,omitempty"`
	NomUtilisateur string `json:"nom_utilisateur,omitempty"`
	DateNaissance string `json:"date_naissance,omitempty"`
	AnneeNaissance string `json:"annee_naissance,omitempty"`
	JourNaissance int    `json:"jour_naissance,omitempty"`
	MoisNaissance int    `json:"mois_naissance,omitempty"`
	Genre         string `json:"genre,omitempty"`
	Civilite      string `json:"civilite,omitempty"`
	
	// Contact
	Email         string `json:"email,omitempty"`
	Telephone     string `json:"telephone,omitempty"`
	Mobile        string `json:"mobile,omitempty"`
	AdresseIP     string `json:"adresse_ip,omitempty"`
	
	// Address
	Adresse           string `json:"adresse,omitempty"`
	ComplementAdresse string `json:"complement_adresse,omitempty"`
	CodePostal        string `json:"code_postal,omitempty"`
	Ville             string `json:"ville,omitempty"`
	VilleNaissance    string `json:"ville_naissance,omitempty"`
	LieuNaissance     string `json:"lieu_naissance,omitempty"`
	Pays              string `json:"pays,omitempty"`
	Region            string `json:"region,omitempty"`
	Departement       string `json:"departement,omitempty"`
	
	// Unique IDs
	NIR     string `json:"nir,omitempty"`
	IBAN    string `json:"iban,omitempty"`
	BIC     string `json:"bic,omitempty"`
	SIRET   string `json:"siret,omitempty"`
	SIREN   string `json:"siren,omitempty"`
	
	// Vehicle
	VINPlaque       string `json:"vin_plaque,omitempty"`
	Immatriculation string `json:"immatriculation,omitempty"`
	NumeroSerie     string `json:"numero_serie,omitempty"`
	Marque          string `json:"marque,omitempty"`
	Modele          string `json:"modele,omitempty"`
	
	// Professional
	Societe    string `json:"societe,omitempty"`
	Profession string `json:"profession,omitempty"`
	Fonction   string `json:"fonction,omitempty"`
	
	// Gaming / FiveM
	SteamID       string `json:"steam_id,omitempty"`
	FiveMLicense  string `json:"fivem_license,omitempty"`
	FiveMLicense2 string `json:"fivem_license2,omitempty"`
	FiveMID       string `json:"fivem_id,omitempty"`
	XboxLiveID    string `json:"xbox_live_id,omitempty"`
	LiveID        string `json:"live_id,omitempty"`
	DiscordID     string `json:"discord_id,omitempty"`
	
	// Metadata
	Sources    []string `json:"_sources,omitempty"`
	Confidence int      `json:"_confidence,omitempty"`
}

// SearchResults contains the search response data
type SearchResults struct {
	Results []Profile `json:"results"`
}

// SearchRequest contains all possible search parameters
type SearchRequest struct {
	// Identity
	NomFamille     string `json:"nom_famille,omitempty"`
	Prenom         string `json:"prenom,omitempty"`
	NomNaissance   string `json:"nom_naissance,omitempty"`
	NomAffichage   string `json:"nom_affichage,omitempty"`
	NomUtilisateur string `json:"nom_utilisateur,omitempty"`
	DateNaissance  string `json:"date_naissance,omitempty"`
	AnneeNaissance string `json:"annee_naissance,omitempty"`
	JourNaissance  int    `json:"jour_naissance,omitempty"`
	MoisNaissance  int    `json:"mois_naissance,omitempty"`
	Genre          string `json:"genre,omitempty"`
	Civilite       string `json:"civilite,omitempty"`
	
	// Contact
	Email     string `json:"email,omitempty"`
	Telephone string `json:"telephone,omitempty"`
	Mobile    string `json:"mobile,omitempty"`
	AdresseIP string `json:"adresse_ip,omitempty"`
	
	// Address
	Adresse           string `json:"adresse,omitempty"`
	ComplementAdresse string `json:"complement_adresse,omitempty"`
	CodePostal        string `json:"code_postal,omitempty"`
	Ville             string `json:"ville,omitempty"`
	VilleNaissance    string `json:"ville_naissance,omitempty"`
	LieuNaissance     string `json:"lieu_naissance,omitempty"`
	Pays              string `json:"pays,omitempty"`
	Region            string `json:"region,omitempty"`
	Departement       string `json:"departement,omitempty"`
	
	// Unique IDs
	NIR   string `json:"nir,omitempty"`
	IBAN  string `json:"iban,omitempty"`
	BIC   string `json:"bic,omitempty"`
	SIRET string `json:"siret,omitempty"`
	SIREN string `json:"siren,omitempty"`
	
	// Vehicle
	VINPlaque       string `json:"vin_plaque,omitempty"`
	Immatriculation string `json:"immatriculation,omitempty"`
	NumeroSerie     string `json:"numero_serie,omitempty"`
	Marque          string `json:"marque,omitempty"`
	Modele          string `json:"modele,omitempty"`
	
	// Professional
	Societe    string `json:"societe,omitempty"`
	Profession string `json:"profession,omitempty"`
	Fonction   string `json:"fonction,omitempty"`
	
	// Gaming
	SteamID       string `json:"steam_id,omitempty"`
	FiveMLicense  string `json:"fivem_license,omitempty"`
	FiveMLicense2 string `json:"fivem_license2,omitempty"`
	FiveMID       string `json:"fivem_id,omitempty"`
	XboxLiveID    string `json:"xbox_live_id,omitempty"`
	LiveID        string `json:"live_id,omitempty"`
	DiscordID     string `json:"discord_id,omitempty"`
	
	// Options
	Page    int  `json:"page,omitempty"`
	PerPage int  `json:"per_page,omitempty"`
	Flexible bool `json:"flexible,omitempty"`
}

// AccountInfo contains account and quota information
type AccountInfo struct {
	Plan              string `json:"plan"`
	DailyQuota        int    `json:"daily_quota"`
	DailyUsed         int    `json:"daily_used"`
	DailyRemaining    int    `json:"daily_remaining"`
	TotalRequests     int    `json:"total_requests"`
	ResultsPerQuery   int    `json:"results_per_query"`
	PaginationEnabled bool   `json:"pagination_enabled"`
}

// UsageLog represents a single API usage log entry
type UsageLog struct {
	Timestamp   time.Time `json:"timestamp"`
	Endpoint    string    `json:"endpoint"`
	Query       string    `json:"query"`
	ResultsCount int      `json:"results_count"`
	TookMs      int       `json:"took_ms"`
}

// UsageResponse contains usage history
type UsageResponse struct {
	Logs []UsageLog `json:"logs"`
}

// HealthResponse contains API health status
type HealthResponse struct {
	Status string `json:"status"`
}