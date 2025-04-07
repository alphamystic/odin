package definers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/alphamystic/odin/lib/utils"
)

// Mothership struct
type Mothership struct {
	OwnerID      string `json:"ownerid"`
	Name         string `json:"name"`
	Password     string
	MSId         string `json:"mothershipid"`
	Address      string `json:"address"`
	IAddress string `json:"iaddress"`
  OAddress string `json:"oaddress"`
	ImplantTunnel string `json:"implant_tunnel"`
	AdminTunnel  string `json:"admin_tunnel"`
	Motherships  string `json:"motherships"` // Keep as a bunch of strings, initialize as needed
	Description  string `json:"description"`
	Tls          bool   `json:"tls"`
	CertPem      string `json:"cert_pem"`
	KeyPem       string `json:"key_pem"`
	Active       bool   `json:"active"`
	GenCommand   string `json:"generate_command"`
	Machinedata  string `json:"machine_data"`
	IsOnline bool `json:"isonline"`
	utils.TimeStamps
}

// MaichineData struct
type MaichineData struct {
	UserName  string
	OsType    string
	HomeDir   string
	Password  string
}

// ProtocolType enum
type ProtocolType int

const (
	HTTP ProtocolType = iota
	HTTPS
	GRPC
	SGRPC // Secure gRPC (TLS encrypted)
	UDP
	SOCKS5
)

// Server struct
type Server struct {
	Address  string
	Port     int
	Protocol ProtocolType
	RootPem  []byte // For tunneling, assigned to the client in question
}

// GetServer extracts a Server object from a JWT token.
func GetServer(data string) (*Server, error) {
	token, err := jwt.Parse(data, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.TOKENKEY), nil
	})
	if err != nil {
		return nil, err
	}
	var server Server
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		dataMap, ok := claims["data"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid data format")
		}
		server = Server{
			Address: dataMap["address"].(string),
			Port:    int(dataMap["port"].(float64)), // JSON numbers are float64
			Protocol: ProtocolType(int(dataMap["protocol"].(float64))),
		}
	}
	return &server, nil
}

// PrintServer generates the URL or address of the server based on its protocol.
func PrintServer(s *Server) string {
	switch s.Protocol {
	case HTTP:
		if s.Port == 0 {
			return s.Address
		} else {
			return fmt.Sprintf("http://%s:%d", s.Address, s.Port)
		}
	case HTTPS:
		if s.Port == 0 {
			return s.Address
		} else {
			return fmt.Sprintf("https://%s:%d", s.Address, s.Port)
		}
	case GRPC, SGRPC:
		return fmt.Sprintf("%s:%d", s.Address, s.Port)
	default:
		return "Unknown Protocol"
	}
}
