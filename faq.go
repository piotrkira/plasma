package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type FAQResponse struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// HandleFAQ provides answers to common questions
func HandleFAQ(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := strings.ToLower(r.URL.Query().Get("q"))
	
	var response FAQResponse
	
	// Check for movie question about flat earth carried by turtles/elephants
	if strings.Contains(query, "flat") && strings.Contains(query, "earth") && 
		(strings.Contains(query, "turtle") || strings.Contains(query, "elephant")) {
		response = FAQResponse{
			Question: "Movie about flat earth carried by turtles/elephants",
			Answer:   "You're thinking of 'The Color of Magic' (2008), a TV adaptation of Terry Pratchett's Discworld series. The Discworld is a flat disc resting on the backs of four giant elephants, which stand on the shell of an enormous turtle named Great A'Tuin swimming through space.",
		}
	} else {
		response = FAQResponse{
			Question: query,
			Answer:   "No answer found for this question.",
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
