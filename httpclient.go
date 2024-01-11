package httpclient

import (
    "encoding/json"
    "net/http"

    "time"
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
    // Create a custom resolver for the provided IP address.
    customResolver := &net.Resolver{
        PreferGo: true,
        Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
            if address == targetURL {
                return net.Dial("tcp", ipAddress)
            }
            return net.Dial(network, address)
        },
    }

    // Configure a custom HTTP transport.
    customTransport := &http.Transport{
        DialContext:           customResolver.DialContext,
        TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
        DisableKeepAlives:     false,
        MaxIdleConns:          100,
        MaxIdleConnsPerHost:   100,
        IdleConnTimeout:       90 * time.Second,
        TLSHandshakeTimeout:   10 * time.Second,
        ExpectContinueTimeout: 1 * time.Second,
    }

    // Create an HTTP client with the custom transport.
    client := &http.Client{
        Transport: customTransport,
        Timeout:   10 * time.Second, // Set your desired timeout.
    }

    // Create an HTTP request.
    req, err := http.NewRequest("GET", targetURL, nil)
    if err != nil {
        return nil, err
    }

    // Set the Bearer token in the request header if provided.
    if bearerToken != "" {
        req.Header.Set("Authorization", "Bearer "+bearerToken)
    }

    // Perform the HTTP request.
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Read the response body.
    htmlBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    // Create a CustomHTTPResponse struct.
    customResponse := CustomHTTPResponse{
        FullHeader:   resp.Header,
        RequestedURL: targetURL,
        ForwardedURL: "",
        HTMLBody:     string(htmlBody),
        HTTPCode:     resp.StatusCode,
    }

    // Check if there were any redirects and add the forwarded URL if available.
    if len(resp.Request.URL.String()) > 0 && resp.Request.URL.String() != targetURL {
        customResponse.ForwardedURL = resp.Request.URL.String()
    }

    // Convert the response struct to JSON.
    jsonResponse, err := json.Marshal(customResponse)
    if err != nil {
        return nil, err
    }

    return jsonResponse, nil
}
