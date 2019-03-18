package ringio_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/antifuchs/o/ringio"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestPropReadWritesOverwrite(t *testing.T) {
	params := gopter.DefaultTestParameters()
	params.MinSize = 1
	params.MinSuccessfulTests = 1000
	properties := gopter.NewProperties(params)
	properties.Property("read a written slice in an overwriting ring buffer", prop.ForAll(
		func(cap uint, str string, times uint) *gopter.PropResult {
			input := []byte(str)
			b := ringio.New(cap, true)

			for i := uint(0); i < times; i++ {
				n, err := b.Write(input)
				if err != nil {
					res := gopter.NewPropResult(false, "writing")
					res.Error = err
					return res
				}
				if n != len(input) {
					return gopter.NewPropResult(false,
						fmt.Sprintf("wrong written length %d!=%d", n, len(input)))
				}

				output := make([]byte, n)
				n, err = b.Read(output)
				if err != nil {
					res := gopter.NewPropResult(false, "reading")
					res.Error = err
					return res
				}
				if cap >= uint(len(input)) {
					if n != len(input) {
						return gopter.NewPropResult(false,
							fmt.Sprintf("wrong read length %d!=%d", n, len(input)))
					}
					if !reflect.DeepEqual(output, input) {
						return gopter.NewPropResult(false,
							fmt.Sprintf("buffers are not equal: %#v %#v", input, output))
					}
				} else {
					if uint(n) != cap {
						return gopter.NewPropResult(false,
							fmt.Sprintf("wrong read length %d!=%d", n, cap))
					}
					writtenInput := input[len(input)-n:]
					if !reflect.DeepEqual(output[0:n], writtenInput) {
						return gopter.NewPropResult(false,
							fmt.Sprintf("buffers are not equal (%d written): %#v %#v", n, output, writtenInput))
					}
				}
			}

			return gopter.NewPropResult(true, "")
		},
		gen.UIntRange(1, 1024).WithLabel("buffer size"),
		gen.AnyString().SuchThat(func(x string) bool { return len(x) > 0 }).WithLabel("text to write"),
		gen.UIntRange(1, 10).WithLabel("time to read&write")),
	)
	properties.TestingRun(t)
}

func TestPropReadWritesBounded(t *testing.T) {
	params := gopter.DefaultTestParameters()
	params.MinSize = 1
	params.MinSuccessfulTests = 1000
	properties := gopter.NewProperties(params)
	properties.Property("read a written slice in a ring buffer that stops at the boundary", prop.ForAll(
		func(cap uint, str string, times uint) *gopter.PropResult {
			input := []byte(str)
			b := ringio.New(cap, false)

			for i := uint(0); i < times; i++ {
				tooLong := uint(len(input)) > cap
				n, err := b.Write(input)

				if tooLong {
					if n != 0 {
						return gopter.NewPropResult(false,
							fmt.Sprintf("should not have written, but: %d", n))
					}
				}

				output := make([]byte, len(input))
				n, err = b.Read(output)
				if !tooLong && err != nil {
					res := gopter.NewPropResult(false, "reading")
					res.Error = err
					return res
				}
				if tooLong && n != 0 {
					return gopter.NewPropResult(false,
						fmt.Sprintf("should not have read, but: %d", n))
				}

				if tooLong {
					return gopter.NewPropResult(true, "write too long")
				}
				if n != len(input) {
					return gopter.NewPropResult(false,
						fmt.Sprintf("wrong read length %d!=%d", n, len(input)))
				}
				if !reflect.DeepEqual(output, input) {
					return gopter.NewPropResult(false,
						fmt.Sprintf("buffers are not equal: %#v %#v", input, output))
				}
			}

			return gopter.NewPropResult(true, "")
		},
		gen.UIntRange(1, 1024).WithLabel("buffer size"),
		gen.AnyString().WithLabel("text to write"),
		gen.UIntRange(1, 10).WithLabel("time to read&write")),
	)
	properties.TestingRun(t)
}
