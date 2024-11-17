package query

import (
	"github.com/threefoldtech/tfgrid-sdk-go/grid-api/internal/db/models"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-api/types"
)

type QueryClient interface {
	Nodes(types.NodeFilter, types.NodeSort, types.Option) ([]models.Node, uint32, error)
	Farms(types.FarmFilter, types.FarmSort, types.Option) ([]models.Farm, uint32, error)
	Twins(types.TwinFilter, types.TwinSort, types.Option) ([]models.Twin, uint32, error)
	Contracts(types.ContractFilter, types.ContractSort, types.Option) ([]models.NodeContract, uint32, error)

	Hardwares(types.HardwareFilter, types.HardwareSort, types.Option) ([]models.Hardware, error)
	Ips(types.PublicIpFilter, types.PublicIpSort, types.Option) ([]models.PublicIP, error)
	Bills(types.BillFilter, types.BillSort, types.Option) ([]models.ContractBill, error)
	Stats(types.StatsFilter) (types.Stats, error)
}
