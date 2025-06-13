package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/shekshuev/gophkeeper/internal/mocks"
	"github.com/shekshuev/gophkeeper/internal/models"
	"github.com/shekshuev/gophkeeper/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	users := mocks.NewMockUserService(ctrl)
	accessTokenSecret := "test"
	os.Setenv("ACCESS_TOKEN_SECRET", accessTokenSecret)
	os.Setenv("ACCESS_TOKEN_EXPIRES", "1h")
	cfg := config.GetConfig()
	accessToken, err := utils.CreateToken(
		cfg.AccessTokenSecret,
		"1",
		cfg.AccessTokenExpires,
	)
	assert.NoError(t, err, "error creating token")
	handler := NewHandler(users, nil, &cfg)
	httpSrv := httptest.NewServer(handler.Router)
	defer httpSrv.Close()

	testCases := []struct {
		name          string
		expectedCode  int
		userID        string
		responseDTO   *models.ReadUserDTO
		serviceError  error
		serviceCalled bool
	}{
		{
			name:          "Success fetching user",
			expectedCode:  http.StatusOK,
			userID:        "1",
			responseDTO:   &models.ReadUserDTO{ID: 1, UserName: "test_user"},
			serviceError:  nil,
			serviceCalled: true,
		},
		{
			name:          "Invalid user ID",
			expectedCode:  http.StatusNotFound,
			userID:        "abc",
			responseDTO:   nil,
			serviceError:  nil,
			serviceCalled: false,
		},
		{
			name:          "User not found",
			expectedCode:  http.StatusNotFound,
			userID:        "2",
			responseDTO:   nil,
			serviceError:  assert.AnError,
			serviceCalled: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.serviceCalled {
				users.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(tc.responseDTO, tc.serviceError)
			}
			req := resty.New().R()
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req.Method = http.MethodGet
			req.URL = httpSrv.URL + "/v1.0/users/" + tc.userID
			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}
