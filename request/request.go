package request

import (
	"net/http"
	url2 "net/url"
)

func Get(url string, headers map[string]string, proxy []byte) (*http.Response, error) {
	// headers := map[string]string{
	// 	"referer":    "https://www.github.com/",
	// 	"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
	// }
	// proxy := []byte("http://127.0.0.1:7890")
	client := http.DefaultClient
	if proxy != nil {
		proxyUrl, _ := url2.Parse(string(proxy))
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
