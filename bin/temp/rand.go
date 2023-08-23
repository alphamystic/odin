package main

import(
  "fmt"
  "net"
  //"sync"
  "strings"
  "bufio"
  "regexp"
  "strconv"
  "os/exec"
//  "odin/lib/utils"
)

func main(){
  trg := &Target {
    Host: "127.0.0.1",
    HostIp: net.ParseIP("127.0.0.1"),
    TargetIp: net.ParseIP("127.0.0.1"),
    Decoys: []net.IP{net.ParseIP("4.4.4.4"),net.ParseIP("127.0.0.2")},
    FireWallName: "ufw",
  }
  services,d := NmapScanForOpenPorts("test",trg)
  for _,service := range services{
    fmt.Println("Service...... ")
    fmt.Println(service)
    fmt.Println("")
  }
  fmt.Println("DAta")

  for _,o := range d{
    fmt.Println(o.Command)
    fmt.Println(o.Output)
    fmt.Println("")
  }
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
      cmnd = `nmap -Pn -A -sC  `+ trg.TargetIp.String()
    case "fe":
      cmnd = `nmap -Pn -A -sC  -f -D ` + newDecoys + " "+ trg.TargetIp.String()
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
      c2 := `nmap -Pn -A -sC  -f -D `+ newDecoys + " "+ trg.TargetIp.String()
      c3 := `nmap -p53,88,135,139,389,445,464,593,636,3268,3269,3389 -sC -sT -Pn `+ trg.TargetIp.String()
      c4 := `nmap -sS -Pn -A -T4 -sV -p- -f -D `+ newDecoys + " " + trg.TargetIp.String()
      commands = append(commands,c1)
      commands = append(commands,c2)
      commands = append(commands,c3)
      commands = append(commands,c4)
      return
    default:
      cmnd =  `nmap -sS -Pn -A -T4 -sV -p- -f -D `+ newDecoys + " " + trg.TargetIp.String()
  }
  commands = append(commands,cmnd)
  return
}

type Target struct{
  Host string //can be null if not specified as a subdomain
  HostIp net.IP
  TargetIp net.IP
  Decoys []net.IP
  FireWallName string
}
type Service struct{
  ServiceName string
  Port int
  Protocol string
	State   bool // open or closed
	Version string
  AT int
  //Data []*Output
}
type Data []*Output
type Output struct{
  Command string
  Output string
}

func  NmapScanForOpenPorts(name string,trg *Target) ([]*Service,Data){
  datas := make(chan *Output)
  svcs := make(chan *Service)
  var services []*Service
  var data []*Output
  commands := NmapCommandBuilder("all",trg)
  //var wg sync.WaitGroup
  //wg.Add(len(commands))
  for _,cmd := range commands {
  //  wg.Add(1)
  // defer wg.Done()
    //go func(cmnd string){
  //    defer wg.Done()
      fmt.Println("Scanning command ",cmd)
      output, err := exec.Command("sh", "-c", cmd).Output()
      if err != nil {
        //utils.Warning(fmt.Sprintf("NMAP: Error scanning for open ports: %s",err))
        fmt.Sprintf("NMAP: Error scanning for open ports: %s",err)
      }
      datum := &Output{
        Command: cmd,
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
      fmt.Println("Done with ",cmd)
  //  }(cmd)
  }
  for datum,ok := <-datas;ok; datum,ok = <-datas {
    if !ok {
      close(datas)
    }
    data = append(data,datum)
  }
  for service,ok := <-svcs;ok; service,ok = <-svcs {
    if !ok {
      close(svcs)
    }
    services = append(services,service)
  }
  //wg.Wait()
  return services,data
}

/*

func NmapScanForOpenPorts(name string, trg *Target) ([]*Service, []*Output) {
    datas := make(chan *Output)
    svcs := make(chan *Service)
    var services []*Service
    var data []*Output
    commands := NmapCommandBuilder("all", trg)
    var wg sync.WaitGroup
    for _, cmd := range commands {
        wg.Add(1)
        go func(cmnd string) {
            defer wg.Done()

            // Set a timeout for the command execution
            ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
            defer cancel()

            fmt.Println("Scanning command ", cmnd)
            output, err := exec.CommandContext(ctx, "sh", "-c", cmnd).Output()
            if err != nil {
                fmt.Sprintf("NMAP: Error scanning for open ports: %s", err)
                return
            }
            datum := &Output{
                Command: cmnd,
                Output:  string(output),
            }
            datas <- datum

            scanner := bufio.NewScanner(strings.NewReader(string(output)))
            r := regexp.MustCompile(`(\d+)\/(\w+)\s+(\w+)\s+([\w\/]+)(.*)`)
            for scanner.Scan() {
                line := scanner.Text()
                matches := r.FindStringSubmatch(line)
                if len(matches) == 6 {
                    portNumber, _ := strconv.Atoi(matches[1])
                    protocol := matches[2]
                    serviceState := matches[3]
                    serviceName := matches[4]
                    version := matches[5]
                    var service Service
                    service.ServiceName = serviceName
                    service.Port = portNumber
                    service.Protocol = protocol
                    if serviceState == "open" {
                        service.State = true
                    } else {
                        service.State = false
                    }
                    service.Version = version
                    switch protocol {
                    case "http", "https", "HTTP", "HTTPS":
                        service.AT = 1
                    default:
                        service.AT = 2
                    }
                    svcs <- &service
                }
            }
            fmt.Println("Done with ", cmnd)
        }(cmd)
    }

    // Wait for all goroutines to complete
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

    return services, data
}

func  NmapScanForOpenPorts(trg *Target) ([]*Service,Data){
  var services []*Service
  var data []*Output
  //reconData := make(chan ReconData)
  commands := NmapCommandBuilder("all",trg)
  for _,cmd := range commands {
    output, err := exec.Command("sh", "-c", cmd).Output()
    if err != nil {
      //utils.Warning(fmt.Sprintf("NMAP: Error scanning for open ports: %s",err))
      fmt.Sprintf("NMAP: Error scanning for open ports: %s",err)
      return nil,nil
    }
    datum := &Output{
      Command: cmd,
      Output: string(output),
    }
    data = append(data,datum)
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
        services = append(services,&service)
      }
    }
  }
  return services,data
}
*/
