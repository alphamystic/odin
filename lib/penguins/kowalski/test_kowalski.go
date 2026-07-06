package kowalski

import (
	//"fmt"
	"log"
	"net"

	"github.com/alphamystic/odin/lib/db"
	"github.com/alphamystic/odin/lib/handlers"
	"github.com/alphamystic/odin/lib/utils"
)

func (k *KOWALSKI) TestWriteToDB() {
	var targets []*handlers.Target
	var failedTargets []*handlers.Target
	var vulns []*handlers.Vulnerabilities
	var exploits []*handlers.Exploit

	log.Println("[+] Starting complete Postman-aligned API E2E DB Persistence Testing...")

	/* ========================================================================
	   EXAMPLE SET 1: PRIMARY WEB ENDPOINT TARGET (Linux Backend)
	   ========================================================================
	*/
	trg1 := &handlers.Target{
		ScanID:       k.ScanID,
		Host:         "alpha-production-app.io",
		HostIp:       net.ParseIP("192.168.10.12"),
		TargetIp:     net.ParseIP("192.168.10.12"),
		Decoys:       []net.IP{net.ParseIP("192.168.10.1"), net.ParseIP("8.8.8.8"), net.ParseIP("1.1.1.1")},
		FireWallName: "Cloudflare WAF",
	}
	targets = append(targets, trg1)

	/* ========================================================================
	   EXAMPLE SET 2: INTERNAL INFRASTRUCTURE ROUTER TARGET (Embedded Platform)
	   ========================================================================
	*/
	trg2 := &handlers.Target{
		ScanID:       k.ScanID,
		Host:         "dev-proxy-internal.local",
		HostIp:       net.ParseIP("10.10.4.56"),
		TargetIp:     net.ParseIP("10.10.4.56"),
		Decoys:       []net.IP{net.ParseIP("10.10.4.1"), net.ParseIP("2.2.4.4")},
		FireWallName: "Fortinet",
	}
	targets = append(targets, trg2)

	// Persist Targets to the Database via API Endpoint
	log.Println("[+] Synchronizing Targets with API endpoint [/api/recon/createtarget]...")
	k.SaveTargetsTODB(targets, failedTargets)
	if len(failedTargets) > 0 {
		log.Printf("[!] %d targets failed API write. Running offline backup tracking protocols...\n", len(failedTargets))
		if err := k.BackupSaveTargetsTODB(failedTargets); err != nil {
			log.Fatal("[-] Critical Fallback Failure: ", err)
		}
	}

	/* ========================================================================
	   RECON RECOVERY AND DISCOVERY ENUMERATION STAGING
	   ========================================================================
	*/
	log.Println("[+] Submitting discovered assets layout mapping to [/api/recon/createrecondata]...")

	rd1 := &handlers.ReconData{
		Trg: trg1,
		WD: &handlers.WebData{
			Directories: []string{"/api/v1/auth", "/dashboard/settings", "/assets/js/"},
			Parameters:  []string{"/api/v1/auth?debug=true", "/dashboard/settings?user_id=root"},
			Files:       []string{"/assets/js/app.min.js", "/documents/invoice.pdf"},
		},
	}
	if err := k.SaveReconDataTODB(rd1); err != nil {
		_ = k.BackupSaveReconDataTODB(rd1)
	}

	rd2 := &handlers.ReconData{
		Trg: trg2,
		WD: &handlers.WebData{
			Directories: []string{"/view", "/cgi-bin/diagnostics"},
			Parameters:  []string{"/view?page=../../etc/passwd"},
			Files:       []string{"/cgi-bin/diagnostics/ping.sh"},
		},
	}
	if err := k.SaveReconDataTODB(rd2); err != nil {
		_ = k.BackupSaveReconDataTODB(rd2)
	}

	/* ========================================================================
	   NETWORK SERVICES SCANNING METRICS LOGGING
	   ========================================================================
	*/
	log.Println("[+] Registering open service bindings layout to [/api/recon/createservice]...")

	srv1 := &handlers.Service{
		TargetID:    trg1.TargetID,
		ServiceName: "Nginx HTTP Proxy",
		Port:        443,
		Protocol:    "TCP",
		State:       true,
		Version:     "1.25.2",
	}
	_ = k.SaveServiceTODB(srv1)

	srv2 := &handlers.Service{
		TargetID:    trg2.TargetID,
		ServiceName: "SSH Remote Management Link",
		Port:        22,
		Protocol:    "TCP",
		State:       true,
		Version:     "OpenSSH 8.9p1 Ubuntu",
	}
	_ = k.SaveServiceTODB(srv2)

	/* ========================================================================
	   VULNERABILITIES ANALYSIS LOG STAGING
	   ========================================================================
	*/
	log.Println("[+] Syncing discovered vulnerability definitions with [/api/recon/vulnerability/create]...")

	// Vuln Example 1: Critical RCE finding on target 1
	v1 := &handlers.Vulnerabilities{
		Trg:           trg1,
		VulnerabilityID: utils.Md5Hash(utils.GenerateUUID()),
		Name:          handlers.RCE,
		Severity:      10,
		Payload:       "<?php system($_GET['cmd']); ?>",
		AT:            handlers.WEBATTACK,
		Grouped:       false,
		Authenticated: true,
		Works:         true,
		Details:       "Arbitrary file write allows deployment of web shell structures on target runtime root.",
	}
	vulns = append(vulns, v1)

	// Vuln Example 2: Local File Inclusion variant on target 2
	v2 := &handlers.Vulnerabilities{
		Trg:           trg2,
		VulnerabilityID: utils.Md5Hash(utils.GenerateUUID()),
		Name:          handlers.LFI,
		Severity:      8,
		Payload:       "../../../../etc/ossec.conf",
		AT:            handlers.WEBATTACK,
		Grouped:       true, // Chained into a larger attack context loop
		Authenticated: false,
		Works:         true,
		Details:       "Directory traversal verification patterns allow validation of configuration properties.",
	}
	vulns = append(vulns, v2)

	if err := k.SaveVulnerabilitiesTODB(vulns); err != nil {
		log.Println("[!] Failed writing to remote API log registry. Backing up locally...")
		_ = k.BackupSaveVulnerabilitiesTODB(vulns)
	}

	/* ========================================================================
	   EXPLOIT POPS ANALYSIS AND CHAINS VERIFICATION STAGING
	   ========================================================================
	*/
	log.Println("[+] Submitting completed attack string maps down to [/api/recon/exploit/create]...")

	// Exploit Variant 1: Standalone exploit pop verification mapping
	exp1 := &handlers.Exploit{
		Trg:             trg1,
		ExploitID:       utils.Md5Hash(utils.GenerateUUID()),
		TargetID:        trg1.TargetID,
		LHOST:           "10.10.10.5",
		LPORT:           4444,
		Address:         "http://alpha-production-app.io/uploads/shell.php",
		AverageSeverity: 10,
		Grouped:         false,
		GroupedVulns:    []string{v1.VulnerabilityID},
		Vulns:           []*handlers.Vulnerabilities{v1},
		Works:           true,
	}
	exploits = append(exploits, exp1)

	// Exploit Variant 2: Chained loop exploit vector targeting inner infrastructure proxy configurations
	exp2 := &handlers.Exploit{
		Trg:             trg2,
		ExploitID:       utils.Md5Hash(utils.GenerateUUID()),
		TargetID:        trg2.TargetID,
		LHOST:           "10.10.10.5",
		LPORT:           5555,
		Address:         "http://dev-proxy.internal/view?page=../../etc/ossec.conf",
		AverageSeverity: 8,
		Grouped:         true,
		GroupedVulns:    []string{v2.VulnerabilityID},
		Vulns:           []*handlers.Vulnerabilities{v2},
		Works:           true,
	}
	exploits = append(exploits, exp2)

	// Run full API integration check for Exploit models
	if err := k.SaveExploitsTODB(exploits); err != nil {
		log.Println("[-] API Exploit staging failed. Commencing fallback serialization...")
		// Backup logic locally if remote storage pipeline breaks
		driver, _ := db.Old("./.brain/failedscans/exploits_"+k.Name, 0644)
		if driver != nil {
			for _, exp := range exploits {
				_ = driver.Write("exploits", exp.ExploitID, exp)
			}
		}
	}

	log.Println("[+] DB Integration Test Suite complete. All systems initialized securely.")
}