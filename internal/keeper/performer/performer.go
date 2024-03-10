package performer

import (
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
	"github.com/nessai1/gophkeeper/internal/keeper/session"
	"go.uber.org/zap"
)

type Sessional interface {
	SetSession(session *session.Session)
	GetSession() *session.Session
}

type Performer interface {
	GetName() string
	GetStruct() string
	GetDescription() string
	GetDetailDescription() string

	// Execute execute command (God please forgive me for this....)
	Execute(
		conn connector.ServiceConnector,
		sessional Sessional,
		logger *zap.Logger,
		args []string,
	) (requireExit bool, err error)
}

var AvailablePerformers = []Performer{
	Help{},
	Exit{},
	Ping{},
	Register{},
	Login{},
	Logout{},
}
