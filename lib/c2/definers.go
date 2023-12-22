package c2


/*
  * Refactoring this to include a definer to make an a definer for api calls
*/

import(
  "github.com/alphamystic/odin/lib/utils"
  dfn"githuib.com/alphamystic/odin/lib/definers"
)

func CreateMinion(msid,name,msAddress,lport,ops,description,msps,port,cmd string) (*dfn.Minion,error) {
  cmd,err := utils.MultipleToToken(cmd)
  if err != nil {
    return nil,err
  }
  var tt utils.TimeStamps
  tt.Touch
  return &dfn.Minion {
    MinionId: utils.Md5Hash(utils.GenerateUUID()),
    Name: name,
    UName: "",
    UserId: "",
    GroupId: "",
    HomeDir: "",
    Os: ops,
    Description: description,
    Installed: false,
    Address: lport,
    Port: port,
    MotherShipId: msid,
    Motherships: msps,
    LastSeen: tt.CreatedAt,
    GenCommand: cmd,
    tt,
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
      OwnerId: uid,
      Name: name,
      Password: hash,
      MSId: msid,
      Address: addr,
      OPort: oPort,
      OProtocol: oProtocol,
      ImplantPort: iPort,
      ImplantProtocol: iProtocol,
      ImplantTunnel: "",
      Tls: true,
      AdminTunnel: "",
      CertPem: cert,
      KeyPem: keyCrt,
      GenCommand: cmd,
      tt,
    },nil
  } else {
    return &dfn.Mothership {
      OwnerId: uid,
      Name: name,
      Password: hash,
      MSId: msid,
      Address: addr,
      OPort: oPort,
      OProtocol: oProtocol,
      ImplantPort: iPort,
      ImplantProtocol: iProtocol,
      ImplantTunnel: "",
      AdminTunnel: "",
      GenCommand: cmd,
      tt,
    }
  },nil
}
