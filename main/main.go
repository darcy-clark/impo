package main

import(
  "fmt"
  "strings"
  "os"
  "time"
  "net/http"
  "encoding/json"
  "io/ioutil"
  "github.com/faiface/beep/wav"
  "github.com/faiface/beep/speaker"
)

func handler(responseWriter http.ResponseWriter, request *http.Request) {
  fmt.Println("Received request")
  if request.Method != "POST" {
    fmt.Println("ERROR: Invalid method")
    return
  }
  hook := Hook{ }
  if err := request.ParseForm(); err != nil {
    fmt.Println(err)
  } else {
    jsonErr := json.NewDecoder(request.Body).Decode(&hook)
    if jsonErr != nil {
      fmt.Println("ERROR: JSON not processed")
      return
    }

    reader := strings.NewReader("{\"text\": \"Hello my friend\"}")
    getRequest, requestErr := http.NewRequest("POST", "https://stream.watsonplatform.net/text-to-speech/api/v1/synthesize?accept=audio/wav&text=Hello%20world&voice=en-US_AllisonVoice", reader)

    if requestErr != nil {
      fmt.Println("ERROR: Request could not be made")
      return
    }
    header := make(http.Header)
    header.Set("Accept", "audio/wav")
    header.Set("Content-Type", "application/json")
    header.Set("Authorization", "Basic NjA2NjA1MjMtNWJhMS00OTM1LWJhYzMtMzBhZWM4MzdiYmRjOlM1Q2FMbFpLT01lcw==")
    getRequest.Header = header
    client := http.Client{ }

    getResponse, responseErr := client.Do(getRequest)
    if responseErr != nil {
      fmt.Println("RESPONSE ERROR: ", responseErr)
      return
    }

    soundBytes, soundErr := ioutil.ReadAll(getResponse.Body)
    if soundErr != nil {
      fmt.Println("SOUND ERROR: ", soundErr)
    }

    audioFile, fileErr := os.Create("/home/darcy/aString.wav")
    if fileErr != nil {
      fmt.Println("FILE ERROR")
      return
    }
    audioFile.Write(soundBytes)

    openFile, openFileErr := os.Open("/home/darcy/aString.wav")
    if openFileErr != nil {
      fmt.Println("FILE OPEN ERROR", openFileErr)
      return
    }
    streamer, format, streamerErr := wav.Decode(openFile)
    if streamerErr != nil {
      fmt.Println("STREAMER ERROR", streamerErr)
      return
    }
    fmt.Println("Before init...")
    speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
    fmt.Println("Before play...")
    speaker.Play(streamer)
    fmt.Println("After play!")
    select{}
    // fmt.Println(hook)
  }
}

func main() {
  http.HandleFunc("/", handler)
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    fmt.Println(err)
  }
}

type Hook struct {
  Actor string
  Repository string
  Commit_Status CommitStatus
}

type CommitStatus struct {
  Name string
  Description string
  State string
  Key string
  Url string
  Type string
  Created_on string
  Updated_on string
  Links Links
}

type Links struct {
  Commit Commit
  Self Self
}

type Commit struct {
  Href string
}

type Self struct {
  Href string
}

