package api

import (
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

/*
Removed:
- policies
- grid version
- certification type
- rentable
- renter
- rent contract
- contract-id in gpu
- farm_free_ips
*/
type Node struct {
	NodeID       uint32       `json:"node_id"`
	FarmID       uint32       `json:"farm_id"`
	TwinID       uint32       `json:"twin_id"`
	FarmName     string       `json:"farm_name"`
	PriceUsd     float64      `json:"price_usd"`
	Capacity     Capacities   `json:"capacity"`
	Location     Location     `json:"location"`
	PublicConfig PublicConfig `json:"public_config"`
	Reservation  Reservation  `json:"reservation"`
	Power        Power        `json:"power"`
	Hardware     Hardware     `json:"hardware"` // moving to /hardware endpoint
	Speed        Speed        `json:"speed"`
	Features     []string     `json:"features"`
}

type Location struct {
	Region    string  `json:"region"`
	Country   string  `json:"country"`
	City      string  `json:"city"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Power struct {
	State     string `json:"state"`
	Target    string `json:"target"`
	UpdatedAt uint32 `json:"updated_at"`
	CreatedAt uint32 `json:"created_at"`
	Uptime    uint32 `json:"uptime"`
	Status    string `json:"status"`
	Healthy   bool   `json:"healthy"`
}

type Capacities struct {
	Total Capacity `json:"total"` // add zos used to both
	Used  Capacity `json:"used"`
	Free  Capacity `json:"free"`
}

type Capacity struct {
	GPU uint8          `json:"gpu"`
	CRU uint64         `json:"cru"`
	SRU gridtypes.Unit `json:"sru"`
	HRU gridtypes.Unit `json:"hru"`
	MRU gridtypes.Unit `json:"mru"`
}

type PublicConfig struct {
	Domain string `json:"domain"`
	Gw4    string `json:"gw4"`
	Gw6    string `json:"gw6"`
	Ipv4   string `json:"ipv4"`
	Ipv6   string `json:"ipv6"`
}

type Reservation struct {
	Dedicated bool    `json:"dedicated"`
	Renter    uint32  `json:"rented"`
	ExtraFee  float64 `json:"extra_fee"`
}

type Hardware struct {
	SerialNumber string      `json:"serial_number"`
	BIOS         BIOS        `json:"bios"`
	Baseboard    Baseboard   `json:"baseboard"`
	Processor    []Processor `json:"processor"`
	Memory       []Memory    `json:"memory"`
	// GPU          []GPU       `json:"gpu"`
}

type BIOS struct {
	Vendor  string `json:"vendor"`
	Version string `json:"version"`
}

type Baseboard struct {
	Manufacturer string `json:"manufacturer"`
	ProductName  string `json:"product_name"`
}

type Processor struct {
	Version     string `json:"version"`
	ThreadCount string `json:"thread_count"`
}

type Memory struct {
	Manufacturer string `json:"manufacturer"`
	Type         string `json:"type"`
}

type Speed struct {
	Upload   float64 `json:"upload"`
	Download float64 `json:"download"`
}

type NodeFilter struct {
	// resources
	FreeMRU  *uint64 `schema:"free_mru,omitempty"`
	FreeHRU  *uint64 `schema:"free_hru,omitempty"`
	FreeSRU  *uint64 `schema:"free_sru,omitempty"`
	FreeGPU  *uint8  `schema:"free_sru,omitempty"`
	TotalMRU *uint64 `schema:"total_mru,omitempty"`
	TotalHRU *uint64 `schema:"total_hru,omitempty"`
	TotalSRU *uint64 `schema:"total_sru,omitempty"`
	TotalCRU *uint8  `schema:"total_cru,omitempty"`
	TotalGPU *uint8  `schema:"total_gpu,omitempty"`

	// location
	Region  *string `schema:"region,omitempty"`
	Country *string `schema:"country,omitempty"`
	City    *string `schema:"city,omitempty"`

	// reservation
	Dedicated          *bool   `schema:"dedicated,omitempty"`
	Renter             *uint32 `schema:"renter,omitempty"`
	Rentable           *bool   `schema:"rentable,omitempty"`
	RentableOrRentedBy *uint64 `schema:"rentable_or_rented_by,omitempty"` // rented by twin or rentable
	SharedOrRentedBy   *uint64 `schema:"available_for,omitempty"`         // rented by twin or free

	// network
	HasIpv6      *bool `schema:"has_ipv6,omitempty"`
	IsGateway    *bool
	AvailableIPs *uint8 `schema:"available_ips,omitempty"`

	// ids
	FarmIDs  []uint32 `schema:"farm_ids,omitempty"`
	NodeIDs  []uint32 `schema:"node_ids,omitempty"`
	TwinIDs  []uint32 `schema:"twin_ids,omitempty"`
	Excluded []uint32 `schema:"excluded,omitempty"`
	OwnedBy  *uint64  `schema:"owned_by,omitempty"`

	// power
	Status       []string `schema:"status,omitempty"`
	Healthy      *bool    `schema:"healthy,omitempty"`
	UpdatedAfter uint32   `schema:"updated_after"`

	// other
	FarmName      *string `schema:"farm_name,omitempty"`
	Certification *string `schema:"certification,omitempty"`

	CostMore *float64 `schema:"price_min,omitempty"`
	CostLess *float64 `schema:"price_max,omitempty"`

	Features []string `schema:"features,omitempty"`
}

type NodeSort struct {
	NodeId    *bool
	FarmId    *bool
	TwinId    *bool
	Price     *bool
	Status    *bool
	Healthy   *bool
	UpdatedAt *bool
	CreatedAt *bool

	FreeMRU  *bool
	FreeHRU  *bool
	FreeSRU  *bool
	FreeCRU  *bool
	FreeGPU  *bool
	TotalMRU *bool
	TotalHRU *bool
	TotalSRU *bool
	TotalCRU *bool
	TotalGPU *bool
	UsedMRU  *bool
	UsedHRU  *bool
	UsedSRU  *bool
	UsedCRU  *bool
	UsedGPU  *bool

	Dedicated *bool
	Renter    *bool
	ExtraFee  *bool
}
