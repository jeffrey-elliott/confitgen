// gen_test.go
package confitgen_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/jeffrey-elliott/confitgen"
)

const verbose = false

const FILE_GUIDEBOOK_SCHEMA = "guidebook.confit.schema.json"
const FILE_GUIDEBOOK_VALUES = "guidebook.confit.values.json"
const FILE_STARSHIP_SCHEMA = "starship.confit.schema.json"
const FILE_STARSHIP_VALUES = "starship.confit.values.json"
const FILE_MAIN_DOT_GO_A = "main.a.go"
const FILE_MAIN_DOT_GO_B = "main.b.go"

func readBytes(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	return data
}

func readBytesTestdata(t *testing.T, name string) []byte {
	t.Helper()
	return readBytes(t, filepath.Join("testdata", name))
}

func readStringGolden(t *testing.T, name string) string {
	t.Helper()
	fb := readBytes(t, filepath.Join("testdata", "golden", name))
	return normalizeLines(string(fb))
}

func normalizeLines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

func writeTempFile(t *testing.T, opts TestOptions, name string, data []byte) {
	t.Helper()
	fp := filepath.Join(opts.Dir, name)
	if err := os.WriteFile(fp, data, 0644); err != nil {
		t.Fatalf("failed to write temp%s: %v", fp, err)
	}
}

func writeTempInternalPackageFile(t *testing.T, opts TestOptions, packageName string, data []byte) {
	t.Helper()
	p := path.Join(opts.DirInternal, packageName)
	fp := filepath.Join(p, fmt.Sprint(packageName, ".go"))
	if err := os.MkdirAll(p, 0755); err != nil {
		t.Fatalf("failed to write folder: temp/interna/confit/packagename%s: %v", p, err)
	}
	if err := os.WriteFile(fp, data, 0644); err != nil {
		t.Fatalf("failed to write file: temp/internal/confit/packagename%s: %v", fp, err)
	}
}

func writeGoMod(t *testing.T, opts TestOptions, name string) {
	t.Helper()
	content := fmt.Sprint("module ", name, "\n\ngo 1.21\n")
	writeTempFile(t, opts, "go.mod", []byte(content))
}

// rune line count
func rlc(s string) string {
	result := struct {
		runes int
		lines int
	}{
		runes: utf8.RuneCountInString(s),
		lines: strings.Count(s, "\n") + 1,
	}

	if len(s) == 0 {
		result.lines = 0
	}

	return fmt.Sprintf("%+v", result)
}

// have want rune line count
func wgrlc(want string, got string) string {
	return fmt.Sprintf("want: %s\n got: %s", rlc(want), rlc(got))
}

func failSummary(t *testing.T, want string, got string) string {
	return fmt.Sprintf(
		"%s failed:\nwant:\n%s\ngot:\n%s\n%s\n",
		t.Name(), want, got, wgrlc(want, got),
	)
}

type TestOptions struct {
	Dir         string
	DirInternal string
}

func writeTempFolders(t *testing.T) (string, string) {
	t.Helper()
	temp := t.TempDir()
	internalPath := filepath.Join(temp, "internal", "confit")
	if err := os.MkdirAll(internalPath, 0755); err != nil {
		t.Fatalf("failed to create internal/confit: %v", err)
	}

	return temp, internalPath
}

// opts := setupTest(t)
func setupTest(t *testing.T) TestOptions {
	tmp, tmpInternal := writeTempFolders(t)

	return TestOptions{
		Dir:         tmp,
		DirInternal: tmpInternal,
	}
}

func TestAppendPackageImport(t *testing.T) {
	want := string(readStringGolden(t, "append-imports.frag"))

	b := confitgen.NewFormattingStringBuilder()
	confitgen.AppendImports("packagename", b)

	got := b.String()

	if got != string(want) {
		t.Fatal(failSummary(t, want, got))
	}
}

