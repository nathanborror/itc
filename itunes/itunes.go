package itunes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const itunesHost = "https://itunesconnect.apple.com/WebObjects/iTunesConnect.woa/ra"

// Client contains credentials to make iTunes Connect requests.
type Client struct {
	ServiceKey    string
	ACN01         string
	MyAccountInfo string
	ITCtx         string
	Session       *Session
}

// NewClient returns an iTunes Connect client.
func NewClient(appleID string, password string) (*Client, error) {
	c := Client{}
	if err := c.serviceConfig(); err != nil {
		return nil, err
	}
	if err := c.signin(appleID, password); err != nil {
		return nil, err
	}
	if err := c.session(); err != nil {
		return nil, err
	}
	return &c, nil
}

type AppSummary struct {
	AdamID   string `json:"adamId"`
	Name     string `json:"name"`
	BundleID string `json:"bundleId"`
}

type Details struct {
	Summaries []AppSummary `json:"summaries"`
}

type DetailsResponse struct {
	Data Details `json:"data"`
}

// Details returns account details.
func (c *Client) Details() (*DetailsResponse, error) {
	req, err := c.NewRequest("GET", itunesHost+"/apps/manageyourapps/summary/v2", nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var details DetailsResponse
	if err := json.Unmarshal(body, &details); err != nil {
		return nil, err
	}
	return &details, nil
}

// Convenience

// NewRequest returns a valid iTunes Connect request with the necessary cookies
// and headers set in advance.
func (c *Client) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	req.AddCookie(&http.Cookie{Name: "acn01", Value: c.ACN01})
	req.AddCookie(&http.Cookie{Name: "myacinfo", Value: c.MyAccountInfo})
	req.AddCookie(&http.Cookie{Name: "itctx", Value: c.ITCtx})
	return req, nil
}

// Paging contains the parameters for paging through results.
type Paging struct {
	Limit int
	Sort  string
	Order string
}

// Encode returns a raw encoded URL query with paging information.
func (p *Paging) Encode(url *url.URL) string {
	if p == nil {
		p = &Paging{Limit: 50, Sort: "email", Order: "asc"}
	}
	q := url.Query()
	q.Add("limit", fmt.Sprintf("%d", p.Limit))
	q.Add("sort", p.Sort)
	q.Add("order", p.Order)
	return q.Encode()
}

// ErrorForServiceErrors returns an error based on the decoded iTunes Connect
// Service Error returned for some requests.
func (c *Client) ErrorForServiceErrors(in []byte) error {
	var errs serviceErrors
	if err := json.Unmarshal(in, &errs); err != nil {
		return err
	}
	var messages []string
	for _, err := range errs.Errors {
		messages = append(messages, fmt.Sprintf("%s (%s)", err.Message, err.Code))
	}
	return fmt.Errorf("iTunes Connect Service Error: %s", strings.Join(messages, "; "))
}

// Private

func (c *Client) serviceConfig() error {
	resp, err := http.Get("https://olympus.itunes.apple.com/v1/app/config?hostname=itunesconnect.apple.com")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var payload serviceConfigResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return err
	}
	c.ServiceKey = payload.AuthServiceKey
	return nil
}

func (c *Client) signin(appleID string, password string) error {
	payload := struct {
		AccountName  string `json:"accountName"`
		Password     string `json:"password"`
		IsRemembered bool   `json:"rememberMe"`
	}{
		appleID,
		password,
		true,
	}
	data, err := json.Marshal(&payload)
	if err != nil {
		return err
	}
	host := "https://idmsa.apple.com/appleauth/auth/signin"
	req, err := http.NewRequest("POST", host, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/javascript")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("X-Apple-Widget-Key", c.ServiceKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "acn01" {
			c.ACN01 = cookie.Value
		}
		if cookie.Name == "myacinfo" {
			c.MyAccountInfo = cookie.Value
		}
	}
	return nil
}

// Session contains information about your current session along with all
// the providers you may belong to.
type Session struct {
	User struct {
		FullName string `json:"fullName"`
		Email    string `json:"emailAddress"`
		PRSID    string `json:"prsId"`
	} `json:"user"`
	Provider           Provider   `json:"provider"`
	AvailableProviders []Provider `json:"availableProviders"`
}

// Provider contains provider information.
type Provider struct {
	ProviderID   int      `json:"providerId"`
	Name         string   `json:"name"`
	ContentTypes []string `json:"contentTypes"`
}

func (c *Client) session() error {
	host := "https://olympus.itunes.apple.com/v1/session"
	req, err := http.NewRequest("GET", host, nil)
	req.AddCookie(&http.Cookie{Name: "acn01", Value: c.ACN01})
	req.AddCookie(&http.Cookie{Name: "myacinfo", Value: c.MyAccountInfo})
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "itctx" {
			c.ITCtx = cookie.Value
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, &c.Session)
}

type serviceConfigResponse struct {
	AuthServiceURL string `json:"authServiceUrl"`
	AuthServiceKey string `json:"authServiceKey"`
}

type serviceErrors struct {
	Errors []serviceError `json:"serviceErrors"`
}

type serviceError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
