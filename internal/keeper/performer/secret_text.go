package performer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
	"github.com/nessai1/gophkeeper/internal/keeper/encrypt"
	"github.com/nessai1/gophkeeper/internal/keeper/secret"
	"github.com/nessai1/gophkeeper/internal/keeper/session"
	"github.com/nessai1/gophkeeper/pkg/command"
	"go.uber.org/zap"
)

type secretTextPerformer struct {
	conn    connector.ServiceConnector
	session session.Session
	logger  *zap.Logger
}

type secretText struct {
	Text string `json:"text"`
}

func askText() (secretText, error) {
	text, err := command.AskText("Enter text")
	if err != nil {
		return secretText{}, fmt.Errorf("cannot read text: %w", err)
	}

	return secretText{
		Text: text,
	}, nil
}

func (p *secretTextPerformer) Set(ctx context.Context, name string) error {
	text, err := askText()
	if err != nil {
		p.logger.Error("Cannot ask user text", zap.Error(err))

		return fmt.Errorf("cannot ask text: %w", err)
	}

	mText, err := json.Marshal(text)
	if err != nil {
		p.logger.Error("Cannot marshal text", zap.Error(err))

		return fmt.Errorf("cannot marshal text: %w", err)
	}

	eText, err := encrypt.EncryptAES256(mText, p.session.SecretKey)
	if err != nil {
		p.logger.Error("Cannot encrypt text", zap.Error(err))

		return fmt.Errorf("cannot encrypt text: %w", err)
	}

	err = p.conn.SetSecret(ctx, name, secret.SecretTypeText, eText)
	if err != nil {
		p.logger.Error("Cannot set text to service", zap.Error(err))

		return fmt.Errorf("cannot set text to service: %w", err)
	}

	fmt.Printf("\033[32mText %s successfuly created!\033[0m\n", name)

	return nil
}

func (p *secretTextPerformer) Get(ctx context.Context, name string) error {
	s, err := p.conn.GetSecret(ctx, name, secret.SecretTypeText)
	if err != nil {
		p.logger.Error("Cannot get text from service", zap.Error(err))

		return fmt.Errorf("cannot get text from service: %w", err)
	}

	dSecret, err := encrypt.DecryptAES256(s, p.session.SecretKey)
	if err != nil {
		p.logger.Error("Cannot decrypt text", zap.Error(err))

		return fmt.Errorf("cannot decrypt text: %w", err)
	}

	var uSecret secretText
	err = json.Unmarshal(dSecret, &uSecret)
	if err != nil {
		p.logger.Error("Cannot unmarshal text", zap.Error(err))

		return fmt.Errorf("cannot unmarshal text: %w", err)
	}

	fmt.Printf("---------\nText %s\n%s\n---------\n", name, uSecret.Text)

	return nil
}

func (p *secretTextPerformer) Update(ctx context.Context, name string) error {
	text, err := askText()
	if err != nil {
		p.logger.Error("Cannot ask user text", zap.Error(err))

		return fmt.Errorf("cannot ask text: %w", err)
	}

	mText, err := json.Marshal(text)
	if err != nil {
		p.logger.Error("Cannot marshal text", zap.Error(err))

		return fmt.Errorf("cannot marshal text: %w", err)
	}

	eText, err := encrypt.EncryptAES256(mText, p.session.SecretKey)
	if err != nil {
		p.logger.Error("Cannot encrypt text", zap.Error(err))

		return fmt.Errorf("cannot encrypt text: %w", err)
	}

	err = p.conn.UpdateSecret(ctx, name, secret.SecretTypeText, eText)
	if err != nil {
		p.logger.Error("Cannot set text to service", zap.Error(err))

		return fmt.Errorf("cannot set text to service: %w", err)
	}

	fmt.Printf("\033[32mText %s successfuly updated!\033[0m\n", name)

	return nil
}

func (p *secretTextPerformer) Delete(ctx context.Context, name string) error {
	err := p.conn.RemoveSecret(ctx, name, secret.SecretTypeText)
	if err != nil {
		p.logger.Error("Cannot remove text from service", zap.Error(err))

		return fmt.Errorf("cannot remove text from service: %w", err)
	}

	fmt.Printf("\033[32mText %s successfuly removed!\033[0m\n", name)

	return nil
}

func (p *secretTextPerformer) List(ctx context.Context) error {
	secrets, err := p.conn.ListSecret(ctx, secret.SecretTypeText)
	if err != nil {
		return fmt.Errorf("cannot get list of text: %w", err)
	}

	if len(secrets) == 0 {
		fmt.Printf("No text found\n")

		return nil
	}

	printPlainSecrets(secrets)
	return nil
}
