package lib
/*
import (
  "net"
  "fmt"
  "net/http"
)
type AdminWrapper struct{
  Conn *net.Conn
	Http_Client *http.Client
	Response *http.Response
}

func NewClient() (*AdminWrapper, error) {
	client := &http.Client{}
	return &AdminWrapper{
    Http_Client: client,
    }, nil
}

func (c *AdminWrapper) SendRequest(url string, requestData interface{}) (*http.Response, error) {
	err := c.Encoder.Encode(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request data: %v", err)
	}
	// Send request to server
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := c.Http_Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	// Save the response for decoding
	c.Response = resp

	return resp, nil
}

// this here is a useless function because I forgot how I was to use it to decode
// the wrapper response and  where to store or return it to.
func (c *AdminWrapper) DecodeResponse(responseData interface{}) error {
	// Decode response body
	return nil
}
*/
