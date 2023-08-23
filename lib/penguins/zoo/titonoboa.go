package zoo

//get-system some process injection
/// use dubois for linux

type Injector interface{
  Inject(pid int, payload []byte) error
}

type ShellCodeRunner struct{
  SC []byte
  Method int
}
