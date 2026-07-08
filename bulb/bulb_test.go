package bulb

import (
	"encoding/json"
	"testing"
)

func strptr(s string) *string { return &s }

func TestMessageMarshalOmitsUnsetParams(t *testing.T) {
	m := Message{Method: strptr("setState"), Params: &Params{}}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	const want = `{"method":"setState","params":{}}`
	if string(b) != want {
		t.Errorf("got %s, want %s", b, want)
	}
}

func TestMessageMarshalTemp(t *testing.T) {
	temp := 2700
	m := Message{Method: strptr("setState"), Params: &Params{Temp: &temp}}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	const want = `{"method":"setState","params":{"temp":2700}}`
	if string(b) != want {
		t.Errorf("got %s, want %s", b, want)
	}
}

// TestBulbParamsTag guards against the Params field being mis-tagged
// (it was once "method" instead of "params"); the list command serializes
// Bulb values to JSON, so the tag must be correct.
func TestBulbMarshal(t *testing.T) {
	ip := "192.168.1.12"
	mac := "6ab731ba1bd5"
	b := Bulb{IP: &ip, Mac: &mac}
	out, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	const want = `{"ip":"192.168.1.12","mac":"6ab731ba1bd5"}`
	if string(out) != want {
		t.Errorf("got %s, want %s", out, want)
	}
}
