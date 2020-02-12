package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	secretToken string
	bearerToken *BearerToken
}

func NewCatwalkClient() *Service {
	//log.WithField("secretToken", secretToken).Debugln("starting catwalk client")
	return &Service{secretToken: secretToken, bearerToken: &BearerToken{}}
}

func (s *Service) GetModel(modelName string) (string, error) {
	return s.GetModels([]string{modelName})
}

func (s *Service) GetModels(models []string) (string, error) {
	token, err := s.GetBearerToken()
	if err != nil {
		return "", err
	}

	req, err := NewCatwalkTimeseriesRequest(models)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to get model data")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "cant parse body of timeseries")
	}

	return string(body), nil
}

func (s *Service) UpdateBearerToken() error {
	req, err := http.NewRequest("GET", authEndpoint, nil)
	if err != nil {
		return errors.Wrap(err, "unable to parse url")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.secretToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "unable to update bearer token")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "cant parse body of bearer token requset")
	}
	bt := &BearerToken{}
	err = json.Unmarshal(body, bt)
	if err != nil {
		log.WithError(err).Debug(string(body))
		log.WithFields(log.Fields{
			"request": req,
			"resp":    resp,
		}).Debug("request object")
		return errors.Wrap(err, "failed to unmarshal bearer token")
	}
	s.bearerToken = bt
	return nil
}

func (s *Service) GetBearerToken() (string, error) {
	if !s.bearerToken.Valid() {
		if err := s.UpdateBearerToken(); err != nil {
			return "", err
		}
	}
	return s.bearerToken.AccessToken, nil
}

type BearerToken struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

func (b *BearerToken) Valid() bool {
	if b.AccessToken == "" {
		return false
	}
	if time.Unix(b.ExpiresAt, 0).Before(time.Now()) {
		return false
	}
	return true
}

type CatwalkTimeseriesPayload struct {
	Grain string   `json:"grain"`
	Model []string `json:"model"`
	Span  string   `json:"span"`
	Id    int      `json:"id"`
	Tz    string   `json:"tz"`
	Start string   `json:"start"`
}

func NewCatwalkTimeseriesRequest(model []string) (*http.Request, error) {
	payload := &CatwalkTimeseriesPayload{
		Grain: "aggregate",
		Model: model,
		Span:  "alltime",
		Id:    100949655599,
		Tz:    time.UTC.String(),
		Start: "2020-12-20T00:00:00Z",
	}
	url := fmt.Sprintf("%s?id=%d&span=%s&tz=%s&grain=%s", tsEndpoint, payload.Id, payload.Span, payload.Tz, payload.Grain)
	for _, v := range model {
		url = url + "&model=" + v
	}
	log.WithField("ReqUri", url).Debugln("created request")
	return http.NewRequest("GET", url, nil)
}
