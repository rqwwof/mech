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

type flags struct {
   email          string
   entity         kanopy.EntityId
   home           string
   min_width      int64
   password       string
   representation string
   s              internal.Stream
}

func (f *flags) New() error {
   var err error
   f.home, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.home = filepath.ToSlash(f.home)
   f.s.ClientId = f.home + "/widevine/client_id.bin"
   f.s.PrivateKey = f.home + "/widevine/private_key.pem"
   return nil
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.Var(&f.entity, "a", "address")
   flag.StringVar(&f.s.ClientId, "c", f.s.ClientId, "client ID")
   flag.StringVar(&f.email, "e", "", "email")
   flag.StringVar(&f.representation, "i", "", "representation")
   flag.StringVar(&f.s.PrivateKey, "k", f.s.PrivateKey, "private key")
   flag.StringVar(&f.password, "p", "", "password")
   flag.Int64Var(&f.min_width, "m", 1280, "min width")
   flag.Parse()
   text.Transport{}.Set()
   switch {
   case f.password != "":
      err := f.authenticate()
      if err != nil {
         panic(err)
      }
   case f.entity.String() != "":
      err := f.download()
      if err != nil {
         panic(err)
      }
   default:
      flag.Usage()
   }
}
