package sendmail

import "testing"

func TestHello(t *testing.T) {
	options := Options{
		Host:     "smtp.yandex.ru",
		Port:     465,
		Login:    "you@login.ru",
		Password: "youPassword",
	}

	client, err := NewClient(options)
	if err != nil {
		t.Error(err)
	}

	message := "From: you@email.ru\r\n" +
		"To: to@email.ru\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Subject: you subject\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		"<p>Message Body</p>"

	err = client.Send("you@email.ru", "to@email.ru", message)
	if err != nil {
		t.Error(err)
	}
}
