package terraform

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	// external packages
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func addSecretRoutes(g *gin.RouterGroup) {
	const path = "/secrets/:name"
	g.GET(path, HandleGetSecret)
	g.POST(path, HandleUpdateSecret)
	g.Handle(httpMethodLock, path, HandleSecretLock)
	g.Handle(httpMethodUnlock, path, HandleSecretUnlock)
}

const tfmain = `
terraform {
  backend local {
    path = "terraform.tfstate"
  }
}

output a {
  value = 1
}

output b {
  value = "b"
  sensitive = true
}

output c {
  value = { "a": 1, "b": "2", "c": [ 1, 2, 3 ], "d": null, "e": { "a": "b" } }
}

output d {
  value = [ 1,2,3,4,5 ]
}

output e {
  value = 1.2
}

output f {
  value = null
}`

type SecretDTO struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	IsString bool   `json:"is_string,omitempty"`
}

func (s *SecretDTO) ToValue() string {
	return s.Value
}

func HandleGetSecret(ctx *gin.Context) {

	const t = `
terraform {
  backend local {
    path = "terraform.tfstate"
  }
}

{{range . -}}
output {{.Name}} {
  value = {{value .Value .IsString}}
}
{{end}}
`

	dto := []SecretDTO{
		{
			Name:  "a",
			Value: "1",
		},
		{
			Name:  "b",
			Value: "b",
		},
		{
			Name:  "c",
			Value: `{ "a": 1, "b": "2", "c": [ 1, 2, 3 ], "d": null, "e": { "a": "b" } }`,
		},
		{
			Name:  "d",
			Value: "[ 1,2,3,4,5 ]",
		},
		{
			Name:  "e",
			Value: "12.34",
		},
		{
			Name:     "f",
			Value:    "12.34",
			IsString: true,
		},
		{
			Name:  "g",
			Value: "null",
		},
	}

	fmap := template.FuncMap{
		"value": func(v string, isString bool) string {
			if strings.HasPrefix(v, "[") || strings.HasPrefix(v, "{") {
				if isString {

				}
				return v
			}
			if regexp.MustCompile(`\d(|\.\d)`).MatchString(v) && !isString {
				return v
			}
			return fmt.Sprintf("\"%s\"", v)
		},
	}

	tmpl, err := template.New("tfmain").Funcs(fmap).Parse(t)
	if err != nil {
		log.Fatal(err)
	}
	if err := tmpl.Execute(os.Stdout, dto); err != nil {
		log.Fatal(err)
	}

	name := ctx.Param("name")

	dir, err := os.MkdirTemp("/tmp", fmt.Sprintf("%s--*", name))
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file := filepath.Join(dir, "main.tf")
	if err := os.WriteFile(file, []byte(tfmain), 0600); err != nil {
		log.Fatal(err)
	}

	init := exec.Command("terraform", "init")
	init.Dir = dir
	init.Stdout = os.Stdout
	if err := init.Run(); err != nil {
		fmt.Println("could not run command: ", err)
	}

	apply := exec.Command("terraform", "apply", "-auto-approve")
	apply.Dir = dir
	apply.Stdout = os.Stdout
	if err := apply.Run(); err != nil {
		fmt.Println("could not run command: ", err)
	}

	tfstate := filepath.Join(dir, "terraform.tfstate")
	if state, err := os.ReadFile(tfstate); err != nil {
		log.Error("Could not read terraform state", err)
	} else {
		log.Info(string(state))
	}

	v := `{
  "version": 4,
  "lineage": "145ec231-90be-4a48-e96d-00b782f31f82",
  "outputs": {
    "a": {
      "value": "a",
      "type": "string"
    },
    "b": {
      "value": "b",
      "type": "string",
      "sensitive": true
    }
  }
}`
	s := Secret{}
	json.Unmarshal([]byte(v), &s)

	ctx.JSON(200, s)
}

func HandleUpdateSecret(ctx *gin.Context) {
	ctx.Status(http.StatusForbidden)
}

func HandleSecretLock(ctx *gin.Context) {
	ctx.Status(http.StatusForbidden)
}

func HandleSecretUnlock(ctx *gin.Context) {
	ctx.Status(http.StatusForbidden)
}
