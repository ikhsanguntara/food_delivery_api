package service

import (
	"food_delivery_api/pkg/model"
)

func (s *service) AddUser(p model.User) (model.User, error) {
	obj, err := s.rmy.CreateUser(p)
	if err != nil {
		return obj, err
	}

	// hide password
	obj.Password = ""

	return obj, nil
}

func (s *service) GetUsers() ([]model.User, error) {
	list, err := s.rmy.ReadUsers()
	if err != nil {
		return list, err
	}

	// hide password
	for i := range list {
		list[i].Password = ""
	}

	return list, nil
}

func (s *service) GetUser(p model.User) (model.User, error) {
	obj, err := s.rmy.ReadUser(p)
	if err != nil {
		return obj, err
	}

	// hide password
	obj.Password = ""

	return obj, nil
}

func (s *service) GetUserByEmailPassword(p model.User) (model.User, error) {
	obj, err := s.rmy.ReadUserByEmailPassword(p)
	if err != nil {
		return obj, err
	}

	// hide password
	obj.Password = ""

	return obj, nil
}

func (s *service) EditUser(p model.User) (model.User, error) {
	obj, err := s.rmy.ReadUser(p)
	if err != nil {
		return obj, err
	}

	obj, err = s.rmy.UpdateUser(p)
	if err != nil {
		return obj, err
	}

	// hide password
	obj.Password = ""

	return obj, nil
}

func (s *service) RemoveUser(p model.User) (model.User, error) {
	obj, err := s.rmy.ReadUser(p)
	if err != nil {
		return obj, err
	}

	_, err = s.rmy.DeleteUser(p)
	if err != nil {
		return obj, err
	}

	// hide password
	obj.Password = ""

	return obj, nil
}
