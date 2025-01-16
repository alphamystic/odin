package cli

import (
  "fmt"
  "time"
  "bytes"
  "strings"
  "io/ioutil"
  "github.com/alphamystic/odin/lib/utils"
  "github.com/spf13/cobra"
  "github.com/alphamystic/odin/lib/penguins/skipper"
)


var cmdCreateTarget = &cobra.Command{
  Use:   "attack",
  Short: "Specify a specific target to attack",
  Long:  "Specify a specific target to attack using the -t flag",
  Run: func(cmd *cobra.Command, args []string) {
    // Retrieve the target flag
    target, err := cmd.Flags().GetString("t")
    if err != nil {
      utils.Warning(fmt.Sprintf("Error retrieving target flag: %v", err))
      return
    }
    // Check if the target flag is set or has default value
    if target == "" || target == "domain.com" {
      utils.Warning("Target not provided or default value used. Use --t to specify a target.")
      return
    }
    // Retrieve the name flag
    name, _ := cmd.Flags().GetString("scansname")
    if !utils.CheckifStringIsEmpty(name) {
      utils.Warning(fmt.Sprintf("Name of scan cannot be empty, use --name \"name_of_scan\": %s",name))
      return
    }
    // Proceed with the scan
    utils.Warning(fmt.Sprintf("Target is %s", target))
    targets := []string{target}
    fmt.Println(targets)
    skp := &skipper.Skipper{
      Name: name,
    }
    t0 := time.Now()
    utils.PrintTextInASpecificColorInBold("white", fmt.Sprintf("Starting scan %s at %s", name, t0.String()))
    skp.Attack(targets)
    t1 := time.Now()
    utils.PrintTextInASpecificColorInBold("white", fmt.Sprintf("The scan %s took %v to run.\n", name, t1.Sub(t0)))
    /*ticker := time.NewTicker(time.Second)
    done := make(chan bool)
    var YES = func(){
      skp.Attack(targets)
      done <- true
    }
    go func(){
      for {
        select{
        case <- done:
          return
        case <- ticker.C:
          YES()
        }
      }
    }()*/
  },
}

var cmdCreateTargets = &cobra.Command {
  Use: "targets",
  Short: "Specify a specfic target to attack",
  Long: "Specify a specfic target to attack",
  Run:func(cmd *cobra.Command,args []string){
    var targets []string
    var err error
    tFile,_ := cmd.Flags().GetString("tF")
    list,_ := cmd.Flags().GetString("tL")
    name,_ :=cmd.Flags().GetString("name")
    if utils.CheckifStringIsEmpty(name) {
      utils.Warning("Name of scan can not be empty.");return
    }
    if !utils.CheckifStringIsEmpty(tFile){
      //open the file get the target list
      targets,err = GetTargetsFromFile(tFile)
      if err != nil{
        utils.CustomError("Error getting targets from file: ",err)
        return
      }
    } else {
      if list != ""{
        targets = ExplodeTargets(list)
      }
    }
    fmt.Println(targets)
    skp := &skipper.Skipper{
      Name: name,
    }
    skp.Attack(targets)
  },
}

func ExplodeTargets(data string)([]string){
  var targets []string
  trgs := strings.Split(data,",")
  for _,trg := range trgs{
    targets = append(targets,trg)
  }
  return targets
}

func GetTargetsFromFile(fileName string)([]string,error){
  var targets []string
  buf,err := ioutil.ReadFile(fileName)
  if err != nil {
    return nil,err
  }
  lines := bytes.Split(buf,[]byte("\n"))
  for _,line := range lines {
    targets = append(targets,string(line))
  }
  return targets,nil
}

func init(){
  cmdCreateTargets.Flags().String("tF","tf","target-filename")
  cmdCreateTargets.Flags().String("tL","tL","A list of targets")
  cmdCreateTarget.Flags().String("t", "domain.com", "A specific target to attack (e.g., `attack --t domain.com --name \"ExampleScan\"`)")
  cmdCreateTargets.Flags().String("name","target","Name for this particular scan say target_one if it's scan for target_one systems")
  cmdCreateTarget.Flags().String("scansname","target","Name for this particular scan say target_one if it's scan for target_one systems")
  cmdCreateTarget.MarkFlagRequired("scansname")
  RootCmd.AddCommand(cmdCreateTarget)
  RootCmd.AddCommand(cmdCreateTargets)
}

/*utils.PrintTextInASpecificColorInBold("white","When you get money, realize that it's not enough but you have to make more save and yet again say thank you......")
utils.PrintTextInASpecificColorInBold("white","Success can only be replicated.....")
*/
