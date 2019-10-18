package cniutils

import(
  "fmt"
  "os"
  "strings"

  gkg "github.com/onsi/ginkgo"
  gmg "github.com/onsi/gomega"
)

var _ = gkg.Describe("Load Configuration", func(){
  gkg.It("Should return expected config according to command line", func(){
    const(
      Host              = "192.168.10.240"
    	Port              = "7443"
      CertFile          = "/opt/tcpwave/tmp/certs/test.cert"
      KeyFile           = "/opt/tcpwave/tmp/certs/test.key"
    	SslVerify         = "false"
      SocketDir         = "/run/tmp/cni"
    	DriverName        = "twcni"
    	NetworkView       = "default"
    	NetworkContainer  = "172.18.0.0/16"
    	MaskLength        = "26"
    	ClusterName       = "cluster-2"
      LogDir            = "/opt/tcpwave/logs"
    )

    //cmdLine := fmt.Sprintf("twcni-daemon --host=%s --port=%s --username=%s --cert=%s --key=%s --cluster-name=%s --network=%s --mask=%s --socket-dir=%s --driver=%s",
    //        Host, Port, UserName, CertFile, KeyFile, ClusterName, NetworkContainer, MaskLength, SocketDir, DriverName)
    cmdLine := fmt.Sprintf("twcni-daemon --host=%s --port=%s --cert=%s --key=%s --socket-dir=%s --driver=%s",
            Host, Port, CertFile, KeyFile, SocketDir, DriverName)

    os.Args = strings.Fields(cmdLine)

    config := LoadConfig()
    gmg.Expect(config.Host).To(gmg.Equal(Host))
    gmg.Expect(config.Port).To(gmg.Equal(Port))
    gmg.Expect(config.CertFile).To(gmg.Equal(CertFile))
    gmg.Expect(config.KeyFile).To(gmg.Equal(KeyFile))
    gmg.Expect(config.SocketDir).To(gmg.Equal(SocketDir))
  })
})
