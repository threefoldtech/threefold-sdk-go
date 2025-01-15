package indexer

import (
	"context"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_rmb "github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/mocks"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/pkg/types"
)

func TestNewGPUWork(t *testing.T) {
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
	gpu := NewGPUWork(2)
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	t.Run("get gpu with valid twin id", func(t *testing.T) {
		twinID := uint32(1)
		expected := []types.NodeGPU{
			{
				NodeTwinID: 1,
				ID:         "gpu-1",
				Vendor:     "NVIDIA",
				Device:     "RTX 3080",
				Contract:   123,
				UpdatedAt:  1234567890,
			},
		}
		client := mock_rmb.NewMockClient(ctrl)
		client.EXPECT().Call(gomock.Any(), twinID, gpuListCmd, nil, gomock.AssignableToTypeOf(&[]types.NodeGPU{})).DoAndReturn(
			func(ctx context.Context, twin uint32, fn string, data, result interface{}) error {
				*(result.(*[]types.NodeGPU)) = expected
				return nil
			},
		)
		got, err := gpu.Get(ctx, client, twinID)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		assert.Len(t, got, 1)
		assert.Equal(t, "gpu-1", got[0].ID)
		assert.Equal(t, "NVIDIA", got[0].Vendor)
	})

	t.Run("get gpu with invalid twin id", func(t *testing.T) {
		twinID := uint32(2)
		client := mock_rmb.NewMockClient(ctrl)
		client.EXPECT().Call(gomock.Any(), twinID, gpuListCmd, nil, gomock.AssignableToTypeOf(&[]types.NodeGPU{})).Return(
			assert.AnError,
		)
		_, err := gpu.Get(ctx, client, twinID)
		assert.Error(t, err)
	})
}
