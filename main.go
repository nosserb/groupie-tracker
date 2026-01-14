package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// Artist structure containing band/artist information
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

// Location structure
type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

// Date structure
type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

// Relation structure
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// Coordinate structure for geocoding
type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
}

// GeocodeRequest structure for API requests
type GeocodeRequest struct {
	Address string `json:"address"`
}

// NominatimResponse structure for OpenStreetMap API response
type NominatimResponse struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

// API response structures
type LocationsIndex struct {
	Index []Location `json:"index"`
}

type DatesIndex struct {
	Index []Date `json:"index"`
}

type RelationIndex struct {
	Index []Relation `json:"index"`
}

// Combined data for template
type ArtistData struct {
	Artist    Artist
	Locations []string
	Dates     []string
	Relations map[string][]string
}

var (
	artists   []Artist
	locations []Location
	dates     []Date
	relations []Relation
)

const apiURL = "https://groupietrackers.herokuapp.com/api"
const nominatimURL = "https://nominatim.openstreetmap.org/search"

// Cache pour éviter les requêtes répétées
var geocodeCache = make(map[string]Coordinate)

func main() {
	// Fetch data from API
	if err := fetchAPIData(); err != nil {
		log.Printf("Warning: Could not fetch API data: %v", err)
	}

	// Serve static files first
	http.HandleFunc("/public/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.HandleFunc("/index.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.css")
	})
	http.HandleFunc("/artists.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "artists.css")
	})
	http.HandleFunc("/map.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "map.css")
	})
	http.HandleFunc("/modal.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "modal.css")
	})
	http.HandleFunc("/credits.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "credits.css")
	})

	// Routes
	http.HandleFunc("/artists.html", artistsHandler)
	http.HandleFunc("/credits.html", creditsHandler)
	http.HandleFunc("/map.html", mapHandler)
	http.HandleFunc("/artist/", artistDetailHandler)
	http.HandleFunc("/api/artists", apiArtistsHandler)
	http.HandleFunc("/api/geocode", geocodeHandler)
	http.HandleFunc("/", indexHandler)

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// Fetch all data from the API
func fetchAPIData() error {
	// Fetch artists
	if err := fetchJSON(apiURL+"/artists", &artists); err != nil {
		return fmt.Errorf("failed to fetch artists: %w", err)
	}

	// Fetch locations
	var locIndex LocationsIndex
	if err := fetchJSON(apiURL+"/locations", &locIndex); err != nil {
		return fmt.Errorf("failed to fetch locations: %w", err)
	}
	locations = locIndex.Index

	// Fetch dates
	var dateIndex DatesIndex
	if err := fetchJSON(apiURL+"/dates", &dateIndex); err != nil {
		return fmt.Errorf("failed to fetch dates: %w", err)
	}
	dates = dateIndex.Index

	// Fetch relations
	var relIndex RelationIndex
	if err := fetchJSON(apiURL+"/relation", &relIndex); err != nil {
		return fmt.Errorf("failed to fetch relations: %w", err)
	}
	relations = relIndex.Index

	log.Printf("Successfully fetched %d artists from API", len(artists))
	return nil
}

// Generic function to fetch JSON from URL
func fetchJSON(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

// Index page handler
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		http.NotFound(w, r)
		return
	}

	// Create a FuncMap with custom functions
	funcMap := template.FuncMap{
		"mod": func(i, j int) int {
			return i % j
		},
		"ge": func(a, b int) bool {
			return a >= b
		},
		"len": func(s []Artist) int {
			return len(s)
		},
		"index": func(s []Artist, i int) Artist {
			return s[i]
		},
	}

	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	// Get top 3 artists for the podium
	topArtists := artists
	if len(topArtists) > 3 {
		topArtists = artists[:3]
	}

	data := struct {
		Artists []Artist
	}{
		Artists: topArtists,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

// Artists page handler
func artistsHandler(w http.ResponseWriter, r *http.Request) {
	// Create a FuncMap with custom functions
	funcMap := template.FuncMap{
		"mod": func(i, j int) int {
			return i % j
		},
	}

	tmpl, err := template.New("artists.html").Funcs(funcMap).ParseFiles("artists.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	data := struct {
		Artists []Artist
	}{
		Artists: artists,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

// Artist detail handler
func artistDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/artist/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 || id > len(artists) {
		http.NotFound(w, r)
		return
	}

	artist := artists[id-1]
	var artistLocs []string
	var artistDates []string
	var artistRels map[string][]string

	// Find corresponding data
	if id <= len(locations) {
		artistLocs = locations[id-1].Locations
	}
	if id <= len(dates) {
		artistDates = dates[id-1].Dates
	}
	if id <= len(relations) {
		artistRels = relations[id-1].DatesLocations
	}

	data := ArtistData{
		Artist:    artist,
		Locations: artistLocs,
		Dates:     artistDates,
		Relations: artistRels,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// API endpoint to get artists as JSON
func apiArtistsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artists)
}

// Credits page handler
func creditsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "credits.html")
}

// Geocode handler - converts address to coordinates
func geocodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, `{"error":"address parameter required"}`, http.StatusBadRequest)
		return
	}

	// Check cache first
	if coord, exists := geocodeCache[address]; exists {
		json.NewEncoder(w).Encode(coord)
		return
	}

	// Query Nominatim API
	coord, err := geocodeAddress(address)
	if err != nil {
		log.Printf("Geocoding error for %s: %v", address, err)
		http.Error(w, `{"error":"geocoding failed"}`, http.StatusInternalServerError)
		return
	}

	// Cache the result
	geocodeCache[address] = coord
	json.NewEncoder(w).Encode(coord)
}

// geocodeAddress converts an address to coordinates using Nominatim API
func geocodeAddress(address string) (Coordinate, error) {
	baseURL, _ := url.Parse(nominatimURL)
	params := url.Values{}
	params.Add("q", address)
	params.Add("format", "json")
	params.Add("limit", "1")
	baseURL.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return Coordinate{}, fmt.Errorf("request creation failed: %w", err)
	}

	// Set User-Agent to comply with Nominatim policy
	req.Header.Set("User-Agent", "Groupie-Tracker/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		return Coordinate{}, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var results []NominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return Coordinate{}, fmt.Errorf("JSON decode failed: %w", err)
	}

	if len(results) == 0 {
		return Coordinate{}, fmt.Errorf("no results found for address: %s", address)
	}

	// Parse coordinates
	lat, err := strconv.ParseFloat(results[0].Lat, 64)
	if err != nil {
		return Coordinate{}, fmt.Errorf("invalid latitude: %w", err)
	}

	lon, err := strconv.ParseFloat(results[0].Lon, 64)
	if err != nil {
		return Coordinate{}, fmt.Errorf("invalid longitude: %w", err)
	}

	return Coordinate{
		Latitude:  lat,
		Longitude: lon,
		Address:   address,
	}, nil
}

// Map page handler
func mapHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "map.html")
}
