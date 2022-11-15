package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

// Coin is the canonical name for money.
const Coin = "coin"

// NotForSale indicates that the item is not meant to be sold.
const NotForSale = -1

// Unchanged indicates that the item's price should not be changed.
const Unchanged = -2

type declinedError struct {
	r record
}

func (e *declinedError) Error() string {
	return fmt.Sprintf("transaction \"%v\" declined", e.r)
}

// record represents a single entry in an inventory.
type record struct {
	count int
	name  string
	price int
}

// addCount adds to a record's count.
// If the count would be made negative, a declinedError is returned.
func (r *record) addCount(count int) error {
	sum := r.count + count
	if sum < 0 {
		return &declinedError{*r}
	}
	r.count = sum
	return nil
}

// String prints out a pretty message describing the record.
func (r record) String() string {
	var buf bytes.Buffer

	// Ignoring error to use 0 as fallback count.
	buf.WriteString(strconv.Itoa(r.count))
	buf.WriteString(" ")
	buf.WriteString(displayName(r.name, r.count))

	if r.price != NotForSale && r.price != Unchanged {
		buf.WriteString(" for sale for $")
		buf.WriteString(humanize.Comma(int64(r.price)))
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
		if r.count == 0 {
			// Skip records with 0 count.
			continue
		}
		counts = append(counts, humanize.Comma(int64(r.count)))

		names = append(
			names,
			strings.TrimSpace(displayName(r.name, r.count)),
		)

		if r.price != NotForSale {
			prices = append(prices, "$"+humanize.Comma(int64(r.price)))
		} else {
			prices = append(prices, "")
		}
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
		if r.price == NotForSale {
			continue
		}
		recs = append(recs, r)
	}
	return recs
}
