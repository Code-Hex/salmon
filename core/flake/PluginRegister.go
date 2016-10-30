package flake

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	pipeline "github.com/mattn/go-pipeline"
	"github.com/pkg/errors"
)

var ExecutorTmpl = `
var run = map[string]func(...string) (string, error){
{{range .}}    "{{.Key}}": plugin.Run{{.Value}},
{{end}}
}

var usage = map[string]string{
{{range .}}    "{{.Key}}": plugin.Detail{{.Value}},
{{end}}
}
`

var TokenTmpl = `{{range .}}    {{.Const}}
{{end}}
)

var CommandNames = map[string]Token{
    "help": HELP,
{{range .}}    "{{.Key}}": {{.Const}},
{{end}}
}
`

type Ptemplate struct {
	Key   string
	Value string
	Const string
}

type Ptemplates []Ptemplate

func PluginRegister() error {

	fileInfos, err := ioutil.ReadDir("plugin")
	if err != nil {
		return errors.Wrapf(err, "Could not open plugin directory")
	}

	var ptemplates Ptemplates
	findSourceCode(&ptemplates, fileInfos)
	if err := migrateExecutorGo(ptemplates); err != nil {
		return err
	}

	if err := migrateTokenGo(ptemplates); err != nil {
		return err
	}

	os.Stdout.WriteString("Migrated executor.go, token.go\n")

	return nil
}

func findSourceCode(ptemplates *Ptemplates, fileInfos []os.FileInfo) {
	regex := regexp.MustCompile(`func Run([A-Za-z]+)`)

	for _, fi := range fileInfos {
		f, err := os.Open("plugin/" + fi.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}
		if pt := grepRunInPlugins(regex, f); pt != (Ptemplate{}) {
			*ptemplates = append(*ptemplates, pt)
		}
		f.Close()
	}
}

func grepRunInPlugins(re *regexp.Regexp, f *os.File) Ptemplate {
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		text := scanner.Text()
		if re.MatchString(text) {
			captured := re.FindStringSubmatch(text)
			name := captured[1]

			return Ptemplate{
				Key:   strings.ToLower(name),
				Value: name,
				Const: strings.ToTitle(name),
			}
		}
	}

	return Ptemplate{}
}

func migrateExecutorGo(ptemplates Ptemplates) error {

	f, err := os.OpenFile("executor.go", os.O_RDWR, 0666)
	if err != nil {
		return errors.Wrapf(err, "Could not open executor.go")
	}
	defer f.Close()

	var src string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, `var run = map[string]func(...string)`) {
			break
		}
		src += text + "\n"
	}

	f.Seek(0, io.SeekStart)
	f.Truncate(0)

	tmpl := template.New("ExecutorTmpl for executor.go")
	template.Must(tmpl.Parse(src + ExecutorTmpl))

	tmpl.Execute(f, ptemplates)
	out, err := pipeline.Output([]string{"goimports", "executor.go"}, []string{"gofmt"})
	if err != nil {
		return errors.Wrapf(err, "Failed to static analyze")
	}

	f.Seek(0, io.SeekStart)
	f.Truncate(0)

	f.Write(out)

	return nil
}

func migrateTokenGo(ptemplates Ptemplates) error {

	f, err := os.OpenFile("token.go", os.O_RDWR, 0666)
	if err != nil {
		return errors.Wrapf(err, "Could not open token.go")
	}
	defer f.Close()

	var src string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, `// Command keywords`) {
			src += text + "\n"
			break
		}
		src += text + "\n"
	}

	f.Seek(0, io.SeekStart)
	f.Truncate(0)

	tmpl := template.New("TokenTmpl for token.go")
	template.Must(tmpl.Parse(src + TokenTmpl))

	tmpl.Execute(f, ptemplates)

	out, err := pipeline.Output([]string{"goimports", "token.go"}, []string{"gofmt"})
	if err != nil {
		return errors.Wrapf(err, "Failed to static analyze")
	}

	f.Seek(0, io.SeekStart)
	f.Truncate(0)

	f.Write(out)

	return nil
}
