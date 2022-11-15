package main

import "testing"

func TestDisplayName(t *testing.T) {
	type test struct {
		name  string
		count int
		want  string
	}

	tests := []test{
		{
			name:  "apple",
			count: 1,
			want:  "Apple",
		},
		{
			name:  "apple pie",
			count: 1,
			want:  "Apple pie",
		},
		{
			name:  "apple-pie",
			count: 1,
			want:  "Apple-pie",
		},
		{
			name:  "apple",
			count: 2,
			want:  "Apples",
		},
		{
			name:  "apple pie",
			count: 3,
			want:  "Apple pies",
		},
		{
			name:  "apple-pie",
			count: 4,
			want:  "Apple-pies",
		},
	}

	for _, test := range tests {
		got := displayName(test.name, test.count)
		if got != test.want {
			t.Fatalf("want: %v got: %v\n", test.want, got)
		}
	}
}
