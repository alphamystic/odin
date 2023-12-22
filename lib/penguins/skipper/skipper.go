package skipper

/*
  * Should be able to chain vulnerabilities and use them to completely compromise a target
*/

//or an input of attack commands that use an exploit
import (
  "fmt"
  "time"
  //"sync"

  "github.com/alphamystic/odin/lib/db"
  "github.com/alphamystic/odin/lib/utils"
  "github.com/alphamystic/odin/lib/handlers"

  "github.com/alphamystic/odin/lib/penguins/ph"
  "github.com/alphamystic/odin/lib/penguins/kowalski"
)

// this should handle pivoting post explloitation basically red and purple teaming functionalities
type Skipper struct {
  Exploits []handlers.Exploit
  Name string
}

//on receiving the exploit, create a mthership with the scan name and generate minions,droppers for it
func (s *Skipper) Attack(targets []string){
  mode :=  ph.InitAttack()
  //var ac AttackCommands
  mode.Recon = true
  //reconCommands := ac.LoadCommands(mode)
  kwsk := &kowalski.KOWALSKI{
    Targets: targets,
    Name: s.Name,
    ScanID: utils.GenerateUUID(),
  }
  exploits := make(chan *handlers.Exploit)
  exploitsDone := make(chan bool)
  go kwsk.Kowalski_Analysis(exploits,exploitsDone)
  //ltg := len(targets)
  //time.Sleep(ltg * time.Minute)// change to a more ideal time
  time.Sleep(100 * time.Millisecond)
  /*if 0 >= len(exploits) {
    utils.PrintTextInASpecificColorInBold("cyan","*******************************************************************************************************************************")
    utils.PrintTextInASpecificColorInBold("white","Sorry ........ No exploit was found for your targets. Try finding another plugin for vulnerability scanning.")
    utils.PrintTextInASpecificColorInBold("white","   OR IT'S JUST SECURE............")
    utils.PrintTextInASpecificColorInBold("cyan","*******************************************************************************************************************************")
    fmt.Println("")
    return
  }*/

  //create a mothership here
  for exploit,ok := <-exploits; ok; exploit,ok = <-exploits {
    if !ok {
      <- exploitsDone
      close(exploits)
    }
    if err := s.SaveExploit(exploit); err != nil {
      utils.Danger(fmt.Errorf("[-]  SKIPPER:  Error saving exploit: %s", err))
      // os.Exit(1) do something else other than exiting app, cache the exploit or something
    }
    if !exploit.Works{
      utils.PrintTextInASpecificColorInBold("yellow","**********************************************************************")
      utils.PrintTextInASpecificColorInBold("blue","      Zero Working exploits were found for:      ")
      utils.PrintTextInASpecificColorInBold("blue",fmt.Sprintf("          Host:   %s",exploit.Trg.Host))
      utils.PrintTextInASpecificColorInBold("blue",fmt.Sprintf("          Host IP Addres:   %s",exploit.Trg.HostIp))
      utils.PrintTextInASpecificColorInBold("blue",fmt.Sprintf("          Target IP Address:   %s",exploit.Trg.TargetIp))
      utils.PrintTextInASpecificColorInBold("yellow","**********************************************************************")
    } else {
      utils.PrintTextInASpecificColorInBold("white",fmt.Sprintf("Popping a shell for %s from %s",exploit.Trg.TargetIp,exploit.Trg.Host))
    }
  }
}

// Save the exploits to a file (port this to db)
func (s *Skipper) SaveExploit(exploit *handlers.Exploit) error {
  driver,err := db.Old("../../.brain/scans/" + s.Name,0644)
  if err != nil{
    return err
  }
  if err := driver.Write("exploits",exploit.Trg.TargetIp.String(),exploit); err != nil{
    return fmt.Errorf("Error saving exploit for %s to db.\nERROR: %v",exploit.Trg.TargetIp.String(),err)
  }
  return nil
}

/*
func ProcessOutput(mode Mode,command string, output []byte) (outputs []CommandOutput, features *mat.VecDense, labels *mat.VecDense, err error) {
  // Initialize the features and labels slices
	var featuresData []float64
	var labelsData []float64
	var next string
  //when done handling and can warrant mooving to the next level, then change the current mode to false and set the next one totrue
  if mode.recon {
    //handle recon
    ReconHandler()
  }
  if mode.pivotting {
    // handle pivotting
  }
  if mode.privilegeEscalation {
    //handle privilegeEscalation
  }
  if mode.postExploitation {
    //handle postExploitation
  }
  if mode.activeDirectory {
    // handle activeDirectory
  }
}
*/
