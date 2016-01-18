package nsot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/hashicorp/go-cleanhttp"
)

// Client provides a client to the  nsot API
type Client struct {
	// User Secret
	Secret string

	// User Email
	Email string

	// URL to the DO API to use
	URL string

	// HttpClient is the client to use. Default will be
	// used if not provided.
	Http *http.Client
}

func NewClient(email string, secret string, url string) (*Client, error) {
	if secret == "" {
		secret = os.Getenv("NSOT_SECRET")
	}

	if email == "" {
		email = os.Getenv("NSOT_EMAIL")
	}

	if url == "" {
		url = os.Getenv("NSOT_URL")
	}

	client := Client{
		Secret: secret,
		Email:  email,
		URL:    url,
		Http:   cleanhttp.DefaultClient(),
	}
	return &client, nil
}

// gets the Auth token for the user with the given secret
func getAuthToken(email string, secret string, url string) (string, error) {
	authPostData := fmt.Sprintf("{\"email\": \"%s\" , \"secret_key\": \"%s\"}", email, secret)
	authUrl := fmt.Sprintf("%s/authenticate/", url)
	authJsonParamStr := []byte(authPostData)
	response, err := http.Post(authUrl, "application/json", bytes.NewBuffer(authJsonParamStr))
	if response.StatusCode != 200 {
		return "", fmt.Errorf("Not expected return code, %s", response.Status)
	}
	if err != nil {
		return "", fmt.Errorf("API error: %s", err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", fmt.Errorf("Error reading response for auth token, %s", err)
		}
		byt := []byte(string(contents))
		var dat map[string]interface{}
		if err := json.Unmarshal(byt, &dat); err != nil {
			return "", fmt.Errorf("Error parsing Auth token response from API: %s", err)
		}
		data := dat["data"]
		data_map, ok := data.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("Error parsing Auth token response from API")
		}
		return data_map["auth_token"].(string), nil
	}
	return "", fmt.Errorf("API error, no auth token returned")
}

// Creates a new request with the params
func (c *Client) NewRequest(comp string, method string, body string) (*http.Request, error) {
	authToken, err := getAuthToken(c.Email, c.Secret, c.URL)
	if err != nil {
		return nil, fmt.Errorf("Error getting authToken %s", err)
	}
	url := fmt.Sprintf("%s/%s", c.URL, comp)
	bodyByteArr := []byte(body)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyByteArr))
	if err != nil {
		return nil, fmt.Errorf("Error crafting http request %s", err)
	}
	authTokenHeader := fmt.Sprintf("AuthToken %s:%s", c.Email, authToken)
	req.Header.Add("Authorization", authTokenHeader)
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

// decodeBody is used to JSON decode a body
func decodeBody(resp *http.Response, out interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &out); err != nil {
		return err
	}

	return nil
}

// checkResp wraps http.Client.Do() and verifies that the
// request was successful. A non-200 request returns an error
func checkResp(resp *http.Response, err error) (*http.Response, error) {
	// If the err is already there, there was an error higher
	// up the chain, so just return that
	if err != nil {
		return resp, err
	}

	switch i := resp.StatusCode; {
	case i == 200:
		return resp, nil
	case i == 201:
		return resp, nil
	case i == 204:
		return resp, nil
	default:
		return nil, fmt.Errorf("API Error: %s", resp.Status)
	}
}
