// Package deployer for project deployment
package deployer

import (
	"fmt"
	"math/rand"
	"net"
	"strings"

	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/pkg/types"
)

// VMSpec struct to hold vm specs
type VMSpec struct {
	CPU     int
	Memory  int
	Storage int
	Public  bool
}

var (
	// Eco spec
	Eco = VMSpec{1, 2, 5, true}
	// Standard spec
	Standard = VMSpec{2, 4, 10, true}
	// Performance spec
	Performance = VMSpec{4, 8, 15, true}
)

var (
	vmFlist = "https://hub.grid.tf/aelawady.3bot/abdulrahmanelawady-gridify-test-latest.flist"

	vmEntrypoint = "/init.sh"
	vmPlanetary  = true
)

func buildNodeFilter(vmSpec VMSpec) types.NodeFilter {
	nodeStatus := "up"
	freeMRU := uint64(vmSpec.Memory * 1024 * 1024 * 1024)
	freeSRU := uint64(vmSpec.Storage * 1024 * 1024 * 1024)
	freeIPs := uint64(0)
	if vmSpec.Public {
		freeIPs = 1
	}
	domain := true

	filter := types.NodeFilter{
		FarmIDs: []uint64{1},
		Status:  &nodeStatus,
		FreeMRU: &freeMRU,
		FreeSRU: &freeSRU,
		FreeIPs: &freeIPs,
		Domain:  &domain,
	}
	return filter
}

func buildNetwork(projectName string, node uint32, generator rand.Rand) workloads.ZNet {
	networkName := randName(10, generator)
	network := workloads.ZNet{
		Name:  networkName,
		Nodes: []uint32{node},
		IPRange: workloads.NewIPRange(net.IPNet{
			IP:   net.IPv4(10, 20, 0, 0),
			Mask: net.CIDRMask(16, 32),
		}),
		SolutionType: projectName,
	}
	return network
}

func buildDeployment(vmSpec VMSpec, networkName, projectName, repoURL string, node uint32, generator rand.Rand) workloads.Deployment {
	vmName := randName(10, generator)
	vm := workloads.VM{
		Name:       vmName,
		Flist:      vmFlist,
		CPU:        vmSpec.CPU,
		Memory:     vmSpec.Memory * 1024,
		RootfsSize: vmSpec.Storage * 1024,
		PublicIP:   vmSpec.Public,
		Planetary:  vmPlanetary,
		Entrypoint: vmEntrypoint,
		EnvVars: map[string]string{
			"REPO_URL": repoURL,
		},
		NetworkName: networkName,
	}

	dl := workloads.NewDeployment(vm.Name, node, projectName, nil, networkName, nil, nil, []workloads.VM{vm}, nil)
	return dl
}

func buildGateway(backend, projectName string, node uint32, generator rand.Rand) workloads.GatewayNameProxy {
	subdomain := randName(10, generator)
	gateway := workloads.GatewayNameProxy{
		NodeID:       node,
		Name:         subdomain,
		Backends:     workloads.NewZosBackends([]string{backend}),
		SolutionType: projectName,
	}
	return gateway
}

func buildPortlessBackend(ip string) string {
	publicIP := strings.Split(ip, "/")[0]
	backend := fmt.Sprintf("http://%s", publicIP)
	return backend
}
