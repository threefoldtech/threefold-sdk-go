package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/pkg/errors"
	zerolog "github.com/rs/zerolog/log"
	client "github.com/threefoldtech/tfgrid-sdk-go/grid-client/node"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	zosTypes "github.com/threefoldtech/tfgrid-sdk-go/grid-client/zos"
	"github.com/threefoldtech/zos/pkg/gridtypes"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

// K8sDeployer for deploying k8s
type K8sDeployer struct {
	tfPluginClient *TFPluginClient
	deployer       MockDeployer
}

// NewK8sDeployer generates new K8s Deployer
func NewK8sDeployer(tfPluginClient *TFPluginClient) K8sDeployer {
	deployer := NewDeployer(*tfPluginClient, true)
	k8sDeployer := K8sDeployer{
		tfPluginClient: tfPluginClient,
		deployer:       &deployer,
	}

	return k8sDeployer
}

// Validate validates K8s deployer
func (d *K8sDeployer) Validate(ctx context.Context, k8sCluster *workloads.K8sCluster) error {
	sub := d.tfPluginClient.SubstrateConn

	if err := validateAccountBalanceForExtrinsics(sub, d.tfPluginClient.Identity); err != nil {
		return d.tfPluginClient.sentry.error(err)
	}

	if err := d.tfPluginClient.State.AssignNodesIPRange(k8sCluster); err != nil {
		return d.tfPluginClient.sentry.error(err)
	}

	if err := k8sCluster.Validate(); err != nil {
		return d.tfPluginClient.sentry.error(err)
	}

	// validate cluster nodes
	var nodes []uint32
	nodes = append(nodes, k8sCluster.Master.NodeID)
	for _, worker := range k8sCluster.Workers {
		if !workloads.Contains(nodes, worker.NodeID) {
			nodes = append(nodes, worker.NodeID)
		}
	}
	return d.tfPluginClient.sentry.error(client.AreNodesUp(ctx, sub, nodes, d.tfPluginClient.NcPool))
}

// generateVersionlessDeployments generates a new deployment without a version
func (d *K8sDeployer) generateVersionlessDeployments(k8sCluster *workloads.K8sCluster) (map[uint32]zosTypes.Deployment, error) {
	err := d.assignNodesIPs(k8sCluster)
	if err != nil {
		return nil, errors.Wrap(err, "failed to assign node ips")
	}
	deployments := make(map[uint32]zosTypes.Deployment)
	nodeWorkloads := make(map[uint32][]zosTypes.Workload)

	masterWorkloads := k8sCluster.Master.MasterZosWorkload(k8sCluster)
	for _, m := range masterWorkloads {
		nodeWorkloads[k8sCluster.Master.NodeID] = append(nodeWorkloads[k8sCluster.Master.NodeID], zosTypes.NewWorkloadFromZosWorkload(m))
	}
	for _, w := range k8sCluster.Workers {
		workerWorkloads := w.WorkerZosWorkload(k8sCluster)
		for _, wr := range workerWorkloads {
			nodeWorkloads[w.NodeID] = append(nodeWorkloads[w.NodeID], zosTypes.NewWorkloadFromZosWorkload(wr))
		}
	}

	for node, ws := range nodeWorkloads {
		dl := workloads.NewGridDeployment(d.tfPluginClient.TwinID, 0, ws)
		dl.Metadata, err = k8sCluster.GenerateMetadata()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to generate deployment %s metadata", k8sCluster.Master.Name)
		}

		deployments[node] = dl
	}
	return deployments, nil
}

