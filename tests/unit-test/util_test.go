package unit_test

import (
	"testing"

	asrt "github.com/stretchr/testify/assert"

	"general_ledger_golang/pkg/util"
)

func TestIncludes(t *testing.T) {
	assert := asrt.New(t)
	t.Run("Test_Includes", func(t *testing.T) {
		str1 := "ACG"
		str2 := "ABCD"
		matchers := []interface{}{"ABC", "C", "ACG", "CHF"}

		ok := util.Includes(str1, matchers)
		notOk := util.Includes(str2, matchers)

		assert.Equal(ok, true)
		assert.Equal(notOk, false)
	})
}
