package db

type Farm struct {
	FarmID      uint64 `gorm:"primaryKey;autoIncrement" json:"farm_id"` // Primary key
	FarmName    string `gorm:"size:100;not null;unique" json:"farm_name"`
	TwinID      uint64 `json:"twin_id" gorm:"not null"`
	Dedicated   bool   `json:"dedicated"`
	FarmFreeIps uint64 `json:"farm_free_ips"`

	Nodes []Node `gorm:"foreignKey:farm_id;constraint:OnDelete:CASCADE" json:"nodes"`
}

type Node struct {
	NodeID      uint64      `gorm:"primaryKey;autoIncrement" json:"node_id"`
	FarmID      uint64      `gorm:"not null" json:"farm_id"`
	TwinID      uint64      `json:"twin_id" gorm:"not null"`
	Features    []string    `gorm:"type:jsonb;serializer:json"`
	PriceUsd    float64     `json:"price_usd"`
	ExtraFee    uint64      `json:"extra_fee"`
	Status      string      `gorm:"size:50" json:"status"`
	Healthy     bool        `json:"healthy"`
	Dedicated   bool        `json:"dedicated"`
	Rented      bool        `json:"rented"`
	Rentable    bool        `json:"rentable"`
	Uptime      Uptime      `json:"uptime"`
	Consumption Consumption `json:"consumption"`
}

type NodeFilter struct {
	NodeID  *uint64 `form:"node_id"`
	FarmID  *uint64 `form:"farm_id"`
	TwinID  *uint64 `form:"twin_id"`
	Status  string  `form:"status"`
	Healthy bool    `form:"healthy"`
}

type FarmFilter struct {
	FarmName *string `form:"farm_name"`
	FarmID   *uint64 `form:"farm_id"`
	TwinID   *uint64 `form:"twin_id"`
}

// Limit used for pagination
type Limit struct {
	Size uint64 `form:"size"`
	Page uint64 `form:"page"`
}

type Uptime int64

type Consumption string

// DefaultLimit returns the default values for the pagination
func DefaultLimit() Limit {
	return Limit{
		Size: 50,
		Page: 1,
	}
}
