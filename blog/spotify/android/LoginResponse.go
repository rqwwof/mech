package android

import (
   "154.pages.dev/protobuf"
   "bytes"
   "errors"
   "io"
   "net/http"
)

const android_client_id = "9a8d2f0ce77a4e248bb71fefcb557637"

type login_response struct {
   m protobuf.Message
}

func (r *login_response) New(username, password string) error {
   var m protobuf.Message
   m.Add(1, func(m *protobuf.Message) {
      m.AddBytes(1, []byte(android_client_id))
   })
   m.Add(101, func(m *protobuf.Message) {
      m.AddBytes(1, []byte(username))
      m.AddBytes(2, []byte(password))
   })
   req, err := http.NewRequest(
      "POST", "https://login5.spotify.com/v3/login", bytes.NewReader(m.Encode()),
   )
   if err != nil {
      return err
   }
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return errors.New(res.Status)
   }
   data, err := io.ReadAll(res.Body)
   if err != nil {
      return err
   }
   return r.m.Consume(data)
}

func (r login_response) login_context() ([]byte, bool) {
   return r.m.GetBytes(5)
}

// github.com/librespot-org/librespot/blob/dev/protocol/proto/spotify/login5/v3/challenges/hashcash.proto
// github.com/librespot-org/librespot/blob/dev/protocol/proto/spotify/login5/v3/login5.proto
func (r login_response) prefix() ([]byte, bool) {
   if v, ok := r.m.Get(3); ok { // Challenges
      if v, ok := v.Get(1); ok { // Challenge
         if v, ok := v.Get(1); ok { // challenges.HashcashChallenge
            return v.GetBytes(1) // prefix
         }
      }
   }
   return nil, false
}
