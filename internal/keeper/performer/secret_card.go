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
	"strconv"
	"strings"
)

type secretCardPerformer struct {
	conn    connector.ServiceConnector
	session session.Session
	logger  *zap.Logger
}

type secretCard struct {
	Number     string `json:"number"`
	CardHolder string `json:"card_holder"`
	CVV        int    `json:"cvv"`
	Expires    string `json:"expires"`
}

func askCard() (secretCard, error) {
	number, err := command.AskText("Enter card number (in format 0000 0000 0000 0000)")
	if err != nil {
		return secretCard{}, fmt.Errorf("cannot read card number: %w", err)
	}

	err = checkNumber(number)
	if err != nil {
		return secretCard{}, fmt.Errorf("invalid card number format: %w", err)
	}

	holder, err := command.AskText("Enter card holder (in format JOHN DOE)")
	if err != nil {
		return secretCard{}, fmt.Errorf("cannot read card holder: %w", err)
	}

	if len(strings.Split(holder, " ")) != 2 {
		return secretCard{}, fmt.Errorf("invalid card holder format: must have first and second names")
	}

	scvv, err := command.AskText("Enter CVV (2-3 digits)")
	if err != nil {
		return secretCard{}, fmt.Errorf("cannot read cvv: %w", err)
	}

	cvv, err := strconv.Atoi(scvv)
	if err != nil {
		return secretCard{}, fmt.Errorf("cvv must be integer: %w", err)
	}

	if cvv > 999 || cvv < 10 {
		return secretCard{}, fmt.Errorf("cvv must be in range 10 <= cvv <= 999")
	}

	expires, err := command.AskText("Enter expires date")
	if err != nil {
		return secretCard{}, fmt.Errorf("cannot read expires date: %w", err)
	}

	return secretCard{
		Number:     number,
		CardHolder: holder,
		CVV:        cvv,
		Expires:    expires,
	}, nil
}

func checkNumber(number string) error {
	// TODO: check Loon alg
	return nil
}

func (p *secretCardPerformer) Set(ctx context.Context, name string) error {
	card, err := askCard()
	if err != nil {
		p.logger.Error("Cannot ask user card", zap.Error(err))

		return fmt.Errorf("cannot ask card: %w", err)
	}

	mCard, err := json.Marshal(card)
	if err != nil {
		p.logger.Error("Cannot marshal card", zap.Error(err))

		return fmt.Errorf("cannot marshal card: %w", err)
	}

	eCard, err := encrypt.EncryptAES256(mCard, p.session.SecretKey)
	if err != nil {
		p.logger.Error("Cannot encrypt card", zap.Error(err))

		return fmt.Errorf("cannot encrypt card: %w", err)
	}

	err = p.conn.SetSecret(ctx, name, secret.SecretTypeCard, eCard)
	if err != nil {
		p.logger.Error("Cannot set card to service", zap.Error(err))

		return fmt.Errorf("cannot set card to service: %w", err)
	}

	fmt.Printf("\033[32mCard %s successfuly created!\033[0m\n", name)

	return nil
}

func (p *secretCardPerformer) Get(ctx context.Context, name string) error {
	s, err := p.conn.GetSecret(ctx, name, secret.SecretTypeCard)
	if err != nil {
		p.logger.Error("Cannot get card from service", zap.Error(err))

		return fmt.Errorf("cannot get card from service: %w", err)
	}

	dSecret, err := encrypt.DecryptAES256(s, p.session.SecretKey)
	if err != nil {
		p.logger.Error("Cannot decrypt card", zap.Error(err))

		return fmt.Errorf("cannot decrypt card: %w", err)
	}

	var uSecret secretCard
	err = json.Unmarshal(dSecret, &uSecret)
	if err != nil {
		p.logger.Error("Cannot unmarshal card", zap.Error(err))

		return fmt.Errorf("cannot unmarshal card: %w", err)
	}

	fmt.Printf("---------\nCard %s\nNumber: %s\nCard Holder: %s\nCVV: %d\nExpires at: %s\n---------\n", name, uSecret.Number, uSecret.CardHolder, uSecret.CVV, uSecret.Expires)

	return nil
}

func (p *secretCardPerformer) Update(ctx context.Context, name string) error {
	card, err := askCard()
	if err != nil {
		p.logger.Error("Cannot ask user card", zap.Error(err))

		return fmt.Errorf("cannot ask card: %w", err)
	}

	mCard, err := json.Marshal(card)
	if err != nil {
		p.logger.Error("Cannot marshal card", zap.Error(err))

		return fmt.Errorf("cannot marshal card: %w", err)
	}

	eCard, err := encrypt.EncryptAES256(mCard, p.session.SecretKey)
	if err != nil {
		p.logger.Error("Cannot encrypt card", zap.Error(err))

		return fmt.Errorf("cannot encrypt card: %w", err)
	}

	err = p.conn.UpdateSecret(ctx, name, secret.SecretTypeCard, eCard)
	if err != nil {
		p.logger.Error("Cannot set card to service", zap.Error(err))

		return fmt.Errorf("cannot set card to service: %w", err)
	}

	fmt.Printf("\033[32mCard %s successfuly updated!\033[0m\n", name)

	return nil
}

func (p *secretCardPerformer) Delete(ctx context.Context, name string) error {
	err := p.conn.RemoveSecret(ctx, name, secret.SecretTypeCard)
	if err != nil {
		p.logger.Error("Cannot remove card from service", zap.Error(err))

		return fmt.Errorf("cannot remove card from service: %w", err)
	}

	fmt.Printf("\033[32mCard %s successfuly removed!\033[0m\n", name)

	return nil
}

func (p *secretCardPerformer) List(ctx context.Context) error {
	secrets, err := p.conn.ListSecret(ctx, secret.SecretTypeCard)
	if err != nil {
		return fmt.Errorf("cannot get list of card: %w", err)
	}

	if len(secrets) == 0 {
		fmt.Printf("No cards found\n")

		return nil
	}

	printPlainSecrets(secrets)
	return nil
}
