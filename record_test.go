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
