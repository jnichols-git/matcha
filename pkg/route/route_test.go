package route

import (
	"net/http"
	"testing"
)

func TestStringRoute(t *testing.T) {
	rt, err := New("/test")
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest(http.MethodGet, "http://url.com/test", nil)
	if req = rt.MatchAndUpdateContext(req); req == nil {
		t.Errorf("Expected route to match")
	}
}

func TestWildcardRoute(t *testing.T) {
	rt, err := New("/[any]")
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest(http.MethodGet, "http://url.com/test", nil)
	if req = rt.MatchAndUpdateContext(req); req == nil {
		t.Errorf("Expected route to match")
	}
}

func TestRegexRoute(t *testing.T) {
	rt, err := New("/{[a-zA-Z]{4}}")
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest(http.MethodGet, "http://url.com/test", nil)
	if req = rt.MatchAndUpdateContext(req); req == nil {
		t.Errorf("Expected route to match")
	}
	req, _ = http.NewRequest(http.MethodGet, "http://url.com/t3st", nil)
	if req = rt.MatchAndUpdateContext(req); req != nil {
		t.Errorf("Expected route to not match")
	}
}

func TestMethodRoute(t *testing.T) {
	rt, err := New("/test", WithMethods(http.MethodPost))
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest(http.MethodPost, "http://url.com/test", nil)
	if req = rt.MatchAndUpdateContext(req); req == nil {
		t.Errorf("Expected POST to match")
	}
	req, _ = http.NewRequest(http.MethodGet, "http://url.com/test", nil)
	if req = rt.MatchAndUpdateContext(req); req != nil {
		t.Errorf("Expected GET not to match")
	}
}
