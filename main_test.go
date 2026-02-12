package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestFetchJSON tests the fetchJSON function with a mock server
func TestFetchJSON(t *testing.T) {
	// Create a mock server
	mockData := []Artist{
		{
			ID:           1,
			Name:         "Test Band",
			CreationDate: 2000,
			Members:      []string{"Member 1", "Member 2"},
			FirstAlbum:   "01-01-2001",
			Image:        "http://example.com/image.jpg",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockData)
	}))
	defer server.Close()

	// Test fetchJSON
	var result []Artist
	err := fetchJSON(server.URL, &result)

	if err != nil {
		t.Fatalf("fetchJSON failed: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 artist, got %d", len(result))
	}

	if result[0].Name != "Test Band" {
		t.Errorf("Expected artist name 'Test Band', got '%s'", result[0].Name)
	}
}

// TestIndexHandler tests the index handler
func TestIndexHandler(t *testing.T) {
	// Mock artists data
	artists = []Artist{
		{ID: 1, Name: "Artist 1", CreationDate: 2000},
		{ID: 2, Name: "Artist 2", CreationDate: 2001},
		{ID: 3, Name: "Artist 3", CreationDate: 2002},
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(indexHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

// TestArtistsHandler tests the artists page handler
func TestArtistsHandler(t *testing.T) {
	// Mock artists data
	artists = []Artist{
		{ID: 1, Name: "Artist 1", CreationDate: 2000},
		{ID: 2, Name: "Artist 2", CreationDate: 2001},
	}

	req, err := http.NewRequest("GET", "/artists", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(artistsHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

// TestArtistDetailHandler tests the artist detail endpoint
func TestArtistDetailHandler(t *testing.T) {
	// Mock data
	artists = []Artist{
		{
			ID:           1,
			Name:         "Test Artist",
			CreationDate: 2000,
			Members:      []string{"Member 1"},
			FirstAlbum:   "01-01-2001",
			Image:        "http://example.com/image.jpg",
		},
	}
	locations = []Location{
		{ID: 1, Locations: []string{"Location 1"}},
	}
	dates = []Date{
		{ID: 1, Dates: []string{"01-01-2020"}},
	}
	relations = []Relation{
		{ID: 1, DatesLocations: map[string][]string{"Location 1": {"01-01-2020"}}},
	}

	req, err := http.NewRequest("GET", "/artist/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(artistDetailHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if response is valid JSON
	var result ArtistData
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Errorf("Response is not valid JSON: %v", err)
	}

	if result.Artist.Name != "Test Artist" {
		t.Errorf("Expected artist name 'Test Artist', got '%s'", result.Artist.Name)
	}
}

// TestArtistDetailHandler_InvalidID tests invalid artist ID
func TestArtistDetailHandler_InvalidID(t *testing.T) {
	artists = []Artist{
		{ID: 1, Name: "Test Artist"},
	}

	req, err := http.NewRequest("GET", "/artist/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(artistDetailHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler should return 404 for invalid ID: got %v want %v", status, http.StatusNotFound)
	}
}

// TestAPIArtistsHandler tests the API endpoint
func TestAPIArtistsHandler(t *testing.T) {
	artists = []Artist{
		{ID: 1, Name: "Artist 1", CreationDate: 2000},
		{ID: 2, Name: "Artist 2", CreationDate: 2001},
	}

	req, err := http.NewRequest("GET", "/api/artists", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apiArtistsHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check Content-Type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Verify JSON response
	var result []Artist
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Errorf("Response is not valid JSON: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 artists in response, got %d", len(result))
	}
}

// TestGeocodeHandler tests the geocode handler
func TestGeocodeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/geocode?address=Paris", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(geocodeHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check Content-Type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Verify JSON response
	var result Coordinate
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Errorf("Response is not valid JSON: %v", err)
	}

	if result.Address != "Paris" {
		t.Errorf("Expected address 'Paris', got '%s'", result.Address)
	}

	if result.Latitude == 0 || result.Longitude == 0 {
		t.Errorf("Coordinates not set properly: lat=%f, lon=%f", result.Latitude, result.Longitude)
	}
}

// TestGeocodeHandler_MissingAddress tests geocode handler with missing address
func TestGeocodeHandler_MissingAddress(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/geocode", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(geocodeHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler should return 400 for missing address: got %v want %v", status, http.StatusBadRequest)
	}
}

// TestExtractYear tests the extractYear function
func TestExtractYear(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{"01-01-2020", 2020},
		{"06-04-2009", 2009},
		{"15-06-1985", 1985},
		{"", 0},
		{"invalid", 0},
	}

	for _, tc := range testCases {
		result := extractYear(tc.input)
		if result != tc.expected {
			t.Errorf("extractYear(%s) = %d, want %d", tc.input, result, tc.expected)
		}
	}
}

