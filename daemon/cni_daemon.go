package main

import(
  "fmt"
  "net"
  "strings"
  "encoding/json"
  "net/http"
  "net/rpc"
  "runtime"
  "github.com/containernetworking/cni/pkg/types"
  "github.com/containernetworking/cni/pkg/version"
  "github.com/containernetworking/cni/pkg/types/current"
  glog "github.com/golang/glog"
  utils "github.com/TCPWAVE/cni-tcpwave/cniutils"
  twc "github.com/TCPWAVE/tims-go-client/twclient"
)

// CNIDaemon driver wrapper
type CNIDaemon struct {
    Drv utils.TCPWaveDriver
}

// NewCNIDaemon creates new instance
func NewCNIDaemon(drv utils.TCPWaveDriver) *CNIDaemon{
  return &CNIDaemon{
    Drv: drv,
  }
}

// Allocate allocates IP
func (cniDem *CNIDaemon) Allocate(args *utils.ExtCmdArgs, result *current.Result) (err error) {
  conf := utils.NetConfig{}
	glog.Infof("allocate: called with args '%s'", *args)
  if err = json.Unmarshal(args.StdinData, &conf); err != nil {
    glog.Errorf("Error parsing netconf: %v", err)
		return fmt.Errorf("Error parsing netconf: %v", err)
	}

	cidr := net.IPNet{IP: conf.IPAM.Subnet.IP, Mask: conf.IPAM.Subnet.Mask}
  gw := conf.IPAM.Gateway
  glog.Infof("Network CIDR: '%s' , Gateway: %s", cidr.String(), gw)
  mac := args.IfMac

  // In Kubernetes to get the container name/hostname
	containerName := ""
	str1 := strings.Split(args.Args, "K8S_POD_NAME=")
	if len(str1) != 1 {
		str2 := strings.Split(str1[1], ";")
		containerName = str2[0]
	}

  ip, err := cniDem.Drv.RequestAddress(conf, cidr, mac, containerName, args.ContainerID)
  if err!=nil{
    glog.Errorf("Error requesting for Ip address : %v", err)
    return fmt.Errorf("Error requesting for Ip address : %v", err)
  }
  glog.Infof("Allocated IP: '%s'", ip)

  ipn, _ := types.ParseCIDR(cidr.String())
	ipn.IP = net.ParseIP(ip)
  iface := 0
	ipConfig := &current.IPConfig{
		Version: "4",
		Address: *ipn,
		Gateway: conf.IPAM.Gateway,
    Interface: &iface,
	}
  netIntrface := &current.Interface{
    Name: args.IfName,
    Mac:mac,
    Sandbox: args.Netns,
  }
  netInterfaces := []*current.Interface{netIntrface}
	routes := convertRoutesToCurrent(conf.IPAM.Routes)
	result.IPs = []*current.IPConfig{ipConfig}
	result.Routes = routes
  result.Interfaces = netInterfaces
	glog.Infof("Allocate result: '%s'", result)
  return err
}

// Release releases IP
func (cniDem *CNIDaemon) Release(args *utils.ExtCmdArgs, reply *struct{}) error {
  conf := utils.NetConfig{}
  // Parse net config.
	glog.Infof("Release: called with args '%s'", *args)
  if err := json.Unmarshal(args.StdinData, &conf); err != nil {
    glog.Errorf("Error parsing netconf: %v", err)
		return fmt.Errorf("Error parsing netconf: %v", err)
	}
  delIP := ""
  // Parse previous result.
  if conf.RawPrevResult != nil {
    resultBytes, err := json.Marshal(conf.RawPrevResult)
    if err != nil {
      return fmt.Errorf("could not serialize prevResult: %v", err)
    }
    res, err := version.NewResult(conf.CNIVersion, resultBytes)
    if err != nil {
      return fmt.Errorf("could not parse prevResult: %v", err)
    }
    conf.RawPrevResult = nil
    conf.PrevResult, err = current.NewResultFromResult(res)
    if err != nil {
      return fmt.Errorf("could not convert result to current version: %v", err)
    }
    delIP = conf.PrevResult.IPs[0].Address.IP.String()
  }
  if delIP == ""{
    glog.Infof("Could not get IP address of container... Deleting container")
    return nil
  }

  glog.Infof("Deleting container  = %s with ip = %s", args.ContainerID, delIP)
  ip, err := cniDem.Drv.ReleaseAddress(conf, delIP, args.IfMac)
	glog.Infof("Address released: '%s'", ip)
  return err
}

func convertRoutesToCurrent(routes []types.Route) []*types.Route {
	var currentRoutes []*types.Route
	for _, r := range routes {
		currentRoutes = append(currentRoutes, &types.Route{
			Dst: r.Dst,
			GW:  r.GW,
		})
	}
	return currentRoutes
}

func createSocketListener(driverSocket *utils.DriverSocket) (net.Listener, error) {
	socketFile := driverSocket.SetupSocket()
	return net.Listen("unix", socketFile)
}

func getTcpwaveDriver(config *utils.Config) (*utils.TCPWaveDriver) {
	hostConfig := config.HostConfig
	objMgr := twc.NewObjectManager(hostConfig, config.KeyFile, config.CertFile,  config.HTTPPoolConnections, config.HTTPRequestTimeout)
	return utils.NewTCPWaveDriver(*objMgr, config.DriverConfig.ClusterName, config.DriverConfig.MaskLength)
}

func startDaemon(conf *utils.Config){
  runtime.LockOSThread()

  glog.Infof("Loaded Config : %v", conf)

  driverSockt := utils.NewDriverSocket(conf.SocketDir, conf.DriverName)
  scktLstnr,err := createSocketListener(driverSockt)
  if err!=nil{
    glog.Errorf("Unable to create socket listener. Error = %v", err)
    return
  }

  tcpDrv := getTcpwaveDriver(conf)
  cni :=NewCNIDaemon(*tcpDrv)
  rpc.Register(cni)
  rpc.HandleHTTP()
  http.Serve(scktLstnr, nil)
  glog.Info("daemon Thread started successfully.")
}

func main(){
  config := utils.LoadConfig()
  glog.Info("Loaded configuration... starting Daemon...")
  startDaemon(config)
}
