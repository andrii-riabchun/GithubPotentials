package main

import (
    "fmt"
    "time"
    potentials "github.com/artisresistance/githubpotentials"
)


func main() {
    a := potentials.New("SECRET",100)
    t:= time.Date(2016,7,15,0,0,0,0,time.Local)
    it := a.SearchIterator(t)

    ticker := time.NewTicker(5 * time.Second)
    quit := make(chan struct{})
    go func() {
        for {
        select {
            case <- ticker.C:
                fmt.Println(a.GetAPIRates())
            case <- quit:
                ticker.Stop()
                return
            }
        }
    }()

    b := a.
        SetCriterias(it,t).
        Dump().
        Sort(potentials.CombinedCriteria)
    fmt.Println(len(b))
     
}