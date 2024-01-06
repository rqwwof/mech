package nbc

import (
   "bytes"
   "encoding/json"
   "errors"
   "net/http"
   "strconv"
   "time"
)

type Metadata struct {
   Air_Date string `json:"airDate"`
   Episode_Number int64 `json:"episodeNumber,string"`
   MPX_Account_ID int64 `json:"mpxAccountId,string"`
   MPX_GUID int64 `json:"mpxGuid,string"`
   Programming_Type string `json:"programmingType"`
   Season_Number int64 `json:"seasonNumber,string"`
   Secondary_Title string `json:"secondaryTitle"`
   Series_Short_Title string `json:"seriesShortTitle"`
}

func New_Metadata(guid int64) (*Metadata, error) {
   body, err := func() ([]byte, error) {
      var p page_request
      p.Variables.Name = strconv.FormatInt(guid, 10)
      p.Query = graphQL_compact(query)
      p.Variables.App = "nbc"
      p.Variables.One_App = true
      p.Variables.Platform = "android"
      p.Variables.Type = "VIDEO"
      return json.MarshalIndent(p, "", " ")
   }()
   if err != nil {
      return nil, err
   }
   res, err := http.Post(
      "https://friendship.nbc.co/v2/graphql", "application/json",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
   var s struct {
      Data struct {
         Bonanza_Page struct {
            Metadata Metadata
         } `json:"bonanzaPage"`
      }
      Errors []struct {
         Message string
      }
   }
   if err := json.NewDecoder(res.Body).Decode(&s); err != nil {
      return nil, err
   }
   if len(s.Errors) >= 1 {
      return nil, errors.New(s.Errors[0].Message)
   }
   return &s.Data.Bonanza_Page.Metadata, nil
}

func (m Metadata) Series() string {
   return m.Series_Short_Title
}

func (m Metadata) Season() (int64, error) {
   return m.Season_Number, nil
}

func (m Metadata) Episode() (int64, error) {
   return m.Episode_Number, nil
}

func (m Metadata) Title() string {
   return m.Secondary_Title
}

func (m Metadata) Date() (time.Time, error) {
   return time.Parse(time.RFC3339, m.Air_Date)
}
