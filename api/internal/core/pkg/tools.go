package pkg

import (
	"encoding/json"
	"io"
	"privet-api/internal/core/models"
)

func DecoderJSON(t string, body io.ReadCloser) (interface{}, error) {
	reqBody, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	if t == "users" {
		var users models.User
		err = json.Unmarshal(reqBody, &users)
		return users, err
	} else {
		var wishes models.WishesList
		err = json.Unmarshal(reqBody, &wishes)
		return wishes, err
	}
}
