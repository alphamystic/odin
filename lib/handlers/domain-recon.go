package handlers

import (
  "fmt"
  "net"
  "sync"
  "bytes"
  "os/exec"
  "bufio"
  "strings"
  "text/template"
  "github.com/alphamystic/odin/lib/utils"
)
/*
1. Use host -d to find the IP associate(if cloud flare ignore it) check for IPv4 or IPv6
2. Find alternate subdomains and their IP associates
3. Use theHarvester to find alternate IPs andemails
4. Use google dorks to find more data on target
*/
// in a way this is misusing channels, they are supposed to help two or moreprocecess comm with each other
// but I think it does that been a minute since I fixed odin :)
// whenever I wanted two proc to talk, I would:
/*type Proc1 struct { data chan*DataType}
type Proc2 struct { data chan*DataType}
proc1 writes to chan while proc2 reads from channel
an innit places the chanl of both as one
var  data1 chan *DataType
func InitProc(data chan *DataType)
InitProc(data)
just like wheagle runs commands
*/

type ReconDataOnSubdomain struct {
  FindSubdomainData  *FindSubdomains
  TheHarvesterData *TheHarvester
}

var temp *template.Template

type TheHarvester struct {
  Command string
  Output string
  Emails []string
  Names []string
  AssociateIps []net.IP
}

func DoReconOnDomain(host string,dtrgs chan<- *Target,done chan<- bool)  {
  hostIp := FindIp(host)
  // persist reconDataOnSubd to db
  reconDataOnSubd := DoReconDataOnSubdomain(host)
  var wg sync.WaitGroup
  wg.Add(len(reconDataOnSubd.FindSubdomainData.Subdomains))
  fmt.Println("Creating domain targets")
  for _,trg := range reconDataOnSubd.FindSubdomainData.SD {
    fmt.Sprintf("Creating target: %s",trg.Address)
    go func(t *SubDomainData){
      defer wg.Done()
      //now we construct a traget from each subdomain
      firewallName := CF(fmt.Sprintf("%s",t.Address))
      target := Target{
        Host: host,
        HostIp: hostIp,// ip of the host
        TargetIp: t.Address,
        Decoys: reconDataOnSubd.TheHarvesterData.AssociateIps,
        FireWallName: firewallName,
      }
      fmt.Sprintf("Created target %s",target.TargetIp)
      // write the traget into the channel
      select {
        case dtrgs <- &target:
            fmt.Println("Sent target to channel:", target.TargetIp)
        default:
            fmt.Println("Could not send target to channel:", target.TargetIp)
        }
    }(trg)
  }
  wg.Wait()
  // make the host a target and write it into the channelmas well
  domFirewallName := CF(host)
  hostTarget := Target{
    Host: host,
    HostIp: hostIp,
    TargetIp: hostIp,
    Decoys: reconDataOnSubd.TheHarvesterData.AssociateIps,
    FireWallName: domFirewallName,
  }
  dtrgs<-&hostTarget
  done<-true
  fmt.Println("Done creating domain targets")
}

// returns an ipv4 Address
func FindIp(domain string) net.IP{
  ips,_ := net.LookupIP(domain)
  for _, ip := range ips {
    if ipv4 := ip.To4(); ipv4 != nil {
      return ipv4
    } else {
      utils.Warning(fmt.Sprintf("RECON ERROR: Unable to find ip for domain name %s",domain))
      continue
    }
  }
  return nil
}

func  DoReconDataOnSubdomain(host string)ReconDataOnSubdomain {
  var rdSubd ReconDataOnSubdomain
  //find the subdomains
  // find their IPs
  th,_,subdomains := ProcessTheHarvester(host)
  fd := FindingSubdomains(host,subdomains)
  rdSubd = ReconDataOnSubdomain {
    FindSubdomainData: fd,
    TheHarvesterData: th,
  }
  return  rdSubd
}

func FindingSubdomains(target string,prevSubdomains []string)(*FindSubdomains){
  cmd := "subfinder -d {{.Domain}}" //later on I did ask myself why not just add a + target
  var buf bytes.Buffer
  temp = template.Must(template.New("cmd").Parse(cmd))
  _ = temp.Execute(&buf,struct{Domain string}{target})
  var subdomains []string
  var sd []*SubDomainData
  output, err := exec.Command("sh", "-c", buf.String()).Output()
  if err != nil {
    utils.Warning(fmt.Sprintf("SUFINDER: Error finding subdomains with subfinder %s: \nERROR: %s",buf.String(),err))
    return nil
  }
  subdomains = ProcessSubdomainOutput(target,string(output))
  subdomains = SanitizeSubdomains(subdomains,prevSubdomains)
  sd = GetSubdomainData(target,subdomains)
  //var fd FindSubdomains
  fd := &FindSubdomains {
    Command: buf.String(),
    Output: string(output),
    Subdomains: subdomains,
    SD: sd,
  }
  return fd
}

