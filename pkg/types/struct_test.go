package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeD(t *testing.T) {

	t1 := TypeD{Val: "[]int"}
	t2 := TypeD{Val: "int"}
	t3 := TypeD{Val: "map[string]bool"}
	t4 := TypeD{Val: "[]pkg.types.User{data=[]int}"}

	t.Run("tidy", func(t *testing.T) {
		assert.Equal(t, t1.Tidy(), "int")
		assert.Equal(t, t2.Tidy(), "int")
		assert.Equal(t, t3.Tidy(), "map[string]bool")
		assert.Equal(t, t4.Tidy(), "pkg.types.User")
	})

	t.Run("get_kind", func(t *testing.T) {
		assert.Equal(t, t1.GetKind(), "int")
		assert.Equal(t, t2.GetKind(), "int")
		assert.Equal(t, t3.GetKind(), "map[string]bool")
		assert.Equal(t, t4.GetKind(), "User")
	})

	t.Run("get_module", func(t *testing.T) {
		assert.Equal(t, t1.GetModule(), "")
		assert.Equal(t, t4.GetModule(), "pkg/types")
	})
}
