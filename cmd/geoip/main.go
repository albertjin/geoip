package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"

    "github.com/albertjin/geoip"
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

func main() {
    if os.Args[1] == "build" {
        filenames := []string {
            "delegated-apnic-latest",
            "delegated-arin-extended-latest",
            "delegated-ripencc-latest",
            "delegated-lacnic-latest",
            "delegated-afrinic-latest",
        }

        index := geoip.NewIndex()
        for _, filename := range filenames {
            ReadGzip(filename + ".gz", func(rd io.Reader) {
                index.LoadRIR(rd)
            })
        }

        WriteGzip("geoip.json.gz", func(w io.Writer) {
            log(index.ToJson(w, nil))
        })
        return
    }

    if os.Args[1] == "lookup" {
        var index *geoip.Index
        ReadGzip("geoip.json.gz", func(rd io.Reader) {
            index, _ = geoip.NewIndexFromJson(rd)
        })

        log(index.Find(os.Args[2]))
        return
    }
}

func log(a... interface{}) {
    fmt.Println(a...)
}