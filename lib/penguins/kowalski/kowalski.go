package kowalski

import (
  "fmt"
  "net"
  "time"
  "sync"
  "github.com/alphamystic/odin/lib/db"
  "github.com/alphamystic/odin/lib/utils"
  "github.com/alphamystic/odin/lib/handlers"
  "github.com/alphamystic/odin/lib/penguins/ph"
)

/*
  * Kowalsi does all the persistance
    1. Scans
    2. Targets
    3. Vulnerabilities
    4. Exploits
*/

type KOWALSKI struct{
  Targets []string
  RMode ph.Mode // should be running modeor something like that
  Ac ph.AttackCommands
  Name string
  ScanID string
}

// this channels should be modified to handle certain  number of buffers or be locked until theyve been written into/out
// alternatively use context and also pass it to the db writer
func (k *KOWALSKI) Kowalski_Analysis(exploits chan<- *handlers.Exploit,exploitsDone chan<- bool){
  var targets = []*handlers.Target{}
  for _,t := range k.Targets{
    if utils.CheckIfStringIsIp(t){
      // add it to the targets array
      firewallName := handlers.CF(t)
      trg := &handlers.Target{
        Host: t,
        HostIp: net.ParseIP(t),
        TargetIp: net.ParseIP(t),
        Decoys: []net.IP{net.ParseIP(t),net.ParseIP("2.2.2.2"),net.ParseIP("4.4.4.4"),net.ParseIP("8.8.8.8")},
        FireWallName: firewallName,
      }
      targets = append(targets,trg)
    } else {
      if utils.CheckIfStringIsDomainName(t) {
        //do recon for domain type to get targets IP addresses
        utils.PrintInformation(fmt.Sprintf("Getting Target --> Recon Data for domain name: %s",t))
        targetChan := make(chan *handlers.Target)
        done := make(chan bool)
        go handlers.DoReconOnDomain(t,targetChan,done)
        time.Sleep(100 * time.Second)
        /*for {
          select {
          case dt,ok:= <- targetChan:
            if !ok {
              return
            }
          default: //again do nothing
          }
        }*/
        for dt,ok := <-targetChan;ok; dt,ok = <-targetChan {
          if !ok {
            <- done
            close(targetChan)
          }
          fmt.Printf("Reading from targets channel. Received %+v\n",dt.TargetIp)
          targets = append(targets,dt)
        }
      } else {
        utils.NoticeError(fmt.Sprintf("Invalid target: %s",t))
      }
    }
  }
  targets = sanitizeTargets(targets)
  //write the targets t db
  err := SaveTargetsTODB(k.Name,targets)
  if err != nil{
    utils.NoticeError(fmt.Sprintf("Scans for %s not saved to db.",k.Name))
    utils.Logerror(err)
  }
  //do recon on the targets
  var wg sync.WaitGroup
  wg.Add(len(targets))
  vulns := make(chan []*handlers.Vulnerabilities)
  outReconData := make(chan *handlers.ReconData)
  reconDone := make(chan bool)
  vulnsDone := make(chan bool)
  doneCreatingExploits := make(chan bool)
  go k.VulnerabilityScanner(outReconData,reconDone,len(targets),vulns,vulnsDone)
  go k.CreateExploits(vulns,exploits,doneCreatingExploits)
  //CreateExploits(vulns<- chan []*handlers.Vulnerabilities, exploits<- *handlers.Exploit,exploitsDone<- chan bool)
  fmt.Println("Ranging through targets")
  for _,trg := range targets { // ad worker group here
    go func(target *handlers.Target){
      defer wg.Done()
      utils.PrintTextInASpecificColorInBold("blue",fmt.Sprintf("Doing recon for IP: %s",target.TargetIp))
      target.Recon(k.Name,outReconData)
      utils.PrintTextInASpecificColorInBold("white",fmt.Sprintf("Done scanning for vulnerabilies for IP: %s",target.TargetIp))
    }(trg)
  }
  wg.Wait()
  <- reconDone
  close(outReconData)
  utils.PrintTextInASpecificColor("yellow","Done doing recon for all targets....................")
  <-vulnsDone
  close(vulns)
  utils.PrintTextInASpecificColor("yellow","Done scanning vulnerabilities for all targets....................")
  for {
    select {
    case <- doneCreatingExploits:
      exploitsDone<- true
      return
    default:
      //do nothing
    }
  }
}

