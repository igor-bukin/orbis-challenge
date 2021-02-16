package http

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/orbis-challenge/src/config"
	"github.com/sirupsen/logrus"
)

// NewClient http client constructor
func NewClient() *http.Client {
	timeout, err := time.ParseDuration(config.Config.HTTPTimeout)
	if err != nil {
		logrus.Error("parse timeout config", "error", err)
		return nil
	}

	return NewClientWithTimeout(timeout)
}

func NewClientWithTimeout(timeout time.Duration) *http.Client {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		MaxIdleConns:          1024,
		MaxIdleConnsPerHost:   1024,
		IdleConnTimeout:       150 * time.Second,
		TLSHandshakeTimeout:   150 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         0x0304,
		},
	}

	return &http.Client{Timeout: timeout, Transport: tr}
}

// CloseResponseBody closes response's body
func CloseResponseBody(resp *http.Response) {
	_, err := io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		logrus.Error("discard response's body", "error", err)
	}

	resp.Body.Close()
}
