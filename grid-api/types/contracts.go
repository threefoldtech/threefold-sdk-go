package types

type Deployment struct {
	Type    string `json:"type"`
	Project string `json:"project"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Hash    string `json:"hash"`
	Raw     string // in case can't parse
}

// NOTE: unified type with optional fields, different types with multiple endpoints
type Contract struct {
	ContractID uint32     `json:"contract_id"`
	TwinID     uint32     `json:"twin_id"`
	NodeID     uint32     `json:"node_id"`
	FarmID     uint32     `json:"farm_id"`
	Name       string     `json:"name,omitempty"`
	Type       string     `json:"type"`
	State      string     `json:"state"`
	CreatedAt  uint       `json:"created_at"`
	NumIps     uint8      `json:"num_ip"`
	Deployment Deployment `json:"deployment"`
}

type ContractFilter struct {
	ContractIDs       []uint32
	TwinIDs           []uint32
	NodeIDs           []uint32
	FarmIDs           []uint32
	Type              []string
	State             []string
	DeploymentProject *string
	DeploymentName    *string
	DeploymentType    *string
	UsedIps           *uint8
	CreatedAfter      *uint32
	CreatedBefore     *uint32
}

type ContractSort struct {
	ContractID *bool
	TwinID     *bool
	NodeID     *bool
	FarmID     *bool
	Type       *bool
	State      *bool
	CreatedAt  *bool
	NumIps     *bool
}
