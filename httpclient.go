package httpclient

import (
    "context"
    "crypto/tls"
    "encoding/json"
    "net/http"
    "net/url"
    "net"
    "time"
    "golang.org/x/net/idna"
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
        Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
            host, port, err := net.SplitHostPort(address)
            if err != nil {
                return nil, err
            }

            ip, err := CustomLookupIP(host, ipAddress)
            if err != nil {
                return nil, err
            }

            return net.Dial(network, ip.String()+":"+port)
        },
    }

    // Rest of the code remains the same...
    // (Configure customTransport, create HTTP client, perform request, etc.)

    return jsonResponse, nil
}
