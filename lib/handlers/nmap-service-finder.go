package handlers

import (
	"bufio"
	"context"
	"fmt"
	"regexp"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alphamystic/odin/lib/db"
	"github.com/alphamystic/odin/lib/utils"
)

// NmapScanForOpenPortsWithContext safe, throttled network scanning pipeline function
func NmapScanForOpenPortsWithContext(ctx context.Context, name string, trg *Target) []*Service {
	commands := NmapCommandBuilder("all", trg)
	utils.PrintInformation(fmt.Sprintf("[*] NMAP: Launching bounded service profiling for target: %s", trg.TargetIp))

	datas := make(chan *Output, len(commands))
	svcs := make(chan *Service, len(commands)*20)

	var wg sync.WaitGroup

	// Throttling semaphore worker limits to prevent socket or descriptor exhaustion rules
	maxConcurrentScans := 2
	semaphore := make(chan struct{}, maxConcurrentScans)

	for _, cmd := range commands {
		wg.Add(1)
		go func(cmnd string) {
			defer wg.Done()

			select {
			case semaphore <- struct{}{}: // Acquire worker slot token
			case <-ctx.Done():
				return
			}
			defer func() { <-semaphore }() // Release slot token

			// Inject a strict timeout bound explicitly per single Nmap command string execution
			cmdCtx, cancel := context.WithTimeout(ctx, 25*time.Minute)
			defer cancel()

			utils.Notice(fmt.Sprintf("[*] Executing sub-process profile: %s", cmnd))

			// Build system command wrapper bound explicitly to our active operational lifecycle context
			osCmd := exec.CommandContext(cmdCtx, "sh", "-c", cmnd)
			output, err := osCmd.CombinedOutput()
			if err != nil {
				utils.Warning(fmt.Sprintf("[-] NMAP process exception caught on query (%s): %v", cmnd, err))
				// Continue processing rather than crashing to allow incomplete buffer loops to clear out safely
			}

			datum := &Output{
				Command: cmnd,
				Output:  string(output),
			}

			select {
			case datas <- datum:
			case <-ctx.Done():
				return
			}

			// Service Parsing Engine Hook Block
			scanner := bufio.NewScanner(strings.NewReader(string(output)))
			r := regexp.MustCompile(`(?i)^(\d+)\/(\w+)\s+(\w+)\s+([\S]+)(.*)?`)

			for scanner.Scan() {
				line := scanner.Text()
				matches := r.FindStringSubmatch(line)
				if len(matches) == 6 {
					portNumber, _ := strconv.Atoi(matches[1])
					protocol := matches[2]
					serviceState := matches[3]
					serviceName := matches[4]
					version := strings.TrimSpace(matches[5])

					service := &Service{
						TargetID:    trg.TargetID,
						ServiceName: serviceName,
						Port:        portNumber,
						Protocol:    strings.ToUpper(protocol),
						State:       (serviceState == "open"),
						Version:     version,
					}

					if strings.EqualFold(protocol, "tcp") && (portNumber == 80 || portNumber == 443 || portNumber == 8080) {
						service.AT = WEBATTACK
					} else {
						service.AT = BRUTEFORCE
					}

					select {
					case svcs <- service:
					case <-ctx.Done():
						return
					}
				}
			}
		}(cmd)
	}

	// Dynamic lifecycle closure manager routine
	go func() {
		wg.Wait()
		close(datas)
		close(svcs)
	}()

	var services []*Service
	var data []*Output

	for datum := range datas {
		data = append(data, datum)
	}
	for service := range svcs {
		services = append(services, service)
	}

	utils.PrintTextInASpecificColor("green", fmt.Sprintf("[+] NMAP: Complete service extraction inventory resolved for host: %s", trg.TargetIp))

	if err := SaveServiceDataTODB(name, trg.TargetIp.String(), data); err != nil {
		utils.Logerror(err)
	}

	return services
}

func NmapCommandBuilder(format string, trg *Target) (commands []string) {
	var decoys []string
	for _, ip := range trg.Decoys {
		if ip != nil {
			decoys = append(decoys, ip.String())
		}
	}
	newDecoys := strings.Join(decoys, ",")

	// Default fallbacks handled gracefully if no route decoys are passed explicitly
	decoyFlag := ""
	if len(newDecoys) > 0 {
		decoyFlag = "-D " + newDecoys + " "
	}

	switch format {
	case "simple":
		commands = append(commands, fmt.Sprintf("nmap -Pn -sV -T4 --open %s", trg.TargetIp.String()))
	case "fe":
		commands = append(commands, fmt.Sprintf("nmap -Pn -A -sC -f %s%s --open", decoyFlag, trg.TargetIp.String()))
	case "all":
		// Streamlined, high-performance scanning signatures to maximize bounty hunting impact metrics
		c1 := fmt.Sprintf("nmap -Pn -sV --version-light -T4 --open %s", trg.TargetIp.String())
		c2 := fmt.Sprintf("nmap -Pn -sV -sC -F -f %s%s --open", decoyFlag, trg.TargetIp.String())
		c3 := fmt.Sprintf("nmap -p22,80,443,445,3306,8080,8443 -sV -Pn %s --open", trg.TargetIp.String())
		commands = append(commands, c1, c2, c3)
	default:
		commands = append(commands, fmt.Sprintf("nmap -Pn -sV -p- -T4 --open %s", trg.TargetIp.String()))
	}
	return commands
}

var SaveServiceDataTODB = func(name, ip string, data []*Output) error {
	driver, err := db.Old("../../.brain/scans/"+name, 0644)
	if err != nil {
		return err
	}
	for _, datum := range data {
		str := utils.Md5Hash(utils.GenerateUUID())
		if err := driver.Write("servicedata", ip+"_"+str, datum); err != nil {
			utils.Logerror(fmt.Errorf("error saving service data for %s to local storage schema: %v", ip, err))
			continue
		}
	}
	return nil
}


//write data to DB for future reference
var SaveServiceDataTODB = func(name,ip string,data ServiceData)error{
  driver,err := db.Old("../../.brain/scans/" + name,0644)
  if err != nil{
    return err
  }
  for _, datum := range data{
    str := utils.Md5Hash(utils.GenerateUUID())
    if err := driver.Write("servicedata",ip + str,datum); err != nil{
      utils.Logerror(fmt.Errorf("Error saving service data for %s to db.\nERROR: %v",ip,err))
      continue
    }
  }
  return nil
}
