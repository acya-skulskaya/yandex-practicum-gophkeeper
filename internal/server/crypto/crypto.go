package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

const (
	iterations = 100000
)

type Crypto struct {
	aesgcm cipher.AEAD
	nonce  []byte
}

func NewCrypto(salt string, keyFile string) (*Crypto, error) {
	//nolint:gosec // ignore
	file, err := os.OpenFile(keyFile, os.O_RDONLY, 0o666)
	if err != nil {
		return nil, fmt.Errorf("error opening key file %s: %w", keyFile, err)
	}
	//nolint:errcheck // ignore err
	defer file.Close()
	fileinfo, _ := file.Stat()
	key := make([]byte, fileinfo.Size())
	_, err = file.Read(key)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", keyFile, err)
	}

	derivedKey := pbkdf2.Key(key, []byte(salt), iterations, 32, sha256.New)

	// NewCipher создает и возвращает новый cipher.Block.
	// Ключевым аргументом должен быть ключ AES, 16, 24 или 32 байта
	// для выбора AES-128, AES-192 или AES-256.
	aesblock, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, fmt.Errorf("could not create AES block cipher: %w", err)
	}

	// NewGCM возвращает заданный 128-битный блочный шифр
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, fmt.Errorf("could not create GCM cipher: %w", err)
	}

	// создаём вектор инициализации
	nonce := derivedKey[len(derivedKey)-aesgcm.NonceSize():]

	return &Crypto{
		aesgcm: aesgcm,
		nonce:  nonce,
	}, nil
}

func NewCryptoMock(salt string, keyFile string) (*Crypto, error) {
	return &Crypto{
		aesgcm: nil,
		nonce:  nil,
	}, nil
}

func (c *Crypto) Encrypt(src []byte) ([]byte, error) {
	if src == nil || len(src) == 0 {
		return nil, nil
	}
	dst := c.aesgcm.Seal(nil, c.nonce, src, nil) // зашифровываем

	return dst, nil
}

func (c *Crypto) Decrypt(dst []byte) ([]byte, error) {
	if len(dst) == 0 {
		return nil, nil
	}

	src, err := c.aesgcm.Open(nil, c.nonce, dst, nil) // расшифровываем
	if err != nil {
		return nil, fmt.Errorf("could not decrypt message: %w", err)
	}

	return src, nil
}

func (c *Crypto) ToHexString(b []byte) string {
	if len(b) == 0 {
		return ""
	}

	return hex.EncodeToString(b)
}
func (c *Crypto) FromHexString(s string) ([]byte, error) {
	if s == "" {
		return nil, nil
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("could not decode hex string: %w", err)
	}
	return b, nil
}
