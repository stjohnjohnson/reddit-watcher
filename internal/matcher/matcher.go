package matcher

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// Buying is for a post that is buying something
	Buying = "buying"
	// Selling is for a post that is selling something
	Selling = "selling"
	// Vendor is for a post from a specific vendor
	Vendor = "vendor"
	// Artisan is for a post from a specific artisan
	Artisan = "artisan"
	// GroupBuy is for a post about a group buy
	GroupBuy = "groupbuy"
	// InterestCheck is for a post about gauging interest
	InterestCheck = "interestcheck"
	// Giveaway is for a post about giving away something
	Giveaway = "giveaway"
	// Types is a list of all types
	Types = []string{Buying, Selling, Vendor, Artisan, GroupBuy, InterestCheck, Giveaway}
)

var typeLookup = map[string]string{
	"vendor":   Vendor,
	"artisan":  Artisan,
	"gb":       GroupBuy,
	"ic":       InterestCheck,
	"giveaway": Giveaway,
}

// ParsedPost represents a parsed Reddit post
type ParsedPost struct {
	Type     string
	Contents string
	Region   string
}

// [TYPE] Something
var nonSalesRex = regexp.MustCompile(`(?i)^\[(vendor|artisan|gb|ic|giveaway)\]\s*(.*)$`)

// [COUNTRY-STATE] [H] Something [W] Something else
var salesRex = regexp.MustCompile(`(?i)^\[(\w+)(?:-\w+)?\]\s*\[H\]\s*(.*)\s*\[W\]\s*(.*)$`)

var moneyRex = regexp.MustCompile(`(?i)(paypal|cash)`)

// ParseTitle returns the a type, content, and error
func ParseTitle(title string) (*ParsedPost, error) {
	if m := nonSalesRex.FindStringSubmatch(title); m != nil {
		return &ParsedPost{
			Type:     typeLookup[strings.ToLower(m[1])],
			Contents: strings.TrimSpace(m[2]),
		}, nil
	}

	sales := salesRex.FindStringSubmatch(title)
	if sales == nil {
		return nil, fmt.Errorf("not parsable: %s", title)
	}
	region, have, want := sales[1], sales[2], sales[3]

	// Ensure it's for sale
	if !moneyRex.MatchString(want) {
		return &ParsedPost{
			Type:     Buying,
			Contents: strings.TrimSpace(want),
			Region:   region,
		}, nil
	}

	// Return the things for sale
	return &ParsedPost{
		Type:     Selling,
		Contents: strings.TrimSpace(have),
		Region:   region,
	}, nil
}

// FindMatching returns list of keywords that match a given title/description
func FindMatching(keywords []string, title, desc string) []string {
	matches := []string{}
	title = strings.ToLower(title)
	desc = strings.ToLower(desc)

	for _, keyword := range keywords {
		if keyword == "*" || strings.Contains(title, keyword) || strings.Contains(desc, keyword) {
			matches = append(matches, keyword)
		}
	}

	return matches
}
