package service

import (
	"food_delivery_api/pkg/model"
	"math"
	"time"
)

func (s *service) GetReportDashboard() (model.Dashboard, error) {
	obj, err := s.rmy.ReadReportDashboard()
	if err != nil {
		return obj, err
	}

	return obj, nil
}

func (s *service) GetReportMetodologi(input model.MetodologiFilter) ([]model.ResponseMetodologi, error) {
	data := []model.ResponseMetodologi{}
	total, err := s.rmy.TotalData(input)
	if err != nil {
		return data, err
	}

	var aBefore float64
	var fBefore float64
	var fPredictions float64
	for _, row := range total {
		month, err := time.Parse("2006-01", row.Month)
		if err != nil {
			return nil, err
		}

		if row.Month == input.StartDate[:7] {
			aBefore = row.TotalQty
			fBefore = row.TotalQty
			data = append(data, model.ResponseMetodologi{
				Month:       month.Format("January 2006"),
				TotalQty:    row.TotalQty,
				Predictions: math.Round(row.TotalQty)})

		} else if row.Month == input.EndDate[:7] {
			fPredictions := fBefore + ((2 / (float64(len(total)) + 1)) * (aBefore - fBefore))
			data = append(data, model.ResponseMetodologi{
				Month:       month.Format("January 2006"),
				TotalQty:    row.TotalQty,
				Predictions: math.Round(fPredictions)})

			aBefore = row.TotalQty
			fBefore = fPredictions
			fPredictions = fBefore + ((2 / (float64(len(total)) + 2)) * (aBefore - fBefore))
			data = append(data, model.ResponseMetodologi{
				Month:       month.AddDate(0, 1, 0).Format("January 2006"),
				TotalQty:    0,
				Predictions: math.Round(fPredictions)})
		} else {
			fPredictions = (fBefore + ((2 / (float64(len(total)) + 1)) * (aBefore - fBefore)))
			data = append(data, model.ResponseMetodologi{
				Month:       month.Format("January 2006"),
				TotalQty:    row.TotalQty,
				Predictions: math.Round(fPredictions)})

			aBefore = row.TotalQty
			fBefore = fPredictions
		}
	}

	return data, nil
}
