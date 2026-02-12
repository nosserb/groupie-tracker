package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

// SearchResult structure for search suggestions
type SearchResult struct {
	Type  string `json:"type"` // "artist", "member", "location", "creationDate", "firstAlbum"
	Name  string `json:"name"`
	Value string `json:"value"`
}

// FilterRequest structure for filter API requests
type FilterRequest struct {
	CreationDateMin int      `json:"creationDateMin"`
	CreationDateMax int      `json:"creationDateMax"`
	FirstAlbumMin   int      `json:"firstAlbumMin"`
	FirstAlbumMax   int      `json:"firstAlbumMax"`
	MemberCounts    []int    `json:"memberCounts"`
	Locations       []string `json:"locations"`
}

// FilterResponse structure for filter API responses
type FilterResponse struct {
	Artists []Artist `json:"artists"`
	Total   int      `json:"total"`
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

// formatLocationName converts API location format (city-country) to readable format
func formatLocationName(location string) string {
	// Replace underscores with spaces
	location = strings.ReplaceAll(location, "_", " ")
	// Replace the last hyphen with a comma to separate city and country
	parts := strings.Split(location, "-")
	if len(parts) >= 2 {
		// Join all parts except the last with hyphen, then add comma and last part (country)
		city := strings.Join(parts[:len(parts)-1], "-")
		country := parts[len(parts)-1]
		return strings.TrimSpace(city) + ", " + strings.TrimSpace(country)
	}
	return location
}

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
	http.HandleFunc("/api/search", searchHandler)
	http.HandleFunc("/api/filter", filterHandler)
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

// API endpoint to get artists
func apiArtistsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artists)
}

// Credits page handler
func creditsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "credits.html")
}

// Geocode handler
func geocodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	address := r.URL.Query().Get("address")
	if address == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "address parameter required"})
		return
	}

	// Check cache first
	if coord, exists := geocodeCache[address]; exists {
		json.NewEncoder(w).Encode(coord)
		return
	}

	// Nominatim API
	coord, err := geocodeAddress(address)
	if err != nil {
		log.Printf("Geocoding error for %s: %v", address, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "geocoding failed"})
		return
	}

	// Cache the result
	geocodeCache[address] = coord
	json.NewEncoder(w).Encode(coord)
}

// geocodeAddress converts coordinates API
func geocodeAddress(address string) (Coordinate, error) {
	// location string
	formattedAddress := formatLocationName(address)

	baseURL, _ := url.Parse(nominatimURL)
	params := url.Values{}
	params.Add("q", formattedAddress)
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
	if err != nil {
		return Coordinate{}, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

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

// Search handler
func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	if query == "" {
		json.NewEncoder(w).Encode([]SearchResult{})
		return
	}

	results := []SearchResult{}
	seen := make(map[string]bool)

	// Search artists
	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), query) {
			key := "artist:" + artist.Name
			if !seen[key] {
				results = append(results, SearchResult{
					Type:  "artist",
					Name:  artist.Name,
					Value: artist.Name,
				})
				seen[key] = true
			}
		}
	}

	// Search members
	for _, artist := range artists {
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
				key := "member:" + member
				if !seen[key] {
					results = append(results, SearchResult{
						Type:  "member",
						Name:  member,
						Value: member,
					})
					seen[key] = true
				}
			}
		}
	}

	// Search locations
	for _, loc := range locations {
		for _, location := range loc.Locations {
			if strings.Contains(strings.ToLower(location), query) {
				key := "location:" + location
				if !seen[key] {
					results = append(results, SearchResult{
						Type:  "location",
						Name:  formatLocationName(location),
						Value: location,
					})
					seen[key] = true
				}
			}
		}
	}

	// Search creation dates
	for _, artist := range artists {
		dateStr := strconv.Itoa(artist.CreationDate)
		if strings.Contains(dateStr, query) {
			key := "creationDate:" + dateStr
			if !seen[key] {
				results = append(results, SearchResult{
					Type:  "creationDate",
					Name:  dateStr,
					Value: dateStr,
				})
				seen[key] = true
			}
		}
	}

	// Search first album dates
	for _, artist := range artists {
		if strings.Contains(artist.FirstAlbum, query) {
			key := "firstAlbum:" + artist.FirstAlbum
			if !seen[key] {
				results = append(results, SearchResult{
					Type:  "firstAlbum",
					Name:  artist.FirstAlbum,
					Value: artist.FirstAlbum,
				})
				seen[key] = true
			}
		}
	}

	// Limit results to 10
	if len(results) > 10 {
		results = results[:10]
	}

	json.NewEncoder(w).Encode(results)
}

