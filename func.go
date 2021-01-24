package src

// A Func is a function or method, depending on the context it appears. E.g. within an interface, the receiver name
// and the body is not evaluated.
type Func struct {
	doc           string
	name          string
	static        bool
	visibility    Visibility
	recName       string
	isPtrReceiver bool
	params        []*Param
	results       []*Param
	body          *Block
	variadic      bool
	annotations   []*Annotation
}

// NewFunc allocates a new parameterless function with a calling name.
func NewFunc(name string) *Func {
	return &Func{name: name}
}

// Name returns the declared identifier which must be unique in its scope.
func (s *Func) Name() string {
	return s.name
}

// SetName updates the functions name which must be unique in its scope (e.g. type or package).
func (s *Func) SetName(name string) *Func {
	s.name = name
	return s
}

// SetVisibility sets the visibility. The default is Public.
func (s *Func) SetVisibility(v Visibility) *Func {
	s.visibility = v
	return s
}

// Visibility returns the current visibility. The default is Public.
func (s *Func) Visibility() Visibility {
	return s.visibility
}

// SetDoc sets functions documentation.
func (s *Func) SetDoc(doc string) *Func {
	s.doc = doc
	return s
}

// Doc returns the functions documentation.
func (s *Func) Doc() string {
	return s.doc
}

// RecName returns the receivers name. The java renderer will ignore this and substitute it implicitly with this.
func (s *Func) RecName() string {
	return s.recName
}

// SetRecName updates the receivers name.
func (s *Func) SetRecName(recName string) *Func {
	s.recName = recName
	return s
}

// PtrReceiver is currently only applicable for Go. In java this may be introduced with Valhalla.
func (s *Func) PtrReceiver() bool {
	return s.isPtrReceiver
}

// SetPtrReceiver updates the receiver to be a pointer type.
func (s *Func) SetPtrReceiver(isPtrReceiver bool) *Func {
	s.isPtrReceiver = isPtrReceiver
	return s
}

// Params returns the backing array of the input parameters.
func (s *Func) Params() []*Param {
	return s.params
}

// SetParams updates the backing array of input parameters.
func (s *Func) SetParams(params ...*Param) *Func {
	s.params = params
	return s
}

// SetParams adds to the backing array of input parameters.
func (s *Func) AddParams(params ...*Param) *Func {
	s.params = append(s.params, params...)
	return s
}

// Results returns the backing array of the out parameters. In languages which only support none (void) or
// one result, all following parameters are treated as Exceptions.
func (s *Func) Results() []*Param {
	return s.results
}

// SetResults updates the backing array of the out parameters. In languages which only support none (void) or
// one result, all following parameters are treated as Exceptions.
func (s *Func) SetResults(results ...*Param) *Func {
	s.results = results
	return s
}

// AddResults appends to the backing array of the out parameters. In languages which only support none (void) or
// one result, all following parameters are treated as Exceptions.
func (s *Func) AddResults(results ...*Param) *Func {
	s.results = append(s.results, results...)
	return s
}

// Body returns the implementation.
func (s *Func) Body() *Block {
	return s.body
}

// SetBody updates the implementation.
func (s *Func) SetBody(body *Block) *Func {
	s.body = body
	return s
}

// Variadic returns true, if the last function parameter should be treated as a variable argument. Language which
// do not support that, fallback to a slice.
func (s *Func) Variadic() bool {
	return s.variadic
}

// SetVariadic updates the variadic state of the last parameter.
func (s *Func) SetVariadic(variadic bool) *Func {
	s.variadic = variadic
	return s
}

// Annotations returns the backing slice of all annotations.
func (s *Func) Annotations() []*Annotation {
	return s.annotations
}

// AddAnnotations appends the given annotations. Note that not all render targets support method annotations, e.g.
// like Go.
func (s *Func) AddAnnotations(a ...*Annotation) *Func {
	s.annotations = append(s.annotations, a...)
	return s
}

// SetStatic updates the static flag of the method. This declares a method to be not part of the according
// instance and it will not be able to modify its receiver instance, so the PtrReceiver flag is ignored.
// In Java, this will cause the renderer to emit a class scoped method.
func (s *Func) SetStatic(b bool) *Func {
	s.static = b
	return s
}

// Static returns the static flag.
func (s *Func) Static() bool {
	return s.static
}