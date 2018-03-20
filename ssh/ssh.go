/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package ssh

import (
	"crypto/rsa"
	"encoding/pem"
	"crypto/x509"
	"bytes"
	"crypto/rand"
	"fmt"
	"net"
	"strings"
	"time"
	"golang.org/x/crypto/ssh"

	"github.com/pkg/errors"
	"git.cm/naiba/gocd"
)

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

func CheckLogin(address string, port int, privateKey string, login string) error {
	conn, err := dial(address, login, privateKey, port)
	if err != nil {
		gocd.Log.Error(err)
		return errors.New("连接服务器失败")
	}
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil {
		return errors.New("建立会话失败")
	}
	defer session.Close()
	opt, err := session.Output("whoami")
	if strings.TrimSpace(string(opt)) == login {
		return nil
	} else {
		gocd.Log.Info(string(opt))
		return errors.New("用户名验证失败")
	}
}

func Deploy(pipeline gocd.Pipeline, who string, saveLog func(plog *gocd.PipeLog) error) {
	var plog gocd.PipeLog
	plog.PipelineID = pipeline.ID
	plog.StartedAt = time.Now()
	plog.Pusher = who
	defer saveLog(&plog)

	conn, err := dial(pipeline.Server.Address, pipeline.Server.Login, pipeline.User.PrivateKey, pipeline.Server.Port)
	if err != nil {
		gocd.Log.Debug(err)
		plog.Status = gocd.PipeLogStatusErrorServerConn
		plog.Log = "连接服务器失败"
		return
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		plog.Status = gocd.PipeLogStatusErrorServerConn
		plog.Log = "建立会话失败"
		return
	}
	defer session.Close()
	session.Wait()

	timer := time.NewTimer(time.Second * 10)
	finish := make(chan bool, 1)
	go func() {
		gocd.Log.Debug("开始执行", pipeline.Shell)
		out, err := session.CombinedOutput(pipeline.Shell)
		if err != nil {
			plog.Log = err.Error()
			plog.Status = gocd.PipeLogStatusErrorShellExec
		} else {
			plog.Log = string(out)
			plog.Status = gocd.PipeLogStatusSuccess
		}
		plog.StoppedAt = time.Now()
		finish <- true
	}()
	select {
	case <-timer.C:
		gocd.Log.Debug("执行失败")
		plog.Log = "Shell 执行超时"
		plog.Status = gocd.PipeLogStatusErrorShellExec
		plog.StoppedAt = time.Now()
		return
	case <-finish:
		gocd.Log.Debug("执行完毕")
		return
	}
}

func dial(address, user, pk string, port int) (*ssh.Client, error) {
	privateKey, err := ssh.ParsePrivateKey([]byte(pk))
	if err != nil {
		gocd.Log.Debug(err, pk)
		return nil, errors.New("解析用户私钥失败")
	}
	return ssh.Dial("tcp", fmt.Sprintf("%s:%d", address, port), &ssh.ClientConfig{
		User:    user,
		Timeout: time.Second * 30,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privateKey),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
}
