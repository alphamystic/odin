package zoo

//get-system some process injection
/// use dubois for linux

import (
  "time"
)

type Injector interface{
  Inject(pid int, payload []byte) error
}

type ShellCodeRunner struct{
  SC []byte
  Method int
}


type Lala interface {
  Tunya()
}

type Dozz struct{ S int}

func (d *Dozz) Tunya() {
  <-time.After(time.Duration(d.S) * time.Second)
}
