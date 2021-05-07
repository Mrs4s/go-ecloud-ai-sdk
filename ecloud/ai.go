package ecloud

import (
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

type (
	AiService struct {
		g           *dataflow.Gout
		opt         *AiOptions
		expireTime  time.Time
		accessToken string
	}

	AiOptions struct {
		AccessKey string
		SecretKey string
		Client    *http.Client
	}

	authMiddleware struct {
		as *AiService
	}
)

func NewAi(opt *AiOptions) (*AiService, error) {
	s := &AiService{
		g:   gout.New(opt.Client),
		opt: opt,
	}
	if err := s.RefreshAccessToken(); err != nil {
		return nil, errors.Wrap(err, "request access token error")
	}
	return s, nil
}

func (as *AiService) RefreshAccessToken() error {
	var rsp gout.H
	if err := as.g.GET("https://smartlib-api-changsha-1.cmecloud.cn:8444/ecloud/ai/oauth/getToken").SetQuery(gout.H{
		"grant_type":    "client_credentials",
		"client_id":     as.opt.AccessKey,
		"client_secret": as.opt.SecretKey,
	}).BindJSON(&rsp).Do(); err != nil {
		return errors.Wrap(err, "send request error")
	}
	if m, ok := rsp["errorCode"]; ok {
		return errors.New(m.(string))
	}
	ei, _ := strconv.ParseInt(rsp["expires_in"].(string), 10, 64)
	as.expireTime = time.Now().Add(time.Second * time.Duration(ei))
	as.accessToken = rsp["access_token"].(string)
	return nil
}

func (as *AiService) auth() *authMiddleware {
	return &authMiddleware{
		as: as,
	}
}

func (a *authMiddleware) ModifyRequest(req *http.Request) error {
	if time.Now().After(a.as.expireTime) {
		if err := a.as.RefreshAccessToken(); err != nil {
			return err
		}
	}
	if req.URL.RawQuery != "" {
		req.URL.RawQuery += "&"
	}
	req.URL.RawQuery += "access_token=" + a.as.accessToken
	return nil
}
