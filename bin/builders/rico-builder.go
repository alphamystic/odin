/*
func CreateNonGroupedExploits() chan []handlers.Exploit{
  var wg sync.WaitGroup
  for vuln := range vulns {
    wg.Add(1)
    go func(name string) {
      defer wg.Done()
      switch name {
      case  LFI:
        exploit := k.ExploitLFI()
        exploits = append(exploits,exploit)
      case RFI:
        exploit := k.ExploitRFI()
        exploits = append(exploits,exploit)
      case SQLI:
        exploit := k.ExploitSQLI()
        exploits = append(exploits,exploit)
      case NOSQL:
        exploit := k.ExploitNOSQL()
        exploits = append(exploits,exploit)
      case SSRF:
        exploit := k.ExploitSSRF()
        exploits = append(exploits,exploit)
      case IDOR:
        exploit := k.ExploitIDOR()
        exploits = append(exploits,exploit)
      case CSRF:
        exploit := k.ExploitCSRF()
        exploits = append(exploits,exploit)
      case STOREDXSS:
        exploit := k.ExploitSTOREDXSS()
        exploits = append(exploits,exploit)
      case REFLECTEDXSS:
        exploit := k.ExploitREFLECTEDXSS()
        exploits = append(exploits,exploit)
      case DOMXSS:
        exploit := k.ExploitDOMXSS()
        exploits = append(exploits,exploit)
      case HTMLINJECTION:
        exploit := k.ExploitHTMLINJECTION()
        exploits = append(exploits,exploit)
      case BUFFEROVERFLOW:
        exploit := k.ExploitBUFFEROVERFLOW()
        exploits = append(exploits,exploit)
      case INTEGEROVERFLOW:
        exploit := k.ExploitINTEGEROVERFLOW()
        exploits = append(exploits,exploit)
      case STACKOVERFLOW:
        exploit := k.ExploitSTACKOVERFLOW()
        exploits = append(exploits,exploit)
      case HEAPOVERFLOW:
        exploit := k.ExploitHEAPOVERFLOW()
        exploits = append(exploits,exploit)
      default:
        return nil,errorUnKnownVulnerability
      }
    }(vuln.Name)
  }
  wg.Wait()
  // iterate from rico to private and createa work group for each then wait for their completion
  return exploits
}
*/
