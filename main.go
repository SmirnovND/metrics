package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

// Функция для генерации пары RSA ключей
func generateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	pubKey := &privKey.PublicKey
	return privKey, pubKey, nil
}

// Функция для сохранения зашифрованного приватного ключа
func saveEncryptedPrivateKey(privKey *rsa.PrivateKey, password []byte) error {
	privKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)

	// Зашифровываем ключ с использованием пароля
	block, err := x509.EncryptPEMBlock(rand.Reader, "ENCRYPTED PRIVATE KEY", privKeyBytes, password, x509.PEMCipherAES256)
	if err != nil {
		return fmt.Errorf("failed to encrypt private key: %v", err)
	}

	// Сохраняем зашифрованный ключ в файл
	privKeyFile, err := os.Create("encrypted_private_key.pem")
	if err != nil {
		return fmt.Errorf("failed to create private key file: %v", err)
	}
	defer privKeyFile.Close()

	err = pem.Encode(privKeyFile, block)
	if err != nil {
		return fmt.Errorf("failed to encode private key: %v", err)
	}

	return nil
}

func main() {
	// Генерируем пару ключей
	privKey, pubKey, err := generateKeyPair(2048)
	if err != nil {
		log.Fatalf("Error generating key pair: %v", err)
	}

	// Сохраняем публичный ключ в файл
	pubKeyFile, err := os.Create("public_key.pem")
	if err != nil {
		log.Fatalf("Error creating public key file: %v", err)
	}
	defer pubKeyFile.Close()

	err = pem.Encode(pubKeyFile, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(pubKey),
	})
	if err != nil {
		log.Fatalf("Error encoding public key: %v", err)
	}

	// Пароль для шифрования приватного ключа
	password := []byte("secretkey")

	// Сохраняем зашифрованный приватный ключ
	err = saveEncryptedPrivateKey(privKey, password)
	if err != nil {
		log.Fatalf("Error saving encrypted private key: %v", err)
	}

	fmt.Println("Private key and public key saved successfully!")
}
