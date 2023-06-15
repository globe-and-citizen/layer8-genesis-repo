package config

import "github.com/globe-and-citizen/layer8-genesis-repo/pkg"

type Config struct {
	Port    int
	KeyPair *pkg.KeyPair
}
