package waha

import "time"

type SessionWebHook struct {
	Event          string    `json:"event"`
	Timestamp      time.Time `json:"timestamp"`
	SessionId      string    `json:"sessionId"`
	IdempotencyKey string    `json:"idempotencyKey"`
	DeliveryId     string    `json:"deliveryId"`
	Data           struct {
		SessionId string `json:"sessionId"`
		Status    string `json:"status"`
	} `json:"data"`
}

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

type GroupResponse struct {
	Id              string      `json:"id"`
	Name            string      `json:"name"`
	LinkedParentJID interface{} `json:"linkedParentJID"`
}
