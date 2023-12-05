package validate

import (
	"context"
	"github.com/eliofery/golang-image/pkg/logging"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	ru_translations "github.com/go-playground/validator/v10/translations/ru"
	"log/slog"
)

var (
	enLang = en.New()
	ruLang = ru.New()

	uni = ut.New(enLang, ruLang)
)

func Rus(ctx context.Context) ut.Translator {
	op := "validate.Rus"

	valid := Validation(ctx)
	l := logging.Logging(ctx)

	trans, _ := uni.GetTranslator("ru")

	err := ru_translations.RegisterDefaultTranslations(valid, trans)
	if err != nil {
		// TODO fix error: conflicting key 'required' rule 'Unknown' with text '{0} обязательное поле' for locale 'ru', value being ignored
		l.Info(op, slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
	}

	return trans
}
