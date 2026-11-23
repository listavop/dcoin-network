
package core

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
)

type Validator struct {
	ID         string `json:"id"`
	PublicKey  string `json:"public_key"`
	PrivateKey ed25519.PrivateKey `json:"-"`
	Active     bool   `json:"active"`
}

func NewValidator() (*Validator, error) {
	publicKey, privateKey, err :=
		ed25519.GenerateKey(rand.Reader)

	if err != nil {
		return nil, err
	}

	return &Validator{
		ID:         hex.EncodeToString(publicKey),
		PublicKey:  hex.EncodeToString(publicKey),
		PrivateKey: privateKey,
		Active:     true,
	}, nil
}

func (v *Validator) Sign(message string) (string, error) {
	if !v.Active {
		return "", errors.New("validator inactive")
	}

	if len(v.PrivateKey) != ed25519.PrivateKeySize {
		return "", errors.New("invalid private key")
	}

	signature := ed25519.Sign(
		v.PrivateKey,
		[]byte(message),
	)

	return hex.EncodeToString(signature), nil
}

func VerifySignature(
	publicKeyHex string,
	message string,
	signatureHex string,
) bool {
	publicKey, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return false
	}

	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false
	}

	if len(publicKey) != ed25519.PublicKeySize ||
		len(signature) != ed25519.SignatureSize {
		return false
	}

	return ed25519.Verify(
		ed25519.PublicKey(publicKey),
		[]byte(message),
		signature,
	)
}

func (v *Validator) SetActive(active bool) {
	v.Active = active
}
