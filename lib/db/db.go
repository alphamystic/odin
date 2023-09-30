package db

import (
  "os"
  "fmt"
  "sync"
  "errors"
  "io/ioutil"
  "path/filepath"
  "encoding/json"
  "github.com/alphamystic/odin/lib/utils"
)

type Driver struct{
  Mutex sync.Mutex
  Mutexes map[string]*sync.Mutex//locks only the database in question
  Dir string
  Perm os.FileMode
}
type Database struct{
  Name string
  Tables int
}

func New(dir string,perm os.FileMode)(*Driver,error){
  dir = filepath.Clean(dir)
  driver := Driver{
    Dir: dir,
    Mutexes: make(map[string]*sync.Mutex),
    Perm: perm,
  }
  //ignore file info n ensrure DB-dir has been created
  if _,err := os.Stat(dir);err == nil{
    return nil,errors.New("DB with dir/name already exists.")
  }
  return &driver,os.MkdirAll(dir,perm)
}

func Old(dir string,perm os.FileMode)(*Driver,error){
  dir = filepath.Clean(dir)
  driver := Driver{
    Dir: dir,
    Mutexes: make(map[string]*sync.Mutex),
    Perm: perm,
  }
  if _,err := os.Stat(dir);err != nil{
    err = os.MkdirAll(dir,perm)
    if err != nil{
      return nil,fmt.Errorf("Error creating directory. %v",err)
    }
  }
  return &driver,nil
}

func (d *Driver) Write(dbName,tableName string,data interface{}) error{
  if !utils.CheckifStringIsEmpty(dbName) {
    return fmt.Errorf("Error db name can not be empty.")
  }
  if  !utils.CheckifStringIsEmpty(tableName) {
    return fmt.Errorf("Error table name can not be empty.")
  }
  mutex := d.getOrCreateMutex(dbName)
  mutex.Lock()
  defer mutex.Unlock()
  dir := filepath.Join(d.Dir,dbName)
  finalPath := filepath.Join(d.Dir,dir,tableName + ".json")
  tempPath := finalPath + ".tmp"
  if err := os.MkdirAll(dir,0755); err != nil{
    return fmt.Errorf("Error creating dir: %v",err)
  }
  b,err := json.MarshalIndent(data,"","\t")
  if err != nil {
    errors.New(fmt.Sprintf("Error marshal indenting: %s\n",err))
  }
  b = append(b,byte('\n'))
  if err := ioutil.WriteFile(tempPath,b,d.Perm); err != nil{
    return errors.New(fmt.Sprintf("Error writing file: %s\n",err))
  }
  return os.Rename(tempPath,finalPath)
}

func (d *Driver) Read(dbName,tableName string,v interface{}) error {
  if !utils.CheckifStringIsEmpty(dbName) {
    return fmt.Errorf("Error db name can not be empty.")
  }
  if  !utils.CheckifStringIsEmpty(tableName) {
    return fmt.Errorf("Error table name can not be empty.")
  }
  record := filepath.Join(d.Dir,dbName,tableName)
  if _,err := stat(record); err != nil{
    return fmt.Errorf("Directory does not exist: %v",err)
  }
  b,err := ioutil.ReadFile(record + ".json")
  if err != nil{
    return errors.New(fmt.Sprintf("Error reading resource %s, ERROR %s\n",record,err))
  }
  return json.Unmarshal(b,&v)
}

/*func (d *Driver) ReadAll(dbName string) ([]string,error){
  if !utils.CheckifStringIsEmpty(dbName) {
    return nil,fmt.Errorf("Error db name can not be empty.")
  }
  dir := filepath.Join(d.Dir,dbName)
  if _,err := stat(dir); err != nil{
    return nil,errors.New(fmt.Sprintf("Database with name %s does not exist: %s\n",dbName,err))
  }
  files,_ := ioutil.ReadDir(dir)
  var records []string
  var wg sync.WaitGroup
  for _,file := range files{
    wg.Add(1)
    go func(name string){
      b,err := ioutil.ReadFile(filepath.Join(dir,name))
      if err != nil {
        utils.NoticeError(fmt.Sprintf("Error reading table name: %s\nERROR: %s\n",name,err))
      }
      records = append(records,string(b))
    }(file.Name())
  }
  wg.Wait()
  return records,nil
}*/

