package main

import (
   "2a.pages.dev/mech/youtube"
   "2a.pages.dev/rosso/http"
   "flag"
   "strings"
)

type flags struct {
   audio string
   height int
   info bool
   refresh bool
   request int
   video_ID string
}

func main() {
   var f flags
   // a
   flag.Func("a", "address", func(s string) error {
      return youtube.Video_ID(s, &f.video_ID)
   })
   // b
   flag.StringVar(&f.video_ID, "b", "", "video ID")
   // f
   flag.IntVar(&f.height, "f", 1080, "target video height")
   // g
   flag.StringVar(&f.audio, "g", "AUDIO_QUALITY_MEDIUM", "target audio")
   // i
   flag.BoolVar(&f.info, "i", false, "information")
   // log
   flag.IntVar(
      &http.Default_Client.Log_Level, "log",
      http.Default_Client.Log_Level, "log level",
   )
   // refresh
   flag.BoolVar(&f.refresh, "refresh", false, "create OAuth refresh token")
   // r
   {
      var b strings.Builder
      b.WriteString("0: Android\n")
      b.WriteString("1: Android embed\n")
      b.WriteString("2: Android check")
      flag.IntVar(&f.request, "r", 0, b.String())
   }
   flag.Parse()
   if f.refresh {
      err := f.do_refresh()
      if err != nil {
         panic(err)
      }
   } else if f.video_ID != "" {
      err := f.download()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}
