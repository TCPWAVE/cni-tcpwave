package cniutils

import(
  "net"
  "strconv"
  "strings"
  "github.com/containernetworking/plugins/pkg/utils/hwaddr"
  glog "github.com/golang/glog"
  twc "github.com/TCPWAVE/tims-go-client/twclient"
)

// TCPWaveDriver tcpwave driver
type TCPWaveDriver struct {
	ObjMgr           twc.ObjectManager
	DefaultCluster   string
	DefaultMaskLen   uint
}

// NewTCPWaveDriver creates new TCPWave Driver
func NewTCPWaveDriver(objMan twc.ObjectManager, defCluster string, defMask uint) (*TCPWaveDriver){
  drv := &TCPWaveDriver{
    ObjMgr: objMan,
    DefaultCluster: defCluster,
    DefaultMaskLen: defMask,
  }
  return drv
}

// ReleaseAddress deletes container address
func (twd *TCPWaveDriver) ReleaseAddress(conf NetConfig, ip string, mac string) (string,error){
  glog.Infof("Ip delete request with ip = %s", ip)
  err := twd.ObjMgr.DeleteIPAddress(ip, "", conf.IPAM.Org)
  if err!=nil{
    glog.Error(err)
    return "", err
  }
  return ip,nil
}

// RequestAddress fetch next free ip from ipam and creates the received object
func (twd *TCPWaveDriver) RequestAddress(config NetConfig, subnetAddr net.IPNet, macAddr string,
  containerName string, containerID string) (string,error){

  // Create network
  var network *twc.Network
  var err error
  networkAddress := strings.Split(config.IPAM.ContainerNetwork, "/")[0]
  network,err = twd.ObjMgr.GetNetwork(config.IPAM.ContainerNetwork, config.IPAM.Org)
  if err!=nil{
    glog.Infof("Creating Network with address : %s", config.IPAM.ContainerNetwork)
    network = &twc.Network{}
    network.Name = config.IPAM.NetworkName
    network.Description = "Kubernetes Network"
    network.Organization = config.IPAM.Org
    addrBits := strings.Split(networkAddress, ".")
    addr1,_ := strconv.Atoi(addrBits[0])
    network.Addr1 = addr1
    addr2,_ := strconv.Atoi(addrBits[1])
    network.Addr2 = addr2
    addr3,_ := strconv.Atoi(addrBits[2])
    network.Addr3 = addr3
    addr4,_ := strconv.Atoi(addrBits[3])
    network.Addr4 = addr4
    network.DMZVisible = "no"
    network.MaskLen = int(config.IPAM.NetMaskLength)
    _,err1 := twd.ObjMgr.CreateNetwork(*network)
    if err1!=nil{
      return "", err1
    }
  }
  // Create Subnet
  var subnet *twc.Subnet
  subnet,err = twd.ObjMgr.GetSubnet(subnetAddr.String(), config.IPAM.Org)
  if err!=nil {
    glog.Infof("Creating Subnet with address : %s", subnetAddr.String())
    subnet = &twc.Subnet{MaskLen: 26}
    subnet.Name = "K8S Subnet"
    subnet.Description = "subnet for kubernetes"
    subnet.Organization = config.IPAM.Org
    subNtAddr := strings.Split(subnetAddr.String(), "/")[0]
    addrBits := strings.Split(subNtAddr, ".")
    glog.Info("Address Bits Array : " + addrBits[3])
    addr1,_ := strconv.Atoi(addrBits[0])
    subnet.Addr1 = addr1
    addr2,_ := strconv.Atoi(addrBits[1])
    subnet.Addr2 = addr2
    addr3,_ := strconv.Atoi(addrBits[2])
    subnet.Addr3 = addr3
    addr4,_ := strconv.Atoi(addrBits[3])
    subnet.Addr4 = addr4
    subnet.RouterAddr = addrBits[0] + "." + addrBits[1] + "." + addrBits[2] + "." + strconv.Itoa(addr4 + 1)
    subnet.NetworkAddr = networkAddress
    subnet.PrimaryDomain = config.IPAM.Domain
    _,err1 := twd.ObjMgr.CreateSubnet(*subnet)
    if err1!=nil{
      return "", err1
    }
  }

  // Fetch available IP from IPAM
  ip,err2 := twd.ObjMgr.GetNextFreeIP(subnetAddr.String(), config.IPAM.Org)
  if err2!=nil{
    return "",err2
  }

  mac := macAddr
  glog.Infof("Free Ip received from IPAM = %s , mac addr = %s", ip, mac)

  if config.Type == "bridge" {
		hwAddr, err := hwaddr.GenerateHardwareAddr4(net.ParseIP(ip), hwaddr.PrivateMACPrefix)
    glog.Infof("Computed Mac addr for bridge type: %s", hwAddr.String())
    mac = hwAddr.String()
		if err != nil {
			glog.Errorf("Problem while generating hardware address using ip: %s", err)
			return "", err
		}
	}
  _, err = twd.ObjMgr.CreateIPAddress(ip, mac, subnetAddr.IP.String(), config.IPAM.Domain, config.IPAM.Org, containerName)
  if err!=nil{
    return "", err
  }
  glog.Infof("Ip Created in IPAM = %s", ip)
  return ip, nil
}