// Deploy deploys a k8s cluster deployment
func (d *K8sDeployer) Deploy(ctx context.Context, k8sCluster *workloads.K8sCluster) error {
	if err := d.tfPluginClient.State.AssignNodesIPRange(k8sCluster); err != nil {
		return d.tfPluginClient.sentry.error(err)
	}

	err := k8sCluster.InvalidateBrokenAttributes(d.tfPluginClient.SubstrateConn)
	if err != nil {
		return d.tfPluginClient.sentry.error(err)
	}

	assignNodesFlistsAndEntryPoints(k8sCluster)

	if err := d.Validate(ctx, k8sCluster); err != nil {
		return d.tfPluginClient.sentry.error(err)
	}

	newDeployments, err := d.generateVersionlessDeployments(k8sCluster)
	if err != nil {
		return d.tfPluginClient.sentry.error(errors.Wrap(err, "could not generate k8s grid deployments"))
	}

	newDeploymentsSolutionProvider := make(map[uint32]*uint64)
	newDeploymentsSolutionProvider[k8sCluster.Master.NodeID] = nil

	k8sCluster.NodeDeploymentID, err = d.deployer.Deploy(ctx, k8sCluster.NodeDeploymentID, newDeployments, newDeploymentsSolutionProvider)

	// update deployments state
	// error is not returned immediately before updating state because of untracked failed deployments
	if contractID, ok := k8sCluster.NodeDeploymentID[k8sCluster.Master.NodeID]; ok && contractID != 0 {
		d.tfPluginClient.State.StoreContractIDs(k8sCluster.Master.NodeID, contractID)
		for _, w := range k8sCluster.Workers {
			d.tfPluginClient.State.StoreContractIDs(w.NodeID, k8sCluster.NodeDeploymentID[w.NodeID])
		}
	}

	return err
}

// BatchDeploy deploys multiple clusters using the deployer
func (d *K8sDeployer) BatchDeploy(ctx context.Context, k8sClusters []*workloads.K8sCluster) error {
	newDeployments := make(map[uint32][]zosTypes.Deployment)
	newDeploymentsSolutionProvider := make(map[uint32][]*uint64)

	for _, k8sCluster := range k8sClusters {
		if err := d.tfPluginClient.State.AssignNodesIPRange(k8sCluster); err != nil {
			return d.tfPluginClient.sentry.error(err)
		}

		err := k8sCluster.InvalidateBrokenAttributes(d.tfPluginClient.SubstrateConn)
		if err != nil {
			return d.tfPluginClient.sentry.error(err)
		}

		assignNodesFlistsAndEntryPoints(k8sCluster)

		if err := d.Validate(ctx, k8sCluster); err != nil {
			return d.tfPluginClient.sentry.error(err)
		}

		dls, err := d.generateVersionlessDeployments(k8sCluster)
		if err != nil {
			return d.tfPluginClient.sentry.error(errors.Wrap(err, "could not generate k8s grid deployments"))
		}

		for nodeID, dl := range dls {
			// solution providers
			newDeploymentsSolutionProvider[nodeID] = nil

			if _, ok := newDeployments[nodeID]; !ok {
				newDeployments[nodeID] = []zosTypes.Deployment{dl}
				continue
			}
			newDeployments[nodeID] = append(newDeployments[nodeID], dl)
		}
	}

	newDls, err := d.deployer.BatchDeploy(ctx, newDeployments, newDeploymentsSolutionProvider)

	// update deployments state
	// error is not returned immediately before updating state because of untracked failed deployments
	for _, k8sCluster := range k8sClusters {
		if err := d.updateStateFromDeployments(k8sCluster, newDls); err != nil {
			return d.tfPluginClient.sentry.error(errors.Wrapf(err, "failed to update cluster with master name '%s' state", k8sCluster.Master.Name))
		}
	}

	return d.tfPluginClient.sentry.error(err)
}

// Cancel cancels a k8s cluster deployment
func (d *K8sDeployer) Cancel(ctx context.Context, k8sCluster *workloads.K8sCluster) (err error) {
	for nodeID, contractID := range k8sCluster.NodeDeploymentID {
		if k8sCluster.Master.NodeID == nodeID {
			err = d.deployer.Cancel(ctx, contractID)
			if err != nil {
				return d.tfPluginClient.sentry.error(errors.Wrapf(err, "could not cancel master %s, contract %d", k8sCluster.Master.Name, contractID))
			}
			d.tfPluginClient.State.CurrentNodeDeployments[nodeID] = workloads.Delete(d.tfPluginClient.State.CurrentNodeDeployments[nodeID], contractID)
			delete(k8sCluster.NodeDeploymentID, nodeID)
			continue
		}
		for _, worker := range k8sCluster.Workers {
			if worker.NodeID == nodeID {
				err = d.deployer.Cancel(ctx, contractID)
				if err != nil {
					return d.tfPluginClient.sentry.error(errors.Wrapf(err, "could not cancel worker %s, contract %d", worker.Name, contractID))
				}
				d.tfPluginClient.State.CurrentNodeDeployments[nodeID] = workloads.Delete(d.tfPluginClient.State.CurrentNodeDeployments[nodeID], contractID)
				delete(k8sCluster.NodeDeploymentID, nodeID)
				break
			}
		}
	}

	return nil
}

