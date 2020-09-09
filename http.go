package wechat

import (
	"errors"

	"github.com/sohaha/zlsgo/zhttp"
	"github.com/sohaha/zlsgo/zjson"
)

var http = zhttp.New()

func init() {
	http.DisableChunke()
}

func (e *Engine) Http() *zhttp.Engine {
	return http
}

func (e *Engine) HttpAccessTokenGet(url string, v ...interface{}) (*zjson.Res, error) {
	token, err := e.GetAccessToken()
	if err != nil {
		return nil, err
	}
	v = append(v, zhttp.QueryParam{"access_token": token})
	return httpResProcess(http.Get(url, v...))
}

func (e *Engine) HttpAccessTokenPost(url string, v ...interface{}) (*zjson.Res, error) {
	token, err := e.GetAccessToken()
	if err != nil {
		return nil, err
	}
	v = append(v, zhttp.QueryParam{"access_token": token})
	return httpResProcess(http.Post(url, v...))
}

func httpResProcess(r *zhttp.Res, e error) (*zjson.Res, error) {
	if e != nil {
		return nil, e
	}
	if r.StatusCode() != 200 {
		return nil, errors.New(r.Response().Status)
	}
	json := zjson.ParseBytes(r.Bytes())
	return &json, nil
}
