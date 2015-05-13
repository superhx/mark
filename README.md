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
go get github.com/superhx/markdown
```
###Usage
---
####Code

Basic usage, read a bytes ,then parse it to `Markdown` , render `Markdown` and output
``` golang
//new and initiate a Marker
marker:=marker.NewMarker()
//input []byte is the markdown input ,
//marker parse the input and return a Markdown object
mark:=marker.Mark(input)
//writer io.Writer
//render markdown to html and output to writer(without style sheet,see markdown_test to pretty)
writer := marker.NewHTMLWriter(mark)
writer.WriteTo(output)
```
The `Markdown` is like a tree (dom tree)

If you want to operate the `Markdown` instead of only simple output it. It is all free for you to modify the Markdown tree as you want. The **markdown.go** contain all `struct` in Markdown tree.

####Markdown transform to HTML

1. Get the **mark** executable file
``` sh
$ cd $GOPATH/src/github.com/superhx/markdown/main
$ go build
```
 or get from google drive
  - [Mac OS X](https://drive.google.com/file/d/0B3wRzs_xbfwQWUt0OEZOMjdjd1U/view?usp=sharing)
  - [Linux](https://drive.google.com/file/d/0B3wRzs_xbfwQTjFCS0M4YTZ5SVE/view?usp=sharing)
  - Windows

2. Tranform markdown file to html file
``` sh
$ ./mark input_file_path output_file_path
```
or you can just emit the output file path
```
$ ./mark input_file_path
```

###Markdown Grammar Support
---
- [Basic Markdown](http://daringfireball.net/projects/markdown/syntax)
- [GitHub Flavored Markdown (gfm)](https://help.github.com/articles/github-flavored-markdown/)


<div><a href="https://github.com/superhx/markdown"><img style="position: absolute; top: 0; right: 0; border: 0; width: 149px; height: 149px;" src="http://aral.github.com/fork-me-on-github-retina-ribbons/right-graphite@2x.png" alt="Fork me on GitHub"></a></div>
