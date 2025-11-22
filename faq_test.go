package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHandleFAQ_MovieQuestion(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		expectAnswer  bool
		expectedTitle string
	}{
		{
			name:          "flat earth with turtle",
			query:         "movie about flat earth carried by turtle",
			expectAnswer:  true,
			expectedTitle: "The Color of Magic",
		},
		{
			name:          "flat earth with elephants",
			query:         "flat earth elephants movie",
			expectAnswer:  true,
			expectedTitle: "The Color of Magic",
		},
		{
			name:          "mage flat earth turtle",
			query:         "mage flat earth turtle",
			expectAnswer:  true,
			expectedTitle: "The Color of Magic",
		},
		{
			name:         "unrelated question",
			query:        "what is the weather today",
			expectAnswer: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/faq?q="+url.QueryEscape(tt.query), nil)
			w := httptest.NewRecorder()

			HandleFAQ(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", w.Code)
			}

			var response FAQResponse
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if tt.expectAnswer {
				if response.Answer == "No answer found for this question." {
					t.Errorf("expected an answer, got none")
				}
				if tt.expectedTitle != "" {
					contains := false
					for _, word := range []string{"Color", "Magic", "Discworld"} {
						if strings.Contains(response.Answer, word) {
							contains = true
							break
						}
					}
					if !contains {
						t.Errorf("expected answer to contain references to Discworld/Color of Magic, got: %s", response.Answer)
					}
				}
			} else {
				if response.Answer != "No answer found for this question." {
					t.Errorf("expected no answer, got: %s", response.Answer)
				}
			}
		})
	}
}

func TestHandleFAQ_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/faq?q=test", nil)
	w := httptest.NewRecorder()

	HandleFAQ(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}
