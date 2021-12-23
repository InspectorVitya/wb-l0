-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
                        orderuid varchar(128) NOT NULL,
                        track_number varchar(128) NULL,
                        entry varchar(128) NULL,
                        locale varchar(128) NULL,
                        internal_signature varchar(128) NULL,
                        customer_id varchar(128) NULL,
                        delivery_service varchar(128) NULL,
                        shardkey varchar(128) NULL,
                        sm_id int4 NULL,
                        date_created timestamp NULL,
                        oof_shard varchar(128) NULL,
                        CONSTRAINT orders_pk PRIMARY KEY (orderuid)
);
CREATE TABLE delivery (
                          "name" varchar(128) NULL,
                          phone varchar(128) NULL,
                          zip varchar(128) NULL,
                          city varchar(128) NULL,
                          address varchar(128) NULL,
                          region varchar(128) NULL,
                          email varchar(128) NULL,
                          order_id varchar(128) NOT NULL,
                          CONSTRAINT delivery_fk FOREIGN KEY (order_id) REFERENCES orders(orderuid)
);
CREATE TABLE items (
                       chrt_id int4 NULL,
                       track_number varchar(256) NULL,
                       price int4 NULL,
                       rid varchar(256) NULL,
                       "name" varchar(128) NULL,
                       sale int4 NULL,
                       "size" varchar(128) NULL,
                       total_price int4 NULL,
                       nm_id int4 NULL,
                       brand varchar(128) NULL,
                       status int4 NULL,
                       order_id varchar NOT NULL,
                       CONSTRAINT items_fk FOREIGN KEY (order_id) REFERENCES orders(orderuid)
);
CREATE TABLE payment (
                         "transaction" varchar(256) NULL,
                         request_id varchar(256) NULL,
                         currency varchar(128) NULL,
                         provider varchar(128) NULL,
                         amount int4 NULL,
                         payment_dt int4 NULL,
                         bank varchar(128) NULL,
                         delivery_cost int4 NULL,
                         goods_total int4 NULL,
                         custom_fee int4 NULL,
                         order_id varchar(128) NOT NULL,
                         CONSTRAINT payment_fk FOREIGN KEY (order_id) REFERENCES orders(orderuid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE payment;
DROP TABLE items;
DROP TABLE delivery;
DROP TABLE orders;
-- +goose StatementEnd
