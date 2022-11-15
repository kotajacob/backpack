package main

import "testing"

func TestRecordString(t *testing.T) {
	type test struct {
		r    record
		want string
	}

	tests := []test{
		{
			r: record{
				count: 1,
				name:  "apple",
				price: 5,
			},
			want: "1 Apple for sale for $5",
		},
		{
			r: record{
				count: 0,
				name:  "apple",
				price: 5,
			},
			want: "0 Apples for sale for $5",
		},
		{
			r: record{
				count: 10,
				name:  "apple",
				price: 5,
			},
			want: "10 Apples for sale for $5",
		},
		{
			r: record{
				count: 1,
				name:  "apple",
				price: -1,
			},
			want: "1 Apple",
		},
		{
			r: record{
				count: 10,
				name:  "apple",
				price: -1,
			},
			want: "10 Apples",
		},
	}
	for _, tc := range tests {
		got := tc.r.String()
		if got != tc.want {
			t.Fatalf("want: %v got: %v\n", tc.want, got)
		}
	}
}

func TestRecordsString(t *testing.T) {
	type test struct {
		rs   records
		want string
	}

	tests := []test{
		{
			rs: records{
				{
					count: 10,
					name:  "Health Potion",
					price: 10,
				},
				{
					count: 10000,
					name:  "Mana Potion",
					price: 8,
				},
				{
					count: 1,
					name:  "Death Potion",
					price: 5000,
				},
			},
			want: "`" + `` + "`" + `` + "`" + `
╔══════════════════════════════════╗
║ Quantity  Item            Price  ║
║──────────────────────────────────║
║ 10        Health Potions  $10    ║
║ 10,000    Mana Potions    $8     ║
║ 1         Death Potion    $5,000 ║
╚══════════════════════════════════╝
` + "`" + `` + "`" + `` + "`",
		},
		{
			rs: records{
				{
					count: 1,
					name:  "Divine Bow",
					price: 0,
				},
				{
					count: 19,
					name:  "Regular arrows",
					price: 0,
				},
				{
					count: 1,
					name:  "Shield",
					price: 0,
				},
			},
			want: "`" + `` + "`" + `` + "`" + `
╔═════════════════════════════════╗
║ Quantity  Item            Price ║
║─────────────────────────────────║
║ 1         Divine Bow      $0    ║
║ 19        Regular arrows  $0    ║
║ 1         Shield          $0    ║
╚═════════════════════════════════╝
` + "`" + `` + "`" + `` + "`",
		},
		{
			rs: records{
				{
					count: 10,
					name:  "Health Potion",
					price: -1,
				},
				{
					count: 10000,
					name:  "Mana Potion",
					price: 8,
				},
				{
					count: 1,
					name:  "Death Potion",
					price: -1,
				},
			},
			want: "`" + `` + "`" + `` + "`" + `
╔═════════════════════════════════╗
║ Quantity  Item            Price ║
║─────────────────────────────────║
║ 10        Health Potions        ║
║ 10,000    Mana Potions    $8    ║
║ 1         Death Potion          ║
╚═════════════════════════════════╝
` + "`" + `` + "`" + `` + "`",
		},
		{
			rs: records{
				{
					count: 0,
					name:  "Health Potion",
				},
				{
					count: 10000,
					name:  "Mana Potion",
					price: 8,
				},
				{
					count: 1,
					name:  "Death Potion",
					price: -1,
				},
			},
			want: "`" + `` + "`" + `` + "`" + `
╔═══════════════════════════════╗
║ Quantity  Item          Price ║
║───────────────────────────────║
║ 10,000    Mana Potions  $8    ║
║ 1         Death Potion        ║
╚═══════════════════════════════╝
` + "`" + `` + "`" + `` + "`",
		},
		{
			rs: records{
				{
					count: 1,
					name:  "apple",
					price: 1,
				},
			},
			want: "`" + `` + "`" + `` + "`" + `
╔════════════════════════╗
║ Quantity  Item   Price ║
║────────────────────────║
║ 1         Apple  $1    ║
╚════════════════════════╝
` + "`" + `` + "`" + `` + "`",
		},
		{
			rs: records{
				{
					count: 1,
					name:  "apple",
					price: -1,
				},
			},
			want: "`" + `` + "`" + `` + "`" + `
╔═════════════════╗
║ Quantity  Item  ║
║─────────────────║
║ 1         Apple ║
╚═════════════════╝
` + "`" + `` + "`" + `` + "`",
		},
	}

	for _, tc := range tests {
		got := tc.rs.String()
		if got != tc.want {
			t.Fatalf("\nwant:\n%v\ngot:\n%v\n", tc.want, got)
		}
	}
}
