/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package ssh

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
	"time"

	"git.cm/naiba/gocd"
	"github.com/pkg/errors"
	"io"
)

//GenKeyPair 创建密钥对
func GenKeyPair() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	var private bytes.Buffer
	if err := pem.Encode(&private, privateKeyPEM); err != nil {
		return "", "", err
	}

	// generate public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	public := ssh.MarshalAuthorizedKey(pub)
	return string(public), private.String(), nil
}

//CheckLogin 检查服务器是否存在
func CheckLogin(address string, port int, privateKey string, login string) error {
	conn, err := dial(address, login, privateKey, port)
	if err != nil {
		gocd.Logger().Errorln("ssh.CheckLogin", err)
		return errors.New("连接服务器失败")
	}
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil {
		return errors.New("建立会话失败")
	}
	defer session.Close()
	opt, err := session.Output("whoami")
	if strings.TrimSpace(string(opt)) != login {
		gocd.Logger().Infoln(string(opt) + "|" + err.Error())
		return errors.New("用户名验证失败")
	}
	return nil
}

//Deploy 进行部署
func Deploy(pipeline gocd.Pipeline, who string, saveLog func(log *gocd.PipeLog) error) {
	var pLog gocd.PipeLog
	pLog.PipelineID = pipeline.ID
	pLog.StartedAt = time.Now()
	pLog.Log = ""
	pLog.Pusher = who
	pLog.Status = gocd.PipeLogStatusRunning
	defer func() {
		// 保留最后5000字
		if len(pLog.Log) > 5000 {
			pLog.Log = pLog.Log[len(pLog.Log)-5000:]
		}
		pLog.Log = "[GoCD]" + pLog.StartedAt.String() + ": 开始执行\r\n" + pLog.Log
		pLog.StoppedAt = time.Now()
		saveLog(&pLog)
	}()

	conn, err := dial(pipeline.Server.Address, pipeline.Server.Login, pipeline.User.PrivateKey, pipeline.Server.Port)
	if err != nil {
		gocd.Logger().Debug(err)
		pLog.Status = gocd.PipeLogStatusErrorServerConn
		pLog.Log += "\r\n[GoCD]" + pLog.StartedAt.String() + ": 连接服务器失败"
		return
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		pLog.Status = gocd.PipeLogStatusErrorServerConn
		pLog.Log += "\r\n[GoCD]" + pLog.StartedAt.String() + ": 建立会话失败"
		return
	}
	defer session.Close()
	session.Wait()

	finish := make(chan bool, 1)
	defer close(finish)

	timer := time.NewTimer(time.Minute * 30)
	buf := new(bytes.Buffer)
	go func() {
		gocd.Logger().Debug("开始执行", pipeline.Shell)
		session.Stdout = buf
		err := session.Run(pipeline.Shell)
		if pLog.Status != gocd.PipeLogStatusRunning {
			return
		}
		if err != nil && err != io.EOF {
			gocd.Logger().Debug("执行失败", err.Error())
			pLog.Log += buf.String()
			pLog.Log += err.Error()
			pLog.Log += "\r\n[GoCD]" + pLog.StartedAt.String() + ": 执行失败"
			pLog.Status = gocd.PipeLogStatusErrorShellExec
		} else {
			pLog.Log += buf.String() + "\r\n [GoCD]" + time.Now().String() + ": 执行完毕"
			pLog.Status = gocd.PipeLogStatusSuccess
		}
		finish <- true
	}()

	select {
	case <-timer.C:
		gocd.Logger().Debug("执行超时", buf.String())
		pLog.Log += buf.String() + "\r\n [GoCD]" + time.Now().String() + ": 执行超时"
		pLog.Status = gocd.PipeLogStatusErrorShellExec
		return
	case <-finish:
		gocd.Logger().Debug("执行完毕")
		return
	}
}

func dial(address, user, pk string, port int) (*ssh.Client, error) {
	privateKey, err := ssh.ParsePrivateKey([]byte(pk))
	if err != nil {
		gocd.Logger().Debug(err, pk)
		return nil, errors.New("解析用户私钥失败")
	}
	return ssh.Dial("tcp", fmt.Sprintf("%s:%d", address, port), &ssh.ClientConfig{
		User:    user,
		Timeout: time.Second * 10,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privateKey),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
}
