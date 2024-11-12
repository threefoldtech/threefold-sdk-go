package api

type PublicIP struct {
	FarmID     uint32 `json:"farm_id"`
	ContractID uint32 `json:"contract_id"`
	IP         string `json:"ip"`
	Gateway    string `json:"gateway"`
}

type PublicIpFilter struct {
	FarmIDs     []uint32 `schema:"farm_ids,omitempty"`
	ContractIDs []uint32
	Ip          *string `schema:"ip,omitempty"`
	Gateway     *string `schema:"gateway,omitempty"`
}

type PublicIpSort struct {
	FarmId     *bool
	ContractId *bool
	Ip         *bool
}
