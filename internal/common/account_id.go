package common

import (
	"errors"

	"github.com/labstack/echo/v4"
)

func AccountIDFromEchoCtx(ectx echo.Context) (int, error) {
	accIdVal := ectx.Get("account_id")
	accId, ok := accIdVal.(int)
	if !ok {
		return 0, errors.New("account id is not int")
	}

	return accId, nil
}
