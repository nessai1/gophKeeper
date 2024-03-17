package session

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/encrypt"
	"github.com/nessai1/gophkeeper/pkg/command"
	"os"
	"path/filepath"
	"strings"
)

type Session struct {
	// Login login of authed user
	Login string
	// AuthToken token for auth on external service
	AuthToken string
	// SecretKey key for decrypt user secrets
	SecretKey [32]byte

	// passwordHash hash for confirm next user session in keeper
	passwordHash string
}

func NewSession(login, password, serviceToken string) Session {
	return Session{
		Login:        login,
		AuthToken:    serviceToken,
		SecretKey:    encrypt.BuildAESKey(login, password),
		passwordHash: hashPassword(password),
	}
}

func hashPassword(password string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password+"keepersession123")))
}

const userDataFilename = "userdata.json"

type userData struct {
	Login        string `json:"login"`
	PasswordHash string `json:"password"`
	ServerToken  string `json:"server_token"`
}

func LoadLocalSession(workDir string) (*Session, error) {
	file, err := os.Open(filepath.Join(workDir, userDataFilename))
	if err != nil {
		return nil, fmt.Errorf("cannot open session file: %w", err)
	}
	defer file.Close()

	b := bytes.Buffer{}
	n, err := b.ReadFrom(file)

	if n == 0 {
		return nil, fmt.Errorf("cannot open session file: is empty")
	} else if err != nil {
		return nil, fmt.Errorf("cannot read session file: %w", err)
	}

	defer os.Remove(file.Name())
	var ud userData

	err = json.Unmarshal(b.Bytes(), &ud)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal session file")
	}

	password, err := command.AskSecret(fmt.Sprintf("Enter %s password (or leave it blank for new session)", ud.Login))
	if strings.TrimSpace(password) == "" {
		return nil, nil
	}

	if hashPassword(password) != ud.PasswordHash {
		fmt.Printf("\033[31mIncorrect password! Session droped\033[0m\n")
		return nil, fmt.Errorf("user enter incorrect password for existing session")
	}

	session := Session{
		Login:        ud.Login,
		AuthToken:    ud.ServerToken,
		SecretKey:    encrypt.BuildAESKey(ud.Login, password),
		passwordHash: ud.PasswordHash,
	}

	return &session, nil
}

func SaveLocalSession(workDir string, session Session) error {
	ud := userData{
		Login:        session.Login,
		PasswordHash: session.passwordHash,
		ServerToken:  session.AuthToken,
	}

	b, err := json.Marshal(ud)
	if err != nil {
		return fmt.Errorf("cannot marshal session: %w", err)
	}

	file, err := os.OpenFile(filepath.Join(workDir, userDataFilename), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("cannot open user data file: %w", err)
	}

	_, err = file.Write(b)
	if err != nil {
		return fmt.Errorf("cannot write user data to file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("cannot close user data file: %w", err)
	}

	return nil
}
