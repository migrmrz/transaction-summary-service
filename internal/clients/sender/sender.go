package sender

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// type money uint64 // total value in cents

// func (m money) Dollars() uint64 { return uint64(m) / 100 }
// func (m money) Cents() uint64   { return uint64(m) % 100 }

type Email struct {
	Email                    string
	Name                     string
	TotalBalance             string
	TotalTransactionsByMonth map[string]int
	AverageCreditAmount      string
	AverageDebitAmount       string
}

type Config struct {
	FromMail      string        `mapstructure:"sender"`
	APIKey        string        `mapstructure:"api-key"`
	TemplateID    string        `mapstructure:"template-id"`
	MaxIdleConns  int           `mapstructure:"max-idle-conns"`
	Timeout       time.Duration `mapstructure:"timeout"`
	SendEmailFlag bool          `mapstructure:"send-email"`
}

type Sendgrid struct {
	FromMail   string
	ToMail     string
	TemplateID string
	APIKey     string
	SendFlag   bool
	httpClient *http.Client
}

func New(conf Config, fromMail string) Sendgrid {
	return Sendgrid{
		FromMail:   conf.FromMail,
		ToMail:     fromMail,
		TemplateID: conf.TemplateID,
		APIKey:     conf.APIKey,
		SendFlag:   conf.SendEmailFlag,
		httpClient: newHTTPClient(conf.MaxIdleConns, conf.Timeout),
	}
}

func newHTTPClient(maxIdleConnsParam int, timeoutParam time.Duration) *http.Client {
	maxIdleConns := maxIdleConnsParam
	maxIdleConnsPerHost := maxIdleConns
	requestTimeout := timeoutParam

	defaultTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConns:          maxIdleConns,
		MaxIdleConnsPerHost:   maxIdleConnsPerHost,
		IdleConnTimeout:       90 * time.Millisecond,
	}

	return &http.Client{Transport: defaultTransport, Timeout: requestTimeout}
}

func (sg Sendgrid) SendEmail(email Email) error {
	if !sg.SendFlag {
		return fmt.Errorf("email will not be sent. If this was intended otherwise, check configuration file")
	}

	from := mail.NewEmail("User", sg.FromMail)
	to := mail.NewEmail(email.Name, sg.ToMail)
	content := mail.NewContent("text/html", " ")
	m := mail.NewV3MailInit(from, "", to, content)
	m.SetTemplateID(sg.TemplateID)

	// Dynamic template data
	m.Personalizations[0].SetDynamicTemplateData("total_balance", email.TotalBalance)
	m.Personalizations[0].SetDynamicTemplateData("debit_amount", email.AverageDebitAmount)
	m.Personalizations[0].SetDynamicTemplateData("credit_amount", email.AverageCreditAmount)

	transactionsByMonth := make([]map[string]string, len(email.TotalTransactionsByMonth))

	count := 0

	for month, numberOfTransactions := range email.TotalTransactionsByMonth {
		transactionsByMonth[count] = map[string]string{
			"month":  month,
			"amount": fmt.Sprintf("%d", numberOfTransactions),
		}

		count += 1
	}

	m.Personalizations[0].SetDynamicTemplateData("transactions", transactionsByMonth)

	// Generate request
	request := sendgrid.GetRequest(sg.APIKey, "/v3/mail/send", "")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)

	// Send request
	client := &rest.Client{HTTPClient: sg.httpClient}
	response, err := client.Send(request)
	if err != nil {
		return err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		err = fmt.Errorf("status: %d (body: %v, headers: %v)", response.StatusCode, response.Body, response.Headers)

		return err
	}

	return nil
}
