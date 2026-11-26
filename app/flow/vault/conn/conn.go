package conn_flow

import (
	"flag"
	vault_flow "pvault/app/flow/vault"
)

type ConnFlow struct {
	localPath  string
	remoteAddr *string
}

func NewConnFlow(flags *flag.FlagSet, local string) ConnFlow {
	return ConnFlow{
		localPath:  local,
		remoteAddr: flags.String("addr", "", "address of a remote vault"),
	}
}

func (f ConnFlow) OpenVaultConnection() (vault_flow.Vault, error) {
	if *f.remoteAddr != "" {
		return vault_flow.OpenRemoteVault(*f.remoteAddr)
	} else {
		return vault_flow.OpenCurrentVault(f.localPath)
	}
}