func (d *Driver) ReadAll(dbName string) ([]string,error){
  if !utils.CheckifStringIsEmpty(dbName) {
    return nil,fmt.Errorf("Error db name can not be empty.")
  }
  dir := filepath.Join(d.Dir,dbName)
  if _,err := stat(dir); err != nil{
    return nil,errors.New(fmt.Sprintf("Database with name %s does not exist: %s\n",dbName,err))
  }
  files,_ := ioutil.ReadDir(dir)
  var records []string
  for _,file := range files{
    b,err := ioutil.ReadFile(filepath.Join(dir,file.Name()))
    if err != nil {
      utils.NoticeError(fmt.Sprintf("Error reading table name: %s\nERROR: %s\n",file.Name(),err))
    }
    records = append(records,string(b))
  }
  return records,nil
}

func (d *Driver) ReadAllWithInterface(dbName string,data interface{})([]interface{},error){
  if !utils.CheckifStringIsEmpty(dbName) {
    return nil,fmt.Errorf("Error db name can not be empty.")
  }
  dir := filepath.Join(d.Dir,dbName)
  if _,err := stat(dir); err != nil{
    return nil,errors.New(fmt.Sprintf("Database with name %s does not exist: %s\n",dbName,err))
  }
  files,_ := ioutil.ReadDir(dir)
  var records []interface{}
  var wg sync.WaitGroup
  for _,file := range files{
    wg.Add(1)
    go func(name string){
      b,err := ioutil.ReadFile(filepath.Join(dir,name))
      if err != nil {
        utils.NoticeError(fmt.Sprintf("Error reading table name: %s\nERROR: %s\n",name,err))
        return
      }
      err = json.Unmarshal(b,&data)
      if err != nil{
        utils.Notice(fmt.Sprintf("Error unmarshaling db data:\n ERROR: %s",err))
        return
      }
      records = append(records,data)
    }(file.Name())
  }
  wg.Wait()
  return records,nil
}


func (d *Driver) Delete(dbName,tableName string) error{
  if !utils.CheckifStringIsEmpty(dbName) {
    return fmt.Errorf("Error db name can not be empty.")
  }
  if  !utils.CheckifStringIsEmpty(tableName) {
    return fmt.Errorf("Error table name can not be empty.")
  }
  path := filepath.Join(dbName,tableName)
  path = filepath.Join(d.Dir,path)
  fmt.Println("Path alone ",path)
  mutex := d.getOrCreateMutex(dbName)
  mutex.Lock()
  defer mutex.Unlock()
  dir := filepath.Join(d.Dir,path)
  fmt.Println("Dir plus path ",dir)
  switch fi,err := stat(dir); {
    case fi == nil,err == nil:
      return fmt.Errorf("Unable to find file or directory named: %v\n. ERROR: %v",path,err)
    case fi.Mode().IsDir():
      return os.RemoveAll(dir)
    case fi.Mode().IsRegular():
      return os.RemoveAll(dir + ".json")
  }
  return nil
}

func stat(path string) (fi os.FileInfo,err error){
  if fi,err = os.Stat(path); os.IsNotExist(err){
    fi,err = os.Stat(path + ".json")
  }
  return
}
func (d *Driver) ExportedStat(path string) (fi os.FileInfo,err error){
  if fi,err = os.Stat(path); os.IsNotExist(err){
    fi,err = os.Stat(path + ".json")
  }
  return
}

func (d *Driver) getOrCreateMutex(dbName string)*sync.Mutex{
  d.Mutex.Lock()
  defer d.Mutex.Unlock()
  m,ok := d.Mutexes[dbName]
  if !ok {
    m = &sync.Mutex{}
    d.Mutexes[dbName] = m
  }
  return m
}
