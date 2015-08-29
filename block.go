package geoip

import (
    "net"
    "regexp"
    "strconv"

    "github.com/albertjin/bitidx"
)

var bit2loc = map[uint32]int{}

func init() {
    for x, i := uint32(1), int(32); i > 0; i-- {
        bit2loc[x] = i
        x <<= 1
    }
}

type BlockKind int

// Ip allocation Block.
type Block struct {
    bits *bitidx.Bits
    kind BlockKind
}

//go:generate stringer -type=BlockKind

const (
    BlockIPv6 BlockKind = 0
    BlockIPv4 BlockKind = 1
)

func (r *Block) Bits() *bitidx.Bits {
    return r.bits
}

func (r *Block) Kind() BlockKind {
    return r.kind
}

func (r *Block) String() string {
    return r.kind.String() + ": " + r.bits.String()
}

var reRIR = regexp.MustCompile(`[a-z]+\|([a-zA-Z]+)\|ipv([46])\|([0-9a-f\.:]+)\|([0-9]+)\|[0-9]+\|(allocated|assigned)`)

// RIR: Regional Internet Registry
// https://ftp.apnic.net/stats/apnic/
// examples,
//   apnic|CN|ipv4|153.36.0.0|131072|20110331|allocated
//   apnic|CN|ipv6|240c:8000::|21|20140905|allocated
//   lacnic|VE|ipv4|190.168.0.0|32768|20070329|assigned
//   lacnic|VE|ipv4|190.168.128.0|16384|20070419|assigned
func ParseRIR(line string) (r *Block, name string) {
    if ss := reRIR.FindStringSubmatch(line); len(ss) == 6 {
        if ip := net.ParseIP(ss[3]); ip != nil {
            mc, _ := strconv.ParseUint(ss[4], 10, 64)
            switch ss[2] {
            case "4":
                if count, ok := bit2loc[uint32(mc)]; ok {
                    r = &Block{bitidx.NewBits([]byte(ip.To4()), count), BlockIPv4}
                    name = ss[1]
                }
            case "6":
                r = &Block{bitidx.NewBits([]byte(ip), int(mc)), BlockIPv6}
                name = ss[1]
            }
        }
    }

    return
}
