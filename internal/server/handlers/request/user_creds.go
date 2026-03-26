package request

const (
	RequestFieldUsername = "username"
	RequestFieldPassword = "password"
)

type UserCredits struct {
	Login    string
	Password string
}

func (o UserCredits) Valid() Problems {
	problems := Problems{
		List: make(map[string]string),
	}

	if o.Login == "" {
		problems.List[RequestFieldUsername] = ErrorMsgEmptyLogin
	}

	if o.Password == "" {
		problems.List[RequestFieldPassword] = ErrorMsgEmptyPassword
	}

	return problems
}
