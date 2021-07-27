CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- courier data
CREATE TABLE IF NOT EXISTS collector_collectors
(
        uuid                        TEXT PRIMARY KEY NOT NULL,
        collector_meta              JSONB,
        phone_number                TEXT,
        status                      TEXT,
        status_change_timestamp     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        created_at                  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at                  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        deleted                     BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS collector_orders
(
    uuid                     text PRIMARY KEY         NOT NULL,
    id                       text                     NOT NULL,
    store_uuid               text,
    store_data               jsonb,
    device_id                text                     NOT NULL,
    client_uuid              text,
    client_data              jsonb,
    collector_uuid           text,
    collector_data           jsonb,
    callback_phone           text                     NOT NULL DEFAULT '',
    comment                  text                     NOT NULL DEFAULT '',
    source                   text                     NOT NULL,
    application              text                     NOT NULL,
    state                    text                     NOT NULL,
    deleted                  boolean                  NOT NULL DEFAULT FALSE,
    items                    jsonb,
    Promotion                jsonb,
    payment_type             text                     NOT NULL DEFAULT '',
    own_delivery             boolean                  NOT NULL DEFAULT FALSE,
    without_delivery         boolean                  NOT NULL DEFAULT FALSE,
    eat_in_store             boolean                  NOT NULL DEFAULT FALSE,
    total_price              numeric(8, 2)            NOT NULL DEFAULT 0,
    delivery_type            text                     NOT NULL DEFAULT '',
    delivery_price           numeric(8, 2)            NOT NULL DEFAULT 0,
    delivery_address         jsonb,
    delivery_address_details jsonb,
    cooking_time             integer                  NOT NULL DEFAULT 0,
    cooking_time_finish      timestamp WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_update_uuid         text                     NOT NULL DEFAULT '',
    last_update_role         text                     NOT NULL DEFAULT '',
    cancel_reason            text                     NOT NULL DEFAULT '',
    cancel_comment           text                     NOT NULL DEFAULT '',
    created_at               timestamp WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at               timestamp WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS collector_order_duplicate
(
    uuid            TEXT PRIMARY KEY NOT NULL,
    collector_uuid  TEXT,
    order_uuid      TEXT,
    client_data     JSONB,
    callback_phone  TEXT,
    comment         TEXT,
    cart_items      JSONB,
    total_price     FLOAT8,
    collect_time    INT
);

CREATE TABLE IF NOT EXISTS collector_products
(
    uuid        text PRIMARY KEY         NOT NULL,
    external_id text                     NOT NULL DEFAULT '',
    name        text                     NOT NULL,
    store_uuid  text                     NOT NULL,
    comment     text                     NOT NULL DEFAULT '',
    url         text                     NOT NULL DEFAULT '',
    deleted     boolean                  NOT NULL DEFAULT FALSE,
    composition text                     NOT NULL DEFAULT '',
    available   boolean                  NOT NULL DEFAULT TRUE,
    stop_list   boolean                  NOT NULL DEFAULT FALSE,
    default_set boolean                  NOT NULL DEFAULT FALSE,
    priority    double precision         NOT NULL DEFAULT 0,
    type        text                     NOT NULL DEFAULT 'single',
    leftover    int                      NOT NULL DEFAULT 0,
    price       numeric(8, 2)            NOT NULL DEFAULT 0,
    meta        jsonb,
    created_at  timestamp WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  timestamp WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS collector_barcodes
(
    uuid            TEXT PRIMARY KEY NOT NULL,
    bar_code        TEXT,
    product_uuid    TEXT REFERENCES collector_products(uuid)
);