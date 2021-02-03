package ast2

// Arch defines the architecture to generate code for.
type Arch string

const (
	ArchAMD64 Arch = "amd64"
	ArchARM64 Arch = "arm64"
	ArchWASM  Arch = "wasm"
)

// OS defines the operating system to generate code for.
type OS string

const (
	OSLinux  OS = "linux"
	OSWin    OS = "windows"
	OSDarwin OS = "darwin"
	OSIOS    OS = "iOS"
)

// Lang defines the target language to generate code for.
type Lang string

const (
	LangJava   Lang = "java"
	LangGo     Lang = "go"
	LangRust   Lang = "rust"
	LangSwift  Lang = "swift"
	LangKotlin Lang = "kotlin"
	LangC      Lang = "c"
	LangCPP    Lang = "c++"
)

// LangVersion specifies an arbitrary version string for a specific language. There is no guarantee of a semantic
// version logic here.
type LangVersion string

const (
	LangVersionJava8 LangVersion = "1.8"
	LangVersionGo16  LangVersion = "1.16"
	LangVersionSwift LangVersion = "5.1"
)

// Framework specifies an arbitrary framework string. This is mostly unique In conjunction with a language.
type Framework string

const (
	FrameworkSDK Framework = ""
)

type Target struct {
	Arch           Arch        // probably empty and of limited use.
	Os             OS          // probably empty and of limited use.
	Lang           Lang        // target generator language.
	MinLangVersion LangVersion // the (inclusive) supported minimum version of the generated code.
	MaxLangVersion LangVersion // the (inclusive) supported maximum version of the generated code.
	Framework      Framework   // the framework to use. Empty means only use the default standard library things.
}

func (t Target) Equals(o Target) bool {
	return t.Lang == o.Lang && t.Os == o.Os && t.Arch == o.Arch && t.MinLangVersion == o.MinLangVersion && t.MaxLangVersion == o.MaxLangVersion && t.Framework == o.Framework
}

// A Mod is the root of a project and describes a module with packages.
//  * Java: denotes a gradle module (build.gradle).
//  * Go: describes a Go module (go.mod).
type Mod struct {
	Target Target
	Pkgs   []*Pkg
	Obj
}

// NewModule allocates a new Module.
func NewModule() *Mod {
	return &Mod{}
}

// Packages appends the given packages and updates the Parent accordingly.
func (n *Mod) Packages(packages ...*Pkg) *Mod {
	n.Pkgs = append(n.Pkgs, packages...)
	for _, pkg := range packages {
		assertNotAttached(pkg)
		pkg.Obj.ObjParent = n
	}

	return n
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *Mod) Children() []Node {
	tmp := make([]Node, 0, len(n.Pkgs)+1)
	tmp = append(tmp, n.Obj.ObjComment)

	for _, pkg := range n.Pkgs {
		tmp = append(tmp, pkg)
	}

	return tmp
}

// Doc sets the nodes comment.
func (n *Mod) Doc(text string) *Mod {
	n.Obj.ObjComment = NewComment(text)
	n.Obj.ObjComment.ObjParent = n
	return n
}
