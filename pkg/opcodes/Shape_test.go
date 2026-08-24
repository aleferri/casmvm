package opcodes

import "testing"

func TestReshapeSigned(t *testing.T) {
	cases := []struct {
		shape    Shape
		in       int64
		expected int64
	}{
		{ByteShape, 200, -56},
		{ByteShape, 127, 127},
		{ByteShape, -1, -1},
		{ShortShape, -1, -1},
		{ShortShape, 40000, -25536},
		{IntShape, -128, -128},
		{IntShape, 1 << 33, 0},
		{LongShape, -5, -5},
	}
	for _, c := range cases {
		got := c.shape.Reshape(c.in)
		if got != c.expected {
			t.Errorf("%s.Reshape(%d): expected %d, got %d", c.shape.Name(), c.in, c.expected, got)
		}
	}
}

func TestReshapeUnsigned(t *testing.T) {
	cases := []struct {
		shape    Shape
		in       int64
		expected int64
	}{
		{UByteShape, 200, 200},
		{UByteShape, -1, 255},
		{UShortShape, 70000, 4464},
		{UShortShape, -1, 65535},
		{UIntShape, -1, 4294967295},
		{ULongShape, -1, -1},
	}
	for _, c := range cases {
		got := c.shape.Reshape(c.in)
		if got != c.expected {
			t.Errorf("%s.Reshape(%d): expected %d, got %d", c.shape.Name(), c.in, c.expected, got)
		}
	}
}
