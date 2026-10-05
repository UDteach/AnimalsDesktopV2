package main

import (
	"reflect"
	"testing"
)

func TestSelectedFrameNumbersFromExplicitList(t *testing.T) {
	got, err := selectedFrameNumbers("0,4,12,26,32,40,52,56", 99, 10)
	if err != nil {
		t.Fatalf("selectedFrameNumbers() error = %v", err)
	}
	want := []int{0, 4, 12, 26, 32, 40, 52, 56}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("selectedFrameNumbers() = %v, want %v", got, want)
	}
}

func TestSelectedFrameNumbersFromRange(t *testing.T) {
	got, err := selectedFrameNumbers("", 4, 12)
	if err != nil {
		t.Fatalf("selectedFrameNumbers() error = %v", err)
	}
	want := []int{12, 13, 14, 15}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("selectedFrameNumbers() = %v, want %v", got, want)
	}
}

func TestParseFrameListRejectsDuplicate(t *testing.T) {
	if _, err := parseFrameList("0,4,4"); err == nil {
		t.Fatalf("parseFrameList() succeeded for duplicate frame")
	}
}

func TestParseFrameListRejectsNegative(t *testing.T) {
	if _, err := parseFrameList("0,-1,4"); err == nil {
		t.Fatalf("parseFrameList() succeeded for negative frame")
	}
}
