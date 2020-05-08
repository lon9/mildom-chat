package mildom

import (
	"log"
	"testing"
)

func TestEnter(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	if err := Listen(10445339); err != nil {
		t.Error(err)
	}
}
