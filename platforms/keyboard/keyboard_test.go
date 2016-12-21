package keyboard

import (
	"testing"

	"gobot.io/x/gobot/gobottest"
)

func TestParseSpace(t *testing.T) {
	gobottest.Assert(t, Parse(bytes{32, 0, 0}).Key, Spacebar)
}

func TestParseEscape(t *testing.T) {
	gobottest.Assert(t, Parse(bytes{27, 0, 0}).Key, Escape)
}

func TestParseNotEscape(t *testing.T) {
	gobottest.Refute(t, Parse(bytes{27, 91, 65}).Key, Escape)
}

func TestParseNumberKeys(t *testing.T) {
	gobottest.Assert(t, Parse(bytes{48, 0, 0}).Key, 48)
	gobottest.Assert(t, Parse(bytes{50, 0, 0}).Key, 50)
	gobottest.Assert(t, Parse(bytes{57, 0, 0}).Key, 57)
}

func TestParseAlphaKeys(t *testing.T) {
	gobottest.Assert(t, Parse(bytes{97, 0, 0}).Key, 97)
	gobottest.Assert(t, Parse(bytes{101, 0, 0}).Key, 101)
	gobottest.Assert(t, Parse(bytes{122, 0, 0}).Key, 122)
}

func TestParseNotAlphaKeys(t *testing.T) {
	gobottest.Refute(t, Parse(bytes{132, 0, 0}).Key, 132)
}

func TestParseArrowKeys(t *testing.T) {
	gobottest.Assert(t, Parse(bytes{27, 91, 65}).Key, 65)
	gobottest.Assert(t, Parse(bytes{27, 91, 66}).Key, 66)
	gobottest.Assert(t, Parse(bytes{27, 91, 67}).Key, 67)
	gobottest.Assert(t, Parse(bytes{27, 91, 68}).Key, 68)
}

func TestParseNotArrowKeys(t *testing.T) {
	gobottest.Refute(t, Parse(bytes{27, 91, 65}).Key, Escape)
	gobottest.Refute(t, Parse(bytes{27, 91, 70}).Key, 70)
}
