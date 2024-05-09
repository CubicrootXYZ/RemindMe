package random_test

import (
	"fmt"
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/random"
	"github.com/stretchr/testify/assert"
)

func TestURLSaveString(t *testing.T) {
	for i := 1; i < 1000; i++ {
		t.Run(fmt.Sprintf("%d characters", i), func(t *testing.T) {
			str := random.URLSaveString(i)
			assert.Len(t, str, i)
		})
	}
}
