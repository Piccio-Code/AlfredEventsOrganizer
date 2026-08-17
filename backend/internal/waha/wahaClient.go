package waha

import (
	"net/http"
	"os"
	"strings"
	"time"
)

type WahaClient struct {
	httpClient  *http.Client
	baseURL     string
	apiKey      string
	session     string
	sessionURL  string
	apiURL      string
	phoneNumber string
}

func NewWahaClient() *WahaClient {
	// WAHA_BASE_URL e' il solo host, il prefisso /api lo aggiungiamo qui perche'
	// e' fisso per tutte le rotte.
	baseURL := strings.TrimSuffix(os.Getenv("WAHA_BASE_URL"), "/")
	if !strings.HasSuffix(baseURL, "/api") {
		baseURL += "/api"
	}

	session := os.Getenv("WAHA_SESSION")

	return &WahaClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    baseURL,
		apiKey:     os.Getenv("WAHA_X_API_KEY"),
		session:    session,
		// Le API di gestione della sessione stanno sotto /sessions/{session},
		// quelle di WhatsApp direttamente sotto /{session}.
		sessionURL:  baseURL + "/sessions/" + session,
		apiURL:      baseURL + "/" + session,
		phoneNumber: os.Getenv("WAHA_PHONE_NUMBER"),
	}
}
