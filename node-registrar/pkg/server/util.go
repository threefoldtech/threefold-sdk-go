package server

import (
	"errors"
	"fmt"

	"github.com/vedhavyas/go-subkey/v2/ed25519"
	"github.com/vedhavyas/go-subkey/v2/sr25519"
)

func verifySignature(publicKey, challenge, signature []byte) error {
	// Verify public key length
	if len(publicKey) != PubKeySize {
		return fmt.Errorf("invalid public key size: expected %d, got %d",
			PubKeySize, len(publicKey))
	}

	// Try ED25519 verification first
	edKey, err := ed25519.Scheme{}.FromPublicKey(publicKey)
	if err == nil && edKey.Verify(challenge, signature) {
		return nil
	}

	// Fallback to SR25519 verification
	srKey, err := sr25519.Scheme{}.FromPublicKey(publicKey)
	if err == nil && srKey.Verify(challenge, signature) {
		return nil
	}

	return errors.New("signature verification failed")
}
