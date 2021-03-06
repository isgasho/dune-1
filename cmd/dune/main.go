package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dunelang/dune"
	"github.com/dunelang/dune/ast"
	"github.com/dunelang/dune/binary"
	"github.com/dunelang/dune/filesystem"
	"github.com/dunelang/dune/parser"

	_ "github.com/dunelang/dune/lib"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func printVersion() {
	fmt.Printf("%s\n\n", dune.VERSION)
}

func main() {
	v := flag.Bool("v", false, "version")
	c := flag.Bool("c", false, "compile")
	s := flag.Bool("s", false, "strip")
	e := flag.Bool("e", false, "eval")
	o := flag.String("o", "", "output file")
	d := flag.Bool("d", false, "decompile")
	a := flag.Bool("a", false, "show AST")
	r := flag.Bool("r", false, "list resources")
	n := flag.Bool("n", false, "no optimizations")
	i := flag.Bool("i", false, "generate native.d.ts and tsconfig.json")
	dts := flag.Bool("dts", false, "generate native.d.ts")
	flag.Parse()

	args := flag.Args()
	aLen := len(args)

	if *v {
		printVersion()
		return
	}

	if *n {
		parser.Optimizations = false
	}

	if *e {
		if aLen != 1 {
			fatal("only one parameter allowed")
		}
		if err := eval(args[0]); err != nil {
			fatal(err)
		}
		return
	}

	if *d {
		p, err := loadProgram(args[0], *s)
		if err != nil {
			fatal(err)
		}
		dune.Print(p)
		return
	}

	if *a {
		at, err := parser.Parse(filesystem.OS, args[0])
		if err != nil {
			fatal(err)
		}
		ast.Print(at)
		return
	}

	if *c {
		p, err := loadProgram(args[0], *s)
		if err != nil {
			fatal(err)
		}

		out := *o
		if out == "" {
			n := filepath.Base(args[0])
			out = strings.TrimSuffix(n, filepath.Ext(n)) + ".bin"
		}
		if err := build(p, out); err != nil {
			fatal(err)
		}
		return
	}

	if *r {
		p, err := loadProgram(args[0], *s)
		if err != nil {
			fatal(err)
		}
		for k, v := range p.Resources {
			fmt.Println(k, len(v))
		}
		return
	}

	if *i {
		var path string
		if aLen == 1 {
			path = args[0]
		}
		generateTsconfig(path)
		generateDts(path)
		return
	}

	if *dts {
		var path string
		if aLen == 1 {
			path = args[0]
		}
		generateDts(path)
		return
	}

	if aLen > 0 {
		if err := exec(args[0], args[1:]); err != nil {
			fatal(err)
		}
		return
	}

	if err := startREPL(); err != nil {
		fatal(err)
	}
}

func build(p *dune.Program, out string) error {
	f, err := os.OpenFile(out, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	if err := binary.Write(f, p); err != nil {
		return err
	}

	return nil
}

func eval(code string) error {
	code = fmt.Sprintf("console.log(%s)", code)
	p, err := dune.CompileStr(code)
	if err != nil {
		return err
	}

	p.AddPermission("trusted")

	vm := dune.NewVM(p)
	vm.FileSystem = filesystem.OS
	_, err = vm.Run()
	return err
}

func exec(programPath string, args []string) error {
	p, err := loadProgram(programPath, false)
	if err != nil {
		return err
	}

	p.AddPermission("trusted")

	vm := dune.NewVM(p)
	vm.FileSystem = filesystem.OS

	ln := len(args)
	values := make([]dune.Value, ln)
	for i := 0; i < ln; i++ {
		values[i] = dune.NewValue(args[i])
	}

	_, err = vm.Run(values...)
	return err
}

func loadProgram(path string, strip bool) (*dune.Program, error) {
	path, ok := find(path)
	if !ok {
		return nil, os.ErrNotExist
	}
	return doLoadProgram(path, strip)
}

func find(path string) (string, bool) {
	if filesystem.Exists(filesystem.OS, path) {
		return path, true
	}

	if filepath.Base(path) == path && filepath.Ext(path) == "" {
		dirs := os.Getenv("DUNE_DIRS")
		if dirs != "" {
			for _, dir := range strings.Split(dirs, ";") {
				t, ok := testPath(dir, path)
				if ok {
					return t, true
				}
			}
		}
	}

	return "", false
}

func testPath(dir, path string) (string, bool) {
	t := filepath.Join(dir, path+".ts")
	if filesystem.Exists(filesystem.OS, t) {
		return t, true
	}

	t = filepath.Join(dir, path+".bin")
	if filesystem.Exists(filesystem.OS, t) {
		return t, true
	}

	return "", false
}

func doLoadProgram(path string, strip bool) (*dune.Program, error) {
	ext := filepath.Ext(path)

	switch ext {
	case ".bin":
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("error opening %s: %w", path, err)
		}
		defer f.Close()
		p, err := binary.Read(f)
		if err != nil {
			return nil, fmt.Errorf("error loading %s: %w", path, err)
		}
		if strip {
			p.Strip()
		}
		return p, nil

	default:
		p, err := dune.Compile(filesystem.OS, path)
		if err != nil {
			return nil, fmt.Errorf("error loading %s: %w", path, err)
		}
		if strip {
			p.Strip()
		}
		return p, err
	}
}

func generateDts(path string) {
	if path == "" {
		path = "."
	}

	filesystem.WritePath(filesystem.OS, filepath.Join(path, "native.d.ts"), []byte(dune.TypeDefs()))
}

func generateTsconfig(path string) {
	if path == "" {
		path = "."
	}

	writeIfNotExists(filesystem.OS, filepath.Join(path, "tsconfig.json"), []byte(`{
	"compilerOptions": {
		"noLib": true,
		"noEmit": true,
		"noImplicitAny": true,
		"baseUrl": ".",
		"paths": {
			"*": [
				"*",
				"vendor/*"
			]
		}
	}
}
`))
}

func writeIfNotExists(fs filesystem.FS, name string, data []byte) {
	if _, err := fs.Stat(name); err == nil {
		return
	}
	filesystem.WritePath(fs, name, data)
}

func fatal(values ...interface{}) {
	fmt.Println(values...)
	os.Exit(1)
}
