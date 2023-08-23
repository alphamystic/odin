package c2

type AdminC2 struct {
  Name string
  Password string
  MSId string
  Address string
  OPort int
  OProtocol string
  ImplantPort int
  ImplantProtocol string
  ImplantTunnel,AdminTunnel string
  Tls bool
  CertPem,KeyPem string
}

type MinionAgent struct {
  Name string
  MuleId string
  C2 bool
  Address string
  Port int
  MotherShip string
  MSId string
  MSession  *Session
}

func CreateMinion(msid,name,msAddress,lport string,port int,session *Session,c2 bool) *MinionAgent{
  var minion MinionAgent
  if c2{
    minion = MinionAgent{
      Name: name,
      MuleId:session.SessionID,
      C2: c2,
      Address: lport,
      Port: port,
      MotherShip: "",
      MSId: msid,
      MSession:session,
    }
  } else {
    minion = MinionAgent{
      Name: name,
      MuleId:session.SessionID,
      C2: c2,
      Address: "0.0.0.0",
      Port: port,
      MotherShip: msAddress,
      MSId: msid,
      MSession:session,
    }
  }
  return &minion
}

type MinionAgents struct {
  Minions []MinionAgent
}

/*
1. Create a template for the mule
2. Add the name and the mothership plus protocol to be used
3. Switch through the file types and generate the appropriate mul
*/

func CreateAdminC2(hash,name,msid,addr,iProtocol,oProtocol,cert,keyCrt string, iPort,oPort int,tls bool) *AdminC2{
  if tls {
    return &AdminC2{
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
    }
  } else {
    return &AdminC2{
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
    }
  }
}
