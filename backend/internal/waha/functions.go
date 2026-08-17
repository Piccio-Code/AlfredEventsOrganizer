package waha

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var GroupNotFoundError = errors.New("the group does not exsist")

type StatusError struct {
	StatusCode int
	Body       string
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("WAHA returned status %d: %s", e.StatusCode, e.Body)
}

func (c *WahaClient) doJSON(method, endpoint string, body any, out any) error {
	var payload io.Reader

	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return err
		}

		payload = bytes.NewReader(encoded)
	}

	req, err := http.NewRequest(method, endpoint, payload)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Api-Key", c.apiKey)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return &StatusError{StatusCode: res.StatusCode, Body: string(raw)}
	}

	if out == nil {
		return nil
	}

	return json.Unmarshal(raw, out)
}

func (c *WahaClient) GetWahaSession() (session SessionResponse, err error) {
	var res newSessionResponse

	if err := c.doJSON(http.MethodGet, c.sessionURL, nil, &res); err != nil {
		return SessionResponse{}, err
	}

	return SessionResponse{
		Name:     res.Name,
		Status:   res.Status,
		Phone:    res.Me.Id,
		PushName: res.Me.PushName,
	}, nil
}

func (c *WahaClient) RestartSession() error {
	return c.doJSON(http.MethodPost, c.sessionURL+"/restart", nil, nil)
}

func (c *WahaClient) GetSessionCode() (sessionCode ParingCodeResponse, err error) {
	var res newParingCode

	body := map[string]string{"phoneNumber": c.phoneNumber}

	if err := c.doJSON(http.MethodPost, c.apiURL+"/auth/request-code", body, &res); err != nil {
		return ParingCodeResponse{}, err
	}

	return ParingCodeResponse{Code: res.Code}, nil
}

func (c *WahaClient) GetGroupsList() (groups []GroupResponse, err error) {
	var res []newGroupResponse

	if err := c.doJSON(http.MethodGet, c.apiURL+"/groups", nil, &res); err != nil {
		return nil, err
	}

	groups = make([]GroupResponse, 0, len(res))

	for _, group := range res {
		name := group.Name
		if name == "" {
			name = group.GroupMetadata.Subject
		}

		groups = append(groups, GroupResponse{
			Id:   group.Id.Serialized,
			Name: name,
		})
	}

	return groups, nil
}

func (c *WahaClient) GetGroupDetail(groupID string) (groupDetail GroupDetailResponse, err error) {
	var res newGroupDetailResponse

	err = c.doJSON(http.MethodGet, c.apiURL+"/groups/"+url.QueryEscape(groupID), nil, &res)

	var statusErr *StatusError
	if errors.As(err, &statusErr) && statusErr.StatusCode == http.StatusNotFound {
		return GroupDetailResponse{}, GroupNotFoundError
	}

	if err != nil {
		return GroupDetailResponse{}, err
	}

	if len(res.GroupMetadata.Participants) == 0 {
		return GroupDetailResponse{}, GroupNotFoundError
	}

	name := res.Name
	if name == "" {
		name = res.GroupMetadata.Subject
	}

	groupDetail = GroupDetailResponse{
		Id:           res.Id.Serialized,
		Name:         name,
		Description:  res.GroupMetadata.Desc,
		Owner:        res.GroupMetadata.Owner.Serialized,
		IsReadOnly:   res.IsReadOnly,
		IsAnnounce:   res.GroupMetadata.Announce,
		Participants: make([]GroupParticipant, 0, len(res.GroupMetadata.Participants)),
	}

	for _, member := range res.GroupMetadata.Participants {
		participant := GroupParticipant{
			Id:           member.Id.Serialized,
			Number:       member.Id.User,
			IsAdmin:      member.IsAdmin,
			IsSuperAdmin: member.IsSuperAdmin,
		}

		contactInfo, err := c.GetContactInfo(participant.Id)
		if err != nil {
			return GroupDetailResponse{}, err
		}

		participant.Name = contactName(contactInfo, participant.Number)

		groupDetail.Participants = append(groupDetail.Participants, participant)
	}

	return groupDetail, nil
}

func (c *WahaClient) GetContactInfo(contactId string) (contactInfo ContactInfo, err error) {
	var res newContactInfo

	if err := c.doJSON(http.MethodGet, c.apiURL+"/contacts/"+url.QueryEscape(contactId), nil, &res); err != nil {
		return ContactInfo{}, err
	}

	return ContactInfo{
		Id:       res.Id,
		Number:   res.Number,
		Name:     res.Name,
		PushName: res.Pushname,
	}, nil
}

func contactName(contact ContactInfo, fallback string) string {
	if contact.Name != "" {
		return contact.Name
	}

	if contact.PushName != "" {
		return contact.PushName
	}

	if contact.Number != "" {
		return contact.Number
	}

	return fallback
}
