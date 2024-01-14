package httpclient

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
)

// CustomHTTPResponse represents the structure of the custom HTTP response.
type CustomHTTPResponse struct {
    FullHeader   http.Header `json:"full_header"`
    RequestedURL string `json:"requested_url"`
    ForwardedURL string `json:"forwarded_url"`
    HTMLBody     string `json:"html_body"`
    HTTPCode     int    `json:"http_code"`
}

// FetchURLWithCustomResponse sends an HTTP GET request to the target URL with the specified headers
// and returns a CustomHTTPResponse struct.
func FetchURLWithCustomResponse(targetURL, hostHeader string, otherHeaders ...map[string]string) (CustomHTTPResponse, error) {
    // Create an HTTP client with custom settings
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Disable SSL certificate verification
    }
    client := &http.Client{Transport: tr}

    // Create an HTTP request with the target URL
    req, err := http.NewRequest("GET", targetURL, nil)
    if err != nil {
        return CustomHTTPResponse{}, err
    }

    // Set the 'Host' header
    req.Header.Set("Host", hostHeader)
    if len(otherHeaders) > 0 {
   // Add additional headers if provided
    for k, v := range otherHeaders[0] {
        req.Header.Set(k, v)
    }
}

    // Send the HTTP request
    resp, err := client.Do(req)
    if err != nil {
        return CustomHTTPResponse{}, err
    }
    defer resp.Body.Close()

    // Read the response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return CustomHTTPResponse{}, err
    }

    // Create a CustomHTTPResponse struct and populate its fields
    // headers,err:=json.Marshal(resp.Header)
    // if err != nil {
    // }
    customResponse := CustomHTTPResponse{
        FullHeader:   resp.Header,
        RequestedURL: targetURL,
        ForwardedURL: resp.Request.URL.String(),
        HTMLBody:     string(body),
        HTTPCode:     resp.StatusCode,
    }

    return customResponse, nil
}

// ToJSON returns a JSON representation of the CustomHTTPResponse struct.
func (c CustomHTTPResponse) ToJSON() (string) {
    jsonResponse, err := json.Marshal(c)
    if err != nil {
        return ""
    }
    return string(jsonResponse)
}
