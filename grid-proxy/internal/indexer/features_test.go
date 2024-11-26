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

func TestNewFeatureWork(t *testing.T) {
	wanted := &FeatureWork{
		findersInterval: map[string]time.Duration{
			"up":  time.Duration(2) * time.Minute,
			"new": newNodesCheckInterval,
		},
	}
	feature := NewFeatureWork(2)
	assert.Exactlyf(t, wanted, feature, "got: %v , expected: %v", feature, wanted)
}

func TestFeatureGet(t *testing.T) {
	feature := NewFeatureWork(2)
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	t.Run("get feature with valid twin id", func(t *testing.T) {
		twinID := uint32(1)
		expected := []string{"zmount", "network", "zdb", "zmachine", "volume", "ipv4",
			"ip",
			"gateway-name-proxy",
			"gateway-fqdn-proxy",
			"qsfs",
			"zlogs",
			"yggdrasil",
			"mycelium",
			"wireguard"}

		client := mock_rmb.NewMockClient(ctrl)
		client.EXPECT().Call(gomock.Any(), twinID, featuresCallCmd, nil, gomock.Any()).DoAndReturn(
			func(ctx context.Context, twinId uint32, fn string, data, result interface{}) error {
				*result.(*[]string) = expected
				return nil
			},
		)
		features, err := feature.Get(ctx, client, twinID)
		assert.NoError(t, err)
		assert.Len(t, features, 1)
		assert.Equal(t, expected, features[0].Features)
		assert.Equal(t, twinID, features[0].NodeTwinId)
		assert.IsTypef(t, features, []types.NodeFeatures{}, "got: %T , expected: %T", features, []types.NodeFeatures{})
	})

	t.Run("get feature with invalid twin id", func(t *testing.T) {
		twinID := uint32(1)

		client := mock_rmb.NewMockClient(ctrl)
		client.EXPECT().Call(gomock.Any(), twinID, featuresCallCmd, nil, gomock.Any()).Return(
			assert.AnError,
		)
		features, err := feature.Get(ctx, client, twinID)
		assert.Error(t, err)
		assert.Len(t, features, 0)
	})

}
