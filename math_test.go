package assert

import "testing"

func TestCompareAbs(t *testing.T) {
	if !compareAbs(12, 14, 4) {
		t.Errorf("compareAbs failed")
	}

	if compareAbs(12, 19, 4) {
		t.Errorf("compareAbs failed")
	}
}

func TestCompareRel(t *testing.T) {
	if !compareRel(12, 14, 4) {
		t.Errorf("compareRel failed")
	}

	if !compareRel(12, 19, 4) {
		t.Errorf("compareRel failed")
	}
}
