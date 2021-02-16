package encryptor

type Service interface {
	EncryptPassword(password string) (string, error)
	CompareHashAndPassword(password, hash string) bool
}

type service struct{}

// New creates new encryptor service
func New() Service {
	return service{}
}
