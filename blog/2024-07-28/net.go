package main

import (
   "net"
   "net/http"
   "os"
)

func main() {
   req, err := http.NewRequest("", "http://example.com", nil)
   if err != nil {
      panic(err)
   }
   req.Close = true
   conn, err := net.Dial("tcp", "example.com:http")
   if err != nil {
      panic(err)
   }
   err = req.Write(conn)
   if err != nil {
      panic(err)
   }
   os.Stdout.ReadFrom(conn)
}
