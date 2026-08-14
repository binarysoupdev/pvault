package vault_flow

import (
	errors_std "errors"
	"fmt"
	"os"
	"path/filepath"
	"pvault/app/config"
	"pvault/vault"
	"time"

	"github.com/binarysoupdev/go-commando/logger"

	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/got-style/style"
)

func OpenCurrentVault(path string) (vault.Vault, error) {
	v, err := OpenLegacyVault(path)
	if err != nil {
		return vault.Vault{}, err
	}

	if v.IsOutOfDate() {
		return vault.Vault{}, errors.Format("vault (@v%d) out-of-date (run \"vault -upgrade\" to repair)", v.GetVersion())
	}

	return v, nil
}

func OpenLegacyVault(path string) (vault.Vault, error) {
	v, err := vault.Open(path)

	if errors_std.Is(err, vault.ErrorNotFound{}) {
		return vault.Vault{}, errors.New("vault not found (run \"vault -init\" to create)")
	}
	if err != nil {
		logger.LogError(err)
		return vault.Vault{}, errors.New("error opening vault (unsupported or malformed)")
	}

	return v, nil
}

func BackupVault(v vault.Vault, cfg config.Config) error {
	err := cfg.ValidateBackupPath()
	if err != nil {
		return errors.Chain(err, "error validating \"config.backup_path\"")
	}

	path := filepath.Join(cfg.BackupPath, fmt.Sprintf("%s (v%d)", time.Now().Format(time.DateTime), v.GetVersion()))

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return errors.Chain(err, "error creating backup directory")
	}

	err = v.Backup(path)
	if err != nil {
		logger.LogError(err)
		return errors.New("error backing vault")
	}

	logger.Logf("[+] create backup %s", path)
	style.BoldCreate.Printf("[+] Created Backup \"%s\"\n", path)

	return nil
}
