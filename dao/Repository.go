package dao

import (
	"database/sql"
	"google.golang.org/protobuf/types/known/emptypb"
	courierProto "stlab.itechart-group.com/go/food_delivery/courier_service/GRPC"
)

type Repository struct {
	OrderRep
	CourierRep
	DeliveryServiceRep
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		NewDeliveryPostgres(db),
		NewCourierPostgres(db),
		NewDeliveryServicePostgres(db),
	}
}

type OrderRep interface {
	GetActiveOrdersFromDB(id int) ([]Order, error)
	GetActiveOrderFromDB(id int) (Order, error)
	ChangeOrderStatusInDB(text string, id uint16) (uint16, error)
	GetOrderFromDB(id int) (Order, error)
	GetCourierCompletedOrdersWithPage_fromDB(limit, page, idCourier int) ([]DetailedOrder, int)
	GetAllOrdersOfCourierServiceWithPageFromDB(limit, page, idService int) ([]DetailedOrder, int)
	GetCourierCompletedOrdersByMouthWithPageFromDB(limit, page, idCourier, Month, Year int) ([]Order, int)
	AssigningOrderToCourierInDB(order Order) error
	GetDetailedOrderByIdFromDB(Id int) (*AllInfoAboutOrder, error)
	CreateOrder(order *courierProto.OrderCourierServer) (*emptypb.Empty, error)
	GetServices(in *emptypb.Empty) (*courierProto.ServicesResponse, error)
	GetCompletedOrdersOfCourierServiceFromDB(limit, page, idService int) ([]Order, int)
	GetCompletedOrdersOfCourierServiceByDateFromDB(limit, page, idService int) ([]Order, int)
	GetCompletedOrdersOfCourierServiceByCourierIdFromDB(limit, page, idService int) ([]Order, int)
	GetOrdersOfCourierServiceForManagerFromDB(limit, page, idService int) ([]DetailedOrder, int)
}

type CourierRep interface {
	SaveCourierInDB(Courier *Courier) error
	GetCouriersFromDB() ([]SmallInfo, error)
	GetCourierFromDB(id int) (Courier, error)
	UpdateCourierInDB(id uint16, status bool) (uint16, error)
	GetCouriersWithServiceFromDB() ([]Courier, error)
	UpdateCourierDB(courier Courier) error
	GetCouriersOfCourierServiceFromDB(limit, page, idService int) ([]Courier, int)
}

type DeliveryServiceRep interface {
	SaveDeliveryServiceInDB(service *DeliveryService) (int, error)
	GetDeliveryServiceByIdFromDB(Id int) (*DeliveryService, error)
	GetAllDeliveryServicesFromDB() ([]DeliveryService, error)
	UpdateDeliveryServiceInDB(service DeliveryService) error
	GetNumberCouriersByServiceFromDB(id int) (int, error)
}
