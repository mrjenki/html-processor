package httpclient

import (
    "encoding/json"
    "net/http"
)

// CustomHTTPResponse represents the JSON response structure.
type CustomHTTPResponse struct {
    FullHeader   http.Header `json:"full_header"`
    RequestedURL string      `json:"requested_url"`
    ForwardedURL string      `json:"forwarded_url"`
    HTMLBody     string      `json:"html_body"`
    HTTPCode     int         `json:"http_code"`
}

// GetHTMLWithCustomOptions sends an HTTP GET request with custom options.
func GetHTMLWithCustomOptions(targetURL, ipAddress, bearerToken string) ([]byte, error) {
  
    // Create a CustomHTTPResponse struct.
    customResponse := CustomHTTPResponse{
        FullHeader:   "a",
        RequestedURL: "b",
        ForwardedURL: "",
        HTMLBody:     "c",
        HTTPCode:     1,
    }

   
    // Convert the response struct to JSON.
    jsonResponse, err := json.Marshal(customResponse)
    if err != nil {
        return nil, err
    }

    return jsonResponse, nil
}
