package geoip

import(
    "net"
    "testing"

    "github.com/albertjin/bitidx"
)

func TestParseRIR(t *testing.T) {
    a, x := ParseRIR("apnic|CN|ipv4|153.36.0.0|131072|20110331|allocated")
    b, y := ParseRIR("apnic|CN|ipv6|240c:8000::|21|20140905|allocated")

    t.Log(a, x)
    t.Log(b, y)

    n := bitidx.NewNode(0)
    n.Put(a.Bits(), 100, false)
    t.Log(n.Find(bitidx.NewBits([]byte((net.IP{153,37,0,1}).To16()), -1)))
    t.Log(n.Find(bitidx.NewBits([]byte(net.IPv4zero.To16()), 12*8)))
}

func TestParseRIR1(t *testing.T) {
    a, x := ParseRIR("apnic|CN|ipv4|153.36.0.0|131072|20110331|allocated")
    b, y := ParseRIR("apnic|CN|ipv6|240c:8000::|21|20140905|allocated")

    t.Log(a, x)
    t.Log(b, y)

    r := bitidx.NewNode(0)
    r.Put(bitidx.NewBits([]byte(net.IPv4zero.To16()), 12*8+1), 0, false)
    s, _ := r.Find(bitidx.NewBits([]byte(net.IPv4zero.To16()), 12*8))

    rs := []bitidx.Node{r, s}

    rs[a.Kind()].Put(a.Bits(), 10, false)
    rs[b.Kind()].Put(b.Bits(), 100, false)

    t.Log(r)
    t.Log(s)

    t.Log(rs[BlockIPv4].Find(bitidx.NewBits([]byte((net.IP{153,37,0,1}).To4()), -1)))
}
