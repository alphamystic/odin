package ph

//Are the various many commands executed durig a pentest from during recon to priv escalation
type AttackCommands struct {
  Recon []string // I think now I do remeber. I really don't need this unless it's for internal like it's incide a network
  PrivilegeEscalation []string
  ActiveDirectory []string
  Api []string
}//just realised I forgot why I needed this

// this are attack stages that happen during any network pentesting, they allo you to make a decision on what next
type Mode struct {
  Recon bool
  Pivotting bool
  PrivilegeEscalation bool
  PostExploitation bool
  ActiveDirectory bool
}

func (ac *AttackCommands) LoadCommands(mode *Mode) *AttackCommands{
  return nil
}

func InitAttack() *Mode {
  return &Mode{
    Recon: false,
    Pivotting: false,
    PrivilegeEscalation: false,
    PostExploitation: false,
    ActiveDirectory: false,
  }
}

func (am *Mode) SetMode(AttackMode string){
  switch AttackMode {
  case "recon":
    am.Recon = true
  case "pivot":
    am.Pivotting =  true
  case "privEsc":
    am.PrivilegeEscalation = true
  case "pe":
    am.PostExploitation = false
  case "ad":
    am.ActiveDirectory = true
  default:
    InitAttack()
  }
}
