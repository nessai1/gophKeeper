package performer

import (
	"context"
	"fmt"
	"github.com/chrusty/go-tableprinter"
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
	"github.com/nessai1/gophkeeper/internal/keeper/secret"
	"go.uber.org/zap"
	"strings"
)

const (
	SecretTypeCredentials = "credentials"
	SecretTypeCard        = "card"
	SecretTypeText        = "text"
	SecretTypeMedia       = "media"
)

const (
	SecretActionSet    = "set"
	SecretActionGet    = "get"
	SecretActionUpdate = "update"
	SecretActionDelete = "delete"
	SecretActionList   = "list"
)

type Secret struct {
}

func (p Secret) GetName() string {
	return "secret"
}

func (p Secret) GetStruct() string {
	return "secret [type] [action] [?name]"
}

func (p Secret) GetDescription() string {
	return "Manipulate with user secrets (credentials, credit card, text, media)"
}

func (p Secret) GetDetailDescription() string {
	return `Mainpulate with user secrets

Available types:

- credentials - credentials of user: his login and password
- card - credit card info: card-holder, cvv, number, expires date
- text - some secret text
- media - some blob files, that contains in keeperData/media dir (keeperData can be renamed in config by 'work_dir' field, check it before

Available actions:

- set - start procedure of create new entity of selected type
	-- For credentials it ask login&password and save it by name
	-- For card it ask card-colder&cvv&number&date and save it by name
	-- For text it ask text and save it by name
	-- For media it find file with path == name and upload it on server

- get - get data of secrets
	-- For credentials/card/text it show information about secret by name
	-- For media it load file from server and save in media dir

- remove - remove secret by name

- update - update secret by name: remove old secret and start 'set' procedure again

- list - get list of secrets
`
}

func (p Secret) Execute(conn connector.ServiceConnector, sessional Sessional, logger *zap.Logger, args []string, workDir string) (requireExit bool, err error) {
	if sessional.GetSession() == nil {
		return false, fmt.Errorf("for working with secrets you need to be authorized")
	}

	err = validateArguments(args)
	if err != nil {
		return false, fmt.Errorf("got invalid arguments for secret command: %w", err)
	}

	var (
		secretType   = args[1]
		secretAction = args[2]
		secretName   = ""
	)

	if secretAction != SecretActionList {
		if len(args) < 4 {
			return false, fmt.Errorf("mismatch arguments count for secret actions: must be 4")
		}
		secretName = args[3]
	}

	var performer secretPerformer
	if secretType == SecretTypeMedia {
		performer = &secretMediaPerformer{
			conn:    conn,
			session: *sessional.GetSession(),
			logger:  logger,
			workDir: workDir,
		}
	} else if secretType == SecretTypeCredentials {
		performer = &secretCredentialsPerformer{
			conn:    conn,
			session: *sessional.GetSession(),
			logger:  logger,
		}
	} else if secretType == SecretTypeText {
		performer = &secretTextPerformer{
			conn:    conn,
			session: *sessional.GetSession(),
			logger:  logger,
		}
	} else if secretType == SecretTypeCard {
		performer = &secretCardPerformer{
			conn:    conn,
			session: *sessional.GetSession(),
			logger:  logger,
		}
	} else {
		return false, fmt.Errorf("invalid secret type: %s", secretType)
	}

	ctx := context.TODO()
	err = nil
	switch secretAction {
	case SecretActionSet:
		err = performer.Set(ctx, secretName)
	case SecretActionGet:
		err = performer.Get(ctx, secretName)
	case SecretActionList:
		err = performer.List(ctx)
	case SecretActionUpdate:
		err = performer.Update(ctx, secretName)
	case SecretActionDelete:
		err = performer.Delete(ctx, secretName)
	default:
		err = fmt.Errorf("undefined action occured: %s", secretAction)
	}

	if err != nil {
		return false, fmt.Errorf("error whle perform %s to type %s: %w", secretAction, secretType, err)
	}

	return false, nil
}

func validateArguments(args []string) error {
	if len(args) > 4 {
		return fmt.Errorf("too many arguments: %d", len(args)-1)
	}

	if err := validateType(args[1]); err != nil {
		return fmt.Errorf("invalid secret arguments: %w", err)
	}

	if err := validateAction(args[2]); err != nil {
		return fmt.Errorf("invalid secret arguments: %w", err)
	}

	if len(args) == 3 {
		return nil
	}

	if strings.TrimSpace(args[3]) == "" && args[1] != SecretActionList {
		return fmt.Errorf("secret name can't be empty")
	}

	return nil
}

func validateType(secretType string) error {
	if secretType != SecretTypeCredentials && secretType != SecretTypeCard && secretType != SecretTypeText && secretType != SecretTypeMedia {
		return fmt.Errorf("invalid secret type assigned: %s", secretType)
	}

	return nil
}

func validateAction(secretAction string) error {
	if secretAction != SecretActionSet && secretAction != SecretActionGet && secretAction != SecretActionUpdate && secretAction != SecretActionDelete && secretAction != SecretActionList {
		return fmt.Errorf("invalid secret action assigned: %s", secretAction)
	}

	return nil
}

type secretPerformer interface {
	Set(ctx context.Context, name string) error
	Get(ctx context.Context, name string) error
	Update(ctx context.Context, name string) error
	Delete(ctx context.Context, name string) error
	List(ctx context.Context) error
}

type printableSecret struct {
	Name        string
	Create_time string // snake case used for table formatter
	Update_time string
}

func printSecrets(secrets []printableSecret) {
	tableprinter.SetBorder(true)
	tableprinter.Print(secrets)
}

func printPlainSecrets(secrets []secret.Secret) {
	printable := make([]printableSecret, len(secrets))
	for i, v := range secrets {
		printable[i] = printableSecret{
			Name:        v.Name,
			Create_time: v.Created.String(),
			Update_time: v.Updated.String(),
		}
	}

	printSecrets(printable)
}
