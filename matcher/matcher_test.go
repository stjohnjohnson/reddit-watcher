package matcher

import (
	"reflect"
	"testing"
)

func TestGetSale(t *testing.T) {
	totalTests := []struct {
		in   string
		mode string
		out  string
	}{
		{
			"[US-TX] [H] PrimeCap / CM PBT L Cherry MX Blues [W] Paypal, Local Cash",
			"selling",
			"PrimeCap / CM PBT L Cherry MX Blues",
		},
		{
			"[CA-ON] [H] BKE Redux Heavy, FC660C 45g Topre Domes, Leopold Keycaps Doubleshot PBT Dolch [W] PayPal ",
			"international",
			"Not in the US: CA",
		},
		{
			"[NO][H]GMK Nautilus, Doomcaps, ETF, Brocaps [W]PayPal, trades",
			"international",
			"Not in the US: NO",
		},
		{
			"[US-PA][H] PayPal, Local Cash [W] 65g r7+ zealios, zeal stabs r2 or newer",
			"buying",
			"Not for sale: 65g r7+ zealios, zeal stabs r2 or newer",
		},
		{
			"[US-FL] [H] PayPal [W] ~60 alps orange or salmon",
			"buying",
			"Not for sale: ~60 alps orange or salmon",
		},
		{
			"[Vendor] GMK Olivia GB MOQ Hit! - Last day! | NovelKeys Restocks | BOX Royals Update",
			"vendor",
			"Not a trade: Vendor, GMK Olivia GB MOQ Hit! - Last day! | NovelKeys Restocks | BOX Royals Update",
		},
		{
			"[Giveaway] KF Orochi - Painted by KeypressGraphics",
			"giveaway",
			"Not a trade: Giveaway, KF Orochi - Painted by KeypressGraphics",
		},
		{
			"May Confirmed Trade Thread",
			"bad",
			"Not parsable: May Confirmed Trade Thread",
		},
		{
			"[MY][H] Fuguthulus, matrix abels, artisans [W] Trades, Paypal ",
			"international",
			"Not in the US: MY",
		},
	}

	for _, tt := range totalTests {
		out, mode, err := GetSale(tt.in)

		if err != nil {
			out = err.Error()
		}
		if out != tt.out {
			t.Errorf("Expected Out %q, got %q", tt.out, out)
		}
		if mode != tt.mode {
			t.Errorf("Expected Mode %q, got %q", tt.mode, mode)
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
	}

	for _, tt := range totalTests {
		out := FindMatching(tt.in, sale, description)

		if !reflect.DeepEqual(out, tt.out) {
			t.Errorf("Expected %q, got %q", tt.out, out)
		}
	}
}
