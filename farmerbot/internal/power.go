package internal

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
)

// powerOn sets the node power state ON
func (f *FarmerBot) powerOn(sub Substrate, nodeID uint32) error {
	log.Info().Uint32("nodeID", nodeID).Msg("POWER ON")
	f.m.Lock()
	defer f.m.Unlock()

	node, ok := f.nodes[nodeID]
	if !ok {
		return fmt.Errorf("node %d is not found", nodeID)
	}

	if node.powerState == on || node.powerState == wakingUp {
		return nil
	}

	_, err := sub.SetNodePowerTarget(f.identity, nodeID, true)
	if err != nil {
		return fmt.Errorf("failed to set node %d power target to up with error: %w", nodeID, err)
	}

	node.powerState = wakingUp
	node.lastTimeAwake = time.Now()
	node.lastTimePowerStateChanged = time.Now()

	f.nodes[nodeID] = node
	return nil
}

// powerOff sets the node power state OFF
func (f *FarmerBot) powerOff(sub Substrate, nodeID uint32) error {
	log.Info().Uint32("nodeID", nodeID).Msg("POWER OFF")
	f.m.Lock()
	defer f.m.Unlock()

	node, ok := f.nodes[nodeID]
	if !ok {
		return fmt.Errorf("node '%d' is not found", nodeID)
	}

	if node.powerState == off || node.powerState == shuttingDown {
		return nil
	}

	if node.neverShutDown {
		return fmt.Errorf("cannot power off node '%d', node is configured to never be shutdown", nodeID)
	}

	if node.PublicConfig.HasValue {
		return fmt.Errorf("cannot power off node '%d', node has public config", nodeID)
	}

	if node.timeoutClaimedResources.After(time.Now()) {
		return fmt.Errorf("cannot power off node '%d', node has claimed resources", nodeID)
	}

	if node.hasActiveRentContract {
		return fmt.Errorf("cannot power off node '%d', node has a rent contract", nodeID)
	}

	if node.hasActiveContracts {
		return fmt.Errorf("cannot power off node '%d', node has active contracts", nodeID)
	}

	if !node.isUnused() {
		return fmt.Errorf("cannot power off node '%d', node is used", nodeID)
	}

	if time.Since(node.lastTimePowerStateChanged) < periodicWakeUpDuration {
		return fmt.Errorf("cannot power off node '%d', node is still in its wakeup duration", nodeID)
	}

	onNodes := f.filterNodesPower([]powerState{on})

	if len(onNodes) < 2 {
		return fmt.Errorf("cannot power off node '%d', at least one node should be on in the farm", nodeID)
	}

	_, err := sub.SetNodePowerTarget(f.identity, nodeID, false)
	if err != nil {
		powerTarget, getErr := sub.GetPowerTarget(nodeID)
		if getErr != nil {
			return fmt.Errorf("failed to get node '%d' power target with error: %w", nodeID, getErr)
		}

		if powerTarget.State.IsDown || powerTarget.Target.IsDown {
			log.Warn().Uint32("nodeID", nodeID).Msg("Node is shutting down although it failed to set power target in tfchain")
			node.powerState = shuttingDown
			node.lastTimePowerStateChanged = time.Now()
			f.nodes[nodeID] = node
		}

		return fmt.Errorf("failed to set node '%d' power target to down with error: %w", nodeID, err)
	}

	node.powerState = shuttingDown
	node.lastTimePowerStateChanged = time.Now()

	f.nodes[nodeID] = node
	return nil
}

// manageNodesPower for power management nodes
func (f *FarmerBot) manageNodesPower(sub Substrate) error {
	nodes := f.filterNodesPower([]powerState{on, wakingUp})

	usedResources, totalResources := calculateResourceUsage(nodes)
	if totalResources == 0 {
		return nil
	}

	resourceUsage := 100 * float32(usedResources) / float32(totalResources)
	if resourceUsage >= float32(f.config.Power.WakeUpThreshold) {
		log.Info().Msgf("Too high resource usage = %.1f%%, threshold = %d%%", resourceUsage, f.config.Power.WakeUpThreshold)
		return f.resourceUsageTooHigh(sub)
	}

	log.Info().Msgf("Too low resource usage = %.1f%%, threshold = %d%%", resourceUsage, f.config.Power.WakeUpThreshold)
	return f.resourceUsageTooLow(sub, usedResources, totalResources)
}

