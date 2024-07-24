package util

// go test -v -count=1 -run ^TestCreateOneLiner$ ./util/

import (
	"strings"
	"testing"
)

func TestCreateOneLiner(t *testing.T) {
	expected := `good bye! bye bye!`
	b := string(CreateOneLiner(xml))
	t.Logf("S --> %s\n", b)
	if !strings.Contains(b, expected) {
		t.Errorf("%s doesn't contain %s", b, expected)
	}
}

func TestExtractAttr(t *testing.T) {
	m := extractattr([]byte(xml))
	t.Logf("%v\n", m)
	if len(m) != 4 {
		t.Errorf("expected 4 values map, got %v", m)
	}

	xml := maptoxml(m)
	t.Logf("%v\n", xml)
	if xml == `` {
		t.Error("expected XML string")
	}
}

var xml = `
<root>
    <greetings idd="456">hello</greetings>
    <greetings>good bye!
bye bye! <times>3</times>
    </greetings>
    <greetings id="123" msg="hello" secret="007"/>
    <smiles>wide</smiles>
</root>

`
