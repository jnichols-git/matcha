package route

import "testing"

func TestStringPart(t *testing.T) {
	sp1, _ := build_stringPart("/sp")
	if ok := sp1.Match(nil, "/sp"); !ok {
		t.Errorf("string part /sp should Match token /sp")
	}
	if ok := sp1.Match(nil, "/notsp"); ok {
		t.Errorf("string part /sp should not Match token /notsp")
	}
	if sp1.val != "/sp" {
		t.Errorf("string part /sp should have val /sp, got %s", sp1.val)
	}
	sp2, _ := build_stringPart("/sp")
	if !sp1.Eq(sp2) {
		t.Errorf("string part /sp should Eq string part /sp")
	}
	sp3, _ := build_stringPart("/notsp")
	if sp1.Eq(sp3) {
		t.Errorf("string part /sp should not Eq string part /notsp")
	}
	rp, _ := build_regexPart("", ".*")
	if sp1.Eq(rp) {
		t.Errorf("string part /sp should not Eq non-string part")
	}
}

func TestWildcardPart(t *testing.T) {
	tp0, err := build_wildcardPart("wc")
	if err != nil {
		t.Errorf("wildcard parts should never return an error on build")
	}
	rmc := newRMC()
	rmc.Allocate("wc")
	if ok := tp0.Match(rmc, "/iahs"); !ok {
		t.Errorf("wildcard parts should match any value")
	} else if param := GetParam(rmc, "wc"); param != "iahs" {
		t.Errorf("expected wildcard to store param 'iahs', got '%s'", param)
	}
	if ok := tp0.Match(rmc, "/"); !ok {
		t.Errorf("wildcard parts should match the root")
	} else if param := GetParam(rmc, "wc"); param != "" {
		t.Errorf("expected wildcard to store param '', got '%s'", param)
	}
	tp1, _ := build_wildcardPart("wc")
	if !tp0.Eq(tp1) {
		t.Error("wildcards with equal params should be Eq")
	}
	tp2, _ := build_wildcardPart("notwc")
	if tp0.Eq(tp2) {
		t.Error("wildcards with different params should not be Eq")
	}
	sp1, _ := build_stringPart("/test")
	if tp0.Eq(sp1) {
		t.Errorf("pr1 should not Eq string part")
	}
	Part(tp2).(paramPart).SetParameterName("wc")
	if !tp0.Eq(tp2) {
		t.Error("wildcards should have same param after SetParameterName")
	}
}

func TestRegexPart(t *testing.T) {
	rp0, err := build_regexPart("", "(.*")
	if err == nil || rp0 != nil {
		t.Errorf("regex part should fail with invalid regex")
	}
	rp1, err := build_regexPart("", "[a-z]{4}")
	if err != nil || rp1 == nil {
		t.Errorf("regex part should build with valid regex, got %s", err)
	}
	if rp1.expr.String() != "[a-z]{4}" {
		t.Errorf("regex part expr.String should be equal to original build expression %s, got %s", "[a-z]{4}", rp1.expr.String())
	}
	if ok := rp1.Match(nil, "/word"); !ok {
		t.Errorf("rp1 should Match token /word")
	}
	if ok := rp1.Match(nil, "/longword"); ok {
		t.Errorf("rp1 should not Match token /longword")
	}
	rp2, _ := build_regexPart("", "[a-z]{4}")
	if !rp1.Eq(rp2) {
		t.Errorf("rp1 should Eq regex part with same param/expr")
	}
	rp3, _ := build_regexPart("param", "[a-z]{4}")
	if rp1.Eq(rp3) {
		t.Errorf("rp1 should not Eq regex part with same expr but different param")
	}
	rp4, _ := build_regexPart("", "[a-z]{5}")
	if rp1.Eq(rp4) {
		t.Errorf("rp1 should not Eq regex part with same param but different expr")
	}
	sp1, _ := build_stringPart("/test")
	if rp1.Eq(sp1) {
		t.Errorf("pr1 should not Eq string part")
	}
	rmc := newRMC()
	rmc.Allocate("param")
	rp2.Match(rmc, "/word")
	if GetParam(rmc, "param") != "" {
		t.Error("param should not be set with empty-param regexPart")
	}
	rp3.Match(rmc, "/word")
	if GetParam(rmc, "param") != "word" {
		t.Errorf("expected param 'word', got '%s'", GetParam(rmc, "param"))
	}
}
