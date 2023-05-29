package mech

import (
   "2a.pages.dev/mech/widevine"
   "2a.pages.dev/rosso/dash"
   "2a.pages.dev/rosso/http"
   "2a.pages.dev/rosso/mp4"
   "encoding/base64"
   "encoding/xml"
   "fmt"
   "io"
   "net/url"
   "os"
)

type Stream struct {
   Client_ID string
   Info bool
   Namer
   Poster widevine.Poster
   Private_Key string
   base *url.URL
}

func (s *Stream) DASH(ref string) (dash.Represents, error) {
   client := http.Default_Client
   client.CheckRedirect = nil
   res, err := client.Get(ref)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   var pre dash.Presentation
   if err := xml.NewDecoder(res.Body).Decode(&pre); err != nil {
      return nil, err
   }
   s.base = res.Request.URL
   return pre.Represents(), nil
}

func (s Stream) DASH_Get(items dash.Represents, index int) error {
   if s.Info {
      for i, item := range items {
         if i == index {
            fmt.Print("!")
         }
         fmt.Println(item)
      }
      return nil
   }
   item := items[index]
   file_name, err := Name(s)
   if err != nil {
      return err
   }
   file, err := os.Create(file_name + item.Ext())
   if err != nil {
      return err
   }
   defer file.Close()
   client := http.Default_Client
   client.CheckRedirect = nil
   req, err := http.Get_Parse(item.Initialization())
   if err != nil {
      return err
   }
   req.URL = s.base.ResolveReference(req.URL)
   res, err := client.Do(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   media := item.Media()
   pro := http.Progress_Chunks(file, len(media))
   dec := mp4.New_Decrypt(pro)
   var key []byte
   if item.Content_Protection != nil {
      private_key, err := os.ReadFile(s.Private_Key)
      if err != nil {
         return err
      }
      client_ID, err := os.ReadFile(s.Client_ID)
      if err != nil {
         return err
      }
      pssh, err := base64.StdEncoding.DecodeString(item.Widevine().PSSH)
      if err != nil {
         return err
      }
      mod, err := widevine.New_Module(private_key, client_ID, pssh)
      if err != nil {
         return err
      }
      keys, err := mod.Post(s.Poster)
      if err != nil {
         return err
      }
      key = keys.Content().Key
      if err := dec.Init(res.Body); err != nil {
         return err
      }
   } else {
      _, err := io.Copy(pro, res.Body)
      if err != nil {
         return err
      }
   }
   client.Log_Level = 0
   for _, ref := range media {
      req.URL, err = s.base.Parse(ref)
      if err != nil {
         return err
      }
      res, err := client.Do(req)
      if err != nil {
         return err
      }
      pro.Add_Chunk(res.ContentLength)
      if item.Content_Protection != nil {
         err = dec.Segment(res.Body, key)
      } else {
         _, err = io.Copy(pro, res.Body)
      }
      if err != nil {
         return err
      }
      if err := res.Body.Close(); err != nil {
         return err
      }
   }
   return nil
}
