package draken

import (
   "fmt"
   "os"
   "testing"
   "time"
)

func TestEntitlement(t *testing.T) {
   var (
      auth auth_login
      err error
   )
   auth.data, err = os.ReadFile("login.json")
   if err != nil {
      t.Fatal(err)
   }
   auth.unmarshal()
   for _, id := range custom_ids {
      movie, err := new_movie(id)
      if err != nil {
         t.Fatal(err)
      }
      title, err := auth.entitlement(movie)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", title)
      time.Sleep(time.Second)
   }
}
