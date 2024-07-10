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

var xml = `
<root>
    <greetings>hello</greetings>
    <greetings>good bye!
bye bye! <times>3</times>
    </greetings>
    <greetings id="123"/>
    <smiles>wide</smiles>
</root>

`
