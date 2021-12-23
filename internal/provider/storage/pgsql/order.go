package pgsql

import (
	"context"
	"github.com/inspectorvitya/wb-l0/internal/model"
)

func (s *StoragePgSql) CreateOrder(ctx context.Context, order *model.Order) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.NamedExecContext(ctx, `INSERT INTO orders (orderuid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
 	VALUES(:orderuid, :track_number, :entry, :locale, :internal_signature, :customer_id, :delivery_service, :shardkey, :sm_id, :date_created, :oof_shard );`, order)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO delivery ("name", phone, zip, city, address, region, email, order_id)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8);`,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
		order.OrderUID,
	)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO payment ("transaction", request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee, order_id)
 	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
		order.OrderUID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}
	for _, val := range order.Items {
		_, err = tx.ExecContext(ctx, `INSERT INTO items (chrt_id, track_number, price, rid, "name", sale, "size", total_price, nm_id, brand, status, order_id)
 		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`,
			val.ChrtID, val.TrackNumber, val.Price, val.Rid, val.Name, val.Sale, val.Size, val.TotalPrice, val.NmID, val.Brand, val.Status, order.OrderUID,
		)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	}
	err = tx.Commit()
	return err
}
func (s *StoragePgSql) GetByID(ctx context.Context, id string) (*model.Order, error) {
	order := &model.Order{}
	tx, err := s.db.BeginTxx(ctx, nil)
	err = tx.GetContext(ctx, order,
		`select orderuid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard from orders where orderuid = $1;`, id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}
	err = tx.GetContext(ctx, order.Delivery,
		`select name, phone, zip, city, address, region, email from delivery where order_id = $1;`, id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}
	err = tx.GetContext(ctx, order.Payment,
		`select "transaction", request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee from payment where order_id = $1;`, id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}
	err = tx.SelectContext(ctx, order.Items, `SELECT chrt_id, track_number, price, rid, "name", sale, "size", total_price, nm_id, brand, status FROM items where order_id = $1;`, id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	return order, nil
}
func (s *StoragePgSql) GetAll(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	tx, err := s.db.BeginTxx(ctx, nil)
	err = tx.SelectContext(ctx, &orders,
		`select orderuid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard from orders`)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	for i := range orders {
		err = tx.GetContext(ctx, &orders[i].Delivery,
			`select name, phone, zip, city, address, region, email from delivery where order_id = $1;`, orders[i].OrderUID)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return nil, rollbackErr
			}
			return nil, err
		}
		err = tx.GetContext(ctx, &orders[i].Payment,
			`select "transaction", request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee from payment where order_id = $1;`, orders[i].OrderUID)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return nil, rollbackErr
			}
			return nil, err
		}
		err = tx.SelectContext(ctx, &orders[i].Items, `SELECT chrt_id, track_number, price, rid, "name", sale, "size", total_price, nm_id, brand, status FROM items where order_id = $1;`, orders[i].OrderUID)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return nil, rollbackErr
			}
			return nil, err
		}
	}
	return orders, nil
}
