package mark

import (
	"bytes"
	//"os"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestMarker(t *testing.T) {
	f, err := ioutil.ReadFile("testdata.json")
	if err != nil {
		fmt.Println(err)
		// assert.Fail(t, err, msgAndArgs ...interface{})
	}
	var testData []map[string]string
	err = json.Unmarshal(f, testData)
	if err != nil {

	}
	// assertEqual(t, "\n\n", "")
	//
	// assertEqual(t, "`inline code`", "<p><pre><code>inline code</code></pre></p>")
	//
	// assertEqual(t, "#heading1", "<h1>heading1</h1>")
	// assertEqual(t, "heading1\n===", "<h1>heading1</h1>")
	// assertEqual(t, "heading2\n---", "<h2>heading2</h2>")
	//
	// assertEqual(t, "first header|second header\n---|---\ncontent cell|content cell\n---|---\ncontent cell", "<table><thead><tr><th>first header</th><th>second header</th></tr></thead><tbody></tbody></table>")
	//
	// assertEqual(t, "- this is a list that only have one li", "<ul><li><p>this is a list that only have one li</p></li></ul>")
	// assertEqual(t, "- this is a list that only have one li with blank append   ", "<ul><li><p>this is a list that only have one li with blank append</p></li></ul>")
	// assertEqual(t, "-    this is a list that only have one li with blank prepend", "<ul><li><p>this is a list that only have one li with blank prepend</p></li></ul>")
	// assertEqual(t, "- parent\n  - child\n  - child", "<ul><li><p>parent</p><ul><li><p>child</p></li><li><p>child</p></li></ul></li></ul>")

}

func assertEqual(t assert.TestingT, input string, oracle string) {
	out := bytes.NewBufferString("")
	NewHTMLWriter(Mark([]byte(input))).WriteTo(out)
	output := out.String()
	assert.EqualValues(t, oracle, output[31:len(output)-10])
}
