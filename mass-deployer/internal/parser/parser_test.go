package parser

import (
	"os"
	"path"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	deployer "github.com/threefoldtech/tfgrid-sdk-go/mass-deployer/pkg/mass-deployer"
	"gopkg.in/yaml.v3"
)

func TestParseConfig(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	tempDir := t.TempDir()
	configPath := path.Join(tempDir, "conig.yaml")
	validFlist := "https://hub.grid.tf/tf-official-apps/base:latest.flist"
	confStruct := deployer.Config{
		NodeGroups: []deployer.NodesGroup{
			{
				Name:       "group_a",
				NodesCount: 5,
				FreeCPU:    2,
				FreeMRU:    256,
				FreeSRU:    50,
				FreeHRU:    50,
				Pubip4:     true,
				Regions:    "Europe",
			},
		},
		Vms: []deployer.Vms{
			{
				Name:       "example-vm",
				Count:      4,
				Nodegroup:  "group_a",
				FreeCPU:    2,
				FreeMRU:    256,
				PublicIP4:  true,
				Flist:      validFlist,
				Entrypoint: "/sbin/zinit init",
				SSHKey:     "example1",
			},
		},
		SSHKeys: map[string]string{
			"example1": "example ssh key",
		},
		Mnemonic: "rival oyster defense garbage fame disease mask mail family wire village vibrant index fuel dolphin",
		Network:  "dev",
	}
	t.Run("invalid yaml format", func(t *testing.T) {
		conf := ` {
  "node_groups": [
    {
      "nodes_count": 10
    }
} `
		err := os.WriteFile(configPath, []byte(conf), 0667)
		assert.NoError(t, err)

		configFile, err := os.Open(configPath)
		assert.NoError(t, err)

		_, err = ParseConfig(configFile)
		assert.Error(t, err)
	})

	t.Run("invalid node_group", func(t *testing.T) {
		conf := confStruct
		conf.NodeGroups = []deployer.NodesGroup{}

		data, err := yaml.Marshal(conf)
		assert.NoError(t, err)

		err = os.WriteFile(configPath, data, 0667)
		assert.NoError(t, err)

		configFile, err := os.Open(configPath)
		assert.NoError(t, err)

		_, err = ParseConfig(configFile)
		assert.Error(t, err)
	})

	t.Run("invalid vms", func(t *testing.T) {
		conf := confStruct
		conf.Vms = []deployer.Vms{}

		data, err := yaml.Marshal(conf)
		assert.NoError(t, err)

		err = os.WriteFile(configPath, data, 0667)
		assert.NoError(t, err)

		configFile, err := os.Open(configPath)
		assert.NoError(t, err)

		_, err = ParseConfig(configFile)
		assert.Error(t, err)
	})

	t.Run("invalid ssh_keys", func(t *testing.T) {
		conf := confStruct
		conf.SSHKeys = map[string]string{}

		data, err := yaml.Marshal(conf)
		assert.NoError(t, err)

		err = os.WriteFile(configPath, data, 0667)
		assert.NoError(t, err)

		configFile, err := os.Open(configPath)
		assert.NoError(t, err)

		_, err = ParseConfig(configFile)
		assert.Error(t, err)
	})

	t.Run("invalid mnemonic", func(t *testing.T) {
		conf := confStruct
		conf.Mnemonic = ""

		data, err := yaml.Marshal(conf)
		assert.NoError(t, err)

		err = os.WriteFile(configPath, data, 0667)
		assert.NoError(t, err)

		configFile, err := os.Open(configPath)
		assert.NoError(t, err)

		_, err = ParseConfig(configFile)
		assert.Error(t, err)
	})

	t.Run("invalid network", func(t *testing.T) {
		conf := confStruct
		conf.Network = ""

		data, err := yaml.Marshal(conf)
		assert.NoError(t, err)

		err = os.WriteFile(configPath, data, 0667)
		assert.NoError(t, err)

		configFile, err := os.Open(configPath)
		assert.NoError(t, err)

		_, err = ParseConfig(configFile)
		assert.Error(t, err)
	})

	t.Run("empty flist", func(t *testing.T) {
		conf := confStruct
		conf.Vms[0].Flist = ""

		data, err := yaml.Marshal(conf)
		assert.NoError(t, err)

		conf.Vms[0].Flist = validFlist

		err = os.WriteFile(configPath, data, 0667)
		assert.NoError(t, err)

		configFile, err := os.Open(configPath)
		assert.NoError(t, err)

		_, err = ParseConfig(configFile)
		assert.Error(t, err)
	})

	t.Run("invalid flist extension", func(t *testing.T) {
		conf := confStruct
		conf.Vms[0].Flist = "https://example-list.list"

		data, err := yaml.Marshal(conf)
		assert.NoError(t, err)

		conf.Vms[0].Flist = validFlist

		err = os.WriteFile(configPath, data, 0667)
		assert.NoError(t, err)

		configFile, err := os.Open(configPath)
		assert.NoError(t, err)

		_, err = ParseConfig(configFile)
		assert.Error(t, err)
	})

	t.Run("invalid md5", func(t *testing.T) {
		conf := confStruct
		conf.Vms[0].Flist = "https://example-flist.flist"

		data, err := yaml.Marshal(conf)
		assert.NoError(t, err)

		conf.Vms[0].Flist = validFlist

		err = os.WriteFile(configPath, data, 0667)
		assert.NoError(t, err)

		configFile, err := os.Open(configPath)
		assert.NoError(t, err)

		_, err = ParseConfig(configFile)
		assert.Error(t, err)
	})

	t.Run("zero cpu in node group", func(t *testing.T) {
		conf := confStruct
		conf.NodeGroups[0].FreeCPU = 0

		data, err := yaml.Marshal(conf)
		assert.NoError(t, err)

		conf.NodeGroups[0].FreeCPU = 2

		err = os.WriteFile(configPath, data, 0667)
		assert.NoError(t, err)

		configFile, err := os.Open(configPath)
		assert.NoError(t, err)

		_, err = ParseConfig(configFile)
		assert.Error(t, err)
	})

	t.Run("zero memory in node group", func(t *testing.T) {
		conf := confStruct
		conf.NodeGroups[0].FreeMRU = 0

		data, err := yaml.Marshal(conf)
		assert.NoError(t, err)

		conf.NodeGroups[0].FreeMRU = 256

		err = os.WriteFile(configPath, data, 0667)
		assert.NoError(t, err)

		configFile, err := os.Open(configPath)
		assert.NoError(t, err)

		_, err = ParseConfig(configFile)
		assert.Error(t, err)
	})

	t.Run("zero cpu in vm", func(t *testing.T) {
		conf := confStruct
		conf.Vms[0].FreeCPU = 0

		data, err := yaml.Marshal(conf)
		assert.NoError(t, err)

		conf.Vms[0].FreeCPU = 2

		err = os.WriteFile(configPath, data, 0667)
		assert.NoError(t, err)

		configFile, err := os.Open(configPath)
		assert.NoError(t, err)

		_, err = ParseConfig(configFile)
		assert.Error(t, err)
	})

	t.Run("zero memory in vm", func(t *testing.T) {
		conf := confStruct
		conf.Vms[0].FreeMRU = 0

		data, err := yaml.Marshal(conf)
		assert.NoError(t, err)

		conf.Vms[0].FreeMRU = 256

		err = os.WriteFile(configPath, data, 0667)
		assert.NoError(t, err)

		configFile, err := os.Open(configPath)
		assert.NoError(t, err)

		_, err = ParseConfig(configFile)
		assert.Error(t, err)
	})

	t.Run("cpu exeed limit in vm", func(t *testing.T) {
		conf := confStruct
		conf.Vms[0].FreeCPU = 35

		data, err := yaml.Marshal(conf)
		assert.NoError(t, err)

		conf.Vms[0].FreeCPU = 2

		err = os.WriteFile(configPath, data, 0667)
		assert.NoError(t, err)

		configFile, err := os.Open(configPath)
		assert.NoError(t, err)

		_, err = ParseConfig(configFile)
		assert.Error(t, err)
	})

	t.Run("memory exeed limit in vm", func(t *testing.T) {
		conf := confStruct
		conf.Vms[0].FreeMRU = 300000

		data, err := yaml.Marshal(conf)
		assert.NoError(t, err)

		conf.Vms[0].FreeMRU = 256

		err = os.WriteFile(configPath, data, 0667)
		assert.NoError(t, err)

		configFile, err := os.Open(configPath)
		assert.NoError(t, err)

		_, err = ParseConfig(configFile)
		assert.Error(t, err)
	})

	t.Run("root size exeed limit in vm", func(t *testing.T) {
		conf := confStruct
		conf.Vms[0].Rootsize = 20000

		data, err := yaml.Marshal(conf)
		assert.NoError(t, err)

		conf.Vms[0].Rootsize = 0

		err = os.WriteFile(configPath, data, 0667)
		assert.NoError(t, err)

		configFile, err := os.Open(configPath)
		assert.NoError(t, err)

		_, err = ParseConfig(configFile)
		assert.Error(t, err)
	})

	t.Run("valid config", func(t *testing.T) {
		conf := confStruct

		data, err := yaml.Marshal(conf)
		assert.NoError(t, err)

		err = os.WriteFile(configPath, data, 0667)
		assert.NoError(t, err)

		configFile, err := os.Open(configPath)
		assert.NoError(t, err)

		_, err = ParseConfig(configFile)
		assert.NoError(t, err)
	})
}
