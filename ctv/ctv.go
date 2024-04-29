package ctv

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "strconv"
   "strings"
   "time"
)

// wikipedia.org/wiki/Geo-blocking
func (a axis_content) manifest(m *media_content) (string, error) {
   address := func() string {
      b := []byte("https://capi.9c9media.com/destinations/")
      b = append(b, a.AxisPlaybackLanguages[0].DestinationCode...)
      b = append(b, "/platforms/desktop/playback/contents/"...)
      b = strconv.AppendInt(b, a.AxisId, 10)
      b = append(b, "/contentPackages/"...)
      b = strconv.AppendInt(b, m.ContentPackages[0].ID, 10)
      b = append(b, "/manifest.mpd?action=reference"...)
      return string(b)
   }()
   res, err := http.Get(address)
   if err != nil {
      return "", err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      var b strings.Builder
      res.Write(&b)
      return "", errors.New(b.String())
   }
   text, err := io.ReadAll(res.Body)
   if err != nil {
      return "", err
   }
   return string(text), nil
}

type poster struct{}

func (poster) RequestHeader() (http.Header, error) {
   return http.Header{}, nil
}

func (poster) RequestBody(b []byte) ([]byte, error) {
   return b, nil
}

func (poster) ResponseBody(b []byte) ([]byte, error) {
   return b, nil
}

func (poster) RequestUrl() (string, bool) {
   return "https://license.9c9media.ca/widevine", true
}
type media_content struct {
   BroadcastDate date
   ContentPackages []struct {
      ID int64
   }
   Episode int
   Media struct {
      Name string
      Type string
   }
   Name string
   Season struct {
      Number int
   }
}

func (d date) MarshalText() ([]byte, error) {
   return d.T.AppendFormat(nil, time.DateOnly), nil
}

func (d *date) UnmarshalText(text []byte) error {
   var err error
   d.T, err = time.Parse(time.DateOnly, string(text))
   if err != nil {
      return err
   }
   return nil
}

func (a axis_content) media() (*media_content, error) {
   address := func() string {
      b := []byte("https://capi.9c9media.com/destinations/")
      b = append(b, a.AxisPlaybackLanguages[0].DestinationCode...)
      b = append(b, "/platforms/desktop/contents/"...)
      b = strconv.AppendInt(b, a.AxisId, 10)
      b = append(b, "?$include=[ContentPackages,Media,Season]"...)
      return string(b)
   }()
   res, err := http.Get(address)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   media := new(media_content)
   err = json.NewDecoder(res.Body).Decode(media)
   if err != nil {
      return nil, err
   }
   return media, nil
}

type namer struct {
   m *media_content
}

type date struct {
   T time.Time
}

func (n namer) Episode() int {
   return n.m.Episode
}

func (n namer) Season() int {
   return n.m.Season.Number
}

func (n namer) Show() string {
   if v := n.m.Media; v.Type == "series" {
      return v.Name
   }
   return ""
}

func (n namer) Title() string {
   if n.m.Media.Type == "movie" {
      return n.m.Name[:len(n.m.Name)-len(" (2024)")]
   }
   return n.m.Name
}

func (n namer) Year() int {
   return n.m.BroadcastDate.T.Year()
}
type axis_content struct {
   AxisId int64
   AxisPlaybackLanguages []struct {
      DestinationCode string
   }
}

const query_axis = `
query axisContent($id: ID!) {
   axisContent(id: $id) {
      axisId
      axisPlaybackLanguages {
         ... on AxisPlayback {
            destinationCode
         }
      }
   }
}
`

func (r resolve_path) axis() (*axis_content, error) {
   body, err := func() ([]byte, error) {
      var s struct {
         OperationName string `json:"operationName"`
         Query string `json:"query"`
         Variables struct {
            ID string `json:"id"`
         } `json:"variables"`
      }
      s.OperationName = "axisContent"
      s.Query = query_axis
      s.Variables.ID = r.id()
      return json.Marshal(s)
   }()
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://www.ctv.ca/space-graphql/apq/graphql",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   // you need this for the first request, then can omit
   req.Header.Set("graphql-client-platform", "entpay_web")
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   var s struct {
      Data struct {
         AxisContent axis_content
      }
   }
   err = json.NewDecoder(res.Body).Decode(&s)
   if err != nil {
      return nil, err
   }
   return &s.Data.AxisContent, nil
}
