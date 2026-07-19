package flow

import (
	"pvault/app/vault/local"

	"github.com/binarysoupdev/go-commando/errors"
)

func OpenLocalVault(path string) (*local.Vault, error) {
	v, err := local.OpenVault(path)
	if err != nil {
		return nil, errors.New("error opening vault (run \"vault -init\" to repair)")
	}

	if v.IsOutOfDate() {
		return nil, errors.Format("vault (@v%d) out-of-date (run \"vault -upgrade\" to repair)", v.GetVersion())
	}

	return &v, nil
}
