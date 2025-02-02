package criterion

import (
   "41.neocities.org/widevine"
   "encoding/hex"
   "fmt"
   "os"
   "os/exec"
   "strings"
   "testing"
)

func TestToken(t *testing.T) {
   data, err := exec.Command("password", "criterionchannel.com").Output()
   if err != nil {
      t.Fatal(err)
   }
   username, password, _ := strings.Cut(string(data), ":")
   data, err = new(AuthToken).Marshal(username, password)
   if err != nil {
      t.Fatal(err)
   }
   os.WriteFile("token.txt", data, os.ModePerm)
}

func TestWrap(t *testing.T) {
   data, err := os.ReadFile("token.txt")
   if err != nil {
      t.Fatal(err)
   }
   var token AuthToken
   err = token.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   item, err := token.Video(video_test.slug)
   if err != nil {
      t.Fatal(err)
   }
   files, err := token.Files(item)
   if err != nil {
      t.Fatal(err)
   }
   file, ok := files.Dash()
   if !ok {
      t.Fatal("VideoFiles.Dash")
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
   key_id, err := hex.DecodeString(video_test.key_id)
   if err != nil {
      t.Fatal(err)
   }
   var pssh widevine.PsshData
   pssh.KeyIds = [][]byte{key_id}
   var module widevine.Cdm
   err = module.New(private_key, client_id, pssh.Marshal())
   if err != nil {
      t.Fatal(err)
   }
   data, err = module.RequestBody()
   if err != nil {
      t.Fatal(err)
   }
   _, err = file.Wrap(data)
   if err != nil {
      t.Fatal(err)
   }
}

func TestVideo(t *testing.T) {
   data, err := os.ReadFile("token.txt")
   if err != nil {
      t.Fatal(err)
   }
   var token AuthToken
   err = token.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   item, err := token.Video(video_test.slug)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", item)
}

var video_test = struct{
   key_id string
   slug string
   url string
}{
   key_id: "e4576465a745213f336c1ef1bf5d513e",
   slug: "my-dinner-with-andre",
   url: "criterionchannel.com/videos/my-dinner-with-andre",
}
