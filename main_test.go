package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestProcess(t *testing.T) {
	testCases := []struct {
		value    string
		offset   time.Duration
		expected string
	}{
		{"00:00:00,000 --> 00:02:34,319", 5 * time.Second, "00:00:05,000 --> 00:02:39,319\n"},
		{"00:01:59,000 --> 00:02:34,319", 5 * time.Minute, "00:06:59,000 --> 00:07:34,319\n"},
		{"00:01:59,000 --> 00:02:34,319", 1 * time.Hour, "01:01:59,000 --> 01:02:34,319\n"},
		{"00:01:59,000 --> 00:02:34,319", -3 * time.Second, "00:01:56,000 --> 00:02:31,319\n"},
		{"01:01:59,000 --> 01:02:34,319", -2 * time.Minute, "00:59:59,000 --> 01:00:34,319\n"},
		{"this string should remain unchanged except new line at end", -2 * time.Minute, "this string should remain unchanged except new line at end\n"},
	}

	for _, tc := range testCases {
		var (
			r = strings.NewReader(tc.value)
			w bytes.Buffer
		)

		if err := process(r, &w, tc.offset); err != nil {
			t.Error(err)
		}

		output := w.String()
		if strings.Compare(output, tc.expected) != 0 {
			t.Errorf("expected %s; got %s", tc.expected, output)
		}
	}
}