func (d *K8sDeployer) updateStateFromDeployments(k8sCluster *workloads.K8sCluster, newDl map[uint32][]zosTypes.Deployment) error {
	k8sNodes := []uint32{k8sCluster.Master.NodeID}
	for _, w := range k8sCluster.Workers {
		k8sNodes = append(k8sNodes, w.NodeID)
	}

	k8sCluster.NodeDeploymentID = map[uint32]uint64{}

	for _, k8sNode := range k8sNodes {
		for _, newDl := range newDl[k8sNode] {
			dlData, err := workloads.ParseDeploymentData(newDl.Metadata)
			if err != nil {
				return errors.Wrapf(err, "could not get deployment %d data", newDl.ContractID)
			}

			if dlData.Name == k8sCluster.Master.Name {
				k8sCluster.NodeDeploymentID[k8sCluster.Master.NodeID] = newDl.ContractID
			}
		}
	}

	if contractID, ok := k8sCluster.NodeDeploymentID[k8sCluster.Master.NodeID]; ok && contractID != 0 {
		if !workloads.Contains(d.tfPluginClient.State.CurrentNodeDeployments[k8sCluster.Master.NodeID], contractID) {
			d.tfPluginClient.State.CurrentNodeDeployments[k8sCluster.Master.NodeID] = append(d.tfPluginClient.State.CurrentNodeDeployments[k8sCluster.Master.NodeID], contractID)
		}
		for _, w := range k8sCluster.Workers {
			if !workloads.Contains(d.tfPluginClient.State.CurrentNodeDeployments[w.NodeID], k8sCluster.NodeDeploymentID[w.NodeID]) {
				d.tfPluginClient.State.CurrentNodeDeployments[w.NodeID] = append(d.tfPluginClient.State.CurrentNodeDeployments[w.NodeID], k8sCluster.NodeDeploymentID[w.NodeID])
			}
		}
	}

	return nil
}

