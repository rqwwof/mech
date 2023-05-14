package amc

import (
   "2a.pages.dev/rosso/http"
   "encoding/json"
   "errors"
   "net/url"
   "os"
   "strings"
)

func Unauth() (*Auth, error) {
   req := http.Post()
   req.Header = http.Header{
      "X-Amcn-Device-ID": {"-"},
      "X-Amcn-Language": {"en"},
      "X-Amcn-Network": {"amcplus"},
      "X-Amcn-Platform": {"web"},
      "X-Amcn-Tenant": {"amcn"},
   }
   req.URL.Host = "gw.cds.amcn.com"
   req.URL.Path = "/auth-orchestration-id/api/v1/unauth"
   req.URL.Scheme = "https"
   res, err := http.Default_Client.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   auth := new(Auth)
   if err := json.NewDecoder(res.Body).Decode(auth); err != nil {
      return nil, err
   }
   return auth, nil
}

func (a *Auth) Login(email, password string) error {
   body, err := json.Marshal(map[string]string{
      "email": email,
      "password": password,
   })
   if err != nil {
      return err
   }
   req := http.Post()
   req.Body_Bytes(body)
   req.Header = http.Header{
      "Authorization": {"Bearer " + a.Data.Access_Token},
      "Content-Type": {"application/json"},
      "X-Amcn-Device-Ad-ID": {"-"},
      "X-Amcn-Device-ID": {"-"},
      "X-Amcn-Language": {"en"},
      "X-Amcn-Network": {"amcplus"},
      "X-Amcn-Platform": {"web"},
      "X-Amcn-Service-Group-ID": {"10"},
      "X-Amcn-Service-ID": {"amcplus"},
      "X-Amcn-Tenant": {"amcn"},
      "X-Ccpa-Do-Not-Sell": {"doNotPassData"},
   }
   req.URL.Host = "gw.cds.amcn.com"
   req.URL.Path = "/auth-orchestration-id/api/v1/login"
   req.URL.Scheme = "https"
   res, err := http.Default_Client.Do(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   return json.NewDecoder(res.Body).Decode(a)
}

// This accepts full URL or path only.
func (a Auth) Content(ref string) (*Content, error) {
   // If trial is active you must add `/watch` here. If trial has expired, you
   // will get `.data.type` of `redirect`. You can remove the `/watch` to
   // resolve this, but the resultant response will still be missing
   // `video-player-ap`.
   url_path := func(r *http.Request) error {
      p, err := url.Parse(ref)
      if err != nil {
         return err
      }
      if strings.HasPrefix(p.Path, "/movies/") {
         r.URL.Path += "/watch"
      }
      r.URL.Path += p.Path
      return nil
   }
   req := http.Get(&url.URL{
      Scheme: "https",
      Host: "gw.cds.amcn.com",
      Path: "/content-compiler-cr/api/v1/content/amcn/amcplus/path",
   })
   err := url_path(req)
   if err != nil {
      return nil, err
   }
   // If you request once with headers, you can request again without any
   // headers for 10 minutes, but then headers are required again
   req.Header = http.Header{
      "Authorization": {"Bearer " + a.Data.Access_Token},
      "X-Amcn-Network": {"amcplus"},
      "X-Amcn-Tenant": {"amcn"},
   }
   res, err := http.Default_Client.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   con := new(Content)
   if err := json.NewDecoder(res.Body).Decode(con); err != nil {
      return nil, err
   }
   return con, nil
}

func (a *Auth) Refresh() error {
   req := http.Post(&url.URL{
      Scheme: "https",
      Host: "gw.cds.amcn.com",
      Path: "/auth-orchestration-id/api/v1/refresh",
   })
   req.Header.Set("Authorization", "Bearer " + a.Data.Refresh_Token)
   res, err := http.Default_Client.Do(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   return json.NewDecoder(res.Body).Decode(a)
}

func (a Auth) Write_File(name string) error {
   data, err := json.MarshalIndent(a, "", " ")
   if err != nil {
      return err
   }
   return os.WriteFile(name, data, 0666)
}

func (a Auth) Playback(ref string) (*Playback, error) {
   path_body := func(r *http.Request) error {
      _, nID, found := strings.Cut(ref, "--")
      if !found {
         return errors.New("nid not found")
      }
      r.URL.Path += nID
      var p playback_request
      p.Ad_Tags.Mode = "on-demand"
      p.Ad_Tags.URL = "-"
      b, err := json.MarshalIndent(p, "", " ")
      if err != nil {
         return err
      }
      r.Body_Bytes(b)
      return nil
   }
   req := http.Post(&url.URL{
      Scheme: "https",
      Host: "gw.cds.amcn.com",
      Path: "/playback-id/api/v1/playback/",
   })
   err := path_body(req)
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Authorization": {"Bearer " + a.Data.Access_Token},
      "Content-Type": {"application/json"},
      "X-Amcn-Device-Ad-ID": {"-"},
      "X-Amcn-Language": {"en"},
      "X-Amcn-Network": {"amcplus"},
      "X-Amcn-Platform": {"web"},
      "X-Amcn-Service-ID": {"amcplus"},
      "X-Amcn-Tenant": {"amcn"},
      "X-Ccpa-Do-Not-Sell": {"doNotPassData"},
   }
   res, err := http.Default_Client.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   var play Playback
   if err := json.NewDecoder(res.Body).Decode(&play.body); err != nil {
      return nil, err
   }
   play.head = res.Header
   return &play, nil
}

type Auth struct {
   Data struct {
      Access_Token string
      Refresh_Token string
   }
}

func Read_Auth(name string) (*Auth, error) {
   data, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   a := new(Auth)
   if err := json.Unmarshal(data, a); err != nil {
      return nil, err
   }
   return a, nil
}

