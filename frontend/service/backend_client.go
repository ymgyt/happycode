package service

import (
	"encoding/binary"
	"net/http"
	"net/url"
)

type BackendClient struct {
	*http.Client
	endpoint *url.URL
}

func NewBackendClient(endpoint *url.URL) *BackendClient {
	return &BackendClient{
		Client:   &http.Client{},
		endpoint: endpoint,
	}
}

func (bc *BackendClient) Get(path string) (*http.Response, error) {
	u := *bc.endpoint
	u.Path = path
	return bc.Client.Get(u.String())

}

func (bc *BackendClient) GetWebSocketPort() int {
	resp, err := bc.Get("/config/server.websocket.port")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var port uint16
	err = binary.Read(resp.Body, binary.LittleEndian, &port)
	if err != nil {
		panic(err)
	}
	return int(port)
}
