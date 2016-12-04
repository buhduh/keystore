package main

import (
	"bytes"
	"crypto/rand"
	b32 "encoding/base32"
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/png"
)

const (
	QR_FMT_STR string = "otpauth://totp/MyPersonal%%3Akeystore@keystore.com?secret=%s&issuer=keystore"
)

type QR struct {
	Secret string
	URI    []byte
}

func NewQR() QR {
	s := make([]byte, 20)
	//Going to assume this doesn't fail
	rand.Read(s)
	secret := b32.StdEncoding.EncodeToString(s)
	qrcode, _ := qr.Encode(fmt.Sprintf(QR_FMT_STR, secret), qr.L, qr.Auto)
	qrcode, _ = barcode.Scale(qrcode, 100, 100)
	buff := new(bytes.Buffer)
	png.Encode(buff, qrcode)
	return QR{Secret: secret, URI: buff.Bytes()}
}
