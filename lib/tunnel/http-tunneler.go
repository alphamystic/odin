package tunnel

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func RunHttpTunnel() {
	proxyURL, err := url.Parse("http://example.com:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	resp, err := client.Get("http://www.google.com/")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}
