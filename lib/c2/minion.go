package c2

// this package is deprecated, use the minions on definers
import(
  "net"
  "time"
)

type Minion struct{
  MinionId string
  Name string
  UName string
  UserId string
  GroupId string
  HomeDir string
  Os  string
  Description string
  Installed bool
  MotherShipId string
  LastSeen string
  CreatedAt string
  UpdatedAt string
}

var now = time.Now()
var currentTime = now.Format("2006-01-02 15:04:05")

func Spear()*Minion{
  return &Minion{
    CreatedAt: currentTime,
    UpdatedAt: currentTime,
  }
}

func (m *Minion) Populate()error{
  return nil
}

// we don't have to register a mutant just use droppers instead
type Mutant struct {
  MutantId string
  Name string
  MotherShipId string
  EXEMinion string
  ELFMinion string
  MACMinion string
  MothersipAddresses []*net.IP // try masking this with net.Mask to avoid/reduce detection
  CreatedAt string
}
