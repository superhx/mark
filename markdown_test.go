package marker

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestMarker(t *testing.T) {
	input, _ := ioutil.ReadFile("markdown_help.md")
	output, err := os.Create("markdown_help.html")
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
	markdown := NewMarker().Mark(input)
	writer := HTMLWriter{markdown}
	writer.WriteTo(output)
	output.WriteString("</body>")
	output.WriteString("<script src=\"http://cdnjs.cloudflare.com/ajax/libs/highlight.js/8.5/highlight.min.js\" ></script> <script>hljs.initHighlightingOnLoad();</script>")
	output.WriteString("</html>")
}
