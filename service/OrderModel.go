package service

import (
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	courierProto "stlab.itechart-group.com/go/food_delivery/courier_service/GRPC"
	"stlab.itechart-group.com/go/food_delivery/courier_service/dao"
)

func (s *CourierService) GetOrders(id int) ([]dao.Order, error) {
	get, err := s.repo.GetActiveOrdersFromDB(id)
	if get == nil {
		return []dao.Order{}, fmt.Errorf("Error in OrderService: %s", err)
	}
	if err != nil {
		return nil, fmt.Errorf("Error with database: %s", err)
	}
	if id == 0 {
		err := errors.New("no id")
		log.Println("id cannot be zero")
		return nil, fmt.Errorf("Error in OrderService: %s", err)
	}
	return get, nil
}

func (s *CourierService) GetOrder(id int) (dao.Order, error) {
	get, err := s.repo.GetActiveOrderFromDB(id)
	if (get == dao.Order{}) {
		return dao.Order{}, fmt.Errorf("Error in OrderService: %s", err)
	}
	if err != nil {
		return dao.Order{}, fmt.Errorf("Error with database: %s", err)
	}
	if id == 0 {
		err := errors.New("no id")
		log.Println("id cannot be zero")
		return dao.Order{}, fmt.Errorf("Error in OrderService: %s", err)
	}
	return get, nil
}

func (s *CourierService) GetOrderForChange(id int) (dao.Order, error) {
	get, err := s.repo.GetOrderFromDB(id)
	if (get == dao.Order{}) {
		return dao.Order{}, fmt.Errorf("Error in OrderService: %s", err)
	}
	if err != nil {
		return dao.Order{}, fmt.Errorf("Error with database: %s", err)
	}
	if id == 0 {
		err := errors.New("no id")
		log.Println("id cannot be zero")
		return dao.Order{}, fmt.Errorf("Error in OrderService: %s", err)
	}
	return get, nil
}

func (s *CourierService) ChangeOrderStatus(text string, id uint16) (uint16, error) {
	_, err := s.GetOrderForChange(int(id))
	if err != nil {
		return 0, fmt.Errorf("Error in OrderService: %s", err)
	}
	orderId, err := s.repo.ChangeOrderStatusInDB(text, id)
	if err != nil {
		return 0, fmt.Errorf("Error with database: %s", err)
	}
	return orderId, nil
}

func (s *CourierService) GetCourierCompletedOrders(limit, page, idCourier int) ([]dao.DetailedOrder, error) {
	var Order = []dao.DetailedOrder{}

	if limit <= 0 || page <= 0 {
		err := errors.New("no page or limit")
		log.Println("no more pages or limit")
		return nil, fmt.Errorf("Error in OrderService: %s", err)
	}
	Order, totalCount := s.repo.GetCourierCompletedOrdersWithPage_fromDB(limit, page, idCourier)
	LimitOfPages := (totalCount / limit) + 1
	if LimitOfPages < page {
		err := errors.New("no page")
		log.Println("no more pages")
		return nil, fmt.Errorf("Error in OrderService: %s", err)
	}
	fmt.Println(Order)
	return Order, nil
}

func (s *CourierService) GetAllOrdersOfCourierService(limit, page, idService int) ([]dao.DetailedOrder, error) {
	var Order = []dao.DetailedOrder{}
	if limit <= 0 || page <= 0 {
		err := errors.New("no page or limit")
		log.Println("no more pages or limit")
		return nil, fmt.Errorf("Error in OrderService: %s", err)
	}
	Order, totalCount := s.repo.GetAllOrdersOfCourierServiceWithPageFromDB(limit, page, idService)
	LimitOfPages := (totalCount / limit) + 1
	if LimitOfPages < page {
		err := errors.New("no page")
		log.Println("no more pages")
		return nil, fmt.Errorf("Error in OrderService: %s", err)
	}
	fmt.Println(Order)
	return Order, nil
}

