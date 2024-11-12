package api

type Stats struct {
	// location
	Countries uint32 `json:"countries"`
	Farms     uint32 `json:"farms"`
	// nodes
	TotalNodes     uint32 `json:"total_nodes"`
	DedicatedNodes uint32 `json:"dedicated_nodes"`
	AccessNodes    uint32 `json:"access_nodes"`
	GatewayNodes   uint32 `json:"gateway_nodes"`
	// resources
	CRU  uint64 `json:"cru"`
	SRU  uint64 `json:"sru"`
	MRU  uint64 `json:"mru"`
	HRU  uint64 `json:"hru"`
	GPUs uint32 `json:"gpus"`
	IPs  uint32 `json:"ips"`
	// usage
	Twins     uint32 `json:"twins"`
	Contracts uint32 `json:"contracts"`
	Workloads uint32 `json:"workloads"`

	Distribution map[string]uint32 `json:"distribution"`
}

type StatsFilter struct {
	Status []string `schema:"status,omitempty"`
}
