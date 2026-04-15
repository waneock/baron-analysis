CREATE TABLE IF NOT EXISTS orders (
    id                  text PRIMARY KEY,
    price               INTEGER NOT NULL DEFAULT 0,
    commission          INTEGER NOT NULL DEFAULT 0,
    tax                 INTEGER NOT NULL DEFAULT 0,

    classid             text NOT NULL,
    instanceid          text NOT NULL,
    appid               INTEGER NOT NULL,
    contextid           text NOT NULL,
    assetid             text NOT NULL,

    name                text NOT NULL,
    offerid             text,
    state               INTEGER NOT NULL DEFAULT 0,

    escrow_end_date     timestamptz,
    list_time           timestamptz,
    last_updated        timestamptz,

    wear                INTEGER NOT NULL DEFAULT 0,
    txid                text,
    trade_locked        boolean NOT NULL DEFAULT false,

    addons              text[] NULL,
    buyer_country_code  varchar(10)
);