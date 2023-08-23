package handlers

import(
  "fmt"
  "sync"
  "time"
  "bufio"
  "regexp"
  "strings"
  "strconv"
  "os/exec"
  "odin/lib/db"
  "odin/lib/utils"
)

//create a channel of datas, let the routines write to it
//read from it, n append to data, close when done
//do the same for services
func  NmapScanForOpenPorts(name string,trg *Target) []*Service{
  datas := make(chan *Output)
  svcs := make(chan *Service)
  var services []*Service
  var data []*Output
  //reconData := make(chan ReconData)
  commands := NmapCommandBuilder("all",trg)
  var wg sync.WaitGroup
  wg.Add(len(commands))
  for _,cmd := range commands {
    //wg.Add(1)
   //defer wg.Done()
    go func(cmnd string){
      defer wg.Done()
      fmt.Println("Scanning command ",cmnd)
      output, err := exec.Command("sh", "-c", cmnd).Output()
      if err != nil {
        utils.Warning(fmt.Sprintf("NMAP: Error scanning for open ports: %s",err))
      }
      datum := &Output{
        Command: cmnd,
        Output: string(output),
      }
      //data = append(data,datum)
      datas <- datum
      // get the service from the nmap scan
      scanner := bufio.NewScanner(strings.NewReader(string(output)))
      r := regexp.MustCompile(`(\d+)\/(\w+)\s+(\w+)\s+([\w\/]+)(.*)`)
      for scanner.Scan() {
        line := scanner.Text()
        matches := r.FindStringSubmatch(line)
        if len(matches) == 6 {
          portNumber,_ := strconv.Atoi(matches[1])
          protocol := matches[2]
          serviceState := matches[3]
          serviceName := matches[4]
          version := matches[5]
          var service  Service
          service.ServiceName = serviceName
          service.Port = portNumber
          service.Protocol = protocol
          if serviceState == "open"{
            service.State = true
          } else {
            service.State = false
          }
          service.Version = version
          switch protocol {
          case "http","https","HTTP","HTTPS":
            service.AT = 1
          default:
            service.AT = 2
          }
          //services = append(services,&service)
          svcs <- &service
        }
      }
      fmt.Println("Done with ",cmnd)
    }(cmd)
  }
  go func() {
      wg.Wait()
      close(datas)
      close(svcs)
  }()

  // Collect data from the datas channel
  for datum := range datas {
      data = append(data, datum)
  }

  // Collect services from the svcs channel
  for service := range svcs {
      services = append(services, service)
  }
  fmt.Println("[+] Done scanning for services for ",trg.TargetIp.String())
  if err := SaveServiceDataTODB(name,trg.TargetIp.String(),data); err != nil{
    utils.Logerror(err)
  }
  return services
}
func NmapCommandBuilder(format string,trg *Target) (commands []string){
  var decoys []string
  for _,ip := range trg.Decoys{
    decoys = append(decoys,ip.String())
  }
  newDecoys := strings.Join(decoys,",")
  var cmnd string
  switch format {
    case "simple":
      cmnd = `nmap -Pn -A -sC  `+ trg.TargetIp.String() +` --open`
    case "fe":
      cmnd = `nmap -Pn -A -sC  -f -D ` + newDecoys + " "+ trg.TargetIp.String() +` --open`
    case "fs","fileshares":
      cmnd = `nmap --script=nfs-ls `+ trg.TargetIp.String()
    case "mdns":
      m1 := `nmap -sS -sV -Pn -sUC -f -p5353 `+ trg.TargetIp.String()
      m2 := `nmap -sU --script=dns-service-discovery -p5353 `+ trg.TargetIp.String()
      commands = append(commands,m1)
      commands = append(commands,m2)
    case "dc":
      cmnd = `nmap -p53,88,135,139,389,445,464,593,636,3268,3269,3389 -sC -sT -Pn `+ trg.TargetIp.String()
    case "all":
      c1 := `nmap -Pn -A -sC  `+ trg.TargetIp.String() +` --open`
      c2 := `nmap -Pn -A -sC  -f -D `+ newDecoys + " "+ trg.TargetIp.String() +` --open`
      c3 := `nmap -p53,88,135,139,389,445,464,593,636,3268,3269,3389 -sC -sT -Pn `+ trg.TargetIp.String()
      c4 := `nmap -sS -Pn -A -T4 -sV -p- -f -D `+ newDecoys + " " + trg.TargetIp.String() +` --open`
      commands = append(commands,c1)
      commands = append(commands,c2)
      commands = append(commands,c3)
      commands = append(commands,c4)
      return
    default:
      cmnd =  `nmap -sS -Pn -A -T4 -sV -p- -f -D `+ newDecoys + " " + trg.TargetIp.String() +` --open`
  }
  commands = append(commands,cmnd)
  return
}


//write data to DB for future reference
var SaveServiceDataTODB = func(name,ip string,data ServiceData)error{
  driver,err := db.Old("../../.brain/scans/" + name,0644)
  if err != nil{
    return err
  }
  for _, datum := range data{
    time.Sleep(1 * time.Second)
    str := utils.RandString(5)
    if err := driver.Write("servicedata",ip + str,datum); err != nil{
      utils.Logerror(fmt.Errorf("Error saving service data for %s to db.\nERROR: %v",ip,err))
      continue
    }
  }
  return nil
}
