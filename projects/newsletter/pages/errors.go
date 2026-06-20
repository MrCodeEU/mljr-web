package pages

import (
	"errors"

	"mljr-web/ui/feedback"
	"mljr-web/ui/primitive"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// errorPage renders an ApiError as a normal HTML page using the same shell
// as everything else, instead of letting PocketBase's default JSON error
// body reach the browser (which is what happened before this existed).
func errorPage(re *core.RequestEvent, message string) g.Node {
	t := translator(re)
	return appPage(re, "", t("newsletter.errors.title"), nil,
		primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text(t("newsletter.errors.heading"))),
		feedback.Alert(
			feedback.AlertProps{Variant: feedback.AlertDanger, Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
			g.Text(message),
		),
		h.P(h.Style("margin-top:var(--sp-4)"),
			h.A(h.Href("/"), g.Text(t("newsletter.errors.back_link"))),
		),
	)
}

// wrapErrors adapts a route handler so that any *router.ApiError it returns
// (re.ForbiddenError/NotFoundError/BadRequestError/...) renders as the page
// above rather than PocketBase's default raw JSON error body. Handlers that
// already redirect with a flash code, or that return non-ApiError errors,
// are unaffected.
func wrapErrors(handler func(re *core.RequestEvent) error) func(re *core.RequestEvent) error {
	return func(re *core.RequestEvent) error {
		err := handler(re)
		if err == nil {
			return nil
		}
		if apiErr, ok := errors.AsType[*router.ApiError](err); ok {
			return renderPage(re, apiErr.Status, errorPage(re, apiErr.Message))
		}
		return err
	}
}
