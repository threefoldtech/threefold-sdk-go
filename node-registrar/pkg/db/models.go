package db

import (
	"time"

	"github.com/lib/pq"
)

type Account struct {
	TwinID    uint64         `gorm:"primaryKey;autoIncrement" json:"twin_id"`
	Relays    pq.StringArray `gorm:"type:text[];default:'{}'" json:"relays" swaggertype:"array,string"` // Optional list of relay domains
	RMBEncKey string         `gorm:"type:text" json:"rmb_enc_key"`                                      // Optional base64 encoded public key for rmb communication
	CreatedAt time.Time
	UpdatedAt time.Time
	// The public key (ED25519 for nodes, ED25519 or SR25519 for farmers) in the more standard base64 since we are moving from substrate echo system?
	// (still SS58 can be used or plain base58 ,TBD)
	PublicKey string `gorm:"type:text;not null;unique" json:"public_key"`
	// Relations | likely we need to use OnDelete:RESTRICT (Prevent Twin deletion if farms exist)
	// @swagger:ignore
	Farms     []Farm     `gorm:"foreignKey:TwinID;references:TwinID;constraint:OnDelete:RESTRICT"`
	Contracts []Contract `json:"contracts" gorm:"foreignKey:ContractID;references:ContractID;constraint:OnDelete:CASCADE"`
}

type Farm struct {
	FarmID    uint64 `gorm:"primaryKey;autoIncrement" json:"farm_id"`
	FarmName  string `gorm:"size:40;not null;unique;check:farm_name <> ''" json:"farm_name"`
	TwinID    uint64 `json:"twin_id" gorm:"not null;check:twin_id > 0"` // Farmer account reference
	Dedicated bool   `json:"dedicated"`
	CreatedAt time.Time
	UpdatedAt time.Time
	// @swagger:ignore
	Nodes []Node `gorm:"foreignKey:FarmID;references:FarmID;constraint:OnDelete:RESTRICT" json:"nodes"`
}

type Node struct {
	NodeID uint64 `json:"node_id" gorm:"primaryKey;autoIncrement"`
	// Constraints set to prevents unintended account deletion if linked Farms/nodes exist.
	FarmID uint64 `json:"farm_id" gorm:"not null;check:farm_id> 0;foreignKey:FarmID;references:FarmID;constraint:OnDelete:RESTRICT"`
	TwinID uint64 `json:"twin_id" gorm:"not null;unique;check:twin_id > 0;foreignKey:TwinID;references:TwinID;constraint:OnDelete:RESTRICT"` // Node account reference

	Location Location `json:"location" gorm:"not null;type:json;serializer:json"`

	// PublicConfig PublicConfig `json:"public_config" gorm:"type:json"`
	Resources    Resources   `json:"resources" gorm:"not null;type:json;serializer:json"`
	Interfaces   []Interface `gorm:"not null;type:json;serializer:json"`
	SecureBoot   bool
	Virtualized  bool
	SerialNumber string

	UptimeReports []UptimeReport `json:"uptime" gorm:"foreignKey:NodeID;references:NodeID;constraint:OnDelete:CASCADE"`
	Contracts     []Contract     `json:"contracts" gorm:"foreignKey:ContractID;references:ContractID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
	Approved  bool
}

type Contract struct {
	ContractID uint64 `json:"contract_id" gorm:"primaryKey;autoIncrement"`
	// gorm:"size:40;not null;unique;check:farm_name <> ''" json:"farm_name"
	State        string `json:"state" gorm:"size:40;not null;check:state<> ''"`
	TwinID       uint64 `json:"twin_id"`
	ContractType string `json:"contract_type"`
	// SolutionProviderID types.OptionU64 `json:"solution_provider_id"`
}

type UptimeReport struct {
	ID         uint64        `gorm:"primaryKey;autoIncrement"`
	NodeID     uint64        `gorm:"index" json:"node_id"`
	Duration   time.Duration `swaggertype:"integer"` // Uptime duration for this period
	Timestamp  time.Time     `gorm:"index"`
	WasRestart bool          // True if this report followed a restart
	CreatedAt  time.Time
}

type ZosVersion struct {
	Key     string `gorm:"primaryKey;size:50"`
	Version string `gorm:"not null"`
}

type Interface struct {
	Name string `json:"name"`
	Mac  string `json:"mac"`
	IPs  string `json:"ips"`
}

type Resources struct {
	HRU uint64 `json:"hru"`
	SRU uint64 `json:"sru"`
	CRU uint64 `json:"cru"`
	MRU uint64 `json:"mru"`
}

type Location struct {
	Country   string `json:"country" gorm:"not null"`
	City      string `json:"city" gorm:"not null"`
	Longitude string `json:"longitude" gorm:"not null"`
	Latitude  string `json:"latitude" gorm:"not null"`
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

// DefaultLimit returns the default values for the pagination
func DefaultLimit() Limit {
	return Limit{
		Size: 50,
		Page: 1,
	}
}
