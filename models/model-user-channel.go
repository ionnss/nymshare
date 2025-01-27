package models

type Channel struct {
	ID        int    `db:"id" json:"id"`
	FullName  string `db:"full_name" json:"full_name"`
	Email     string `db:"email" json:"email"`
	PublicKey string `db:"public_key" json:"public_key"`
	PublicUrl string `db:"public_url" json:"public_url"`
	PublicQR  string `db:"public_qr" json:"public_qr"`
	Verified  bool   `db:"verified" json:"verified"`
}
