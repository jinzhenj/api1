package api1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	parser := Parser{}
	var err error

	t1 := `
    group t1
		
		scalar T1
		
		enum T1 {
			O1
			O2
		}
	`
	t01 := `
	  group t01

		struct T01 {
			f1: int
		}

		interface T01 {
			f1(): int
		}
	`
	t2 := `
	  group t2
		
		enum T2 {
			O1
			O1
			O2
			O2
		}
	`
	t3 := `
	  group t3

		struct T3 {
			F1: int
			F1: string
		}
	`
	t4 := `
	  group t4

		interface T4 {
			f1(): int
			f1(): string
		}
	`
	t5 := `
	  group t5

		interface T5 {
			f1(param1: int, param1: string): int
		}
	`
	t6 := `
	  group t6

		struct T6 {
			F1: unknown
		}
	`
	t7 := `
	  group t7

		struct T7 {
			F1: [[unknown]]
		}
	`
	t8 := `
	  group t8

		interface T8 {
			f1(): unknown
		}
	`
	t9 := `
	  group t9

		interface T9 {
			f1(param1: [unknown]): int
		}
	`
	t10 := `
		# test void is no longger valid
	  group t10

		interface T10 {
			f1(): void
		}
	`
	t11 := `
    # test enum value type not mixed
		group t11

		enum E11 {
			O1 = 1
			O2 = 2
		}
	`

	t12 := `
    # test enum value type not mixed
		group t12

		enum E12 {
			O1 = "str1"
			O2 = "str2"
		}
	`

	t13 := `
    # test enum value type not mixed
		group t13

		enum E13 {
			O1 = 1
			O2 = 2
		}
	`

	t14 := `
    # test enum value type not mixed
		group t14

		enum E14 {
			O1 = 1
			O2 = "str2"
		}
	`

	t15 := `
    # test enum value type not mixed
		group t15

		enum E15 {
			O1 = 1
			O2
		}
	`

	testcases := []string{t1, t01, t2, t3, t4, t5, t6, t7, t8, t9, t10, t14, t15}
	for _, testcase := range testcases {
		_, err = parser.Parse(testcase)
		t.Log(err)
		assert.Error(t, err)
	}

	testcases2 := []string{t11, t12, t13}
	for _, testcase := range testcases2 {
		_, err = parser.Parse(testcase)
		assert.NoError(t, err)
	}
}
