package types

import (
	"errors"
	"strings"

	"github.com/docker/go-connections/nat"
)

// ProtoPort represents a string port/proto Eg. "8080/tcp"
type ProtoPort string

func (p ProtoPort) ParsePort() (port string, proto string, err error) {
	if string(p) == "" {
		return "", "", errors.New("invalid port: value is empty")
	}

	port, proto, found := strings.Cut(string(p), "/")
	if !found {
		return "", "", errors.New("invalid port missing seperator '/'")
	}
	return port, port, nil

}

type PortMap nat.PortMap

// Expose Maps a single port proto to host to container
func (m PortMap) Expose(port ProtoPort, hostIP string) error {
	cport, cproto, err := port.ParsePort()
	if err != nil {
		return err
	}
	p, err := nat.NewPort(cproto, cport)
	if err != nil {
		return err
	}

	if _, ok := m[p]; !ok {
		m[p] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: p.Port(),
			},
		}
		return nil
	}
	m[p] = append(m[p], nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: p.Port(),
	})

	return nil

}
