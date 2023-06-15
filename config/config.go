package config

import "github.com/globe-and-citizen/layer8-genesis-repo/pkg"

type Config struct {
	RESTPort int
	GRPCPort int
	KeyPair  *pkg.KeyPair
}
