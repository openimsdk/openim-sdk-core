package testv3new

type RegisterManager struct {
}

func NewRegisterManager() *RegisterManager {
	return &RegisterManager{}
}

func (r *RegisterManager) RegisterOne(userID string) error {
	return nil
}

func (r *RegisterManager) RegisterBatch(userIDs []string) error {
	return nil
}

func (r *RegisterManager) GetTokens(userIDs ...string) []string {
	return nil
}

func (p *PressureTester) CreateGroup(groupID string, ownerUserID string, userIDs []string) error {
	return nil
}
