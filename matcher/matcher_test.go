package matcher

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseTitle(t *testing.T) {
	totalTests := []struct {
		in  string
		out *ParsedPost
		err error
	}{
		{
			"[US-TX] [H] PrimeCap / CM PBT L Cherry MX Blues [W] Paypal, Local Cash",
			&ParsedPost{
				Type:     Selling,
				Contents: "PrimeCap / CM PBT L Cherry MX Blues",
				Region:   "US",
			},
			nil,
		},
		{
			"[CA-ON] [H] BKE Redux Heavy, FC660C 45g Topre Domes, Leopold Keycaps Doubleshot PBT Dolch [W] PayPal ",
			&ParsedPost{
				Type:     Selling,
				Contents: "BKE Redux Heavy, FC660C 45g Topre Domes, Leopold Keycaps Doubleshot PBT Dolch",
				Region:   "CA",
			},
			nil,
		},
		{
			"[NO][H]GMK Nautilus, Doomcaps, ETF, Brocaps [W]PayPal, trades",
			&ParsedPost{
				Type:     Selling,
				Contents: "GMK Nautilus, Doomcaps, ETF, Brocaps",
				Region:   "NO",
			},
			nil,
		},
		{
			"[US-PA][H] PayPal, Local Cash [W] 65g r7+ zealios, zeal stabs r2 or newer",
			&ParsedPost{
				Type:     Buying,
				Contents: "65g r7+ zealios, zeal stabs r2 or newer",
				Region:   "US",
			},
			nil,
		},
		{
			"[US-FL] [H] PayPal [W] ~60 alps orange or salmon",
			&ParsedPost{
				Type:     Buying,
				Contents: "~60 alps orange or salmon",
				Region:   "US",
			},
			nil,
		},
		{
			"[Vendor] GMK Olivia GB MOQ Hit! - Last day! | NovelKeys Restocks | BOX Royals Update",
			&ParsedPost{
				Type:     Vendor,
				Contents: "GMK Olivia GB MOQ Hit! - Last day! | NovelKeys Restocks | BOX Royals Update",
			},
			nil,
		},
		{
			"[Giveaway] KF Orochi - Painted by KeypressGraphics",
			&ParsedPost{
				Type:     Giveaway,
				Contents: "KF Orochi - Painted by KeypressGraphics",
			},
			nil,
		},
		{
			"May Confirmed Trade Thread",
			nil,
			fmt.Errorf("not parsable: May Confirmed Trade Thread"),
		},
	}

	for _, tt := range totalTests {
		out, err := ParseTitle(tt.in)

		if !reflect.DeepEqual(out, tt.out) {
			t.Errorf("Expected Out %q, got %q", tt.out, out)
		}
		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("Expected Err %q, got %q", tt.err, err)
		}
	}
}

func TestFindMatching(t *testing.T) {
	sale := "PrimeCap / CM PBT L Cherry MX Blues"
	description := "something something timestamp something tao hao"

	totalTests := []struct {
		in  []string
		out []string
	}{
		{
			[]string{"brocap"},
			[]string{},
		},
		{
			[]string{"primecap", "brocap"},
			[]string{"primecap"},
		},
		{
			[]string{"tao hao", "primecap"},
			[]string{"tao hao", "primecap"},
		},
		{
			[]string{"*"},
			[]string{"*"},
		},
	}

	for _, tt := range totalTests {
		out := FindMatching(tt.in, sale, description)

		if !reflect.DeepEqual(out, tt.out) {
			t.Errorf("Expected %q, got %q", tt.out, out)
		}
	}
}
