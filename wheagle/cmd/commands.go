package cmd

import (
  "fmt"
)


var cmdAdminCli = &cobra.Command {
  Use: "im",
  Long: "Interact with a admin server",
	Run: func(cmd *cobra.Command, args []string){
    id,err := cmd.Flags().GetString("id")
    if err != nil {
      utils.Logerror(fmt.Errorf("Honey-badger id can not be nill"));return
    }
    pass,err := cmd.Flags().GetString("pass")
    if err != nil {
      utils.Logerror(fmt.Errorf(" Password for Honey-badger can not be nill"));return
    }
    var (
      err error
      url string
      client *http.Client
      requestBody bytes.Buffer
      encoder *gob.Encoder
      decoder *gob.Decoder
      work *core.Work
    )
    work = new(core.Work)
    /// @TODO initialize connection (tcp or tls)
    con,err := c2.GetConn(id)
    if err != nil{
      utils.Logerror(err);return
    }
    client = new(http.Client)
    if con.Tls {
      client = &http.CLient{
  			Transport: &http.Transport{
  				TLSClientConfig: &tls.Config{
  					InsecureSkipVerify: true,
  				},
  			},
  		}
      url = fmt.Sprintf("https://%s",con.Address)
    } else {
      client = &http.Client{}
      url = fmt.Sprintf("http://%s",con.Address)
    }
    work.UserId = pass
    work.OperatorId = id
    work.ForMS = true
    go func(){
      encoder = gob.NewEncoder(&requestBody)
      decoder = gob.NewDecoder(&requestBody)
    }()
    if err = encioder.Encode(work);err != nil{
      utils.Logerror(fmt.Errorf("Error encoding registration.\nERROR: %q",err));return
    }
    req,err := http.NewRequest("PUT",url,&requestBody)
    if err != nil{
      utils.Logerror(err);return
    }
    //utils.ClientAddHeaderVal(req,"Register","")
    utils.Interactor(id,false)
		fmt.Println("")
    var iarg string
		reader := bufio.NewReader(os.Stdin)
	  for {
      START:
      switch iarg {
        case "shell":
          conn,err net.Dial("tcp",address)
          if err != nil {
            utils.Logerror(err)
          }
        }
      END:
    }
  },
}

func AdminErrorChecker = func(res http.Response)error{
  val := res.Header.Get("ERROR")
  if len(Val) > 0 {
    return errors.New(val)
  }
  return nil
}
