package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

type data struct {
	Name          string
	Package       string
	TypeName      string
	ServiceID     string
	BinaryName    string
	ModuleName    string
	RouteBase     string
	FrameworkRoot string
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) < 1 || args[0] != "new" {
		fmt.Fprintln(stderr, "usage: chenile-servicegen new --name <service-name> [--out <dir>] [--framework-root <dir>]")
		return 2
	}
	flags := flag.NewFlagSet("new", flag.ContinueOnError)
	flags.SetOutput(stderr)
	name := flags.String("name", "", "service name")
	out := flags.String("out", ".", "output directory")
	frameworkRoot := flags.String("framework-root", "../../chenile-framework", "relative path from generated service to the framework root")
	if err := flags.Parse(args[1:]); err != nil {
		return 2
	}
	if *name == "" {
		fmt.Fprintln(stderr, "--name is required")
		return 2
	}
	d := derive(*name, *frameworkRoot)
	target := filepath.Join(*out, d.BinaryName)
	if err := generate(target, d); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := updateWorkspace(target, d); err != nil {
		fmt.Fprintf(stderr, "warning: could not update go.work: %v\n", err)
	}
	fmt.Fprintf(stdout, "created %s\n", target)
	fmt.Fprintf(stdout, "next:\n  cd %s\n  go test ./...\n  go run ./cmd/%s\n", target, d.BinaryName)
	return 0
}

func derive(name string, frameworkRoot string) data {
	parts := splitName(name)
	typeName := ""
	for _, part := range parts {
		typeName += strings.Title(part)
	}
	pkg := strings.Join(parts, "")
	kebab := strings.Join(parts, "-")
	return data{
		Name:          name,
		Package:       pkg,
		TypeName:      typeName,
		ServiceID:     pkg + "Service",
		BinaryName:    kebab + "-service",
		ModuleName:    kebab + "-service",
		RouteBase:     "/" + kebab + "s",
		FrameworkRoot: frameworkRoot,
	}
}

func splitName(name string) []string {
	fields := strings.FieldsFunc(strings.ToLower(name), func(r rune) bool {
		return r == '-' || r == '_' || unicode.IsSpace(r)
	})
	clean := []string{}
	for _, field := range fields {
		if field != "" {
			clean = append(clean, field)
		}
	}
	if len(clean) == 0 {
		return []string{"service"}
	}
	return clean
}

func generate(target string, d data) error {
	files := map[string]string{
		"go.mod":                                      goModTemplate,
		"cmd/" + d.BinaryName + "/main.go":            mainTemplate,
		d.Package + "/model.go":                       modelTemplate,
		d.Package + "/service.go":                     serviceTemplate,
		d.Package + "/controller.go":                  controllerTemplate,
		d.Package + "/module.go":                      moduleTemplate,
		d.Package + "/service_test.go":                serviceUnitTestTemplate,
		d.Package + "/controller_test.go":             controllerUnitTestTemplate,
		"test/" + d.Package + "_service_test.go":      testTemplate,
		"test/features/" + d.Package + ".feature":     featureTemplate,
		"test/fixtures/create_" + d.Package + ".json": fixtureTemplate,
	}
	for path, source := range files {
		if err := render(filepath.Join(target, path), source, d); err != nil {
			return err
		}
	}
	return nil
}

func render(path string, source string, d data) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	tpl, err := template.New(filepath.Base(path)).Parse(source)
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	if err := tpl.Execute(&buffer, d); err != nil {
		return err
	}
	return os.WriteFile(path, buffer.Bytes(), 0644)
}

func updateWorkspace(target string, d data) error {
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return err
	}
	workspaceRoot := filepath.Clean(filepath.Join(absTarget, d.FrameworkRoot))
	workFile := filepath.Join(workspaceRoot, "go.work")
	content, err := os.ReadFile(workFile)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(workspaceRoot, absTarget)
	if err != nil {
		return err
	}
	rel = "./" + filepath.ToSlash(rel)
	text := string(content)
	if strings.Contains(text, rel) {
		return nil
	}
	line := "\t" + rel + "\n"
	if index := strings.LastIndex(text, ")"); index >= 0 {
		text = text[:index] + line + text[index:]
	} else {
		text += "\nuse " + rel + "\n"
	}
	return os.WriteFile(workFile, []byte(text), 0644)
}

const goModTemplate = `module {{.ModuleName}}

go 1.22

require (
	core v0.0.0
	http v0.0.0
	packager v0.0.0
	test v0.0.0
)

replace base => {{.FrameworkRoot}}/base
replace core => {{.FrameworkRoot}}/core
replace http => {{.FrameworkRoot}}/http
replace owiz => {{.FrameworkRoot}}/owiz
replace packager => {{.FrameworkRoot}}/packager
replace test => {{.FrameworkRoot}}/test
`