// change to read from inReconData instead of ranging. THen send done signal when done
// load the recon data into a scanner as they come in.
// for brevity, take n a count such that it's the length of the target and increase it as rd comes in and close when done
func (k *KOWALSKI) VulnerabilityScanner(inReconData <- chan *handlers.ReconData,reconDone chan<- bool,val int,vulns chan<- []*handlers.Vulnerabilities,vulnsDone chan<- bool){
  for rd,ok := <- inReconData;ok; rd,ok = <- inReconData{
    var count int
    if !ok{
      if count == val {
        reconDone<-true
      }
    }
    count = count + 1
    fmt.Sprintf("Recon Data on: %s",rd.Trg.TargetIp)
    if err := SaveReconDataTODB(k.Name, rd); err != nil{
      utils.Logerror(err)
    }
    fmt.Println(rd)
    //vulns<- vulnerabilityFound
  }
  vulnsDone<-true
}

func (k *KOWALSKI) CreateExploits(vulns<- chan []*handlers.Vulnerabilities, exploits chan<- *handlers.Exploit,doneCreatingExploits chan<- bool) {
  var wg sync.WaitGroup
  var vlns  []*handlers.Vulnerabilities
//  var grpdVulns = make(map[net.IP][]*Vulnerabilities)
  for mvulns,ok := <- vulns;ok; mvulns,ok = <- vulns{
    wg.Add(1)
    if !ok{
      fmt.Println("We are now reading from a closed channels or something........")
      continue
    //  exploitsDone<-true //or something of the sought
    //if you write done and wg isn't done then it willl close some vulns/exploits out
    }
    if err := SaveVulnerabilitiesTODB(k.Name, mvulns); err != nil{
      utils.Logerror(err)
    }
    fmt.Println(mvulns)
    for _,vuln := range mvulns{
      go func(vln *handlers.Vulnerabilities){
      utils.PrintTextInASpecificColor("cyan",fmt.Sprintf("Vulnerailities on on: %s",vuln.Trg.TargetIp.String()))
        wg.Add(1)
        defer wg.Done()
        if vln.Grouped{
          fmt.Println("Doing the groued exploits thingy")
          vlns = append(vlns,vln)
          // so if I do this reeatedly it will overwrite, how do I keep them all in?
        } else {
          var expVulns  []*handlers.Vulnerabilities
          expVulns = append(expVulns,vln)
          exploit := &handlers.Exploit {
          	Trg: vuln.Trg,
          	LHOST: "msid",
          	LPORT: 5000,
          	Address: "0.0.0.0",
          	Target: vln.Trg.TargetIp.String(),
          	AverageSeverity: 9,
          	Grouped: false,
          	Vulns: expVulns,
          	Works: true,
          }
          exploits<- exploit
        }
      }(vuln)
    }
  }
  wg.Wait()
  doneCreatingExploits<- true
  fmt.Println(vlns)// will sort this later
}

