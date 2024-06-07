package rtbf

import (
   "fmt"
   "os"
   "testing"
)

func TestSix(t *testing.T) {
   text, err := os.ReadFile("account.json")
   if err != nil {
      t.Fatal(err)
   }
   var account accounts_login
   err = account.unmarshal(text)
   if err != nil {
      t.Fatal(err)
   }
   token, err := account.token()
   if err != nil {
      t.Fatal(err)
   }
   gigya, err := token.login()
   if err != nil {
      t.Fatal(err)
   }
   title, err := gigya.entitlement()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", title)
}
