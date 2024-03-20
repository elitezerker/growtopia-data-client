package proxy

import (
	"golang.org/x/net/proxy"
	"log"
	"net/http"
	"net/url"
)

func GetPassword(u *url.URL) string {
	if password, isSet := u.User.Password(); isSet {
		return password
	}
	return ""
}

func SetUpHttpTransport(proxyRawURL string) http.Transport {
	var httpTransport http.Transport
	if proxyRawURL != "" {
		parsedURL, err := url.Parse(proxyRawURL)
		if err != nil {
			log.Fatalf("Failed to parse proxy URL: %v", err)
		}
		SetupProxy(parsedURL, &httpTransport)
	}
	return httpTransport
}

func SetupProxy(parsedURL *url.URL, httpTransport *http.Transport) {
	var auth *proxy.Auth
	if username := parsedURL.User.Username(); username != "" {
		auth = &proxy.Auth{
			User:     username,
			Password: GetPassword(parsedURL),
		}
	}
	switch parsedURL.Scheme {
	case "http", "https":
		httpTransport.Proxy = http.ProxyURL(parsedURL)
	case "socks5":
		dialer, err := proxy.SOCKS5("tcp", parsedURL.Host, auth, proxy.Direct)
		if err != nil {
			log.Fatalf("Failed to connect to the proxy: %v", err)
		}
		httpTransport.Dial = dialer.Dial
	default:
		log.Fatalf("Unsupported proxy type: %v", parsedURL.Scheme)
	}
}