func sanitizeTargets(targets []*handlers.Target) []*handlers.Target {
  fmt.Println("Sanitizing targets.........s")
  // Create a map to store unique IP addresses
  uniqueIPs := make(map[string]bool)
  //ignore default decoys
  assumedDecoys := []net.IP{net.ParseIP("2.2.2.2"), net.ParseIP("4.4.4.4"), net.ParseIP("8.8.8.8")}
  // Create a new slice to store the sanitized targets
  sanitizedTargets := make([]*handlers.Target, 0)
  for _, target := range targets {
    // Check if the host IP address is unique
    if _, ok := uniqueIPs[target.HostIp.String()]; !ok {
      uniqueIPs[target.HostIp.String()] = true
      // Add the target to the sanitized targets slice
      sanitizedTargets = append(sanitizedTargets, target)
    }
    // Check if the target IP address is unique
    if _, ok := uniqueIPs[target.TargetIp.String()]; !ok {
      uniqueIPs[target.TargetIp.String()] = true
      // Add the target to the sanitized targets slice
      sanitizedTargets = append(sanitizedTargets, target)
    }
    // Check if any of the decoy IP addresses are unique
    for _, decoyIP := range target.Decoys {
      //if they are in default set them to be seen so as they can be assumed
      if decoyIP.Equal(assumedDecoys[0]) || decoyIP.Equal(assumedDecoys[1]) || decoyIP.Equal(assumedDecoys[2]) {
        //target.Decoys[i] = assumedDecoys[i]
        uniqueIPs[decoyIP.String()] = true
      }
      if _, ok := uniqueIPs[decoyIP.String()]; !ok {
        uniqueIPs[decoyIP.String()] = true
        // Add the target to the sanitized targets slice
        trg := &handlers.Target{
          Host: target.Host,
          HostIp: target.HostIp,
          TargetIp: decoyIP,
          Decoys: target.Decoys,
          FireWallName: "",
          //RoSD *ReconDataOnSubdomain persist this to DB When done
        }
        sanitizedTargets = append(sanitizedTargets, trg)
      }
    }
  }
  return sanitizedTargets
}

/* well I just realised to chain vulnerabilities we are going to have to learn how to sought them as per target that way they can be chained
// to do recon on a target we need to ensure each is an IP address or atleast not some akamai or cloudflare
func Kowalski_Analysis(target []string) chan handlers.Exploit {
  exploits := make(chan handlers.Exploit)
  var targets = []*handlers.Target{}
  //reconData := make(chan *handlers.ReconData)
  //var kowalski KOWALSKI
  for _,t := range target{
    if utils.CheckIfStringIsIp(t){
      // add it to the targets array
      firewallName := handlers.CF(t)
      trg := &handlers.Target{
        Host: t,
        HostIp: net.ParseIP(t),
        TargetIp: net.ParseIP(t),
        Decoys: []net.IP{net.ParseIP(t),net.ParseIP("2.2.2.2"),net.ParseIP("4.4.4.4"),net.ParseIP("8.8.8.8")},
        FireWallName: firewallName,
      }
      targets = append(targets,trg)
    } else {
      if utils.CheckIfStringIsDomainName(t) {
        //do recon for domain type to get targets IP addresses
        utils.PrintInformation(fmt.Sprintf("Getting Target --> Recon Data for domain name: %s",t))
        dtrg := handlers.DoReconOnDomain(t)
        for dt,ok := <-dtrg;ok; dt,ok = <-dtrg {
          if !ok {
            close(dtrg)
          }
          targets = append(targets,dt)
        }
      } else {
        utils.NoticeError(fmt.Sprintf("Invalid target: %s",t))
      }
    }
  }
  outReconDatas := make(chan *handlers.ReconData)
  var wg sync.WaitGroup
  wg.Add(len(targets))
  for _,trg := range targets{ // ad worker group here
    go func(target *handlers.Target){
      defer wg.Done()
      utils.PrintTextInASpecificColorInBold("blue",fmt.Sprintf("Doing recon for IP: %s",target.TargetIp))
      reconData := target.Recon(outReconData)
      for val,ok := <- reconData;ok; val,ok = <- reconData{
        if !ok{
          //do something
        }
        fmt.Println("Can we get here")
        reconDatas<- val
      }
    }(trg)
  }
  wg.Wait()
  //check reconDatas to ensure it's all done the close that channel
  for rd,ok := <-reconDatas;ok; rd,ok = <-reconDatas{
    if !ok {
      close(reconDatas)
    }
    fmt.Println("I want us here")
    //enumerate and scan for vulnerabilities
    // scan for vulnerabilities for each
    utils.PrintTextInASpecificColorInBold("white",fmt.Sprintf("Scanning vulnerabilities for IP: %s",rd.Trg.TargetIp))
    return nil
  }
  /*rico := RICO{RD:&reconData}
  private := PRIVATE{RD:&reconData}
  go func(){
    //create a channel for getting vulnerabilities from both rico and private
    // when both are done, we calll create create exploits and place the vulns found into an exploit
    // let CreaeExploit write into the exploits into a channel
    vulns := make(chan handlers.Vulnerability)
    var vulnerabilities []handlers.Vulnerability
    go func() { vulns <- rico.Scanner() }()
    go func() { vulns <- private.Scanner() }()
    for i := 0; i < 3;i++{
      v := <-vulns
      vulnerabilities = append(vulnerabilities,v...)
    }
    close(vulns)// this is perenial not sure but logic says functions are done and should be closed
    exploits = kowalski.CreateExploits(vulnerabilities)
  }()
  close(exploits)
  return exploits
}

// do recon from recon interface
// for bruteforce vulnerabilities call rico
// for other vulnerabilities call private
// now this is where that classifier thingy comes in.
type ExploitManager struct{
  Exploits map[net.IP][]*handlers.Vulnerabilities
}
*/

