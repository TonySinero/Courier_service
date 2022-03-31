package dao

import (
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestRepository_GetCourierCompletedOrdersWithPage_fromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db)

	testTable := []struct {
		name          string
		mock          func(courier_id, limit, page int)
		courier_id    int
		limit         int
		page          int
		expectedOrder []DetailedOrder
	}{
		{
			name: "OK",
			mock: func(courier_id, limit, page int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"order_date", "courier_id", "id", "delivery_service_id", "delivery_time", "status", "customer_address", "restaurant_address", "name", "phone_number"}).
					AddRow("2022-02-02", 1, 1, 1, time.Date(2020, time.May, 2, 2, 2, 2, 2, time.UTC), "completed", "address", "address", "name", "1234567")

				mock.ExpectQuery(`SELECT delivery.order_date, delivery.courier_id,delivery.id,delivery.delivery_service_id,delivery.delivery_time,delivery.status,delivery.customer_address,delivery.restaurant_address,couriers.name,couriers.phone_number FROM delivery JOIN couriers ON`).
					WillReturnRows(rows)

				rows2 := sqlmock.NewRows([]string{"courier_id"}).
					AddRow(1)
				mock.ExpectQuery(`SELECT courier_id FROM delivery WHERE status='completed' (.+)`).
					WillReturnRows(rows2)

				mock.ExpectCommit()
			},
			courier_id: 1,
			limit:      1,
			page:       1,
			expectedOrder: []DetailedOrder{
				{
					IdDeliveryService:  1,
					IdOrder:            1,
					IdCourier:          1,
					DeliveryTime:       time.Date(2020, time.May, 2, 2, 2, 2, 2, time.UTC),
					CustomerAddress:    "address",
					Status:             "completed",
					CourierName:        "name",
					CourierPhoneNumber: "1234567",
					RestaurantAddress:  "address",
					OrderDate:          "2022-02-02",
				},
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.courier_id, tt.limit, tt.page)
			got, _ := r.GetCourierCompletedOrdersWithPage_fromDB(tt.courier_id, tt.limit, tt.page)

			assert.Equal(t, tt.expectedOrder, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
func TestRepository_GetCourierCompletedOrdersByMouthWithPage_fromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db)

	testTable := []struct {
		name          string
		mock          func(courier_id, limit, page int)
		courier_id    int
		limit         int
		page          int
		month         int
		year          int
		expectedOrder []Order
	}{
		{
			name: "OK",
			mock: func(courier_id, limit, page int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"courier_id", "id", "delivery_service_id", "delivery_time", "order_date", "status", "customer_address", "restaurant_address"}).
					AddRow(1, 1, 1, time.Date(2020, time.May, 2, 2, 2, 2, 2, time.UTC), "2022-02-02", "completed", "address", "restaurant_address")

				mock.ExpectQuery(`SELECT courier_id ,id ,delivery_service_id ,delivery_time ,order_date ,status ,customer_address, restaurant_address FROM delivery where (.+)`).
					WillReturnRows(rows)

				rows2 := sqlmock.NewRows([]string{"courier_id"}).
					AddRow(1)
				mock.ExpectQuery(`SELECT courier_id FROM delivery WHERE (.+)`).
					WillReturnRows(rows2)

				mock.ExpectCommit()
			},
			courier_id: 1,
			limit:      1,
			page:       1,
			month:      1,
			year:       2022,
			expectedOrder: []Order{
				{
					IdDeliveryService: 1,
					Id:                1,
					IdCourier:         1,
					DeliveryTime:      time.Date(2020, time.May, 2, 2, 2, 2, 2, time.UTC),
					OrderDate:         "2022-02-02",
					CustomerAddress:   "address",
					RestaurantAddress: "restaurant_address",
					Status:            "completed",
				},
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.courier_id, tt.limit, tt.page)
			got, _ := r.GetCourierCompletedOrdersByMouthWithPageFromDB(tt.courier_id, tt.limit, tt.page, tt.month, tt.year)

			assert.Equal(t, tt.expectedOrder, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
