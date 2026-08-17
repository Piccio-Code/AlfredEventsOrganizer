package waha

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (c *WahaClient) setAPIHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)
}

func (c *WahaClient) GetWahaSession() (session SessionResponse, err error) {
	req, err := http.NewRequest("GET", c.baseSessionURL, nil)
	if err != nil {
		return SessionResponse{}, err
	}
	c.setAPIHeaders(req)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return SessionResponse{}, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&session); err != nil {
		return SessionResponse{}, err
	}

	return session, nil
}

func (c *WahaClient) RestartSession() error {
	req, err := http.NewRequest("GET", c.baseSessionURL+"/stop", nil)
	if err != nil {
		return err
	}
	c.setAPIHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	req, err = http.NewRequest("POST", c.baseSessionURL+"/start", nil)
	if err != nil {
		return err
	}
	c.setAPIHeaders(req)

	resp, err = c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *WahaClient) GetSessionCode() (sessionCode ParingCodeResponse, err error) {
	postBody, err := json.Marshal(map[string]string{
		"phoneNumber": c.phoneNumber,
	})
	if err != nil {
		return ParingCodeResponse{}, err
	}

	req, err := http.NewRequest(
		"POST",
		c.baseSessionURL+"/pairing-code",
		bytes.NewBuffer(postBody),
	)
	if err != nil {
		return ParingCodeResponse{}, err
	}

	c.setAPIHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ParingCodeResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
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
		return ParingCodeResponse{}, err
	}

	return sessionCode, nil
}

func (c *WahaClient) GetGroupsList() (groups []GroupResponse, err error) {
	req, err := http.NewRequest("GET", c.baseSessionURL+"/groups", nil)

	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"OpenWA returned status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	if err := json.Unmarshal(body, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

func (c *WahaClient) GetGroupDetail(groupID string) (groupDetail GroupDetailResponse, err error) {
	req, err := http.NewRequest("GET", c.baseSessionURL+"/groups/"+groupID, nil)

	if err != nil {
		return GroupDetailResponse{}, err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return GroupDetailResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GroupDetailResponse{}, err
	}

	if resp.StatusCode == 404 {
		return GroupDetailResponse{}, errors.New("the group does not exsist")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return GroupDetailResponse{}, fmt.Errorf(
			"OpenWA returned status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	if err := json.Unmarshal(body, &groupDetail); err != nil {
		return GroupDetailResponse{}, err
	}

	if len(groupDetail.Participants) == 0 {
		return GroupDetailResponse{}, errors.New("the group does not exsist")
	}

	for i := 0; i < len(groupDetail.Participants); i++ {
		participant := groupDetail.Participants[i]

		contactInfo, err := c.GetContactInfo(participant.Id)
		if err != nil {
			return GroupDetailResponse{}, err
		}

		if contactInfo.Name == "" {
			participant.Name = contactInfo.Number
		} else {
			participant.Name = contactInfo.Name
		}

		groupDetail.Participants[i] = participant
	}

	return groupDetail, nil
}

func (c *WahaClient) GetContactInfo(contactId string) (contactInfo ContactInfo, err error) {
	req, err := http.NewRequest("GET", c.baseSessionURL+"/contacts/"+contactId, nil)

	if err != nil {
		return ContactInfo{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ContactInfo{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ContactInfo{}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ContactInfo{}, fmt.Errorf(
			"OpenWA returned status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	if err := json.Unmarshal(body, &contactInfo); err != nil {
		return ContactInfo{}, err
	}

	return contactInfo, nil
}
