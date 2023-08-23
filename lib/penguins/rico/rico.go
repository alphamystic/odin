package rico

import (
  //"sync"
  "odin/lib/handlers"
)
/*
  * Implements bruteforce handlers.Vulnerabilities
  * Before each bruteforce, try looking for a common or potential CVE(We will try a neural net for this)
*/
type RICO struct{
  RD *handlers.ReconData
}

type BruteForcer interface{
  BruteForce() Vulnerabilities
}
// rico creates a work group then a channel of handlers.Vulnerabilities to be written into
// create a go routine for each of the scanners
// each scanner writes into the vulns channel and when each is done they signal the work group
/*
func RunRico()([]handlers.Vulnerabilities){
  vulns := make(chan []*handlers.Vulnerabilities)
  go func(){
    defer wg.Done()
    //put all scanners in here
  }()
  var Vscanner handlers.VulnerabilityScanner
  Vscanner = RICO{
    RD: reconData,
  }
  vulns := VScanner.VulnScanner()
  return vulns
}

// all this willl go into a goroutine(one for each)
func (r *RICO) VulnScannerw() (vulns chan []handlers.Vulnerabilities,err error){
  return nil,nil
}

func (r *RICO) RICOWebAttack(vulns chan []handlers.Vulnerabilities){
  vulns <- EnumearteWeb() //whatever vulnerability u find
  return
}

func (r *RICO) RICOFTPBruteforce(vulns chan []handlers.Vulnerabilities){
  return
}

func (r *RICO) SSHBruteforce(vulns chan []handlers.Vulnerabilities){
  return
}

func (r *RICO) TelnetBruteforce(vulns chan []handlers.Vulnerabilities){
  return
}

func (r *RICO) ADBruteforce(vulns chan []handlers.Vulnerabilities){
  return
}

func (r *RICO) RDPBruteforce(vulns chan []handlers.Vulnerabilities){
  return
}
*/

/*
var wg sync.WaitGroup
wg.Add(1)
vulns := make(chan []*handlers.Vulnerabilities)
go func(){
  defer wg.Done()
  //put all scanners in here
}()
wg.Wait()
*/
