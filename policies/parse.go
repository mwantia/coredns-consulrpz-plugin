package policies

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func GetPolicyTargetReader(target string) (io.Reader, error) {
	if strings.HasPrefix(target, "https://") || strings.HasPrefix(target, "http://") {
		response, err := http.Get(target)
		if err != nil {
			return nil, err
		}

		return response.Body, nil
	}

	return os.Open(target)
}
