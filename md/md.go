package md

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	ast "github.com/gomarkdown/markdown/ast"
	ps "github.com/gomarkdown/markdown/parser"
	"github.com/rivo/tview"
)

// HighlightStyle determines the syntax highlighting colorstyle:
// https://xyproto.github.io/splash/docs/all.html
var HighlightStyle = "vs"

// ExtensionFlags is the flags used for markdown parsing
var ExtensionFlags = ps.FencedCode |
	ps.Autolink |
	ps.HardLineBreak |
	ps.BackslashLineBreak |
	ps.Strikethrough |
	ps.NoEmptyLineBeforeBlock

var trashyCodeBlockMatching = regexp.MustCompile("(.)```")

// Parse parses md into tview string
func Parse(s string) string {
	return ParseNoEscape(tview.Escape(s))
}

func getTextNodeContainer(t *ast.Text) *ast.Container {
	parent := t.Parent
	if parent != nil {
		return parent.AsContainer()
	}

	return nil
}

// ParseNoEscape parses md into tview string without escaping it
func ParseNoEscape(s string) string {
	b := strings.Builder{}

	parser := ps.NewWithExtensions(ExtensionFlags)
	if parser == nil {
		return s
	}

	s = trashyCodeBlockMatching.ReplaceAllString(s, "$1\n```")
	s = fixQuotes(s)

	doc := parser.Parse([]byte(s))
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		switch node := node.(type) {
		case *ast.Softbreak:
			b.WriteRune('\n')
		case *ast.Hardbreak:
			b.WriteRune('\n')
		case *ast.Emph:
			b.WriteString(isFormatEnter(entering, "b"))
			b.Write(node.Content)
		case *ast.Strong:
			b.WriteString(isFormatEnter(entering, "b"))
			b.Write(node.Content)
		case *ast.Del:
			b.WriteString(isFormatEnter(entering, "d"))
			b.Write(node.Content)
		case *ast.BlockQuote:
			// skipped, check ast.Text
		case *ast.Link:
			b.WriteString(isFormatEnter(entering, "u"))
			b.Write(node.Title)
		case *ast.Code:
			b.WriteString("[:black:]")
			b.Write(node.Literal)
			b.WriteString("[:-:]")
		case *ast.CodeBlock:
			var lexer = lexers.Fallback
			if lang := string(node.Info); lang != "" {
				if l := lexers.Get(lang); l != nil {
					lexer = l
				}
			}

			var fmtter = formatters.Get("tview-256bit")
			if fmtter == nil {
				fmtter = formatters.Fallback
			}

			var style = styles.Get(HighlightStyle)
			if style == nil {
				style = styles.Fallback
			}

			content := strings.TrimFunc(
				string(node.Literal),
				func(r rune) bool {
					return r == '\n'
				},
			)

			iterator, err := lexer.Tokenise(nil, content)
			if err != nil {
				b.Write(node.Literal)
				break
			}

			code := strings.Builder{}

			if err := fmtter.Format(&code, style, iterator); err != nil {
				b.Write(node.Literal)
				break
			}

			var s string
			for _, l := range strings.Split(code.String(), "\n") {
				s += "\n[grey]┃[-] " + l
			}

			if !strings.HasSuffix(s, "\n") {
				s += "\n"
			}

			b.WriteString(s)
		case *ast.ListItem:
			if entering {
				// if it's an ordered list and not bullets
				if (node.ListFlags & ast.ListTypeOrdered) != 0 {
					if p, ok := node.Parent.(*ast.List); ok {
						for i, c := range p.Children {
							if c != node {
								continue
							}

							b.WriteString(
								fmt.Sprintf(
									"%d%s %s", p.Start+i,
									string(node.Delimiter), string(node.Literal),
								),
							)
						}
					}
				} else {
					b.WriteRune(rune(node.BulletChar))
					b.WriteRune(' ')
				}
			}
		case *ast.Aside:
			b.Write(node.Literal)
		case *ast.CrossReference:
			b.Write(node.Literal)
		case *ast.Citation:
			b.Write(node.Literal)
		case *ast.Image:
			b.Write(node.Literal)
		case *ast.Caption:
			b.Write(node.Literal)
		case *ast.CaptionFigure:
			b.Write(node.Literal)
		case *ast.Document:
			b.Write(node.Literal)
		case *ast.Paragraph:
			b.Write(node.Literal)
		case *ast.HTMLSpan:
			b.Write(node.Literal)
		case *ast.HTMLBlock:
			b.Write(node.Literal)
		case *ast.Heading:
			b.Write(node.Literal)
		case *ast.HorizontalRule:
			b.Write(node.Literal)
		case *ast.List:
			b.Write(node.Literal)
		case *ast.Table:
			b.Write(node.Literal)
		case *ast.TableCell:
			b.Write(node.Literal)
		case *ast.TableHeader:
			b.Write(node.Literal)
		case *ast.TableBody:
			b.Write(node.Literal)
		case *ast.TableRow:
			b.Write(node.Literal)
		case *ast.TableFooter:
			b.Write(node.Literal)
		case *ast.Math:
			b.Write(node.Literal)
		case *ast.MathBlock:
			b.Write(node.Literal)
		case *ast.DocumentMatter:
			b.Write(node.Literal)
		case *ast.Callout:
			b.Write(node.Literal)
		case *ast.Index:
			b.Write(node.Literal)
		case *ast.Subscript:
			b.Write(node.Literal)
		case *ast.Superscript:
			b.Write(node.Literal)
		case *ast.Footnotes:
			b.Write(node.Literal)
		case *ast.Text:
			if ct := getTextNodeContainer(node); ct != nil {
				switch ct.Parent.(type) {
				case *ast.BlockQuote:
					b.WriteString("[green]>")
					b.Write(node.Literal)
					b.WriteString("[-]")
				default:
					b.Write(node.Literal)
				}
			} else {
				b.Write(node.Literal)
			}
		default:

		}

		return ast.GoToNext
	})

	return b.String()
}