type SubDomainData struct{
  Host string //parent Host
  SubdName string
  Address net.IP
}
type FindSubdomains struct{
  Command string
  Output string
  Subdomains []string
  SD []*SubDomainData
}

func GetSubdomainData(host string, subdomains []string)(results []*SubDomainData) {
  var wg sync.WaitGroup
  // Create a channel to receive the results
  ch := make(chan *SubDomainData, len(subdomains))
  // Start a goroutine for each subdomain to find the IP address and create a SubDomainData struct
  for _, subd := range subdomains {
    wg.Add(1)
    go func(subdomain string) {
      defer wg.Done()
    //  START:
        ip := FindIp(subdomain)
        data := &SubDomainData {
          Host: host,
          SubdName: subdomain,
          Address: ip,
        }
        if data.Address == nil{
          utils.PrintTextInASpecificColor("red",fmt.Sprintf("Nill address for %s",data.SubdName))
          data.Address = ip
          utils.PrintTextInASpecificColor("red",fmt.Sprintf("possible ip is %s",data.Address))
          if data.Address == nil{ utils.PrintTextInASpecificColor("red",fmt.Sprintf("Ignoring subdomain %s for lack of IPAddress",data.SubdName));return }// we ignore this subdomain
        //  goto START
        }
      utils.PrintTextInASpecificColor("cyan",fmt.Sprintf("Found subdomain data for %s with IP %s\n",data.SubdName,data.Address))
      ch <- data
    }(subd)
  }
  // Close the channel once all goroutines have finished
  go func(){
    wg.Wait()
    close(ch)
  }()
  // Collect the results from the channel into a slice of SubDomainData structs
  for data := range ch {
    results = append(results, data)
  }
  return results
}
/*
func GetSubdomainData(host string,PrevSubdomains []string) []SubDomainData{
  var lsbd []SubDomainData
  var sbd SubDomainData //push this into the for loop
  var wg sync.WaitGroup
  subdomains := FindingSubdomains(host,PrevSubdomains)
  wg.Add(len(subdomains.Subdomains))
  for _, subd := range subdomains.Subdomains{
    go func(subdomain string){
      defer wg.Done()
      // get the Ip of each of those subd
      ip := FindIp(subdomain)
      sbd = SubDomainData{
        Host: host,
        SubdName: subd,
        Address: ip,
      }
      //append to list of subd
      lsbd = append(lsbd,sbd)
    }(subd)
  }
  wg.Wait()
  return lsbd
}*/

func ProcessSubdomainOutput(target,output string)(subdomains []string) {
  // Create a scanner to read the command output
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	// Iterate through the scanner's lines
	for scanner.Scan() {
		line := scanner.Text()
		// Check if the line is a subdomain
		if strings.HasSuffix(line, target) {
      subdomains = append(subdomains,line)
		}
	}
	if err := scanner.Err(); err != nil {
    utils.PrintTextInASpecificColorInBold("red",fmt.Sprintf("%s",err))
	}
  return
}

func SanitizeSubdomains(subdomains,prevSubdomains []string)[]string{
  // add the two arrays togetther
  for _,subd := range prevSubdomains {
    subdomains = append(subdomains,subd)
  }
  subdomains = utils.RemoveStringDuplicates(subdomains)
  return subdomains
}


