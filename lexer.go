package marker

import (
	"fmt"
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"html"
	"regexp"
	"strings"
	"sync"
)

//Flag for Markdown
const (
	compileFlag    = pcre.UTF8
	matcherFlag    = pcre.PARTIAL_SOFT
	goroutineCount = 100
)

var block BlockRex
var inline InlineRex
var once sync.Once

//BlockRex ...
type BlockRex struct {
	newline, code, hr, heading, nptable, lheading, def, table, paragraph, text *regexp.Regexp
	fences, list, html, item, li, blockquote                                   pcre.Regexp
}

//InlineRex ...
type InlineRex struct {
	escape, autolink, url, tag, nolink             *regexp.Regexp
	link, reflink, strong, em, code, br, del, text pcre.Regexp
}

//Marker is a mark parser for markdown
type Marker struct {
	defs   map[string]Def
	relink []Node
	pool   chan bool
}

//Mark parse the markdown file,then return MarkDown Obj
func Mark(bytes []byte) (markdown *MarkDown) {
	once.Do(setUp)
	mark := &Marker{
		defs:   make(map[string]Def),
		relink: []Node{},
		pool:   make(chan bool, goroutineCount),
	}

	bytes = regexp.MustCompile("\r\n|\r").ReplaceAll(bytes, []byte("\n"))
	bytes = regexp.MustCompile("\u00a0").ReplaceAll(bytes, []byte("    "))
	bytes = regexp.MustCompile("\u2424").ReplaceAll(bytes, []byte("\n"))
	markdown = &MarkDown{Parts: mark.parse(bytes)}
	for i := 0; i < goroutineCount; i++ {
		mark.pool <- false
	}
	mark.link()
	return
}

