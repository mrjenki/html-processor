package htmlprocessor

import (
	"net/http"
	"os"
    curl "github.com/andelf/go-curl"
)
func getHTML(url string) (map[string]string, error) {
	// init the curl object
	easy := curl.EasyInit()
	defer easy.Cleanup()

	// set options
	easy.Setopt(curl.OPT_URL, url)
	easy.Setopt(curl.OPT_SSL_VERIFYPEER, false)
	easy.Setopt(curl.OPT_FOLLOWLOCATION, true)

	// create a buffer
	html := ""
	header := ""
	// set the callback function
	fooTest := func (buf []byte, userdata interface{}) bool {
		html += string(buf)
		return true
	}
	fooHeader := func (buf []byte, userdata interface{}) bool {
		header += string(buf)
		return true
	}
	easy.Setopt(curl.OPT_WRITEFUNCTION, fooTest)
	easy.Setopt(curl.OPT_HEADERFUNCTION, fooHeader)

	// execute
	if err := easy.Perform(); err != nil {
		return nil, err
	}

	// return the html
	return map[string]string{"url": url, "html": html, "headers": header}, nil
}