// Filter handler for advanced filtering
func filterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters
	creationDateMin := 0
	creationDateMax := 9999
	firstAlbumMin := 0
	firstAlbumMax := 9999
	memberCountsStr := r.URL.Query().Get("memberCounts")
	locationsStr := r.URL.Query().Get("locations")

	// Parse creation date range
	if val := r.URL.Query().Get("creationDateMin"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			creationDateMin = parsed
		}
	}
	if val := r.URL.Query().Get("creationDateMax"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			creationDateMax = parsed
		}
	}

	// Parse first album date range
	if val := r.URL.Query().Get("firstAlbumMin"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			firstAlbumMin = parsed
		}
	}
	if val := r.URL.Query().Get("firstAlbumMax"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			firstAlbumMax = parsed
		}
	}

	// Parse member counts filter
	var memberCounts []int
	if memberCountsStr != "" {
		parts := strings.Split(memberCountsStr, ",")
		for _, part := range parts {
			if count, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
				memberCounts = append(memberCounts, count)
			}
		}
	}

	// Parse locations filter
	var filterLocations []string
	if locationsStr != "" {
		filterLocations = strings.Split(locationsStr, ",")
		for i := range filterLocations {
			filterLocations[i] = strings.TrimSpace(filterLocations[i])
		}
	}

	// Filter artists
	filteredArtists := filterArtists(
		creationDateMin,
		creationDateMax,
		firstAlbumMin,
		firstAlbumMax,
		memberCounts,
		filterLocations,
	)

	response := FilterResponse{
		Artists: filteredArtists,
		Total:   len(filteredArtists),
	}

	json.NewEncoder(w).Encode(response)
}

// filterArtists applies all filters and returns matching artists
func filterArtists(creationDateMin, creationDateMax, firstAlbumMin, firstAlbumMax int, memberCounts []int, filterLocations []string) []Artist {
	var result []Artist

	for i, artist := range artists {
		// Check creation date range
		if artist.CreationDate < creationDateMin || artist.CreationDate > creationDateMax {
			continue
		}

		// Check first album year range
		albumYear := extractYear(artist.FirstAlbum)
		if albumYear < firstAlbumMin || albumYear > firstAlbumMax {
			continue
		}

		// Check member count filter
		if len(memberCounts) > 0 {
			memberCount := len(artist.Members)
			found := false
			for _, count := range memberCounts {
				if count == 7 && memberCount >= 7 {
					found = true
					break
				} else if count == memberCount {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Check locations filter
		if len(filterLocations) > 0 {
			artistLocs := []string{}
			if i+1 <= len(locations) {
				artistLocs = locations[i].Locations
			}
			found := false
			for _, filterLoc := range filterLocations {
				for _, artistLoc := range artistLocs {
					if strings.Contains(strings.ToLower(artistLoc), strings.ToLower(filterLoc)) {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found && len(filterLocations) > 0 {
				continue
			}
		}

		result = append(result, artist)
	}

	return result
}

// extractYear extracts the year from a date string (e.g., "06-04-2009" -> 2009)
func extractYear(dateStr string) int {
	parts := strings.Split(dateStr, "-")
	if len(parts) > 0 {
		if year, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
			return year
		}
	}
	return 0
}

// Map page handler
func mapHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "map.html")
}