func (mark *Marker) parse(bytes []byte) []Node {
	nodes := []Node{}
	for len(bytes) > 0 {
		//newline
		if node := block.newline.Find(bytes); node != nil {
			bytes = bytes[len(node):]
			if len(node) > 1 {
				nodes = append(nodes, &Space{})
			}
		}

		//code
		if node := block.code.Find(bytes); node != nil {
			bytes = bytes[len(node):]
			node = regexp.MustCompile(`^ {4}`).ReplaceAll(node, []byte(""))
			node := removeEndNewline(node)
			nodes = append(nodes, &Code{Text: html.EscapeString(string(node))})
			continue
		}

		//fences
		if matcher := block.fences.Matcher(bytes, matcherFlag); matcher.Matches() {
			bytes = bytes[len(matcher.Group(0)):]
			nodes = append(nodes, &Code{Lang: html.EscapeString(matcher.GroupString(2)), Text: html.EscapeString(matcher.GroupString(3))})
			continue
		}

		//heading
		if node := block.heading.FindSubmatch(bytes); node != nil {
			bytes = bytes[len(node[0]):]
			text := &Text{}
			mark.inlineParse(node[2], text)
			nodes = append(nodes, &Heading{Depth: len(node[1]), Text: text})
			continue
		}

		//nptable is table no leading pipe
		if node := block.nptable.FindSubmatch(bytes); node != nil {
			bytes = bytes[len(node[0]):]
			header := regexp.MustCompile(` *\| *`).Split(string(regexp.MustCompile(`^ *| *\| *$`).ReplaceAll(node[1], []byte(""))), -1)
			table := Nptable{
				Header: header,
			}

			align := regexp.MustCompile(` *\| *`).Split(string(regexp.MustCompile(`^ *|\| *$`).ReplaceAll(node[2], []byte(""))), -1)
			for i, length := 0, len(align); i < length; i++ {
				if matched, _ := regexp.MatchString(`^ *-+: *$`, align[i]); matched {
					align[i] = "right"
				} else if matched, _ = regexp.MatchString(`^ *:-+: *$`, align[i]); matched {
					align[i] = "center"
				} else if matched, _ = regexp.MatchString(`^ *:-+ *$`, align[i]); matched {
					align[i] = "left"
				} else {
					align[i] = ""
				}
			}
			table.Align = align

			line := strings.Split(string(regexp.MustCompile(`\n$`).ReplaceAll(node[3], []byte(""))), "\n")
			cells := make([][]*Text, len(line))
			for i, length, spliter := 0, len(line), regexp.MustCompile(` *\| *`); i < length; i++ {
				temp := spliter.Split(line[i], -1)
				cells[i] = make([]*Text, len(temp))
				for j, length := 0, len(temp); j < length; j++ {
					text := &Text{}
					mark.inlineStringParse(temp[j], text)
					cells[i][j] = text
				}
			}
			table.Cells = cells
			nodes = append(nodes, &table)
			continue
		}

		//lheading
		if node := block.lheading.FindSubmatch(bytes); node != nil {
			bytes = bytes[len(node[0]):]
			depth := 2
			if string(node[2]) == "=" {
				depth = 1
			}
			text := &Text{}
			mark.inlineParse(node[1], text)
			nodes = append(nodes, &Heading{
				Text:  text,
				Depth: depth,
			})
			continue
		}

		//hr
		if node := block.hr.Find(bytes); node != nil {
			bytes = bytes[len(node):]
			nodes = append(nodes, &Hr{})
			continue
		}

		//blockquote
		if matcher := block.blockquote.Matcher(bytes, matcherFlag); matcher.Matches() {
			bytes = bytes[len(matcher.Group(0)):]
			children := mark.parse(regexp.MustCompile(`^ *> ?`).ReplaceAll(matcher.Group(0), []byte("")))
			nodes = append(nodes, &Blockquote{Parts: children})
			continue
		}

		// list
		if matcher := block.list.Matcher(bytes, matcherFlag); matcher.Matches() {
			listBytes := append(matcher.Group(1), '\n')
			bytes = bytes[len(matcher.Group(0)):]
			ordered := false
			bull := matcher.GroupString(3)
			if !(bull == "*" || bull == "+" || bull == "-") {
				ordered = true
			}
			list := List{Items: []*Item{}, Ordered: ordered}
			for matcher := block.item.Matcher(listBytes, matcherFlag); len(listBytes) > 0 && matcher.Matches(); matcher.Match(listBytes, matcherFlag) {
				listBytes = listBytes[len(matcher.Group(0)):]
				list.Items = append(list.Items, mark.subList([]byte(matcher.GroupString(0))))
			}
			nodes = append(nodes, &list)
			continue
		}

		//html
		if matcher := block.html.Matcher(bytes, matcherFlag); matcher.Matches() {
			bytes = bytes[len(matcher.Group(0)):]
			nodes = append(nodes, &HTML{Text: matcher.GroupString(0)})
			continue
		}

		//def
		if node := block.def.FindSubmatch(bytes); node != nil {
			bytes = bytes[len(node[0]):]
			mark.defs[string(node[1])] = Def{Href: html.EscapeString(string(node[2])), Title: html.EscapeString(string(node[3]))}
			continue
		}

		//table
		if node := block.table.FindSubmatch(bytes); node != nil {
			bytes = bytes[len(node[0]):]
			table := Nptable{
				Header: regexp.MustCompile(` *\| *`).Split(string(regexp.MustCompile(`^ *| *\| *$`).ReplaceAll(node[1], []byte(""))), -1),
			}

			align := regexp.MustCompile(` *\| *`).Split(string(regexp.MustCompile(`^ *|\| *$`).ReplaceAll(node[2], []byte(""))), -1)
			for i, length := 0, len(align); i < length; i++ {
				if matched, _ := regexp.MatchString(`^ *-+: *$`, align[i]); matched {
					align[i] = "right"
				} else if matched, _ = regexp.MatchString(`^ *:-+: *$`, align[i]); matched {
					align[i] = "center"
				} else if matched, _ = regexp.MatchString(`^ *:-+ *$`, align[i]); matched {
					align[i] = "left"
				} else {
					align[i] = ""
				}
			}
			table.Align = align

			line := strings.Split(string(regexp.MustCompile(`(?: *\| *)?\n$`).ReplaceAll(node[3], []byte(""))), "\n")
			cells := make([][]*Text, len(line))
			for i, length, replacer, spliter := 0, len(line), regexp.MustCompile(`^ *\| *| *\| *$`), regexp.MustCompile(` *\| *`); i < length; i++ {
				temp := spliter.Split(replacer.ReplaceAllString(line[i], ""), -1)
				cells[i] = make([]*Text, len(temp))
				for j, length := 0, len(temp); j < length; j++ {
					text := &Text{}
					mark.inlineStringParse(temp[j], text)
					cells[i][j] = text
				}
			}
			table.Cells = cells

			nodes = append(nodes, &table)
			continue
		}

		//text
		if node := block.text.Find(bytes); node != nil {
			bytes = bytes[len(node):]
			text := &Text{}
			mark.inlineParse(node, text)
			nodes = append(nodes, &BlockText{Text: text})
			continue
		}

		if len(bytes) > 0 {
			fmt.Println("Infinite loop on byte:\n" + string(bytes))
			break
		}
	}

	return nodes
}

