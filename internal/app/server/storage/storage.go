package storage

type Type int

const (
	SBinary Type = iota
	SCard
	SPassword
	SText
)

type IRepository interface {
	IDataRepository
	ISessionRepository
	IUserRepository
}

type IDataRepository interface {
	GetAllData(uid string) ([]SecureData, error)
	GetAllDataByType(uid string, t Type) ([]SecureData, error)
	GetDataByID(uid, id string) (SecureData, error)
	StoreData(data SecureData) (string, error)
}

type ISessionRepository interface {
	DeleteSession(cid string) error
	GetSession(cid string) (string, error)
	StoreSession(cid, token string) error
}

type IUserRepository interface {
	AddUser(name, pwd string) (User, error)
	GetUserByID(string) (User, error)
	GetUserByName(string) (User, error)
}

type SecureData struct {
	UID  string `json:"-"`
	ID   string `json:"id"`
	Data []byte `json:"data"`
	Type Type   `json:"-"`
}

type User struct {
	ID       string `json:"-"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func NewStorage() IRepository {
	return NewBasicStorage()
}
