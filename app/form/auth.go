package form

type Register struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required"`
}

type Login struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required"`
}
