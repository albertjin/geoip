package geoip

import (
    "bufio"
    "encoding/json"
    "io"
    "net"

    "github.com/albertjin/bitidx"
)

type Index struct {
    index []bitidx.Node
    ids map[string]int
    names []string
    lastId int
}

type jsonIndex struct {
    Index bitidx.Node `json:"index"`
    Ids map[string]int `json:"ids"`
    Names []string `json:"names"`
    LastId int `json:"lastId"`
    Extra interface{} `json:"extra"`
}

// New Index object.
func NewIndex() *Index {
    v6 := bitidx.NewNode(0)
    z := []byte(net.IPv4zero.To16())
    v6.Put(bitidx.NewBits(z, 12*8+1), 0, false)
    v4, _, _ := v6.Find(bitidx.NewBits(z, 12*8))

    return &Index{[]bitidx.Node{v6, v4}, map[string]int{"": 0}, []string{""}, 0}
}

// Restore structure from json previously exported.
func NewIndexFromJson(rd io.Reader) (index *Index, extra interface{}) {
    var v *jsonIndex
    dec := json.NewDecoder(rd)
    if dec.Decode(&v) == nil {
        v.Index.ConsolidateNum()
        s, _, _ := v.Index.Find(bitidx.NewBits([]byte(net.IPv4zero.To16()), 12*8))
        return &Index{[]bitidx.Node{v.Index, s}, v.Ids, v.Names, v.LastId}, v.Extra
    }
    return nil, nil
}

// Put a block with name.
func (index *Index) Put(block *Block, name string) {
    id, ok := index.ids[name]
    if !ok {
        index.lastId++
        id = index.lastId
        index.ids[name] = id
        index.names = append(index.names, name)
    }

    index.index[block.Kind()].Put(block.Bits(), id, false)
}

// Find ip for name. When no ip is found empty string is returned.
func (index *Index) Find(ip string) (name string, length int) {
    if i := net.ParseIP(ip); i != nil {
        var n bitidx.Node
        if j := i.To4(); j != nil {
            n = index.index[BlockIPv4]
            i = j
        } else {
            n = index.index[BlockIPv6]
        }

        _, x, l := n.Find(bitidx.NewBits([]byte(i), -1))
        if id, ok := x.(int); ok {
            name = index.names[id]
            length = l
        }
    }

    return
}

// Load RIR records into index.
func (index *Index) LoadRIR(rd io.Reader) {
    i := bufio.NewReader(rd)
    for {
        line, err := i.ReadString('\n')
        if err != nil {
            break
        }
        if r, n := ParseRIR(line); r != nil {
            index.Put(r, n)
        }
    }
}

// Marshal content in json and export with extra data.
func (index *Index) ToJson(w io.Writer, extra interface{}) (err error) {
    enc := json.NewEncoder(w)
    return enc.Encode(&jsonIndex{index.index[BlockIPv6], index.ids, index.names, index.lastId, extra})
}
