// License: AGPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuyItem(t *testing.T) {
	type test struct {
		request string
		coins   string
		seller  string

		wantReply  string
		buyerWant  string
		sellerWant string
	}

	tests := []test{
		{
			request: "10 apples",
			coins:   "50",
			seller:  "20,apple,1",
			wantReply: "buyer bought 10 apples for $10\n" +
				"buyer has 10 apples\n" +
				"seller has 10 apples for sale for $1",
			buyerWant:  "40,coin,-1\n10,apple,-1",
			sellerWant: "10,apple,1",
		},
		{
			request: "1 apple",
			coins:   "50",
			seller:  "",
			wantReply: "seller does not have 1 apple in stock\n" +
				"Please choose one of the following items:\n" +
				"```\n" +
				"╔════════════════╗\n" +
				"║ Quantity  Item ║\n" +
				"║────────────────║\n" +
				"╚════════════════╝\n" +
				"```",
			buyerWant:  "50,coin,-1",
			sellerWant: "",
		},
		{
			request: "1 apple",
			coins:   "50",
			seller:  "1,apple,-1\n10,arrow,-1\n1,sword,100",
			wantReply: "seller does not have 1 apple for sale\n" +
				"Please choose one of the following items:\n" +
				"```\n" +
				"╔════════════════════════╗\n" +
				"║ Quantity  Item   Price ║\n" +
				"║────────────────────────║\n" +
				"║ 1         sword  $100  ║\n" +
				"╚════════════════════════╝\n" +
				"```",
			buyerWant:  "50,coin,-1",
			sellerWant: "1,apple,-1\n10,arrow,-1\n1,sword,100",
		},
		{
			request: "10 apples",
			coins:   "50",
			seller:  "1,apple,1",
			wantReply: "seller does not have 10 apples in stock\n" +
				"Please choose one of the following items:\n" +
				"```\n" +
				"╔════════════════════════╗\n" +
				"║ Quantity  Item   Price ║\n" +
				"║────────────────────────║\n" +
				"║ 1         apple  $1    ║\n" +
				"╚════════════════════════╝\n" +
				"```",
			buyerWant:  "50,coin,-1",
			sellerWant: "1,apple,1",
		},
		{
			request: "10 apples",
			coins:   "2",
			seller:  "100,apple,1",
			wantReply: "buyer has insufficient funds\n" +
				"10 apples costs $10\n" +
				"buyer only has 2 coins",
			buyerWant:  "2,coin,-1",
			sellerWant: "100,apple,1",
		},
		{
			request:    "0 apples",
			coins:      "50",
			seller:     "20,apple,1",
			wantReply:  "You can't buy 0 of an item, silly!",
			buyerWant:  "50,coin,-1",
			sellerWant: "20,apple,1",
		},
		{
			request: "-10 regular arrows",
			coins:   "50",
			seller:  "20,regular arrow,1",
			wantReply: "You've requested to give away your items?\n" +
				"Try again with: \"10 regular arrows\"",
			buyerWant:  "50,coin,-1",
			sellerWant: "20,regular arrow,1",
		},
		{
			request:    "10",
			coins:      "50",
			seller:     "20,apple,1\n18,arrow,4\n1,sword,40",
			wantReply:  "You can't buy coins. Make sure you request an item.",
			buyerWant:  "50,coin,-1",
			sellerWant: "20,apple,1\n18,arrow,4\n1,sword,40",
		},
		{
			request: "",
			coins:   "50",
			seller:  "20,apple,1\n18,arrow,4\n1,sword,40",
			wantReply: "```\n" +
				"╔═════════════════════════╗\n" +
				"║ Quantity  Item    Price ║\n" +
				"║─────────────────────────║\n" +
				"║ 20        apples  $1    ║\n" +
				"║ 18        arrows  $4    ║\n" +
				"║ 1         sword   $40   ║\n" +
				"╚═════════════════════════╝\n" +
				"```",
			buyerWant:  "50,coin,-1",
			sellerWant: "20,apple,1\n18,arrow,4\n1,sword,40",
		},
	}

	for _, tc := range tests {
		dir := t.TempDir()

		// Write the buyer's inventory.
		buyerPath := filepath.Join(dir, "buyer.csv")
		err := os.WriteFile(buyerPath, []byte(tc.coins+",coin,-1"), 0777)
		if err != nil {
			t.Fatal(err)
		}

		// Write the seller's inventory.
		sellerPath := filepath.Join(dir, "seller.csv")
		err = os.WriteFile(sellerPath, []byte(tc.seller), 0777)
		if err != nil {
			t.Fatal(err)
		}

		b := backpack{
			dir: dir,
		}
		reply := b.buyItem(tc.request, "buyer", "seller")
		if tc.wantReply != reply {
			t.Logf(
				"incorrect reply:\nwant:\n%v\ngot:\n%v\n",
				tc.wantReply,
				reply,
			)
			t.Fail()
		}

		buyerGot, err := os.ReadFile(buyerPath)
		if err != nil {
			t.Fatal(err)
		}

		if strings.TrimSpace(string(buyerGot)) != tc.buyerWant {
			t.Logf(
				"incorrect buyer inventory:\nwant:\n%v\ngot:\n%v\n",
				tc.buyerWant,
				string(buyerGot),
			)
			t.Fail()
		}

		sellerGot, err := os.ReadFile(sellerPath)
		if err != nil {
			t.Fatal(err)
		}

		if strings.TrimSpace(string(sellerGot)) != tc.sellerWant {
			t.Logf(
				"incorrect seller inventory:\nwant:\n%v\ngot:\n%v\n",
				tc.sellerWant,
				string(sellerGot),
			)
			t.Fail()
		}
	}
}

