package v1

import "pvault/crypt"

const LEGACY_HASH_SIZE = 60

func MarshalFromLegacy(name string, bytes []byte) []byte {
	return append(buildHeader(name), bytes[LEGACY_HASH_SIZE:]...)
}

func (r Record) MarshalToLegacy(password string) ([]byte, error) {
	hash := make([]byte, LEGACY_HASH_SIZE)

	ciphertext, err := crypt.Marshal(password, r)
	if err != nil {
		return nil, err
	}

	return append(hash, ciphertext...), nil
}
