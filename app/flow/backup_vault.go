package flow

import (
	"fmt"
	"os"
	"path/filepath"
	"pvault/app/vault/local"
	"pvault/config"
	"time"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/got-style/style"
)

func BackupVault(v local.Vault, cfg config.Config) error {
	err := cfg.ValidateBackupPath()
	if err != nil {
		return errors.Chain(err, "error validating backup path")
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
