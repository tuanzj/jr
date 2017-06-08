package jr

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type GAEURLFetch struct {
	ctx *Context
}

func (uf *GAEURLFetch) getCLient() *http.Client {
	client := &http.Client{
		Transport: &urlfetch.Transport{
			Context: uf.ctx.ctx,
			AllowInvalidServerCertificate: appengine.IsDevAppServer(),
		},
	}
	return client
}

func (uf *GAEURLFetch) encode(parameters map[string]string) string {
	if len(parameters) == 0 {
		return ""
	}
	params := url.Values{}
	for k, v := range parameters {
		params.Add(k, v)
	}
	return params.Encode()
}

func (uf *GAEURLFetch) Do(method, urlString, query string, headers ...map[string]string) (string, error) {
	data := ""
	if method == http.MethodGet {
		if query != "" {
			urlString += "?" + query
		}
	} else if method == http.MethodPost {
		data = query
	}

	reader := strings.NewReader(data)

	request, err := http.NewRequest(method, urlString, reader)
	if err != nil {
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
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", errors.New("")
	}

	bodyResp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(bodyResp[:]), nil

}

func (uf *GAEURLFetch) Get(urlString string, query map[string]string, headers ...map[string]string) (string, error) {
	return uf.Do(http.MethodGet, urlString, uf.encode(query), headers...)
}

func (uf *GAEURLFetch) Post(urlString string, data map[string]string, headers ...map[string]string) (string, error) {
	return uf.Do(http.MethodPost, urlString, uf.encode(data), headers...)
}

func (uf *GAEURLFetch) PostJSON(urlString string, data interface{}, headers ...map[string]string) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	addHeader := make(map[string]string, 0)
	addHeader["Cache-Control"] = "max-age=0, must-revalidate"
	addHeader["Content-Type"] = "application/json"
	headers = append(headers, addHeader)
	return uf.Do(http.MethodPost, urlString, string(jsonData), headers...)
}

func (uf *GAEURLFetch) PostXML(urlString string, data interface{}, headers ...map[string]string) (string, error) {
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

func (uf *GAEURLFetch) URLString(urlString string, parameters map[string]string) string {
	params := url.Values{}
	for k, v := range parameters {
		params.Add(k, v)
	}
	urlString += "?" + params.Encode()
	return urlString
}
