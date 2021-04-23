package lemmatizer

import (
	"fmt"
	"testing"
)

func TestLemmatize(t *testing.T) {
	lemmas := itLemmatizer.Lemmatize("finivano")
	fmt.Println(lemmas)
}
