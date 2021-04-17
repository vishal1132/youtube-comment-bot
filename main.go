package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// RefreshToken is
type RefreshToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

// RefBody is
type RefBody struct {
	ClientID     string `json:"client_id"`
	RefreshToken string `json:"refresh_token"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

//Payload is the body of the request for comment request
type Payload struct {
	Vsnippet Vsnippet `json:"snippet"`
}

// Snippet is the
type Snippet struct {
	TextOriginal string `json:"textOriginal"`
}

// TopLevelComment is
type TopLevelComment struct {
	Snippet Snippet `json:"snippet"`
}

// Vsnippet is
type Vsnippet struct {
	VideoID         string          `json:"videoId"`
	TopLevelComment TopLevelComment `json:"topLevelComment"`
}

// Request is the
type Request struct {
	VideoID     string `json:"videoId"`
	NumComments int    `json:"num_comments"`
}

type server struct {
	logger zerolog.Logger
}

// var token string = "ya29.a0AfH6SMDTSNgpzRnUhx_GNrSc2WhGGXZBd_SHVSUPFiC592GZEQLRpZsP_5RXdn1qFw3Ik3voPc-qX9L2EkELJmBqPd7aZ9YjJrpc5gphwnVqTL6z4bUxys0ZI3WCxllqrPESL9xUPKNJHTwRIjkP4j3PXMibYa7djj1Kg119ckg"
var token string = "ya29.a0AfH6SMBHaggQpZDxcXxYaZnTGr-5ZM5WPy--Zkp5DQ1NBTg7sv2GqfM3b2Tq5a2qwJtyDYwkxWjEzG_YYR2S9DXrfeapBGqJ55Wd_DLQzxEFB5_JIvxVJooV_UYuHNL-UCQ63Gaom1rtOSeJpCjegJd3hIP_sgDUWtiZhkS592I"

func (s *server) getRandomComment() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1).Intn(650000)
	f, err := os.OpenFile("comments.txt", os.O_RDONLY, 0644)
	lg := s.logger.With().Str("getcomment", "file").Logger()
	if err != nil {
		lg.Error().Err(err).Msg("unable to open commens file")
		return ""
	}
	scanner := bufio.NewScanner(f)
	for i := 0; i < r1; i++ {
		scanner.Scan()
	}
	return scanner.Text()
}

func (s *server) makeComments(videoID string, num int) {
	lg := s.logger.With().Str("make", "comments").Logger()
	for i := 0; i < num; i++ {
		var cmt string = s.getRandomComment()
		if cmt == "" {
			lg.Error().Msg("Unable to get comment")
			return
		}
		vsnippet := Vsnippet{
			VideoID:         videoID,
			TopLevelComment: TopLevelComment{Snippet: Snippet{TextOriginal: cmt}},
		}
		data := Payload{
			// fill struct
			Vsnippet: vsnippet,
		}
		payloadBytes, err := json.Marshal(data)
		if err != nil {
			// handle err
			lg.Error().Err(err).Msg("failed to marshal payload for comment")
			return
		}
		newbody := bytes.NewReader(payloadBytes)

		req, err := http.NewRequest("POST", "https://youtube.googleapis.com/youtube/v3/commentThreads?part=snippet&key=AIzaSyDlt_ceDvJGRj6aVt7urPVHaaurEecQbnE", newbody)
		if err != nil {
			lg.Error().Err(err).Msg("error creating new http request")
		}
		btoken := fmt.Sprintf("Bearer %s", token)
		req.Header.Set("Authorization", btoken)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// handle err
			lg.Error().Err(err).Msg("failed to make http request")
			return
		}
		defer resp.Body.Close()
		res, err := ioutil.ReadAll(resp.Body)
		log.Println(string(res))
	}
}

func (s *server) v1Handler(w http.ResponseWriter, r *http.Request) {
	lg := s.logger.With().Str("", "").Logger()
	// lg.Info().Msg(string(r.))
	log.Println(*r)
	var vid = Request{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &vid)
	if err != nil {
		lg.Error().Err(err).Msg("failed to unmarshal request body")
		return
	}
	go s.makeComments(vid.VideoID, vid.NumComments)
	fmt.Fprintf(w, "Success")

}

func (s *server) refreshToken() {
	lg := s.logger.With().Str("scope", "refreshtoken").Logger()

	// ref, err := json.Marshal(refbody)
	// refreader := bytes.NewReader(ref)
	data := url.Values{}
	data.Set("client_id", "338733535217-pfd8tbmfi7n3sug19uqtbiufbo05pamf.apps.googleusercontent.com")
	data.Set("refresh_token", "1//04AHHfB0UXgxjCgYIARAAGAQSNwF-L9Ir6vNYj2ykpKGe0qQOZWkNdyIReOL44l9rhOpoTVvBn0Mux_-gYcRYPSPH4ivF25XJMQQ")
	data.Set("client_secret", "AnKw8bGnTCCEHbAwLxO_lkC4")
	data.Set("grant_type", "refresh_token")

	r, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", strings.NewReader(data.Encode()))
	if err != nil {
		lg.Error().Err(err).Msg("error creating new http request")
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		lg.Error().Err(err).Msg("failed to marshal body for url encoded")
		return
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		// handle err
		lg.Error().Err(err).Msg("failed to make http request")
		return
	}
	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)
	reftoken := RefreshToken{}
	err = json.Unmarshal(res, &reftoken)
	if err != nil {
		lg.Error().Err(err).Msg("error unmarshaling refresh token response")
	}
	// log.Printf("%+v", reftoken)
	token = reftoken.AccessToken
	log.Println(reftoken.AccessToken)
	time.Sleep(50 * time.Minute)
}

func (s *server) v2Handler(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	vidID := keys["vidID"][0]
	numC := keys["numc"][0]
	num, err := strconv.Atoi(numC)
	if err != nil {
		return
	}
	go s.makeComments(vidID, num)
	fmt.Fprintf(w, "success")
}

func main() {
	s := server{logger: zerolog.New(os.Stdout)}
	go func() {
		for {
			s.refreshToken()
		}
	}()
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/comment", s.v1Handler)
	mux.HandleFunc("/v2/comment", s.v2Handler)
	port := os.Getenv("PORT")
	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)

	// http.ListenAndServe(":8080", mux)
}
