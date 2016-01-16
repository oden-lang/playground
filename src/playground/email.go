package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
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

var successTemplate *template.Template
var odenErrorTemplate *template.Template
var goRunErrorTemplate *template.Template

type SuccessViewModel struct {
	Time       string
	OdenSource string
	GoCode     string
	Events     []Event
}

type OdenErrorViewModel struct {
	Time       string
	OdenSource string
	Output     string
}

type GoRunErrorViewModel struct {
	Time       string
	OdenSource string
	GoCode     string
	Events     string
}

func (m *Mailer) SendTemplated(subject string, tmpl *template.Template, data interface{}) {
	var out bytes.Buffer
	err := tmpl.Execute(&out, data)
	if err != nil {
		fmt.Printf("Could not execute %s template on data (%v): %s\n", tmpl.Name(), data, err)
		return
	}

	err = m.Send(m.defaultFrom, m.defaultTo, subject, "", out.String())
	if err == nil {
		fmt.Println("Sent", tmpl.Name(), "email")
	} else {
		fmt.Println("Failed to send", tmpl.Name(), "email:", err)
	}
}

func (m *Mailer) SendCodeSuccess(odenCode, goCode string, output *PlayResponse) {
	data := SuccessViewModel{time.Now().Format("Mon Jan _2 15:04:05 2006"), odenCode, goCode, output.Events}
	m.SendTemplated("Program successfully run!", successTemplate, data)
}

func (m *Mailer) SendOdenCompilationError(odenCode, output string) {
	data := OdenErrorViewModel{time.Now().Format("Mon Jan _2 15:04:05 2006"), odenCode, output}
	m.SendTemplated("Program did not compile!", odenErrorTemplate, data)
}

func (m *Mailer) SendGoRunError(odenCode, goCode, output string) {
	data := GoRunErrorViewModel{time.Now().Format("Mon Jan _2 15:04:05 2006"), odenCode, goCode, output}
	m.SendTemplated("Program did not run!", goRunErrorTemplate, data)
}

func NewMailer(key, domain, defaultFrom, defaultTo string) *Mailer {
	m := Mailer{key, domain, defaultFrom, defaultTo}
	return &m
}

var mailer *Mailer

func mustParse(name string, contents string) *template.Template {
	tmpl, err := template.New(name).Parse(contents)

	if err != nil {
		panic(fmt.Sprintf("Could not parse %s template:", name, err))
	} else {
		return tmpl
	}
}

func init() {
	apiKey := os.Getenv("MAILGUN_API_KEY")
	mailer = NewMailer(apiKey, "oden-lang.org", "Oden Playground <playground@oden-lang.org>", "odenlanguage@gmail.com")

	successTemplate = mustParse("success", `
		<p>
			<strong>The following program was successfully run in the Oden Playground at <time>{{.Time}}</time>.</strong>
		</p>
		<p><pre>{{.OdenSource}}</pre></p>
		<p><strong>The compiled Go code was:</strong></p>
		<p><pre>{{.GoCode}}</pre></p>
		<p><strong>The output was:</strong></p>
		<p>
			<pre>{{range .Events}}{{.Message}}{{end}}</pre>
		</p>
	`)

	odenErrorTemplate = mustParse("Oden compilation error", `
		<p>
			<strong color="red">The following program failed to compile in the Oden Playground at <time>{{.Time}}</time>.</strong>
		</p>
		<p><pre>{{.OdenSource}}</pre></p>
		<p><strong>The output was:</strong></p>
		<p>
			<pre>{{.Output}}</pre>
		</p>
	`)

	goRunErrorTemplate = mustParse("Go run error", `
		<p>
			<strong color="red">The following program did not run in play.golang.org at <time>{{.Time}}</time>.</strong>
		</p>
		<p><pre>{{.OdenSource}}</pre></p>
		<p><strong>The compiled Go code was:</strong></p>
		<p><pre>{{.GoCode}}</pre></p>
		<p><strong>The output was:</strong></p>
		<p>
			<pre>{{.Output}}</pre>
		</p>
	`)
}
