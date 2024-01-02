package responses

const (
	ErrInvalidCredentials  = "invalid username or password"
	ErrUnauthorized        = "authentication error"
	ErrInternalServerError = "internal server error"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type DataResponse[T any] struct {
	Data T `json:"data"`
}

func NewDataResponse[T any](data T) DataResponse[T] {
	return DataResponse[T]{
		Data: data,
	}
}

type Success struct {
	Success bool `json:"success"`
}

func NewSuccessResponse(success bool) DataResponse[Success] {
	return NewDataResponse(Success{Success: success})
}
