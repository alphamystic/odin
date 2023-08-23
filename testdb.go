package main

/*
  * Write Read Delete and Update
  * Updates* Writes,Reads and d data of any format into json or from a given json format
  * A collection is a database name whilst a resource is a table name

  * Compresion
  * Streamer
  * Read write
*/


import (
  "os"
  "fmt"
  "sync"
  "errors"
  "io/ioutil"
  "path/filepath"
  "encoding/json"
  "encoding/gob"
  "odin/lib/utils"
)


type Driver struct{
  mutex sync.Mutex
  mutexes map[string]*sync.Mutex
  Dir string
  Perm os.FileMode//0644
}

func New(dir string,perm os.FileMode)(*Driver,error){
  dir = filepath.Clean(dir)
  driver := Driver{
    Dir:dir,
    mutexes: make(map[string]*sync.Mutex),
    Perm:perm,
  }
  if _,err := os.Stat(dir); err == nil{
    return nil,errors.New("Directory/db already exists")
  }
  return &driver,os.MkdirAll(dir,perm)
}

func (d Driver) Write(collection,resource string,v interface{}) error{
  if collection == "" || len(collection) <= 0{
    return fmt.Errorf("Missing collection, Can't be nil!")
  }
  if resource == "" || len(resource) <= 0{
    return fmt.Errorf("Missing resource, Can't be nil!")
  }
  mutex := d.getOrCreateMutex(collection)
  mutex.Lock()
  defer mutex.Unlock()
  dir := filepath.Join(d.Dir,collection)
  finalPath := filepath.Join(dir,resource +".json")
  tempPath := finalPath  + ".tmp"
  if err := os.MkdirAll(dir,0755); err != nil{
    return fmt.Errorf("Error creating directory: %s",err)
  }
  fmt.Println("trying to create user")
  b,err := json.MarshalIndent(v,"","\t")
  if err != nil{
    return errors.New(fmt.Sprintf("Error marshal indenting: %s\n",err))
  }
  b = append(b,byte('\n'))
  if err := ioutil.WriteFile(tempPath,b,d.Perm); err != nil{
    return errors.New(fmt.Sprintf("Error writing file: %s\n",err))
  }
  fmt.Println("You aint shit")
  return os.Rename(tempPath,finalPath)
}

func (d Driver) Read(collection,resource string,v interface{}) error{
  if collection == "" || len(collection) <= 0{
    return fmt.Errorf("Missing collection, Can't be nil!")
  }
  if resource == "" || len(resource) <= 0{
    return fmt.Errorf("Missing resource, Can't be nil!")
  }
  record := filepath.Join(d.Dir,collection,resource)
  if _,err := stat(record); err != nil{
    return errors.New(fmt.Sprintf("Directory does not exist: %s",err))
  }
  b,err := ioutil.ReadFile(record + ".json")
  if err != nil{
    return errors.New(fmt.Sprintf("Error reading resource %s, ERROR %s\n",record,err))
  }
  return json.Unmarshal(b,&v)
}

func (d Driver) ReadAll(collection string)([]string,error){
  if collection == "" || len(collection) <= 0{
    return nil,fmt.Errorf("Missing collection, Can't be nil!")
  }
  dir := filepath.Join(d.Dir,collection)
  if _,err := stat(dir); err != nil{
    return nil,errors.New(fmt.Sprintf("Database with name %s does not exist: %s\n",collection,err))
  }
  files,_ := ioutil.ReadDir(dir)
  var records []string
  for _,file := range files{
    b,err := ioutil.ReadFile(filepath.Join(dir,file.Name()))
    if err != nil{
      return nil,errors.New(fmt.Sprintf("Error reading table name %s, ERROR: %s\n",file.Name(),err))
    }
    records = append(records,string(b))
  }
  return records,nil
}

func (d Driver) Delete(collection,resource string) error{
  path := filepath.Join(collection,resource)
  mutex := d.getOrCreateMutex(collection)
  mutex.Lock()
  defer mutex.Unlock()
  dir := filepath.Join(d.Dir,path)
  switch fi,err := stat(dir);{
  case fi == nil,err == nil:
    return errors.New(fmt.Sprintf("Unable to find file or directory named: %v\n",path))
  case fi.Mode().IsDir():
    return os.RemoveAll(dir)
  case fi.Mode().IsRegular():
    os.RemoveAll(dir + ".json")
  }
  return nil
}

func stat(path string) (fi os.FileInfo, err error){
  if fi,err = os.Stat(path); os.IsNotExist(err){
    fi,err = os.Stat(path +".json")
  }
  return
}

func (d *Driver) getOrCreateMutex(collection string)*sync.Mutex{
  d.mutex.Lock()
  defer d.mutex.Unlock()
  m,ok := d.mutexes[collection]
  if !ok {
    m = &sync.Mutex{}
    d.mutexes[collection] = m
  }
  return m
}

func main(){
  d,err := New("vulns1",0755)
  if err != nil{
    utils.Logerror(fmt.Errorf("Error creating dir: %v",err))
  }
  type Vulerability struct{
    Name string
    CVE int
  }
  data := Vulerability{
    Name: "SQL",
    CVE: 9,
  }
  err = d.Write("127-0-0-1","SQL",data)
  if err != nil{
    fmt.Println(err)
  }
  data9 := Vulerability{
    Name: "gegq343  q4y3  4y34  q",
    CVE: 9,
  }
  err = d.Write("127-0-0-1","SSRF",data9)
  if err != nil{
    fmt.Println(err)
  }

  data0 := Vulerability{
    Name: "erherhgqehqhe5h",
    CVE: 9,
  }
  err = d.Write("127-0-0-1","RCE",data0)
  if err != nil{
    fmt.Println(err)
  }

  path := "vulns1/127-0-0-1/test.bin"
  err = SaveSessions(path,data)
  if err != nil{
    utils.Logerror(err)
  }
  data2 := Vulerability{
    Name: "fvewet",
    CVE: 3,
  }
  err = SaveSessions(path,data2)
  if err != nil{
    utils.Logerror(err)
  }

  data3 := Vulerability{
    Name: "kguyfyf87t78ogugugguguguggvukguy",
    CVE: 3,
  }
  err = SaveSessions(path,data3)
  if err != nil{
    utils.Logerror(err)
  }
}

func SaveSessions(path string,data interface{}) error {
	file, err := os.OpenFile(path,os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	return encoder.Encode(data)
}
