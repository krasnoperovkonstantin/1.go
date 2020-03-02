package main

import (
   "fmt"
   "os"
   "sync"
   "bufio"
   "strings"
   "net/http"
   "io/ioutil"
)

func getcontent(url string) string {
   resp, err := http.Get(url)
   if err != nil {
      fmt.Fprintln(os.Stderr, err)
      return ""
   }
   defer resp.Body.Close()
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      fmt.Fprintln(os.Stderr, err)
      return ""
   }
   return string(body)
}

func process(url string, ch_total chan int) {
   count := strings.Count(getcontent(url), "Go")
   fmt.Printf("Count for %s: %d\n", url, count)
   ch_total <- count
}

func main() {
   var wg sync.WaitGroup
   ch := make(chan bool, 5)
   ch_total := make (chan int)
   ch_done := make (chan bool)
   var total int
   go func () {
      for true {
         count, flag := <- ch_total
         if flag {
            total += count
         } else {
            fmt.Printf("Total: %d", total)
            ch_done <- true
            return
         }
      }
   }()
   scanner := bufio.NewScanner(os.Stdin)
   for scanner.Scan() {
      ch <- true
      wg.Add(1)
      go func(url string) {
         defer wg.Done()
         process (url, ch_total)
         <- ch
      }(scanner.Text())
   }
   wg.Wait()
   close(ch_total)
   <- ch_done
}
