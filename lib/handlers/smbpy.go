package handlers

import (
    //"fmt"
    "github.com/stacktitan/smb/smb"
)

type SmbConfig struct {
    TargetIp   string
    TargetPort int
    Username   string
    Password   string
    Domain     string
}

type SmbClient interface {
    Connect(config *SmbConfig) error
    Disconnect() error
    ExecuteCommand(cmd string) ([]byte, error)
}

type SmbPyClient struct {
    SmbSession *smb.Session
}

func (client *SmbPyClient) Connect(config *SmbConfig) error {
    options := smb.Options{
        Host:     config.TargetIp,
        Port:     config.TargetPort,
        User:     config.Username,
        Password: config.Password,
        Domain:   config.Domain,
    }
    session, err := smb.NewSession(options, false)
    if err != nil {
        return err
    }
    client.SmbSession = session
    return nil
}

func (client *SmbPyClient) Disconnect() error {
    if client.SmbSession == nil {
        return nil
    }
    client.SmbSession.Close()
    return nil
}
/*
func (client *SmbPyClient) ExecuteCommand(cmd string) ([]byte, error) {
    if client.SmbSession == nil {
        return nil, fmt.Errorf("No session present or not connected")
    }
    err := client.SmbSession.Mount("/")
    if err != nil {
        return nil, err
    }
    defer client.SmbSession.Unmount("/")
    _, err = client.SmbSession.Exec(fmt.Sprintf("echo '%s'", cmd))
    if err != nil {
        return nil, err
    }
    out, err := client.SmbSession.ReadFile("/.out")
    if err != nil {
        return nil, err
    }
    return out, nil
}

func EnumSMB() {
    config := &SmbConfig{
        TargetIp:   "192.168.0.1",
        TargetPort: 445,
        Username:   "user",
        Password:   "password",
        Domain:     "domain",
    }
    client := &SmbPyClient{}
    err := client.Connect(config)
    if err != nil {
        panic(err)
    }
    defer client.Disconnect()

    output, err := client.ExecuteCommand("ls")
    if err != nil {
        panic(err)
    }

    fmt.Println(string(output))
}
*/
