package mysql

import (
	"food_delivery_api/cfg"
	"food_delivery_api/pkg/model"
	"sort"
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

	obj = model.Dashboard{
		Summary:     sum,
		CustomerTrx: ct,
		StockAlert:  sa,
		Top10Trx:    t10,
		Top5Product: t5p,
	}

	return obj, nil
}

func (s *Storage) ReadReportChart() (model.Chart, error) {
	var obj model.Chart

	// forming daily trx chart
	var dr []model.DailyRow
	s.db.Raw(`select name as product, sum(qty) as qty, sum(sub_total) as amount from transaction_lines
		where DATE(created_at) = CURDATE() group by name`).
		Scan(&dr)

	var quantities []model.Dataset
	var product []string
	var qtySold []int
	for _, v := range dr {
		quantities = append(quantities, model.Dataset{Label: v.Product, Data: []int{v.Qty}})
		product = append(product, v.Product)
		qtySold = append(qtySold, v.Qty)
	}
	qty := []model.Dataset{{Label: "# of Qty", Data: qtySold}}

	dtac := model.ChartData{ChartType: "Vertical Bar Chart", Labels: []string{"Today Transaction (Qty)"}, Datasets: quantities}
	dtqc := model.ChartData{ChartType: "Doughnut Chart", Labels: product, Datasets: qty}

	// forming monthly trx chart
	var mr []model.MonthlyRow
	s.db.Raw(`select month(created_at) as month, code as category, name as product, sum(qty) as qty, sum(sub_total) as amount
		from transaction_lines group by month, code, name`).
		Scan(&mr)

	var months []string
	var mproducts []string
	for _, v := range mr {
		m := convertMonth(v.Month)
		if !contains(months, m) {
			months = append(months, m)
		}

		if !contains(mproducts, v.Product) {
			mproducts = append(mproducts, v.Product)
		}
	}

	var mtacd []model.Dataset
	for _, v := range mproducts {
		var data []int
		for _, row := range mr {
			if row.Product == v {
				data = append(data, row.Qty)
			}
		}

		mtacd = append(mtacd, model.Dataset{Label: v, Data: data})
	}

	// Sort the dataset alphabetically based on the label
	sort.Slice(mtacd, func(i, j int) bool {
		return mtacd[i].Label < mtacd[j].Label
	})

	mtac := model.ChartData{ChartType: "Line Chart", Labels: months, Datasets: mtacd}

	obj = model.Chart{
		DailyTrxAmountChart:   dtac,
		DailyTrxQtyChart:      dtqc,
		MonthlyTrxAmountChart: mtac,
	}

	return obj, nil
}

func (s *Storage) ReadReportExponentialSmoothing(qp model.QueryGetExponentialSmoothing) ([]model.ExponentialSmoothingRow, error) {
	var list []model.ExponentialSmoothingRow

	s.db.Raw(`select date_format(created_at, '%Y-%m') as month, product_id, name, uom, sum(qty) as qty
		from transaction_lines where created_at between date(?) AND date(?) group by month, product_id, name, uom`,
		qp.StartDate, qp.EndDate).
		Scan(&list)

	return list, nil
}

func convertMonth(strMonth string) string {
	date, _ := time.Parse("1", strMonth)
	return date.Month().String()
}

func contains(list []string, value string) bool {
	for _, str := range list {
		if str == value {
			return true
		}
	}
	return false
}

func (s *Storage) ReadReportArima(qp model.QueryGetArima) ([]model.ArimaRow, error) {
	var list []model.ArimaRow

	s.db.Raw(`select date_format(date_sub(t.created_at, interval weekday(t.created_at) day), '%Y-%m-%d') as start_of_week,
		date_format(date_add(t.created_at, interval 6 - weekday(t.created_at) day), '%Y-%m-%d') as end_of_week,
       	sum(tl.qty) as total_qty
		from transactions t left join transaction_lines tl on t.id = tl.transaction_id
		where tl.product_id = ?
		group by start_of_week, end_of_week
		order by start_of_week`, qp.ProductID).
		Scan(&list)

	return list, nil
}
