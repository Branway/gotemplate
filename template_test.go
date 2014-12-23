// Tests for template

package main

import (
	"bytes"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"testing"
)

type TestTemplate struct {
	title   string
	args    string
	pkg     string
	in      string
	outName string
	out     string
}

const basicTest = `package tt

import "fmt"

// template type Set(A)
type A int

type Set struct { a A }
func NewSet(a A) A { return A(0) }
func NewSizedSet(a A) A { return A(1) }
func UtilityFunc1() {}
func utilityFunc() {}
func (a A) f0() {}
func (a *A) F1() {}
var AVar1 int
var aVar2 int
var (
	AVar3 int
	aVar4 int
)
`

var tests = []TestTemplate{
	{
		title:   "Simple test public",
		args:    "MySet(int)",
		pkg:     "main",
		in:      basicTest,
		outName: "gotemplate_MySet.go",
		out: `package main

import "fmt"

// template type Set(A)

type MySet struct{ a int }

func NewMySet(a int) int      { return int(0) }
func NewSizedMySet(a int) int { return int(1) }
func UtilityFunc1MySet()      {}
func utilityFuncMySet()       {}
func (a int) f0()             {}
func (a *int) F1()            {}

var AVar1MySet int
var aVar2MySet int
var (
	AVar3MySet int
	aVar4MySet int
)
`,
	},
	{
		title:   "Simple test private",
		args:    "mySet(float64)",
		pkg:     "main",
		in:      basicTest,
		outName: "gotemplate_mySet.go",
		out: `package main

import "fmt"

// template type Set(A)

type mySet struct{ a float64 }

func newMySet(a float64) float64      { return float64(0) }
func newSizedMySet(a float64) float64 { return float64(1) }
func utilityFunc1MySet()              {}
func utilityFuncMySet()               {}
func (a float64) f0()                 {}
func (a *float64) F1()                {}

var aVar1MySet int
var aVar2MySet int
var (
	aVar3MySet int
	aVar4MySet int
)
`,
	},
	{
		title: "Test function",
		args:  "Min(int8, func(a int8, b int8) bool { return a < b })",
		pkg:   "main",
		in: `package tt

// template type TT(A, Less)
type A int
func Less(a, b A) { return a < b }

func TT(a, b A) A { return Less(a, b) }
func TTone(a A) A { return !Less(a, b) }
`,
		outName: "gotemplate_Min.go",
		out: `package main

// template type TT(A, Less)

func Min(a, b int8) int8 {
	return func(a int8, b int8) bool {
		return a < b
	}(a, b)
}
func Minone(a int8) int8 {
	return !func(a int8, b int8) bool {
		return a < b
	}(a, b)
}
`,
	},
	{
		title: "Simple Test constants",
		args:  "Vector2(float32, 2)",
		pkg:   "main",
		in: `package tt

// template type Vector(A, n)
type A float32
const n = 3

type Vector [n]A

func (v Vector) Add(b Vector) {
	for i := range v {
		v += b[i]
	}
}
`,
		outName: "gotemplate_Vector2.go",
		out: `package main

// template type Vector(A, n)

type Vector2 [2]float32

func (v Vector2) Add(b Vector2) {
	for i := range v {
		v += b[i]
	}
}
`,
	},
	{
		title: "Test constants",
		args:  "Matrix22(float32, 2, 2)",
		pkg:   "main",
		in: `package mat

// template type Matrix(A, n, m)
type A float32

const (
	n, a, b, m = 1, 2, 3, 1
)

type Matrix [n][m]A

func (mat Matrix) Add(x Matrix) {
	for i := range mat {
		for j := range mat[i] {
			mat[i][j] += x[i][j]
		}
	}
}
`,
		outName: "gotemplate_Matrix22.go",
		out: `package main

// template type Matrix(A, n, m)

const (
	aMatrix22, bMatrix22 = 2, 3
)

type Matrix22 [2][2]float32

func (mat Matrix22) Add(x Matrix22) {
	for i := range mat {
		for j := range mat[i] {
			mat[i][j] += x[i][j]
		}
	}
}
`,
	},
	{
		title: "Test vars",
		args:  "ProgXX(xx1, xx2, xx3, xx4, xx5, xx6)",
		pkg:   "main",
		in: `package prog

// template type Prog(a, b, c, d, e, f)
type A float32

var (
	a, z = 1, 2
	b, n, m, c = 3, 4, 5, 6
	d = 7
)

var (
	o = 8
	e = 8
)

var (
	e = 8
)

var (
	oo = 9
	e = 10
)

var (
	p, f, q = 11, 12, 13
)

func Prog() int {return a+b+c+d+e+f}
`,
		outName: "gotemplate_ProgXX.go",
		out: `package main

// template type Prog(a, b, c, d, e, f)
type AProgXX float32

var (
	zProgXX          = 2
	nProgXX, mProgXX = 4, 5
)

var (
	oProgXX = 8
)

var (
	ooProgXX = 9
)

var (
	pProgXX, qProgXX = 11, 13
)

func ProgXX() int { return xx1 + xx2 + xx3 + xx4 + xx5 + xx6 }
`,
	},
	{
		title: "Test complex type decls",
		args:  "tmpl(int, string, map[string]map[string]chan int, float32, rune, chan []string)",
		pkg:   "main",
		in: `package tt

// template type TMPL(A, B, C, D, E, F)
type A int

type TMPL struct {
	a A
	b B
	c C
	d D
	e E
	f F
}

type ImportantType bool

type (
	ImportantType1 int
	B struct {
		v map[int][][][]rune
	}
	importantType2 map[int]int
	C chan struct {
		x []string
	}
)

type (
	D rune
	importantType3 struct{}
	E string
	F map[string]int
)
`,
		outName: "gotemplate_tmpl.go",
		out: `package main

// template type TMPL(A, B, C, D, E, F)

type tmpl struct {
	a int
	b string
	c map[string]map[string]chan int
	d float32
	e rune
	f chan []string
}

type importantTypeTmpl bool

type (
	importantType1Tmpl int

	importantType2Tmpl map[int]int
)

type (
	importantType3Tmpl struct{}
)
`,
	},
}

