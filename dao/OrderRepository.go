package dao

import (
	"database/sql"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	courierProto "stlab.itechart-group.com/go/food_delivery/courier_service/GRPC"
	"time"
)

type OrderPostgres struct {
	db *sql.DB
}

func NewDeliveryPostgres(db *sql.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

type Order struct {
	IdDeliveryService int       `json:"delivery_service_id,omitempty"`
	Id                int       `json:"id"`
	IdCourier         int       `json:"courier_id,omitempty"`
	DeliveryTime      time.Time `json:"delivery_time,omitempty"`
	CustomerAddress   string    `json:"customer_address,omitempty"`
	Status            string    `json:"status"`
	OrderDate         string    `json:"order_date"`
	RestaurantAddress string    `json:"restaurant_address"`
	Picked            bool      `json:"picked"`
}

type DetailedOrder struct {
	IdDeliveryService     int       `json:"delivery_service_id,omitempty"`
	IdOrder               int       `json:"id"`
	IdCourier             int       `json:"courier_id,omitempty"`
	DeliveryTime          time.Time `json:"delivery_time,omitempty"`
	CustomerAddress       string    `json:"customer_address,omitempty"`
	Status                string    `json:"status"`
	OrderDate             string    `json:"order_date,omitempty"`
	RestaurantAddress     string    `json:"restaurant_address,omitempty"`
	Picked                bool      `json:"picked"`
	CourierName           string    `json:"name"`
	CourierSurname        string    `json:"surname"`
	CourierPhoneNumber    string    `json:"phone_number"`
	OrderIdFromRestaurant int       `json:"id_from_restaurant"`
}

type AllInfoAboutOrder struct {
	IdDeliveryService     int       `json:"delivery_service_id,omitempty"`
	IdOrder               int       `json:"id"`
	IdCourier             int       `json:"courier_id,omitempty"`
	DeliveryTime          time.Time `json:"delivery_time,omitempty"`
	CustomerAddress       string    `json:"customer_address,omitempty"`
	Status                string    `json:"status"`
	OrderDate             string    `json:"order_date,omitempty"`
	RestaurantAddress     string    `json:"restaurant_address,omitempty"`
	RestaurantName        string    `json:"restaurant_name"`
	Picked                bool      `json:"picked"`
	CourierName           string    `json:"name"`
	CourierSurname        string    `json:"surname"`
	CourierPhoneNumber    string    `json:"phone_number"`
	OrderIdFromRestaurant int       `json:"id_from_restaurant"`
	CustomerName          string    `json:"customer_name"`
	CustomerPhone         string    `json:"customer_phone"`
	PaymentType           int       `json:"payment_type"`
}

func (r *OrderPostgres) GetActiveOrdersFromDB(id int) ([]Order, error) {
	var Orders []Order

	insertValue := `Select delivery_service_id,id,courier_id,delivery_time,customer_address,status,order_date,restaurant_address,picked from delivery where courier_id = $1 and status = 'ready to delivery'`
	get, err := r.db.Query(insertValue, id)
	if err != nil {
		log.Println("Error with getting list of orders: " + err.Error())
		return nil, err
	}

	for get.Next() {
		var order Order
		err = get.Scan(&order.IdDeliveryService, &order.Id, &order.IdCourier, &order.DeliveryTime, &order.CustomerAddress, &order.Status, &order.OrderDate, &order.RestaurantAddress, &order.Picked)
		Orders = append(Orders, order)
	}
	return Orders, nil
}

func (r *OrderPostgres) GetActiveOrderFromDB(id int) (Order, error) {
	var Ord Order

	insertValue := `Select delivery_service_id,id,courier_id,delivery_time,customer_address,status,order_date,restaurant_address,picked from delivery where id = $1`
	get, err := r.db.Query(insertValue, id)
	if err != nil {
		log.Println("Error with getting order by id: " + err.Error())
		return Order{}, err
	}

	for get.Next() {
		var order Order
		err = get.Scan(&order.IdDeliveryService, &order.Id, &order.IdCourier, &order.DeliveryTime, &order.CustomerAddress, &order.Status, &order.OrderDate, &order.RestaurantAddress, &order.Picked)
		Ord = order
	}
	return Ord, nil
}
func (r *OrderPostgres) GetOrderFromDB(id int) (Order, error) {
	var Ord Order

	insertValue := `Select delivery_service_id,id,courier_id,delivery_time,customer_address,status,order_date,restaurant_address,picked from delivery where id = $1`
	get, err := r.db.Query(insertValue, id)
	if err != nil {
		log.Println("Error with getting order by id: " + err.Error())
		return Order{}, err
	}

	for get.Next() {
		var order Order
		err = get.Scan(&order.IdDeliveryService, &order.Id, &order.IdCourier, &order.DeliveryTime, &order.CustomerAddress, &order.Status, &order.OrderDate, &order.RestaurantAddress, &order.Picked)
		Ord = order
	}
	return Ord, nil
}

func (r *OrderPostgres) ChangeOrderStatusInDB(text string, id uint16) (uint16, error) {

	UpdateValue := `UPDATE "delivery" SET "status" = $1 WHERE "id" = $2`
	_, err := r.db.Exec(UpdateValue, text, id)
	if err != nil {
		log.Println("Error with getting order by id: " + err.Error())
		return 0, fmt.Errorf("updateOrder: error while scanning for order:%w", err)
	}
	return id, nil
}

func (r *OrderPostgres) GetCourierCompletedOrdersWithPage_fromDB(limit, page, idCourier int) ([]DetailedOrder, int) {
	var Orders []DetailedOrder
	transaction, err := r.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer transaction.Commit()
	res, err := transaction.Query(fmt.Sprintf("SELECT delivery.order_date, delivery.courier_id,delivery.id,delivery.delivery_service_id,delivery.delivery_time,delivery.status,delivery.customer_address,delivery.restaurant_address,couriers.name,couriers.phone_number FROM delivery JOIN couriers ON couriers.id_courier=delivery.courier_id Where delivery.status='completed' and delivery.courier_id=%d LIMIT %d OFFSET %d", idCourier, limit, limit*(page-1)))
	if err != nil {
		log.Fatal(err)
	}
	for res.Next() {
		var order DetailedOrder
		err = res.Scan(&order.OrderDate, &order.IdCourier, &order.IdOrder, &order.IdDeliveryService, &order.DeliveryTime, &order.Status, &order.CustomerAddress, &order.RestaurantAddress, &order.CourierName, &order.CourierPhoneNumber)
		if err != nil {
			panic(err)
		}
		Orders = append(Orders, order)
	}

	var Ordersss []Order
	resl, err := transaction.Query(fmt.Sprintf("SELECT courier_id FROM delivery WHERE status='completed' and courier_id=%d ", idCourier))
	if err != nil {
		log.Println(err)
	}
	for resl.Next() {
		var order1 Order
		err = resl.Scan(&order1.IdCourier)
		if err != nil {
			panic(err)
		}

		Ordersss = append(Ordersss, order1)
	}
	return Orders, len(Ordersss)
}

func (r *OrderPostgres) GetAllOrdersOfCourierServiceWithPageFromDB(limit, page, idService int) ([]DetailedOrder, int) {
	var Orders []DetailedOrder
	transaction, err := r.db.Begin()
	if err != nil {
		log.Println(err)
	}
	defer transaction.Commit()
	res, err := transaction.Query("SELECT d.id_from_restaurant, d.order_date, d.courier_id,d.id,d.delivery_service_id,d.delivery_time,d.status,d.customer_address,d.restaurant_address,co.name, co.surname,co.phone_number FROM delivery AS d JOIN couriers AS co ON co.id_courier=d.courier_id Where d.delivery_service_id=$1 and status = 'ready to delivery' ORDER BY d.id LIMIT $2 OFFSET $3", idService, limit, limit*(page-1))
	if err != nil {
		log.Println(err)
	}
	for res.Next() {
		var order DetailedOrder
		err = res.Scan(&order.OrderIdFromRestaurant, &order.OrderDate, &order.IdCourier, &order.IdOrder, &order.IdDeliveryService, &order.DeliveryTime, &order.Status, &order.CustomerAddress, &order.RestaurantAddress, &order.CourierName, &order.CourierSurname, &order.CourierPhoneNumber)
		if err != nil {
			log.Println(err)
		}

		Orders = append(Orders, order)
	}

	var length int
	resl, err := transaction.Query("SELECT count(*) FROM delivery WHERE delivery_service_id=$1", idService)
	if err != nil {
		log.Println(err)
	}
	for resl.Next() {
		err = resl.Scan(&length)
		if err != nil {
			log.Println(err)
		}
	}
	return Orders, length
}

func (r *OrderPostgres) GetCourierCompletedOrdersByMouthWithPageFromDB(limit, page, idCourier, Month, Year int) ([]Order, int) {
	var Orders []Order
	transaction, err := r.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to db")

	defer transaction.Commit()
	res, err := transaction.Query("SELECT courier_id ,id ,delivery_service_id ,delivery_time ,order_date ,status ,customer_address, restaurant_address FROM delivery where status='completed' and courier_id=$1 and Extract(MONTH from order_date )=$2 and Extract(Year from order_date )=$3 LIMIT $4 OFFSET $5", idCourier, Month, Year, limit, limit*(page-1))
	if err != nil {
		panic(err)
	}
	for res.Next() {
		var order Order
		err = res.Scan(&order.IdCourier, &order.Id, &order.IdDeliveryService, &order.DeliveryTime, &order.OrderDate, &order.Status, &order.CustomerAddress, &order.RestaurantAddress)
		if err != nil {
			panic(err)
		}

		Orders = append(Orders, order)
	}
	var Ordersss []Order
	resl, err := transaction.Query(fmt.Sprintf("SELECT courier_id FROM delivery WHERE courier_id=%d and Extract(MONTH from order_date )=%d", idCourier, Month))
	if err != nil {
		panic(err)
	}
	for resl.Next() {
		var order Order
		err = resl.Scan(&order.IdCourier)
		if err != nil {
			panic(err)
		}

		Ordersss = append(Ordersss, order)
	}

	return Orders, len(Ordersss)
}

func (r *OrderPostgres) AssigningOrderToCourierInDB(order Order) error {
	log.Println("connected to db")
	s := "UPDATE delivery SET courier_id = $1 WHERE id = $2"
	log.Println(s)
	insert, err := r.db.Query(s, order.IdCourier, order.Id)
	defer insert.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *OrderPostgres) GetDetailedOrderByIdFromDB(Id int) (*AllInfoAboutOrder, error) {
	var order AllInfoAboutOrder
	transaction, err := r.db.Begin()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer transaction.Commit()
	res, err := transaction.Query(fmt.Sprintf("SELECT d.payment_type,d.customer_name,d.customer_phone,d.id_from_restaurant,d.id, d.order_date, d.courier_id,d.id,d.delivery_service_id,d.delivery_time,d.status,d.customer_address,d.restaurant_name,d.restaurant_address,co.name,co.surname,co.phone_number FROM delivery AS d JOIN couriers AS co ON co.id_courier=d.courier_id Where d.id=%d", Id))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for res.Next() {
		err = res.Scan(&order.PaymentType, &order.CustomerName, &order.CustomerPhone, &order.OrderIdFromRestaurant, &order.IdOrder, &order.OrderDate, &order.IdCourier, &order.IdOrder, &order.IdDeliveryService, &order.DeliveryTime, &order.Status, &order.CustomerAddress, &order.RestaurantName, &order.RestaurantAddress, &order.CourierName, &order.CourierSurname, &order.CourierPhoneNumber)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return &order, nil
}

func (r *OrderPostgres) CreateOrder(order *courierProto.OrderCourierServer) (*emptypb.Empty, error) {
	timestamp1 := time.Now()
	timestamp2 := time.Now().Add(45 * time.Minute)
	_, err := r.db.Exec("INSERT INTO delivery (delivery_service_id, customer_address, order_date, restaurant_address, delivery_time, restaurant_name, id_from_restaurant,customer_name,payment_type,customer_phone) VALUES ($1, $2, $3, $4, $5, $6, $7,$8,$9,$10)", order.CourierServiceID, order.ClientAddress, timestamp1, order.RestaurantAddress, timestamp2, order.RestaurantName, order.OrderID, order.ClientFullName, order.PaymentType, order.ClientPhoneNumber)
	if err != nil {
		log.Fatalf("CreateOrder:%s", err)
		return &emptypb.Empty{}, fmt.Errorf("CreateOrder:%w", err)
	}
	return &emptypb.Empty{}, nil
}

func (r *OrderPostgres) GetServices(in *emptypb.Empty) (*courierProto.ServicesResponse, error) {
	var Services courierProto.ServicesResponse

	res, err := r.db.Query("SELECT id, name, email, photo, description, phone_number, manager_id, status FROM delivery_service")
	if err != nil {
		log.Println("Error with getting list of orders: " + err.Error())
		return nil, err
	}
	for res.Next() {
		var service courierProto.DeliveryService
		err = res.Scan(&service.Id, &service.Name, &service.Email, &service.Photo, &service.Description, &service.Phone, &service.ManagerId, &service.Status)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		Services.Services = append(Services.Services, &service)
	}
	return &Services, nil
}

func (r *OrderPostgres) GetCompletedOrdersOfCourierServiceFromDB(limit, page, idService int) ([]Order, int) {
	var Orders []Order
	transaction, err := r.db.Begin()
	if err != nil {
		log.Println(err)
	}
	defer transaction.Commit()

	res, err := transaction.Query("SELECT order_date,courier_id,id,delivery_time,status,customer_address FROM delivery WHERE status='completed' and delivery_service_id=$1 ORDER BY id LIMIT $2 OFFSET $3",
		idService, limit, limit*(page-1))
	if err != nil {
		log.Println(err)
	}
	for res.Next() {
		var order Order
		err = res.Scan(&order.OrderDate, &order.IdCourier, &order.Id, &order.DeliveryTime, &order.Status, &order.CustomerAddress)
		if err != nil {
			log.Println(err)
		}

		Orders = append(Orders, order)
	}

	resl, err := transaction.Query("SELECT count(*) FROM delivery WHERE status='completed' and delivery_service_id=$1", idService)
	if err != nil {
		log.Println(err)
	}
	var length int
	for resl.Next() {

		err = resl.Scan(&length)
		if err != nil {
			log.Println(err)
		}
	}
	return Orders, length
}

func (r *OrderPostgres) GetCompletedOrdersOfCourierServiceByDateFromDB(limit, page, idService int) ([]Order, int) {
	var Orders []Order
	transaction, err := r.db.Begin()
	if err != nil {
		log.Println(err)
	}
	defer transaction.Commit()

	res, err := transaction.Query("SELECT delivery_service_id, order_date, courier_id,id,delivery_time,status,customer_address FROM delivery WHERE status='completed' and delivery_service_id=$1 ORDER BY order_date LIMIT $2 OFFSET $3",
		idService, limit, limit*(page-1))
	if err != nil {
		log.Println(err)
	}
	for res.Next() {
		var order Order
		err = res.Scan(&order.IdDeliveryService, &order.OrderDate, &order.IdCourier, &order.Id, &order.DeliveryTime, &order.Status, &order.CustomerAddress)
		if err != nil {
			log.Println(err)
		}

		Orders = append(Orders, order)
	}

	resl, err := transaction.Query("SELECT count(*) FROM delivery WHERE status='completed' and delivery_service_id=$1", idService)
	if err != nil {
		log.Println(err)
	}
	var length int
	for resl.Next() {
		err = resl.Scan(&length)
		if err != nil {
			log.Println(err)
		}
	}
	return Orders, length
}

func (r *OrderPostgres) GetCompletedOrdersOfCourierServiceByCourierIdFromDB(limit, page, idService int) ([]Order, int) {
	var Orders []Order
	transaction, err := r.db.Begin()
	if err != nil {
		log.Println(err)
	}
	defer transaction.Commit()

	res, err := transaction.Query("SELECT delivery_service_id, order_date, courier_id,id,delivery_time,status,customer_address FROM delivery WHERE status='completed' and delivery_service_id=$1 ORDER BY courier_id LIMIT $2 OFFSET $3",
		idService, limit, limit*(page-1))
	if err != nil {
		log.Println(err)
	}
	for res.Next() {
		var order Order
		err = res.Scan(&order.IdDeliveryService, &order.OrderDate, &order.IdCourier, &order.Id, &order.DeliveryTime, &order.Status, &order.CustomerAddress)
		if err != nil {
			log.Println(err)
		}

		Orders = append(Orders, order)
	}
	reslen, err := transaction.Query("SELECT count(*) FROM delivery WHERE status='completed' and delivery_service_id=$1", idService)
	if err != nil {
		log.Println(err)
	}
	var length int
	for reslen.Next() {
		err = reslen.Scan(&length)
		if err != nil {
			log.Println(err)
		}
	}
	return Orders, length
}

func (r *OrderPostgres) GetOrdersOfCourierServiceForManagerFromDB(limit, page, idService int) ([]DetailedOrder, int) {
	var Orders []DetailedOrder
	transaction, err := r.db.Begin()
	if err != nil {
		log.Println(err)
	}
	defer transaction.Commit()
	res, err := transaction.Query("SELECT d.order_date, d.courier_id,d.id,d.delivery_service_id,d.delivery_time,d.status,d.customer_address,d.restaurant_address,co.name, co.surname,co.phone_number FROM delivery AS d JOIN couriers AS co ON co.id_courier=d.courier_id Where d.delivery_service_id=$1 and status != 'completed' ORDER BY d.id LIMIT $2 OFFSET $3",
		idService, limit, limit*(page-1))
	if err != nil {
		log.Println(err)
	}
	for res.Next() {
		var order DetailedOrder
		err = res.Scan(&order.OrderDate, &order.IdCourier, &order.IdOrder, &order.IdDeliveryService, &order.DeliveryTime,
			&order.Status, &order.CustomerAddress, &order.RestaurantAddress, &order.CourierName, &order.CourierSurname,
			&order.CourierPhoneNumber)
		if err != nil {
			log.Println(err)
		}

		Orders = append(Orders, order)
	}

	var length int
	resl, err := transaction.Query("SELECT count(*) FROM delivery WHERE delivery_service_id=$1", idService)
	if err != nil {
		panic(err)
	}
	for resl.Next() {
		err = resl.Scan(&length)
		if err != nil {
			log.Println(err)
		}
	}
	return Orders, length
}
