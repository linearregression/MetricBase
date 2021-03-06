package http

import (
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestParseHttpTimespan_Table(t *testing.T) {
	// Generate requests and expected output

	var timespan_tests = []struct {
		queryString string
		start, end  time.Time
	}{
		// Intervals
		{"interval=2014",
			time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2014, 12, 31, 23, 59, 59, 0, time.UTC)},

		// Start and end times
		{"start=2013&end=2014",
			time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2014, 12, 31, 23, 59, 59, 0, time.UTC)},
		{"start=201405&end=20140526",
			time.Date(2014, 5, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2014, 5, 26, 23, 59, 59, 0, time.UTC)},
	}

	for _, tt := range timespan_tests {
		url, _ := url.Parse("/rpc/query?" + tt.queryString)
		req := &http.Request{
			URL: url,
		}
		req.ParseForm()

		start, end, err := ParseHttpTimespan(req)
		if err != nil {
			t.Errorf("Got unexpected error parsing %s: %s", tt.queryString, err)
			continue
		}

		if !start.Equal(tt.start) {
			t.Errorf("Expected %v to start at %v, but got %v.", tt.queryString, tt.start, start)
		}

		if !end.Equal(tt.end) {
			t.Errorf("Expected %v to end at %v, but got %v.", tt.queryString, tt.end, end)
		}
	}
}

func TestparseInterval(t *testing.T) {
	now := time.Now()
	var timespan_tests = []struct {
		interval   string
		start, end time.Time
	}{
		// Intervals
		{"2014",
			time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2014, 12, 31, 23, 59, 59, 0, time.UTC)},
		{"201405",
			time.Date(2014, 5, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2014, 5, 31, 23, 59, 59, 0, time.UTC)},
		{"20140526",
			time.Date(2014, 5, 26, 0, 0, 0, 0, time.UTC),
			time.Date(2014, 5, 26, 23, 59, 59, 0, time.UTC)},

		{"-1w",
			now.Add(-time.Hour * 24 * 7),
			now},
		{"1h10s",
			now.Add(time.Hour + time.Second*10),
			now},

		// TODO: Tests with months that doesn't have 31 days
		/*
			{"interval=201402",
			time.Date(2014, 2, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2014, 2, 28, 23, 59, 59, 0, time.UTC)},
			// TODO: Leap years
		*/
	}

	for _, tt := range timespan_tests {
		start, end, err := parseInterval(tt.interval, now)
		if err != nil {
			t.Errorf("Got unexpected error parsing %s: %s", tt.interval, err)
			continue
		}

		if !start.Equal(tt.start) {
			t.Errorf("Expected %v to start at %v, but got %v.", tt.interval, tt.start, start)
		}

		if !end.Equal(tt.end) {
			t.Errorf("Expected %v to end at %v, but got %v.", tt.interval, tt.end, end)
		}
	}
}

func TestparseIntervalErrors(t *testing.T) {
	now := time.Now()
	var timespan_errors = []string{
		"1",
		"20145",
	}

	for _, tt := range timespan_errors {
		_, _, err := parseInterval(tt, now)
		if err == nil {
			t.Errorf("Expected an error for input '%s', got nothing.", tt)
		}
	}
}
