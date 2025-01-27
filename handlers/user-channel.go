package handlers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"nymshare/db"
	"nymshare/models"

	"github.com/skip2/go-qrcode"
)

func ChannelRegister(w http.ResponseWriter, r *http.Request) {
	// Method verification
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse form input
	fullName := r.FormValue("full_name")
	email := r.FormValue("email")

	// Validate input
	if fullName == "" || email == "" {
		http.Error(w, "Full name and Email are required", http.StatusBadRequest)
		return
	}

	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		http.Error(w, "Error generating keys", http.StatusInternalServerError)
		return
	}

	// Extract public key
	publicKey := &privateKey.PublicKey

	// Encode private key to PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyBytes := pem.EncodeToMemory(privateKeyPEM)

	// Encode public key to PEM format
	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	}
	publicKeyBytes := pem.EncodeToMemory(publicKeyPEM)

	// Generate public URL
	publicUrl := fmt.Sprintf("/channel/%s", email) // We'll make this more secure later

	// Generate QR code
	qr, err := qrcode.Encode(publicUrl, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "Error generating QR code", http.StatusInternalServerError)
		return
	}

	// Create channel record
	channel := &models.Channel{
		FullName:  fullName,
		Email:     email,
		PublicKey: string(publicKeyBytes),
		PublicUrl: publicUrl,
		PublicQR:  string(qr),
		Verified:  false,
	}

	// Save to database
	query := `
		INSERT INTO channels (full_name, email, public_key, public_url, public_qr, verified)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	err = db.DB.QueryRow(
		query,
		channel.FullName,
		channel.Email,
		channel.PublicKey,
		channel.PublicUrl,
		channel.PublicQR,
		channel.Verified,
	).Scan(&channel.ID)

	if err != nil {
		http.Error(w, "Error saving channel", http.StatusInternalServerError)
		return
	}

	// Set headers for private key download
	w.Header().Set("Content-Type", "application/x-pem-file")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pem", email))

	// Write private key to response
	w.Write(privateKeyBytes)
}
