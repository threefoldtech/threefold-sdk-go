package indexer

import (
	"context"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/mocks"
	types "github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/pkg/types"
)

func TestNewSpeedWork(t *testing.T) {
	wanted := &SpeedWork{
		findersInterval: map[string]time.Duration{
			"up": 2 * time.Minute,
		},
	}
	speed := NewSpeedWork(2)
	assert.Exactlyf(t, wanted, speed, "got: %v , expected: %v", speed, wanted)
}

func TestSpeedGet(t *testing.T) {
	speed := NewSpeedWork(2)
	ctrl := gomock.NewController(t)
	t.Run("get speed with valid twin id", func(t *testing.T) {
		twinID := uint32(1)
		expected := TaskResult{
			Name: "iperf",
			Result: []IperfResult{{
				UploadSpeed:   float64(200),
				DownloadSpeed: float64(100),
				NodeID:        uint32(1),
				CpuReport: CPUUtilizationPercent{
					HostTotal:    float64(3),
					HostUser:     float64(1),
					HostSystem:   float64(2),
					RemoteTotal:  float64(3),
					RemoteUser:   float64(1),
					RemoteSystem: float64(2),
				},
			},
			},
		}
		client := mocks.NewMockClient(ctrl)
		client.EXPECT().Call(gomock.Any(), twinID, perfTestCallCmd, gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, twin uint32, fn string, data, result interface{}) error {
				*(result.(*TaskResult)) = expected
				return nil
			},
		)
		got, err := speed.Get(context.Background(), client, twinID)
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.IsTypef(t, got, []types.Speed{}, "got: %v , expected: %v", got, []types.Speed{})
		assert.Exactlyf(t, got[0].Upload, expected.Result.([]IperfResult)[0].UploadSpeed, "got: %v , expected: %v", got[0].Upload, expected.Result.([]IperfResult)[0].UploadSpeed)
		assert.Exactlyf(t, got[0].Download, expected.Result.([]IperfResult)[0].DownloadSpeed, "got: %v , expected: %v", got[0].Download, expected.Result.([]IperfResult)[0].DownloadSpeed)
		assert.Exactlyf(t, got[0].NodeTwinId, twinID, "got: %v , expected: %v", got[0].NodeTwinId, twinID)
	})

	t.Run("get speed with invalid twin id", func(t *testing.T) {
		twinID := uint32(1)
		client := mocks.NewMockClient(ctrl)
		client.EXPECT().Call(gomock.Any(), twinID, perfTestCallCmd, gomock.Any(), gomock.Any()).Return(assert.AnError)
		got, err := speed.Get(context.Background(), client, twinID)
		assert.Error(t, err)
		assert.Len(t, got, 0)
	})
}
