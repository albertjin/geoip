package main

import (
    "compress/gzip"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "sync"
    "time"

    "github.com/albertjin/geoip"
)

func ReadGzip(filename string, handle func(io.Reader)) (err error) {
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
    return
}

func WriteGzip(filename string, level int, handle func(io.Writer)) (err error) {
    f, err := os.Create(filename)
    if err != nil {
        return
    }
    defer f.Close()

    z, err := gzip.NewWriterLevel(f, level)
    if err != nil {
        return
    }
    defer z.Close()

    handle(z)
    return
}

func main() {
    if os.Args[1] == "build" {
        var sources []*struct {
            Name string `json:"name"`
            Source string `json:"source"`
        }
        data, _ := ioutil.ReadFile("geoip.sources.json")
        json.Unmarshal(data, &sources)

        index := geoip.NewIndex()
        for _, s := range sources {
            ReadGzip(s.Name + ".gz", func(rd io.Reader) {
                index.LoadRIR(rd)
            })
        }

        WriteGzip("geoip.index.json.gz", gzip.BestCompression, func(w io.Writer) {
            log(index.ToJson(w, nil))
        })
        return
    }

    if os.Args[1] == "download" {
        var sources []*struct {
            Name string `json:"name"`
            Source string `json:"source"`
        }
        data, _ := ioutil.ReadFile("geoip.sources.json")
        json.Unmarshal(data, &sources)

        var wait sync.WaitGroup
        for _, s := range sources {
            wait.Add(1); go func(name, source string) {
                defer wait.Done()
                r, _ := http.DefaultClient.Head(source)
                lm := r.Header.Get("Last-Modified")
                lmt, _ := time.Parse(time.RFC1123, lm)

                filename := name + ".gz"
                if st, err := os.Stat(filename); (err != nil) || (st.ModTime().Before(lmt)) {
                    log("[start]", filename, source)
                    r, _ := http.DefaultClient.Get(source)
                    if r.StatusCode == http.StatusOK {
                        WriteGzip(filename, gzip.BestCompression, func(w io.Writer) {
                            buf := make([]byte, 1024)
                            for {
                                if n, err := r.Body.Read(buf); (err != nil) || (n == 0) {
                                    break
                                } else {
                                    w.Write(buf[:n])
                                }
                            }
                        })
                        log("[finished]", filename, source)
                    } else {
                        log("[error]", filename, source)
                    }
                }
            }(s.Name, s.Source)
        }

        wait.Wait()
        return
    }

    if os.Args[1] == "lookup" {
        var index *geoip.Index
        ReadGzip("geoip.index.json.gz", func(rd io.Reader) {
            index, _ = geoip.NewIndexFromJson(rd)
        })

        log(index.Find(os.Args[2]))
        return
    }
}

func log(a... interface{}) {
    fmt.Println(a...)
}