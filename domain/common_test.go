package domain

import (
	"testing"
)

func TestGetEthernetTyp(t *testing.T) {

	etype := make([]byte, 20)
	etype[12] = 12 // 12 ==> c
	etype[13] = 13 // 13 ==> d

	ethertyp := getEthernetTyp(etype)

	if ethertyp != "cd" {
		t.Errorf("ethertyp decoding failed. Decoding: %s", ethertyp)
	}
}
