package util

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func HandleError(err error, resp *resty.Response, msg string) error {
	if err != nil {
		return fmt.Errorf("%s: %v", msg, err)
	}

	if resp.IsError() {
		return HttpError(msg, resp)
	}

	return nil
}

func HttpError(message string, resp *resty.Response) error {
	return fmt.Errorf("%s. Status code %d, response: %s", message, resp.StatusCode(), resp.Body())
}