//InlineParse ...
func (mark *Marker) InlineParse(bytes []byte) (text *Text) {
	mark.inlineSynParse(bytes, text)
	mark.link()
	return
}

func (mark *Marker) inlineSynParse(bytes []byte, text *Text) {
	parts := []Node{}
	for len(bytes) > 0 {
		//escape
		if cap := inline.escape.FindSubmatch(bytes); cap != nil {
			bytes = bytes[len(cap[0]):]
			parts = append(parts, &InlineText{Text: html.EscapeString(string(cap[1]))})
			continue
		}

		//autolink
		if cap := inline.autolink.FindSubmatch(bytes); cap != nil {
			bytes = bytes[len(cap[0]):]
			var text, href string
			if string(cap[2]) == "@" {
				if cap[1][6] == ':' {
					text = string(cap[1][7:])
				} else {
					text = string(cap[1])
				}
				href = "mainto:" + text
			} else {
				text = string(cap[1])
				href = text
			}
			parts = append(parts, &Link{Text: html.EscapeString(text), Href: html.EscapeString(href)})
			continue
		}

		//url
		if cap := inline.url.Find(bytes); cap != nil {
			bytes = bytes[len(cap):]
			text := string(cap)
			href := text
			parts = append(parts, &Link{Text: html.EscapeString(text), Href: html.EscapeString(href)})
			continue
		}

		//tag unsolved

		//link
		if matcher := inline.link.Matcher(bytes, matcherFlag); matcher.Matches() {
			bytes = bytes[len(matcher.Group(0)):]
			text := matcher.GroupString(1)
			href := matcher.GroupString(2)
			if matcher.Group(0)[0] != '!' {
				parts = append(parts, &Link{Text: text, Href: html.EscapeString(href), Title: html.EscapeString(matcher.GroupString(3))})
			} else {
				parts = append(parts, &Image{Text: html.EscapeString(text), Href: html.EscapeString(href), Title: matcher.GroupString(3)})
			}
			continue
		}

		//relink nolink unsolved
		if matcher := inline.reflink.Matcher(bytes, matcherFlag); matcher.Matches() {
			bytes = bytes[len(matcher.Group(0)):]
			text := matcher.GroupString(1)
			href := matcher.GroupString(2)
			var node Node
			if matcher.Group(0)[0] != '!' {
				node = &Link{Text: text, Href: html.EscapeString(href)}
			} else {
				node = &Image{Text: html.EscapeString(text), Href: html.EscapeString(href)}
			}
			parts = append(parts, node)
			mark.relink = append(mark.relink, node)

		}

		//strong
		if matcher := inline.strong.Matcher(bytes, matcherFlag); matcher.Matches() {
			bytes = bytes[len(matcher.Group(0)):]
			text := matcher.GroupString(1) + matcher.GroupString(2)
			parts = append(parts, &Strong{Text: html.EscapeString(text)})
			continue
		}

		//em
		if matcher := inline.em.Matcher(bytes, matcherFlag); matcher.Matches() {
			bytes = bytes[len(matcher.Group(0)):]
			text := matcher.GroupString(1) + matcher.GroupString(2)
			parts = append(parts, &Em{Text: html.EscapeString(text)})
			continue
		}

		//code
		if matcher := inline.code.Matcher(bytes, matcherFlag); matcher.Matches() {
			bytes = bytes[len(matcher.Group(0)):]
			parts = append(parts, &InlineCode{Text: html.EscapeString(matcher.GroupString(2))})
			continue
		}

		//br
		if matcher := inline.br.Matcher(bytes, matcherFlag); matcher.Matches() {
			bytes = bytes[len(matcher.Group(0)):]
			parts = append(parts, &Br{})
			continue
		}

		//del
		if matcher := inline.del.Matcher(bytes, matcherFlag); matcher.Matches() {
			bytes = bytes[len(matcher.Group(0)):]
			parts = append(parts, &Del{Text: html.EscapeString(matcher.GroupString(1))})
		}

		//text
		if matcher := inline.text.Matcher(bytes, matcherFlag); matcher.Matches() {
			bytes = bytes[len(matcher.Group(0)):]
			parts = append(parts, &InlineText{Text: html.EscapeString(matcher.GroupString(0))})
		}
	}
	text.Parts = parts
}

