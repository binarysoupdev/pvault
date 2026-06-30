package vault

import "pvault/errors"

func (v Vault) IsOutOfDate() bool {
	return v.Version() < CURRENT_VERSION
}

func (v *Vault) Upgrade() error {
	if !v.IsOutOfDate() {
		return errors.New("vault is up-to-date")
	}

	db, err := v.initNewDatabaseVersion2()
	if err != nil {
		return err
	}

	err = v.Database.Upgrade(v.Index, db)
	if err != nil {
		return errors.Chain(err, "error upgrading database")
	}

	v.Database = db
	return nil
}
