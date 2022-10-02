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
			want:    "1,apples,-1",
		},
		{
			op:      opAdd,
			request: "100",
			begin:   "",
			want:    "100,coins,-1",
		},
		{
			op:      opAdd,
			request: "10 apples",
			begin:   "",
			want:    "10,apples,-1",
		},
		{
			op:      opDel,
			request: "apples",
			begin:   "",
			want:    "0,apples,-1",
		},
		{
			op:      opDel,
			request: "1 apples",
			begin:   "10,apples,-1",
			want:    "9,apples,-1",
		},
		{
			op:      opDel,
			request: "8 apples",
			begin:   "10,apples,-1",
			want:    "2,apples,-1",
		},
		{
			op:      opAdd,
			request: "apples",
			begin:   "1,apples,-1",
			want:    "2,apples,-1",
		},
		{
			op:      opAdd,
			request: "10 apples",
			begin:   "1,apples,-1",
			want:    "11,apples,-1",
		},
		{
			op:      opAdd,
			request: "apples",
			begin:   "1,pears,-1",
			want:    "1,pears,-1\n1,apples,-1",
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
			t.Fatalf("want: %v got: %v\n", tc.want, got)
		}
	}
}
