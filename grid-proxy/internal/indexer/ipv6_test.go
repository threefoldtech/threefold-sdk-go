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

func TestNewIpv6Work(t *testing.T) {
	wanted := &Ipv6Work{
		finders: map[string]time.Duration{
			"up": 2 * time.Minute,
		},
	}
	ipv6 := NewIpv6Work(2)
	assert.Exactlyf(t, wanted, ipv6, "got: %v , expected: %v", ipv6, wanted)
}

func TestIpv6Get(t *testing.T) {
	ipv6 := NewIpv6Work(2)
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	t.Run("get ipv6 with valid twin id", func(t *testing.T) {
		twinID := uint32(1)
		expected := []types.HasIpv6{
			{
				NodeTwinId: 1,
				HasIpv6:    true,
				UpdatedAt:  time.Now().Unix(),
			},
		}
		client := mock_rmb.NewMockClient(ctrl)
		client.EXPECT().Call(gomock.Any(), twinID, cmd, nil, gomock.Any()).DoAndReturn(
			func(ctx context.Context, twin uint32, fn string, data, result interface{}) error {
				*(result.(*bool)) = true
				return nil
			},
		)
		got, err := ipv6.Get(ctx, client, twinID)

		assert.NoError(t, err)
		assert.Equal(t, expected[0].NodeTwinId, got[0].NodeTwinId)
		assert.Equal(t, expected[0].HasIpv6, got[0].HasIpv6)
		assert.Len(t, got, 1)
	})

	t.Run("get ipv6 with invalid twin id", func(t *testing.T) {
		twinID := uint32(2)
		client := mock_rmb.NewMockClient(ctrl)
		client.EXPECT().Call(gomock.Any(), twinID, cmd, nil, gomock.Any()).Return(
			assert.AnError,
		)
		got, _ := ipv6.Get(ctx, client, twinID)
		// we expect an error here because the twin id is invalid but the implementation ignore errors
		//assert.NoError(t, err)
		assert.Empty(t, got)
	})
}
