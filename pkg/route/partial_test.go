package route

import (
	"testing"
)

func TestPartialPart(t *testing.T) {
	pp0, err := parse_partialEndPart(`/[part]{a-zA-Z++`)
	if err == nil {
		t.Error("parse_partialEndPart should fail when the partial is invalid")
		t.Errorf("%+v", pp0.subPart)
	}
	pp0, err = parse_partialEndPart(`/[part]{a-zA-Z+}+`)
	if err != nil {
		t.Error(err)
	}
	// Eq
	pp1, _ := parse_partialEndPart(`/[part]{a-zA-Z+}+`)
	if !pp0.Eq(pp1) {
		t.Error("pp0 should Eq pp1")
	}
	Part(pp1).(paramPart).SetParameterName("otherpart")
	if pp1.ParameterName() != "otherpart" {
		t.Errorf("expected otherpart, got %s", pp1.ParameterName())
	}
	if pp0.Eq(pp1) {
		t.Error("pp0 should not Eq pp1 when param is different")
	}
	pp1, _ = parse_partialEndPart(`/[part]{a-zA-Z_+}+`)
	if pp0.Eq(pp1) {
		t.Error("pp0 should not Eq pp1 when expr is different")
	}
	sp1, _ := build_stringPart("/part")
	if pp0.Eq(sp1) {
		t.Error("pp0 should not Eq a non-partial Part")
	}
	// Expr
	if expr := pp0.Expr(); expr != "*" {
		t.Errorf("pp0 should have expr '*', got '%s'", expr)
	}
}
