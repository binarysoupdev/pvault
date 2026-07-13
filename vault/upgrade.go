package vault

import (
	"pvault/vault/database/version2"

	"github.com/binarysoupdev/go-commando/errors"
)

func (v Vault) IsOutOfDate() bool {
	return v.Version() < CURRENT_VERSION
}

func (v *Vault) Upgrade() error {
	if !v.IsOutOfDate() {
		return errors.New("vault is up-to-date")
	}

	db := version2.NewDatabase(v.Path)

	err := db.Initialize(v.Index)
	if err != nil {
		return err
	}

	err = v.Database.Upgrade(v.Index)
	if err != nil {
		return errors.Chain(err, "error upgrading database")
	}

	v.Database = db
	return nil
}
