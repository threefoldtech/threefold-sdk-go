package models

import (
	"time"
)

type Twin struct {
	ID          int    `gorm:"primaryKey"`
	GridVersion string `gorm:"type:grid_version"`
	AccountID   string `gorm:"size:100"`
	PublicKey   string `gorm:"size:100"`
	Relay       string `gorm:"size:100;index:idx_twin_relay"`
	Farms       []Farm `gorm:"foreignKey:TwinID"`
	Nodes       []Node `gorm:"foreignKey:TwinID"`
}

type Farm struct {
	ID              int        `gorm:"primaryKey"`
	Name            string     `gorm:"size:100"`
	TwinID          int        `gorm:"index:idx_farm_twin"`
	StellarAddress  string     `gorm:"size:100"`
	Dedicated       bool       `gorm:"default:false"`
	PricingPolicyID int        `gorm:"foreignKey:PricingPolicyID"`
	Certification   string     `gorm:"type:certification"`
	GridVersion     string     `gorm:"type:grid_version"`
	Nodes           []Node     `gorm:"foreignKey:FarmID"`
	PublicIPs       []PublicIP `gorm:"foreignKey:FarmID"`
}

type Node struct {
	ID            int `gorm:"primaryKey"`
	TwinID        int `gorm:"foreignKey:TwinID"`
	FarmID        int `gorm:"foreignKey:FarmID;index:idx_node_farm"`
	TotalHRU      int
	TotalCRU      int
	TotalSRU      int
	TotalMRU      int
	SerialNumber  int
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	ExtraFee      int
	GridVersion   string         `gorm:"type:grid_version"`
	Certification string         `gorm:"type:certification"`
	Contracts     []NodeContract `gorm:"foreignKey:NodeID"`
	RentContracts []RentContract `gorm:"foreignKey:NodeID"`
	Location      Location       `gorm:"foreignKey:NodeID"`
	Power         Power          `gorm:"foreignKey:NodeID"`
	Interfaces    []Interface    `gorm:"foreignKey:NodeID"`
	Info          NodeInfo       `gorm:"foreignKey:NodeID"`
	// Dmi           Dmi            `gorm:"foreignKey:NodeID"`
	Gpus []GPU `gorm:"foreignKey:NodeID"`
}

type NodeContract struct {
	ID                 int       `gorm:"primaryKey"`
	NodeID             int       `gorm:"foreignKey:NodeID"`
	TwinID             int       `gorm:"foreignKey:TwinID"`
	ContractResourceID int       `gorm:"foreignKey:ContractResourceID"`
	SolutionProviderID int       `gorm:"foreignKey:SolutionProviderID"`
	State              string    `gorm:"type:contract_state"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	DeploymentData     string
	DeploymentHash     string
	IPsNum             int
}

type RentContract struct {
	ID                 int       `gorm:"primaryKey"`
	NodeID             int       `gorm:"foreignKey:NodeID"`
	TwinID             int       `gorm:"foreignKey:TwinID"`
	ContractResourceID int       `gorm:"foreignKey:ContractResourceID"`
	SolutionProviderID int       `gorm:"foreignKey:SolutionProviderID"`
	State              string    `gorm:"type:contract_state"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
}

type NameContract struct {
	ID                 int       `gorm:"primaryKey"`
	TwinID             int       `gorm:"foreignKey:TwinID"`
	ContractResourceID int       `gorm:"foreignKey:ContractResourceID"`
	SolutionProviderID int       `gorm:"foreignKey:SolutionProviderID"`
	State              string    `gorm:"type:contract_state"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
}

type ContractBill struct {
	ContractID int       `gorm:"primaryKey"`
	Discount   string    `gorm:"type:discount"`
	Timestamp  time.Time `gorm:"autoCreateTime"`
	Amount     int
}

type NodeInfo struct {
	NodeID        int `gorm:"primaryKey"`
	HasIPv6       bool
	NumWorkloads  int
	UploadSpeed   float64
	DownloadSpeed float64
}

type Hardware struct {
	NodeID    int    `gorm:"primaryKey"`
	Bios      string `gorm:"type:jsonb"`
	Baseboard string `gorm:"type:jsonb"`
	Processor string `gorm:"type:jsonb"`
	Memory    string `gorm:"type:jsonb"`
}

type GPU struct {
	NodeID   int    `gorm:"primaryKey"`
	ID       string `gorm:"size:100"`
	Vendor   string `gorm:"size:100"`
	Device   string `gorm:"size:100"`
	Contract int    `gorm:"foreignKey:Contract"`
}

type Location struct {
	NodeID    int `gorm:"primaryKey"`
	Latitude  int
	Longitude int
	Country   string `gorm:"size:100"`
	City      string `gorm:"size:100"`
	Region    string `gorm:"size:20"`
}

type PublicConfig struct {
	NodeID int    `gorm:"primaryKey"`
	IPv4   string `gorm:"size:50"`
	IPv6   string `gorm:"size:50"`
	GW4    string `gorm:"size:50"`
	GW6    string `gorm:"size:50"`
	Domain string `gorm:"size:50"`
}

type Interface struct {
	NodeID int    `gorm:"primaryKey"`
	Name   string `gorm:"size:50"`
	Mac    string `gorm:"size:50"`
	IPs    string `gorm:"size:50"`
}

type Power struct {
	NodeID           int    `gorm:"primaryKey"`
	State            string `gorm:"type:NodeState"`
	Target           string `gorm:"type:NodeState"`
	Status           string `gorm:"type:NodeState"`
	Healthy          bool
	LastUptimeReport time.Time
	TotalUptime      int
}

type ContractResource struct {
	ID  int `gorm:"primaryKey"`
	HRU int
	SRU int
	MRU int
	CRU int
}

type SolutionProvider struct {
	ID          int `gorm:"primaryKey"`
	Description string
	Link        string
	Approved    bool
	Providers   string
}

type PricingPolicy struct {
	ID                    int    `gorm:"primaryKey"`
	GridVersion           string `gorm:"type:grid_version"`
	Name                  string `gorm:"size:100"`
	FoundationAccount     string `gorm:"size:100"`
	CertifiedSalesAccount string `gorm:"size:100"`
	DedicationDiscount    int
	SUValue               int
	SUUnit                string `gorm:"size:10"`
	CUValue               int
	CUUnit                string `gorm:"size:10"`
	NUValue               int
	NUUnit                string `gorm:"size:10"`
	IPUValue              int
	IPUUnit               string `gorm:"size:10"`
}

type PublicIP struct {
	FarmID     int    `gorm:"foreignKey:FarmID"`
	ContractID int    `gorm:"foreignKey:ContractID"`
	Gateway    string `gorm:"size:25"`
	IP         string `gorm:"size:25"`
}

type ContractState string

const (
	Created     ContractState = "created"
	Deleted     ContractState = "deleted"
	GracePeriod ContractState = "grace_period"
	OutOfFunds  ContractState = "out_of_funds"
)

type NodeState string

const (
	Up      NodeState = "Up"
	Down    NodeState = "Down"
	Standby NodeState = "Standby"
)
