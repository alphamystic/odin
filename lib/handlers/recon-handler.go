package handlers

import(
  "fmt"
  "github.com/alphamystic/odin/lib/utils"
)
/*
  *Handles recon output
  * Does more recon if needed
*/

func (trg *Target) Recon(name string,outRecon chan<- *ReconData){
  //check the presence of a firewall
  present,name := CheckFirewall(fmt.Sprintf("%s",trg.TargetIp))
  if !present{
    //firewall not present
    utils.PrintInformation(fmt.Sprintf("Firewall not present on: %s",trg.TargetIp))
    wd := CreateWebData()
    services := NmapScanForOpenPorts(name,trg)
    //construct rd and write to reconData
    rd := &ReconData {
      Trg: trg,
      Services: services,
      WD: wd,
    }
    outRecon <- rd
    utils.PrintTextInASpecificColor("magenta","Done doing recon for "+ trg.TargetIp.String())
  } else {
    // doing enumeration, we will scan for firwall again this time with ports
    if present && utils.CheckifStringIsEmpty(name) {
      // set the firewall name
      utils.PrintInformation(fmt.Sprintf("Firewall present on %s Name: %s",trg.TargetIp,name))
      //ignore if  cloudflare or akamai
      if name == "Akamai" || name =="Cloudflare" || name == "Cloudfront" || name ==  "AWS/Cloudfront" {
        //find the parameters,directories and files
        wd := CreateWebData()
        rd := &ReconData {
          Trg: trg,
          Services: nil,
          WD: wd,
        }
        outRecon <-rd
      } else {
        //check cve for that firewall or ignore the target/ip
      }
    }
  }
}

/*
nmapscannetwork
nmap scan for file shares nmap --script=nfs-ls 192.169.0.103
func NmapCommandBuilder()[]string{
  normal
  withDecoys

}*/
//should return open ports and potential services
//name port protocol service information


func SanitizeNmapServices(){}

func CreateWebData()*WebData{
  //dirb http://10.5.5.25:8080/ -w
  return nil
}

func FindDirectories()[]string{
  var directories []string
  return directories
}

func FindParameters()[]string{
  var parameters []string
  return parameters
}

func Crawl()([]string,[]string){
  var files []string
  var parameters []string
  return files,parameters
}
