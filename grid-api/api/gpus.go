package api

type GPU struct {
	ID     string `json:"id"`
	NodeId uint32
	Vendor string `json:"vendor"`
	Device string `json:"device"`
}

type GPUFilter struct {
	Device *string `schema:"gpu_device,omitempty"`
	Vendor *string `schema:"gpu_vendor,omitempty"`
}

type GPUSort struct {
	ID     *bool
	Vendor *bool
	Device *bool
}
