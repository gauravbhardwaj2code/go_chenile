package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

type data struct {
	Name            string
	Package         string
	TypeName        string
	ServiceID       string
	BinaryName      string
	ModuleName      string
	RouteBase       string
	FrameworkRoot   string
	UseLocalReplace bool
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
	publicDeps := flags.Bool("public-deps", false, "generate public Chenile module dependencies without local replace directives")
	if err := flags.Parse(args[1:]); err != nil {
		return 2
	}
	if *name == "" {
		fmt.Fprintln(stderr, "--name is required")
		return 2
	}
	d := derive(*name, *frameworkRoot)
	d.UseLocalReplace = !*publicDeps
	target := filepath.Join(*out, d.BinaryName)
	if err := generate(target, d); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := tidyGeneratedModule(target, d); err != nil {
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

func tidyGeneratedModule(target string, d data) error {
	if !d.UseLocalReplace {
		return nil
	}
	frameworkRoot := filepath.Clean(filepath.Join(target, d.FrameworkRoot))
	if _, err := os.Stat(filepath.Join(frameworkRoot, "bdd-utils", "go.mod")); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = target
	cmd.Env = append(os.Environ(), "GOWORK=off")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go mod tidy generated service: %w\n%s", err, string(output))
	}
	return nil
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
		Name:            name,
		Package:         pkg,
		TypeName:        typeName,
		ServiceID:       pkg + "Service",
		BinaryName:      kebab + "-service",
		ModuleName:      kebab + "-service",
		RouteBase:       "/" + kebab + "s",
		FrameworkRoot:   frameworkRoot,
		UseLocalReplace: true,
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
		"go.mod":                                                   goModTemplate,
		"cmd/" + d.BinaryName + "/main.go":                         mainTemplate,
		d.Package + "/contract/request.go":                         contractRequestTemplate,
		d.Package + "/contract/response.go":                        contractResponseTemplate,
		d.Package + "/contract/controller.go":                      contractControllerTemplate,
		d.Package + "/contract/controller_test.go":                 contractControllerTestTemplate,
		d.Package + "/domain/model.go":                             domainModelTemplate,
		d.Package + "/domain/errors.go":                            domainErrorsTemplate,
		d.Package + "/repository/repository.go":                    repositoryTemplate,
		d.Package + "/repository/memory_repository.go":             memoryRepositoryTemplate,
		d.Package + "/service/service.go":                          serviceTemplate,
		d.Package + "/service/service_test.go":                     serviceTestTemplate,
		d.Package + "/module/module.go":                            moduleTemplate,
		"config/application.yaml":                                  configTemplate,
		"test/" + d.Package + "_service_test.go":                   testTemplate,
		"test/features/" + d.Package + ".feature":                  featureTemplate,
		"test/fixtures/create_" + d.Package + ".json":              fixtureTemplate,
		"test/fixtures/create_" + d.Package + "_missing_name.json": invalidFixtureTemplate,
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

go 1.26

toolchain go1.26.3

require (
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/bdd-utils v0.1.0
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/base v0.1.0
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/chenile v0.1.0
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/config v0.1.0
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/packager v0.1.0
)

require (
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/core v0.1.0 // indirect
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/http v0.1.0 // indirect
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/owiz v0.1.0 // indirect
	github.com/cucumber/gherkin/go/v26 v26.2.0 // indirect
	github.com/cucumber/godog v0.15.1 // indirect
	github.com/cucumber/messages/go/v21 v21.0.1 // indirect
	github.com/gofrs/uuid v4.3.1+incompatible // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-memdb v1.3.4 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/spf13/pflag v1.0.7 // indirect
)

{{if .UseLocalReplace}}replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/bdd-utils => {{.FrameworkRoot}}/bdd-utils
replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/base => {{.FrameworkRoot}}/base
replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/chenile => {{.FrameworkRoot}}/chenile
replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/config => {{.FrameworkRoot}}/config
replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/core => {{.FrameworkRoot}}/core
replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/http => {{.FrameworkRoot}}/http
replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/owiz => {{.FrameworkRoot}}/owiz
replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/packager => {{.FrameworkRoot}}/packager
{{end}}
`

const mainTemplate = `package main

import (
	"log"

	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/config"
	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/packager"

	"{{.ModuleName}}/{{.Package}}/module"
)

func main() {
	cfg, err := config.Load("config/application.yaml")
	if err != nil {
		log.Fatal(err)
	}
	app, err := packager.NewChenileWebApp(module.New())
	if err != nil {
		log.Fatal(err)
	}
	address := ":" + cfg.String("server.port", "8080")
	log.Println("listening on " + address)
	log.Fatal(app.ListenAndServe(address))
}
`

const contractRequestTemplate = `package contract

type Create{{.TypeName}}Request struct {
	Name string ` + "`json:\"name\"`" + `
}
`

const contractResponseTemplate = `package contract

import "{{.ModuleName}}/{{.Package}}/domain"

type {{.TypeName}} struct {
	ID   string ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}

func New{{.TypeName}}Response(entity domain.{{.TypeName}}) {{.TypeName}} {
	return {{.TypeName}}{
		ID:   entity.ID,
		Name: entity.Name,
	}
}
`

const contractControllerTemplate = `package contract

import (
	"context"

	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/chenile"

	"{{.ModuleName}}/{{.Package}}/domain"
	servicepkg "{{.ModuleName}}/{{.Package}}/service"
)

type Controller struct {
	service servicepkg.Service
}

func NewController(service servicepkg.Service) Controller {
	return Controller{service: service}
}

func Routes(service servicepkg.Service) []chenile.Route {
	return NewController(service).Routes()
}

func (c Controller) Routes() []chenile.Route {
	return []chenile.Route{
		chenile.POST("{{.RouteBase}}", "create", func() *Create{{.TypeName}}Request {
			return &Create{{.TypeName}}Request{}
		}, c.Create),
	}
}

func (c Controller) Create(ctx context.Context, request Create{{.TypeName}}Request) ({{.TypeName}}, error) {
	entity, err := c.service.Create(ctx, domain.Create{{.TypeName}}Command{Name: request.Name})
	if err != nil {
		return {{.TypeName}}{}, err
	}
	return New{{.TypeName}}Response(entity), nil
}
`

const contractControllerTestTemplate = `package contract

import (
	"context"
	"testing"

	"{{.ModuleName}}/{{.Package}}/domain"
)

type fakeService struct{}

func (fakeService) Create(ctx context.Context, command domain.Create{{.TypeName}}Command) (domain.{{.TypeName}}, error) {
	return domain.{{.TypeName}}{ID: "{{.Package}}-1", Name: command.Name}, nil
}

func TestRoutesDeclareCreateOperation(t *testing.T) {
	routes := Routes(fakeService{})

	if len(routes) != 1 {
		t.Fatalf("expected one route, got %d", len(routes))
	}
	if routes[0].Name != "create" {
		t.Fatalf("expected create operation, got %q", routes[0].Name)
	}
	if routes[0].Path != "{{.RouteBase}}" {
		t.Fatalf("expected {{.RouteBase}}, got %q", routes[0].Path)
	}
	if routes[0].NewInput == nil {
		t.Fatal("expected input factory")
	}
}

func TestCreateInvokesService(t *testing.T) {
	controller := NewController(fakeService{})

	payload, err := controller.Create(context.Background(), Create{{.TypeName}}Request{Name: "Alice"})
	if err != nil {
		t.Fatal(err)
	}
	if payload.Name != "Alice" {
		t.Fatalf("expected Alice, got %q", payload.Name)
	}
}
`

const domainModelTemplate = `package domain

type Create{{.TypeName}}Command struct {
	Name string
}

type {{.TypeName}} struct {
	ID   string
	Name string
}
`

const domainErrorsTemplate = `package domain

import (
	"net/http"

	chenileerrors "github.com/gauravbhardwaj2code/go_chenile/chenile-framework/base/errors"
)

const (
	Error{{.TypeName}}NameRequired = 1001
)

func NameRequired() error {
	return chenileerrors.Builder().
		Status(http.StatusBadRequest).
		Code(Error{{.TypeName}}NameRequired).
		MessageKey("{{.Package}}.name.required").
		Description("name is required").
		Field("name", "required", "name is required").
		Build()
}
`

const repositoryTemplate = `package repository

import (
	"context"

	"{{.ModuleName}}/{{.Package}}/domain"
)

type Repository interface {
	Create(context.Context, domain.{{.TypeName}}) (domain.{{.TypeName}}, error)
}
`

const memoryRepositoryTemplate = `package repository

import (
	"context"
	"fmt"
	"sync"

	"{{.ModuleName}}/{{.Package}}/domain"
)

type MemoryRepository struct {
	mu     sync.Mutex
	nextID int
	values map[string]domain.{{.TypeName}}
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextID: 1,
		values: map[string]domain.{{.TypeName}}{},
	}
}

func (r *MemoryRepository) Create(ctx context.Context, entity domain.{{.TypeName}}) (domain.{{.TypeName}}, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if entity.ID == "" {
		entity.ID = fmt.Sprintf("{{.Package}}-%d", r.nextID)
		r.nextID++
	}
	r.values[entity.ID] = entity
	return entity, nil
}
`

const serviceTemplate = `package service

import (
	"context"

	"{{.ModuleName}}/{{.Package}}/domain"
	"{{.ModuleName}}/{{.Package}}/repository"
)

type Service interface {
	Create(context.Context, domain.Create{{.TypeName}}Command) (domain.{{.TypeName}}, error)
}

type service struct {
	repository repository.Repository
}

func New(repository repository.Repository) Service {
	return service{repository: repository}
}

func (s service) Create(ctx context.Context, command domain.Create{{.TypeName}}Command) (domain.{{.TypeName}}, error) {
	if command.Name == "" {
		return domain.{{.TypeName}}{}, domain.NameRequired()
	}
	return s.repository.Create(ctx, domain.{{.TypeName}}{
		Name: command.Name,
	})
}
`

const serviceTestTemplate = `package service

import (
	"context"
	"testing"

	"{{.ModuleName}}/{{.Package}}/domain"
	"{{.ModuleName}}/{{.Package}}/repository"
)

func TestServiceCreates{{.TypeName}}(t *testing.T) {
	service := New(repository.NewMemoryRepository())

	{{.Package}}, err := service.Create(context.Background(), domain.Create{{.TypeName}}Command{Name: "Alice"})
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

func TestServiceRejectsMissingName(t *testing.T) {
	service := New(repository.NewMemoryRepository())

	_, err := service.Create(context.Background(), domain.Create{{.TypeName}}Command{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
`

const moduleTemplate = `package module

import (
	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/chenile"

	"{{.ModuleName}}/{{.Package}}/contract"
	"{{.ModuleName}}/{{.Package}}/repository"
	"{{.ModuleName}}/{{.Package}}/service"
)

func New() chenile.Module {
	return chenile.NewModule("{{.Package}}", func(builder *chenile.Builder) error {
		repo := repository.NewMemoryRepository()
		svc := service.New(repo)
		return builder.Service("{{.ServiceID}}").Routes(contract.Routes(svc)...)
	})
}
`

const configTemplate = `service.name: {{.ServiceID}}
server.port: 8080
chenile.profile: local
`

const testTemplate = `package test

import (
	"io"
	"testing"

	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/packager"
	godogtest "github.com/gauravbhardwaj2code/go_chenile/chenile-framework/bdd-utils/godog"

	"{{.ModuleName}}/{{.Package}}/module"
)

func TestCreate{{.TypeName}}(t *testing.T) {
	app, err := packager.NewChenileWebApp(module.New())
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

  Scenario: Reject {{.Package}} without name
    When I POST a REST request to URL "{{.RouteBase}}" with payload
      """
      {
        "name": ""
      }
      """
    Then the http status code is 400
    And success is false
    And the error array size is 2
`

const fixtureTemplate = `{
  "name": "Alice"
}
`

const invalidFixtureTemplate = `{
  "name": ""
}
`
