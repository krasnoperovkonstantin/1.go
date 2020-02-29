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
      fmt.Println(err)
      return ""
   }
   defer resp.Body.Close()
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      fmt.Println(err)
      return ""
   }
   return string(body)
}

func process(url string) int {
   count := strings.Count(getcontent(url), "Go")
   fmt.Printf("Count for %s: %d\n", url, count)
   return count
}

func main() { 
   var wg sync.WaitGroup
   ch := make(chan bool, 5)
   var total int
   scanner := bufio.NewScanner(os.Stdin)
   for scanner.Scan() {
      ch <- true
      wg.Add(1)
      go func(url string) {
         defer wg.Done()
         total += process (url)
         <- ch
      }(scanner.Text())
   }
   wg.Wait()
   fmt.Printf("Total: %d",total)
}
