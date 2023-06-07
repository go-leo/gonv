package stringx_test

import (
	"github.com/go-leo/gox/stringx"
	"testing"
)

type str string

func TestBlank(t *testing.T) {
	t.Log(stringx.IsNotBlank(""))
	t.Log(stringx.IsNotBlank(" "))
	t.Log(stringx.IsNotBlank("	 "))
	t.Log(stringx.IsNotBlank("1"))
	t.Log(stringx.IsNotBlank("2 "))

	t.Log(stringx.IsNotBlank(str("")))
	t.Log(stringx.IsNotBlank(str(" ")))
	t.Log(stringx.IsNotBlank(str("	 ")))
	t.Log(stringx.IsNotBlank(str("1")))
	t.Log(stringx.IsNotBlank(str("2 ")))
}
