package vault_flow

import (
	"pvault/remote/client"

	"github.com/binarysoupdev/got-style/style"
)

func OpenRemoteVault(path string) (client.Vault, error) {
	v, err := client.ConnectToVault(path)
	if err != nil {
		return client.Vault{}, nil
	}

	style.Create.Printf("Connected to: %s\n", v.RemoteAddress())
	return v, nil
}
