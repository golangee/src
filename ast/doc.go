// Package ast provides a Node API around the src.* types. Intentionally, we do not want a complex
// parent-children relation in the base model, which would be very error prone. Indeed, the first draft
// was build exactly that way and was usually misused, causing all kinds of weired side effects, due to
// reused instances. Also, polluting the model with render semantics avoid a clear separation of concerns and
// the support of other language renderers. That is the reason why this ast package provides a short-lived
// parent-children Node API. You may just want to start with NewModNode.
package ast
