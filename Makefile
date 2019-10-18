# Delete local binaries
clean:
	rm bin/tcpwave-cni
	rm bin/tcpwave-cni-daemon

# Ensure go dependencies
deps:
	dep ensure

# Build binaries
build: deps
	go build -o bin/tcpwave-cni ./plugin/
	go build -o bin/tcpwave-cni-daemon ./daemon/
