package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type SessionResponse struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Phone       string    `json:"phone"`
	PushName    string    `json:"pushName"`
	ConnectedAt time.Time `json:"connectedAt"`
	LastActive  time.Time `json:"lastActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	LastError   string    `json:"lastError"`
	Restriction struct {
		Kind      string    `json:"kind"`
		Code      string    `json:"code"`
		ExpiresAt time.Time `json:"expiresAt"`
	} `json:"restriction"`
	EngineLoaded bool `json:"engineLoaded"`
}

type ParingCodeResponse struct {
	PairingCode string `json:"pairingCode"`
	Status      string `json:"status"`
}

func (s *Server) GetWahaSession() (session SessionResponse, err error) {
	url := os.Getenv("OPENWA_BASE_URL") + "/api/sessions/" + os.Getenv("OPENWA_ALFRED_ID")
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-API-Key", os.Getenv("OPENWA_X_API_KEY"))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&session); err != nil {
		s.infoLog.Println(err)
		return SessionResponse{}, err
	}

	return session, nil
}

func (s *Server) RestartSession() error {
	req, err := http.NewRequest(
		"GET",
		"https://waha.coursetracker.it/api/sessions/"+os.Getenv("OPENWA_ALFRED_ID")+"/stop",
		nil,
	)

	if err != nil {
		return err
	}

	req.Header.Add("Accept", `application/json`)
	req.Header.Add("X-Api-Key", os.Getenv("OPENWA_X_API_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	req, err = http.NewRequest(
		"POST",
		"https://waha.coursetracker.it/api/sessions/"+os.Getenv("OPENWA_ALFRED_ID")+"/start",
		nil,
	)

	if err != nil {
		return err
	}

	req.Header.Add("Accept", `application/json`)
	req.Header.Add("X-Api-Key", os.Getenv("OPENWA_X_API_KEY"))

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (s *Server) GetSessionCode() (sessionCode ParingCodeResponse, err error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	postBody, err := json.Marshal(map[string]string{
		"phoneNumber": os.Getenv("OPENWA_PHONE_NUMBER"),
	})
	if err != nil {
		s.infoLog.Println(err)
		return ParingCodeResponse{}, err
	}

	req, err := http.NewRequest(
		"POST",
		os.Getenv("OPENWA_BASE_URL")+"/sessions/"+os.Getenv("OPENWA_ALFRED_ID")+"/pairing-code",
		bytes.NewBuffer(postBody),
	)
	if err != nil {
		s.infoLog.Println(err)
		return ParingCodeResponse{}, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", os.Getenv("OPENWA_X_API_KEY"))

	resp, err := httpClient.Do(req)
	if err != nil {
		s.infoLog.Println(err)
		return ParingCodeResponse{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.infoLog.Println(err)
		return ParingCodeResponse{}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ParingCodeResponse{}, fmt.Errorf(
			"OpenWA returned status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	if err := json.Unmarshal(body, &sessionCode); err != nil {
		s.infoLog.Println(err)
		return ParingCodeResponse{}, err
	}

	return sessionCode, nil
}
