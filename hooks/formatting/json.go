package formatting

import (
	"encoding/json"

	"github.com/danielkrainas/csense/api/v1"
)

func JSON(r *v1.Reaction) ([]byte, string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, "", err
	}

	return b, "application/json", nil
}
