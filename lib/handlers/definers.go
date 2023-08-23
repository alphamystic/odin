package handlers

// Vulnerability represents a type of vulnerability.
type Vulnerability int

//you can use a url shortener to hide payload @Random thoughts
const (
	// LFI represents local file inclusion vulnerability.
	LFI Vulnerability = iota
	// RFI represents remote file inclusion vulnerability.
	RFI
	// SQLI represents SQL injection vulnerability.
	SQLI
	// COMMANDINJECTION represents a grouped type of command injection
	COMMANDINJECTION
  // NOSQL Represents an injection in a nosql database
  NOSQL
	// SSRF represents server-side request forgery vulnerability.
	SSRF	// http://ctf.securityboat.in:81/
	// IDOR represents insecure direct object reference vulnerability.
	IDOR
	// CSRF represents cross-site request forgery vulnerability.
	CSRF
  // STOREDXSS represents stored xss type of  cross-site scripting vulnerability.
	STOREDXSS
  // REFLECTEDXSS represents reflected xss type of  cross-site scripting vulnerability.
  REFLECTEDXSS
  // DOMXSS represents dom xss type of  cross-site scripting vulnerability.
  DOMXSS
  // HTMLINJECTION represents a HTMLINJECTION type of vulnerability
  HTMLINJECTION
  // BUFFEROVERFLOW represents a buffer overflow
  BUFFEROVERFLOW
  // INTEGEROVERFLOW represents an integer overflow
  INTEGEROVERFLOW
  // STACKOVERFLOW represents a stack based overflow
  STACKOVERFLOW
  // HEAPOVERFLOW represents a heap based overflow
  HEAPOVERFLOW
	// RACECONDITION represents a racecondition type of vulnerability
	RACECONDITION
	// RCE represents a remote code execution type of vulnerability
	// All vulnerabilities should either lead to this or chained to lead to RCE
	RCE
  // when printed it represents a new attack vector
  UNKNOWN
)

func (v Vulnerability) GetVulnerability() string {
	switch v {
	case LFI:
		return "Local File Inclusion"
	case RFI:
		return "Remote File Inclusion"
	case SQLI:
		return "SQL Injection"
	case COMMANDINJECTION:
		return "Command Injection"
	case NOSQL:
		return "NoSQL Injection"
	case SSRF:
		return "Server-Side Request Forgery"
	case IDOR:
		return "Insecure Direct Object Reference"
	case CSRF:
		return "Cross-Site Request Forgery"
	case STOREDXSS:
	  return "Stored Xss"
	case REFLECTEDXSS:
	  return "Reflected Xss"
	case DOMXSS:
	  return "DOM Xss"
	case HTMLINJECTION:
	  return "Html Injection"
	case BUFFEROVERFLOW:
	  return "Buffer Overflow"
	case INTEGEROVERFLOW:
	  return "Integer Overflow"
	case STACKOVERFLOW:
	  return "Stack OverFlow"
	case HEAPOVERFLOW:
		return "Heap Overflow"
	case RACECONDITION:
		return "Run Race"
	case RCE:
	  return "Remote Code Execution"
	default:
		return "Unknown Vulnerability"
	}
	return ""
}
/*
vuln := handlers.SQLI
fmt.Println(vuln.String()) // Output: SQL Injection
*/
// Represents a type of attack that can be used to compromise a target
// they can all be chained to vreate an attack style
type AttackType int

const (
  // PHISHING represents an attack where a target is compromised through phishing
  PHISHING AttackType = iota
  // WEBATTACK represents an attack where a target is compromised through exploiting a web app vulnerability
  WEBATTACK
  // BRUTEFORCE represents an attack where a target is compromised through bruteforcing a given service
  BRUTEFORCE
  // ZERODAY represents an attack where a target is compromised through a zero day
  ZERODAY
)

func (at AttackType) GetAttackType() string{
	switch at {
		case PHISHING:
			return "PHISHING"
		case WEBATTACK:
			 return "WEBATTACK"
		case BRUTEFORCE:
			return "BRUTEFORCE"
		case ZERODAY:
			return "ZERODAY"
		default:
			return ""
	}
	return ""
}

type Vulnerabilities struct {
	Trg *Target
  Name Vulnerability
	Severity int
	Target string
  Payload string
  AT AttackType
	Grouped bool
  Authenticated bool
	Works bool
  Details string // you can plug in a client here to be used for http or login creds
}

type Exploit struct {
	Trg *Target
	LHOST string
	LPORT int
	Address string
	Target string
	AverageSeverity int
	Grouped bool
	Vulns []*Vulnerabilities
	Works bool
}

// implemented by private and rico
/*
type Scanner interface{
  Scan(reconDataChannel chan *ReconData) []Vulnerabilities// should return a channel of this
}
*/
type Scanner interface{
	Scan(reconData ReconData) []Vulnerabilities
}
