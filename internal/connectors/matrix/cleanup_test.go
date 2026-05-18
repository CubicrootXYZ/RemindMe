package matrix

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestService_Cleanup(t *testing.T) {
	service, fx := testService(t)

	fx.matrixDB.EXPECT().Cleanup().Return(nil).Once()

	err := service.Cleanup()
	require.NoError(t, err)
}
