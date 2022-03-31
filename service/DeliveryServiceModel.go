package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/minio/minio-go"
	"log"
	"stlab.itechart-group.com/go/food_delivery/courier_service/dao"
	"strconv"
)

func (s *CourierService) CreateDeliveryService(DeliveryService dao.DeliveryService) (int, error) {
	id, err := s.repo.SaveDeliveryServiceInDB(&DeliveryService)
	if err != nil {
		log.Println(err)
		return 0, fmt.Errorf("Error in DeliveryServiceService: %s", err)
	}
	return id, nil
}

func (s *CourierService) GetDeliveryServiceById(Id int) (*dao.DeliveryService, error) {
	var service *dao.DeliveryService
	service, err := s.repo.GetDeliveryServiceByIdFromDB(Id)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Error in DeliveryService: %s", err)
	}
	service.NumOfCouriers, err = s.repo.GetNumberCouriersByServiceFromDB(Id)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Error in DeliveryService: %s", err)
	}
	if service.Id == 0 {
		err = errors.New("not found")
		log.Println(err)
		return nil, fmt.Errorf("Error in DeliveryService: %s", err)
	}
	if service.Status == "inactive" {
		err := errors.New("account deleted")
		log.Println("account deleted")
		return nil, fmt.Errorf("Error in DeliveryService: %s", err)
	}
	return service, nil
}

func (s *CourierService) GetAllDeliveryServices() ([]dao.DeliveryService, error) {
	var Services = []dao.DeliveryService{}
	Services, err := s.repo.GetAllDeliveryServicesFromDB()
	if err != nil {
		log.Println(err)
		return []dao.DeliveryService{}, fmt.Errorf("Error in DeliveryService: %s", err)
	}
	Couriers, err := s.repo.GetCouriersWithServiceFromDB()
	if err != nil {
		log.Println(err)
		return []dao.DeliveryService{}, fmt.Errorf("Error in DeliveryService: %s", err)
	}
	for i, service := range Services {
		count := 0
		for _, courier := range Couriers {
			if service.Id == int(courier.DeliveryServiceId) {
				count++
			}
		}
		Services[i].NumOfCouriers = count
	}
	return Services, nil
}

func (s *CourierService) UpdateDeliveryService(service dao.DeliveryService) error {
	if err := s.repo.UpdateDeliveryServiceInDB(service); err != nil {
		log.Println(err)
		return fmt.Errorf("Error in DeliveryService: %s", err)
	}
	return nil
}
func (s *CourierService) SaveLogoFile(cover []byte, id int) error {
	client, err := InitClientDO()
	if err != nil {
		log.Println(err)
		return err
	}

	_, err1 := client.PutObject("storage-like-s3", fmt.Sprintf("logo_img/%s", strconv.Itoa(id)),
		bytes.NewReader(cover), int64(len(cover)), minio.PutObjectOptions{ContentType: "image/jpeg", UserMetadata: map[string]string{"x-amz-acl": "public-read"}})
	if err1 != nil {
		log.Println(err1)
		return err1
	}

	var service dao.DeliveryService
	service.Id = id
	service.Photo = "https://storage-like-s3.fra1.digitaloceanspaces.com/logo_img/" + strconv.Itoa(id)

	if err := s.repo.UpdateDeliveryServiceInDB(service); err != nil {
		log.Println(err)
		return fmt.Errorf("Error in DeliveryService: %s", err)
	}

	log.Println("Uploaded logo with link https://storage-like-s3.fra1.digitaloceanspaces.com/logo_img/" + strconv.Itoa(id))
	return nil
}
