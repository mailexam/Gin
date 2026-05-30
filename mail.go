package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strconv"
)

func sendTest(to, subject, body string) error {
	login := os.Getenv("MAILEXAM_LOGIN")
	password := os.Getenv("MAILEXAM_PASSWORD")
	if login == "" || password == "" {
		return fmt.Errorf("MAILEXAM_LOGIN and MAILEXAM_PASSWORD must be set")
	}

	port, _ := strconv.Atoi(os.Getenv("MAILEXAM_PORT"))
	if port == 0 {
		port = 587
	}

	from := os.Getenv("MAIL_FROM")
	if from == "" {
		from = "noreply@example.test"
	}
	if to == "" {
		to = "user@example.test"
	}
	if subject == "" {
		subject = "Gin + Mailexam"
	}
	if body == "" {
		body = "Mailexam test from Gin"
	}

	host := login + ".mailexam.ru"
	addr := fmt.Sprintf("%s:%d", host, port)
	auth := smtp.PlainAuth("", login, password, host)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		from, to, subject, body,
	))

	if port == 587 || port == 2525 {
		return sendWithSTARTTLS(addr, host, auth, from, []string{to}, msg)
	}

	return smtp.SendMail(addr, auth, from, []string{to}, msg)
}

func sendWithSTARTTLS(addr, host string, auth smtp.Auth, from string, to []string, msg []byte) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err = client.StartTLS(&tls.Config{ServerName: host}); err != nil {
			return err
		}
	}

	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return err
		}
	}

	if err = client.Mail(from); err != nil {
		return err
	}
	for _, rcpt := range to {
		if err = client.Rcpt(rcpt); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err = w.Write(msg); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}

	return client.Quit()
}
