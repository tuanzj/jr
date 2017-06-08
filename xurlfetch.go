package jr

import (
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type XURLFetch struct {
	ctx *Context
}

func (uf *XURLFetch) getCLient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: uf.ctx.config.DevENV},
		},
	}
	return client
}

func (uf *XURLFetch) encode(parameters map[string]string) string {
	if len(parameters) == 0 {
		return ""
	}
	params := url.Values{}
	for k, v := range parameters {
		params.Add(k, v)
	}
	return params.Encode()
}

func (uf *XURLFetch) Do(method, urlString, query string, headers ...map[string]string) (string, error) {
	data := ""
	if method == http.MethodGet {
		if query != "" {
			urlString += "?" + query
		}
	} else if method == http.MethodPost {
		data = query
	}

	// uf.ctx.ILog("urlString", urlString)
	// uf.ctx.ILog("headers", headers)
	reader := strings.NewReader(data)

	request, err := http.NewRequest(method, urlString, reader)
	if err != nil {
		// uf.ctx.ILog("err1 "+urlString, err)
		return "", err
	}

	for i := 0; i < len(headers); i++ {
		for k, v := range headers[i] {
			request.Header.Add(k, v)
		}
	}

	client := uf.getCLient()

	response, err := client.Do(request)
	if err != nil {
		// uf.ctx.ILog("err2 "+urlString, err)
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		// uf.ctx.ILog("err3 "+urlString, response)
		return "", errors.New("Status Code Error: " + response.Status)
	}

	bodyResp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// uf.ctx.ILog("err4 "+urlString, err)
		return "", err
	}

	return string(bodyResp[:]), nil

}

func (uf *XURLFetch) Get(urlString string, query map[string]string, headers ...map[string]string) (string, error) {
	return uf.Do(http.MethodGet, urlString, uf.encode(query), headers...)
}

func (uf *XURLFetch) Post(urlString string, data map[string]string, headers ...map[string]string) (string, error) {
	return uf.Do(http.MethodPost, urlString, uf.encode(data), headers...)
}

func (uf *XURLFetch) PostJSON(urlString string, data interface{}, headers ...map[string]string) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	addHeader := make(map[string]string)
	addHeader["Cache-Control"] = "max-age=0, must-revalidate"
	addHeader["Content-Type"] = "application/json"
	headers = append(headers, addHeader)
	// uf.ctx.ILog("JSONPOST", string(jsonData))
	return uf.Do(http.MethodPost, urlString, string(jsonData), headers...)
}

func (uf *XURLFetch) PostXML(urlString string, data interface{}, headers ...map[string]string) (string, error) {
	xmlData, err := xml.Marshal(data)
	if err != nil {
		return "", err
	}
	addHeader := make(map[string]string, 0)
	addHeader["Cache-Control"] = "max-age=0, must-revalidate"
	addHeader["Content-Type"] = "application/xml"
	headers = append(headers, addHeader)

	return uf.Do(http.MethodPost, urlString, string(xmlData), headers...)
}

func (uf *XURLFetch) URLString(urlString string, parameters map[string]string) string {
	params := url.Values{}
	for k, v := range parameters {
		params.Add(k, v)
	}
	urlString += "?" + params.Encode()
	return urlString
}
