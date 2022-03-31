package dao

import (
	"database/sql"
	"fmt"
	"log"
)

type DeliveryServicePostgres struct {
	db *sql.DB
}

func NewDeliveryServicePostgres(db *sql.DB) *DeliveryServicePostgres {
	return &DeliveryServicePostgres{db: db}
}

type DeliveryService struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Photo         string `json:"photo"`
	Description   string `json:"description"`
	PhoneNumber   string `json:"phone_number"`
	ManagerId     int    `json:"manager_id"`
	Status        string `json:"status"`
	NumOfCouriers int
}

func (r *DeliveryServicePostgres) SaveDeliveryServiceInDB(service *DeliveryService) (int, error) {
	row := r.db.QueryRow(`INSERT INTO delivery_service (name, email, photo, description,
                              phone_number,manager_id, status) VALUES ($1, $2, $3, $4, $5,$6, $7) RETURNING id`,
		service.Name, service.Email, service.Photo, service.Description,
		service.PhoneNumber, service.ManagerId, service.Status)
	var id int
	if err := row.Scan(&id); err != nil {
		log.Println(fmt.Sprintf("Create Delivery : error:%s", err))
		return 0, fmt.Errorf("Create Delivery Service: error:%s", err)
	}
	return id, nil
}

func (r *DeliveryServicePostgres) GetDeliveryServiceByIdFromDB(Id int) (*DeliveryService, error) {
	var service DeliveryService
	res, err := r.db.Query("SELECT id, name,email,photo,description,phone_number,manager_id,status FROM delivery_service Where manager_id=$1", Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for res.Next() {
		err = res.Scan(&service.Id, &service.Name, &service.Email, &service.Photo, &service.Description,
			&service.PhoneNumber, &service.ManagerId, &service.Status)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return &service, nil
}

//SELECT count(*) FROM couriers AS co JOIN delivery_service AS d ON co.delivery_service_id=d.id WHERE d.id=2
func (r *DeliveryServicePostgres) GetAllDeliveryServicesFromDB() ([]DeliveryService, error) {
	var services []DeliveryService
	res, err := r.db.Query(`SELECT id, name, email, photo, description, phone_number, manager_id, status
                                  FROM delivery_service ORDER BY id`)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for res.Next() {
		var service DeliveryService
		err = res.Scan(&service.Id, &service.Name, &service.Email, &service.Photo, &service.Description,
			&service.PhoneNumber, &service.ManagerId, &service.Status)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		services = append(services, service)
	}

	return services, nil
}

func (r *DeliveryServicePostgres) UpdateDeliveryServiceInDB(service DeliveryService) error {
	var oldService DeliveryService
	transaction, err := r.db.Begin()
	if err != nil {
		log.Println(err)
	}
	defer transaction.Commit()
	res, err := transaction.Query(`SELECT id, name,email,photo,description,phone_number,manager_id,status 
                                  FROM delivery_service Where id=$1`, service.Id)
	if err != nil {
		log.Println(err)
		return err
	}
	for res.Next() {
		err = res.Scan(&oldService.Id, &oldService.Name, &oldService.Email,
			&oldService.Photo, &oldService.Description, &oldService.PhoneNumber,
			&oldService.ManagerId, &oldService.Status)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	if service.Name == "" {
		service.Name = oldService.Name
	}
	if service.Email == "" {
		service.Email = oldService.Email
	}
	if service.Photo == "" {
		service.Photo = oldService.Photo
	}
	if service.Description == "" {
		service.Description = oldService.Description
	}
	if service.PhoneNumber == "" {
		service.PhoneNumber = oldService.PhoneNumber
	}
	if service.Status == "" {
		service.Status = oldService.Status
	}

	s := `UPDATE delivery_service SET name = $1, email = $2, description = $3, 
                            phone_number = $4, status = $5, photo=$6 WHERE id = $7`
	log.Println(s)
	insert, err := transaction.Query(s, service.Name, service.Email, service.Description,
		service.PhoneNumber, service.Status, &service.Photo, service.Id)
	defer insert.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *DeliveryServicePostgres) GetNumberCouriersByServiceFromDB(id int) (int, error) {

	selectValue := `SELECT count(*) FROM couriers AS co JOIN delivery_service AS d ON co.delivery_service_id=d.id WHERE d.id=$1`
	get, err := r.db.Query(selectValue, id)

	if err != nil {
		log.Println("Error of getting list of couriers :" + err.Error())
		return 0, err
	}

	var numOfCouriers int
	for get.Next() {
		err = get.Scan(&numOfCouriers)
	}
	return numOfCouriers, nil
}
