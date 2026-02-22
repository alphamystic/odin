package kowalski


import (
  "fmt"
  "log"
  "net"
  "github.com/alphamystic/odin/lib/utils"
  "github.com/alphamystic/odin/lib/handlers"
  "github.com/alphamystic/odin/lib/penguins/ph"
)


func (k *KOWALSKI) TestWriteToDB(){
  var targets = []*handlers.Target{}
  var failedTarges = []*handlers.Target{}
  var vulns =  []*handlers.Vulnerabilities{}
  trg := &handlers.Target {
    ScanID: k.ScanID,
    Host: "host.com",
    HostIp: net.ParseIP("3.3.3.3"),
    TargetIp: net.ParseIP("3.3.3.3"),
    Decoys: []net.IP{net.ParseIP(t),net.ParseIP("2.2.2.2"),net.ParseIP("4.4.4.4"),net.ParseIP("8.8.8.8")},
    FireWallName: "Cloudflare",
  }
  append(targets,trg)
  append(failedTarges,trg)
  trg2 := &handlers.Target {
    ScanID: scan_id,
    Host: "exampleHost.com",
    HostIp: net.ParseIP("3.4.4.3"),
    TargetIp: net.ParseIP("3.4.4.3"),
    Decoys: []net.IP{net.ParseIP(t),net.ParseIP("2.2.2.2"),net.ParseIP("4.4.4.4"),net.ParseIP("8.8.8.8")},
    FireWallName: "Cloudflare",
  }
  append(targets,trg2)
  append(failedTarges,trg2)
  k.SaveTargetsTODB(targets, failedTarges)
  if err := k.BackupSaveTargetsTODB(failedTarges); err != nil{
    log.fatal(err)
  }
  // write recon data to db
  log.Println("Writing web data to body...")
  rd := &handlers.WebDataBody{
		TargetID:    trg.TargetID,
		Directories: []string{"target.com/trg", "target.com/dire/two"},
		Parameters:  []string{"target.com/trg?par=one", "target.com/dire/two/par?=two"},
		Filepaths:   []string{"target.com/trg/target.pdf", "target.com/dire/two.png"},
	}
  if err := k.SaveReconDataTODB(rd); err != nil{
    log.Println("Error api saving recondata: ",rd)
    if err := k.BackupSaveReconDataTODB(rd); err != nil {
      log.Fatal(err)
    }
  }
  // write vulnerabilities to DB
  vuln := &handlers.Vulnerailities{
    Trg: trg,
    Name: handlers.RCE,
  	Severity : 10,
  	Target: trg.TargetIp.String(),
    Payload : "<? exec (gethsll.pp:1.1.1.1) ?>",
    AT : handlers.WEBATTACK,
  	Grouped: false,
    Authenticated : false,
  	Works : true,
    Details : "A shell execution via system call in php",
  }
  append(vulns, vuln)
  k.SaveVulnerabilityToDB(vuln)
  if err := k.BackupSaveVulnerabilitiesTODB(vulns); err != nil {
    log.Fatal(err)
  }
  //Chain the vulnerability to an exploit and write an exploit to db
}
