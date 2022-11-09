package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuyItem(t *testing.T) {
	type test struct {
		count  int
		item   string
		coins  string
		seller string

		wantReply  string
		buyerWant  string
		sellerWant string
	}

	tests := []test{
		{
			count:  10,
			item:   "apples",
			coins:  "50",
			seller: "20,apple,1",
			wantReply: "buyer bought 10 apples for $10\n" +
				"buyer has 10 apples\n" +
				"seller has 10 apples for sale for $1",
			buyerWant:  "40,coin,-1\n10,apple,-1",
			sellerWant: "10,apple,1\n10,coin,-1",
		},
		{
			count:  10,
			item:   "apples",
			coins:  "50",
			seller: "20,apple,1",
			wantReply: "buyer bought 10 apples for $10\n" +
				"buyer has 10 apples\n" +
				"seller has 10 apples for sale for $1",
			buyerWant:  "40,coin,-1\n10,apple,-1",
			sellerWant: "10,apple,1\n10,coin,-1",
		},
		{
			count:  10,
			item:   "apples",
			coins:  "50",
			seller: "20,apple,1\n2,coin,-1",
			wantReply: "buyer bought 10 apples for $10\n" +
				"buyer has 10 apples\n" +
				"seller has 10 apples for sale for $1",
			buyerWant:  "40,coin,-1\n10,apple,-1",
			sellerWant: "10,apple,1\n12,coin,-1",
		},
		{
			count:  10,
			item:   "massive catapults",
			coins:  "5000",
			seller: "11,massive catapult,150",
			wantReply: "buyer bought 10 massive catapults for $1500\n" +
				"buyer has 10 massive catapults\n" +
				"seller has 1 massive catapult for sale for $150",
			buyerWant:  "3500,coin,-1\n10,massive catapult,-1",
			sellerWant: "1,massive catapult,150\n1500,coin,-1",
		},
		{
			count:  1,
			item:   "apple",
			coins:  "50",
			seller: "",
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
			count:  1,
			item:   "apple",
			coins:  "50",
			seller: "1,apple,-1\n10,arrow,-1\n1,sword,100",
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
			count:  10,
			item:   "apples",
			coins:  "50",
			seller: "1,apple,1",
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
			count:  10,
			item:   "apples",
			coins:  "2",
			seller: "100,apple,1",
			wantReply: "buyer has insufficient funds\n" +
				"10 apples costs $10\n" +
				"buyer only has 2 coins",
			buyerWant:  "2,coin,-1",
			sellerWant: "100,apple,1",
		},
		{
			count:      0,
			item:       "apples",
			coins:      "50",
			seller:     "20,apple,1",
			wantReply:  "You can't buy 0 of an item, silly!",
			buyerWant:  "50,coin,-1",
			sellerWant: "20,apple,1",
		},
		{
			count:  -10,
			item:   "regular arrows",
			coins:  "50",
			seller: "20,regular arrow,1",
			wantReply: "You've requested to give away your items?\n" +
				"Try again with: \"10 regular arrows\"",
			buyerWant:  "50,coin,-1",
			sellerWant: "20,regular arrow,1",
		},
		{
			count:      10,
			item:       "coins",
			coins:      "50",
			seller:     "10,coins,1\n20,apple,1\n18,arrow,4\n1,sword,40",
			wantReply:  "You can't buy coins silly!",
			buyerWant:  "50,coin,-1",
			sellerWant: "10,coins,1\n20,apple,1\n18,arrow,4\n1,sword,40",
		},
		{
			count:      10,
			item:       "",
			coins:      "50",
			seller:     "10,coins,1\n20,apple,1\n18,arrow,4\n1,sword,40",
			wantReply:  "You forgot to request an item.",
			buyerWant:  "50,coin,-1",
			sellerWant: "10,coins,1\n20,apple,1\n18,arrow,4\n1,sword,40",
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
		reply := b.buyItem(tc.count, tc.item, "buyer", "seller")
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
