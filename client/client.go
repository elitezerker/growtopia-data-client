package client

import (
	"github.com/elitezerker/growtopia-data-client/config"
	"github.com/elitezerker/growtopia-data-client/httpclient"
	"github.com/elitezerker/growtopia-data-client/proxy"
	"log"
	"net/http"
	"strings"
)

type Config struct {
	ServerDataUrl  string
	HttpMethod     string
	PayloadContent string
	ProxyRawURL    string
	UserAgent      string
	ContentType    string
}

func NewConfig(serverDataUrl, httpMethod, payloadContent, proxyRawURL, userAgent string, contentType string) Config {
	return Config{
		ServerDataUrl:  serverDataUrl,
		HttpMethod:     httpMethod,
		PayloadContent: payloadContent,
		ProxyRawURL:    proxyRawURL,
		UserAgent:      userAgent,
		ContentType:    contentType,
	}
}

type Client struct {
	Cfg      Config
	Client   http.Client
	HttpClnt httpclient.HttpClient
}

func NewClient(cfg Config) Client {
	httpTransport := proxy.SetUpHttpTransport(cfg.ProxyRawURL)
	client := http.Client{Transport: &httpTransport}
	httpClnt := httpclient.HttpClient{
		Client: &client,
	}
	return Client{
		Cfg:      cfg,
		Client:   client,
		HttpClnt: httpClnt,
	}
}

func (c Client) BuildRequest() (*http.Request, error) {
	payload := strings.NewReader(c.Cfg.PayloadContent)
	request, err := http.NewRequest(c.Cfg.HttpMethod, c.Cfg.ServerDataUrl, payload)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", c.Cfg.UserAgent)
	request.Header.Add("Content-Type", c.Cfg.ContentType)
	return request, nil
}

func (c Client) ExecuteRequest() (config.Data, error) {
	request, err := c.BuildRequest()
	if err != nil {
		log.Fatalf("Error building request: %v", err)
	}
	response, err := c.HttpClnt.FetchResponseBody(request)
	if err != nil {
		return config.Data{}, err
	}
	defer response.Body.Close()
	body := c.HttpClnt.ReadResponseBody(response)
	return config.ParseConfig(body), nil
}
