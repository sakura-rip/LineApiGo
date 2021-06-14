// Copyright (c) 2020 @sakura-rip
// Version 1.1 beta
// LastUpdate 2020/08/28

package lineapigo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/apache/thrift/lib/go/thrift"
	sqlogin "github.com/sakura-rip/lineapigo/secondaryqrcodeloginservice"
	ser "github.com/sakura-rip/lineapigo/talkservice"
	talk "github.com/sakura-rip/lineapigo/talkservice"
)

// LoginRequest struct of qrlogin
type LoginRequest struct {
	login1     *sqlogin.SecondaryQRCodeLoginServiceClient
	loginCheck *sqlogin.SecondaryQRCodeLoginServiceClient
	sessionID  string
}

// SavedInfo struct of saved
type SavedInfo struct {
	allChatMids []string
	allChatInfo map[string]ser.Chat
	friendmids  []string
	profile     *talk.Profile
}

// LineClient struct of line
type LineClient struct {
	talk          *talk.TalkServiceClient
	poll          *talk.TalkServiceClient
	qrLogin       LoginRequest
	revision      int64
	ctx           context.Context
	appType       string
	reqSeq        int32
	reqSeqMessage map[string]int32
	globalRev     int64
	individualRev int64
	data          SavedInfo
}

// NewLineClient return new LineClient
func NewLineClient(appType string) *LineClient {
	return &LineClient{
		ctx:           context.Background(),
		revision:      -1,
		reqSeq:        0,
		qrLogin:       LoginRequest{},
		reqSeqMessage: map[string]int32{},
		appType:       appType,
	}
}

// LoginWithAuthToken >_<
func (cl *LineClient) LoginWithAuthToken(token string) {
	cl.talk = createTalkService(token, cl.appType, lineHostURL+normal)
	cl.poll = createTalkService(token, cl.appType, lineHostURL+longPolling)
}

// LoginWithQrCode >_<
func (cl *LineClient) LoginWithQrCode() {
	cl.createLogionSession1()
	cl.CreateQrSession()
	cl.createLogionSession2()
	url, er := cl.CreateQrCode()
	if er != nil {
		log.Printf("%+v\n", er)
	}
	fmt.Println("login this url on your mobile : " + url)

	cl.WaitForQrCodeVerified()
	_, err := cl.CertificateLogin("")
	if err != nil {
		pin, _ := cl.CreatePinCode()
		fmt.Println("input this pin code on your mobile : " + pin)
		cl.WaitForInputPinCode()
	}
	token, _, _ := cl.QrLogin()
	fmt.Println(token)

	cl.talk = createTalkService(token, cl.appType, lineHostURL+normal)
	cl.poll = createTalkService(token, cl.appType, lineHostURL+longPolling)
}

// createTalkService create thrift session using autoToken
func createTalkService(authToken, appType, targetURL string) *talk.TalkServiceClient {
	var transport thrift.TTransport

	option := thrift.THttpClientOptions{
		Client: &http.Client{
			Transport: &http.Transport{},
		},
	}
	transport, _ = thrift.NewTHttpClientWithOptions(targetURL, option)
	connect := transport.(*thrift.THttpClient)
	connect.SetHeader("X-Line-Access", authToken)
	connect.SetHeader("X-Line-Application", GetLineApplication(appType))
	connect.SetHeader("User-Agent", GetUserAgent(appType))
	connect.SetHeader("x-lal", GetXLal(appType))
	pcol := thrift.NewTCompactProtocol(transport)
	tstc := thrift.NewTStandardClient(pcol, pcol)
	return talk.NewTalkServiceClient(tstc)
}

func createLoginService(targetURL string, headers map[string]string) *sqlogin.SecondaryQRCodeLoginServiceClient {
	var transport thrift.TTransport
	option := thrift.THttpClientOptions{
		Client: &http.Client{
			Transport: &http.Transport{},
		},
	}
	transport, _ = thrift.NewTHttpClientWithOptions(targetURL, option)
	connect := transport.(*thrift.THttpClient)
	//set header
	for key, val := range headers {
		connect.SetHeader(key, val)
	}
	pcol := thrift.NewTCompactProtocol(transport)
	tstc := thrift.NewTStandardClient(pcol, pcol)
	return sqlogin.NewSecondaryQRCodeLoginServiceClient(tstc)
}

func (cl *LineClient) createLogionSession2() {
	cl.qrLogin.loginCheck = createLoginService(lineHostURL+newRegistration, map[string]string{
		"X-Line-Application": GetLineApplication(cl.appType) + ";SECONDARY",
		"User-Agent":         GetUserAgent(cl.appType),
		"x-lal":              GetXLal(cl.appType),
		"X-Line-Access":      cl.qrLogin.sessionID,
	})
}

func (cl *LineClient) createLogionSession1() {
	cl.qrLogin.login1 = createLoginService(lineHostURL+secondaryQrLogin, map[string]string{
		"X-Line-Application": GetLineApplication(cl.appType) + ";SECONDARY",
		"User-Agent":         GetUserAgent(cl.appType),
		"x-lal":              GetXLal(cl.appType),
	})
}
