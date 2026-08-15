package vault

import (
	"pvault/vault/database"

	"github.com/binarysoupdev/go-extensions/errors"
)

func (v Vault) IsOutOfDate() bool {
	return v.GetVersion() < database.CURRENT_VERSION
}

func (v *Vault) Upgrade() error {
	if !v.IsOutOfDate() {
		return errors.New("vault is up-to-date")
	}

	new, err := v.DatabaseEncoder.Upgrade(v.Path)
	if err != nil {
		return errors.New("error upgrading database")
	}

	v.DatabaseEncoder = new
	v.Meta.DatabaseVersion = new.GetVersion()

	err = v.SaveMetadata()
	if err != nil {
		return err
	}

	return nil
}
