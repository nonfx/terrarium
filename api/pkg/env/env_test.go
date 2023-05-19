package env

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvString(t *testing.T) {
	os.Setenv(PREFIX+"TEST_VAR", "test_value")
	defer os.Unsetenv(PREFIX + "TEST_VAR")

	val := GetEnvString("TEST_VAR", "default_value")
	assert.Equal(t, "test_value", val)

	val = GetEnvString("NON_EXISTENT_VAR", "default_value")
	assert.Equal(t, "default_value", val)
}

func TestGetEnvInt(t *testing.T) {
	os.Setenv(PREFIX+"TEST_VAR", "5")
	defer os.Unsetenv(PREFIX + "TEST_VAR")

	val := GetEnvInt("TEST_VAR", 10)
	assert.Equal(t, 5, val)

	val = GetEnvInt("NON_EXISTENT_VAR", 10)
	assert.Equal(t, 10, val)

	os.Setenv(PREFIX+"INVALID_VAR", "invalid")
	defer os.Unsetenv(PREFIX + "INVALID_VAR")

	val = GetEnvInt("INVALID_VAR", 10)
	assert.Equal(t, 10, val)
}

func TestGetEnvBool(t *testing.T) {
	os.Setenv(PREFIX+"TEST_VAR", "true")
	defer os.Unsetenv(PREFIX + "TEST_VAR")

	val := GetEnvBool("TEST_VAR", false)
	assert.True(t, val)

	val = GetEnvBool("NON_EXISTENT_VAR", false)
	assert.False(t, val)

	os.Setenv(PREFIX+"INVALID_VAR", "invalid")
	defer os.Unsetenv(PREFIX + "INVALID_VAR")

	val = GetEnvBool("INVALID_VAR", false)
	assert.False(t, val)
}

func TestGetEnvInt64(t *testing.T) {
	os.Setenv(PREFIX+"TEST_VAR", strconv.FormatInt(500, 10))
	defer os.Unsetenv(PREFIX + "TEST_VAR")

	val := GetEnvInt64("TEST_VAR", 1000)
	assert.Equal(t, int64(500), val)

	val = GetEnvInt64("NON_EXISTENT_VAR", 1000)
	assert.Equal(t, int64(1000), val)

	os.Setenv(PREFIX+"INVALID_VAR", "invalid")
	defer os.Unsetenv(PREFIX + "INVALID_VAR")

	val = GetEnvInt64("INVALID_VAR", 1000)
	assert.Equal(t, int64(1000), val)
}
