package mildom

import (
	"log"
	"testing"
)

func TestEnter(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	listener, err := GetListener(10467370)
	if err != nil {
		t.Fatal(err)
	}
	for msg := range listener {
		t.Log(msg)
	}
}
