package deployer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	zosTypes "github.com/threefoldtech/tfgrid-sdk-go/grid-client/zos"
)

func TestDeploymentUtils(t *testing.T) {
	tfPluginClient, err := setup()
	assert.NoError(t, err)

	identity := tfPluginClient.Identity
	twinID := tfPluginClient.TwinID

	dl := workloads.NewGridDeployment(twinID, 0, []zosTypes.Workload{})

	dlName, err := deploymentWithNameGateway(identity, twinID, true, 0, backendURLWithTLSPassthrough)
	assert.NoError(t, err)

	t.Run("deployments count public ips", func(t *testing.T) {
		count, err := CountDeploymentPublicIPs(dl)
		assert.NoError(t, err)
		assert.Equal(t, count, uint32(0))
	})

	t.Run("deployments hash", func(t *testing.T) {
		got, err := HashDeployment(dl)
		assert.NoError(t, err)
		assert.NotEmpty(t, got)
	})

	t.Run("deployments workloads hashes", func(t *testing.T) {
		wlHash := "\xa9>\a\xaf\x04\x10\xca\xc1\xac\b\xe1\x177\xf9\xf6D"

		hashes, err := GetWorkloadHashes(dlName)
		assert.NoError(t, err)
		assert.Equal(t, hashes["name"], wlHash)
		assert.Equal(t, len(hashes), 1)
	})

	t.Run("deployments workloads same names", func(t *testing.T) {
		same := SameWorkloadsNames(dl, dlName)
		assert.NoError(t, err)
		assert.Equal(t, same, false)
	})

	t.Run("deployments workloads versions", func(t *testing.T) {
		versions := ConstructWorkloadVersions(dlName)
		assert.Equal(t, versions["name"], uint32(0))
	})

	t.Run("deployments workloads exist", func(t *testing.T) {
		exists := HasWorkload(&dlName, zosTypes.GatewayFQDNProxyType)
		assert.Equal(t, exists, false)

		exists = HasWorkload(&dlName, zosTypes.GatewayNameProxyType)
		assert.Equal(t, exists, true)
	})

	t.Run("deployments capacity", func(t *testing.T) {
		cap, err := Capacity(dlName)
		assert.NoError(t, err)
		assert.Equal(t, cap.CRU, uint64(0))
		assert.Equal(t, cap.SRU, uint64(0))
		assert.Equal(t, cap.MRU, uint64(0))
		assert.Equal(t, cap.HRU, uint64(0))
	})
}
