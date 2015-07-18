##Mark
---
**Mark** is a [markdown](http://daringfireball.net/projects/markdown/syntax) parser implemented in [Go](http://golang.org/). It's fast and support extension such as Table,fence code,etc.It is safe for utf-8 input.

The parser idea comes from [chjj/marked](https://github.com/chjj/marked). But it isn't as complete as chjj/marked

###Installation
---
Mark is compatible with go1.4.2 (I don't know whether compatible with go version under 1.4.2)

With go and git installed

Install pcre
```
sudo apt-get install libpcre++-dev // In Ubuntu
brew install pcre				   // In Mac OS X
```
Get markdown parser

``` golang
go get github.com/superhx/mark
```
###Usage
---
####Code

Basic usage, read a bytes then parse it to `Markdown` , render `Markdown` and output
``` golang
//marker parse the input and return a Markdown object
mark:=mark.Mark(input)
//writer io.Writer
//render markdown to html and output to writer(without style sheet,see markdown_test to pretty)
writer := mark.NewHTMLWriter(mark)
writer.WriteTo(output)
```
The `Markdown` is like a tree (dom tree)

If you want to operate the `Markdown` instead of only simple output it. It is all free for you to modify the Markdown tree as you want. The **markdown.go** contain all `struct` in Markdown tree.

###Markdown Grammar Support
---
- [Basic Markdown](http://daringfireball.net/projects/markdown/syntax)
- [GitHub Flavored Markdown (gfm)](https://help.github.com/articles/github-flavored-markdown/)
