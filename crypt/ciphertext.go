package crypt

import "github.com/binarysoupdev/cryptool/crypt"

const SALT_SIZE = crypt.SALT_SIZE

type Ciphertext []byte

func (ct Ciphertext) Salt() []byte {
	return ct[:SALT_SIZE]
}

func (ct Ciphertext) Text() []byte {
	return ct[SALT_SIZE:]
}
