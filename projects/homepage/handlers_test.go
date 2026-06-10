package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mljr-web/internal/mail"

	altchalib "github.com/altcha-org/altcha-lib-go"
	"github.com/labstack/echo/v4"
)

type fakeMailer struct {
	err      error
	calls    int
	messages []mail.ContactMessage
}

func (m *fakeMailer) SendContact(_ context.Context, msg mail.ContactMessage) error {
	m.calls++
	m.messages = append(m.messages, msg)
	return m.err
}

func TestContactSubmitSendsMailAfterAltcha(t *testing.T) {
	key := "test-altcha-key"
	mailer := &fakeMailer{}

	rec := postContact(t, key, mailer, contactSignals{
		Name:    "Test User",
		Email:   "test@example.com",
		Message: "This is a valid test message.",
		Altcha:  validAltchaPayload(t, key),
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if mailer.calls != 1 {
		t.Fatalf("mailer calls = %d, want 1", mailer.calls)
	}
	if got := mailer.messages[0].Email; got != "test@example.com" {
		t.Fatalf("message email = %q", got)
	}
	if !strings.Contains(rec.Body.String(), "Message sent!") {
		t.Fatalf("response did not contain success patch: %s", rec.Body.String())
	}
}

func TestContactSubmitReportsMailFailure(t *testing.T) {
	key := "test-altcha-key"
	mailer := &fakeMailer{err: errors.New("smtp unavailable")}

	rec := postContact(t, key, mailer, contactSignals{
		Name:    "Test User",
		Email:   "test@example.com",
		Message: "This is a valid test message.",
		Altcha:  validAltchaPayload(t, key),
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if mailer.calls != 1 {
		t.Fatalf("mailer calls = %d, want 1", mailer.calls)
	}
	if !strings.Contains(rec.Body.String(), "Message delivery failed") {
		t.Fatalf("response did not contain delivery error: %s", rec.Body.String())
	}
}

func TestContactSubmitDoesNotSendMailWithInvalidAltcha(t *testing.T) {
	mailer := &fakeMailer{}

	rec := postContact(t, "test-altcha-key", mailer, contactSignals{
		Name:    "Test User",
		Email:   "test@example.com",
		Message: "This is a valid test message.",
		Altcha:  "invalid",
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if mailer.calls != 0 {
		t.Fatalf("mailer calls = %d, want 0", mailer.calls)
	}
	if !strings.Contains(rec.Body.String(), "Captcha verification failed") {
		t.Fatalf("response did not contain captcha error: %s", rec.Body.String())
	}
}

func postContact(t *testing.T, key string, mailer mail.ContactMailer, signals contactSignals) *httptest.ResponseRecorder {
	t.Helper()

	body, err := json.Marshal(signals)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/contact", strings.NewReader(string(body)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("User-Agent", "handler-test")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := contactSubmit(key, mailer)(c); err != nil {
		t.Fatal(err)
	}
	return rec
}

func validAltchaPayload(t *testing.T, key string) string {
	t.Helper()

	number := int64(7)
	ch, err := altchalib.CreateChallenge(altchalib.ChallengeOptions{
		HMACKey: key,
		Number:  &number,
		Salt:    "test-salt",
	})
	if err != nil {
		t.Fatal(err)
	}
	payload := altchalib.Payload{
		Algorithm: ch.Algorithm,
		Challenge: ch.Challenge,
		Number:    number,
		Salt:      ch.Salt,
		Signature: ch.Signature,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(raw)
}
