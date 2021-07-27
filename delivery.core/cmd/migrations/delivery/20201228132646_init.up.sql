CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- courier data
CREATE TABLE IF NOT EXISTS delivery_couriers
(
    uuid                            TEXT PRIMARY KEY         NOT NULL,
    courier_meta                    JSONB,
    chat_id                         TEXT,
    phone_number                    TEXT                     NOT NULL,
    courier_type                    TEXT,
    status                          TEXT                     NOT NULL DEFAULT 'on moderation',
    status_change_timestamp         TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_lat                        FLOAT,
    last_lon                        FLOAT,
    latlon_timestamp                TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at                      TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS delivery_orders
(
    uuid                        TEXT PRIMARY KEY         NOT NULL,
    order_number                TEXT                     NOT NULL,
    delivery_items              JSONB,
    comment                     TEXT,
    pickup_request_time         TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    drop_off_request_time       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    pickup_person_contacts      JSONB,
    drop_off_person_contacts    JSONB,
    pickup_route                JSONB,
    dropoff_route               JSONB,
    service                     TEXT                     NOT NULL,
    created_at                  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    finished_at                 TIMESTAMP WITH TIME ZONE,
    canceled_at                 TIMESTAMP WITH TIME ZONE,
    cancel_reason               TEXT,
    payment_type                TEXT                     NOT NULL,
    payment_status              TEXT                     NOT NULL,
    state                       TEXT                     NOT NULL DEFAULT 'created',
    deleted_at                  TIMESTAMP WITH TIME ZONE,
    delivery_price              NUMERIC(8,2)             NOT NULL,
    courier_uuid                TEXT,
    courier_data                JSONB
);

CREATE TABlE IF NOT EXISTS delivery_events
(
    uuid            TEXT PRIMARY KEY         NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    event_type      TEXT                     NOT NULL,
    payload         JSONB,
    courier_uuid    TEXT                     REFERENCES delivery_couriers(uuid),
    order_uuid      TEXT                     REFERENCES delivery_orders(uuid)
);


CREATE TABLE IF NOT EXISTS delivery_tasks
(
    uuid                    TEXT PRIMARY KEY            NOT NULL,
    order_number            TEXT,
    client_phone            TEXT,
    client_name             TEXT                        NOT NULL,
    order_uuid              TEXT,
    type                    TEXT                        NOT NULL,
    route                   JSONB                       NOT NULL,
    state                   TEXT                        ,
    created_at              TIMESTAMP WITH TIME ZONE                            DEFAULT CURRENT_TIMESTAMP,
    expected_time           TIMESTAMP WITH TIME ZONE                            DEFAULT CURRENT_TIMESTAMP,
    last_update_time        TIMESTAMP WITH TIME ZONE,
    finished_time           TIMESTAMP WITH TIME ZONE,
    courier_uuid            TEXT
);

CREATE TABLE IF NOT EXISTS delivery_users
(
    uuid                    TEXT PRIMARY KEY            NOT NULL,
    login                   TEXT                        NOT NULL,
    state                   TEXT                        NOT NULL                DEFAULT 'на модерации',
    meta                    JSONB,
    deleted                 BOOLEAN                     NOT NULL                DEFAULT false
);

CREATE TABLE IF NOT EXISTS delivery_courier_coordinates
(
  uuid                      TEXT PRIMARY KEY            NOT NULL,
  courier_uuid              TEXT                        NOT NULL        REFERENCES delivery_couriers(uuid),
  lat                       FLOAT,
  lon                       FLOAT,
  latlon_timestamp          TIMESTAMP WITH TIME ZONE                  DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS delivery_courier_queue
(
    uuid                    TEXT    PRIMARY KEY         NOT NULL,
    courier_uuid            TEXT                        NOT NULL        REFERENCES delivery_couriers(uuid),
    tasks                   JSONB
);
