package sendmail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
)

// Options параметры подключения
type Options struct {
	Host     string
	Port     uint16
	Login    string
	Password string
}

// Client client
type Client struct {
	c *smtp.Client
}

// NewClient create connection
func NewClient(options Options) (Client, error) {
	var client Client

	servername := fmt.Sprintf("%s:%d", options.Host, options.Port)
	splitHost, _, _ := net.SplitHostPort(servername)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         splitHost,
	}
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return client, err
	}

	client.c, err = smtp.NewClient(conn, splitHost)
	if err != nil {
		return client, err
	}

	// auth
	auth := smtp.PlainAuth("", options.Login, options.Password, splitHost)
	if err = client.c.Auth(auth); err != nil {
		return client, err
	}

	return client, nil
}

// Send send email
func (client Client) Send(from, to, email string) error {
	var err error

	addrFrom := mail.Address{Address: from}
	addrTo := mail.Address{Address: to}

	// from & to
	if err = client.c.Mail(addrFrom.Address); err != nil {
		return err
	}
	if err = client.c.Rcpt(addrTo.Address); err != nil {
		return err
	}

	// data
	w, err := client.c.Data()
	if err != nil {
		return err
	}
	// write
	_, err = w.Write([]byte(email))
	if err != nil {
		return err
	}

	return w.Close()
}

// Quit QUIT command and closes the connection
func (client Client) Quit() error {
	return client.c.Quit()
}
