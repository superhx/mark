package markdown

import (
	"io"
	"strconv"
)

//Node is a interface
type Node interface {
	IsBlock() bool
	WriteToHTML(w io.Writer)
}

//Block is a Node which has child Nodes(Non-Terminal)
type Block struct {
}

// func (b *Block) String() string {
// 	mark := new(bytes.Buffer)
// 	for _, str := range b.children {
// 		mark.WriteString(str.String())
// 	}
// 	return mark.String()
// }

//IsBlock ...
func (b *Block) IsBlock() bool {
	return true
}

//Inline is a Node which do not have child(Terminal)
type Inline struct {
}

//IsBlock ...
func (inline Inline) IsBlock() bool {
	return false
}

//MarkDown ...
type MarkDown struct {
	*Block
	Parts []Node
}

// WriteToHTML ...
func (markdown MarkDown) WriteToHTML(w io.Writer) {
	w.Write([]byte("<article class=\"markdown-body\">\n"))
	for _, part := range markdown.Parts {
		part.WriteToHTML(w)
	}
	w.Write([]byte("\n</article>\n"))
}

//Space ...
type Space struct {
	*Inline
}

// WriteToHTML ...
func (space Space) WriteToHTML(w io.Writer) {
}

//Code ...
type Code struct {
	*Inline
	Lang string
	Text string
}

// WriteToHTML ...
func (code Code) WriteToHTML(w io.Writer) {
	w.Write([]byte("<pre><code class=\"" + code.Lang + "\">\n"))
	w.Write([]byte(code.Text))
	w.Write([]byte("\n</code></pre>\n"))
}

//Heading ...
type Heading struct {
	*Block
	Depth int
	Text  Text
}

// WriteToHTML ...
func (heading Heading) WriteToHTML(w io.Writer) {
	head := "h" + strconv.Itoa(heading.Depth)
	w.Write([]byte("<" + head + ">\n"))
	heading.Text.WriteToHTML(w)
	w.Write([]byte("\n</" + head + ">\n"))
}

//Nptable ...
type Nptable struct {
	*Block
	Header []string
	Align  []string
	Cells  [][]Text
}

// WriteToHTML ...
func (table Nptable) WriteToHTML(w io.Writer) {
	w.Write([]byte("<table>\n<thead>\n<tr>\n"))
	for _, head := range table.Header {
		w.Write([]byte("<td>" + head + "</td>\n"))
	}
	w.Write([]byte("\n</tr>\n</thead>\n<tbody>\n"))
	for _, line := range table.Cells {
		w.Write([]byte("<tr>\n"))
		for _, cell := range line {
			w.Write([]byte("<td>"))
			cell.WriteToHTML(w)
			w.Write([]byte("</td>\n"))
		}
		w.Write([]byte("</tr>\n"))
	}
	w.Write([]byte("</tbody>\n</table>\n"))
}

//Hr ...
type Hr struct {
	*Block
}

// WriteToHTML ...
func (hr Hr) WriteToHTML(w io.Writer) {
	w.Write([]byte("<hr>\n"))
}

//Blockquote ...
type Blockquote struct {
	*Block
	Parts []Node
}

// WriteToHTML ...
func (blockquote Blockquote) WriteToHTML(w io.Writer) {
	w.Write([]byte("<blockquote>\n"))
	for _, part := range blockquote.Parts {
		part.WriteToHTML(w)
	}
	w.Write([]byte("\n</blockquote>\n"))
}

//List ...
type List struct {
	*Block
	Items   []Node
	Ordered bool
}

// WriteToHTML ...
func (list List) WriteToHTML(w io.Writer) {
	if list.Ordered {
		w.Write([]byte("<ol>\n"))
	} else {
		w.Write([]byte("<ul>\n"))
	}
	for _, item := range list.Items {
		item.WriteToHTML(w)
	}
	if list.Ordered {
		w.Write([]byte("<ol>\n"))
	} else {
		w.Write([]byte("<ul>\n"))
	}
}

//Item ...
type Item struct {
	*List
	Parts []Node
}

// WriteToHTML ...
func (item Item) WriteToHTML(w io.Writer) {
	w.Write([]byte("<li>\n"))
	for _, part := range item.Parts {
		part.WriteToHTML(w)
	}
	if item.List != nil {
		item.List.WriteToHTML(w)
	}
	w.Write([]byte("\n</li>\n"))
}

//HTML ...
type HTML struct {
	*Block
	Text string
}

// WriteToHTML ...
func (html HTML) WriteToHTML(w io.Writer) {
	w.Write([]byte(html.Text))
}

//Def ...
type Def struct {
	*Block
	Href  string
	Title string
}

// WriteToHTML ...
func (def Def) WriteToHTML(w io.Writer) {

}

//Text ...
type Text struct {
	*Block
	Parts []Node
}

// WriteToHTML ...
func (text Text) WriteToHTML(w io.Writer) {
	w.Write([]byte("<p>\n"))
	for _, part := range text.Parts {
		part.WriteToHTML(w)
	}
	w.Write([]byte("\n</p>\n"))
}

//InlineText ...
type InlineText struct {
	*Inline
	Text string
}

// WriteToHTML ...
func (text InlineText) WriteToHTML(w io.Writer) {
	w.Write([]byte(text.Text))
}

//Link ...
type Link struct {
	*Inline
	Text string
	Href string
}

// WriteToHTML ...
func (link Link) WriteToHTML(w io.Writer) {
	w.Write([]byte("<a href=\"" + link.Href + "\">" + link.Text + "</a>"))
}

//Image ...
type Image struct {
	*Inline
	Text  string
	Href  string
	Title string
}

// WriteToHTML ...
func (image Image) WriteToHTML(w io.Writer) {
	w.Write([]byte("<img src=\"" + image.Href + "\" alt=\"" + image.Text + "\" title=\"" + image.Title + "\" >"))
}

//Strong ...
type Strong struct {
	*Inline
	Text string
}

// WriteToHTML ...
func (strong Strong) WriteToHTML(w io.Writer) {
	w.Write([]byte("<strong>" + strong.Text + "</strong>"))
}

//Em ...
type Em struct {
	*Inline
	Text string
}

// WriteToHTML ...
func (em Em) WriteToHTML(w io.Writer) {
	w.Write([]byte("<em>" + em.Text + "</em>"))
}

//Br ...
type Br struct {
	*Inline
}

// WriteToHTML ...
func (br Br) WriteToHTML(w io.Writer) {
	w.Write([]byte("<br/>"))
}

//InlineCode ...
type InlineCode struct {
	*Inline
	Text string
}

// WriteToHTML ...
func (code InlineCode) WriteToHTML(w io.Writer) {
	w.Write([]byte("<code>" + code.Text + "</code>"))
}

//Del ...
type Del struct {
	*Inline
	Text string
}

// WriteToHTML ...
func (del Del) WriteToHTML(w io.Writer) {
	w.Write([]byte("<del>" + del.Text + "</del>"))
}
