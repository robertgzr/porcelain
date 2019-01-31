package main

import (
	"fmt"
	"strings"
	"testing"
)

func printVersion(v GitVersion) string {
	return fmt.Sprintf("%d.%d.%d", v[0], v[1], v[2])
}

func TestParseVersion(t *testing.T) {
	r := strings.NewReader("git version 2.20.1")
	p := NewParser(r)
	v, err := p.ParseVersion()
	if err != nil {
		t.Fatal(err)
	}
	if v[0] != 2 || v[1] != 20 || v[2] != 1 {
		t.Fatalf("wrong version, expected 2.20.1, got %s", printVersion(v))
	}
	if isMinimumVersion(v) == false {
		t.Fatalf("minium version detection failed, %s < 2.20.1", printVersion(v))
	}
}

func TestParseVersionFail(t *testing.T) {
	r := strings.NewReader("git version 1.8.13")
	p := NewParser(r)
	v, err := p.ParseVersion()
	if err != nil {
		t.Fatal(err)
	}
	if v[0] != 1 || v[1] != 8 || v[2] != 13 {
		t.Fatalf("wrong version, expected 1.8.13, got %s", printVersion(v))
	}
	if isMinimumVersion(v) == true {
		t.Fatalf("minium version detection failed, %s < 2.20.1", printVersion(v))
	}
}
