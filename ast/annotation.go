package ast

import (
	"github.com/golangee/src"
)

// AnnotationNode represents a field, type or method annotation.
type AnnotationNode struct {
	parent        Node
	srcAnnotation *src.Annotation
	*payload
}

// NewAnnotationNode wraps the given instance and creates a sub tree with parent/children relations to
// create a foundation for context-aware renderers.
func NewAnnotationNode(parent Node, a *src.Annotation) *AnnotationNode {
	return &AnnotationNode{
		parent:        parent,
		srcAnnotation: a,
		payload:       newPayload(),
	}
}

// SrcAnnotation returns the wrapped instance.
func (n *AnnotationNode) SrcAnnotation() *src.Annotation {
	return n.srcAnnotation
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *AnnotationNode) Parent() Node {
	return n.parent
}
