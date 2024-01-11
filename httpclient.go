package httpclient
import (
    "crypto/tls"
    "io/ioutil"
    "net"
    "net/http"
    "time"
)

type ResponseData struct {
    FullHeader   http.Header `json:"full_header"`
    RequestedURL string      `json:"requested_url"`
    ForwardedURL string      `json:"forwarded_url"`
    HTMLBody     string      `json:"html_body"`
    HTTPCode     int         `json:"http_code"`
}

func fetchURL(url, bearerToken string) (*ResponseData, error) {
    // Create HTTP client with custom transport
    client := &http.Client{
        Transport: &http.Transport{
            DialContext: (&net.Dialer{
                Timeout: 5 * time.Second,
            }).DialContext,
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return nil // Follow redirects
        },
    }

    // Create request
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    // Add bearer token
    req.Header.Add("Authorization", "Bearer "+bearerToken)

    // Send request
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Read HTML body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    // Prepare the response data
    responseData := &ResponseData{
        FullHeader:   resp.Header,
        RequestedURL: url,
        ForwardedURL: resp.Request.URL.String(),
        HTMLBody:     string(body),
        HTTPCode:     resp.StatusCode,
    }

    return responseData, nil
}
