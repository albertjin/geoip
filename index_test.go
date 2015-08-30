package geoip

import(
    "compress/gzip"
    "io"
    "os"
    "testing"
)

func ReadGzip(filename string, handle func(io.Reader)) {
    f, err := os.Open(filename)
    if err != nil {
        return
    }
    defer f.Close()

    z, err := gzip.NewReader(f)
    if err != nil {
        return
    }
    defer z.Close()

    handle(z)
}

func WriteGzip(filename string, handle func(io.Writer)) {
    f, err := os.Create(filename)
    if err != nil {
        return
    }
    defer f.Close()

    z := gzip.NewWriter(f)
    defer z.Close()

    handle(z)
}

func TestIndex(t *testing.T) {
    index := NewIndex()

    rirs := []string{
        "lacnic|VE|ipv4|190.168.0.0|32768|20070329|assigned",
        "lacnic|VE|ipv4|190.168.128.0|16384|20070419|assigned",
        "lacnic|VE|ipv4|190.168.192.0|16384|20070503|assigned",
    }

    for _, rir := range rirs {
        r, n := ParseRIR(rir)
        index.Put(r, n)
    }

    if n := index.Find("190.168.2.1"); n != "VE" {
        t.Error("not expected", n)
    }
}

func TestIndexLoad(t *testing.T) {
    var index *Index
    ReadGzip("geoip.index.json.gz", func(rd io.Reader) {
        index, _ = NewIndexFromJson(rd)
    })

    if n := index.Find("153.16.0.1"); n != "US" {
        t.Error("not expected", n)
    }

    if n := index.Find("153.37.0.1"); n != "CN" {
        t.Error("not expected", n)
    }

    if n := index.Find("122.250.0.1"); n != "JP" {
        t.Error("not expected", n)
    }
}
