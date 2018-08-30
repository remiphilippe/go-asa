package goasa

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/golang/glog"
)

// ASA struct holding ASA object
type ASA struct {
	// Hostname or IP address
	Hostname           string
	Insecure           bool
	basicAuthorization *basicAuthorization
	debug              bool
}

// ReferenceObject Reference ASA Object
type ReferenceObject struct {
	Kind     string `json:"kind"`
	ObjectID string `json:"objectId,omitempty"`
	Name     string `json:"name,omitempty"`
	SelfLink string `json:"selfLink,omitempty"`
	Value    string `json:"value,omitempty"`
}

type rangeInfo struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

type basicAuthorization struct {
	Username string
	Password string
}

type requestParameters struct {
	// Request for POST / PUT
	ASARequest interface{}
	// URI Query if needed (GET)
	URIQuery map[string]string
	// Paging parameters for GET
	PageStart int
	PageLimit int
}

// NewASA returns a initizialized ASA struct
func NewASA(hostname string, param map[string]string) (*ASA, error) {
	a := new(ASA)
	a.Hostname = hostname
	a.debug = false
	a.Insecure = false

	if _, ok := param["debug"]; ok {
		if param["debug"] == "true" {
			a.debug = true
		}
	}

	if _, ok := param["insecure"]; ok {
		if param["insecure"] == "true" {
			a.Insecure = true
		}
	}

	a.basicAuthorization = new(basicAuthorization)
	if _, ok := param["username"]; ok {
		a.basicAuthorization.Username = param["username"]
	} else {
		if a.debug {
			glog.Errorf("username is mandatory\n")
		}
		return nil, fmt.Errorf("username is mandatory")
	}

	if _, ok := param["password"]; ok {
		a.basicAuthorization.Password = param["password"]
	} else {
		if a.debug {
			glog.Errorf("password is mandatory\n")
		}
		return nil, fmt.Errorf("password is mandatory")
	}

	return a, nil
}

func (a *ASA) request(endpoint, method string, r *requestParameters) (bodyText []byte, err error) {
	var req *http.Request
	var jsonReq []byte
	var body io.Reader

	uri := url.URL{
		Host:   a.Hostname,
		Scheme: "https",
		Path:   apiBasePath + endpoint,
	}

	switch method {
	case apiPOST, apiPUT:
		if r != nil && r.ASARequest != nil {
			jsonReq, err = json.Marshal(r.ASARequest)
			if err != nil {
				glog.Errorf("request - marshall error: %s\n", err)
				return nil, err
			}
			body = bytes.NewBuffer(jsonReq)
		} else {
			body = nil
		}

		req, err = http.NewRequest(method, uri.String(), body)
		if err != nil {
			glog.Errorln(err)
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

	case apiGET, apiDELETE:
		req, err = http.NewRequest(method, uri.String(), nil)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		q := req.URL.Query()
		if r != nil {
			for k, v := range r.URIQuery {
				q.Add(k, v)
			}
		}
		// if method == apiGET {
		// 	if r != nil && r.PageLimit > 0 {
		// 		q.Add("limit", string(r.PageLimit))
		// 	}

		// 	if r != nil && r.PageStart > 0 {
		// 		q.Add("start", string(r.PageStart))
		// 	}
		// }

		req.URL.RawQuery = q.Encode()

	default:
		return nil, fmt.Errorf("Unknown Method %s", method)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: a.Insecure},
	}
	client := &http.Client{Transport: tr}

	req.SetBasicAuth(a.basicAuthorization.Username, a.basicAuthorization.Password)

	resp, err := client.Do(req)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	bodyText, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("request - readall error: %s\n", err)
		return nil, err
	}

	// spew.Dump(string(jsonReq))
	// spew.Dump(string(bodyText))

	if a.debug {
		if string(jsonReq) != "" {
			glog.Infof("request: %s\n", jsonReq)
		}

		glog.Infof("response: %s\n", bodyText)
	}

	glog.Infof("Response: %s\n", strconv.Itoa(resp.StatusCode))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if resp.StatusCode == 401 {
			return nil, fmt.Errorf("Authentication Failed")
		}

		err = parseResponse(bodyText)
		if err != nil {
			if a.debug {
				glog.Errorf("parse response error: %s\n", err)
			}
			return nil, err
		}

		return nil, fmt.Errorf("response code: %d", resp.StatusCode)
	}

	return bodyText, nil
}

// Post POST to ASA API
func (a *ASA) Post(endpoint string, asaReq interface{}) (bodyText []byte, err error) {
	r := requestParameters{
		ASARequest: asaReq,
	}
	return a.request(endpoint, apiPOST, &r)
}

// Put PUT to ASA API
func (a *ASA) Put(endpoint string, asaReq interface{}) (bodyText []byte, err error) {
	r := requestParameters{
		ASARequest: asaReq,
	}
	return a.request(endpoint, apiPUT, &r)
}

// Get GET to ASA API
func (a *ASA) Get(endpoint string, uriQuery map[string]string) (bodyText []byte, err error) {
	r := requestParameters{
		URIQuery: uriQuery,
	}
	return a.request(endpoint, apiGET, &r)
}

// Delete DELETE to ASA API
func (a *ASA) Delete(endpoint string) (err error) {
	_, err = a.request(endpoint, apiDELETE, nil)
	return err
}
