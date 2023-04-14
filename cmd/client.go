package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

const base_url = "https://prod-api.hingeaws.net"

type hingeClient struct {
	client http.Client
	auth   string
}

func (c *hingeClient) getSelf() {
	url := fmt.Sprintf("%s/user/v2", base_url)

	req, err := http.NewRequest("get", url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req = c.addHeaders(req)
	req.Header.Add("Authorization", "Bearer "+c.auth)

	res, err := c.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func (c *hingeClient) getRecs() (string, error) {

	url := fmt.Sprintf("%s/rec", base_url)
	method := "POST"

	payload := strings.NewReader(`{
  		"excludedFeedIds": [],
  		"genderId": 0,
  		"genderPrefId": 1,
  		"latitude": 41.9038795,
  		"longitude": 12.4520834
	}`)
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", err
	}
	req = c.addHeaders(req)
	req.Header.Add("Authorization", "Bearer "+c.auth)

	res, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *hingeClient) getUsers(ids string) (string, error) {

	url := fmt.Sprintf("%s/user/v2/public?ids=%s", base_url, ids)
	method := "GET"

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return "", err
	}
	req = c.addHeaders(req)
	req.Header.Add("Authorization", "Bearer "+c.auth)

	res, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	//fmt.Println(string(body))
	return string(body), nil
}

func (c *hingeClient) doRating(rating string, token string, sess string, subject string) {
	rid := uuid.New().String()
	url := fmt.Sprintf("%s/rate/v1/initiate", base_url)
	method := "POST"
	layout := "2006-01-02T15:04:05.123Z"
	time := time.Now()
	timestr := time.Format(layout)
	payloadstr := fmt.Sprintf(`{"created":"%s","hasPairing":false,"initiatedWith":"standard","origin":"compatibles","rating":"%s","ratingId":"%s","ratingToken":"%s","sessionId":"%s","subjectId":"%s"}`, timestr, rating, rid, token, sess, subject)
	payload := strings.NewReader(payloadstr)
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req = c.addHeaders(req)
	req.Header.Add("Authorization", "Bearer "+c.auth)

	res, err := c.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Printf("Rating says: %s with args %s", string(body), payloadstr)
	fmt.Printf("Rating says: %s", string(body))
}

func (c *hingeClient) addHeaders(req *http.Request) *http.Request {
	req.Header.Add("x-device-model", "sdk_gphone64_x86_64")
	req.Header.Add("x-app-version", "9.18.0")
	req.Header.Add("User-Agent", "okhttp/4.10.0")
	req.Header.Add("x-os-version", "13")
	req.Header.Add("x-install-id", "4d2620b8-6be2-4c23-ba55-13ccae1bc910")
	req.Header.Add("x-device-model-code", "sdk_gphone64_x86_64")
	req.Header.Add("x-device-platform", "android")
	req.Header.Add("Content-Type", "application/json")
	return req
}

func (c *hingeClient) promptIdToText(id string) (string, error) {
	switch id {
	case "63052f597621ac000162dcb1":
		return "We'll instantly hit it off if", nil
	case "630530154bb20e0001f76454":
		return "Ask me anything about", nil
	case "6305309f64696b0001cdc562":
		return "Give me your honest opinion about", nil
	case "630526c7662dd8000165ea95":
		return "Instead of grabbing drinks, let's", nil
	case "6304feb964696b0001cdc561":
		return "Choose out first date", nil
	case "630535ab64696b0001cdc563":
		return "Pick the best one", nil
	case "63052eee662dd8000165ea99":
		return "Which do we have in common", nil
	case "630530df7621ac000162dcb4":
		return "Which event should we time travel to?", nil
	case "630530767621ac000162dcb3":
		return "Pick the one that's gotta go", nil
	case "63052f977621ac000162dcb2":
		return "A dream home must include", nil
	case "630526e6662dd8000165ea96":
		return "Would you rather", nil
	case "63052dd090123a000193e734":
		return "Let's chat about", nil
	case "63052606662dd8000165ea94":
		return "Let's break the ice by", nil
	case "6305334b90123a000193e735":
		return "The best spot in town for pizza is", nil
	case "63052eb14bb20e0001f76453":
		return "Pick our first getaway", nil
	case "6305272f2fd2b50001f2c8e5":
		return "Pick the most underrated", nil
	case "63052ff52fd2b50001f2c8e6":
		return "If we won the lottery, let's spend it on", nil
	case "63052e923afd1500016c990f":
		return "Guess my secret talent", nil
	case "6304ff373afd1500016c990e":
		return "Two truths and a lie", nil
	case "63052e55662dd8000165ea97":
		return "Which is worth splurging on", nil
	case "63052ecd662dd8000165ea98":
		return "Like for tips about", nil
	default:
		return "", errors.New("prompt id not matched")
	}
}
