// Package o365verify uses the autodiscover JSON API of Office 365 to enumerate valid email addresses
package o365verify

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Result contains all available information about the email address
type Result struct {
	EmailAddress       string // the email address searched
	CalculatedBETarget string // the mailbox server to which the request is routed. Reference: https://docs.microsoft.com/en-us/exchange/management/health/troubleshooting-autodiscover-health-set
	MailboxGUID        string // the GUID of the mailbox
	ValidAddress       bool   // is the email address valid
	DomainIsO365       bool   // is the domain an office365 domain
}

// VerifyAddress looks up the specified email address in O365 and returns all available information
func VerifyAddress(email string) (*Result, error) {
	var (
		r   Result
		err error
	)

	domain := strings.Split(email, "@")[1]
	r.EmailAddress = email

	client := &http.Client{
		Timeout: 45 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { // don't follow redirects
			return http.ErrUseLastResponse
		},
	}

	url := fmt.Sprintf("https://outlook.office365.com/autodiscover/autodiscover.json/v1.0/%s@%s?Protocol=Autodiscoverv1", randomString(15), domain)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Microsoft Office/16.0 (Windows NT 10.0; Microsoft Outlook 16.0.12026; Pro)")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	rawBody, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	if strings.Contains(string(rawBody), "outlook.office365.com") {
		r.DomainIsO365 = true
	}

	url = fmt.Sprintf("https://outlook.office365.com/autodiscover/autodiscover.json/v1.0/%s?Protocol=Autodiscoverv1", email)
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Microsoft Office/16.0 (Windows NT 10.0; Microsoft Outlook 16.0.12026; Pro)")
	req.Header.Set("Accept", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	rawBody, _ = ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	if resp.Header.Get("X-MailboxGuid") != "" {
		if resp.StatusCode == 200 {
			r.ValidAddress = true
		} else if resp.StatusCode == 302 && !strings.Contains(string(rawBody), "outlook.office365.com") {
			r.ValidAddress = true
		}
	}

	r.MailboxGUID = resp.Header.Get("X-MailboxGuid")
	r.CalculatedBETarget = resp.Header.Get("X-CalculatedBETarget")

	return &r, nil
}

// Returns an int >= min, < max
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// Generate a random string of A-Z chars with len = l
func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}
