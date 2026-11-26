package host

import (
	"pvault/vault"

	"github.com/binarysoupdev/cryptool/net/conn"
	"github.com/binarysoupdev/cryptool/net/host"
	"github.com/binarysoupdev/got-style/style"
)

type Host struct {
	Vault vault.Vault
}

func NewHost(v vault.Vault) Host {
	return Host{
		Vault: v,
	}
}

func (Host) Log(status int, _ int, v any) {
	switch status {
	case host.S_ERROR:
		style.Error.Printf("[X] %v\n", v)
	case host.S_ACCEPT_CLIENT:
		style.Create.Printf("[+] %v\n", v)
	case host.S_LOST_CLIENT:
		style.Delete.Printf("[-] %v\n", v)
	}
}

func (h Host) Handle(_ int, c *conn.Conn) error {
	c.SendMessage([]byte(h.Vault.Meta.Nickname))
	return nil
}
