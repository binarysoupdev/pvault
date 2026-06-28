package vault

import "pvault/errors"

func (v Vault) IsOutOfDate() bool {
	return v.Version() < CURRENT_VERSION
}

func (v *Vault) Upgrade() error {
	if !v.IsOutOfDate() {
		return errors.New("vault is up-to-date")
	}

	err := v.Database.Upgrade()
	if err != nil {
		return errors.Chain(err, "error upgrading database")
	}
	v.Database = v.NewCurrentDatabase()

	err = v.Database.SaveIndex(v.Index)
	if err != nil {
		return errors.Chain(err, "error saving new index version")
	}

	return nil
}
