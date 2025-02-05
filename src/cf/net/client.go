package net

import (
	"cf/terminal"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"
)

const (
	PRIVATE_DATA_PLACEHOLDER = "[PRIVATE DATA HIDDEN]"
)

func newHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyFromEnvironment,
	}
	return &http.Client{
		Transport:     tr,
		CheckRedirect: PrepareRedirect,
	}
}

func PrepareRedirect(req *http.Request, via []*http.Request) error {
	if len(via) > 1 {
		return errors.New("stopped after 1 redirect")
	}

	prevReq := via[len(via)-1]

	req.Header.Set("Authorization", prevReq.Header.Get("Authorization"))

	if traceEnabled() {
		dumpRequest(req)
	}

	return nil
}

func Sanitize(input string) (sanitized string) {
	re := regexp.MustCompile(`(?m)^Authorization: .*`)
	sanitized = re.ReplaceAllString(input, "Authorization: "+PRIVATE_DATA_PLACEHOLDER)
	re = regexp.MustCompile(`password=[^&]*&`)
	sanitized = re.ReplaceAllString(sanitized, "password="+PRIVATE_DATA_PLACEHOLDER+"&")
	re = regexp.MustCompile(`"access_token":"[^"]*"`)
	sanitized = re.ReplaceAllString(sanitized, `"access_token":"`+PRIVATE_DATA_PLACEHOLDER+`"`)
	re = regexp.MustCompile(`"refresh_token":"[^"]*"`)
	sanitized = re.ReplaceAllString(sanitized, `"refresh_token":"`+PRIVATE_DATA_PLACEHOLDER+`"`)
	return
}

func doRequest(request *http.Request) (response *http.Response, err error) {
	httpClient := newHttpClient()

	if traceEnabled() {
		dumpRequest(request)
	}

	response, err = httpClient.Do(request)

	if err != nil {
		return
	}

	if traceEnabled() {
		dumpedResponse, err := httputil.DumpResponse(response, true)
		if err != nil {
			fmt.Println("Error dumping response")
		} else {
			fmt.Printf("\n%s\n%s\n", terminal.HeaderColor("RESPONSE:"), Sanitize(string(dumpedResponse)))
		}
	}

	return
}

func traceEnabled() bool {
	traceEnv := strings.ToLower(os.Getenv("CF_TRACE"))
	return traceEnv == "true" || traceEnv == "yes"
}

func dumpRequest(req *http.Request) {
	dumpedRequest, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println("Error dumping request")
	} else {
		fmt.Printf("\n%s\n%s\n", terminal.HeaderColor("REQUEST:"), Sanitize(string(dumpedRequest)))
	}
}
