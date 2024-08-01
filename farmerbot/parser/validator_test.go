package parser

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/threefoldtech/tfgrid-sdk-go/farmerbot/internal"
	"github.com/threefoldtech/tfgrid-sdk-go/farmerbot/mocks"
)

var testCases = []struct {
	name               string
	includedNodes      []uint32
	priorityNodes      []uint32
	excludedNodes      []uint32
	neverShutdownNodes []uint32
	shouldError        bool
}{
	{
		name:               "test valid include, exclude, priority nodes and never shutdown nodes",
		includedNodes:      []uint32{20, 21, 22, 30, 31, 32, 40, 41},
		priorityNodes:      []uint32{20, 21},
		excludedNodes:      []uint32{23, 24, 34},
		neverShutdownNodes: []uint32{22, 30},
		shouldError:        false,
	},
	{
		name:          "test invalid include",
		includedNodes: []uint32{26, 27},
		shouldError:   true,
	},
	{
		name:          "test invalid priority",
		includedNodes: []uint32{21},
		priorityNodes: []uint32{20, 21},
		shouldError:   true,
	},
	{
		name:               "test invalid never shutdown node",
		includedNodes:      []uint32{21},
		neverShutdownNodes: []uint32{20, 21},
		shouldError:        true,
	}, {
		name:          "test overlapping nodes in include and exclude",
		includedNodes: []uint32{21},
		excludedNodes: []uint32{20, 21},
		shouldError:   true,
	}, {
		name:          "test overlapping nodes in priority and exclude",
		includedNodes: []uint32{21},
		priorityNodes: []uint32{21},
		excludedNodes: []uint32{20, 21},
		shouldError:   true,
	}, {
		name:               "test overlapping nodes in never shutdown and exclude",
		includedNodes:      []uint32{21},
		excludedNodes:      []uint32{20, 21},
		neverShutdownNodes: []uint32{21},
		shouldError:        true,
	}, {
		name:          "test invalid exclude",
		excludedNodes: []uint32{26, 27},
		shouldError:   true,
	}, {
		name:               "test all nodes included and other nodes are valid",
		priorityNodes:      []uint32{21},
		excludedNodes:      []uint32{22},
		neverShutdownNodes: []uint32{20},
		shouldError:        false,
	}, {
		name:          "test all nodes included and invalid priority nodes",
		priorityNodes: []uint32{27, 26},
		shouldError:   true,
	}, {
		name:               "test all nodes included and invalid shutdown nodes",
		neverShutdownNodes: []uint32{27, 26},
		shouldError:        true,
	}, {
		name:               "test all nodes included and overlapping shutdown nodes and excluded",
		neverShutdownNodes: []uint32{21, 20},
		excludedNodes:      []uint32{21, 20},
		shouldError:        true,
	}, {
		name:          "test all nodes included and overlapping priority and excluded nodes",
		excludedNodes: []uint32{21, 20},
		priorityNodes: []uint32{21, 20},
		shouldError:   true,
	},
}
var nodesMap = map[uint32]bool{
	20: true, 21: true, 22: true, 23: true, 24: true, 30: true, 31: true, 32: true, 34: true, 40: true, 41: true,
}

func TestValidateInput(t *testing.T) {
	config := internal.Config{FarmID: uint32(25)}
	ctrl := gomock.NewController(t)
	mockGetNodes := mocks.NewMockSubstrate(ctrl)
	mockGetNodes.EXPECT().GetNodes(config.FarmID).Times(13).Return([]uint32{20, 21, 22, 23, 24, 30, 31, 32, 34, 40, 41}, nil)
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			config.IncludedNodes = test.includedNodes
			config.ExcludedNodes = test.excludedNodes
			config.PriorityNodes = test.priorityNodes
			config.NeverShutDownNodes = test.neverShutdownNodes
			got := validateInput(config, mockGetNodes)
			if test.shouldError {
				assert.Error(t, got)
			} else {
				assert.NoError(t, got)
			}
		})
	}

}

func TestValidateWithIncludedNodes(t *testing.T) {
	config := internal.Config{}
	for i := 0; i < 7; i++ {
		t.Run(testCases[i].name, func(t *testing.T) {
			config.IncludedNodes = testCases[i].includedNodes
			config.PriorityNodes = testCases[i].priorityNodes
			config.NeverShutDownNodes = testCases[i].neverShutdownNodes
			config.ExcludedNodes = testCases[i].excludedNodes

			got := validateWithIncludedNodes(config, nodesMap)
			if testCases[i].shouldError {
				assert.Error(t, got)
			} else {
				assert.NoError(t, got)
			}
		})
	}
}
func TestValidateWithAllNodes(t *testing.T) {
	config := internal.Config{}
	for i := 8; i < 13; i++ {
		t.Run(testCases[i].name, func(t *testing.T) {
			config.IncludedNodes = testCases[i].includedNodes
			config.PriorityNodes = testCases[i].priorityNodes
			config.NeverShutDownNodes = testCases[i].neverShutdownNodes
			config.ExcludedNodes = testCases[i].excludedNodes
			got := validateWithAllNodes(config, nodesMap)
			if testCases[i].shouldError {
				assert.Error(t, got)
			} else {
				assert.NoError(t, got)
			}
		})
	}
}
