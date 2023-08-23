package lib

type AdminWrapper struct{
  Conn *net.Conn
	Http_Client *http.Client
	Encoder *gob.Encoder
	Decoder *gob.Decoder
	Response *http.Response
}

func NewClient() (*Client, error) {
	client := &http.Client{}
	encoder := gob.NewEncoder(nil)
	decoder := gob.NewDecoder(nil)
	return &AdminWrapper{
    Http_Client: client,
    Encoder: encoder,
    Decoder: decoder
    }, nil
}

func (c *AdminWrapper) SendRequest(url string, requestData interface{}) (*http.Response, error) {
  // Encode request data
	c.Encoder.Reset(c.Http_Client.Transport)
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

func (c *AdminWrapper) DecodeResponse(responseData *interface{}) error {
	// Decode response body
	err := c.Decoder.Decode(responseData)
	if err != nil {
		return fmt.Errorf("failed to decode response data: %v", err)
	}

	return nil
}
