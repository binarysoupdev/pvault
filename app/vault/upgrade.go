package vault

import (
	"pvault/app/vault/database"

	"github.com/binarysoupdev/go-commando/errors"
)

func (v Vault) IsOutOfDate() bool {
	return v.Meta.DatabaseVersion < database.CURRENT_VERSION
}

func (v *Vault) Upgrade() error {
	if !v.IsOutOfDate() {
		return errors.New("vault is up-to-date")
	}

	new, err := v.Database.Upgrade(v.Map)
	if err != nil {
		return err
	}

	v.Database = new
	return nil
}
