package main

import (
   "154.pages.dev/media/internal"
   "154.pages.dev/media/pluto"
   "errors"
   "fmt"
   "net/http"
)

func (f flags) download() error {
   video, err := f.web.Video(f.forward)
   if err != nil {
      return err
   }
   clip, err := video.Clip()
   if err != nil {
      return err
   }
   source, ok := clip.DASH()
   if !ok {
      return errors.New("pluto.EpisodeClip.DASH")
   }
   var req http.Request
   req.URL, err = source.Parse(f.base)
   if err != nil {
      return err
   }
   media, err := internal.DASH(&req)
   if err != nil {
      return err
   }
   for _, medium := range media {
      if medium.ID == f.representation {
         f.s.Name = pluto.Namer{video}
         f.s.Poster = pluto.Poster{}
         return f.s.Download(medium)
      }
   }
   // 2 MPD all
   for i, medium := range media {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(medium)
   }
   return nil
}
