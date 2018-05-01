package matcher

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// [TYPE] Something
var nonSalesRex = regexp.MustCompile(`(?i)^\[(vendor|artisan|gb)\]\s*(.*)$`)

// [COUNTRY-STATE] [H] Something [W] Something else
var salesRex = regexp.MustCompile(`(?i)^\[(\w+)(?:-\w+)?\]\s*\[H\]\s*(.*)\s*\[W\]\s*(.*)$`)

var moneyRex = regexp.MustCompile(`(?i)(paypal|cash)`)

func GetSale(title string) (string, string, error) {
	nonSales := nonSalesRex.FindStringSubmatch(title)
	if nonSales != nil {
		return "", strings.ToLower(nonSales[1]), errors.New(fmt.Sprintf("Not a trade: %s, %s", nonSales[1], nonSales[2]))
	}

	sales := salesRex.FindStringSubmatch(title)
	if sales == nil {
		return "", "bad", errors.New(fmt.Sprintf("Not parsable: %s", title))
	}
	cntry, have, want := sales[1], sales[2], sales[3]

	// Check location
	if cntry != "US" {
		return "", "international", errors.New(fmt.Sprintf("Not in the US: %s", cntry))
	}

	// Ensure it's for sale
	if !moneyRex.MatchString(want) {
		return "", "buying", errors.New(fmt.Sprintf("Not for sale: %s", want))
	}

	// Return the things for sale
	return strings.TrimSpace(have), "selling", nil
}

func FindMatching(keywords []string, forsale, desc string) []string {
	matches := []string{}
	forsale = strings.ToLower(forsale)
	desc = strings.ToLower(desc)

	for _, keyword := range keywords {
		if strings.Contains(forsale, keyword) || strings.Contains(desc, keyword) {
			matches = append(matches, keyword)
		}
	}

	return matches
}
