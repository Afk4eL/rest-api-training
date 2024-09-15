package utilhttp

import (
	"io"
	"rest-arch-training/internal/server/utils/validator"

	"github.com/go-chi/render"
)

func ReadRequest(body io.Reader, s any) error {
	if err := render.DecodeJSON(body, s); err != nil {
		return err
	}

	return validator.ValidateStruct(s)
}
