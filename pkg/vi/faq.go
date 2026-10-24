package vi

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

type FAQ interface {
	app.Composer

	ID(string) FAQ
	Class(...string) FAQ
	Question(string) FAQQuestion
	Questionf(string, ...any) FAQQuestion
}

func NewFAQ() FAQ {
	return &faq{}
}

type faq struct {
	app.Compo

	IID        string
	IClass     string
	IQuestions []FAQQuestion
}

func (f *faq) ID(v string) FAQ {
	f.IID = v
	return f
}

func (f *faq) Class(v ...string) FAQ {
	f.IClass = app.AppendClass(f.IClass, v...)
	return f
}

func (f *faq) Question(v string) FAQQuestion {
	return FAQQuestion{
		faq:      f,
		question: v,
	}
}

func (f *faq) Questionf(format string, v ...any) FAQQuestion {
	return f.Question(fmt.Sprintf(format, v...))
}

func (f *faq) Render() app.UI {
	return app.Div().
		ID(f.IID).
		Class("vi-faq").
		Class(f.IClass).
		Attr("itemscope", "").
		Attr("itemtype", "https://schema.org/FAQPage").
		Body(
			app.Range(f.IQuestions).Slice(func(i int) app.UI {
				q := f.IQuestions[i]
				return app.Div().
					Attr("itemscope", "").
					Attr("itemprop", "mainEntity").
					Attr("itemtype", "https://schema.org/Question").
					Body(
						app.H3().
							Attr("itemprop", "name").
							Text(q.question),
						app.Div().
							Attr("itemscope", "").
							Attr("itemprop", "acceptedAnswer").
							Attr("itemtype", "https://schema.org/Answer").
							Body(
								app.P().
									Attr("itemprop", "text").
									Body(q.answer...),
							),
					)
			}),
		)
}

type FAQQuestion struct {
	faq      *faq
	question string
	answer   []app.UI
}

func (q FAQQuestion) Answer(v ...app.UI) FAQ {
	q.answer = app.FilterUIElems(v...)
	q.faq.IQuestions = append(q.faq.IQuestions, q)
	return q.faq
}

func (q FAQQuestion) TextAnswer(v string) FAQ {
	return q.Answer(app.Text(v))
}

func (q FAQQuestion) TextAnswerf(format string, v ...any) FAQ {
	return q.TextAnswer(fmt.Sprintf(format, v...))
}
