package password_encoder

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func Encoder(password string) (string, error) {

	salt := make([]byte, 16)

	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	time := uint32(1)
	memory := uint32(64 * 1024)
	threads := uint8(4)
	keyLen := uint32(32)

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		memory, time, threads, b64Salt, b64Hash)

	return encoded, nil
}

func Match(password, encodedHash string) (bool, error) {

	parts := strings.Split(encodedHash, "$")

	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	var memory uint32
	var time uint32
	var threads uint8

	if parts[2] != "v=19" {
		return false, errors.New("incompatible argon2 version")
	}
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, err
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	computedHash := argon2.IDKey(
		[]byte(password),
		salt,
		time,
		memory,
		threads,
		uint32(len(hash)),
	)
	if subtle.ConstantTimeCompare(hash, computedHash) == 1 {
		return true, nil
	}
	return false, nil
}
