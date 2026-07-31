package vault

import (
	"pvault/app/vault/database"

	"github.com/binarysoupdev/go-commando/errors"
)

func (v Vault) IsOutOfDate() bool {
	return v.GetVersion() < database.CURRENT_VERSION
}

func (v *Vault) Upgrade() error {
	if !v.IsOutOfDate() {
		return errors.New("vault is up-to-date")
	}

	new, err := v.Database.Upgrade(v.Path)
	if err != nil {
		return errors.New("error upgrading database")
	}

	v.Database = new
	return nil
}