func TestModifyItem(t *testing.T) {
	type test struct {
		op      operation
		request string

		begin string
		want  string
	}

	tests := []test{
		{
			op:      opAdd,
			request: "apples",
			begin:   "",
			want:    "1,apple,-1",
		},
		{
			op:      opSet,
			request: "key lime pie -1",
			begin:   "1,key lime pie,10",
			want:    "1,key lime pie,-1",
		},
		{
			op:      opAdd,
			request: "apple",
			begin:   "",
			want:    "1,apple,-1",
		},
		{
			op:      opAdd,
			request: "100",
			begin:   "",
			want:    "100,coin,-1",
		},
		{
			op:      opAdd,
			request: "10 apples",
			begin:   "",
			want:    "10,apple,-1",
		},
		{
			op:      opDel,
			request: "apples",
			begin:   "",
			want:    "0,apple,-1",
		},
		{
			op:      opDel,
			request: "1 apples",
			begin:   "10,apple,-1",
			want:    "9,apple,-1",
		},
		{
			op:      opDel,
			request: "8 apples",
			begin:   "10,apple,-1",
			want:    "2,apple,-1",
		},
		{
			op:      opAdd,
			request: "apples",
			begin:   "1,apple,-1",
			want:    "2,apple,-1",
		},
		{
			op:      opAdd,
			request: "10 apples",
			begin:   "1,apple,-1",
			want:    "11,apple,-1",
		},
		{
			op:      opAdd,
			request: "apples",
			begin:   "1,pear,-1",
			want:    "1,pear,-1\n1,apple,-1",
		},
		{
			op:      opSet,
			request: "0 Mana Potions 10",
			begin:   "10,Mana Potion,-1",
			want:    "0,Mana Potion,10",
		},
		{
			op:      opSet,
			request: "0 Mana Potions",
			begin:   "10,Mana Potion,-1",
			want:    "0,Mana Potion,-1",
		},
	}

	for _, tc := range tests {
		dir := t.TempDir()
		path := filepath.Join(dir, "owner.csv")

		if tc.begin != "" {
			err := os.WriteFile(path, []byte(tc.begin), 0777)
			if err != nil {
				t.Fatal(err)
			}
		}

		b := backpack{
			dir: dir,
		}
		b.modifyItem(tc.request, "owner", tc.op)

		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		got := strings.TrimSpace(string(data))

		if tc.want != got {
			t.Fatalf(
				"want: %v got: %v\n",
				tc.want,
				got,
			)
		}
	}
}
