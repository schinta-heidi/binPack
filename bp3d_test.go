package bp3d

import (
	"fmt"
	"reflect"
	"testing"
)

type result struct {
	packed Bin
}

type testData struct {
	bins        []*Bin
	items       []*Item
	expectation result
}

func TestPack(t *testing.T) {
	testCases := []testData{
		// Edge case that needs rotation.
		// from https://github.com/dvdoug/BoxPacker/issues/20
		{
			bins: []*Bin{
				NewBin("Le grande box", 100, 100, 300, 1500),
			},
			items: []*Item{
				NewItem("Item 1", 150, 50, 50, 20),
			},
			expectation: result{
				packed: Bin{
					"Le grande box", 100, 100, 300, 1500,
					[]*Item{
						{"Item 1", 150, 50, 50, 20, RotationType_HDW, Pivot{0, 0, 0}},
					},
				},
			},
		},

		// test three items fit into smaller bin.
		// from https://github.com/dvdoug/BoxPacker/blob/master/tests/PackerTest.php#L12
		{
			bins: []*Bin{
				NewBin("Le petite box", 296, 296, 8, 1000),
				NewBin("Le grande box", 2960, 2960, 80, 10000),
			},
			items: []*Item{
				NewItem("Item 1", 250, 250, 2, 200),
				NewItem("Item 2", 250, 250, 2, 200),
				NewItem("Item 3", 250, 250, 2, 200),
			},
			expectation: result{
				packed: Bin{
					"Le petite box", 296, 296, 8, 1000,
					[]*Item{
						{"Item 1", 250, 250, 2, 200, RotationType_WHD, Pivot{0, 0, 0}},
						{"Item 2", 250, 250, 2, 200, RotationType_WHD, Pivot{0, 0, 2}},
						{"Item 3", 250, 250, 2, 200, RotationType_WHD, Pivot{0, 0, 4}}},
				},
			},
		},

		// test three items fit into larger bin.
		// from https://github.com/dvdoug/BoxPacker/blob/master/tests/PackerTest.php#L36
		{
			bins: []*Bin{
				NewBin("Le petite box", 296, 296, 8, 1000),
				NewBin("Le grande box", 2960, 2960, 80, 10000),
			},
			items: []*Item{
				NewItem("Item 1", 2500, 2500, 20, 2000),
				NewItem("Item 2", 2500, 2500, 20, 2000),
				NewItem("Item 3", 2500, 2500, 20, 2000),
			},
			expectation: result{
				packed: Bin{
					"Le grande box", 2960, 2960, 80, 10000,
					[]*Item{
						{"Item 1", 2500, 2500, 20, 2000, RotationType_WHD, Pivot{0, 0, 0}},
						{"Item 2", 2500, 2500, 20, 2000, RotationType_WHD, Pivot{0, 0, 20}},
						{"Item 3", 2500, 2500, 20, 2000, RotationType_WHD, Pivot{0, 0, 40}}},
				},
			},
		},

		// TODO(gedex): five items packed into two large bins and one small bin.
		// from https://github.com/dvdoug/BoxPacker/blob/master/tests/PackerTest.php#L60

		// 1 bin that 7 items fit into.
		// from https://github.com/bom-d-van/binpacking/blob/master/binpacking_test.go
		{
			bins: []*Bin{
				NewBin("Bin 1", 220, 160, 100, 110),
			},
			items: []*Item{
				NewItem("Item 1", 20, 100, 30, 10),
				NewItem("Item 2", 100, 20, 30, 10),
				NewItem("Item 3", 20, 100, 30, 10),
				NewItem("Item 4", 100, 20, 30, 10),
				NewItem("Item 5", 100, 20, 30, 10),
				NewItem("Item 6", 100, 100, 30, 10),
				NewItem("Item 7", 100, 100, 30, 10),
			},
			expectation: result{
				packed: Bin{
					"Bin 1", 220, 160, 100, 110,
					[]*Item{
						{"Item 7", 100, 100, 30, 10, RotationType_WHD, Pivot{0, 0, 0}},
						{"Item 6", 100, 100, 30, 10, RotationType_WHD, Pivot{100, 0, 0}},
						{"Item 2", 100, 20, 30, 10, RotationType_HWD, Pivot{200, 0, 0}},
						{"Item 3", 20, 100, 30, 10, RotationType_HWD, Pivot{0, 100, 0}},
						{"Item 4", 100, 20, 30, 10, RotationType_WHD, Pivot{100, 100, 0}},
						{"Item 5", 100, 20, 30, 10, RotationType_HDW, Pivot{200, 100, 0}},
						{"Item 1", 20, 100, 30, 10, RotationType_HWD, Pivot{100, 120, 0}}},
				},
			},
		},
	}

	for _, tc := range testCases {
		testPack(t, tc)
	}
}

func testPack(t *testing.T, td testData) {
	packer := NewPacker()
	packer.AddBin(td.bins...)
	packer.AddItem(td.items...)

	out := packer.Pack()
	if out == nil {
		t.Fatalf("Got error: %v", out)
	}

	if !reflect.DeepEqual(out, &td.expectation.packed) {
		t.Errorf("\nGot:\n%+v\nwant:\n%+v", formatBins(out), formatBins(&td.expectation.packed))
	}
}

func formatBins(bin *Bin) string {
	var s string
	s += fmt.Sprintln(bin)
	s += fmt.Sprintln(" packed items:")
	for _, i := range bin.Items {
		s += fmt.Sprintln("  ", i)
	}
	return s
}
