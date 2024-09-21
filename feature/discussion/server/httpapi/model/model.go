package model

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
