package types

type Hardware struct {
	SerialNumber string      `json:"serial_number"`
	BIOS         BIOS        `json:"bios"`
	Baseboard    Baseboard   `json:"baseboard"`
	Processor    []Processor `json:"processor"`
	Memory       []Memory    `json:"memory"`
	GPU          []GPU       `json:"gpu"`
}

type GPU struct {
	ID     string `json:"id"`
	NodeId uint32
	Vendor string `json:"vendor"`
	Device string `json:"device"`
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

type HardwareFilter struct {
	NodeIds   []uint32
	GpuDevice *string
	GpuVendor *string
}

type HardwareSort struct {
	NodeId    *bool
	GpuID     *bool
	GpuVendor *bool
	GpuDevice *bool
}
