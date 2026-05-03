package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsOldSectionRef(t *testing.T) {
	assert.True(t, IsOldSectionRef("115BAC"))
	assert.True(t, IsOldSectionRef("80C"))
	assert.True(t, IsOldSectionRef("80D"))
	assert.True(t, IsOldSectionRef("80E"))
	assert.True(t, IsOldSectionRef("80G"))
	assert.True(t, IsOldSectionRef("80GG"))
	assert.True(t, IsOldSectionRef("80TTA"))
	assert.True(t, IsOldSectionRef("80TTB"))
	assert.True(t, IsOldSectionRef("80CCD"))
	assert.True(t, IsOldSectionRef("10(13A)"))
	assert.True(t, IsOldSectionRef("24(b)"))
}

func TestIsOldSectionRef_ITA2025Sections(t *testing.T) {
	assert.False(t, IsOldSectionRef("202"))
	assert.False(t, IsOldSectionRef("197"))
	assert.False(t, IsOldSectionRef("392"))
	assert.False(t, IsOldSectionRef("393(1)"))
}

func TestITA2025Equivalent(t *testing.T) {
	assert.NotEmpty(t, ITA2025Equivalent("115BAC"))
	assert.Contains(t, ITA2025Equivalent("115BAC"), "202")
	assert.Empty(t, ITA2025Equivalent("392"))
}
