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
	"time"

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

	rec := postContact(t, mailer, contactSignals{
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

	rec := postContact(t, mailer, contactSignals{
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

	rec := postContact(t, mailer, contactSignals{
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

func TestContactSubmitRejectsHeaderInjectionAndOversizedFields(t *testing.T) {
	tests := []struct {
		name    string
		signals contactSignals
		want    string
	}{
		{
			name:    "name header injection",
			signals: contactSignals{Name: "Attacker\r\nBcc: victim@example.com", Email: "test@example.com", Message: "A valid message body"},
			want:    "Name must be at most",
		},
		{
			name:    "email header injection",
			signals: contactSignals{Name: "Test User", Email: "test@example.com\r\nBcc: victim@example.com", Message: "A valid message body"},
			want:    "Valid email required",
		},
		{
			name:    "oversized message",
			signals: contactSignals{Name: "Test User", Email: "test@example.com", Message: strings.Repeat("x", contactMaxMessageLength+1)},
			want:    "at most 5,000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mailer := &fakeMailer{}
			rec := postContact(t, mailer, tt.signals)
			if mailer.calls != 0 {
				t.Fatalf("mailer calls = %d, want 0", mailer.calls)
			}
			if !strings.Contains(rec.Body.String(), tt.want) {
				t.Fatalf("response %q does not contain %q", rec.Body.String(), tt.want)
			}
		})
	}
}

func TestContactRateLimiter(t *testing.T) {
	limiter := &contactRateLimiter{}
	now := time.Unix(1_700_000_000, 0)
	for i := 0; i < contactRateMax; i++ {
		if !limiter.allow("192.0.2.1", now) {
			t.Fatalf("attempt %d unexpectedly rejected", i+1)
		}
	}
	if limiter.allow("192.0.2.1", now) {
		t.Fatal("attempt above rate limit was accepted")
	}
	if !limiter.allow("192.0.2.1", now.Add(contactRateWindow)) {
		t.Fatal("attempt after rate window was rejected")
	}
	if !limiter.consumeSolution("solution", now) {
		t.Fatal("new ALTCHA solution was rejected")
	}
	if limiter.consumeSolution("solution", now) {
		t.Fatal("replayed ALTCHA solution was accepted")
	}
	if !limiter.consumeSolution("solution", now.Add(contactReplayWindow)) {
		t.Fatal("expired replay-cache entry was not removed")
	}
}

func TestContactSubmitRejectsOversizedRequestBody(t *testing.T) {
	e := echo.New()
	body := `{"name":"Test","email":"test@example.com","message":"` + strings.Repeat("x", contactMaxBodyBytes) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/contact", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	err := contactSubmit("test-altcha-key", &fakeMailer{})(e.NewContext(req, rec))
	var httpErr *echo.HTTPError
	if !errors.As(err, &httpErr) || httpErr.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("error = %#v, want HTTP %d", err, http.StatusRequestEntityTooLarge)
	}
}

func postContact(t *testing.T, mailer mail.ContactMailer, signals contactSignals) *httptest.ResponseRecorder {
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

	if err := contactSubmit("test-altcha-key", mailer)(c); err != nil {
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
