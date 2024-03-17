package secret

import "time"

type SecretType int

const (
	SecretTypeCredentials SecretType = iota
	SecretTypeCard
	SecretTypeText
	SecretTypeMedia
)

type Secret struct {
	SecretType SecretType
	Name       string
	Created    time.Time
	Updated    time.Time
}
