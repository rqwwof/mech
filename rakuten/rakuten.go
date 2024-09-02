package rakuten

import (
   "bytes"
   "encoding/json"
   "errors"
   "net/http"
)

func (GizmoMovie) Show() string {
   return ""
}

func (GizmoMovie) Season() int {
   return 0
}

func (GizmoMovie) Episode() int {
   return 0
}

type GizmoMovie struct {
   Data struct {
      Title string
      Year  int
   }
}

func (g *GizmoMovie) Title() string {
   return g.Data.Title
}

func (g *GizmoMovie) Year() int {
   return g.Data.Year
}

type OnDemand struct {
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
func (o *OnDemand) Info() (*StreamInfo, error) {
   body, err := json.Marshal(o)
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://gizmo.rakuten.tv/v3/avod/streamings",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      var b bytes.Buffer
      resp.Write(&b)
      return nil, errors.New(b.String())
   }
   var value struct {
      Data struct {
         StreamInfos []StreamInfo `json:"stream_infos"`
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.Data.StreamInfos[0], nil
}

func (StreamInfo) WrapRequest(b []byte) ([]byte, error) {
   return b, nil
}

// github.com/mitmproxy/mitmproxy/blob/main/mitmproxy/contentviews/protobuf.py
func (StreamInfo) RequestHeader() (http.Header, error) {
   head := http.Header{}
   head.Set("content-type", "application/x-protobuf")
   return head, nil
}

func (StreamInfo) UnwrapResponse(b []byte) ([]byte, error) {
   return b, nil
}

type StreamInfo struct {
   LicenseUrl   string `json:"license_url"`
   Url          string
   VideoQuality string `json:"video_quality"`
}

func (s *StreamInfo) RequestUrl() (string, bool) {
   return s.LicenseUrl, true
}
