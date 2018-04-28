package matcher

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var titleRex = regexp.MustCompile(`(?i)^\[(\w+)-(\w+)\]\s*\[H\]\s*(.*)\s*\[W\]\s*(.*)$`)
var moneyRex = regexp.MustCompile(`(?i)(paypal|cash)`)

func FindValidSale(title string) (string, error) {
	fields := titleRex.FindStringSubmatch(title)

	if fields == nil {
		return "", errors.New(fmt.Sprintf("Not parsable: %s", title))
	}
	cntry, state, have, want := fields[1], fields[2], fields[3], fields[4]

	// Check location
	if cntry != "US" {
		return "", errors.New(fmt.Sprintf("Not in the US: %s, %s", cntry, state))
	}

	// Ensure it's for sale
	if !moneyRex.MatchString(want) {
		return "", errors.New(fmt.Sprintf("Not for sale: %s", want))
	}

	// Return the things for sale
	return strings.TrimSpace(have), nil
}
