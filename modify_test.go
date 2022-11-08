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
			op:      opAdd,
			request: "5  apples",
			begin:   "",
			want:    "5,apple,-1",
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
