package client

import (
	"github.com/binarysoupdev/cryptool/net/client"
	"github.com/binarysoupdev/go-extensions/errors"
)

type Vault struct {
	client client.Client
}

func ConnectToVault(address string) (Vault, error) {
	c, err := client.Connect(address)
	if err != nil {
		return Vault{}, errors.Chain(err, "error connecting to remote vault")
	}

	return Vault{
		client: c,
	}, nil
}

func (v Vault) Close() error {
	return v.client.Close()
}

func (v Vault) RemoteAddress() string {
	return v.client.RemoteAddr().String()
}
