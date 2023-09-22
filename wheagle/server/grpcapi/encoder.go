package grpcapi

import (
  "fmt"
  "bytes"
  "encoding/gob"
)

func WorkEncode(cmd *Command) ([]byte, error) {
  var buf bytes.Buffer
  enc := gob.NewEncoder(&buf)
  err := enc.Encode(cmd)
  if err != nil {
    return nil, fmt.Errorf("Error encoding command: %v",err)
  }
  return buf.Bytes(), nil
}


func WorkDecode(body []byte) (*Command,error) {
  var receivedCmd Command
	dec := gob.NewDecoder(bytes.NewReader(body))
	if err := dec.Decode(&receivedCmd); err != nil {
		return nil, fmt.Errorf("Failed to decode received data: %s", err)
	}
	return &receivedCmd, nil
}


func OPWorkEncode(cmd *C2Command) ([]byte, error) {
  var buf bytes.Buffer
  enc := gob.NewEncoder(&buf)
  err := enc.Encode(cmd)
  if err != nil {
    return nil, fmt.Errorf("Error encoding command: %v",err)
  }
  return buf.Bytes(), nil
}


func OPWorkDecode(body []byte) (*C2Command,error) {
  var receivedCmd C2Command
	dec := gob.NewDecoder(bytes.NewReader(body))
	if err := dec.Decode(&receivedCmd); err != nil {
		return nil, fmt.Errorf("Failed to decode received data: %s", err)
	}
	return &receivedCmd, nil
}

func EncodeScreenShot(scnrshts *Screenshots) ([]byte,error) {
  var buf bytes.Buffer
  enc := gob.NewEncoder(&buf)
  if err := enc.Encode(scnrshts); err != nil {
    return nil,fmt.Errorf("Error encoding screenshot: %q",err)
  }
  return buf.Bytes(),nil
}


func DecodeScreenshot(body []byte) (*Screenshots,error) {
  var scrnshts Screenshots
  dec := gob.NewDecoder(bytes.NewReader(body))
  if err := dec.Decode(&scrnshts); err != nil {
    return nil,fmt.Errorf("Error decoding screenshot: %q",err)
  }
  return &scrnshts,nil
}
