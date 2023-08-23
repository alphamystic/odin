package handlers

/*
  * Does recon on a particular target returning
    1. Find Subdomains and return them with their IP Addreses {SubDomain:Ip_ADDRESS}
    2. Find Any other Alternative IP Address (Checkk if it's associated with an IP or they are already in place)
    3. Services running
        For each subd and alternate IP check for the services running and open ports
          {IP:{SERVICES/PORT}} //also check for iot/scada systems
    4. Try findning the vulnerabilities available and the payloads to be used >> Send to classifier to do this
*/

import (
  "net"
)

/*
// taking every aaddress as a subdomain, the recon therefore has a rocess of:
  1.Get the target.
  2. Check if it's a domain or IP
  3. If it's an IP, get the recon data straight, if not?
  4. Use  (You can create a Host,AssociateIP struct to assist in data handling)
      1. Use host -d to find the IP associate thn check for firewall presence(if cloud flare ignore it) check for IPv4 or IPv6
      2. Find alternate subdomains and their IP associates
      3. Use theHarvester to find alternate IPs andemails
      4. Use google dorks to find more data on target
  ## DO RECON
  5. Check for firewall presence in each.
  6. Scan for open ports and services
  ## Now we should know the service version and scan for CVE's
  7. If web service, check for CMS,directories,parameters,cgi
  If others, sent to rico for bruteforce
  The neural net should enumerate all the remaining services like smb,ad and find an exploit to create another way in
      this can be directory mounting,DNS,NFS,SMTP,SNMP,RDP,SSH,FTP,Responder for hashes and NTLM

    *In a way this has skipper written all over it but only time will answer it*
*/

type Target struct{
  Host string //can be null if not specified as a subdomain
  HostIp net.IP
  TargetIp net.IP
  Decoys []net.IP
  FireWallName string
  //RoSD *ReconDataOnSubdomain persist this to DB When done
}

type ReconData struct{
  Trg *Target
  /*CommandUsed string
  CommandOutput string*/
  Services []*Service
  WD  *WebData
}

type WebData struct{
  Directories []string
  Parameters []string// should also be associated with the directory it came from
  Files string
}

type Service struct{
  ServiceName string
  Port int
  Protocol string
	State   bool // open or closed
	Version string
  AT AttackType
  //Data []*Output
}

type ServiceData []*Output

type Output struct{
  Command string
  Output string
}

// Not sure on what kind of data it should return but a vulnerabilities will do for now
type CVE struct {
  CVEID string
  Present bool
  Confirmed bool
  Vuln Vulnerabilities
  POC Exploit
}

type CVEChecker interface {
  CVECheck() *CVE
}

// find associatie and verifies ownership
type Whois struct{}