func (mark *Marker) inlineParse(bytes []byte, text *Text) {
	mark.pool <- true
	go func() {
		mark.inlineSynParse(bytes, text)
		<-mark.pool
	}()
}

func (mark *Marker) inlineStringParse(str string, text *Text) {
	mark.inlineParse([]byte(str), text)
}

func (mark *Marker) subList(bytes []byte) *Item {
	matcher := block.li.Matcher(bytes, matcherFlag)
	node := matcher.Group(0)
	bytes = bytes[len(node):]
	if len(node) == len(bytes) {
		return &Item{Parts: mark.parse(matcher.Group(3))}
	}
	item := &Item{&List{Items: []*Item{}}, mark.parse([]byte(matcher.GroupString(3) + "\n"))}
	matcher = block.item.Matcher(bytes, matcherFlag)
	bull := matcher.GroupString(2)
	if !(bull == "*" || bull == "+" || bull == "-") {
		item.Ordered = true
	}
	for ; len(bytes) > 0 && matcher.Matches(); matcher.Match(bytes, matcherFlag) {
		bytes = bytes[len(matcher.Group(0)):]
		item.Items = append(item.Items, mark.subList(matcher.Group(0)))
	}
	return item
}

func (mark *Marker) link() {
	for _, node := range mark.relink {
		switch node.(type) {
		case *Link:
			link := node.(*Link)
			def := mark.defs[link.Href]
			link.Href = def.Href
			link.Title = def.Title
		case *Image:
			image := node.(*Image)
			def := mark.defs[image.Href]
			image.Href = def.Href
			image.Title = def.Title
		}
	}
}

func removeEndNewline(bytes []byte) []byte {
	for i := len(bytes) - 1; i >= 0; i-- {
		if bytes[i] != '\n' {
			bytes = bytes[:i]
			return bytes
		}
	}
	return nil
}

