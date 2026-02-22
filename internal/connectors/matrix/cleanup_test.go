package matrix

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestService_Cleanup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, fx := testService(ctrl)

	fx.matrixDB.EXPECT().Cleanup().Return(nil).MinTimes(1)

	err := service.Cleanup()
	require.NoError(t, err)
}
