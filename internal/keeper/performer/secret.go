package performer

import (
	"context"
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
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
	-- For media it find file in media dir by name and upload it on server

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
		secretName = args[3]
	}

	var performer secretPerformer
	if secretType == SecretTypeMedia {
		performer = &secretMediaPerformer{
			conn:      conn,
			sessional: sessional,
			logger:    logger,
			workDir:   workDir,
		}
	}

	if secretAction == SecretActionSet {
		err := performer.Set(context.TODO(), secretName)
		if err != nil {
			return false, fmt.Errorf("error whle perform %s to type %s: %w", secretAction, secretType, err)
		}
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
