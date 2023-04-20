package utils

import (
	"math"
	"testing"

	"github.com/nbio/st"
)

func TestMinInt(t *testing.T) {
	//TODO I know, I should mock math.Min really
	var top int = 25
	var bottom int = 12
	want := int(math.Min(float64(top), float64(bottom)))
	got := Min(top, bottom)
	st.Expect(t, want, got)
	got = Min(bottom, top)
	st.Expect(t, want, got)
}

func TestMinFloat(t *testing.T) {
	//TODO I know, I should mock math.Min really
	var top float32 = 25
	var bottom float32 = 12
	want := float32(math.Min(float64(top), float64(bottom)))
	got := Min(top, bottom)
	st.Expect(t, want, got)
	got = Min(bottom, top)
	st.Expect(t, want, got)
}
