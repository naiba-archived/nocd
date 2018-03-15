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
	"golang.org/x/crypto/ssh"
	"crypto/rand"
	"github.com/pkg/errors"
	"fmt"
	"net"
	"git.cm/naiba/gocd"
	"strings"
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
	pk, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return errors.New("解析用户私钥失败")
	}
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", address, port), &ssh.ClientConfig{
		User: login,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(pk),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
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
