package tests

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	authProto "stlab.itechart-group.com/go/food_delivery/courier_service/GRPCC"
	"stlab.itechart-group.com/go/food_delivery/courier_service/controller"
	"stlab.itechart-group.com/go/food_delivery/courier_service/dao"
	"stlab.itechart-group.com/go/food_delivery/courier_service/service"
	"stlab.itechart-group.com/go/food_delivery/courier_service/service/mocks"
	"testing"
	"time"
)

func TestHandler_GetCouriers(t *testing.T) {
	type mockBehaviorCheck func(s *mock_service.MockAllProjectApp, role string)
	type mockBehaviorParseToken func(s *mock_service.MockAllProjectApp, token string)
	type mockBehavior func(s *mock_service.MockAllProjectApp, courier dao.SmallInfo)

	var couriers []dao.SmallInfo
	cour := dao.SmallInfo{
		Id:          1,
		CourierName: "test",
		PhoneNumber: "1038812",
		Photo:       "my fav photo",
		Surname:     "Shorokhov",
		Deleted:     true,
	}
	couriers = append(couriers, cour)

	testTable := []struct {
		name                   string
		inputBody              string
		inputCourier           dao.SmallInfo
		inputRole              string
		inputToken             string
		mockBehaviorParseToken mockBehaviorParseToken
		mockBehavior           mockBehavior
		mockBehaviorCheck      mockBehaviorCheck
		expectedStatusCode     int
		expectedRequestBody    string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Test","id_courier":1,"courier_name":"test","phone_number":"1038812","photo":"my fav photo","surname":"Shorokhov","deleted":true}`,
			inputCourier: dao.SmallInfo{
				Id:          1,
				CourierName: "test",
				PhoneNumber: "1038812",
				Photo:       "my fav photo",
				Surname:     "Shorokhov",
			},
			inputRole:  "Courier manager",
			inputToken: "testToken",
			mockBehaviorParseToken: func(s *mock_service.MockAllProjectApp, token string) {
				s.EXPECT().ParseToken(token).Return(&authProto.UserRole{
					UserId:      1,
					Role:        "Courier manager",
					Permissions: "",
				}, nil)
			},
			mockBehaviorCheck: func(s *mock_service.MockAllProjectApp, role string) {
				s.EXPECT().CheckRole([]string{"Superadmin", "Courier manager"}, role).Return(nil)
			},
			mockBehavior: func(s *mock_service.MockAllProjectApp, courier dao.SmallInfo) {
				s.EXPECT().GetCouriers().Return(couriers, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id_courier":1,"courier_name":"test","phone_number":"1038812","photo":"my fav photo","surname":"Shorokhov","deleted":true}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			get := mock_service.NewMockAllProjectApp(c)
			testCase.mockBehavior(get, testCase.inputCourier)
			testCase.mockBehaviorParseToken(get, testCase.inputToken)
			testCase.mockBehaviorCheck(get, testCase.inputRole)

			services := &service.Service{AllProjectApp: get}
			handler := controller.NewHandler(services)
			r := handler.InitRoutesGin()

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/couriers/", bytes.NewBufferString(testCase.inputBody))
			req.Header.Set("Authorization", "Bearer testToken")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			//assert.Equal(t, testCase.expectedRequestBody,w.Body.String())
			assert.Contains(t, w.Body.String(), testCase.expectedRequestBody)

		})
	}
}

