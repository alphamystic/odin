package transport
/*
import (
  "io"
  "os"
)

type Transporter struct {}

func (t *Transporter) Read(data []byte) (int,error){
  copy(data,os.Stdout)
  return len(fromhere),nil
}

func (t *Transporter) Write(data []byte) (int,error){
  copy(os.Stdin,data)
  return len(data),nil
}
/*
create buf
 let read read from stdout
 let write write into stdin
 now find a way to send the buffer around and back
*/