func (s *CourierService) GetCourierCompletedOrdersByMonth(limit, page, idService, Month, Year int) ([]dao.Order, error) {
	var Order = []dao.Order{}
	if limit <= 0 || page <= 0 {
		err := errors.New("no page or limit")
		log.Println("no more pages or limit")
		return nil, fmt.Errorf("Error in OrderService: %s", err)
	}
	Order, totalCount := s.repo.GetCourierCompletedOrdersByMouthWithPageFromDB(limit, page, idService, Month, Year)
	LimitOfPages := (totalCount / limit) + 1
	if LimitOfPages < page {
		err := errors.New("no page")
		log.Println("no more pages")
		return nil, fmt.Errorf("Error in OrderService: %s", err)
	}
	if Month >= 13 || Month < 1 {
		err := errors.New("enter correct month")
		log.Println("enter correct month")
		return nil, fmt.Errorf("Error in OrderService: %s", err)
	}

	return Order, nil
}

func (s *CourierService) AssigningOrderToCourier(order dao.Order) error {
	if err := s.repo.AssigningOrderToCourierInDB(order); err != nil {
		log.Println(err)
		return fmt.Errorf("Error in OrderService: %s", err)
	}
	return nil
}

func (s *CourierService) GetDetailedOrderById(Id int) (*dao.AllInfoAboutOrder, error) {
	var Order *dao.AllInfoAboutOrder
	Order, err := s.repo.GetDetailedOrderByIdFromDB(Id)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Error in OrderService: %s", err)
	}
	if Order.IdOrder == 0 {
		err = errors.New("not found")
		log.Println(err)
		return nil, fmt.Errorf("Error in OrderService: %s", err)
	}
	return Order, nil
}

func (s *CourierService) CreateOrder(order *courierProto.OrderCourierServer) (*emptypb.Empty, error) {
	return s.repo.OrderRep.CreateOrder(order)
}

func (s *CourierService) GetServices(in *emptypb.Empty) (*courierProto.ServicesResponse, error) {
	return s.repo.OrderRep.GetServices(in)
}
func (s *CourierService) GetCompletedOrdersOfCourierService(limit, page, idService int) ([]dao.Order, error) {
	var Order = []dao.Order{}
	Order, totalCount := s.repo.GetCompletedOrdersOfCourierServiceFromDB(limit, page, idService)
	LimitOfPages := (totalCount / limit) + 1
	if LimitOfPages < page {
		err := errors.New("no page")
		log.Println("no more pages")
		return nil, fmt.Errorf("Error in OrderService: %s", err)
	}
	return Order, nil
}

func (s *CourierService) GetCompletedOrdersOfCourierServiceByDate(limit, page, idService int) ([]dao.Order, error) {
	var Order = []dao.Order{}
	Order, totalCount := s.repo.GetCompletedOrdersOfCourierServiceByDateFromDB(limit, page, idService)
	LimitOfPages := (totalCount / limit) + 1
	if LimitOfPages < page {
		err := errors.New("no page")
		log.Println("no more pages")
		return nil, fmt.Errorf("Error in OrderService: %s", err)
	}
	return Order, nil
}

func (s *CourierService) GetCompletedOrdersOfCourierServiceByCourierId(limit, page, idService int) ([]dao.Order, error) {
	var Order = []dao.Order{}
	Order, totalCount := s.repo.GetCompletedOrdersOfCourierServiceByCourierIdFromDB(limit, page, idService)
	LimitOfPages := (totalCount / limit) + 1
	if LimitOfPages < page {
		err := errors.New("no page")
		log.Println("no more pages")
		return nil, fmt.Errorf("Error in OrderService: %s", err)
	}
	return Order, nil
}

func (s *CourierService) GetOrdersOfCourierServiceForManager(limit, page, idService int) ([]dao.DetailedOrder, error) {
	var Order = []dao.DetailedOrder{}
	Order, totalCount := s.repo.GetOrdersOfCourierServiceForManagerFromDB(limit, page, idService)
	LimitOfPages := (totalCount / limit) + 1
	if LimitOfPages < page {
		err := errors.New("no page")
		log.Println("no more pages")
		return nil, fmt.Errorf("Error in OrderService: %s", err)
	}
	fmt.Println(Order)
	return Order, nil
}
