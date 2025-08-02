// main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jeffrey-elliott/confitgen"
)

func main() {
	hrule := func() {
		fmt.Println(strings.Repeat("-", 60))
	}

	packageName := flag.String("schema", "", "the name of the package")
	longHelp := flag.Bool("more", false, "show extended help")

	flag.Usage = func() {
		hrule()
		fmt.Println("Usage:")
		fmt.Println("  confitgen -schema <package>")
		fmt.Println()
		fmt.Println("Run with -more for more documentation.")
		hrule()
	}

	flag.Parse()

	if *longHelp {
		hrule()
		fmt.Println("Example:")
		fmt.Println("  confitgen -schema mypackage")
		fmt.Println()
		fmt.Println("       required schema: mypackage.confit.schema.json")
		fmt.Println("  required schema root: MyPackageConfitSchema")
		fmt.Println("        file generated: mypackage.go")
		fmt.Println()
		fmt.Println("In your own app, drop that file here:")
		fmt.Println("  internal/confit/mypackage.go")
		fmt.Println()
		fmt.Println("For more, visit:")
		fmt.Println("  github.com/jeffrey-elliott/confitgen")
		hrule()
		os.Exit(0)
	}

	schemaFilename := fmt.Sprintf("%s.confit.schema.json", *packageName)

	schemaData, err := os.ReadFile(schemaFilename)
	if err != nil {
		log.Fatal(err)
	}

	outFile := fmt.Sprintf("%s.go", *packageName)
	f, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := confitgen.Generate(f, schemaData); err != nil {
		log.Fatal(err)
	}
}
