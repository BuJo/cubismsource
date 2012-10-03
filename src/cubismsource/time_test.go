package main

import "time"
import "testing"

// This tests various time parsing/formatting because this is poorly documented
func TestTime(t *testing.T) {
	t1 := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	s1 := t1.Format("2006-01-02 10:04:00")
	s2 := t1.Format("2006-01-02 15:04:00-07:00")
	s3 := t1.Format("2006-01-02 15:04:05.999999999 -0700 MST")
	s4 := t1.Format(time.RFC3339)
	if s1 != "2009-11-10 110:00:00" {
		t.Errorf("Expected Bad Minutes: %v, want %v", s1, "2009-11-10 110:00:00")
	}
	if s2 != "2009-11-10 23:00:00+00:00" {
		t.Errorf("Expected good parse: %v, want %v", s2, "2009-11-10 23:00:00+00:00")
	}
	if s3 != "2009-11-10 23:00:00 +0000 UTC" {
		t.Errorf("Expected good parse: %v, want %v", s3, "2009-11-10 23:00:00 +0000 UTC")
	}
	if s4 != "2009-11-10T23:00:00Z" {
		t.Errorf("Expected good parse: %v, want %v", s4, "2009-11-10T23:00:00Z")
	}

	t2, _ := time.Parse("2006-01-02T15:04:05.999Z", "2012-10-01T15:15:40.000Z")
	s5 := t2.Format(time.RFC3339)
	if s5 != "0001-01-01T00:00:00Z" {
		t.Errorf("Expected Bad Parse/Format: %v, want %v", s5, "0001-01-01T00:00:00Z")
	}

	t2, _ = time.Parse("2006-01-02T15:04:05.000Z", "2012-10-01T15:48:40.000Z")
	s6 := t2.Format(time.RFC3339)
	if s6 != "2012-10-01T15:48:40Z" {
		t.Errorf("Expected good parse/format: %v, want %v", s6, "2012-10-01T15:48:40Z")
	}

	if t2.Month() != 10 {
		t.Errorf("Expected month to be 10, was: %d\n", t2.Month())
	}
	if t2.Year() != 2012 {
		t.Errorf("Expected year to be 2012, was: %d\n", t2.Year())
	}
	if t2.Day() != 1 {
		t.Errorf("Expected day to be 1, was: %d\n", t2.Day())
	}
}
