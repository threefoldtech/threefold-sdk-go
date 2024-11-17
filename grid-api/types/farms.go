package types

/*
Removed
public ips

What
pricing vs farming policies

TODO: ask sameh about farming policies
*/
type Farm struct {
	Name           string   `json:"name"`
	FarmID         uint32   `json:"farm_id"`
	TwinID         uint32   `json:"twin_id"`
	PricingPolicy  uint8    `json:"pricing_policy"`
	Certification  string   `json:"certification"`
	StellarAddress string   `json:"stellar_address"`
	Dedicated      bool     `json:"dedicated"`
	Nodes          []uint32 `json:"nodes"`
	IpsCount       uint8    `json:"ips"`
}

// TODO: use the node filter
// type FarmNodeFilter struct {
// 	NodeFreeMRU      *uint64  `schema:"node_free_mru,omitempty"`
// 	NodeFreeHRU      *uint64  `schema:"node_free_hru,omitempty"`
// 	NodeFreeSRU      *uint64  `schema:"node_free_sru,omitempty"`
// 	NodeTotalCRU     *uint64  `schema:"node_total_cru,omitempty"`
// 	NodeRentedBy     *uint64  `schema:"node_rented_by,omitempty"`
// 	NodeAvailableFor *uint64  `schema:"node_available_for,omitempty"`
// 	NodeHasGPU       *bool    `schema:"node_has_gpu,omitempty"`
// 	NodeHasIpv6      *bool    `schema:"node_has_ipv6,omitempty"`
// 	NodeCertified    *bool    `schema:"node_certified,omitempty"`
// 	NodeStatus       []string `schema:"node_status,omitempty"`
// 	NodeFeatures     []string `schema:"node_features,omitempty"`
// }

type FarmFilter struct {
	Name            *string  `schema:"name,omitempty"`
	Certification   *string  `schema:"certification,omitempty"`
	Country         *string  `schema:"country,omitempty"`
	Region          *string  `schema:"region,omitempty"`
	StellarAddress  *string  `schema:"stellar_address,omitempty"`
	PricingPolicyID *uint8   `schema:"pricing_policy_id,omitempty"`
	Dedicated       *bool    `schema:"dedicated,omitempty"`
	FarmIDs         []uint32 `schema:"farm_id,omitempty"`
	TwinIDs         []uint32 `schema:"twin_id,omitempty"`
	FreeIPs         *uint8   `schema:"free_ips,omitempty"`
	TotalIPs        *uint8   `schema:"total_ips,omitempty"`
	NodeFilter               // filter farm based on its nodes
}

type FarmSort struct {
	Name          *bool
	FarmId        *bool
	TwinId        *bool
	Dedicated     *bool
	Ips           *bool
	Nodes         *bool
	Certification *bool
}