func calculateResourceUsage(nodes map[uint32]node) (uint64, uint64) {
	usedResources := capacity{}
	totalResources := capacity{}

	for _, node := range nodes {
		if node.hasActiveRentContract {
			usedResources.add(node.resources.total)
		} else {
			usedResources.add(node.resources.used)
		}
		totalResources.add(node.resources.total)
	}

	used := usedResources.cru + usedResources.hru + usedResources.mru + usedResources.sru
	total := totalResources.cru + totalResources.hru + totalResources.mru + totalResources.sru

	return used, total
}

func (f *FarmerBot) resourceUsageTooHigh(sub Substrate) error {
	nodesKeys := make([]uint32, 0, len(f.nodes))
	for k := range f.nodes {
		nodesKeys = append(nodesKeys, k)
	}

	for i := 0; i < len(f.nodes); i++ {
		nodeID := nodesKeys[i]
		node := f.nodes[nodeID]

		if node.powerState == off {
			return f.powerOn(sub, nodeID)
		}
	}

	return fmt.Errorf("no available node to wake up, resources usage is high")
}

func (f *FarmerBot) resourceUsageTooLow(sub Substrate, usedResources, totalResources uint64) error {
	onNodes := f.filterNodesPower([]powerState{on})

	// nodes with public config can't be shutdown
	// Do not shutdown a node that just came up (give it some time `periodicWakeUpDuration`)
	nodesAllowedToShutdown := f.filterAllowedNodesToShutDown()

	if len(onNodes) <= 1 {
		log.Debug().Msg("Nothing to shutdown")
		return nil
	}

	if len(nodesAllowedToShutdown) == 0 {
		log.Debug().Msg("No nodes are allowed to shutdown")
		return nil
	}

	log.Debug().Uints32("nodes IDs", maps.Keys(nodesAllowedToShutdown)).Msg("Nodes allowed to shutdown")

	newUsedResources := usedResources
	newTotalResources := totalResources
	nodesLeftOnline := len(onNodes)

	// use keys to keep nodes order
	nodesAllowedToShutdownKeys := make([]uint32, 0, len(nodesAllowedToShutdown))
	for k := range nodesAllowedToShutdown {
		nodesAllowedToShutdownKeys = append(nodesAllowedToShutdownKeys, k)
	}

	// shutdown a node if there is more than an unused node (aka keep at least one node online)
	for i := 0; i < len(nodesAllowedToShutdown); i++ {
		node := nodesAllowedToShutdown[nodesAllowedToShutdownKeys[i]]

		if nodesLeftOnline == 1 {
			break
		}
		nodesLeftOnline -= 1
		newUsedResources -= node.resources.used.hru + node.resources.used.sru +
			node.resources.used.mru + node.resources.used.cru
		newTotalResources -= node.resources.total.hru + node.resources.total.sru +
			node.resources.total.mru + node.resources.total.cru

		if newTotalResources == 0 {
			break
		}

		newResourceUsage := 100 * float32(newUsedResources) / float32(newTotalResources)
		if newResourceUsage < float32(f.config.Power.WakeUpThreshold) {
			// we need to keep the resource percentage lower then the threshold
			log.Info().Uint32("nodeID", uint32(node.ID)).Msgf("Too low resource usage = %.1f%%. Turning off unused node", newResourceUsage)
			err := f.powerOff(sub, uint32(node.ID))
			if err != nil {
				log.Error().Err(err).Uint32("nodeID", uint32(node.ID)).Msg("Failed to power off node")

				if node.powerState == shuttingDown {
					continue
				}

				nodesLeftOnline += 1
				newUsedResources += node.resources.used.hru + node.resources.used.sru +
					node.resources.used.mru + node.resources.used.cru
				newTotalResources += node.resources.total.hru + node.resources.total.sru +
					node.resources.total.mru + node.resources.total.cru
			}
		}
	}

	return nil
}
