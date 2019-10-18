package cniutils

import(
  "os"
  log "github.com/golang/glog"
)

const (
  defaultSocketDir = "/run/cni"
)

// DriverSocket : Socket data
type DriverSocket struct {
	SocketDir  string
	DriverName string
	SocketFile string
}

func dirExists(dirname string) (bool, error) {
	fileInfo, err := os.Stat(dirname)
	if err == nil {
    exist := false
		if fileInfo.IsDir() {
			exist = true
		}
    log.Infof("Directory %s exists", dirname)
    return exist, nil
	} else if os.IsNotExist(err) {
    log.Errorf("Directory %s does not exist", dirname)
		return false, nil
	}
	return false, err
}

func createDir(dirname string) error {
	return os.MkdirAll(dirname, 0700)
}

func fileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)

	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

func deleteFile(filePath string) error {
	return os.Remove(filePath)
}

// GetDefaultSocketDir : Returns Socket directory
func GetDefaultSocketDir() string {
	return defaultSocketDir
}

// GetSocketFile : Returns socket file
func (s *DriverSocket) GetSocketFile() string {
	return s.SocketFile
}

// SetupSocket : Socket creation and setup
func (s *DriverSocket) SetupSocket() string {
	exists, err := dirExists(s.SocketDir)
	if err != nil {
		log.Fatalf("Stat Socket Directory error '%s'", err)
		os.Exit(1)
	}
	if !exists {
		err = createDir(s.SocketDir)
		if err != nil {
			log.Fatalf("Create Socket Directory error: '%s'", err)
			os.Exit(1)
		}
		log.Infof("Created Socket Directory: '%s'", s.SocketDir)
	}

	log.Infof("SocketFile: '%s'", s.SocketFile)
	exists, err = fileExists(s.SocketFile)
	if err != nil {
		log.Fatalf("Stat Socket File error: '%s'", err)
		os.Exit(1)
	}
	if exists {
		err = deleteFile(s.SocketFile)
		if err != nil {
			log.Fatalf("Delete Socket File error: '%s'", err)
			os.Exit(1)
		}
		log.Infof("Deleted Old Socket File: '%s'", s.SocketFile)
	}
	return s.SocketFile
}

// NewDriverSocket : return a new Socket Driver Config
func NewDriverSocket(socketDir string, driverName string) *DriverSocket {
	if socketDir == "" {
		socketDir = GetDefaultSocketDir()
	}
	return &DriverSocket{
		SocketDir:  socketDir,
		DriverName: driverName,
		SocketFile: socketDir + "/" + driverName + ".sock"}
}