func testTemplate(t *testing.T, test *TestTemplate) {
	// Disable logging
	log.SetOutput(ioutil.Discard)

	// Make temporary directory
	dir, err := ioutil.TempDir("", "gotemplate_test")
	if err != nil {
		t.Fatalf("Failed to make temp dir: %v", err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Make subdirectories
	src := path.Join(dir, "src")
	err = os.Mkdir(src, 0700)
	if err != nil {
		t.Fatalf("Failed to make dir %q: %v", src, err)
	}
	input := path.Join(src, "input")
	err = os.Mkdir(input, 0700)
	if err != nil {
		t.Fatalf("Failed to make dir %q: %v", input, err)
	}
	output := path.Join(src, "output")
	err = os.Mkdir(output, 0700)
	if err != nil {
		t.Fatalf("Failed to make dir %q: %v", output, err)
	}

	// Change directory to output directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to read cwd: %v", err)
	}
	err = os.Chdir(output)
	if err != nil {
		t.Fatalf("Failed to cd %q dir: %v", output, err)
	}
	defer func() {
		err := os.Chdir(cwd)
		if err != nil {
			t.Logf("Failed to change back to cwd: %v", err)
		}
	}()

	// Set GOPATH to directory
	build.Default.GOPATH = dir

	// Write template input
	tmpl := path.Join(input, "main.go")
	err = ioutil.WriteFile(tmpl, []byte(test.in), 0600)
	if err != nil {
		t.Fatalf("Failed to write %q: %v", tmpl, err)
	}

	// Write main.go for output
	main := path.Join(output, "main.go")
	err = ioutil.WriteFile(main, []byte("package main"), 0600)
	if err != nil {
		t.Fatalf("Failed to write %q: %v", main, err)
	}

	// Instantiate template
	template := newTemplate(output, "input", test.args)
	template.instantiate()

	// Check output
	expectedFile := path.Join(output, test.outName)
	actualBytes, err := ioutil.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("Failed to read %q: %v", expectedFile, err)
	}
	actual := string(actualBytes)
	if actual != test.out {
		t.Errorf(`Output is wrong
Got
-------------
%s
-------------
Expected
-------------
%s
-------------
`, actual, test.out)
		actualFile := expectedFile + ".actual"
		err = ioutil.WriteFile(actualFile, []byte(test.out), 0600)
		if err != nil {
			t.Fatalf("Failed to write %q: %v", actualFile, err)
		}
		cmd := exec.Command("diff", "-u", actualFile, expectedFile)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		_ = cmd.Run()
		t.Errorf("Diff\n----\n%s", out.String())
	}

}

func TestSub(t *testing.T) {
	fatalf = func(format string, args ...interface{}) {
		t.Fatalf(format, args...)
	}
	for i := range tests {
		t.Logf("Test[%d] %q", i, tests[i].title)
		testTemplate(t, &tests[i])
	}
}