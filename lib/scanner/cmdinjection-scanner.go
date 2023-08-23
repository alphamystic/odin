package scanner

import (
	"fmt"
	"regexp"
	"strings"
	"os/exec"
	"net/http"
)

// I think we should add a try counter that just encodes the params differently so we swith it up and see default out at some point.
func CommandInjection(input net.URL,client *http.Client) {
	resp, err := http.Get(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	// Scan the page for potential command injection vulnerabilities
	// Look for parameters that could be used for injection
	pattern := `(\?|&)([^=]+)(=|$)`
	re := regexp.MustCompile(pattern)
	parameters := re.FindAllStringSubmatch(input, -1)
	for _, param := range parameters {
		// Check for injection vulnerabilities in each parameter value
		injectionPattern := `(;|--|\|\||&&)[^\s]*`
		injectionRe := regexp.MustCompile(injectionPattern)
		param[2] = strings.ToLower(param[2])
		if injectionRe.MatchString(param[2]) {
			fmt.Println("Possible injection vulnerability detected in parameter", param[2])
			// Mitigate the vulnerability by removing potentially malicious input
			param[2] = injectionRe.ReplaceAllString(param[2], "")
			fmt.Println("Sanitized parameter value:", param[2])
		} else {
			fmt.Println("No injection vulnerabilities detected in parameter", param[2])
		}
	}
	// Execute the command
	output, err := client.Do(req)
	if err != nil {
		utils.Warning(fmt.Sprintf("Error trying out command injection in request: %s",err))
		return
	}
	fmt.Println(string(output))
}
