package gin

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"

	"github.com/delfimarime/fx/problems"
	"github.com/gin-gonic/gin"
	"schneider.vip/problem"
)

func sendProblem(c *gin.Context, title string, status int, detail string, attr map[string]interface{}) {
	opts := []problem.Option{
		problem.Title(title), problem.Status(status), problem.Detail(detail),
	}
	if attr != nil {
		for k, v := range attr {
			opts = append(opts, problem.Custom(k, v))
		}
	}
	p := problem.New(opts...)
	if _, err := p.WriteTo(c.Writer); err != nil {
		zap.L().Error("Error sending problem:", zap.Error(err))
	}
}

func TryAndCatch(c *gin.Context, f func() error) {
	if err := f(); err != nil {
		switch e := err.(type) {
		case problems.CannotPerformOperationError:
			Send403(c)
		case problems.EntityNotFoundError:
			Send404(c)
		case *problems.JsonValidationError:
			Send400ForSchemaValidation(c, e)
		default:
			Send500(c)
		}
	}
}

func Send403(c *gin.Context) {
	detail := fmt.Sprintf("Cannot perform %s %s", c.Request.Method, c.Request.RequestURI)
	sendProblem(c, "Forbidden", http.StatusForbidden, detail, nil)
}

func Send404(c *gin.Context) {
	detail := fmt.Sprintf("%s %s not found", c.Request.Method, c.Request.RequestURI)
	sendProblem(c, "Not Found", http.StatusNotFound, detail, nil)
}

func Send405(c *gin.Context) {
	detail := fmt.Sprintf("%s %s not allowed", c.Request.Method, c.Request.RequestURI)
	sendProblem(c, "Method not Allowed", http.StatusMethodNotAllowed, detail, nil)
}

func SendHttpNotAcceptable(c *gin.Context) {
	detail := fmt.Sprintf(`Header['Accept']='%s' not allowed`, c.Request.Header.Get("Accept"))
	sendProblem(c, "Not Acceptable", http.StatusNotAcceptable, detail, nil)
}

func SendHttpUnsupportedMediaType(c *gin.Context) {
	detail := fmt.Sprintf(`Header['Content-Type']='%s' not allowed`, c.Request.Header.Get("Content-Type"))
	sendProblem(c, "Unsupported Media Type", http.StatusUnsupportedMediaType, detail, nil)
}

func Send500(c *gin.Context) {
	sendProblem(c, "An error occurred", http.StatusInternalServerError, "Something went wrong!", nil)
}

func Send400ForSchemaValidation(c *gin.Context, e *problems.JsonValidationError) {
	var violations []map[string]any
	for _, each := range e.Problems {
		violations = append(violations, map[string]any{
			"field":   each.Field(),
			"message": each.Description(),
		})
	}
	detail := fmt.Sprintf(`The body isn't compliant. The parameter must be compliant with schema %s`, e.Schema)
	sendProblem(c, "Bad Request", http.StatusBadRequest, detail, map[string]any{"violations": violations})
}

func Send400(c *gin.Context, objectType, name, reason string) {
	detail := fmt.Sprintf(`The %s["%s"] isn't compliant. The parameter must be %s`, objectType, name, reason)
	sendProblem(c, "Bad Request", http.StatusBadRequest, detail, nil)
}

func Send400ForQueryParameter(c *gin.Context, name, reason string) {
	Send400(c, "query", name, reason)
}

func Send400ForHeaderParameter(c *gin.Context, name, reason string) {
	Send400(c, "header", name, reason)
}

func Send400ForBodyParameter(c *gin.Context, name, reason string) {
	Send400(c, "body", name, reason)
}
