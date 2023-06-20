package mysql

import (
	"food_delivery_api/cfg"
	"food_delivery_api/pkg/model"
	"time"
)

func (s *Storage) ReadReportDashboard() (model.Dashboard, error) {
	var obj model.Dashboard
	var err error

	// forming summary
	var ttlCtg, ttlPrd, ttlTrx, ttlUom, ttlCst int64
	err = s.db.Model(&model.Category{}).Count(&ttlCtg).Error
	err = s.db.Model(&model.Product{}).Count(&ttlPrd).Error
	err = s.db.Model(&model.Transaction{}).Count(&ttlTrx).Error
	err = s.db.Model(&model.UOM{}).Count(&ttlUom).Error
	err = s.db.Model(&model.Transaction{}).Distinct("customer").Count(&ttlCst).Error
	if err != nil {
		return obj, err
	}

	sum := model.Summary{
		TotalCategory:    ttlCtg,
		TotalProduct:     ttlPrd,
		TotalTransaction: ttlTrx,
		TotalUOM:         ttlUom,
		TotalCustomer:    ttlCst,
	}

	// forming customer trx
	var ct []model.CustomerTrx
	s.db.Raw(`select customer, count(*) as total_trx, sum(total) as amount_trx, avg(total) as average_trx 
		from transactions group by customer`).
		Scan(&ct)

	// forming stock alert
	var sa []model.StockAlert
	s.db.Raw(`select p.id, p.code, p.name, p.qty, u.name as UOM from products p join uoms u on u.id = p.uom_id
		where qty < 100 order by qty;`).
		Scan(&sa)

	// forming top 10 trx
	var t10 []model.Top10Trx
	s.db.Raw(`select created_at, trx_id, created_by, customer, total from transactions order by total desc limit 10`).
		Scan(&t10)

	for i, v := range t10 {
		t10[i].TrxDate = v.CreatedAt.Format(cfg.AppTLayout)
	}

	// forming top 5 product
	var t5p []model.Top5Product
	s.db.Raw(`select product_id as id, name, price, sum(qty) as total_qty, sum(sub_total) / sum(qty) as average_amount,
    	sum(sub_total) as total_amount from transaction_lines group by product_id, name, price limit 5`).
		Scan(&t5p)

	// forming daily trx chart
	var dc []model.DailyTrxChart
	s.db.Raw(`select day(created_at) as day, month(created_at) as month, year(created_at) as year, sum(total) as amount
		from transactions group by day(created_at), month(created_at), year(created_at)`).
		Scan(&dc)

	for i, v := range dc {
		dc[i].Month = convertMonth(v.Month)
	}

	// forming monthly trx chart
	var mc []model.MonthlyTrxChart
	s.db.Raw(`select month(created_at) as month, year(created_at) as year, sum(total) as amount from transactions group by year(created_at), month(created_at)`).
		Scan(&mc)

	for i, v := range mc {
		mc[i].Month = convertMonth(v.Month)
	}

	obj = model.Dashboard{
		Summary:         sum,
		CustomerTrx:     ct,
		StockAlert:      sa,
		Top10Trx:        t10,
		Top5Product:     t5p,
		DailyTrxChart:   dc,
		MonthlyTrxChart: mc,
	}

	return obj, nil
}

func convertMonth(strMonth string) string {
	date, _ := time.Parse("1", strMonth)
	return date.Month().String()
}

func (s *Storage) ReadReportMetodologi(qp model.MetodologiFilter) ([]model.ResponseMetodologi, error) {

	return []model.ResponseMetodologi{}, nil
}