// if given subdomains have an associate IP, read the IP and add it to the list of targetsAssociateIP
func ProcessTheHarvester(domain string) (*TheHarvester,error,[]string) {
  var buf bytes.Buffer
  var aips []net.IP
  cmd := `theHarvester -d {{.Domain}} -r -n -c -b anubis,certspotter,crtsh,baidu,brave`
  temp = template.Must(template.New("cmd").Parse(cmd))
  _ = temp.Execute(&buf,struct{Domain string}{domain})
  var th TheHarvester
  output, err := exec.Command("sh", "-c", buf.String()).Output()
  if err != nil {
    utils.Warning(fmt.Sprintf("THEHARVESTER: Error Executing theHarvester command %s\nERROR %s",buf.String(),err))
    return nil,err,nil
  }
  /*getNames := func(output string) []string{
    //get names from the output
  }
  getEmails := func(output string) []string{
    // get emails from the output
    var emails make([]string)
    line1 := `[*] No emails found.`
  }
  // get the emails
  emails := getEmails(output)
  //get any names found
  names := getNames(output)*/
  // get the SubDomains and associate IPs
  SeperateIpFromSubDomain := func(input string)[]string{
    var toBeSplit []string
    //split on new line
    lines := strings.Split(input, "\n")
    //iterate through the lines to seperate all IPs and Potential Subdomains
    for _, line := range lines {
     if strings.Contains(line,domain){
        toBeSplit = append(toBeSplit,line)
      }
    }
    return toBeSplit
  }
  Splitter := func(lines []string)([]string,[]string){
    // map to keep track of unique IPs
    uniqueIPs := make(map[string]bool)
    var uniqueIPList []string
    var stringList []string
    for _, line := range lines {
      if strings.Contains(line,":"){
        ///sepearte ip and subdomain
        parts := strings.Split(line, ":")
        stringPart := parts[0]
        ipPart := parts[1]
        stringList = append(stringList, stringPart)
        uniqueIPList = append(uniqueIPList, ipPart)
        if !uniqueIPs[ipPart] {
          uniqueIPs[ipPart] = true
          uniqueIPList = append(uniqueIPList, ipPart)
        }
      } else {
        //append to subdomains
        stringList = append(stringList,line)
      }
    }
    return stringList,uniqueIPList
  }
  lines := SeperateIpFromSubDomain(string(output))
  subdomains,ips := Splitter(lines)
  ips = utils.RemoveStringDuplicates(ips)
  for _,ip := range ips {
    ipn := net.ParseIP(ip)
    aips = append(aips,ipn)
  }
  th = TheHarvester {
    Command: buf.String(),
    Output: string(output),
    Emails: nil,
    Names: nil,
    AssociateIps: aips,
  }
  return &th,nil,subdomains
}

func CF(target string)string{
	data := HttpChecker(target)
	isFireWallPresent,name := CheckForFirewall(data,target)
	if !isFireWallPresent {
    return ""
	}
  if isFireWallPresent && utils.CheckifStringIsEmpty(name) {
  		return name
  }
  return ""
}
// don't know why I made this two redundant
func CheckFirewall(target string)(present bool,FWname string){
  data := HttpChecker(target)
	isFireWallPresent,name := CheckForFirewall(data,target)
	if !isFireWallPresent {
    return false,""
	}
  if isFireWallPresent && utils.CheckifStringIsEmpty(name) {
  		return true,name
  }
  return false,""
}
func CheckForFirewall(input,target string) (present bool,FWname string){
	lines := strings.Split(input, "\n")
  count := 0
	for _,line := range lines {
    if strings.Contains(line,"ERROR") {
      count += 2
      data := HttpsChecker(target,count)
      present,FWname = CheckForFirewall(data,target)
      break
    } else {
      if strings.Contains(line,"The site"){
  			//fmt.Println(line)
  			// search keyword
  			keyword := "behind"
  			// Find the index of keyword in the input string
  			index := strings.Index(line, keyword)
  			// Extract the substring starting from the index of keyword to the end of the line
  			name := line[index:]
  			unwanted := []string{"behind","WAF."}
  			for _,word := range unwanted {
  				name = strings.ReplaceAll(name,word,"")
  				}
  			present = true
  			FWname = name
  			break
  		}
		}
	}
	return present,FWname
}

var HttpsChecker = func(target string,count int)string{
  if count > 1{
    return "The site is behind Unamed WAF."
  }
	cmnd := `wafw00f -a -v https://`+target
	output, err := exec.Command("sh", "-c", cmnd).Output()
	if err != nil{
		utils.NoticeError(fmt.Sprintf("WAFW00F: Error executing command: %s \nERROR: %s\n",cmnd,err))
    return ""
	}
	return string(output)
}

var HttpChecker = func(target string) string{
	cmnd := `wafw00f -a -v http://`+target
	output, err := exec.Command("sh", "-c", cmnd).Output()
	if err != nil{
		utils.NoticeError(fmt.Sprintf("WAFW00F: Error executing command: %s \nERROR: %s\n",cmnd,err))
    return ""
	}
	return string(output)
}
