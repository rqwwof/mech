package main

import (
   "41.neocities.org/dash"
   "41.neocities.org/media/kanopy"
   "41.neocities.org/media/internal"
   "41.neocities.org/text"
   "flag"
   "fmt"
   "io"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flags) download() error {
   data, err := os.ReadFile(f.home + "/kanopy.txt")
   if err != nil {
      return err
   }
   var auth kanopy.Authenticate
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   deep, err := auth.DeepLink(&f.entity)
   if err != nil {
      return err
   }
   play, err := auth.Playlist(deep)
   if err != nil {
      return err
   }
   resp, err := http.Get(play.StreamUrl)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return err
   }
   var mpd dash.Mpd
   mpd.Unmarshal(data)
   for represent := range mpd.Representation() {
      if *represent.Width < f.min_width {
         if *represent.MimeType == "video/mp4" {
            continue
         }
      }
      switch f.representation {
      case "":
         fmt.Print(&represent, "\n\n")
      case represent.Id:
         f.s.Wrapper = play
         return f.s.Download(&represent)
      }
   }
   return nil
}

func (f *flags) authenticate() error {
   data, err := kanopy.Authenticate{}.Marshal(f.email, f.password)
   if err != nil {
      return err
   }
   return os.WriteFile(f.home+"/kanopy.txt", data, os.ModePerm)
}
