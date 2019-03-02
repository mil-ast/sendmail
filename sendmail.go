package sendmail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
)

// Host host
type Host struct {
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
func NewClient(host Host) (Client, error) {
	var client Client

	servername := fmt.Sprintf("%s:%d", host.Host, host.Port)
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
	auth := smtp.PlainAuth("", host.Login, host.Password, splitHost)
	if err = client.c.Auth(auth); err != nil {
		return client, err
	}

	return client, nil
}

// Send send email
func (client Client) Send(from, to, body string) error {
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
	_, err = w.Write([]byte(body))
	if err != nil {
		return err
	}

	return w.Close()
}

// Quit QUIT command and closes the connection
func (client Client) Quit() error {
	return client.c.Quit()
}
