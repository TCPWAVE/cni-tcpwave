package main

import(
  "os"
  "fmt"
  "log"
  "encoding/json"
  "net/rpc"
  "path/filepath"
  "github.com/containernetworking/cni/pkg/skel"
  "github.com/containernetworking/cni/pkg/types"
  "github.com/containernetworking/cni/pkg/version"
  "github.com/containernetworking/cni/pkg/types/current"

  utils "github.com/TCPWAVE/cni-tcpwave/cniutils"
)

var(
  logger    *log.Logger
)

// constants
const(
  DefaultLogDir = "/opt/tcpwave/logs"
)

func initLogger(){
  // Create Log Dir
  dirInfo, err := os.Stat(DefaultLogDir)
  if err == nil{
    if! dirInfo.IsDir(){
      os.MkdirAll(DefaultLogDir, 0700)
    }
  }else if os.IsNotExist(err){
    os.MkdirAll(DefaultLogDir, 0700)
  }

  // Create Log file
  file, err := os.OpenFile(DefaultLogDir+"/twcni.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
  if err != nil {
      log.Fatalln("Failed to open log file : ", err)
  }
  logger = log.New(file,"[INFO]: ",log.Ldate|log.Ltime|log.Lshortfile)
}

func runPlugin() {
	skel.PluginMain(cmdAdd, cmdCheck, cmdDel, version.All, "TcpWave CNI Plugin")
}

func cmdAdd(args *skel.CmdArgs) error {
  logger.Printf("Received Add Commmand with args : %v", args)
  versionDecoder := &version.ConfigDecoder{}
	confVersion, err := versionDecoder.Decode(args.StdinData)
	if err != nil {
		return err
	}

	result := &current.Result{}
	extArgs := &utils.ExtCmdArgs{CmdArgs: *args}
	mac:= GetMacAddress(args.Netns, args.IfName, logger)
	extArgs.IfMac = mac
  logger.Println("Calling CNIDaemon.Allocate")
	if err := rpcCall("CNIDaemon.Allocate", extArgs, result); err != nil {
		return err
	}
  logger.Printf("Command : Add, Result : %v, version : %v", result, confVersion)
	return types.PrintResult(result, confVersion)
}

func cmdCheck(args *skel.CmdArgs) error {
  logger.Printf("Received Check Commmand with args : %v", args)
  versionDecoder := &version.ConfigDecoder{}
	confVersion, err := versionDecoder.Decode(args.StdinData)
	if err != nil {
		return err
	}
	result := &current.Result{}
  logger.Printf("Command : Check, Result : %v, version : %v", result, confVersion)
  return types.PrintResult(result, confVersion)
}

func cmdDel(args *skel.CmdArgs) error {
  logger.Printf("Received Del Commmand with args : %v", args)
  result := struct{}{}
	extArgs := &utils.ExtCmdArgs{CmdArgs: *args}

	mac:= GetMacAddress(args.Netns, args.IfName, logger)
	extArgs.IfMac = mac
  logger.Printf("Deleting container with id = %v", args.ContainerID)
  logger.Println("Calling CNIDaemon.Release")
	if err := rpcCall("CNIDaemon.Release", extArgs, &result); err != nil {
		return fmt.Errorf("error dialing CNIDaemon daemon: %v", err)
	}
  logger.Printf("Command : Del, Result : %v", result)
	return nil
}

func rpcCall(method string, args *utils.ExtCmdArgs, result interface{}) error {
  conf := utils.NetConfig{}
	if err := json.Unmarshal(args.StdinData, &conf); err != nil {
		return fmt.Errorf("Error parsing netconf: %v", err)
	}

  client, err := rpc.DialHTTP("unix", utils.NewDriverSocket(conf.IPAM.SocketDir, conf.IPAM.Type).GetSocketFile())
	if err != nil {
		return fmt.Errorf("Error dialing CNI Daemon: %v", err)
	}

  // The daemon may be running under a different working dir
	// so make sure the netns path is absolute.
	netns, err := filepath.Abs(args.Netns)
	if err != nil {
		return fmt.Errorf("Failed to make %q an absolute path: %v", args.Netns, err)
	}
	args.Netns = netns
  err = client.Call(method, args, result)
	if err != nil {
		return fmt.Errorf("Error calling %v: %v", method, err)
	}
  return nil
}

func main(){
  initLogger()
  logger.Println("CNI Plugin logger initialized")
  logger.Println("Starting Plugin")
  runPlugin()
}
