package model

type UserAuthorization struct {
	Username string
	Password string
}

type UserAuthorizationResponse struct {
	Token string
}

type CreateMeeting struct {
	Name string
}

type Meeting struct {
	ID   string
	Name string
}
