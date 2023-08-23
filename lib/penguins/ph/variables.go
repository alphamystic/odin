package ph

import (
  "errors"
)
// RunningMode refers to ways to ensure constatnt learning
type RunningMode int

const (
  // PENTESTING used to set odin to a pentesting mode, proper documentaion and full scan and enumeration
  PENTESTING RunningMode = iota
  // CTF is used to train odin on vulneralibilities and train it on known unknowns
  CTF
  // RESEARCH is used not nescescarilly to createa new neural net but to find other possible ways
  RESEARCH
  // BUGHUNTING is used in bug hunting, a memory can be stored to be used at a later stage
  BUGHUNTING
  // used to train the predictor on particular vulnerabilities
  TRAINNING
)

var (
  ErrorUnKnownVulnerability = errors.New("A New UNKNONWN Vulnerability has been found")
  RicoMode = false
  runnningMode RunningMode
  pentesting bool
  ctf bool
  research bool
  bughutning bool
  training bool
)

func SetRunningMode(rm RunningMode){
  switch rm {
    case PENTESTING:
      pentesting = true
    case CTF:
      ctf = true
    case BUGHUNTING:
      bughutning = true
    case TRAINNING:
      training =true
    default:
      research = true
  }
}

/*
rm := penguins.RunningMode
rm.SetRunningMode(PENTESTING)
*/