func TestAppendFunctions(t *testing.T) {
	want := string(readStringGolden(t, "append-functions.frag"))

	b := confitgen.NewFormattingStringBuilder()
	confitgen.AppendFunctions("TypeName", b)

	got := b.String()

	if got != string(want) {
		t.Fatal(failSummary(t, want, got))
	}
}

func TestGenerate_Structs(t *testing.T) {
	schema := readBytesTestdata(t, FILE_GUIDEBOOK_SCHEMA)

	var buf bytes.Buffer
	err := confitgen.Generate(&buf, []byte(schema))
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	out := buf.String()
	if verbose {
		t.Logf("Output:\n%s", out)
	}

	checks := []string{
		"type Guidebook struct {",
		"type Narration struct {",
		"Britishness  float64 `json:\"Britishness\"`",
		"func New() Guidebook { return Guidebook{} }",
		"func Load(path string) (Guidebook, error)",
	}
	for _, substr := range checks {
		if !strings.Contains(out, substr) {
			t.Errorf("expected output to contain %q", substr)
		}
	}
}

func TestGenerate_Compiles(t *testing.T) {
	opts := setupTest(t)
	schema := readBytesTestdata(t, "guidebook.confit.schema.json")

	var buf bytes.Buffer
	err := confitgen.Generate(&buf, []byte(schema))
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	writeTempFile(t, opts, "confit.go", buf.Bytes())
	writeGoMod(t, opts, "example")

	cmd := exec.Command("go", "build", ".")
	cmd.Dir = opts.Dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generated code failed to compile: %v\n%s", err, out)
	}
}

func TestGenerate_IntegratedValues_A(t *testing.T) {
	opts := setupTest(t)

	setupIntegratedValues(t, opts, FILE_GUIDEBOOK_SCHEMA, FILE_GUIDEBOOK_VALUES, "guidebook")

	main := readBytesTestdata(t, FILE_MAIN_DOT_GO_A)

	writeTempFile(t, opts, "main.go", main)
	writeGoMod(t, opts, "hhgttg")

	cmd := exec.Command("go", "run", ".")
	cmd.Dir = opts.Dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("runner failed: %v\n%s", err, out)
	}

	if !strings.Contains(string(out), "panic-free-blue") {
		t.Errorf("expected theme value, got: \n%s", out)
	}
}

func TestGenerate_IntegratedValues_B(t *testing.T) {
	opts := setupTest(t)

	setupIntegratedValues(t, opts, FILE_GUIDEBOOK_SCHEMA, FILE_GUIDEBOOK_VALUES, "guidebook")
	setupIntegratedValues(t, opts, FILE_STARSHIP_SCHEMA, FILE_STARSHIP_VALUES, "starship")

	main := readBytesTestdata(t, FILE_MAIN_DOT_GO_B)

	writeTempFile(t, opts, "main.go", main)
	writeGoMod(t, opts, "hhgttg")

	cmd := exec.Command("go", "run", ".")
	cmd.Dir = opts.Dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("runner failed: %v\n%s", err, out)
	}

	if !strings.Contains(string(out), "panic-free-blue") {
		t.Errorf("expected theme value, got: \n%s", out)
	}
	if !strings.Contains(string(out), "Ford Prefect") {
		t.Errorf("expected crew value, got: \n%s", out)
	}
}

func setupIntegratedValues(t *testing.T, opts TestOptions, schemaName string, valuesName string, packageName string) {
	t.Helper()

	schema := readBytesTestdata(t, schemaName)
	data := readBytesTestdata(t, valuesName)

	writeTempFile(t, opts, valuesName, data)

	var buf bytes.Buffer
	if err := confitgen.Generate(&buf, schema); err != nil {
		t.Fatal(err)
	}

	if verbose {
		fmt.Printf("%v", buf.String())
	}

	writeTempInternalPackageFile(t, opts, packageName, buf.Bytes())
}
