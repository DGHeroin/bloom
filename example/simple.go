package main

import (
    "fmt"
    "github.com/DGHeroin/bloom"
)

func main() {
    e := bloom.New("my-data")
    r, err := e.Bucket("user-device")
    if err != nil {
        panic(err)
    }
    defer e.CloseAll()

    var result []string
    for i := 0; i < 10*1000*1000; i++ {
        result = append(result, fmt.Sprint(i))
    }
    r.AddStrings(result...)

    r.AddString("k1")
    r.AddString("k2")
    r.AddString("k3")
    r.AddString("k4")

    r.RemoveString("k3")

    fmt.Printf(`
            k1:%v
            k2:%v
            k3:%v
            k4:%v
`,
        r.Exist("k1"),
        r.Exist("k2"),
        r.Exist("k3"),
        r.Exist("k4"),
    )
    e.CloseBucket("user-device")
}
