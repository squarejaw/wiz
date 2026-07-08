package cmd

import (
	"encoding/json"
	"net"
	"os"
	"testing"
)

type fakeResp struct {
	addr net.Addr
	data []byte
}

type fakeReader struct {
	responses []fakeResp
	i         int
}

func (f *fakeReader) Read(buf []byte) (int, net.Addr, error) {
	if f.i >= len(f.responses) {
		return 0, nil, os.ErrDeadlineExceeded
	}
	r := f.responses[f.i]
	f.i++
	return copy(buf, r.data), r.addr, nil
}

// bogusAddr is a net.Addr that is not a *net.UDPAddr, used to verify that
// collectBulbs skips unexpected address types without panicking.
type bogusAddr struct{}

func (bogusAddr) Network() string { return "bogus" }
func (bogusAddr) String() string  { return "bogus" }

func ipAddr(ip string) *net.UDPAddr {
	return &net.UDPAddr{IP: net.ParseIP(ip)}
}

func TestCollectBulbsDedup(t *testing.T) {
	r := &fakeReader{responses: []fakeResp{
		{addr: ipAddr("192.168.1.12"), data: []byte(`{"result":{"mac":"6ab731ba1bd5"}}`)},
		{addr: ipAddr("192.168.1.34"), data: []byte(`{"result":{"mac":"86082d45c947"}}`)},
		// same MAC from a different IP — must be deduped
		{addr: ipAddr("192.168.1.99"), data: []byte(`{"result":{"mac":"6ab731ba1bd5"}}`)},
		// no MAC in response — must be skipped
		{addr: ipAddr("192.168.1.50"), data: []byte(`{"result":{}}`)},
	}}
	bulbs, err := collectBulbs(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bulbs) != 2 {
		t.Fatalf("expected 2 bulbs, got %d: %+v", len(bulbs), bulbs)
	}
	if *bulbs[0].IP != "192.168.1.12" || *bulbs[0].Mac != "6ab731ba1bd5" {
		t.Errorf("first bulb = %s/%s, want 192.168.1.12/6ab731ba1bd5", *bulbs[0].IP, *bulbs[0].Mac)
	}
	if *bulbs[1].IP != "192.168.1.34" || *bulbs[1].Mac != "86082d45c947" {
		t.Errorf("second bulb = %s/%s, want 192.168.1.34/86082d45c947", *bulbs[1].IP, *bulbs[1].Mac)
	}
}

func TestCollectBulbsEmptyMarshalsAsArray(t *testing.T) {
	bulbs, err := collectBulbs(&fakeReader{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := json.Marshal(bulbs)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != "[]" {
		t.Errorf("got %s, want []", out)
	}
}

func TestCollectBulbsNonUDPAddrSkipped(t *testing.T) {
	r := &fakeReader{responses: []fakeResp{
		{addr: bogusAddr{}, data: []byte(`{"result":{"mac":"6ab731ba1bd5"}}`)},
	}}
	// A non-*net.UDPAddr addr should be skipped, not panic.
	bulbs, err := collectBulbs(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bulbs) != 0 {
		t.Fatalf("expected 0 bulbs, got %d", len(bulbs))
	}
}
