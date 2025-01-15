package indexer

import (
	"context"
	"errors"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_rmb "github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/mocks"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/pkg/types"
)

func TextNewWorkloadWork(t *testing.T) {
	wanted := &WorkloadWork{
		findersInterval: map[string]time.Duration{
			"up": time.Duration(2) * time.Minute,
		},
	}
	workload := NewWorkloadWork(2)
	assert.Exactly(t, wanted, workload)
}

func TestWorkloadGet(t *testing.T) {
	workload := NewWorkloadWork(2)
	ctrl := gomock.NewController(t)

	t.Run("get workload with valid twin id", func(t *testing.T) {
		twinID := uint32(1)
		var response struct {
			Users struct {
				Workloads uint32 `json:"workloads"`
			} `json:"users"`
		}
		ctx := context.Background()

		client := mock_rmb.NewMockClient(ctrl)
		client.EXPECT().Call(gomock.Any(), twinID, statsCall, nil, gomock.Any()).DoAndReturn(
			func(ctx context.Context, twinId uint32, fn string, data, result interface{}) error {
				response.Users.Workloads = 10
				*result.(*struct {
					Users struct {
						Workloads uint32 `json:"workloads"`
					} `json:"users"`
				}) = response
				return nil
			},
		)
		workloads, err := workload.Get(ctx, client, twinID)
		assert.NoError(t, err)
		assert.Len(t, workloads, 1)
		assert.Equal(t, response.Users.Workloads, workloads[0].WorkloadsNumber)
		assert.Equal(t, twinID, workloads[0].NodeTwinId)
		assert.IsType(t, workloads, []types.NodesWorkloads{})
	})

	t.Run("get workload with invalid twin id", func(t *testing.T) {
		twinID := uint32(1)
		ctx := context.Background()

		client := mock_rmb.NewMockClient(ctrl)
		client.EXPECT().Call(gomock.Any(), twinID, statsCall, nil, gomock.Any()).Return(errors.New("error"))
		workloads, err := workload.Get(ctx, client, twinID)
		assert.Error(t, err)
		assert.Len(t, workloads, 0)
	})
}
