package hashivault

type vault struct {
}

func (v *vault) GetSecret(path string) (Secret, error) {
	panic("implement me")
}

func (v *vault) SetDefaultGoogleCredentials(path, key string) error {
	panic("implement me")
}
