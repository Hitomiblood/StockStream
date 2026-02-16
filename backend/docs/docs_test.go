package docs

import (
	"strings"
	"testing"

	"github.com/swaggo/swag"
)

func TestSwaggerInfoIsRegistered(t *testing.T) {
	if SwaggerInfo == nil {
		t.Fatalf("SwaggerInfo should not be nil")
	}

	registered := swag.GetSwagger(SwaggerInfo.InstanceName())
	if registered == nil {
		t.Fatalf("swagger spec should be registered")
	}
}

func TestSwaggerTemplateContainsAPIPaths(t *testing.T) {
	doc := SwaggerInfo.ReadDoc()
	if !strings.Contains(doc, "/api/v1/stocks") {
		t.Fatalf("swagger doc missing expected path")
	}
}