// UpdateFromRemote update a k8s cluster
func (d *K8sDeployer) UpdateFromRemote(ctx context.Context, k8sCluster *workloads.K8sCluster) error {
	if err := d.removeDeletedContracts(k8sCluster); err != nil {
		return d.tfPluginClient.sentry.error(errors.Wrap(err, "failed to remove deleted contracts"))
	}
	currentDeployments, err := d.deployer.GetDeployments(ctx, k8sCluster.NodeDeploymentID)
	if err != nil {
		return d.tfPluginClient.sentry.error(errors.Wrap(err, "failed to fetch remote deployments"))
	}
	zerolog.Debug().Msg("calling updateFromRemote")

	keyUpdated, tokenUpdated, networkUpdated := false, false, false
	// calculate k's properties from the currently deployed deployments
	for _, dl := range currentDeployments {
		for _, w := range dl.Workloads {
			if w.Type == zosTypes.ZMachineType {
				d, err := w.Workload3().WorkloadData()
				if err != nil {
					zerolog.Error().Err(err).Msg("failed to get workload data")
				}
				SSHKey := d.(*zos.ZMachine).Env["SSH_KEY"]
				token := d.(*zos.ZMachine).Env["K3S_TOKEN"]
				networkName := string(d.(*zos.ZMachine).Network.Interfaces[0].Network)
				if !keyUpdated && SSHKey != k8sCluster.SSHKey {
					k8sCluster.SSHKey = SSHKey
					keyUpdated = true
				}
				if !tokenUpdated && token != k8sCluster.Token {
					k8sCluster.Token = token
					tokenUpdated = true
				}
				if !networkUpdated && networkName != k8sCluster.NetworkName {
					k8sCluster.NetworkName = networkName
					networkUpdated = true
				}
			}
		}
	}

	nodeDeploymentID := make(map[uint32]uint64)
	for node, dl := range currentDeployments {
		nodeDeploymentID[node] = dl.ContractID
	}
	k8sCluster.NodeDeploymentID = nodeDeploymentID
	// maps from workload name to (public ip, node id, disk size, actual workload)
	workloadNodeID := make(map[string]uint32)
	workloadDiskSize := make(map[string]uint64)
	workloadComputedIP := make(map[string]string)
	workloadComputedIP6 := make(map[string]string)
	workloadObj := make(map[string]gridtypes.Workload)

	publicIPs := make(map[string]string)
	publicIP6s := make(map[string]string)
	diskSize := make(map[string]uint64)
	for node, dl := range currentDeployments {
		for _, w := range dl.Workloads {
			if w.Type == zosTypes.ZMachineType {
				workloadNodeID[w.Name] = node
				workloadObj[w.Name] = *w.Workload3()

			} else if w.Type == zosTypes.PublicIPType {
				ipResult := zos.PublicIPResult{}
				if err := json.Unmarshal(w.Result.Data, &d); err != nil {
					return d.tfPluginClient.sentry.error(errors.Wrap(err, "failed to load public ip data"))
				}
				publicIPs[w.Name] = ipResult.IP.String()
				publicIP6s[w.Name] = ipResult.IPv6.String()
			} else if w.Type == zosTypes.ZMountType {
				wl, err := w.Workload3().WorkloadData()
				if err != nil {
					return d.tfPluginClient.sentry.error(errors.Wrap(err, "failed to load disk data"))
				}
				diskSize[w.Name] = uint64(wl.(*zos.ZMount).Size) / zosTypes.Gigabyte
			}
		}
	}
	for _, dl := range currentDeployments {
		for _, w := range dl.Workloads {
			if w.Type == zosTypes.ZMachineType {
				publicIPKey := fmt.Sprintf("%sip", w.Name)
				diskKey := fmt.Sprintf("%sdisk", w.Name)
				workloadDiskSize[w.Name] = diskSize[diskKey]
				workloadComputedIP[w.Name] = publicIPs[publicIPKey]
				workloadComputedIP6[w.Name] = publicIP6s[publicIPKey]
			}
		}
	}
	// update master
	masterNodeID, ok := workloadNodeID[k8sCluster.Master.Name]
	if !ok {
		k8sCluster.Master = nil
	} else {
		masterWorkload := workloadObj[k8sCluster.Master.Name]
		masterIP := workloadComputedIP[k8sCluster.Master.Name]
		masterIP6 := workloadComputedIP6[k8sCluster.Master.Name]
		masterDiskSize := workloadDiskSize[k8sCluster.Master.Name]

		m, err := workloads.NewK8sNodeFromWorkload(masterWorkload, masterNodeID, masterDiskSize, masterIP, masterIP6)
		if err != nil {
			return d.tfPluginClient.sentry.error(errors.Wrap(err, "failed to get master node from workload"))
		}
		k8sCluster.Master = &m
	}
	// update workers
	workers := make([]workloads.K8sNode, 0)
	for _, w := range k8sCluster.Workers {
		workerNodeID, ok := workloadNodeID[w.Name]
		if !ok {
			// worker doesn't exist in any deployment, skip it
			continue
		}
		delete(workloadNodeID, w.Name)
		workerWorkload := workloadObj[w.Name]
		workerIP := workloadComputedIP[w.Name]
		workerIP6 := workloadComputedIP6[w.Name]

		workerDiskSize := workloadDiskSize[w.Name]
		w, err := workloads.NewK8sNodeFromWorkload(workerWorkload, workerNodeID, workerDiskSize, workerIP, workerIP6)
		if err != nil {
			return d.tfPluginClient.sentry.error(errors.Wrap(err, "failed to get worker data from workload"))
		}
		workers = append(workers, w)
	}
	// add missing workers (in case of failed deletions)
	for name, workerNodeID := range workloadNodeID {
		if name == k8sCluster.Master.Name {
			continue
		}
		workerWorkload := workloadObj[name]
		workerIP := workloadComputedIP[name]
		workerIP6 := workloadComputedIP6[name]
		workerDiskSize := workloadDiskSize[name]
		w, err := workloads.NewK8sNodeFromWorkload(workerWorkload, workerNodeID, workerDiskSize, workerIP, workerIP6)
		if err != nil {
			return d.tfPluginClient.sentry.error(errors.Wrap(err, "failed to get worker data from workload"))
		}
		workers = append(workers, w)
	}
	k8sCluster.Workers = workers
	zerolog.Debug().Msg("after updateFromRemote\n")
	enc := json.NewEncoder(log.Writer())
	enc.SetIndent("", "  ")
	err = enc.Encode(k8sCluster)
	if err != nil {
		return d.tfPluginClient.sentry.error(errors.Wrap(err, "failed to encode k8s deployer"))
	}

	return nil
}

