package tpl

import (
	"bytes"
	"context"
	"fmt"
	"github.com/eliofery/golang-image/internal/resources"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/ydb-platform/ydb-go-sdk/v3/log"
	"html/template"
	"io"
	"path"
)

type Tpl struct {
	layout string
	page   string
	parts  []string
}

func New(page string) *Tpl {
	parts, err := getParts()
	if err != nil {
		log.Error(err)
		parts = []string{}
	}

	return &Tpl{
		layout: getLayout(layoutDefault),
		page:   getPage(page),
		parts:  parts,
	}
}

func (t *Tpl) SetLayout(layout string) *Tpl {
	return &Tpl{
		layout: getLayout(layout),
		page:   t.page,
		parts:  t.parts,
	}
}

func (t *Tpl) Render(ctx context.Context, data Data) error {
	op := "tpl.Render"

	w := router.ResponseWriter(ctx)
	r := router.Request(ctx)

	layoutFileName := path.Base(t.getAllFiles()[0])
	tpl := template.New(layoutFileName)
	tpl = tpl.Funcs(getFuncMap(r, data))

	tpl, err := tpl.ParseFS(resources.Views, t.getAllFiles()...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = io.Copy(w, &buf); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func Render(ctx context.Context, page string, data Data) error {
	return New(page).Render(ctx, data)
}
