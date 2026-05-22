package steps

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/mail"
	"sync"
	"sync/atomic"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// SmtpCapturedMessage carries the fields the SMTP sink decodes from a
// successfully delivered message. Stored on the bag so asserts can
// validate subject/to/from/body without re-parsing RFC 5322.
type SmtpCapturedMessage struct {
	From    string
	To      string
	Subject string
	Body    string
	Raw     []byte
}

// smtpBackend implements smtp.Backend by routing every NewSession to
// a session that delegates message capture to the bag-stored counter
// and last-message slot.
type smtpBackend struct {
	mu      *sync.Mutex
	hits    *atomic.Int64
	lastMsg **SmtpCapturedMessage
}

func (b *smtpBackend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &smtpSession{backend: b}, nil
}

type smtpSession struct {
	backend *smtpBackend
	from    string
	to      string
}

// AuthMechanisms advertises AUTH PLAIN so Go's net/smtp.SendMail with
// PlainAuth is accepted. Required by go-smtp v0.21+ to surface AUTH in
// EHLO.
func (s *smtpSession) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

// Auth returns a SASL server that accepts any credentials — the sink
// is for tests, not security.
func (s *smtpSession) Auth(mech string) (sasl.Server, error) {
	return sasl.NewPlainServer(func(identity, username, password string) error {
		return nil
	}), nil
}

func (s *smtpSession) Mail(from string, _ *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *smtpSession) Rcpt(to string, _ *smtp.RcptOptions) error {
	s.to = to
	return nil
}

func (s *smtpSession) Data(r io.Reader) error {
	raw, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("smtp sink: read data: %w", err)
	}
	captured := &SmtpCapturedMessage{From: s.from, To: s.to, Raw: raw}
	if msg, perr := mail.ReadMessage(bytes.NewReader(raw)); perr == nil {
		captured.Subject = msg.Header.Get("Subject")
		body, _ := io.ReadAll(msg.Body)
		captured.Body = string(body)
	}
	s.backend.mu.Lock()
	*s.backend.lastMsg = captured
	s.backend.mu.Unlock()
	s.backend.hits.Add(1)
	return nil
}

func (s *smtpSession) Reset()        { s.from = ""; s.to = "" }
func (s *smtpSession) Logout() error { return nil }

// StartSmtpSink boots an in-process SMTP server on
// constants.SmtpSinkBindAddr that accepts any AUTH PLAIN and captures
// the most recent message into the bag.
//
// Writes (bag):
//   - BagKeySmtpServer       *smtp.Server
//   - BagKeySmtpHits         *atomic.Int64
//   - BagKeySmtpLastMessage  **SmtpCapturedMessage
//
// Compensate: smtp.Server.Close() with a 2s deadline.
func StartSmtpSink() saga.Step {
	return saga.Step{
		Name: "triggers/triggers.StartSmtpSink",
		Do: func(c *saga.Context) error {
			hits := &atomic.Int64{}
			lastPtr := new(*SmtpCapturedMessage)

			be := &smtpBackend{mu: &sync.Mutex{}, hits: hits, lastMsg: lastPtr}

			srv := smtp.NewServer(be)
			srv.Addr = constants.SmtpSinkBindAddr
			srv.Domain = "saga.test"
			srv.ReadTimeout = 10 * time.Second
			srv.WriteTimeout = 10 * time.Second
			srv.MaxMessageBytes = 1024 * 1024
			srv.MaxRecipients = 10
			srv.AllowInsecureAuth = true

			ready := make(chan error, 1)
			go func() {
				if err := srv.ListenAndServe(); err != nil && err.Error() != "smtp: server closed" {
					ready <- err
					return
				}
				ready <- nil
			}()

			// Fail fast when the bind address is already taken.
			select {
			case err := <-ready:
				if err != nil {
					return fmt.Errorf("listen %s: %w", constants.SmtpSinkBindAddr, err)
				}
			case <-time.After(150 * time.Millisecond):
			}

			c.Set(BagKeySmtpServer, srv)
			c.Set(BagKeySmtpHits, hits)
			c.Set(BagKeySmtpLastMessage, lastPtr)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			v, ok := c.Get(BagKeySmtpServer)
			if !ok {
				return nil
			}
			srv, ok := v.(*smtp.Server)
			if !ok {
				return nil
			}
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			done := make(chan error, 1)
			go func() { done <- srv.Close() }()
			select {
			case <-done:
			case <-ctx.Done():
			}
			return nil
		},
	}
}
