package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type SessionsWaha struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Config struct {
		Webhooks []struct {
			Url    string   `json:"url"`
			Events []string `json:"events"`
			Hmac   struct {
				Key interface{} `json:"key"`
			} `json:"hmac"`
			Retries struct {
				DelaySeconds int    `json:"delaySeconds"`
				Attempts     int    `json:"attempts"`
				Policy       string `json:"policy"`
			} `json:"retries"`
			CustomHeaders interface{} `json:"customHeaders"`
		} `json:"webhooks"`
		Metadata struct {
		} `json:"metadata"`
		Webjs struct {
			TagsEventsOn bool `json:"tagsEventsOn"`
		} `json:"webjs"`
		Client interface{} `json:"client"`
	} `json:"config"`
	Me struct {
		Id               string      `json:"id"`
		Lid              string      `json:"lid"`
		PushName         string      `json:"pushName"`
		ReachoutTimelock interface{} `json:"reachoutTimelock"`
	} `json:"me"`
	Presence   interface{} `json:"presence"`
	Timestamps struct {
	} `json:"timestamps"`
	Engine struct {
		Engine      string `json:"engine"`
		WWebVersion string `json:"WWebVersion"`
		State       string `json:"state"`
	} `json:"engine"`
}

type SessionCode struct {
	Code string `json:"code"`
}

func (s *Server) GetWahaSession() (session SessionsWaha, err error) {
	httpClient := http.Client{Timeout: time.Duration(5) * time.Second}
	req, err := http.NewRequest("GET", "https://waha.coursetracker.it/api/sessions/default", nil)
	if err != nil {
		s.infoLog.Println(err)
		return SessionsWaha{}, err
	}
	req.Header.Add("Accept", `application/json`)
	req.Header.Add("X-Api-Key", os.Getenv("WAHA_API_KEY"))

	resp, err := httpClient.Do(req)
	if err != nil {
		s.infoLog.Println(err)
		return SessionsWaha{}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return SessionsWaha{}, errors.New("bad Request")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		s.infoLog.Println(err)
		return SessionsWaha{}, err
	}

	return session, nil
}

func (s *Server) RestartSession() error {
	req, err := http.NewRequest(
		"POST",
		"https://waha.coursetracker.it/api/sessions/default/logout",
		nil,
	)

	if err != nil {
		return err
	}

	req.Header.Add("Accept", `application/json`)
	req.Header.Add("X-Api-Key", os.Getenv("WAHA_API_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	req, err = http.NewRequest(
		"POST",
		"https://waha.coursetracker.it/api/sessions/default/start",
		nil,
	)

	if err != nil {
		return err
	}

	req.Header.Add("Accept", `application/json`)
	req.Header.Add("X-Api-Key", os.Getenv("WAHA_API_KEY"))

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (s *Server) getSessionCode() (sessionCode SessionCode, err error) {
	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}

	postBody, err := json.Marshal(map[string]string{
		"phoneNumber": "393917031610",
		"method":      "",
	})
	if err != nil {
		s.infoLog.Println(err)
		return SessionCode{}, err
	}

	req, err := http.NewRequest(
		"POST",
		"https://waha.coursetracker.it/api/default/auth/request-code",
		bytes.NewBuffer(postBody),
	)
	if err != nil {
		s.infoLog.Println(err)
		return SessionCode{}, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", os.Getenv("WAHA_API_KEY"))

	resp, err := httpClient.Do(req)
	if err != nil {
		s.infoLog.Println(err)
		return SessionCode{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			s.infoLog.Printf(
				"WAHA request failed: status=%d, body=<failed to read: %v>",
				resp.StatusCode,
				err,
			)
		} else {
			s.infoLog.Printf(
				"WAHA request failed: status=%d, body=%s",
				resp.StatusCode,
				strings.TrimSpace(string(body)),
			)
		}

		return SessionCode{}, fmt.Errorf("WAHA returned status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&sessionCode); err != nil {
		s.infoLog.Println(err)
		return SessionCode{}, err
	}

	return sessionCode, nil
}
