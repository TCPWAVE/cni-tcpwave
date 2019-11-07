package cniutils

import(
  "flag"
  "net"
  "time"
  "github.com/containernetworking/cni/pkg/types"
  "github.com/containernetworking/cni/pkg/types/current"
  "github.com/TCPWAVE/tims-go-client/twclient"
)

// Constants
const (
	HTTPRequestTimeout  = 60
	HTTPPoolConnections = 10
)

// DriverConfig : Driver specefic configuration
type DriverConfig struct {
	SocketDir          string
	DriverName         string
	DomainName         string
	ContainerNetwork   string  // CIDR
	MaskLength         uint
	ClusterName        string
}

// CustomerConfig holds customer info
type CustomerConfig struct{
  OrganizationName      string
  OrganizationID        string
}

// Config : Overall Configuration
type Config struct {
	twclient.ClientConfig
	DriverConfig
  CustomerConfig
}

// LoadConfig : Read input from execution command or stdin and construct Config Object
func LoadConfig() (config *Config) {
  config = new(Config)
  flag.StringVar(&config.Host, "host", "192.168.0.109", "IP of TCPWave IPAM Host")
	flag.StringVar(&config.Port, "port", "7443", "TCPWave IPAM Port")

	flag.StringVar(&config.CertFile, "cert", "/opt/tcpwave/certs/client.crt", "Client Certificate")
	flag.StringVar(&config.KeyFile, "key", "/opt/tcpwave/certs/client.key", "Client Certificate Key")

  flag.StringVar(&config.SocketDir, "socket-dir", "/run/cni", "Directory where tcpwave IPAM daemon sockets are created")
  flag.StringVar(&config.DriverName, "driver", "tcpwave-cni", "Name of Tcpwave IPAM driver")

  flag.Parse()
  //flag.Lookup("log_dir").Value.Set("/path/to/log/dir")
  config.HTTPRequestTimeout = time.Duration(HTTPRequestTimeout)
  config.HTTPPoolConnections = HTTPPoolConnections
  return config
}

// IPAMConfig : CNI IPAM specefic configuration
type IPAMConfig struct {
	Type             string        `json:"type"`
	SocketDir        string        `json:"socket_dir"`
	ContainerNetwork string        `json:"network_addr"`
  NetworkName      string        `json:"network_name"`
	NetMaskLength    uint          `json:"network_mask"`
	Subnet           types.IPNet   `json:"subnet"`
	Gateway          net.IP        `json:"gateway"`
	Routes           []types.Route `json:"routes"`
  Domain           string        `json:"primary_domain"`
  Org              string        `json:"organization_name"`
}

// NetConfig : Complete CNI Configuration
type NetConfig struct {
	Name             string                      `json:"name"`
	Type             string                      `json:"type"`
	Bridge           string                      `json:"bridge"`
	IsGateway        bool                        `json:"is_gateway"`
	IPAM             *IPAMConfig                 `json:"ipam"`
  RawPrevResult    *map[string]interface{}     `json:"prevResult"`
	PrevResult       *current.Result             `json:"-"`
  CNIVersion       string                      `json:"cniVersion"`
}
