package cniutils

import(
  gkg "github.com/onsi/ginkgo"
  gmg "github.com/onsi/gomega"
)

var _ = gkg.Describe("Socket Creation", func(){
  gkg.It("New Socket Creation", func(){
    const(
      SocketDir      = "/tmp/cni/run"
      DriverName     = "twcni"
      SocketFileName = "/tmp/cni/run/twcni.sock"
      DefSocketFile  = "/run/cni/twcni.sock"
    )
    defSocket := NewDriverSocket("", DriverName)
    gmg.Expect(defSocket.GetSocketFile()).To(gmg.Equal(DefSocketFile))

    socket := NewDriverSocket(SocketDir, DriverName)
    gmg.Expect(socket.SocketFile).To(gmg.Equal(SocketFileName))
    gmg.Expect(socket.DriverName).To(gmg.Equal(DriverName))
    gmg.Expect(socket.SocketDir).To(gmg.Equal(SocketDir))
    gmg.Expect(socket.GetSocketFile()).To(gmg.Equal(SocketFileName))

    socketFile := socket.SetupSocket()
    gmg.Expect(socketFile).To(gmg.Equal(SocketFileName))
    gmg.Expect(dirExists(SocketDir)).To(gmg.Equal(true))
    gmg.Expect(dirExists(SocketDir)).To(gmg.Equal(true))

  })
})
