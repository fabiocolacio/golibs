package crypto

import (
    "golang.org/x/crypto/scrypt"
)

var (
    KeyHashLength int = 32
)

// HashAndSaltPassword takes a password and salt, and creates a hash
// of the password concatenated with the salt using Scrypt.
func HashAndSaltPassword(passwd, salt []byte) ([]byte, error) {
    return scrypt.Key(passwd, salt, 32768, 8, 1, KeyHashLength)
}
