// License: AGPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestModifyItem(t *testing.T) {
	type test struct {
		op    string
		count int
		price int
		item  string

		begin string
		want  string
	}

	tests := []test{
		{
			op:    "add",
			count: 1,
			item:  "apples",
			price: -1,
			begin: "",
			want:  "1,apple,-1",
		},
		{
			op:    "add",
			count: 5,
			item:  "apples",
			price: -1,
			begin: "",
			want:  "5,apple,-1",
		},
		{
			op:    "set",
			count: 1,
			item:  "key lime pie",
			price: -1,
			begin: "1,key lime pie,10",
			want:  "1,key lime pie,-1",
		},
		{
			op:    "add",
			count: 1,
			item:  "apple",
			price: -1,
			begin: "",
			want:  "1,apple,-1",
		},
		{
			op:    "remove",
			item:  "apples",
			price: -1,
			begin: "",
			want:  "0,apple,-1",
		},
		{
			op:    "remove",
			count: 1,
			item:  "apples",
			price: -1,
			begin: "10,apple,-1",
			want:  "9,apple,-1",
		},
		{
			op:    "remove",
			count: 8,
			item:  "apples",
			price: -1,
			begin: "10,apple,-1",
			want:  "2,apple,-1",
		},
		{
			op:    "add",
			count: 1,
			item:  "apples",
			price: -1,
			begin: "1,apple,-1",
			want:  "2,apple,-1",
		},
		{
			op:    "add",
			count: 10,
			item:  "apples",
			price: -1,
			begin: "1,apple,-1",
			want:  "11,apple,-1",
		},
		{
			op:    "add",
			count: 1,
			item:  "apples",
			price: -1,
			begin: "1,pear,-1",
			want:  "1,pear,-1\n1,apple,-1",
		},
		{
			op:    "set",
			count: 0,
			price: 10,
			item:  "Mana Potions",
			begin: "10,Mana Potion,-1",
			want:  "0,Mana Potion,10",
		},
		{
			op:    "set",
			count: 0,
			item:  "Mana Potions",
			price: -1,
			begin: "10,Mana Potion,-1",
			want:  "0,Mana Potion,-1",
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
		b.modifyItem(tc.count, tc.price, tc.item, "owner", tc.op)

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
