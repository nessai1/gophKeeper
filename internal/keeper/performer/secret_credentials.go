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

type secretCredentialsPerformer struct {
	conn    connector.ServiceConnector
	session session.Session
	logger  *zap.Logger
}

type secretCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func askCredentials() (secretCredentials, error) {
	login, err := command.AskText("Enter login")
	if err != nil {
		return secretCredentials{}, fmt.Errorf("cannot read credentials login: %w", err)
	}

	password, err := command.AskText("Enter password")
	if err != nil {
		return secretCredentials{}, fmt.Errorf("cannot read credentials password: %w", err)
	}

	return secretCredentials{
		Login:    login,
		Password: password,
	}, nil
}

func (p *secretCredentialsPerformer) Set(ctx context.Context, name string) error {
	creds, err := askCredentials()
	if err != nil {
		p.logger.Error("Cannot ask user credentials", zap.Error(err))

		return fmt.Errorf("cannot ask credentials: %w", err)
	}

	mCreds, err := json.Marshal(creds)
	if err != nil {
		p.logger.Error("Cannot marshal credentials", zap.Error(err))

		return fmt.Errorf("cannot marshal credentials: %w", err)
	}

	eCreds, err := encrypt.EncryptAES256(mCreds, p.session.SecretKey)
	if err != nil {
		p.logger.Error("Cannot encrypt credentials", zap.Error(err))

		return fmt.Errorf("cannot encrypt credentials: %w", err)
	}

	err = p.conn.SetSecret(ctx, name, secret.SecretTypeCredentials, eCreds)
	if err != nil {
		p.logger.Error("Cannot set credentials to service", zap.Error(err))

		return fmt.Errorf("cannot set credentials to service: %w", err)
	}

	fmt.Printf("\033[32mCredentials %s successfuly created!\033[0m\n", name)

	return nil
}

func (p *secretCredentialsPerformer) Get(ctx context.Context, name string) error {
	s, err := p.conn.GetSecret(ctx, name, secret.SecretTypeCredentials)
	if err != nil {
		p.logger.Error("Cannot get credentials from service", zap.Error(err))

		return fmt.Errorf("cannot get credentials from service: %w", err)
	}

	dSecret, err := encrypt.DecryptAES256(s, p.session.SecretKey)
	if err != nil {
		p.logger.Error("Cannot decrypt credentials", zap.Error(err))

		return fmt.Errorf("cannot decrypt credentials: %w", err)
	}

	var uSecret secretCredentials
	err = json.Unmarshal(dSecret, &uSecret)
	if err != nil {
		p.logger.Error("Cannot unmarshal credentials", zap.Error(err))

		return fmt.Errorf("cannot unmarshal credentials: %w", err)
	}

	fmt.Printf("---------\nCredentials\nname: %s\nLogin: %s\tPassword: %s\n---------\n", name, uSecret.Login, uSecret.Password)

	return nil
}

func (p *secretCredentialsPerformer) Update(ctx context.Context, name string) error {
	creds, err := askCredentials()
	if err != nil {
		p.logger.Error("Cannot ask user credentials", zap.Error(err))

		return fmt.Errorf("cannot ask credentials: %w", err)
	}

	mCreds, err := json.Marshal(creds)
	if err != nil {
		p.logger.Error("Cannot marshal credentials", zap.Error(err))

		return fmt.Errorf("cannot marshal credentials: %w", err)
	}

	eCreds, err := encrypt.EncryptAES256(mCreds, p.session.SecretKey)
	if err != nil {
		p.logger.Error("Cannot encrypt credentials", zap.Error(err))

		return fmt.Errorf("cannot encrypt credentials: %w", err)
	}

	err = p.conn.UpdateSecret(ctx, name, secret.SecretTypeCredentials, eCreds)
	if err != nil {
		p.logger.Error("Cannot set credentials to service", zap.Error(err))

		return fmt.Errorf("cannot set credentials to service: %w", err)
	}

	fmt.Printf("\033[32mCredentials %s successfuly updated!\033[0m\n", name)

	return nil
}

func (p *secretCredentialsPerformer) Delete(ctx context.Context, name string) error {
	err := p.conn.RemoveSecret(ctx, name, secret.SecretTypeCredentials)
	if err != nil {
		p.logger.Error("Cannot remove credentials from service", zap.Error(err))

		return fmt.Errorf("cannot remove credentials from service: %w", err)
	}

	fmt.Printf("\033[32mCredentials %s successfuly removed!\033[0m\n", name)

	return nil
}

func (p *secretCredentialsPerformer) List(ctx context.Context) error {
	secrets, err := p.conn.ListSecret(ctx, secret.SecretTypeCredentials)
	if err != nil {
		return fmt.Errorf("cannot get list of credentials: %w", err)
	}

	if len(secrets) == 0 {
		fmt.Printf("No credentials found\n")

		return nil
	}

	printPlainSecrets(secrets)
	return nil
}
