package tpl

import (
	"github.com/eliofery/golang-image/internal/app/models/user"
	"github.com/eliofery/golang-image/pkg/errors"
	"github.com/eliofery/golang-image/pkg/validate"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
	"time"
)

type funcTemplate any

var (
	funcMap = template.FuncMap{
		"csrfInput":   csrfInput,
		"errors":      errorsMsg,
		"error":       errorMsg,
		"currentUser": currentUser,
		"now":         timeNow,
	}
)

func csrfInput(r *http.Request, _ Data) funcTemplate {
	return func() template.HTML {
		return csrf.TemplateField(r)
	}
}

/*
Пример использования:
{{if errors}}
<ul>

	{{range errors}}
	<li>{{.}}</li>
	{{end}}

</ul>
{{end}}
*/
func errorsMsg(r *http.Request, data Data) funcTemplate {
	var (
		ErrSomeWrong = errors.New("Непредвиденная ошибка на стороне сервера")

		errMessages  []string
		pubErr       errors.PublicError
		validatorErr validator.ValidationErrors
	)

	trans := validate.Rus(r.Context())

	for _, err := range data.Errors {
		switch {
		case errors.As(err, &validatorErr):
			for _, err := range validatorErr {
				errMessages = append(errMessages, err.Translate(trans))
			}
		case errors.As(err, &pubErr):
			errMessages = append(errMessages, pubErr.Public())
		default:
			errMessages = append(errMessages, ErrSomeWrong.Error())
		}
	}

	return func() []string {
		return errMessages
	}
}

/*
Пример использования:

	type Dto struct {
	    Email    string `validate:"required,email"`
	}

<input id="email"r type="text" name="email" placeholder="Введите ваш email">
{{if error .Errors "Email"}}

	<p>{{error .Errors "Email"}}</p>

{{end}}
*/
func errorMsg(r *http.Request, _ Data) funcTemplate {
	var (
		validatorErr validator.ValidationErrors
	)

	errMessages := map[string]string{}
	trans := validate.Rus(r.Context())

	return func(errs []error, key string) string {
		for _, err := range errs {
			if errors.As(err, &validatorErr) {
				for _, err := range validatorErr {
					errMessages[err.Field()] = err.Translate(trans)
				}
			}
		}

		return errMessages[key]
	}
}

func currentUser(r *http.Request, _ Data) funcTemplate {
	return func() *user.User {
		return user.CtxUser(r.Context())
	}
}

func timeNow(_ *http.Request, _ Data) funcTemplate {
	return func() time.Time {
		return time.Now()
	}
}
