package api

import (
	"github.com/threefoldtech/tfgrid-sdk-go/grid-api/types"
)

type API interface {
	Nodes(types.NodeFilter, types.NodeSort, types.Option) ([]types.Node, error)
	Node(uint32) (types.Node, error)
	Hardwares(types.HardwareFilter, types.HardwareSort, types.Option) ([]types.Hardware, error)

	Farms(types.FarmFilter, types.FarmSort, types.Option) ([]types.Farm, error)
	Farm(uint32) (types.Farm, error)
	Ips(types.PublicIpFilter, types.PublicIpSort, types.Option) ([]types.PublicIP, error)

	Twins(types.TwinFilter, types.TwinSort, types.Option) ([]types.Twin, error)
	Twin(uint32) (types.Twin, error)
	TwinSpending(uint32) (types.Spending, error)

	Contracts(types.ContractFilter, types.ContractSort, types.Option) ([]types.Contract, error)
	Contract(uint32) (types.Contract, error)
	Bills(types.BillFilter, types.BillSort, types.Option) ([]types.Bill, error)

	Stats(types.StatsFilter) (types.Stats, error)
	Status() (types.Status, error)
	Call(types.CallParam) (any, error)
}

type ApiClient struct{}

var _ API = (*ApiClient)(nil)

func (api *ApiClient) Nodes(filter types.NodeFilter, sort types.NodeSort, option types.Option) ([]types.Node, error) {
	return []types.Node{}, nil
}

func (api *ApiClient) Node(id uint32) (types.Node, error) {
	return types.Node{}, nil
}

func (api *ApiClient) Hardwares(filter types.HardwareFilter, sort types.HardwareSort, option types.Option) ([]types.Hardware, error) {
	return []types.Hardware{}, nil
}

func (api *ApiClient) Farms(filter types.FarmFilter, sort types.FarmSort, option types.Option) ([]types.Farm, error) {
	return []types.Farm{}, nil
}

func (api *ApiClient) Farm(id uint32) (types.Farm, error) {
	return types.Farm{}, nil
}

func (api *ApiClient) Ips(filter types.PublicIpFilter, sort types.PublicIpSort, option types.Option) ([]types.PublicIP, error) {
	return []types.PublicIP{}, nil
}

func (api *ApiClient) Twins(filter types.TwinFilter, sort types.TwinSort, option types.Option) ([]types.Twin, error) {
	return []types.Twin{}, nil
}

func (api *ApiClient) Twin(id uint32) (types.Twin, error) {
	return types.Twin{}, nil
}

func (api *ApiClient) TwinSpending(id uint32) (types.Spending, error) {
	return types.Spending{}, nil
}

func (api *ApiClient) Contracts(filter types.ContractFilter, sort types.ContractSort, option types.Option) ([]types.Contract, error) {
	return []types.Contract{}, nil
}

func (api *ApiClient) Contract(id uint32) (types.Contract, error) {
	return types.Contract{}, nil
}

func (api *ApiClient) Bills(filter types.BillFilter, sort types.BillSort, option types.Option) ([]types.Bill, error) {
	return []types.Bill{}, nil
}

func (api *ApiClient) Stats(filter types.StatsFilter) (types.Stats, error) {
	return types.Stats{}, nil
}

func (api *ApiClient) Status() (types.Status, error) {
	return types.Status{}, nil
}

func (api *ApiClient) Call(param types.CallParam) (any, error) {
	return nil, nil
}
