package wechat

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	netHttp "net/http"
	"sort"
	"strings"

	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/zutil"
)

type (
	XMLData map[string]string
	request struct {
		request *netHttp.Request
		rawData []byte
	}
)

func FormatMap2XML(m XMLData) (string, error) {
	buf := zutil.GetBuff()
	defer zutil.PutBuff(buf)
	if _, err := io.WriteString(buf, "<xml>"); err != nil {
		return "", err
	}
	for k, v := range m {
		if _, err := io.WriteString(buf, fmt.Sprintf("<%s>", k)); err != nil {
			return "", err
		}
		if err := xml.EscapeText(buf, zstring.String2Bytes(v)); err != nil {
			return "", err
		}
		if _, err := io.WriteString(buf, fmt.Sprintf("</%s>", k)); err != nil {
			return "", err
		}
	}
	if _, err := io.WriteString(buf, "</xml>"); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ParseXML2Map parse xml to map
func ParseXML2Map(b []byte) (m XMLData, err error) {
	var (
		d     = xml.NewDecoder(bytes.NewReader(b))
		depth = 0
		tk    xml.Token
		key   string
		buf   bytes.Buffer
	)
	m = XMLData{}
	for {
		tk, err = d.Token()
		if err != nil {
			if err == io.EOF {
				err = nil
				return
			}
			return
		}
		switch v := tk.(type) {
		case xml.StartElement:
			depth++
			switch depth {
			case 2:
				key = v.Name.Local
				buf.Reset()
			case 3:
				if err = d.Skip(); err != nil {
					return
				}
				depth--
				key = "" // key == "" indicates that the node with depth==2 has children
			}
		case xml.CharData:
			if depth == 2 && key != "" {
				buf.Write(v)
			}
		case xml.EndElement:
			if depth == 2 && key != "" {
				m[key] = buf.String()
			}
			depth--
		}
	}
}

// SignWithMD5 生成MD5签名
// func SignWithMD5(m XMLData, apikey string) string {
// 	return strings.ToUpper(zstring.Md5(BuildSignStr(m, apikey)))
// }

// SignWithHMacSHA256 生成HMAC-SHA256签名
// func SignWithHMacSHA256(m XMLData, apikey string) string {
// 	signature := utils.HMAC(utils.AlgoSha256, BuildSignStr(m, apikey), apikey)
//
// 	return strings.ToUpper(signature)
// }

// Sign 生成签名
func BuildSignStr(m XMLData) string {
	l := len(m)
	ks := make([]string, 0, l)
	kvs := make([]string, 0, l)
	for k := range m {
		if k == "sign" {
			continue
		}
		ks = append(ks, k)
	}
	sort.Strings(ks)
	for _, k := range ks {
		if v, ok := m[k]; ok && v != "" {
			kvs = append(kvs, fmt.Sprintf("%s=%s", k, v))
		}
	}
	// kvs = append(kvs, fmt.Sprintf("key=%s", apikey))
	return strings.Join(kvs, "&")
}

func CheckResError(v []byte) (*zjson.Res, error) {
	data := zjson.ParseBytes(v)
	code := data.Get("errcode").String()
	if code == "0" || code == "" {
		return &data, nil
	}
	msg := data.Get("errmsg").String()
	return &zjson.Res{}, errors.New(code + ": " + msg)
}

// func replyXmlUnmarshal(plaintext []byte, data *ReplySt) error {
// 	log.Debug(zstring.Bytes2String(plaintext))
// 	err := xml.Unmarshal(plaintext, &data)
// 	// log.Error(err, data)
// 	// if data.Event == "LOCATION" {
// 	// 	if data.LocationX == "" {
// 	// 		data.LocationY = data.latitude
// 	// 		data.LocationX = data.longitude
// 	// 	}
// 	// }
// 	return err
// }

// func (r *request) getQuery(key string) (string, bool) {
// 	if values, ok := r.getQueryArray(key); ok {
// 		return values[0], ok
// 	}
// 	return "", false
// }
//
// func (r *request) getQueryArray(key string) ([]string,
// 	bool) {
// 	if values, ok := r.request.URL.Query()[key]; ok && len(values) > 0 {
// 		return values, true
// 	}
// 	return []string{}, false
// }

// func (r *request) getDataRaw() (result []byte, err error) {
// 	if len(r.rawData) > 0 {
// 		return r.rawData, nil
// 	}
// 	if r.request.Body == nil {
// 		err = errors.New("request.Body is nil")
// 		return
// 	}
// 	var buf bytes.Buffer
// 	tee := io.TeeReader(r.request.Body, &buf)
// 	r.request.Body = tee
// 	result, err = ioutil.ReadAll(buf)
// 	if err == nil {
// 		r.rawData = result
// 	}
// 	return
// }
