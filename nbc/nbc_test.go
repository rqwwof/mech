package nbc

import (
   "154.pages.dev/widevine"
   "encoding/base64"
   "fmt"
   "os"
   "testing"
   "time"
)

const raw_pssh = "AAAAV3Bzc2gAAAAA7e+LqXnWSs6jyCfc1R0h7QAAADcIARIQBVLkSEJlSk6BsyYAS+R74BoLYnV5ZHJta2V5b3MiEAVS5EhCZUpOgbMmAEvke+AqAkhE"

func Test_License(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   private_key, err := os.ReadFile(home + "/widevine/private_key.pem")
   if err != nil {
      t.Fatal(err)
   }
   client_ID, err := os.ReadFile(home + "/widevine/client_id.bin")
   if err != nil {
      t.Fatal(err)
   }
   pssh, err := base64.StdEncoding.DecodeString(raw_pssh)
   if err != nil {
      t.Fatal(err)
   }
   mod, err := widevine.New_Module(private_key, client_ID, nil, pssh)
   if err != nil {
      t.Fatal(err)
   }
   key, err := mod.Key(Core)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%x\n", key)
}

func Test_On_Demand(t *testing.T) {
   for _, mpx_guid := range mpx_guids {
      meta, err := New_Metadata(mpx_guid)
      if err != nil {
         t.Fatal(err)
      }
      video, err := meta.On_Demand()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", video)
      time.Sleep(time.Second)
   }
}

var mpx_guids = []int64 {
   // episode unlocked
   // nbc.com/saturday-night-live/video/october-21-bad-bunny/9000283422
   9000283422,
   // movie locked
   // nbc.com/john-wick/video/john-wick/3448375
   3448375,
}
