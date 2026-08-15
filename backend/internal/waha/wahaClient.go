package waha

import (
	"net/http"
	"os"
	"time"
)

type WahaClient struct {
	httpClient     *http.Client
	baseURL        string
	apiKey         string
	baseSessionURL string
	phoneNumber    string
}

func NewWahaClient() *WahaClient {
	baseURL := os.Getenv("OPENWA_BASE_URL")
	return &WahaClient{
		httpClient:     &http.Client{Timeout: 10 * time.Second},
		baseURL:        baseURL,
		apiKey:         os.Getenv("OPENWA_X_API_KEY"),
		baseSessionURL: baseURL + "/sessions/" + os.Getenv("OPENWA_ALFRED_ID"),
		phoneNumber:    os.Getenv("OPENWA_PHONE_NUMBER"),
	}
}