func setUp() {
	bullet := `(?:[*+-]|\d+\.)`
	tag := `(?!(?:a|em|strong|small|s|cite|q|dfn|abbr|data|time|code|var|samp|kbd|sub|sup|i|b|u|mark|ruby|rt|rp|bdi|bdo|span|br|wbr|ins|del|img)\b)\w+(?!:/|[^\w\s@]*@)\b`
	item := `^( *)(bull) ([^\n]*(\n(?!\1bull )[^\n]*)*)(?:\n+|$)`
	item = strings.Replace(item, "bull", bullet, -1)

	li := `^( *)(bull) ([^\n]*(\n(?! *bull )[^\n]*)*)(?:\n+|$)`
	li = strings.Replace(li, "bull", bullet, -1)

	hr := `^( *[-*_]){3,} *(?:\n+|$)`

	heading := `^ *(#{1,6}) *([^\n]+?) *#* *(?:\n+|$)`

	lheading := `^([^\n]+)\n *(=|-){2,} *(?:\n+|$)`

	def := `^ *\[([^\]]+)\]: *<?([^\s>]+)>?(?: +["(]([^\n]+)[")])? *(?:\n+|$)`

	blockquote := `^( *>[^\n]+(\n(?!def)[^\n]+)*\n*)+`
	blockquote = strings.Replace(blockquote, "def", def, -1)

	list := `^(( *)(bull) [\s\S]+?)(?:hr|def|\n{2,}(?! )(?!\2bull )\n*|\s*$)`
	list = strings.Replace(list, "bull", bullet, -1)
	list = strings.Replace(list, "hr", `\n+(?=\2?(?:[-*_] *){3,}(?:\n+|$))`, -1)
	list = strings.Replace(list, "def", `\n+(?=`+def+`)`, -1)

	html := `^ *(?:comment *(?:\n|\s*$)|closed *(?:\n{2,}|\s*$)|closing *(?:\n{2,}|\s*$))`
	html = strings.Replace(html, "comment", `<!--[\s\S]*?-->`, -1)
	html = strings.Replace(html, "closed", `<(tag)[\s\S]+?<\/\1>`, -1)
	html = strings.Replace(html, "closing", `<tag(?:"[^"]*"|'[^']*'|[^'">])*?>`, -1)
	html = strings.Replace(html, `tag`, tag, -1)

	paragraph := `^((?:[^\n]+\n?(?!hr|heading|lheading|blockquote|tag|def))+)\n*`
	paragraph = strings.Replace(paragraph, "hr", hr, -1)
	paragraph = strings.Replace(paragraph, "heading", heading, -1)
	paragraph = strings.Replace(paragraph, "lheading", lheading, -1)
	paragraph = strings.Replace(paragraph, "blockquote", blockquote, -1)
	paragraph = strings.Replace(paragraph, "tag", tag, -1)
	paragraph = strings.Replace(def, "def", def, -1)

	inside := `(?:\[[^\]]*\]|[^\[\]]|\](?=[^\[]*\]))*`
	href := `\s*<?([\s\S]*?)>?(?:\s+['"]([\s\S]*?)['"])?\s*`

	link := `^!?\[(inside)\]\(href\)`
	link = strings.Replace(link, "inside", inside, -1)
	link = strings.Replace(link, "href", href, -1)

	relink := `^!?\[(inside)\]\s*\[([^\]]*)\]`
	relink = strings.Replace(relink, "inside", inside, -1)

	block = BlockRex{
		newline:    regexp.MustCompile(`^\n+`),
		code:       regexp.MustCompile(`^(( {4}| {2}|\t)[^\n]+\n*)+`),
		fences:     pcre.MustCompile("^ *(`{3,}|~{3,}) *(\\S+)? *\\n([\\s\\S]+?)\\s*\\1 *(?:\\n+|$)", compileFlag),
		hr:         regexp.MustCompile(hr),
		heading:    regexp.MustCompile(heading),
		nptable:    regexp.MustCompile(`^ *(\S.*\|.*)\n *([-:]+ *\|[-| :]*)\n((?:.*\|.*(?:\n|$))*)\n*`),
		lheading:   regexp.MustCompile(`^([^\n]+)\n *(=|-){2,} *(?:\n+|$)`),
		blockquote: pcre.MustCompile(blockquote, compileFlag),
		list:       pcre.MustCompile(list, compileFlag),
		html:       pcre.MustCompile(html, compileFlag),
		def:        regexp.MustCompile(def),
		table:      regexp.MustCompile(`^ *\|(.+)\n *\|( *[-:]+[-| :]*)\n((?: *\|.*(?:\n|$))*)\n`),
		paragraph:  regexp.MustCompile(paragraph),
		text:       regexp.MustCompile(`^[^\n]+`),
		item:       pcre.MustCompile(item, compileFlag),
		li:         pcre.MustCompile(li, compileFlag),
	}

	inline = InlineRex{
		escape:   regexp.MustCompile("^\\\\([\\\\`*{}\\[\\]()#+\\-.!_>])"),
		autolink: regexp.MustCompile(`^<([^ >]+(@|:\/)[^ >]+)>`),
		url:      regexp.MustCompile(`^(https?:\/\/[^\s<]+[^<.,:;"')\]\s])`),
		tag:      regexp.MustCompile(`^<!--[\s\S]*?-->|^<\/?\w+(?:"[^"]*"|'[^']*'|[^'">])*?>`),
		link:     pcre.MustCompile(link, compileFlag),
		reflink:  pcre.MustCompile(relink, compileFlag),
		nolink:   regexp.MustCompile(`^!?\[((?:\[[^\]]*\]|[^\[\]])*)\]`),
		strong:   pcre.MustCompile(`^__([\s\S]+?)__(?!_)|^\*\*([\s\S]+?)\*\*(?!\*)`, compileFlag),
		em:       pcre.MustCompile(`^\b_((?:__|[\s\S])+?)_\b|^\*((?:\*\*|[\s\S])+?)\*(?!\*)`, compileFlag),
		code:     pcre.MustCompile("^(`+)\\s*([\\s\\S]*?[^`])\\s*\\1(?!`)", compileFlag),
		br:       pcre.MustCompile(`^ {2,}\n(?!\s*$)`, compileFlag),
		del:      pcre.MustCompile(`^~~(?=\S)([\s\S]*?\S)~~`, compileFlag),
		text:     pcre.MustCompile("^[\\s\\S]+?(?=[\\\\<!\\[_*`]| {2,}\\n|$)", compileFlag),
	}
}
