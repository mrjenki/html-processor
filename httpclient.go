package httpclient

import (
    "context"
    "crypto/tls"
    "encoding/json"
    "net/http"
    "net/url"
    "golang.org/x/net/idna"
    "io/ioutil"
    "fmt"
)

// CustomHTTPResponse represents the JSON response structure.
type CustomHTTPResponse struct {
    FullHeader   http.Header `json:"full_header"`
    RequestedURL string      `json:"requested_url"`
    ForwardedURL string      `json:"forwarded_url"`
    HTMLBody     string      `json:"html_body"`
    HTTPCode     int         `json:"http_code"`
}

// CustomLookupIP is a custom function for DNS resolution with the specified IP address.
func CustomLookupIP(hostname, ipAddress string) (net.IP, error) {
    if ipAddress == "" {
        return nil, fmt.Errorf("IP address is required")
    }

    punyHostname, err := idna.ToASCII(hostname)
    if err != nil {
        return nil, err
    }

    if punyHostname == ipAddress {
        return net.ParseIP(ipAddress), nil
    }

    return nil, fmt.Errorf("Invalid IP address for hostname")
}

// GetHTMLWithCustomOptions sends an HTTP GET request with custom options.
func GetHTMLWithCustomOptions(targetURL, ipAddress, bearerToken string) ([]byte, error) {
    // Create a custom DNS resolver.
    customResolver := &net.Resolver{
        PreferGo: true,
    }

    // Resolve the IP address for the specified hostname.
    ip, err := CustomLookupIP(targetURL, ipAddress)
    if err != nil {
        return nil, err
    }

    // Configure a custom HTTP transport.
    customTransport := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
            return net.Dial(network, ip.String()+":80")
        },
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