func (d *K8sDeployer) removeDeletedContracts(k8sCluster *workloads.K8sCluster) error {
	sub := d.tfPluginClient.SubstrateConn
	nodeDeploymentID := make(map[uint32]uint64)
	for nodeID, deploymentID := range k8sCluster.NodeDeploymentID {
		cont, err := sub.GetContract(deploymentID)
		if err != nil {
			return errors.Wrap(err, "failed to get deployments")
		}
		if !cont.IsDeleted() {
			nodeDeploymentID[nodeID] = deploymentID
		}
	}
	k8sCluster.NodeDeploymentID = nodeDeploymentID
	return nil
}

// TODO: integrate new list private ips function
func (d *K8sDeployer) getK8sUsedIPs(k8s *workloads.K8sCluster) map[uint32][]byte {
	usedIPs := make(map[uint32][]byte)

	if k8s.Master.IP != "" {
		ip := net.ParseIP(k8s.Master.IP).To4()
		if ip != nil {
			usedIPs[k8s.Master.NodeID] = append(usedIPs[k8s.Master.NodeID], ip[3])
		}
	}

	for _, w := range k8s.Workers {
		if w.IP != "" {
			ip := net.ParseIP(w.IP).To4()
			if ip != nil {
				usedIPs[w.NodeID] = append(usedIPs[w.NodeID], ip[3])
			}
		}
	}

	return usedIPs
}

func (d *K8sDeployer) getK8sFreeIP(ipRange gridtypes.IPNet, nodeID uint32, k8s *workloads.K8sCluster) (string, error) {
	nodeUsedIPs := d.getK8sUsedIPs(k8s)

	ip := ipRange.IP.To4()
	if ip == nil {
		return "", errors.Errorf("the provided ip range (%s) is not a valid ipv4", ipRange.String())
	}

	for i := 2; i < 255; i++ {
		hostID := byte(i)
		if !workloads.Contains(nodeUsedIPs[nodeID], hostID) {
			nodeUsedIPs[nodeID] = append(nodeUsedIPs[nodeID], hostID)
			ip[3] = hostID
			return ip.String(), nil
		}
	}
	return "", errors.New("all ips are used")
}

func (d *K8sDeployer) assignNodesIPs(k8sCluster *workloads.K8sCluster) error {
	masterNodeRange := k8sCluster.NodesIPRange[k8sCluster.Master.NodeID]
	if k8sCluster.Master.IP == "" || !masterNodeRange.Contains(net.ParseIP(k8sCluster.Master.IP)) {
		ip, err := d.getK8sFreeIP(masterNodeRange, k8sCluster.Master.NodeID, k8sCluster)
		if err != nil {
			return errors.Wrap(err, "failed to find free ip for master")
		}
		k8sCluster.Master.IP = ip
	}
	for idx, w := range k8sCluster.Workers {
		workerNodeRange := k8sCluster.NodesIPRange[w.NodeID]
		if w.IP != "" && workerNodeRange.Contains(net.ParseIP(w.IP)) {
			continue
		}
		ip, err := d.getK8sFreeIP(workerNodeRange, w.NodeID, k8sCluster)
		if err != nil {
			return errors.Wrap(err, "failed to find free ip for worker")
		}
		k8sCluster.Workers[idx].IP = ip
	}
	return nil
}

func assignNodesFlistsAndEntryPoints(k *workloads.K8sCluster) {
	if k.Flist == "" {
		k.Flist = k.Master.Flist
	}
	if k.Entrypoint == "" {
		if k.Master.Entrypoint != "" {
			k.Entrypoint = k.Master.Entrypoint
		} else {
			k.Entrypoint = "/sbin/zinit init" // set default value
		}
	}

	k.Master.Flist = k.Flist
	k.Master.Entrypoint = k.Entrypoint
	for i := range k.Workers {
		k.Workers[i].Flist = k.Flist
		k.Workers[i].Entrypoint = k.Entrypoint
	}
}
