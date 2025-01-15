package indexer

import (
	"context"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/internal/explorer/db"
	mocks "github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/mocks"
)

func TestNewNodeFinder(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockDB := mocks.NewMockDatabase(ctrl)

	idsChan := make(chan uint32, 10)

	mockDB.EXPECT().GetLastNodeTwinID(gomock.Any()).Return(uint32(1), nil).Times(1)
	mockDB.EXPECT().GetNodeTwinIDsAfter(gomock.Any(), gomock.Any()).Return([]uint32{2, 3}, nil).Times(1)

	ctx, cancel := context.WithCancel(context.Background())

	go newNodesFinder(ctx, time.Second, mockDB, idsChan)

	time.Sleep(1 * time.Second)

	cancel()

	var ids []uint32
	for {
		select {
		case id := <-idsChan:
			ids = append(ids, id)
		default:
			goto done
		}
	}

done:

	assert.Equal(t, []uint32{2, 3}, ids)

}

func TestHealthyNodesFinder(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDB := mocks.NewMockDatabase(ctrl)

	idsChan := make(chan uint32, 10)

	mockDB.EXPECT().GetHealthyNodeTwinIds(gomock.Any()).Return([]uint32{1, 2, 3}, nil).Times(1)

	ctx, cancel := context.WithCancel(context.Background())

	go healthyNodesFinder(ctx, time.Second, mockDB, idsChan)

	time.Sleep(500 * time.Millisecond)

	cancel()

	var ids []uint32
	for {
		select {
		case id := <-idsChan:
			ids = append(ids, id)
		default:
			goto done
		}
	}

done:

	assert.Equal(t, []uint32{1, 2, 3}, ids)
	assert.Len(t, ids, 3)

}

func TestUpNodesFinder(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDB := mocks.NewMockDatabase(ctrl)

	idsChan := make(chan uint32, 10)

	mockDB.EXPECT().GetNodes(gomock.Any(), gomock.Any(), gomock.Any()).Return([]db.Node{
		{
			TwinID:   1,
			NodeID:   1,
			FarmID:   3,
			FarmName: "farm",
		},
		{
			TwinID:   2,
			NodeID:   2,
			FarmID:   3,
			FarmName: "farm",
		},
	}, uint(2), nil)

	ctx, cancel := context.WithCancel(context.Background())

	go upNodesFinder(ctx, time.Second, mockDB, idsChan)

	time.Sleep(500 * time.Millisecond)

	cancel()

	var ids []uint32
	for {
		select {
		case id := <-idsChan:
			ids = append(ids, id)
		default:
			goto done
		}
	}

done:

	assert.Equal(t, []uint32{1, 2}, ids)
	assert.Len(t, ids, 2)

}