// TestFilterArtists tests the filterArtists function
func TestFilterArtists(t *testing.T) {
	// Setup test data
	artists = []Artist{
		{
			ID:           1,
			Name:         "The Beatles",
			CreationDate: 1960,
			FirstAlbum:   "22-03-1963",
			Members:      []string{"John", "Paul", "George", "Ringo"},
		},
		{
			ID:           2,
			Name:         "Queen",
			CreationDate: 1970,
			FirstAlbum:   "13-02-1973",
			Members:      []string{"Freddie", "Brian", "Roger", "John"},
		},
		{
			ID:           3,
			Name:         "Pink Floyd",
			CreationDate: 1965,
			FirstAlbum:   "05-08-1967",
			Members:      []string{"Syd", "Roger", "Rick", "David", "Nick"},
		},
	}

	locations = []Location{
		{
			ID:        1,
			Locations: []string{"London_England", "Liverpool_England"},
		},
		{
			ID:        2,
			Locations: []string{"London_England"},
		},
		{
			ID:        3,
			Locations: []string{"London_England", "Cambridge_England"},
		},
	}

	testCases := []struct {
		name             string
		creationDateMin  int
		creationDateMax  int
		firstAlbumMin    int
		firstAlbumMax    int
		memberCounts     []int
		filterLocations  []string
		expectedCount    int
		expectedArtistID int
	}{
		{
			name:             "No filters - return all",
			creationDateMin:  0,
			creationDateMax:  9999,
			firstAlbumMin:    0,
			firstAlbumMax:    9999,
			memberCounts:     []int{},
			filterLocations:  []string{},
			expectedCount:    3,
			expectedArtistID: 0,
		},
		{
			name:             "Filter by creation date range",
			creationDateMin:  1960,
			creationDateMax:  1969,
			firstAlbumMin:    0,
			firstAlbumMax:    9999,
			memberCounts:     []int{},
			filterLocations:  []string{},
			expectedCount:    2,
			expectedArtistID: 1,
		},
		{
			name:             "Filter by member count",
			creationDateMin:  0,
			creationDateMax:  9999,
			firstAlbumMin:    0,
			firstAlbumMax:    9999,
			memberCounts:     []int{4},
			filterLocations:  []string{},
			expectedCount:    2,
			expectedArtistID: 1,
		},
		{
			name:             "Filter by location",
			creationDateMin:  0,
			creationDateMax:  9999,
			firstAlbumMin:    0,
			firstAlbumMax:    9999,
			memberCounts:     []int{},
			filterLocations:  []string{"London"},
			expectedCount:    3,
			expectedArtistID: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filterArtists(
				tc.creationDateMin,
				tc.creationDateMax,
				tc.firstAlbumMin,
				tc.firstAlbumMax,
				tc.memberCounts,
				tc.filterLocations,
			)

			if len(result) != tc.expectedCount {
				t.Errorf("Expected %d artists, got %d", tc.expectedCount, len(result))
			}

			if tc.expectedCount > 0 && len(result) > 0 && tc.expectedArtistID > 0 {
				if result[0].ID != tc.expectedArtistID {
					t.Errorf("Expected first artist ID %d, got %d", tc.expectedArtistID, result[0].ID)
				}
			}
		})
	}
}

// TestFilterHandler tests the filter API handler
func TestFilterHandler(t *testing.T) {
	// Setup test data
	artists = []Artist{
		{
			ID:           1,
			Name:         "Test Band",
			CreationDate: 2000,
			FirstAlbum:   "01-01-2001",
			Members:      []string{"Member 1", "Member 2"},
		},
	}
	locations = []Location{
		{
			ID:        1,
			Locations: []string{"Test_City"},
		},
	}

	req, err := http.NewRequest("GET", "/api/filter?creationDateMin=1990&creationDateMax=2010", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(filterHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check Content-Type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Verify JSON response
	var result FilterResponse
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Errorf("Response is not valid JSON: %v", err)
	}

	if result.Total != 1 {
		t.Errorf("Expected 1 artist in filter result, got %d", result.Total)
	}

	if len(result.Artists) != 1 {
		t.Errorf("Expected 1 artist, got %d", len(result.Artists))
	}
}
