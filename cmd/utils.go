package cmd

import (
	"os"
	"os/exec"
)

func Execute(name string, args ...string) error {
	c := exec.Command(name, args...)

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Run()
}

var DefaultProtocArgs = []string{
	"-I", "protos/includes/googleapis",
	"-I", "protos/includes/grpc_ecosystem",
	"-I", "protos/includes/gnostic",
}
