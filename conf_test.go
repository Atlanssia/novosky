package main

import "testing"

func testInit(t *testing.T) {
	initConf()
	if len(appConf.Uris) <= 0 {
		t.Fail()
	}
}
