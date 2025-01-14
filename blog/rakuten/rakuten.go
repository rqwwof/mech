package rakuten

import (
   "bytes"
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

type on_demand struct {
   AudioLanguage            string `json:"audio_language"`
   AudioQuality             string `json:"audio_quality"`
   ClassificationId         int    `json:"classification_id"`
   ContentId                string `json:"content_id"`
   ContentType              string `json:"content_type"`
   DeviceIdentifier         string `json:"device_identifier"`
   DeviceSerial             string `json:"device_serial"`
   DeviceStreamVideoQuality string `json:"device_stream_video_quality"`
   Player                   string `json:"player"`
   SubtitleLanguage         string `json:"subtitle_language"`
   VideoType                string `json:"video_type"`
}

// geo block
func (o on_demand) streamings() ([]stream_info, error) {
   o.AudioQuality = "2.0"
   o.ContentType = "movies"
   o.DeviceIdentifier = "atvui40"
   o.DeviceStreamVideoQuality = "FHD"
   o.Player = "atvui40:DASH-CENC:WVM"
   o.SubtitleLanguage = "MIS"
   o.VideoType = "stream"
   o.DeviceSerial = "not implemented"
   o.AudioLanguage = "ENG"
   data, err := json.Marshal(o)
   if err != nil {
      return nil, err
   }
   resp, err := http.Post(
      "https://gizmo.rakuten.tv/v3/avod/streamings",
      "application/json", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   var value struct {
      Data struct {
         StreamInfos []stream_info `json:"stream_infos"`
      }
      Errors []struct {
         Message string
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   if err := value.Errors; len(err) >= 1 {
      return nil, errors.New(err[0].Message)
   }
   return value.Data.StreamInfos, nil
}

type stream_info struct {
   LicenseUrl   string `json:"license_url"`
   Url          string
   VideoQuality string `json:"video_quality"`
}

type address struct {
   market_code string
   season_id   string
   content_id  string
}

func (a *address) Set(data string) error {
   data = strings.TrimPrefix(data, "https://")
   data = strings.TrimPrefix(data, "www.")
   data = strings.TrimPrefix(data, "rakuten.tv")
   data = strings.TrimPrefix(data, "/")
   var found bool
   a.market_code, data, found = strings.Cut(data, "/")
   if !found {
      return errors.New("market code not found")
   }
   data, a.content_id, found = strings.Cut(data, "movies/")
   if !found {
      data = strings.TrimPrefix(data, "player/episodes/stream/")
      a.season_id, a.content_id, found = strings.Cut(data, "/")
      if !found {
         return errors.New("episode not found")
      }
   }
   return nil
}

func (a *address) classification_id() (string, error) {
   var v int
   switch a.market_code {
   case "cz":
      v = 272
   case "dk":
      v = 283
   case "fi":
      v = 284
   case "fr":
      v = 23
   case "ie":
      v = 41
   case "it":
      v = 36
   case "nl":
      v = 323
   case "no":
      v = 286
   case "pt":
      v = 64
   case "se":
      v = 282
   case "ua":
      v = 276
   case "uk":
      v = 18
   default:
      return "", errors.New(a.market_code)
   }
   return strconv.Itoa(v), nil
}

type hello_movie struct {
   ViewOptions  struct {
      Private struct {
         Streams []struct {
            AudioLanguages []struct {
               Id string
            } `json:"audio_languages"`
         }
      }
   } `json:"view_options"`
   Title        string
   Year int
}

type hello_episode struct {
   ViewOptions  struct {
      Private struct {
         Streams []struct {
            AudioLanguages []struct {
               Id string
            } `json:"audio_languages"`
         }
      }
   } `json:"view_options"`
   TvShowTitle  string `json:"tv_show_title"`
   SeasonNumber int `json:"season_number"`
   Number       int
   Title        string
}

func (a *address) movie() (*hello_movie, error) {
   classification, err := a.classification_id()
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/v3/movies/" + a.content_id
   req.URL.RawQuery = url.Values{
      "classification_id": {classification},
      "device_identifier": {"atvui40"},
      "market_code":       {a.market_code},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      var b strings.Builder
      resp.Write(&b)
      return nil, errors.New(b.String())
   }
   var value struct {
      Data hello_movie
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.Data, nil
}

func (a *address) season() ([]hello_episode, error) {
   classification, err := a.classification_id()
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/v3/seasons/" + a.season_id
   req.URL.RawQuery = url.Values{
      "classification_id": {classification},
      "device_identifier": {"atvui40"},
      "market_code":       {a.market_code},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      var b strings.Builder
      resp.Write(&b)
      return nil, errors.New(b.String())
   }
   var value struct {
      Data struct {
         Episodes []hello_episode
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return value.Data.Episodes, nil
}
