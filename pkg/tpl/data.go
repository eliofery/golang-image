package tpl

import "github.com/eliofery/golang-image/pkg/errors"

type Data struct {
	Meta     Meta
	Data     any
	Messages []any
	Errors   []error
}

type Meta struct {
	Title       string
	Description string
}

/*
PublicErrors
Пример использования:

	data := tpl.Data{
	    Errors: tpl.PublicErrors(
	        "ошибка 1",
	        "ошибка 2",
	        "ошибка 3",
	    ),
	}
*/
func PublicErrors(err ...string) []error {
	var errMsg []error

	for _, e := range err {
		errMsg = append(errMsg, errors.Public(errors.New(e), e))
	}

	return errMsg
}
