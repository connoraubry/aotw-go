package server

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

type XMLDoc struct {
	Results struct {
		Albummatches struct {
			Album []struct {
				Text   string `xml:",chardata"`
				Name   string `xml:"name"`
				Artist string `xml:"artist"`
				URL    string `xml:"url"`
				Mbid   string `xml:"mbid"`
			} `xml:"album"`
		} `xml:"albummatches"`
	} `xml:"results"`
}

func dummyOptions() *Options {
	o := Options{
		FormOption{Title: "Test 1", Artist: "Test 2"},
		FormOption{Title: "Test 2", Artist: "Artist 33"},
		FormOption{Title: "Test 3", Artist: "Artist 2"},
	}
	return &o
}

func queryOptions(query string) *Options {
	logrus.WithField("query", query).Debug("Querying options")

	api_key := "1c8265c48e6df7a5f2d0db26a32eeda3"
	base := "http://ws.audioscrobbler.com/2.0/"
	method := "method=album.search"

	urlOpts := fmt.Sprintf("%s&album=%s&api_key=%s&limit=10", method, url.QueryEscape(query), api_key)
	url := fmt.Sprintf("%s?%s", base, urlOpts)

	logrus.WithField("url", url).Debug("Querying url")

	res, err := http.Get(url)
	if err != nil {
		logrus.WithField("url", url).Error("Received error ", err)
	}
	defer res.Body.Close()

	var target XMLDoc
	xml.NewDecoder(res.Body).Decode(&target)

	return target.ToOptions()
}

func (self *XMLDoc) ToOptions() *Options {
	o := Options{}
	for _, match := range self.Results.Albummatches.Album {
		newOpt := FormOption{
			Title:  match.Name,
			Artist: match.Artist,
		}
		o = append(o, newOpt)
	}

	return &o
}

func ParseSubmission(submission string) (AlbumInfo, error) {
	logrus.WithField("submission", submission).Debug("ParseSubmission")
	aInfo := AlbumInfo{}

	idx := strings.LastIndex(submission, "(")
	if idx == -1 {
		return aInfo, fmt.Errorf("Error finding '(' in str %s", submission)
	}

	aInfo.Title = strings.TrimRight(submission[:idx], " ")
	aInfo.Artist = strings.TrimRight(submission[idx+1:], ")")

	return aInfo, nil
}

func GetIPHash(r *http.Request) string {
	logrus.Debug("Generating IP Hash")
	var ipAddr string

	if ips, ok := r.Header["X-Real-IP"]; ok {
		logrus.Debug("Found X-Real-IP header")
		fmt.Println(ips)
		ipAddr = ips[0]
	} else {
		logrus.Debug("Using 'remoteAddr'")
		ipAddr = r.RemoteAddr
	}

	for key, val := range r.Header {
		fmt.Println(key, val)
	}

	h := sha256.New()
	h.Write([]byte(ipAddr))
	res := h.Sum(nil)

	resString := fmt.Sprintf("%x", res)

	logrus.WithField("ip", resString).Debug("Generated IP Hash")
	return resString
}
