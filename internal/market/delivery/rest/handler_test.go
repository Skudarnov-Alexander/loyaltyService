package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/market"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/market/service/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestHandler_PostOrder(t *testing.T) {
	type mockBehaviorCheckOrder func(s *mock.MockMarketService, userID, orderID string)
	type mockBehaviorSaveOrder func(s *mock.MockMarketService, userID, orderID string)

	type ResponseTest struct {
		Message string `json:"message"`
	}

	type want struct {
		contentType string
		statusCode  int
		response    ResponseTest
		error       bool
		errorMsg    echo.HTTPError
	}

	type input struct {
		orderID     string
		contentType string
		userID      string
	}

	testTable := []struct {
		name           string
		input          input
		mockCheckOrder mockBehaviorCheckOrder
		mockSaveOrder  mockBehaviorSaveOrder
		want           want
	}{
		{
			name: "Add existed order (User) - Status:200 OK",
			input: input{
				orderID:     "567867585",
				contentType: echo.MIMETextPlain,
				userID:      "testUUID",
			},
			mockCheckOrder: func(s *mock.MockMarketService, userID, orderID string) {
				s.EXPECT().CheckOrder(gomock.Any(), userID, orderID).
					Return(true, nil)
			},
			mockSaveOrder: func(s *mock.MockMarketService, userID, orderID string) {},
			want: want{
				contentType: echo.MIMEApplicationJSON,
				statusCode:  http.StatusOK,
				response:    ResponseTest{Message: "order is loaded yet"},
				error:       false,
			},
		},
		{
			name: "Add new order - Status:202 Accepted",
			input: input{
				orderID:     "567867585",
				contentType: echo.MIMETextPlain,
				userID:      "testUUID",
			},
			mockCheckOrder: func(s *mock.MockMarketService, userID, orderID string) {
				s.EXPECT().CheckOrder(gomock.Any(), userID, orderID).
					Return(false, nil)
			},
			mockSaveOrder: func(s *mock.MockMarketService, userID, orderID string) {
				s.EXPECT().SaveOrder(gomock.Any(), userID, orderID).
					Return(nil)
			},
			want: want{
				contentType: echo.MIMEApplicationJSON,
				statusCode:  http.StatusAccepted,
				response:    ResponseTest{Message: "new order is loaded"},
				error:       false,
			},
		},

		{
			name: "Add existed order (another User) - Status:407 Conflict",
			input: input{
				orderID:     "567867585",
				contentType: echo.MIMETextPlain,
				userID:      "testUUID",
			},
			mockCheckOrder: func(s *mock.MockMarketService, userID, orderID string) {
				s.EXPECT().CheckOrder(gomock.Any(), userID, orderID).
					Return(false, errors.New("blalal"))
			},
			mockSaveOrder: func(s *mock.MockMarketService, userID, orderID string) {
				s.EXPECT().SaveOrder(gomock.Any(), userID, orderID).
					Return(errors.New("blalal"))
			},
			want: want{
				contentType: echo.MIMEApplicationJSON,
				statusCode:  http.StatusConflict,
				error:       true,
				errorMsg: echo.HTTPError{
					Code:     409,
					Message:  market.ErrOrderIsExist,
					Internal: nil,
				},
			},
		},

		/*
			{
				name:        "Get counter metric value - Status:200 OK",
				inputBody:   `{"id":"testCounter","type":"counter"}`,
				inputMetric: model.Metric{ID: "testCounter", MType: "counter"},
				mockBehavior: func(s *mocks.MockService, metric model.Metric) {
					s.EXPECT().FetchMetric(gomock.Any(), metric).
						Return(&model.Metric{
							ID:    "testCounter",
							MType: "counter",
							Delta: &dMock[0],
						}, nil)
				},
				want: want{
					contentType: "application/json",
					statusCode:  http.StatusOK,
					body:        `{"id":"testCounter","type":"counter","delta":7676}`,
				},
			},
			{
				name:        "Service failure - Status:404 Not Found",
				inputBody:   `{"id":"testGauge","type":"gauge"}`,
				inputMetric: model.Metric{ID: "testGauge", MType: "gauge"},
				mockBehavior: func(s *mocks.MockService, metric model.Metric) {
					s.EXPECT().FetchMetric(gomock.Any(), metric).
						Return(nil, errors.New("service failure"))
				},
				want: want{
					contentType: "application/json",
					statusCode:  http.StatusNotFound,
					body:        `{404 Metric is not found}`,
				},
			},
			{
				name:        "Service returns nil struct - Status:404 Not Found",
				inputBody:   `{"id":"testGauge","type":"gauge"}`,
				inputMetric: model.Metric{ID: "testGauge", MType: "gauge"},
				mockBehavior: func(s *mocks.MockService, metric model.Metric) {
					s.EXPECT().FetchMetric(gomock.Any(), metric).
						Return(nil, nil)
				},
				want: want{
					contentType: "application/json",
					statusCode:  http.StatusNotFound,
					body:        `{404 Metric is not found}`,
				},
			},
			{
				name:         "Empty fields in Metric - Status:400 BadRequest",
				inputBody:    `{"id":"testCounter"}`,
				mockBehavior: func(s *mocks.MockService, metric model.Metric) {},
				want: want{
					contentType: "application/json",
					statusCode:  http.StatusBadRequest,
					body:        `{400 BadRequest}`,
				},
			},
			{
				name:         "Empty Body - Status:400 BadRequest",
				mockBehavior: func(s *mocks.MockService, metric model.Metric) {},
				want: want{
					contentType: "application/json",
					statusCode:  http.StatusBadRequest,
					body:        `{400 BadRequest}`,
				},
			},
			{
				name:         "not JSON Body - Status:400 BadRequest",
				inputBody:    "testNoJSON",
				mockBehavior: func(s *mocks.MockService, metric model.Metric) {},
				want: want{
					contentType: "application/json",
					statusCode:  http.StatusBadRequest,
					body:        `{400 BadRequest}`,
				},
			},
			{
				name:         "Get metric with incorrect type - Status:400 BadRequest",
				inputBody:    `{"id":"testGauge","type":"falseType"}`,
				mockBehavior: func(s *mocks.MockService, metric model.Metric) {},
				want: want{
					contentType: "application/json",
					statusCode:  http.StatusBadRequest,
					body:        `{400 BadRequest}`,
				},
			},
			{
				name:         "Get metric with extra fields - Status:400 BadRequest",
				inputBody:    `{"id":"testGauge","type":"gauge","hash":"56ab45bc673"}`,
				mockBehavior: func(s *mocks.MockService, metric model.Metric) {},
				want: want{
					contentType: "application/json",
					statusCode:  http.StatusBadRequest,
					body:        `{400 BadRequest}`,
				},
			},

			{
				name:         "Nil Body - Status:400 BadRequest",
				mockBehavior: func(s *mocks.MockService, metric model.Metric) {},
				want: want{
					contentType: "text/plain; charset=utf-8",
					statusCode:  http.StatusBadRequest,
					body:        "metric type is incorrect\n", //как убрать перенос каретки?
				},
			},*/

	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			//---Init Deps
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock.NewMockMarketService(ctrl)
			tt.mockCheckOrder(mockService, tt.input.userID, tt.input.orderID)
			tt.mockSaveOrder(mockService, tt.input.userID, tt.input.orderID)

			h := New(mockService)

			//---Test Router
			e := echo.New()
			e.POST("/api/user/orders", h.PostOrder)

			//---Test Request
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString(tt.input.orderID))

			req.Header.Set(echo.HeaderContentType, echo.MIMETextPlain)
			c := e.NewContext(req, rec)
			c.Set("uuid", tt.input.userID)

			//---Assets
			if !tt.want.error {
				h.PostOrder(c)
				assert.Equalf(t, tt.want.statusCode, rec.Code,
					"Incorrect status code: want = %d, have = %d", tt.want.statusCode, rec.Code)

				ct := rec.Header().Get(echo.HeaderContentType)
				assert.Equalf(t, tt.want.contentType, ct,
					"Content-Type doesn't contain correct value: got %d want %d", ct, tt.want.contentType)

				body := ResponseTest{}
				if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body)) {
					assert.Equal(t, tt.want.response, body)
				}
			}

			if tt.want.error {
				assert.ErrorIs(t, h.PostOrder(c), market.ErrOrderIsExist)
				//fmt.Printf("hhh %s\n", rec.Body.Bytes())

				assert.Equalf(t, tt.want.statusCode, rec.Code,
					"Incorrect status code: want = %d, have = %d", tt.want.statusCode, rec.Code)

				ct := rec.Header().Get(echo.HeaderContentType)
				assert.Equalf(t, tt.want.contentType, ct,
					"Content-Type doesn't contain correct value: got %d want %d", ct, tt.want.contentType)

			}
		})
	}

}
