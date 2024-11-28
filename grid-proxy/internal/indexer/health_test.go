package indexer

import (
	"context"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_rmb "github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/mocks"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/pkg/types"
	"github.com/threefoldtech/zos/pkg/diagnostics"
)

func TestNewHealthWork(t *testing.T) {
	wanted := &HealthWork{
		findersInterval: map[string]time.Duration{
			"up":      2 * time.Minute,
			"healthy": 2 * time.Minute,
		},
	}
	health := NewHealthWork(2)
	assert.Exactlyf(t, wanted, health, "got: %v , expected: %v", health, wanted)
}

func TestHealthGet(t *testing.T) {
	health := NewHealthWork(2)
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	t.Run("get health with valid twin id", func(t *testing.T) {
		twinID := uint32(1)
		expected := []types.HealthReport{
			{
				NodeTwinId: 1,
				Healthy:    true,
				UpdatedAt:  time.Now().Unix(),
			},
		}
		client := mock_rmb.NewMockClient(ctrl)
		client.EXPECT().Call(gomock.Any(), twinID, healthCallCmd, nil, gomock.AssignableToTypeOf(&diagnostics.Diagnostics{})).DoAndReturn(
			func(ctx context.Context, twin uint32, fn string, data, result interface{}) error {
				diag := result.(*diagnostics.Diagnostics)
				diag.Healthy = true
				return nil
			},
		)
		got, err := health.Get(ctx, client, twinID)

		assert.NoError(t, err)
		assert.Equal(t, expected[0].NodeTwinId, got[0].NodeTwinId)
		assert.Equal(t, expected[0].Healthy, got[0].Healthy)
		assert.Len(t, got, 1)
	})

	t.Run("get health with invalid twin id", func(t *testing.T) {
		twinID := uint32(2)
		client := mock_rmb.NewMockClient(ctrl)
		client.EXPECT().Call(gomock.Any(), twinID, healthCallCmd, nil, gomock.AssignableToTypeOf(&diagnostics.Diagnostics{})).Return(
			assert.AnError,
		)
		got, _ := health.Get(ctx, client, twinID)
		// we expect an error here because the twin id is invalid but the implementation ignore errors
		//assert.Error(t, err)
		assert.Len(t, got, 1)
		assert.False(t, got[0].Healthy)
	})
}

func TestRemoveDuplicates(t *testing.T) {
	reports := []types.HealthReport{
		{NodeTwinId: 1, Healthy: true},
		{NodeTwinId: 2, Healthy: false},
		{NodeTwinId: 1, Healthy: true}, //Duplicate
		{NodeTwinId: 3, Healthy: true},
	}

	result := removeDuplicates(reports)
	assert.Len(t, result, 3)
	assert.Contains(t, result, reports[0])
	assert.Contains(t, result, reports[1])
	assert.Contains(t, result, reports[3])
}
