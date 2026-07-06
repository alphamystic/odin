package builder

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	dfn "github.com/alphamystic/odin/lib/definers"
	"github.com/alphamystic/odin/lib/utils"
)

type Generator interface {
	Generate() error
}

type MinionGenerate struct {
	Ma  *dfn.Minion
	Bob *Builder
}

type MSGenerate struct {
	Ac2 *dfn.Mothership
	Bob *Builder
}

func (m *MinionGenerate) Generate() error {
	muleData := struct {
		ID            string
		MotherShipID  string
		Expiry        string
		Active        bool
		SessionID     string
		Address       string
		TunnelAddress string
		EntryPoint    string
	}{
		ID:            m.Ma.OwnerID, // Aligned to updated Minion field
		MotherShipID:  m.Ma.MothershipID,
		Expiry:        m.Ma.LastSeen, // Aligned to updated Minion field
		Active:        m.Ma.Active,
		SessionID:     m.Ma.MinionID,
		Address:       m.Ma.Address + ":" + m.Ma.Port, // m.Ma.Port is a string; utils.IntToString removed
		TunnelAddress: m.Ma.TunnelAddress,
		EntryPoint:    m.Bob.EntryPoint,
	}

	tmpl := template.Must(template.New("source").Parse(m.Bob.Template))
	tempFile, err := CreateTempFile(m.Bob.Dir, tmpl, muleData)
	if err != nil {
		return fmt.Errorf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if err = m.Bob.BuildGoBinary(tempFile.Name()); err != nil {
		return err
	}

	utils.PrintTextInASpecificColorInBold("white", fmt.Sprintf("Created mule: %s.", m.Ma.Name))
	return nil
}

func (msg *MSGenerate) Generate() error {
	ac2 := struct {
		Name          string
		Password      string
		MSId          string
		Address       string
		OProtocol     string
		ImplantTunnel string
		AdminTunnel   string
		EntryPoint    string
		IAddress      string
		OAddress      string
	}{
		Name:          msg.Ac2.Name,
		Password:      msg.Ac2.Password,
		MSId:          msg.Ac2.MSId,
		Address:       msg.Ac2.Address,
		OProtocol:     msg.Ac2.OProtocol, // Aligned to updated Mothership field
		ImplantTunnel: msg.Ac2.ImplantTunnel,
		AdminTunnel:   msg.Ac2.AdminTunnel,
		EntryPoint:    msg.Bob.EntryPoint,
		IAddress:      msg.Ac2.IAddress,
		OAddress:      msg.Ac2.OAddress,
	}

	tpl := template.Must(template.New("source").Parse(msg.Bob.Template))
	tempFile, err := CreateTempFile(msg.Bob.Dir, tpl, ac2)
	if err != nil {
		return fmt.Errorf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if err = msg.Bob.BuildGoBinary(tempFile.Name()); err != nil {
		return err
	}

	utils.PrintTextInASpecificColorInBold("white", fmt.Sprintf("Created AdminC2: %s.", msg.Ac2.Name))
	return nil
}

func (b *Builder) BuildGoBinary(loc string) error {
	var cmd *exec.Cmd
	currOs := utils.GetCurrentOS()
	if currOs == "windows" {
		cmd = exec.Command("cmd", "/c", b.BuildCommand+" "+loc)
	} else {
		cmd = exec.Command("sh", "-c", b.BuildCommand+" "+loc)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CreateTempFile(dir string, tmpl *template.Template, message interface{}) (*os.File, error) {
	file, err := os.Create(dir + utils.RandString(6) + ".go")
	if err != nil {
		return nil, err
	}

	// Separate deletion from generation to prevent resource exhaustion issues
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, message)
	if err != nil {
		file.Close()
		return nil, err
	}

	_, err = file.Write(buf.Bytes())
	file.Close() // Explicit close prior to running compiler tasks
	if err != nil {
		return nil, err
	}
	return file, nil
}

/*func (msg *MSGenerate) Generate() error{
    tmplChan := make(chan *template.Template)
    errChan := make(chan error)
    go func() {
        tmpl,err := template.New("source").Parse(msg.Bob.Template)
        if err != nil{
            errChan <- err
            return
        }
        tmplChan <- tmpl
    }()
    select {
    case tmpl := <- tmplChan:
        tempFile,err := CreateTempFile(tmpl,msg.Ac2)
        if err != nil{
            return fmt.Errorf("Error creating temporary file: %v",err)
        }
        defer os.Remove(tempFile.Name())
        if err = msg.Bob.BuildGoBinary(); err != nil{
            return err
        }
        utils.PrintTextInASpecificColorInBold("white",fmt.Sprintf("Created AdminC2: %s.",msg.Ac2.Name))
        return nil
    case err := <- errChan:
        return err
    case <- time.After(time.Second * 5):
        return fmt.Errorf("Template parsing timed out")
    }
}
*/
