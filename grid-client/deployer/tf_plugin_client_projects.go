package deployer

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// CancelByProjectName cancels a deployed project
func (t *TFPluginClient) CancelByProjectName(projectName string) error {
	log.Info().Str("project name", projectName).Msg("canceling contracts")
	contracts, err := t.ContractsGetter.ListContractsOfProjectName(projectName)
	if err != nil {
		return errors.Wrapf(err, "could not load contracts for project %s", projectName)
	}
	contractIDS := make([]uint64, 0)

	contractsSlice := append(contracts.NameContracts, contracts.NodeContracts...)
	for _, contract := range contractsSlice {
		contractID, err := strconv.ParseUint(contract.ContractID, 0, 64)
		if err != nil {
			return errors.Wrapf(err, "could not parse contract %s into uint64", contract.ContractID)
		}
		contractIDS = append(contractIDS, contractID)
	}
	if err := t.BatchCancelContract(contractIDS); err != nil {
		return fmt.Errorf("failed to cancel contracts for project %s: %w", projectName, err)
	}

	log.Info().Str("project name", projectName).Msg("project is canceled")
	return nil
}
