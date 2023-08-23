package utils

/*
  Ported from my old code
      https://github.com/3l0racle/golog
*/
import (
  "io"
  "os"
  "log"
  "io/fs"
  "net/http"
  "path/filepath"
)

type ErrorLogger struct{
  Filename string
  Dir string
  Perm fs.FileMode
  Text interface{}
}

type RequestLogger struct {
  FileName string
  Dir string
  Perm fs.FileMode
  Handle http.Handler
}

func LogErrorToNamedFileInDir(name,dir string,text ...interface{}) {
  dir = filepath.Clean(dir)
  if dir != "" || len(dir) <= 0{
    name = name + ".log"
    name = filepath.Join(dir,name)
    f,err := os.OpenFile(name,os.O_RDWR|os.O_CREATE|os.O_APPEND,0666)
    if err != nil {
      Logerror(err)
    }
    defer f.Close()
    writer := io.MultiWriter(os.Stdout,f)
    log.SetOutput(writer)
    log.Println(text)
  }
}

func Start(l RequestLogger) http.Handler{
  //open file for loging
  l.OpenLogFile()
  //set the flags
  log.SetFlags(log.Ldate|log.Ltime|log.Lshortfile)
  //return the handler
  return l.LogRequest(l.Handle)
}

//opens or creates a file for logging
// To be addedfile permisions
func (l RequestLogger) OpenLogFile(){
  //clean the directory
  dir := filepath.Clean(l.Dir)
  if dir != "" || len(l.Dir) != 0{
    if l.FileName != "" && len(l.FileName) >= 0{
      name := l.FileName + ".log"
      name = filepath.Join(l.Dir,name)
      dataLog,err := os.OpenFile(name,os.O_WRONLY|os.O_CREATE|os.O_APPEND,l.Perm)
      if err != nil{
        log.Fatal("[-] Error logging to file: ",err)
      }
      log.SetOutput(dataLog)
    }
  }
}


func (l RequestLogger)LogRequest(handler http.Handler) http.Handler{
  return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request){
    log.Printf("%s %s %s\n",req.RemoteAddr,req.Method,req.URL)
    handler.ServeHTTP(res,req)
  })
}
