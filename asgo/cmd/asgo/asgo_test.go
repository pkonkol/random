package main

import (
	// "go.mongodb.org/mongo-driver/bson"
	"testing"
)

// var fullFlag = flag.Bool("full", false, "start full scanning")
// var testFlag = flag.String("test", "", "test the tests for flags")

func TestTesting(t *testing.T) {
	t.Log("TEST WORKS")
}

func TestFlags(t *testing.T) {
	t.Log(*fullFlag)
	t.Log(*testFlag)
}
