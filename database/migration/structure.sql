CREATE TABLE IF NOT EXISTS orders
(
    order_uid          varchar PRIMARY KEY NOT NULL,
    track_number       varchar,
    entry              varchar,
    locale             varchar,
    internal_signature varchar,
    customer_id        varchar,
    delivery_service   varchar,
    shardkey           varchar,
    sm_id              bigint,
    date_created       timestamp,
    oof_shard          varchar
);

CREATE TABLE IF NOT EXISTS deliveries
(
    order_uid varchar NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
    name      varchar,
    phone     varchar,
    zip       varchar,
    city      varchar,
    address   varchar,
    region    varchar,
    email     varchar
);

CREATE TABLE IF NOT EXISTS payments
(
    order_uid     varchar NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction   varchar,
    request_id    varchar,
    currency      varchar,
    provider      varchar,
    amount        bigint,
    payment_dt    bigint,
    bank          varchar,
    delivery_cost bigint,
    goods_total   bigint,
    custom_fee    bigint
);

CREATE TABLE IF NOT EXISTS items
(
    order_uid    varchar NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id      bigint,
    track_number varchar,
    price        bigint,
    rid          varchar,
    name         varchar,
    sale         bigint,
    size         varchar,
    total_price  bigint,
    nm_id        bigint,
    brand        varchar,
    status       bigint
);