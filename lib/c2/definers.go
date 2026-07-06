package c2


/*
  * Refactoring this to include a definer to make an a definer for api calls
*/

import(
  "github.com/alphamystic/odin/lib/utils"
  dfn"github.com/alphamystic/odin/lib/definers"
)

func CreateMinion(msid,name,msAddress,lport,ops,description,msps,port,cmd string) (*dfn.Minion,error) {
  cmd,err := utils.MultipleToToken(cmd)
  if err != nil {
    return nil,err
  }
  var tt utils.TimeStamps
  tt.Touch()
  return &dfn.Minion {
    MinionID: utils.Md5Hash(utils.GenerateUUID()),
    Name: name,
    UName: "",
    UserID: "",
    GroupID: "",
    HomeDir: "",
    Os: ops,
    Description: description,
    Installed: false,
    MothershipID: msid,
    Address: "127.0.0.1",
    Port: port,
    Motherships: msps,
    TunnelAddress: "",
    Tls: false,
    Active: true,
    OwnerID: "",
    IsDropper: true,
    LastSeen: tt.CreatedAt.String(),
    GenCommand: cmd,
    TimeStamps: tt,
  },nil
}

func CreateMothership(hash,name,msid,addr,iProtocol,oProtocol,cert,keyCrt,cmd string, iPort,oPort int,tls bool) (*dfn.Mothership,error){
  var tt utils.TimeStamps
  cmd,err := utils.MultipleToToken(cmd)
  if err != nil {
    return nil,err
  }
  uid := utils.Md5Hash(utils.GenerateUUID())
  tt.Touch()
  if tls {
    return &dfn.Mothership {
      OwnerID: uid,
      Name: name,
      Password: hash,
      MSId: msid,
      Address: addr,
      IAddress: "",
      OAddress: "",
      ImplantTunnel: "",
      AdminTunnel: "",
      Motherships: "",
      Description: "",
      Tls: true,
      CertPem: cert,
      KeyPem: keyCrt,
      Active: true,
      GenCommand: cmd,
      Machinedata: "",
      IsOnline: true,
      TimeStamps: tt,
    },nil
  } else {
    return &dfn.Mothership {
       OwnerID: uid,
        Name: name,
        Password: hash,
        MSId: msid,
        Address: addr,
        IAddress: "",
        OAddress: "",
        ImplantTunnel: "",
        AdminTunnel: "",
        Motherships: "",
        Description: "",
        Tls: false,
        CertPem: cert,
        KeyPem: keyCrt,
        Active: true,
        GenCommand: cmd,
        Machinedata: "",
        IsOnline: true,
        TimeStamps: tt,
    }, nil
  }
}
