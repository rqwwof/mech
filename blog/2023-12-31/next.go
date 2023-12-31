package main

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
   "os"
   "strings"
   "time"
)

func try(name, version string) error {
   var req http.Request
   req.Header = make(http.Header)
   req.Method = "POST"
   req.ProtoMajor = 1
   req.ProtoMinor = 1
   req.URL = new(url.URL)
   req.URL.Host = "www.youtube.com"
   req.URL.Scheme = "https"
   //req.URL.RawQuery = "prettyPrint=false"
   req.URL.RawQuery = "prettyPrint=true"
   req.URL.Path = "/youtubei/v1/next"
   req.Body = func() io.ReadCloser {
      s := fmt.Sprintf(`
      {
         "videoId": "2ZcDwdXEVyI",
         "context": {
            "client": {
               "clientName": %q,
               "clientVersion": %q
            }
         }
      }
      `, name, version)
      return io.NopCloser(strings.NewReader(s))
   }()
   res, err := new(http.Transport).RoundTrip(&req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      fmt.Println(res.Status, name)
      return nil
   }
   body, err := io.ReadAll(res.Body)
   if err != nil {
      return err
   }
   lower := strings.ToLower(string(body))
   if strings.Contains(lower, "in the heat of the night") {
      fmt.Println("pass", name, len(body))
      os.WriteFile(name+".json", body, 0666)
   } else {
      fmt.Println("fail", name, len(body))
   }
   return nil
}

func main() {
   for _, line := range strings.Split(clients, "\n") {
      fields := strings.Split(line, ";")
      err := try(fields[1], fields[2])
      if err != nil {
         panic(err)
      }
      time.Sleep(99*time.Millisecond)
   }
}

const clients = `1;WEB;2.20220918
2;MWEB;2.20220918
3;ANDROID;17.36.4
5;IOS;17.36.4
7;TVHTML5;7.20220918
8;TVLITE;2
10;TVANDROID;1.0
13;XBOXONEGUIDE;1.0
14;ANDROID_CREATOR;22.36.102
15;IOS_CREATOR;22.36.102
16;TVAPPLE;1.0
18;ANDROID_KIDS;7.36.1
19;IOS_KIDS;7.36.1
21;ANDROID_MUSIC;5.26.1
23;ANDROID_TV;2.19.1.303051424
26;IOS_MUSIC;5.26.1
27;MWEB_TIER_2;9.20220918
28;ANDROID_VR;1.37
29;ANDROID_UNPLUGGED;6.36
30;ANDROID_TESTSUITE;1.9
31;WEB_MUSIC_ANALYTICS;0.2
33;IOS_UNPLUGGED;6.36
38;ANDROID_LITE;3.26.1
39;IOS_EMBEDDED_PLAYER;2.4
41;WEB_UNPLUGGED;1.20220918
42;WEB_EXPERIMENTS;1
43;TVHTML5_CAST;1.1
55;ANDROID_EMBEDDED_PLAYER;17.36.4
56;WEB_EMBEDDED_PLAYER;9.20220918
57;TVHTML5_AUDIO;2.0
58;TV_UNPLUGGED_CAST;0.1
59;TVHTML5_KIDS;3.20220918
60;WEB_HEROES;0.1
61;WEB_MUSIC;1.0
62;WEB_CREATOR;1.20220918
63;TV_UNPLUGGED_ANDROID;1.37
64;IOS_LIVE_CREATION_EXTENSION;17.36.4
65;TVHTML5_UNPLUGGED;6.36
66;IOS_MESSAGES_EXTENSION;17.36.4
67;WEB_REMIX;1.20220918
68;IOS_UPTIME;1.0
69;WEB_UNPLUGGED_ONBOARDING;0.1
70;WEB_UNPLUGGED_OPS;0.1
71;WEB_UNPLUGGED_PUBLIC;0.1
72;TVHTML5_VR;0.1
74;ANDROID_TV_KIDS;1.19.1
75;TVHTML5_SIMPLY;1.0
76;WEB_KIDS;2.20220918
77;MUSIC_INTEGRATIONS;0.1
80;TVHTML5_YONGLE;0.1
84;GOOGLE_ASSISTANT;0.1
85;TVHTML5_SIMPLY_EMBEDDED_PLAYER;2.0
87;WEB_INTERNAL_ANALYTICS;0.1
88;WEB_PARENT_TOOLS;1.20220918
89;GOOGLE_MEDIA_ACTIONS;0.1
90;WEB_PHONE_VERIFICATION;1.0.0
92;IOS_PRODUCER;0.1
93;TVHTML5_FOR_KIDS;7.20220918
94;GOOGLE_LIST_RECS;0.1
95;MEDIA_CONNECT_FRONTEND;0.1`
