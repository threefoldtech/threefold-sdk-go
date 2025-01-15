package indexer

import (
	"context"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_rmb "github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/mocks"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/pkg/types"
	zosDmiTypes "github.com/threefoldtech/zos/pkg/capacity/dmi"
)

func TestNewDMIWork(t *testing.T) {
	wanted := &DMIWork{
		findersInterval: map[string]time.Duration{
			"up":  2 * time.Minute,
			"new": 5 * time.Minute,
		},
	}
	dmi := NewDMIWork(2)
	assert.Exactlyf(t, wanted, dmi, "got: %v , expected: %v", dmi, wanted)
}

func TestDMIGet(t *testing.T) {
	dmi := NewDMIWork(2)
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	t.Run("get dmi with valid twin id", func(t *testing.T) {
		twinID := uint32(1)
		expected := zosDmiTypes.DMI{
			Sections: []zosDmiTypes.Section{
				{
					TypeStr: "BIOS",
					SubSections: []zosDmiTypes.SubSection{
						{
							Title: "BIOS Information",
							Properties: map[string]zosDmiTypes.PropertyData{
								"Vendor":  {Val: "American Megatrends Inc."},
								"Version": {Val: "3.2"},
							},
						},
					},
				},
				{
					TypeStr: "Baseboard",
					SubSections: []zosDmiTypes.SubSection{
						{
							Title: "Base Board Information",
							Properties: map[string]zosDmiTypes.PropertyData{
								"Manufacturer": {Val: "Supermicro"},
								"Product Name": {Val: "X9DRi-LN4+/X9DR3-LN4+"},
							},
						},
					},
				},
				{
					TypeStr: "Processor",
					SubSections: []zosDmiTypes.SubSection{
						{
							Title: "Processor Information",
							Properties: map[string]zosDmiTypes.PropertyData{
								"Version":      {Val: "Intel(R) Xeon(R) CPU E5-2620 0 @ 2.00GHz"},
								"Thread Count": {Val: "12"},
							},
						},
						{
							Title: "Processor Information",
							Properties: map[string]zosDmiTypes.PropertyData{
								"Version":      {Val: "Intel(R) Xeon(R) CPU E5-2620 0 @ 2.00GHz"},
								"Thread Count": {Val: "12"},
							},
						},
					},
				},
				{
					TypeStr: "MemoryDevice",
					SubSections: []zosDmiTypes.SubSection{
						{
							Title: "Memory Device",
							Properties: map[string]zosDmiTypes.PropertyData{
								"Manufacturer": {Val: "Samsung"},
								"Type":         {Val: "DDR3"},
							},
						},
						{
							Title: "Memory Device",
							Properties: map[string]zosDmiTypes.PropertyData{
								"Manufacturer": {Val: "Samsung"},
								"Type":         {Val: "DDR3"},
							},
						},
						{
							Title: "Memory Device",
							Properties: map[string]zosDmiTypes.PropertyData{
								"Manufacturer": {Val: "Samsung"},
								"Type":         {Val: "DDR3"},
							},
						},
					},
				},
			},
		}
		client := mock_rmb.NewMockClient(ctrl)
		client.EXPECT().Call(gomock.Any(), twinID, DmiCallCmd, nil, &zosDmiTypes.DMI{}).DoAndReturn(
			func(ctx context.Context, twin uint32, fn string, data, result interface{}) error {
				*(result.(*zosDmiTypes.DMI)) = expected
				return nil
			},
		)
		got, err := dmi.Get(ctx, client, twinID)
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, expected.Sections[0].SubSections[0].Properties["Vendor"].Val, got[0].BIOS.Vendor)
		assert.Equal(t, expected.Sections[0].SubSections[0].Properties["Version"].Val, got[0].BIOS.Version)
		assert.Equal(t, expected.Sections[1].SubSections[0].Properties["Manufacturer"].Val, got[0].Baseboard.Manufacturer)
		assert.Equal(t, expected.Sections[1].SubSections[0].Properties["Product Name"].Val, got[0].Baseboard.ProductName)
		assert.Equal(t, expected.Sections[2].SubSections[0].Properties["Version"].Val, got[0].Processor[0].Version)
		assert.Equal(t, expected.Sections[2].SubSections[0].Properties["Thread Count"].Val, got[0].Processor[0].ThreadCount)
		assert.Equal(t, expected.Sections[3].SubSections[0].Properties["Manufacturer"].Val, got[0].Memory[0].Manufacturer)
		assert.Equal(t, expected.Sections[3].SubSections[0].Properties["Type"].Val, got[0].Memory[0].Type)
		assert.IsTypef(t, got, []types.Dmi{}, "got: %T , expected: %T", got, []types.Dmi{})

	})

	t.Run("get dmi with invalid twin id", func(t *testing.T) {
		twinID := uint32(2)
		client := mock_rmb.NewMockClient(ctrl)
		client.EXPECT().Call(gomock.Any(), twinID, DmiCallCmd, nil, &zosDmiTypes.DMI{}).Return(
			assert.AnError,
		)
		_, err := dmi.Get(ctx, client, twinID)
		assert.Error(t, err)
	})

}
