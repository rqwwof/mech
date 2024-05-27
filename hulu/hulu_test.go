package hulu

import (
   "154.pages.dev/widevine"
   "encoding/hex"
   "fmt"
   "os"
   "testing"
)

// hulu.com/watch/023c49bf-6a99-4c67-851c-4c9e7609cc1d
const default_kid = "21b82dc2ebb24d5aa9f8631f04726650"

// hulu.com/watch/023c49bf-6a99-4c67-851c-4c9e7609cc1d
var test_deep = &DeepLink{
   "EAB::023c49bf-6a99-4c67-851c-4c9e7609cc1d::196861183::262714326",
}

func TestLicense(t *testing.T) {
   var auth Authenticate
   err := auth.getenv()
   if err != nil {
      t.Fatal(err)
   }
   play, err := auth.Playlist(test_deep)
   if err != nil {
      t.Fatal(err)
   }
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   private_key, err := os.ReadFile(home + "/widevine/private_key.pem")
   if err != nil {
      t.Fatal(err)
   }
   client_id, err := os.ReadFile(home + "/widevine/client_id.bin")
   if err != nil {
      t.Fatal(err)
   }
   key_id, err := hex.DecodeString(default_kid)
   if err != nil {
      t.Fatal(err)
   }
   var module widevine.CDM
   err = module.New(private_key, client_id, widevine.PSSH(key_id, nil))
   if err != nil {
      t.Fatal(err)
   }
   key, err := module.Key(play, key_id)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%x\n", key)
}
