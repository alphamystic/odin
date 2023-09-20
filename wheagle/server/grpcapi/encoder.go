package grpcapi

import (
  "bytes"
  "encoding/gob"
)
func WorkEncode(cmd *grpcapi.Command) ([]byte, error) {
  var buf bytes.Buffer
  enc := gob.NewEncoder(&buf)
  err := enc.Encode(cmd)
  if err != nil {
    return nil, fmt.Errorf("Error encoding command: %v",err)
  }
  return buf.Bytes(), nil
}


func WorkDecode(body []byte) (*grpcapi.Command,error) {
  var receivedCmd grpcapi.Command
	dec := gob.NewDecoder(bytes.NewReader(body))
	if err := dec.Decode(&receivedCmd); err != nil {
		return nil, fmt.Errorf("Failed to decode received data: %s", err)
	}
	return &receivedCmd, nil
}
