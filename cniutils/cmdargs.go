package cniutils

import (
	"github.com/containernetworking/cni/pkg/skel"
)

// ExtCmdArgs extend skel.CmdArgs to include IfMac
type ExtCmdArgs struct {
	skel.CmdArgs
	IfMac string
}
