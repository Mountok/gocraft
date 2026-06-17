package update

import "testing"

func TestIsOutdated(t *testing.T) {
	tests := []struct {
		current string
		latest  string
		want    bool
	}{
		{current: "v1.0.0", latest: "v1.1.0", want: true},
		{current: "1.1.0", latest: "v1.1.0", want: false},
		{current: "v1.2.0", latest: "v1.1.0", want: false},
		{current: "dev", latest: "v1.1.0", want: false},
	}

	for _, test := range tests {
		got := isOutdated(test.current, test.latest)
		if got != test.want {
			t.Fatalf("isOutdated(%q, %q) = %v, want %v", test.current, test.latest, got, test.want)
		}
	}
}