func TestHandler_GetOneCourier(t *testing.T) {
	type mockBehaviorCheck func(s *mock_service.MockAllProjectApp, role string)
	type mockBehaviorParseToken func(s *mock_service.MockAllProjectApp, token string)
	type mockBehavior func(s *mock_service.MockAllProjectApp, courier dao.Courier)
	cour := dao.Courier{
		Id:          1,
		UserId:      1,
		CourierName: "test",
		PhoneNumber: "1038812",
		Photo:       "my fav photo",
		Surname:     "Shorokhov",
		Deleted:     true,
	}

	testTable := []struct {
		name                   string
		inputBody              string
		inputCourier           dao.Courier
		inputRole              string
		inputToken             string
		mockBehaviorParseToken mockBehaviorParseToken
		mockBehavior           mockBehavior
		mockBehaviorCheck      mockBehaviorCheck
		expectedStatusCode     int
		expectedRequestBody    string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Test","id_courier":1,"user_id":1,"courier_name":"test","phone_number":"1038812","photo":"my fav photo","surname":"Shorokhov","deleted":true}`,
			inputCourier: dao.Courier{
				Id:                1,
				UserId:            1,
				CourierName:       "test",
				PhoneNumber:       "1038812",
				Photo:             "my fav photo",
				Surname:           "Shorokhov",
				DeliveryServiceId: 0,
			},
			inputRole:  "Courier manager",
			inputToken: "testToken",
			mockBehaviorParseToken: func(s *mock_service.MockAllProjectApp, token string) {
				s.EXPECT().ParseToken(token).Return(&authProto.UserRole{
					UserId:      1,
					Role:        "Courier manager",
					Permissions: "",
				}, nil)
			},
			mockBehaviorCheck: func(s *mock_service.MockAllProjectApp, role string) {
				s.EXPECT().CheckRole([]string{"Superadmin", "Courier", "Courier manager"}, role).Return(nil)
			},
			mockBehavior: func(s *mock_service.MockAllProjectApp, courier dao.Courier) {
				s.EXPECT().GetCourier(1).Return(cour, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id_courier":1,"user_id":1,"courier_name":"test","ready_to_go":false,"phone_number":"1038812","email":"","rating":0,"photo":"my fav photo","surname":"Shorokhov","number_of_failures":0,"deleted":true,"delivery_service_id":0}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			get := mock_service.NewMockAllProjectApp(c)
			testCase.mockBehavior(get, testCase.inputCourier)
			testCase.mockBehaviorParseToken(get, testCase.inputToken)
			testCase.mockBehaviorCheck(get, testCase.inputRole)

			services := &service.Service{AllProjectApp: get}
			handler := controller.NewHandler(services)
			r := handler.InitRoutesGin()

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/courier/1", bytes.NewBufferString(testCase.inputBody))
			req.Header.Set("Authorization", "Bearer testToken")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Contains(t, w.Body.String(), testCase.expectedRequestBody)

		})
	}
}

func TestHandler_GetCourierCompletedOrders(t *testing.T) {
	type mockBehaviorCheck func(s *mock_service.MockAllProjectApp, role string)
	type mockBehaviorParseToken func(s *mock_service.MockAllProjectApp, token string)
	type mockBehavior func(s *mock_service.MockAllProjectApp, order []dao.DetailedOrder)

	var orders []dao.DetailedOrder
	ord := dao.DetailedOrder{
		IdDeliveryService:     1,
		IdOrder:               1,
		IdCourier:             1,
		DeliveryTime:          time.Date(2020, time.May, 2, 2, 2, 2, 2, time.UTC),
		CustomerAddress:       "Some address",
		Status:                "ready to delivery",
		OrderDate:             "11.11.2022",
		CourierPhoneNumber:    "",
		CourierName:           "",
		CourierSurname:        "",
		OrderIdFromRestaurant: 0,
		Picked:                false,
	}
	orders = append(orders, ord)

	testTable := []struct {
		name                   string
		inputBody              string
		inputOrder             []dao.DetailedOrder
		inputRole              string
		inputToken             string
		mockBehaviorParseToken mockBehaviorParseToken
		mockBehavior           mockBehavior
		mockBehaviorCheck      mockBehaviorCheck
		expectedStatusCode     int
		expectedRequestBody    string
	}{
		{
			name: "OK",
			//inputBody: `{"name":"Test","delivery_service_id":1,"id":1,"courier_id":1,"delivery_time":"15:00","customer_address":"Some address","status":"ready to delivery","order_date":"11.11.2022"}`,
			inputBody: `{"name":"Test","courier_id":1}`,
			inputOrder: []dao.DetailedOrder{
				{
					IdCourier: 1,
				},
			},
			inputRole:  "Courier",
			inputToken: "testToken",
			mockBehaviorParseToken: func(s *mock_service.MockAllProjectApp, token string) {
				s.EXPECT().ParseToken(token).Return(&authProto.UserRole{
					UserId:      1,
					Role:        "Courier",
					Permissions: "",
				}, nil)
			},
			mockBehaviorCheck: func(s *mock_service.MockAllProjectApp, role string) {
				s.EXPECT().CheckRole([]string{"Superadmin", "Courier", "Courier manager"}, role).Return(nil)
			},
			mockBehavior: func(s *mock_service.MockAllProjectApp, order []dao.DetailedOrder) {
				s.EXPECT().GetCourierCompletedOrders(1, 1, 1).Return(orders, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"data":[{"delivery_service_id":1,"id":1,"courier_id":1,"delivery_time":"2020-05-02T02:02:02.000000002Z","customer_address":"Some address","status":"ready to delivery","order_date":"11.11.2022","picked":false,"name":"","surname":"","phone_number":"","id_from_restaurant":0}]}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			get := mock_service.NewMockAllProjectApp(c)
			testCase.mockBehavior(get, testCase.inputOrder)
			testCase.mockBehaviorParseToken(get, testCase.inputToken)
			testCase.mockBehaviorCheck(get, testCase.inputRole)

			services := &service.Service{AllProjectApp: get}
			handler := controller.NewHandler(services)
			r := handler.InitRoutesGin()

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/orders/completed?limit=1&page=1&idcourier=1", bytes.NewBufferString(testCase.inputBody))
			req.Header.Set("Authorization", "Bearer testToken")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
			assert.Contains(t, w.Body.String(), testCase.expectedRequestBody)

		})
	}
}