/*
type Vuln struct{
  Id net.IP
  vulns []*handlers.Vulnerability
}
func (bip Vuln) ByIp()net.IP{return bip.Id}
var Sorter = func(vlns []*handlers.Vulnerability)[]Vuln{
  sort.Sort(Vuln(vlns))
}
*/
var SaveTargetsTODB = func(name string,targets []*handlers.Target)error{
  driver,err := db.Old("../.brain/scans/" + name,0644)
  if err != nil{
    return err
  }
  for _, target := range targets{
    if err := driver.Write("targets",target.TargetIp.String(),target); err != nil{
      utils.Logerror(fmt.Errorf("Error saving target %s to db.\nERROR: %v",target.TargetIp.String(),err))
      continue
    }
  }
  return nil
}

var SaveVulnerabilitiesTODB = func(name string,vulns []*handlers.Vulnerabilities)error{
  driver,err := db.Old("../.brain/scans/" + name,0644)
  if err != nil{
    return err
  }
  for _, vuln := range vulns{
    time.Sleep(1 * time.Second)
    str := utils.RandString(5)
    if err := driver.Write("vulnerabilities",str + vuln.Trg.TargetIp.String(),vuln); err != nil{
      utils.Logerror(fmt.Errorf("Error saving vulnerability for %s to db.\nERROR: %v",vuln.Trg.TargetIp.String(),err))
      continue
    }
  }
  return nil
}

var SaveReconDataTODB = func(name string,rd *handlers.ReconData) error{
  driver,err := db.Old("../.brain/scans/" + name + "/" + "recondata" + "/" + rd.Trg.TargetIp.String(),0644)
  if err != nil{
    return err
  }
  err = driver.Write("services",rd.Trg.TargetIp.String(),rd.Services)
  if err != nil{
    utils.Logerror(fmt.Errorf("Error writing services to db for %s.\nERROR: %v",rd.Trg.TargetIp.String(),err))
  }
  err = driver.Write("webdata",rd.Trg.TargetIp.String(),rd.WD)
  if err != nil{
    utils.Logerror(fmt.Errorf("Error writing webdata to db for %s.\nERROR: %v",rd.Trg.TargetIp.String(),err))
  }
  return nil
}
/*
for rd := range inReconData{
  //save to db
  // scanfor vulns and find a way to group them fo exploit creation
  if rd == nil{
    fmt.Println("We have nil recon data.")
    return
  }
  fmt.Sprintf("Recon Data on: %s",rd.Trg.TargetIp)
  fmt.Println(rd)
  done<-true
}
*/
