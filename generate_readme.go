package main

import (
	"fmt"
	"github.com/elliotchance/pie/functions"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Function struct {
	Name string
	For  int
	Doc  string
	BigO string
}

func (f Function) BriefDoc() (brief string) {
	for _, line := range strings.Split(f.Doc, "\n") {
		if line == "" {
			break
		}

		brief += line + " "
	}

	return
}

func main() {
	var funcs []Function

	for _, function := range functions.Functions {
		file, err := ioutil.ReadFile("functions/" + function.File)
		if err != nil {
			panic(err)
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", file, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		funcs = append(funcs, Function{
			Name: function.Name,
			For:  function.For,
			Doc:  getDoc(f.Decls),
			BigO: function.BigO,
		})
	}

	longestFunctionName := 0
	for _, function := range funcs {
		if newLen := len(function.Name); newLen > longestFunctionName {
			longestFunctionName = newLen
		}
	}

	newDocs := fmt.Sprintf("| Function%s | String | Number | Struct | Maps | Big-O    | Description |\n", strings.Repeat(" ", longestFunctionName-6))
	newDocs += fmt.Sprintf("| %s | :----: | :----: | :----: | :--: | :------: | ----------- |\n", strings.Repeat("-", longestFunctionName+2))

	for _, function := range funcs {
		newDocs += fmt.Sprintf("| `%s`%s | %s      | %s      | %s      | %s    | %s%s | %s |\n",
			function.Name,
			strings.Repeat(" ", longestFunctionName-len(function.Name)),
			tick(function.For&functions.ForStrings),
			tick(function.For&functions.ForNumbers),
			tick(function.For&functions.ForStructs),
			tick(function.For&functions.ForMaps),
			function.BigO,
			strings.Repeat(" ", 8-utf8.RuneCountInString(function.BigO)),
			function.BriefDoc())
	}

	newDocs += "\n"

	for _, function := range funcs {
		newDocs += fmt.Sprintf("## %s\n\n%s\n",
			function.Name, function.Doc)
	}

	readme, err := ioutil.ReadFile("README.md")
	if err != nil {
		panic(err)
	}

	newReadme := regexp.MustCompile(`(?s)\| Function.*# FAQ`).
		ReplaceAllString(string(readme), newDocs + "# FAQ")
	//fmt.Printf("%s\n", newReadme)

	err = ioutil.WriteFile("README.md", []byte(newReadme), 0644)
	if err != nil {
		panic(err)
	}
}

func tick(x int) string {
	if x != 0 {
		return "✓"
	}

	return " "
}

func getDoc(decls []ast.Decl) (doc string) {
	for _, decl := range decls {
		if f, ok := decl.(*ast.FuncDecl); ok && f.Doc != nil {
			for _, comment := range f.Doc.List {
				if len(comment.Text) < 3 {
					doc += "\n"
				} else {
					doc += comment.Text[3:] + "\n"
				}
			}
		}
	}

	//doc = strings.TrimLeft(doc, " ")

	return
}
