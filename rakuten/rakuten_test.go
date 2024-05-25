package rakuten

import (
   "154.pages.dev/widevine"
   "encoding/hex"
   "fmt"
   "os"
   "testing"
)

func (m movie_test) license() ([]byte, error) {
   home, err := os.UserHomeDir()
   if err != nil {
      return nil, err
   }
   private_key, err := os.ReadFile(home + "/widevine/private_key.pem")
   if err != nil {
      return nil, err
   }
   client_id, err := os.ReadFile(home + "/widevine/client_id.bin")
   if err != nil {
      return nil, err
   }
   key_id, err := hex.DecodeString(m.key_id)
   if err != nil {
      return nil, err
   }
   var module widevine.CDM
   err = module.New(private_key, client_id, widevine.PSSH(
      key_id, []byte(m.content_id),
   ))
   if err != nil {
      return nil, err
   }
   var web WebAddress
   web.Set(m.url)
   info, err := web.HD().Info()
   if err != nil {
      return nil, err
   }
   return module.Key(info, key_id)
}

func TestLicenseSe(t *testing.T) {
   key, err := tests["se"].license()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%x\n", key)
}

func TestFr(t *testing.T) {
   var web WebAddress
   web.Set(tests["fr"].url)
   stream, err := web.FHD().Info()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", stream)
}

func TestLicenseFr(t *testing.T) {
   key, err := tests["fr"].license()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%x\n", key)
}

func TestSe(t *testing.T) {
   var web WebAddress
   web.Set(tests["se"].url)
   stream, err := web.FHD().Info()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", stream)
}
