package helper

import "testing"

func TestNetworkForAddress(t *testing.T) {
	testCases := []struct {
		address string
		want    string
	}{
		{address: "0.0.0.0:7777", want: "tcp4"},
		{address: "127.0.0.1:7777", want: "tcp4"},
		{address: "[::1]:7777", want: "tcp6"},
		{address: ":7777", want: "tcp"},
		{address: "localhost:7777", want: "tcp"},
	}

	for _, testCase := range testCases {
		got := NetworkForAddress(testCase.address)
		if got != testCase.want {
			t.Fatalf("NetworkForAddress(%q) = %q, want %q", testCase.address, got, testCase.want)
		}
	}
}

func TestCopySlice(t *testing.T) {
	values := []int{1, 2, 3}

	got := CopySlice(values)
	got[0] = 9

	if values[0] != 1 {
		t.Fatalf("CopySlice changed original slice: got %v", values)
	}
}

func TestTruncateText(t *testing.T) {
	got := TruncateText("127.0.0.1:7777", 10)
	want := "127.0.0.1."
	if got != want {
		t.Fatalf("TruncateText = %q, want %q", got, want)
	}
}
