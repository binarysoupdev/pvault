package local

import (
	"github.com/binarysoupdev/go-commando/errors"
)

func (v Vault) IsOutOfDate() bool {
	return v.GetVersion() < CURRENT_VERSION
}

func (v *Vault) Upgrade() error {
	if !v.IsOutOfDate() {
		return errors.New("vault is up-to-date")
	}
	var err error

	v.Database, err = v.Database.Upgrade(v.Map)
	if err != nil {
		return errors.Chain(err, "error upgrading index")
	}

	return nil
}
