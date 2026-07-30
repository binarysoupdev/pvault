package flow

import (
	"fmt"
	"os"
	"path/filepath"
	"pvault/app/vault"
	"pvault/config"
	"time"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/got-style/style"
)

func OpenVault(path string) (vault.Vault, error) {
	v, err := vault.Open(path)
	if err != nil {
		return vault.Vault{}, errors.New("error opening vault (run \"vault -init\" to repair)")
	}

	if v.IsOutOfDate() {
		return vault.Vault{}, errors.Format("vault (@v%d) out-of-date (run \"vault -upgrade\" to repair)", v.GetVersion())
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
		return errors.Chain(err, "error backing vault")
	}

	style.BoldCreate.Printf("[+] Created Backup \"%s\"\n", path)
	return nil
}
