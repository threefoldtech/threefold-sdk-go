package types

type Bill struct {
	ContractId uint32 `json:"contract_id"`
	Amount     uint64 `json:"amount"`
	Discount   string `json:"discount"`
	Timestamp  uint64 `json:"timestamp"`
}

type BillFilter struct {
	ContractId    *uint32
	TwinId        *uint32
	Discount      *string
	CreatedAfter  *uint32
	CreatedBefore *uint32
	AmountMore    *uint64
	AmountLess    *uint64
}

type BillSort struct {
	ContractId *bool
	Amount     *bool
	Discount   *bool
	Timestamp  *bool
}