func TestHandler_GetAllOrdersOfCourierService(t *testing.T) {
	type mockBehaviorCheck func(s *mock_service.MockAllProjectApp, role string)
	type mockBehaviorParseToken func(s *mock_service.MockAllProjectApp, token string)
	type mockBehavior func(s *mock_service.MockAllProjectApp, order []dao.DetailedOrder)

	var orders []dao.DetailedOrder
	ord := dao.DetailedOrder{
		IdDeliveryService: 1,
		IdOrder:           1,
		IdCourier:         1,
		DeliveryTime:      time.Date(2020, time.May, 2, 2, 2, 2, 2, time.UTC),
		CustomerAddress:   "Some address",
		Status:            "ready to delivery",
		OrderDate:         "11.11.2022",
		RestaurantAddress: "",
		Picked:            false,
	}
	orders = append(orders, ord)

	testTable := []struct {
		name                   string
		inputBody              string
		inputOrder             []dao.DetailedOrder
		inputRole              string
		inputToken             string
		mockBehaviorParseToken mockBehaviorParseToken
		mockBehavior           mockBehavior
		mockBehaviorCheck      mockBehaviorCheck
		expectedStatusCode     int
		expectedRequestBody    string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Test","delivery_service_id":1,"id":1,"courier_id":1,"delivery_time":"2020-05-02T02:02:02.000000002Z","customer_address":"Some address","status":"ready to delivery","order_date":"11.11.2022","restaurant_address":"","picked":false}}`,
			inputOrder: []dao.DetailedOrder{
				{
					IdDeliveryService:     1,
					IdOrder:               1,
					IdCourier:             1,
					DeliveryTime:          time.Date(2020, time.May, 2, 2, 2, 2, 2, time.UTC),
					CustomerAddress:       "Some address",
					Status:                "ready to delivery",
					OrderDate:             "11.11.2022",
					RestaurantAddress:     "",
					OrderIdFromRestaurant: 0,
					Picked:                false,
				},
			},
			inputRole:  "Courier",
			inputToken: "testToken",
			mockBehaviorParseToken: func(s *mock_service.MockAllProjectApp, token string) {
				s.EXPECT().ParseToken(token).Return(&authProto.UserRole{
					UserId:      1,
					Role:        "Courier",
					Permissions: "",
				}, nil)
			},
			mockBehaviorCheck: func(s *mock_service.MockAllProjectApp, role string) {
				s.EXPECT().CheckRole([]string{"Superadmin", "Courier", "Courier manager"}, role).Return(nil)
			},
			mockBehavior: func(s *mock_service.MockAllProjectApp, order []dao.DetailedOrder) {
				s.EXPECT().GetAllOrdersOfCourierService(1, 1, 1).Return(orders, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"data":[{"delivery_service_id":1,"id":1,"courier_id":1,"delivery_time":"2020-05-02T02:02:02.000000002Z","customer_address":"Some address","status":"ready to delivery","order_date":"11.11.2022","picked":false,"name":"","surname":"","phone_number":"","id_from_restaurant":0}]}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			get := mock_service.NewMockAllProjectApp(c)
			testCase.mockBehavior(get, testCase.inputOrder)
			testCase.mockBehaviorParseToken(get, testCase.inputToken)
			testCase.mockBehaviorCheck(get, testCase.inputRole)

			services := &service.Service{AllProjectApp: get}
			handler := controller.NewHandler(services)
			r := handler.InitRoutesGin()

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/orders/?limit=1&page=1&iddeliveryservice=1", bytes.NewBufferString(testCase.inputBody))
			req.Header.Set("Authorization", "Bearer testToken")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
			assert.Contains(t, w.Body.String(), testCase.expectedRequestBody)

		})
	}
}

