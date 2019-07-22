package snowflake_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viostream/terraform-provider-snowflake/pkg/snowflake"
)

func TestEscapeString(t *testing.T) {
	a := assert.New(t)

	a.Equal(`\'`, snowflake.EscapeString(`'`))
	a.Equal(`\\\'`, snowflake.EscapeString(`\'`))
}
