package ast2

// A Func is a function or method, depending on the context it appears. E.g. within an interface, the receiver name
// and the body is not evaluated.
type Func struct {
	FunName         string
	FunStatic       bool
	FunVisibility   Visibility
	FunReceiverName string
	FunPtrReceiver  bool
	FunParams       []*Param
	FunResults      []*Param
	FunBody         *Block
	FunVariadic     bool
	FunAnnotations  []*Annotation
	Obj
}

// NewFunc allocates a new parameterless function with a calling name.
func NewFunc(name string) *Func {
	return &Func{FunName: name}
}

// Name returns the declared identifier which must be unique in its scope.
func (s *Func) Name() string {
	return s.FunName
}

// SetName updates the functions name which must be unique in its scope (e.g. type or package).
func (s *Func) SetName(name string) *Func {
	s.FunName = name
	return s
}

// SetVisibility sets the visibility. The default is Public.
func (s *Func) SetVisibility(v Visibility) *Func {
	s.FunVisibility = v
	return s
}

// Visibility returns the current visibility. The default is Public.
func (s *Func) Visibility() Visibility {
	return s.FunVisibility
}

// RecName returns the receivers name. The java renderer will ignore this and substitute it implicitly with this.
func (s *Func) RecName() string {
	return s.FunReceiverName
}

// SetRecName updates the receivers name.
func (s *Func) SetRecName(recName string) *Func {
	s.FunReceiverName = recName
	return s
}

// PtrReceiver is currently only applicable for Go. In java this may be introduced with Valhalla.
func (s *Func) PtrReceiver() bool {
	return s.FunPtrReceiver
}

// SetPtrReceiver updates the receiver to be a pointer type.
func (s *Func) SetPtrReceiver(isPtrReceiver bool) *Func {
	s.FunPtrReceiver = isPtrReceiver
	return s
}

// Params returns the backing array of the input parameters.
func (s *Func) Params() []*Param {
	return s.FunParams
}

// SetParams updates the backing array of input parameters.
func (s *Func) SetParams(params ...*Param) *Func {
	s.FunParams = params
	return s
}

// SetParams adds to the backing array of input parameters.
func (s *Func) AddParams(params ...*Param) *Func {
	s.FunParams = append(s.FunParams, params...)
	return s
}

// Results returns the backing array of the out parameters. In languages which only support none (void) or
// one result, all following parameters are treated as Exceptions.
func (s *Func) Results() []*Param {
	return s.FunResults
}

// SetResults updates the backing array of the out parameters. In languages which only support none (void) or
// one result, all following parameters are treated as Exceptions.
func (s *Func) SetResults(results ...*Param) *Func {
	s.FunResults = nil

	return s.AddResults(results...)
}

// AddResults appends to the backing array of the out parameters. In languages which only support none (void) or
// one result, all following parameters are treated as Exceptions.
func (s *Func) AddResults(results ...*Param) *Func {
	for _, result := range results {
		assertNotAttached(result)
		assertSettableParent(result).SetParent(result)
	}

	return s
}

// Body returns the implementation.
func (s *Func) Body() *Block {
	return s.FunBody
}

// SetBody updates the implementation.
func (s *Func) SetBody(body *Block) *Func {
	assertNotAttached(body)
	assertSettableParent(body).SetParent(s)
	s.FunBody = body

	return s
}

// Variadic returns true, if the last function parameter should be treated as a variable argument. Language which
// do not support that, fallback to a slice.
func (s *Func) Variadic() bool {
	return s.FunVariadic
}

// SetVariadic updates the variadic state of the last parameter.
func (s *Func) SetVariadic(variadic bool) *Func {
	s.FunVariadic = variadic

	return s
}

// Annotations returns the backing slice of all annotations.
func (s *Func) Annotations() []*Annotation {
	return s.FunAnnotations
}

// AddAnnotations appends the given annotations. Note that not all render targets support method annotations, e.g.
// like Go.
func (s *Func) AddAnnotations(a ...*Annotation) *Func {
	for _, annotation := range a {
		assertNotAttached(annotation)
		assertSettableParent(annotation).SetParent(s)
	}

	return s
}

// SetStatic updates the static flag of the method. This declares a method to be not part of the according
// instance and it will not be able to modify its receiver instance, so the PtrReceiver flag is ignored.
// In Java, this will cause the renderer to emit a class scoped method.
func (s *Func) SetStatic(b bool) *Func {
	s.FunStatic = b
	return s
}

// Static returns the static flag.
func (s *Func) Static() bool {
	return s.FunStatic
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (s *Func) Children() []Node {
	tmp := make([]Node, 0, len(s.FunParams)+len(s.FunResults)+1)
	for _, param := range s.FunParams {
		tmp = append(tmp, param)
	}

	for _, param := range s.FunResults {
		tmp = append(tmp, param)
	}

	if s.FunBody != nil {
		tmp = append(tmp, s.FunBody)
	}

	return tmp
}
