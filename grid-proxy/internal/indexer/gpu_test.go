package indexer

import (
	"testing"
	"time"
	"context"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_rmb "github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/mocks"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/pkg/types"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/pkg/types"
)




func TestNewGPUWork(t *testing.T){
	wanted := &GPUWork{
		findersInterval: map[string]time.Duration{
			"up":  2 * time.Minute,
			"new": 5 * time.Minute,
		},
	}
	gpu := NewGPUWork(2)
	assert.Exactlyf(t, wanted, gpu, "got: %v , expected: %v", gpu, wanted)
}

func TestGPUGet(t *testing.T) {
	dmi := NewGPUWork(2)
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	t.Run("get gpu with valid twin id", func(t *testing.T) {
		expected := zosGpu
	})
}