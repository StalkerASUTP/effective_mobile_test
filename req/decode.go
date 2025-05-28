package req

import (
	"encoding/json"
	"io"
)

func Decode[T any](body io.ReadCloser) (T, error) {
	var payLoad T
	err := json.NewDecoder(body).Decode(&payLoad)
	if err != nil {
		return payLoad, err
	}
	return payLoad, nil
}