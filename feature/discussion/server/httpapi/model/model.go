package model

type UserAuthorization struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserAuthorizationResponse struct {
	Token string `json:"token"`
}

type CreateMeeting struct {
	Title    string `json:"title"`
	Password string `json:"password"`
}

type Meeting struct {
	ID    string `json:"id"`
	Title string `json:"name"`
}
