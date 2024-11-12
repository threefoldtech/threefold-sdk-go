package api

type Status struct {
	Version string `json:"version"`
	Healthy bool   `json:"healthy"`
	DBConn  string `json:"db_conn"`
	RMBConn string `json:"rmb_conn"`
}

type CallParam struct {
	NodeId  *uint32
	TwinId  *uint32
	Command *string
	Payload *string
}

type Option struct {
	Size      *uint32  `schema:"size,omitempty"`
	Page      *uint32  `schema:"page,omitempty"`
	RetCount  *bool    `schema:"ret_count,omitempty"`
	Randomize *bool    `schema:"randomize,omitempty"`
	SortBy    *string  `schema:"sort_by,omitempty"`
	SortOrder *string  `schema:"sort_order,omitempty"`
	Balance   *float64 `schema:"balance,omitempty"`
}
