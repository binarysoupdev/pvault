package vault

func (v Vault) IsOutOfDate() bool {
	return v.Version < CURRENT_VERSION
}
