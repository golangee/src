package ast2

// A Macro provides a dynamic amount of children and does not cause emitted source itself. It can be used
// to emit generator specific nodes at calling time based on the current ast state (especially its parents).
type Macro struct {
	Func func(m *Macro) []Node // Func should always return a new defensive copy with shared Node instances.
	Obj
}

func NewMacro() *Macro {
	return &Macro{}
}

// Target returns the available module target information. If there is no target available, all fields are default.
func (n *Macro) Target() Target {
	mod := &Mod{}
	if ok := ParentAs(n, &mod); ok {
		return mod.Target
	}

	return Target{}
}

// Children just delegates to Func.
func (n *Macro) Children() []Node {
	if n.Func == nil {
		return nil
	}

	return n.Func(n)
}

type targetFunc struct {
	target Target
	f      func(m *Macro) []Node
}

func (tf targetFunc) matches(t Target) bool {
	if tf.target.Equals(t) {
		return true
	}

	// if only lang is set, just match that, ignore everything else
	if t.MinLangVersion == "" && t.MaxLangVersion == "" && t.Framework == "" && tf.target.Lang == t.Lang && t.Arch == "" && t.Os == "" {
		return true
	}

	return false
}

// A MacroBuilder helps to create a target dependent macro.
type MacroBuilder struct {
	macros []targetFunc
	m      *Macro
}

// NewMacroBuilder create a new builder to pick between multiple targets.
func NewMacroBuilder() *MacroBuilder {
	b := &MacroBuilder{
		m: NewMacro(),
	}

	b.m.Func = func(m *Macro) []Node {
		target := m.Target()
		for _, macro := range b.macros {
			if macro.matches(target) {
				return macro.f(m)
			}
		}

		return []Node{NewComment("no macro match found")}
	}

	return b
}

// Add appends a target and the applied func. Empty fields will be ignored (like a * glob).
func (b *MacroBuilder) Add(t Target, f func(m *Macro) []Node) *MacroBuilder {
	b.macros = append(b.macros, targetFunc{
		target: t,
		f:      f,
	})

	return b
}

// AddLang appends a simplified catch-all language mapper which ever returns the given fixed set of nodes.
func (b *MacroBuilder) AddLang(lang Lang, nodes ...Node) *MacroBuilder {
	return b.Add(Target{
		Lang: lang,
	}, func(m *Macro) []Node {
		return nodes
	})
}

// Macro returns the internal macro instance.
func (b *MacroBuilder) Macro() *Macro {
	return b.m
}
