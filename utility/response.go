package utility

type JsonResponse struct {
	Success bool  `json:"success"`
	Data    any   `json:"data"`
	Error   error `json:"error"`
}
