package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	substrate "github.com/threefoldtech/tfchain/clients/tfchain-client-go"
)

type Account struct {
	TwinID    uint64 `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time
	PublicKey string `gorm:"type:text;not null;unique"` // ED25519 public key in the more standard base64 since we are moving from substarte echo system? (still SS58 can be used or plain base58 ,TBD)
	// Relations
	Farms []Farm `gorm:"foreignKey:TwinID;references:TwinID"`
	Nodes []Node `gorm:"foreignKey:TwinID;references:TwinID"`
}

type Farm struct {
	FarmID      uint64 `gorm:"primaryKey;autoIncrement" json:"farm_id"`
	FarmName    string `gorm:"size:40;not null;unique;check:farm_name <> ''" json:"farm_name"`
	TwinID      uint64 `json:"twin_id" gorm:"not null;check:twin_id > 0"`
	Dedicated   bool   `json:"dedicated"`
	FarmFreeIps uint64 `json:"farm_free_ips"`

	CreatedAt time.Time

	Nodes []Node `gorm:"foreignKey:farm_id;constraint:OnDelete:CASCADE" json:"nodes"`
}

type Node struct {
	NodeID uint64 `json:"node_id" gorm:"primaryKey;autoIncrement"`
	FarmID uint64 `json:"farm_id" gorm:"not null;check:farm_id> 0"`
	TwinID uint64 `json:"twin_id" gorm:"not null;check:twin_id > 0"`

	ZosVersion string `json:"zos_version" gorm:"not null"`
	NodeType   string `json:"node_type" gorm:"not null"`

	Location Location `json:"location" gorm:"not null;type:json"`

	// PublicConfig PublicConfig `json:"public_config" gorm:"type:json"`
	Resources    Resources `json:"resources" gorm:"not null;type:json"`
	Interface    Interface `json:"interface" gorm:"not null;type:json"`
	SecureBoot   bool
	Virtualized  bool
	SerialNumber string

	Uptime      Uptime      `json:"uptime"`
	Consumption Consumption `json:"consumption" gorm:"type:jsonb;serializer:json"`

	PriceUsd float64 `json:"price_usd"`
	Status   string  `json:"status"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Consumption []substrate.NruConsumption

type Uptime int64

type Interface struct {
	Name string `json:"name"`
	Mac  string `json:"mac"`
	IPs  string `json:"ips"`
}

// Value implements the Valuer interface for storing Interface in the database
func (i Interface) Value() (driver.Value, error) {
	bytes, err := json.Marshal(i)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Interface: %w", err)
	}
	return string(bytes), nil
}

// Scan implements the Scanner interface for retrieving Interface from the database
func (i *Interface) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid data type for Interface: %T", value)
	}

	if err := json.Unmarshal(bytes, i); err != nil {
		return fmt.Errorf("failed to unmarshal Interface: %w", err)
	}
	return nil
}

type Resources struct {
	HRU uint64 `json:"hru"`
	SRU uint64 `json:"sru"`
	CRU uint64 `json:"cru"`
	MRU uint64 `json:"mru"`
}

// Value implements the Valuer interface for storing Resources in the database
func (r Resources) Value() (driver.Value, error) {
	bytes, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resources: %w", err)
	}
	return string(bytes), nil
}

// Scan implements the Scanner interface for retrieving Resources from the database
func (r *Resources) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid data type for resources: %T", value)
	}

	if err := json.Unmarshal(bytes, r); err != nil {
		return fmt.Errorf("failed to unmarshal resources: %w", err)
	}
	return nil
}

type Location struct {
	Country   string `json:"country" gorm:"not null"`
	City      string `json:"city" gorm:"not null"`
	Longitude string `json:"longitude" gorm:"not null"`
	Latitude  string `json:"latitude" gorm:"not null"`
}

// Value implements the Valuer interface for storing Location in the database
func (l Location) Value() (driver.Value, error) {
	bytes, err := json.Marshal(l)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Location: %w", err)
	}
	return string(bytes), nil
}

// Scan implements the Scanner interface for retrieving Location from the database
func (l *Location) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid data type for Location: %T", value)
	}

	if err := json.Unmarshal(bytes, l); err != nil {
		return fmt.Errorf("failed to unmarshal Location: %w", err)
	}
	return nil
}

//	type PublicConfig struct {
//		PublicIPV4 string `json:"public_ip_v4"`
//		PublicIPV6 string `json:"public_ip_v6"`
//		Domain     string `json:"domain"`
//	}
//
// // Value implements the Valuer interface for storing PublicConfig in the database
//
//	func (c PublicConfig) Value() (driver.Value, error) {
//		bytes, err := json.Marshal(c)
//		if err != nil {
//			return nil, fmt.Errorf("failed to marshal PublicConfig: %w", err)
//		}
//		return string(bytes), nil
//	}
//
// // Scan implements the Scanner interface for retrieving PublicConfig from the database
//
//	func (c *PublicConfig) Scan(value any) error {
//		bytes, ok := value.([]byte)
//		if !ok {
//			return fmt.Errorf("invalid data type for PublicConfig: %T", value)
//		}
//
//		if err := json.Unmarshal(bytes, c); err != nil {
//			return fmt.Errorf("failed to unmarshal PublicConfig: %w", err)
//		}
//		return nil
//	}
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

// DefaultLimit returns the default values for the pagination
func DefaultLimit() Limit {
	return Limit{
		Size: 50,
		Page: 1,
	}
}
