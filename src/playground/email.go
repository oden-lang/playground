package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Mailer struct {
	key         string
	domain      string
	defaultFrom string
	defaultTo   string
}

func (m *Mailer) Send(from, to, subject, text, body string) error {
	client := &http.Client{}
	form := url.Values{
		"from":    {from},
		"to":      {to},
		"subject": {subject},
		"text":    {text},
		"html":    {body},
	}
	encodedForm := form.Encode()

	url := "https://api:" + m.key + "@api.mailgun.net/v3/" + m.domain + "/messages"

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(encodedForm))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len([]byte(encodedForm))))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		return nil
	}

	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return errors.New("Failed to send email: " + string(body))
	} else {
		return errors.New("Failed to send email, HTTP status: " + resp.Status)
	}
}

func codeParagraph(code string) string {
	return "<p><pre>" + code + "</pre></p>"
}

func (m *Mailer) SendCodeSuccess(odenCode, goCode string, output *PlayResponse) {
	var out string
	for _, e := range output.Events {
		out += e.Message
	}
	err := m.Send(m.defaultFrom, m.defaultTo, "Program successfully run!", "",
		"<html>"+
			"<p><strong>The following program was successfully run in the Oden Playground at <time>"+time.Now().String()+"</time>.</strong></p>"+
			"<p><pre>"+codeParagraph(odenCode)+"</pre></p>"+
			"<p><strong>The compiled Go code was:</strong></p>"+
			"<p><pre>"+codeParagraph(goCode)+"</pre></p>"+
			"<p><strong>The output was:</strong></p>"+
			"<p><pre>"+codeParagraph(out)+"</pre></p>"+
			"</html>")

	if err == nil {
		fmt.Println("Sent success email")
	} else {
		fmt.Println("Failed to send email:", err)
	}
}

func NewMailer(key, domain, defaultFrom, defaultTo string) *Mailer {
	m := Mailer{key, domain, defaultFrom, defaultTo}
	return &m
}

var mailer *Mailer

func init() {
	apiKey := os.Getenv("MAILGUN_API_KEY")
	mailer = NewMailer(apiKey, "oden-lang.org", "Oden Playground <playground@oden-lang.org>", "odenlanguage@gmail.com")
}
