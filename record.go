package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/gertd/go-pluralize"
)

type declinedError struct {
	r record
}

func (e *declinedError) Error() string {
	return fmt.Sprintf("transaction \"%v\" declined", e.r)
}

// record represents a single entry in an inventory.
type record struct {
	count string
	name  string
	price string
}

// addCount adds to a record's count.
// If the count would be made negative, a declinedError is returned.
func (r *record) addCount(count string) error {
	current, err := strconv.Atoi(r.count)
	if err != nil {
		return fmt.Errorf(
			"failed parsing count of %v: %v\n",
			r,
			err,
		)
	}
	addend, err := strconv.Atoi(count)
	if err != nil {
		return fmt.Errorf(
			"failed parsing count %v: %v\n",
			count,
			err,
		)
	}
	sum := current + addend
	if sum < 0 {
		return &declinedError{*r}
	}
	r.count = strconv.Itoa(sum)
	return nil
}

// String prints out a pretty message describing the record.
func (r record) String() string {
	var buf bytes.Buffer

	// Ignoring error to use 0 as fallback count.
	count, _ := strconv.Atoi(r.count)
	plur := pluralize.NewClient()
	buf.WriteString(strings.TrimSpace(plur.Pluralize(r.name, count, true)))

	price, err := strconv.Atoi(r.price)
	if r.price != NULL_PRICE && err == nil {
		buf.WriteString(" for sale for $")
		buf.WriteString(humanize.Comma(int64(price)))
	}

	return buf.String()
}

// records represents a collection of entries in an inventory.
type records []record

var recordsTable = lipgloss.NewStyle().
	BorderStyle(lipgloss.DoubleBorder())
var recordsColumn = lipgloss.NewStyle().
	PaddingLeft(1).
	PaddingRight(1)

// String prints out a pretty table showing the records.
// It looks like this, but wrapped in 3 backticks for discord:
//
//	╔══════════════════════════════════╗
//	║ Quantity  Item            Price  ║
//	║──────────────────────────────────║
//	║ 10        Health Potions  $10    ║
//	║ 10,000    Mana Potions    $8     ║
//	║ 1         Death Potion    $5,000 ║
//	╚══════════════════════════════════╝
//
// The price column is omitted if no items contain a price.
func (rs records) String() string {
	var buf bytes.Buffer
	buf.WriteString("```\n")

	// Gather column data.
	var counts []string
	var names []string
	var prices []string
	for _, r := range rs {
		count, _ := strconv.Atoi(r.count)
		if count == 0 {
			// Skip records with 0 count.
			continue
		}
		r.count = humanize.Comma(int64(count))
		counts = append(counts, r.count)

		plur := pluralize.NewClient()
		names = append(
			names,
			strings.TrimSpace(plur.Pluralize(r.name, count, false)),
		)

		price, err := strconv.Atoi(r.price)
		if r.price != NULL_PRICE && err == nil {
			r.price = "$" + humanize.Comma(int64(price))
		} else {
			r.price = ""
		}
		prices = append(prices, r.price)
	}

	// Add headings.
	counts = append([]string{"Quantity"}, counts...)
	names = append([]string{"Item"}, names...)

	// Exclude the price column if it is empty (other than the header).
	var found bool
	for _, p := range prices {
		if p != "" {
			found = true
		}
	}
	if found {
		prices = append([]string{"Price"}, prices...)
	} else {
		prices = []string{}
	}

	countCol := recordsColumn.Render(
		lipgloss.JoinVertical(lipgloss.Top, counts...),
	)
	nameCol := recordsColumn.Render(
		lipgloss.JoinVertical(lipgloss.Top, names...),
	)
	priceCol := recordsColumn.Render(
		lipgloss.JoinVertical(lipgloss.Top, prices...),
	)

	if lipgloss.Height(priceCol) <= 1 {
		priceCol = ""
	}

	table := lipgloss.JoinHorizontal(lipgloss.Left, countCol, nameCol, priceCol)
	// Add a line under the header.
	line := strings.Repeat("─", lipgloss.Width(table))
	rows := strings.Split(table, "\n")
	header := rows[0]
	body := rows[1:]
	rows = append([]string{header, line}, body...)
	table = strings.Join(rows, "\n")

	buf.WriteString(recordsTable.Render(table))
	buf.WriteString("\n```")
	return buf.String()
}

// forSale returns records which have a price.
func (rs records) forSale() records {
	var recs records
	for _, r := range rs {
		if r.price == NULL_PRICE {
			continue
		}
		recs = append(recs, r)
	}
	return recs
}