func TestHandler_GetCourierCompletedOrdersByMonth(t *testing.T) {
	type mockBehaviorCheck func(s *mock_service.MockAllProjectApp, role string)
	type mockBehaviorParseToken func(s *mock_service.MockAllProjectApp, token string)
	type mockBehavior func(s *mock_service.MockAllProjectApp, order []dao.Order)

	var orders []dao.Order
	ord := dao.Order{
		IdDeliveryService: 1,
		Id:                1,
		IdCourier:         1,
		DeliveryTime:      time.Date(2020, time.May, 2, 2, 2, 2, 2, time.UTC),
		CustomerAddress:   "Some address",
		Status:            "ready to delivery",
		OrderDate:         "11.11.2022",
		RestaurantAddress: "",
		Picked:            false,
	}
	orders = append(orders, ord)

	testTable := []struct {
		name                   string
		inputBody              string
		inputOrder             []dao.Order
		inputRole              string
		inputToken             string
		mockBehaviorParseToken mockBehaviorParseToken
		mockBehavior           mockBehavior
		mockBehaviorCheck      mockBehaviorCheck
		expectedStatusCode     int
		expectedRequestBody    string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Test","delivery_service_id":1,"id":1,"courier_id":1,"delivery_time":"2020-05-02T02:02:02.000000002Z","customer_address":"Some address","status":"ready to delivery","order_date":"11.11.2022","restaurant_address":"","picked":false}`,
			inputOrder: []dao.Order{
				{
					IdDeliveryService: 1,
					Id:                1,
					IdCourier:         1,
					DeliveryTime:      time.Date(2020, time.May, 2, 2, 2, 2, 2, time.UTC),
					CustomerAddress:   "Some address",
					Status:            "ready to delivery",
					OrderDate:         "11.11.2022",
					RestaurantAddress: "",
					Picked:            false,
				},
			},
			inputRole:  "Courier",
			inputToken: "testToken",
			mockBehaviorParseToken: func(s *mock_service.MockAllProjectApp, token string) {
				s.EXPECT().ParseToken(token).Return(&authProto.UserRole{
					UserId:      1,
					Role:        "Courier",
					Permissions: "",
				}, nil)
			},
			mockBehaviorCheck: func(s *mock_service.MockAllProjectApp, role string) {
				s.EXPECT().CheckRole([]string{"Superadmin", "Courier", "Courier manager"}, role).Return(nil)
			},
			mockBehavior: func(s *mock_service.MockAllProjectApp, order []dao.Order) {
				s.EXPECT().GetCourierCompletedOrdersByMonth(1, 1, 1, 11, 2022).Return(orders, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"data":[{"delivery_service_id":1,"id":1,"courier_id":1,"delivery_time":"2020-05-02T02:02:02.000000002Z","customer_address":"Some address","status":"ready to delivery","order_date":"11.11.2022","restaurant_address":"","picked":false}]}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			get := mock_service.NewMockAllProjectApp(c)
			testCase.mockBehavior(get, testCase.inputOrder)
			testCase.mockBehaviorParseToken(get, testCase.inputToken)
			testCase.mockBehaviorCheck(get, testCase.inputRole)

			services := &service.Service{AllProjectApp: get}
			handler := controller.NewHandler(services)
			r := handler.InitRoutesGin()

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/orders/bymonth?limit=1&page=1&idcourier=1&month=11&year=2022", bytes.NewBufferString(testCase.inputBody))
			req.Header.Set("Authorization", "Bearer testToken")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
			assert.Contains(t, w.Body.String(), testCase.expectedRequestBody)
		})
	}
}
