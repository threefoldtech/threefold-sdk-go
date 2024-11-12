package api

type API interface {
	Nodes(NodeFilter, NodeSort, Option) []Node
	Node(uint32) Node
	Hardware(HWF) []Hardware
	Farms(FarmFilter, FarmSort, Option) []Farm
	Farm(uint32) Farm
	Twins(TwinFilter, TwinSort, Option) []Twin
	Twin(uint32) Twin
	TwinSpending(uint32) Spending
	Contracts(ContractFilter, ContractSort, Option) []Contract
	Contract(uint32) Contract
	Bills(BillFilter, BillSort, Option) []Bill
	Ips(PublicIpFilter, PublicIpSort, Option) []PublicIP
	Gpus(GPUFilter, GPUSort, Option) []GPU
	Stats(StatsFilter) Stats

	Status() Status
	Call(CallParam) any
}

type ApiClient struct{}

var _ API = (*ApiClient)(nil)

func (api *ApiClient) Nodes(NodeFilter, NodeSort, Option) []Node {
	return []Node{}
}

func (api *ApiClient) Node(uint32) Node {
	return Node{}
}

func (api *ApiClient) Farms(FarmFilter, FarmSort, Option) []Farm {
	return []Farm{}
}

func (api *ApiClient) Farm(uint32) Farm {
	return Farm{}
}

func (api *ApiClient) Twins(TwinFilter, TwinSort, Option) []Twin {
	return []Twin{}
}

func (api *ApiClient) Twin(uint32) Twin {
	return Twin{}
}

func (api *ApiClient) TwinSpending(uint32) Spending {
	return Spending{}
}

func (api *ApiClient) Contracts(ContractFilter, ContractSort, Option) []Contract {
	return []Contract{}
}

func (api *ApiClient) Contract(uint32) Contract {
	return Contract{}
}

func (api *ApiClient) Bills(BillFilter, BillSort, Option) []Bill {
	return []Bill{}
}

func (api *ApiClient) Ips(PublicIpFilter, PublicIpSort, Option) []PublicIP {
	return []PublicIP{}
}

func (api *ApiClient) Gpus(GPUFilter, GPUSort, Option) []GPU {
	return []GPU{}
}

func (api *ApiClient) Stats(StatsFilter) Stats {
	return Stats{}
}

func (api *ApiClient) Status() Status {
	return Status{}
}

func (api *ApiClient) Call(CallParam) any {
	return nil
}
