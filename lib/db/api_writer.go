package db


import (
  "fmt"
  "errors"
  "encoding/json"
  "github.com/alphamystic/odin/lib/utils"
)

type ApiWriter struct {
  WriterClient *utils.OdinAPIClient
}

func NewApiWriter(writer *utils.OdinAPIClient) (*ApiWriter,error) {
  if writer == nil{
    return nil, errors.New("No API  Client Present.")
  }
  return &ApiWriter{
    WriterClient : writer,
  },nil
}


func (api_writer *ApiWriter) WriteToAPI(method,directory string, minimalService any) (*utils.APIResponse,error) {
  jsonData, err := json.Marshal(minimalService)
	if err != nil {
		utils.Notice(fmt.Sprintf("Error encoding Service JSON: %s", err))
		return nil,err
	}
  //payload := bytes.NewReader(jsonData)
  api_resp,err := api_writer.WriterClient.DoRequest(method,directory,string(jsonData))
  if err != nil {
    utils.NoticeError(fmt.Sprintf("%s",err))
    return nil,err
  }
  return api_resp,nil
}
