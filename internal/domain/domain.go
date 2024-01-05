package domain

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
}

type NoteRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type NoteResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedBy string `json:"created_by"`
}

type SharedNoteRequest struct {
	ToUsersID []string `json:"to_users_id"`
}
