package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
)

func Sha1() HashBuilder {
	return NewWrappedHasher(sha1.New())
}

func Sha256() HashBuilder {
	return NewWrappedHasher(sha256.New())
}

func MD5() HashBuilder {
	return NewWrappedHasher(md5.New())
}
