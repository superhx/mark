package markdown

//Node is a interface
type Node interface {
}

//MarkDown ...
type MarkDown struct {
	Parts []Node
}

//Space ...
type Space struct {
}

//Code ...
type Code struct {
	Lang string
	Text string
}

//Heading ...
type Heading struct {
	Depth int
	Text  *Text
}

//Nptable ...
type Nptable struct {
	Header []string
	Align  []string
	Cells  [][]*Text
}

//Hr ...
type Hr struct {
}

//Blockquote ...
type Blockquote struct {
	Parts []Node
}

//List ...
type List struct {
	Items   []*Item
	Ordered bool
}

//Item ...
type Item struct {
	*List
	Parts []Node
}

//HTML ...
type HTML struct {
	Text string
}

//Def ...
type Def struct {
	Href  string
	Title string
}

//BlockText ...
type BlockText struct {
	*Text
}

//Text ...
type Text struct {
	Parts []Node
}

//InlineText ...
type InlineText struct {
	Text string
}

//Link ...
type Link struct {
	Text  string
	Href  string
	Title string
}

//Image ...
type Image struct {
	Text  string
	Href  string
	Title string
}

//Strong ...
type Strong struct {
	Text string
}

//Em ...
type Em struct {
	Text string
}

//Br ...
type Br struct {
}

//InlineCode ...
type InlineCode struct {
	Text string
}

//Del ...
type Del struct {
	Text string
}
