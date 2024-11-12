package api

type Twin struct {
	TwinID    uint32 `json:"twin_id"`
	Relay     string `json:"relay"`
	AccountID string `json:"account_id"`
	PublicKey string `json:"public_key"`
}

type Spending struct {
	LastHour float64 `json:"last_hour"`
	LifeTime float64 `json:"life_time"`
}

type TwinFilter struct {
	TwinIDs   []uint32 `schema:"twin_id,omitempty"`
	Relays    []string `schema:"relay,omitempty"`
	AccountID *string  `schema:"account_id,omitempty"`
	PublicKey *string  `schema:"public_key,omitempty"`
}

type TwinSort struct {
	TwinId *bool
	Relay  *bool
}
