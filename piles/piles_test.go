package piles

import (
	"image"
	"reflect"
	"testing"

	"dr2w.com/progenart/img"
)

var resolveTests = []struct {
	name   string
	config *Config
	input  []string
	output []string
}{
	{
		"Basic 4",
		&Config{
			Wrap:         false,
			Connectivity: Four,
		},
		[]string{
			"00000",
			"00100",
			"00820",
			"00000",
			"00000",
		},
		[]string{
			"00000",
			"00310",
			"02101",
			"00210",
			"00000",
		},
	},
}

func TestResolve(t *testing.T) {
	for _, test := range resolveTests {
		input := img.NewFromStrings(test.input)
		want := img.NewFromStrings(test.output)
		if got := test.config.Resolve(input); !reflect.DeepEqual(got, want) {
			t.Errorf("%s:\ngot:\n%s\nwant:\n%s\n", test.name, img.ToSimpleString(got), img.ToSimpleString(want))
		}

	}
}

var spillTests = []struct {
	name   string
	input  []string
	point  image.Point
	config *Config
	want   []string
}{
	{
		"Basic 4",
		[]string{
			"000",
			"040",
			"000",
		},
		image.Pt(1, 1),
		&Config{
			Wrap:         false,
			Connectivity: Four,
		},
		[]string{
			"010",
			"101",
			"010",
		},
	},
	{
		"Basic 8",
		[]string{
			"000",
			"080",
			"000",
		},
		image.Pt(1, 1),
		&Config{
			Wrap:         false,
			Connectivity: Eight,
		},
		[]string{
			"111",
			"101",
			"111",
		},
	},
	{
		"Extra 4",
		[]string{
			"000",
			"060",
			"000",
		},
		image.Pt(1, 1),
		&Config{
			Wrap:         false,
			Connectivity: Four,
		},
		[]string{
			"010",
			"121",
			"010",
		},
	},
	{
		"Overflow 4, No Wrap",
		[]string{
			"000",
			"400",
			"000",
		},
		image.Pt(0, 1),
		&Config{
			Wrap:         false,
			Connectivity: Four,
		},
		[]string{
			"100",
			"010",
			"100",
		},
	},
	{
		"Overflow 4, Horizontal Wrap",
		[]string{
			"000",
			"400",
			"000",
		},
		image.Pt(0, 1),
		&Config{
			Wrap:         true,
			Connectivity: Four,
		},
		[]string{
			"100",
			"011",
			"100",
		},
	},
	{
		"Overflow 4, Vertical Wrap - Up",
		[]string{
			"040",
			"000",
			"000",
		},
		image.Pt(1, 0),
		&Config{
			Wrap:         true,
			Connectivity: Four,
		},
		[]string{
			"101",
			"010",
			"010",
		},
	},
	{
		"Overflow 4, Vertical Wrap - Down",
		[]string{
			"000",
			"000",
			"040",
		},
		image.Pt(1, 2),
		&Config{
			Wrap:         true,
			Connectivity: Four,
		},
		[]string{
			"010",
			"010",
			"101",
		},
	},
	{
		"Complex Case 1",
		[]string{
			"400",
			"354",
			"901",
		},
		image.Pt(0, 2),
		&Config{
			Wrap:         true,
			Connectivity: Eight,
		},
		[]string{
			"511",
			"465",
			"112",
		},
	},
}

func TestSpill(t *testing.T) {
	for _, test := range spillTests {
		input := img.NewFromStrings(test.input)
		want := img.NewFromStrings(test.want)
		if got, _ := test.config.spill(test.point, input); !reflect.DeepEqual(got, want) {
			t.Errorf("%s:\ngot:\n%s,\nwant:\n%s", test.name, img.ToSimpleString(got), img.ToSimpleString(want))
		}
	}
}
