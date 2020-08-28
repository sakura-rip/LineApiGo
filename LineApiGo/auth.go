// Copyright (c) 2020 @Ch31212y
// Version 1.1 bata
// LastUpdate 2020/08/28

package lineapigo

import (
	"log"

	sqlogin "github.com/ch31212y/lineapigo/secondaryqrcodeloginservice"
)

// CraeteQrSession create qr session for SecondaryQRLogin
func (cl *LineClient) CraeteQrSession() (string, error) {
	req := sqlogin.NewCreateQrSessionRequest()
	res, err := cl.qrLogin.login1.CreateSession(cl.ctx, req)
	if err != nil {
		log.Printf("%+v\n", err)
	}
	cl.qrLogin.sessionID = res.AuthSessionId
	return res.AuthSessionId, err
}

// CreateQrCode create qr code for SecondaryQRLogin
func (cl *LineClient) CreateQrCode() (string, error) {
	req := sqlogin.NewCreateQrCodeRequest()
	req.AuthSessionId = cl.qrLogin.sessionID
	res, err := cl.qrLogin.login1.CreateQrCode(cl.ctx, req)
	// FIXME return res.CallbackUrl + createSecret(), err
	return res.CallbackUrl, err
}

// WaitForQrCodeVerified wait for qr code verfied
func (cl *LineClient) WaitForQrCodeVerified() {
	req := sqlogin.NewCheckQrCodeVerifiedRequest()
	req.AuthSessionId = cl.qrLogin.sessionID
	_, err := cl.qrLogin.loginCheck.CheckQrCodeVerified(cl.ctx, req)
	if err != nil {
		log.Printf("%+v\n", err)
	}
}

// CertificateLogin try login with certificate
func (cl *LineClient) CertificateLogin(certificate string) (*sqlogin.VerifyCertificateResponse, error) {
	req := sqlogin.NewVerifyCertificateRequest()
	req.AuthSessionId = cl.qrLogin.sessionID
	res, err := cl.qrLogin.login1.VerifyCertificate(cl.ctx, req)
	return res, err
}

// CreatePinCode create pin code
func (cl *LineClient) CreatePinCode() (string, error) {
	req := sqlogin.NewCreatePinCodeRequest()
	req.AuthSessionId = cl.qrLogin.sessionID
	res, err := cl.qrLogin.login1.CreatePinCode(cl.ctx, req)
	if err != nil {
		log.Printf("%+v\n", err)
	}
	return res.PinCode, err
}

// WaitForInputPinCode wait for input pin
func (cl *LineClient) WaitForInputPinCode() {
	req := sqlogin.NewCheckPinCodeVerifiedRequest()
	req.AuthSessionId = cl.qrLogin.sessionID
	cl.qrLogin.loginCheck.CheckPinCodeVerified(cl.ctx, req)
}

// QrLogin qr login
func (cl *LineClient) QrLogin() (string, string, error) {
	req := sqlogin.NewQrCodeLoginRequest()
	req.AuthSessionId = cl.qrLogin.sessionID
	req.SystemName = systemName
	req.AutoLoginIsRequired = true
	callback, err := cl.qrLogin.login1.QrCodeLogin(cl.ctx, req)
	return callback.AccessToken, callback.Certificate, err
}

// FIXME
// func createSecret() string {
// 	privateKey := curve25519.GeneratePrivateKey()
// 	curve25519.PublicKey(privateKey)
// }
