package keyboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSpace(t *testing.T) {
	assert.Equal(t, Spacebar, Parse(bytes{32, 0, 0}).Key)
}

func TestParseEscape(t *testing.T) {
	assert.Equal(t, Escape, Parse(bytes{27, 0, 0}).Key)
}

func TestParseHyphen(t *testing.T) {
	assert.Equal(t, Hyphen, Parse(bytes{45, 0, 0}).Key)
}

func TestParseAsterisk(t *testing.T) {
	assert.Equal(t, Asterisk, Parse(bytes{42, 0, 0}).Key)
}

func TestParsePlus(t *testing.T) {
	assert.Equal(t, Plus, Parse(bytes{43, 0, 0}).Key)
}

func TestParseSlash(t *testing.T) {
	assert.Equal(t, Slash, Parse(bytes{47, 0, 0}).Key)
}

func TestParseDot(t *testing.T) {
	assert.Equal(t, Dot, Parse(bytes{46, 0, 0}).Key)
}

func TestParseNotEscape(t *testing.T) {
	assert.NotEqual(t, Escape, Parse(bytes{27, 91, 65}).Key)
}

func TestParseNumberKeys(t *testing.T) {
	assert.Equal(t, 48, Parse(bytes{48, 0, 0}).Key)
	assert.Equal(t, 50, Parse(bytes{50, 0, 0}).Key)
	assert.Equal(t, 57, Parse(bytes{57, 0, 0}).Key)
}

func TestParseAlphaKeys(t *testing.T) {
	assert.Equal(t, 97, Parse(bytes{97, 0, 0}).Key)
	assert.Equal(t, 101, Parse(bytes{101, 0, 0}).Key)
	assert.Equal(t, 122, Parse(bytes{122, 0, 0}).Key)
}

func TestParseNotAlphaKeys(t *testing.T) {
	assert.NotEqual(t, 132, Parse(bytes{132, 0, 0}).Key)
}

func TestParseArrowKeys(t *testing.T) {
	assert.Equal(t, 65, Parse(bytes{27, 91, 65}).Key)
	assert.Equal(t, 66, Parse(bytes{27, 91, 66}).Key)
	assert.Equal(t, 67, Parse(bytes{27, 91, 67}).Key)
	assert.Equal(t, 68, Parse(bytes{27, 91, 68}).Key)
}

func TestParseNotArrowKeys(t *testing.T) {
	assert.NotEqual(t, Escape, Parse(bytes{27, 91, 65}).Key)
	assert.NotEqual(t, 70, Parse(bytes{27, 91, 70}).Key)
}
