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
				count: "1",
				name:  "apple",
				price: "5",
			},
			want: "1 apple for sale for $5",
		},
		{
			r: record{
				count: "0",
				name:  "apple",
				price: "5",
			},
			want: "0 apples for sale for $5",
		},
		{
			r: record{
				count: "10",
				name:  "apple",
				price: "5",
			},
			want: "10 apples for sale for $5",
		},
		{
			r: record{
				count: "1",
				name:  "apple",
				price: "-1",
			},
			want: "1 apple",
		},
		{
			r: record{
				count: "10",
				name:  "apple",
				price: "-1",
			},
			want: "10 apples",
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
					count: "10",
					name:  "Health Potion",
					price: "10",
				},
				{
					count: "10000",
					name:  "Mana Potion",
					price: "8",
				},
				{
					count: "1",
					name:  "Death Potion",
					price: "5000",
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
					count: "1",
					name:  "Divine Bow",
				},
				{
					count: "19",
					name:  "Regular arrows",
				},
				{
					count: "1",
					name:  "Shield",
				},
			},
			want: "`" + `` + "`" + `` + "`" + `
╔══════════════════════════╗
║ Quantity  Item           ║
║──────────────────────────║
║ 1         Divine Bow     ║
║ 19        Regular arrows ║
║ 1         Shield         ║
╚══════════════════════════╝
` + "`" + `` + "`" + `` + "`",
		},
		{
			rs: records{
				{
					count: "10",
					name:  "Health Potion",
				},
				{
					count: "10000",
					name:  "Mana Potion",
					price: "8",
				},
				{
					count: "1",
					name:  "Death Potion",
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
					count: "0",
					name:  "Health Potion",
				},
				{
					count: "10000",
					name:  "Mana Potion",
					price: "8",
				},
				{
					count: "1",
					name:  "Death Potion",
					price: "-1",
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
					count: "1",
					name:  "apple",
					price: "1",
				},
			},
			want: "`" + `` + "`" + `` + "`" + `
╔════════════════════════╗
║ Quantity  Item   Price ║
║────────────────────────║
║ 1         apple  $1    ║
╚════════════════════════╝
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
