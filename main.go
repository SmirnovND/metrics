package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
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

// Функция для создания AES-GCM ключа из пароля
func createAESKeyFromPassword(password []byte) ([]byte, error) {
	// Используем SHA-256 для создания ключа из пароля
	hash := sha256.Sum256(password)
	return hash[:], nil
}

// Функция для шифрования данных с использованием AES-GCM
func encryptWithAESGCM(data []byte, key []byte) ([]byte, error) {
	// Генерируем случайный IV (инициализирующий вектор)
	iv := make([]byte, 12) // Для AES-GCM IV должен быть длиной 12 байт
	_, err := rand.Read(iv)
	if err != nil {
		return nil, fmt.Errorf("failed to generate IV: %v", err)
	}

	// Создаем объект AES в режиме GCM
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %v", err)
	}

	// Шифруем данные
	ciphertext := gcm.Seal(nil, iv, data, nil)

	// Возвращаем IV + зашифрованные данные
	return append(iv, ciphertext...), nil
}

// Функция для сохранения зашифрованного приватного ключа
func saveEncryptedPrivateKey(privKey *rsa.PrivateKey, password []byte) error {
	privKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)

	// Генерируем AES ключ из пароля
	aesKey, err := createAESKeyFromPassword(password)
	if err != nil {
		return fmt.Errorf("failed to create AES key: %v", err)
	}

	// Шифруем приватный ключ с использованием AES-GCM
	encryptedKey, err := encryptWithAESGCM(privKeyBytes, aesKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt private key: %v", err)
	}

	// Сохраняем зашифрованный ключ в файл
	privKeyFile, err := os.Create("encrypted_private_key.pem")
	if err != nil {
		return fmt.Errorf("failed to create private key file: %v", err)
	}
	defer privKeyFile.Close()

	// Кодируем зашифрованный приватный ключ в PEM формат
	err = pem.Encode(privKeyFile, &pem.Block{
		Type:  "ENCRYPTED PRIVATE KEY",
		Bytes: encryptedKey,
	})
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
