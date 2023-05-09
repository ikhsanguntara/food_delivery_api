package mysql

import (
	"food_delivery_api/pkg/model"

	"github.com/gin-gonic/gin"
)

func (s *Storage) CreateUOM(obj model.UOM) (model.UOM, error) {
	err := s.db.Create(&obj).Error
	if err != nil {
		return obj, err
	}

	return obj, nil
}

func (s *Storage) ReadUOMs(c *gin.Context) ([]model.UOM, int64, error) {
	var list []model.UOM
	var ttl int64

	s.db.Find(&list).Count(&ttl)
	err := s.db.Scopes(Paginate(c)).Find(&list).Error
	if err != nil {
		return list, ttl, err
	}

	return list, ttl, nil
}

func (s *Storage) ReadUOM(obj model.UOM) (model.UOM, error) {
	err := s.db.First(&obj, obj.ID).Error
	if err != nil {
		return obj, err
	}

	return obj, nil
}

func (s *Storage) UpdateUOM(obj model.UOM) (model.UOM, error) {
	err := s.db.Model(&obj).Updates(obj).Error
	if err != nil {
		return obj, err
	}

	return obj, nil
}

func (s *Storage) DeleteUOM(obj model.UOM) (model.UOM, error) {
	err := s.db.Delete(&obj, obj.ID).Error
	if err != nil {
		return obj, err
	}

	return obj, nil
}