const mainTemplate = `package main

import (
	"log"

	"packager"

	"{{.ModuleName}}/{{.Package}}"
)

func main() {
	app, err := packager.NewWebApp(packager.Module{Name: "{{.Package}}", Register: {{.Package}}.Register})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening on :8080")
	log.Fatal(app.ListenAndServe(":8080"))
}
`

const modelTemplate = `package {{.Package}}

type Create{{.TypeName}}Request struct {
	Name string ` + "`json:\"name\"`" + `
}

type {{.TypeName}} struct {
	ID   string ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`

const serviceTemplate = `package {{.Package}}

import "context"

type Service interface {
	Create(context.Context, Create{{.TypeName}}Request) ({{.TypeName}}, error)
}

type service struct{}

func NewService() Service {
	return service{}
}

func (service) Create(ctx context.Context, request Create{{.TypeName}}Request) ({{.TypeName}}, error) {
	return {{.TypeName}}{
		ID:   "{{.Package}}-1",
		Name: request.Name,
	}, nil
}
`

const controllerTemplate = `package {{.Package}}

import (
	"context"
	"net/http"

	"core"
)

func Register(registry *core.Registry) error {
	service := NewService()
	return registry.RegisterService(core.ServiceDefinition{
		ID:   "{{.ServiceID}}",
		Name: "{{.ServiceID}}",
		Operations: []core.OperationDefinition{
			{
				Name:   "create",
				Method: http.MethodPost,
				Path:   "{{.RouteBase}}",
				NewInput: func() any {
					return &Create{{.TypeName}}Request{}
				},
				Handler: func(ctx context.Context, exchange *core.Exchange) (any, error) {
					request := exchange.Input.(*Create{{.TypeName}}Request)
					return service.Create(ctx, *request)
				},
			},
		},
	})
}
`

const moduleTemplate = `package {{.Package}}

// Package {{.Package}} is a generated Chenile service module.
`

const serviceUnitTestTemplate = `package {{.Package}}

import (
	"context"
	"testing"
)

func TestServiceCreates{{.TypeName}}(t *testing.T) {
	service := NewService()

	{{.Package}}, err := service.Create(context.Background(), Create{{.TypeName}}Request{Name: "Alice"})
	if err != nil {
		t.Fatal(err)
	}
	if {{.Package}}.ID == "" {
		t.Fatal("expected generated id")
	}
	if {{.Package}}.Name != "Alice" {
		t.Fatalf("expected Alice, got %q", {{.Package}}.Name)
	}
}
`

const controllerUnitTestTemplate = `package {{.Package}}

import (
	"context"
	"testing"

	"core"
)

func TestRegisterAddsCreateOperation(t *testing.T) {
	registry := core.NewRegistry()

	if err := Register(registry); err != nil {
		t.Fatal(err)
	}

	operation, ok := registry.Operation("{{.ServiceID}}", "create")
	if !ok {
		t.Fatal("expected create operation")
	}
	if operation.Path != "{{.RouteBase}}" {
		t.Fatalf("expected {{.RouteBase}}, got %q", operation.Path)
	}
	if operation.NewInput == nil {
		t.Fatal("expected input factory")
	}
}

func TestRegisteredCreateHandlerInvokesService(t *testing.T) {
	registry := core.NewRegistry()
	if err := Register(registry); err != nil {
		t.Fatal(err)
	}
	operation, ok := registry.Operation("{{.ServiceID}}", "create")
	if !ok {
		t.Fatal("expected create operation")
	}

	payload, err := operation.Handler(context.Background(), &core.Exchange{
		Input: &Create{{.TypeName}}Request{Name: "Alice"},
	})
	if err != nil {
		t.Fatal(err)
	}
	{{.Package}} := payload.({{.TypeName}})
	if {{.Package}}.Name != "Alice" {
		t.Fatalf("expected Alice, got %q", {{.Package}}.Name)
	}
}
`

const testTemplate = `package test

import (
	"io"
	"testing"

	"packager"
	godogtest "test/godog"

	"{{.ModuleName}}/{{.Package}}"
)

func TestCreate{{.TypeName}}(t *testing.T) {
	app, err := packager.NewWebApp(packager.Module{Name: "{{.Package}}", Register: {{.Package}}.Register})
	if err != nil {
		t.Fatal(err)
	}

	status := godogtest.Suite{
		Name:         "{{.Package}}-service",
		Router:       app.Router,
		FeaturePaths: []string{"features/{{.Package}}.feature"},
		TestingT:     t,
		Output:       io.Discard,
	}.Run()
	if status != 0 {
		t.Fatalf("godog suite failed with status %d", status)
	}
}
`

const featureTemplate = `Feature: {{.TypeName}} service

  Scenario: Create {{.Package}}
    When I POST a REST request to URL "{{.RouteBase}}" with payload
      """
      {
        "name": "Alice"
      }
      """
    Then the http status code is 200
    And success is true
    And the REST response key "name" is "Alice"
`

const fixtureTemplate = `{
  "name": "Alice"
}
`
