package golang

import (
	"fmt"
	"github.com/golangee/src/ast"
	"github.com/golangee/src/golang/validate"
	"github.com/golangee/src/render"
	"strings"
)

// renderError emits an idiomatic way of representing behavior-based errors instead types. See
// also https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully and
// https://dave.cheney.net/2014/12/24/inspecting-errors.
func (r *Renderer) renderError(node *ast.Error, w *render.BufferedWriter) error {
	importer := r.importer(node)

	// TODO actually emitting this interface is pointless and verbose in Go, however it describes the architectural intention clearly.
	r.writeCommentNode(w, false, node.Identifier(), node.Comment())
	if err := validate.ExportedIdentifier(node.Visibility(), node.Identifier()); err != nil {
		return err
	}

	sumTypeMarkerFunc := MakePublic(node.GroupName)
	const markerDoc = "...is the according marker function."

	w.Printf("type %s interface {\n", node.GroupName)
	r.writeComment(w, false, sumTypeMarkerFunc, markerDoc)
	w.Printf("  %s() bool\n", sumTypeMarkerFunc)
	w.Printf("  error")
	w.Printf("}\n")

	for _, eCase := range node.Cases {
		r.writeCommentNode(w, false, eCase.Identifier(), eCase.Comment())
		if err := validate.ExportedIdentifier(eCase.Visibility(), eCase.Identifier()); err != nil {
			return err
		}

		w.Printf("type %s interface {\n", node.GroupName+eCase.TypeName)

		fName := MakePublic(eCase.TypeName)
		r.writeComment(w, false, fName, markerDoc)
		w.Printf("  %s() bool\n", fName)
		for _, property := range eCase.Properties {
			if property.Read.Enabled {
				fName := MakePublic(property.FieldName)
				r.writeComment(w, false, fName, "...returns the "+property.FieldName+" and "+strings.TrimSpace(DeEllipsis("", property.CommentText())))
				w.Printf("%s()", fName)
				if err := r.renderTypeDecl(property.FieldType, w); err != nil {
					return fmt.Errorf("unable to render property getter: %w", err)
				}
				w.Printf("\n")
			}
		}

		w.Printf("  %s\n", node.GroupName)
		w.Printf("}\n")
	}

	for _, eCase := range node.Cases {
		caseTypeMarkerFunc := MakePublic(eCase.TypeName)
		tName := MakePrivate(node.Identifier() + eCase.TypeName)
		recName := strings.ToLower(tName[:1])

		// renders fields
		r.writeComment(w, false, tName, "...is a package-private implementation for a "+caseTypeMarkerFunc+" error.")
		w.Printf("type %s struct {\n", tName)
		r.writeComment(w, false, "cause", "...optional error cause to unwrap.")
		w.Printf("  cause error\n")
		for _, property := range eCase.Properties {
			if err := validate.ExportedIdentifier(ast.Private, property.FieldName); err != nil {
				return err
			}

			r.writeCommentNode(w, false, property.FieldName, property.Comment())
			w.Printf("%s ", property.FieldName)
			if err := r.renderTypeDecl(property.FieldType, w); err != nil {
				return fmt.Errorf("unable to render property field: %w", err)
			}

			w.Printf("\n")
		}

		w.Printf("}\n")

		// render marker interface
		r.writeComment(w, false, sumTypeMarkerFunc, markerDoc)
		w.Printf("func (%s *%s) %s() bool{\nreturn true\n}\n", recName, tName, sumTypeMarkerFunc)
		r.writeComment(w, false, caseTypeMarkerFunc, markerDoc)
		w.Printf("func (%s *%s) %s() bool{\nreturn true\n}\n\n", recName, tName, caseTypeMarkerFunc)
		r.writeComment(w, false, "Unwrap", "...unpacks the cause or returns nil.")
		w.Printf("func (%s *%s) Unwrap() error {\nreturn %s.cause\n}\n", recName, tName, recName)

		// render error func
		r.writeComment(w, false, "Error", "...conforms to the error behavior.")
		w.Printf("func (%s *%s) Error() string {\n", recName, tName)
		w.Printf(" return %s(", importer.shortify("fmt.Sprintf"))
		w.Printf(`"`)
		w.Printf(eCase.TypeName)
		for _, property := range eCase.Properties {
			w.Printf(" %s=%%v", property.FieldName)
		}
		w.Printf(`",`)
		for _, property := range eCase.Properties {
			w.Printf("%s.%s,", recName, property.FieldName)
		}
		w.Printf(")\n")
		w.Printf("}\n")

		// render getters and setters
		for _, property := range eCase.Properties {
			if property.Read.Enabled {
				fName := MakePublic(property.FieldName)
				r.writeComment(w, false, fName, "...returns the "+property.FieldName+" and "+strings.TrimSpace(DeEllipsis("", property.CommentText())))
				w.Printf("func (%s *%s) %s() ", recName, tName, fName)
				if err := r.renderTypeDecl(property.FieldType, w); err != nil {
					return fmt.Errorf("unable to render getter: %w", err)
				}
				w.Printf("{\n")
				w.Printf("return %s.%s\n", recName, property.FieldName)
				w.Printf("}\n")
			}

			if property.Write.Enabled {
				fName := "Set" + MakePublic(property.FieldName)
				pName := strings.ToLower(property.FieldName[:1])
				if pName == recName {
					pName = "value"
				}

				r.writeComment(w, false, fName, "...set the "+property.FieldName+" and "+strings.TrimSpace(DeEllipsis("", property.CommentText())))
				r.writeCommentNode(w, false, fName, property.Comment())
				w.Printf("func (%s *%s) %s( ", recName, tName, fName)
				if err := r.renderTypeDecl(property.FieldType, w); err != nil {
					return fmt.Errorf("unable to render getter: %w", err)
				}

				w.Printf(" %s)\n{", pName)
				w.Printf("  %s.%s = %s\n", recName, property.FieldName, pName)
				w.Printf("}\n")
			}
		}

	}

	return nil
}
