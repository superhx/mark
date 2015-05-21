package main

import (
	"fmt"
	"github.com/superhx/mark"
	"io/ioutil"
	"os"
	"regexp"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("lack arguement!")
		return
	}
	inputFile := os.Args[1]
	outputFile := ""
	if len(os.Args) == 2 {
		name := regexp.MustCompile(`^([\s\S]*)(?:\.[\s\S]*)$`).FindStringSubmatch(inputFile)
		if len(name) != 2 {
			fmt.Println("Invalid file name!")
			return
		}
		outputFile = name[1]
	} else {
		outputFile = os.Args[2]
	}
	outputFile += ".html"

	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Println("file read fail")
		return
	}
	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("file open fail")
	}
	defer func() {
		output.Close()
	}()
	output.WriteString("<html><head>")
	output.WriteString("<link rel=\"stylesheet\" href=\"https://cdn.lukej.me/github-markdown-css/0.2.0/github-markdown.min.css\">")
	output.WriteString("<link rel=\"stylesheet\" href=\"http://cdnjs.cloudflare.com/ajax/libs/highlight.js/8.5/styles/default.min.css\">")
	output.WriteString("<style> .markdown-body {  min-width: 200px;max-width: 790px;margin: 0 auto;padding: 30px;}</style>")
	output.WriteString("</head><body>")
	writer := mark.NewHTMLWriter(mark.Mark(input))
	writer.WriteTo(output)
	output.WriteString("</body>")
	output.WriteString("<script src=\"http://cdnjs.cloudflare.com/ajax/libs/highlight.js/8.5/highlight.min.js\" ></script> <script>hljs.initHighlightingOnLoad();</script>")
	output.WriteString("</html>")
}
