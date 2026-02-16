package crypto

import (
	"encoding/base64"
	"testing"
)

func TestDeriveKey(t *testing.T) {
	pass := "secret"
	salt := "salt"
	key, err := DeriveKey(pass, salt)
	if err != nil {
		t.Fatalf("DeriveKey failed: %v", err)
	}
	if len(key) != 32 {
		t.Errorf("Expected key length 32, got %d", len(key))
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key, _ := DeriveKey("pass", "salt")
	plaintext := "Hello World"
	
	encrypted, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	
	decrypted, err := Decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}
	
	if decrypted != plaintext {
		t.Errorf("Expected %s, got %s", plaintext, decrypted)
	}
}

func TestDecryptTampered(t *testing.T) {
	key, _ := DeriveKey("pass", "salt")
	plaintext := "Hello World"
	
	encrypted, _ := Encrypt(key, plaintext)
	
	bytes, _ := base64.StdEncoding.DecodeString(encrypted)
	bytes[len(bytes)-1] ^= 0x01
	tampered := base64.StdEncoding.EncodeToString(bytes)
	
	_, err := Decrypt(key, tampered)
	if err == nil {
		t.Error("Expected error for tampered ciphertext, got nil")
	}
}
