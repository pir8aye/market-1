CREATE TABLE IF NOT EXISTS pricecoinmarketcap (
	id SERIAL PRIMARY KEY,
	asset_id  character varying(32) NOT NULL,  
	name character varying(32) NOT NULL,
	symbol character varying(32) NOT NULL,
	rank character varying(32) NOT NULL,
	price_usd character varying(64) NOT NULL,
	price_btc character varying(64) NOT NULL,
	volume_usd_24h character varying(64) NOT NULL,
	market_cap_usd character varying(64) NOT NULL,
	available_supply character varying(64) NOT NULL,
	total_supply character varying(64) NOT NULL,
	percent_change_1h character varying(64) NOT NULL,
	percent_change_24h character varying(64) NOT NULL,
	percent_change_7d character varying(64) NOT NULL,
	last_updated character varying(64) NOT NULL,
	price_cny character varying(64) NOT NULL,
	volume_cny_24h character varying(64) NOT NULL,
	market_cap_cny character varying(64) NOT NULL 
);

CREATE INDEX pricecoinmarketcap_last_updated_symbol_index ON pricecoinmarketcap USING btree(last_updated COLLATE "default" DESC NULLS FIRST, symbol COLLATE "default" ASC NULLS LAST);
CREATE INDEX pricecoinmarketcap_symbol ON pricecoinmarketcap USING btree(symbol COLLATE "default" ASC NULLS LAST);

CREATE TABLE IF NOT EXISTS coinmarketcapmin (
	id SERIAL PRIMARY KEY,
	asset_id character varying(32) NOT NULL,  
	name character varying(32) NOT NULL, 
	symbol character varying(32) NOT NULL, 
	rank integer  NOT NULL, 
	price_usd_first real  NOT NULL, 
	price_usd_last real NOT NULL, 
	price_usd_low real  NOT NULL, 
	price_usd_high  real NOT NULL, 
	price_btc_first real  NOT NULL, 
	price_btc_last real  NOT NULL, 
	price_btc_low real  NOT NULL, 
	price_btc_high real  NOT NULL, 
	price_cny_first real  NOT NULL, 
	price_cny_last real  NOT NULL, 
	price_cny_low real  NOT NULL, 
	price_cny_high real  NOT NULL, 
	last_updated bigint NOT NULL,
	timestamp bigint NOT NULL,
       _group character varying(32) 
); 

CREATE INDEX index_timestamp_coinmarketcapmin ON coinmarketcapmin (timestamp);
CREATE INDEX index_last_updated_coinmarketcapmin ON coinmarketcapmin (last_updated);
CREATE INDEX index_symbol_coinmarketcapmin ON coinmarketcapmin USING hash (symbol);

CREATE TABLE IF NOT EXISTS coinmarketcap5min (
	id SERIAL PRIMARY KEY,
	asset_id character varying(32) NOT NULL,  
	name character varying(32) NOT NULL, 
	symbol character varying(32) NOT NULL, 
	rank integer  NOT NULL, 
	price_usd_first real  NOT NULL, 
	price_usd_last real NOT NULL, 
	price_usd_low real  NOT NULL, 
	price_usd_high  real NOT NULL, 
	price_btc_first real  NOT NULL, 
	price_btc_last real  NOT NULL, 
	price_btc_low real  NOT NULL, 
	price_btc_high real  NOT NULL, 
	price_cny_first real  NOT NULL, 
	price_cny_last real  NOT NULL, 
	price_cny_low real  NOT NULL, 
	price_cny_high real  NOT NULL, 
	last_updated bigint NOT NULL,
	timestamp bigint NOT NULL,
       _group character varying(32) 
); 

CREATE INDEX index_timestamp_coinmarketcap5min ON coinmarketcap5min (timestamp);
CREATE INDEX index_last_updated_coinmarketcap5min ON coinmarketcap5min (last_updated);
CREATE INDEX index_symbol_coinmarketcap5min ON coinmarketcap5min USING hash (symbol);

CREATE TABLE IF NOT EXISTS coinmarketcap10min (
	id SERIAL PRIMARY KEY,
	asset_id character varying(32) NOT NULL,  
	name character varying(32) NOT NULL, 
	symbol character varying(32) NOT NULL, 
	rank integer  NOT NULL, 
	price_usd_first real  NOT NULL, 
	price_usd_last real NOT NULL, 
	price_usd_low real  NOT NULL, 
	price_usd_high  real NOT NULL, 
	price_btc_first real  NOT NULL, 
	price_btc_last real  NOT NULL, 
	price_btc_low real  NOT NULL, 
	price_btc_high real  NOT NULL, 
	price_cny_first real  NOT NULL, 
	price_cny_last real  NOT NULL, 
	price_cny_low real  NOT NULL, 
	price_cny_high real  NOT NULL, 
	last_updated bigint NOT NULL,
	timestamp bigint NOT NULL,
       _group character varying(32) 
); 

CREATE INDEX index_timestamp_coinmarketcap10min ON coinmarketcap10min (timestamp);
CREATE INDEX index_last_updated_coinmarketcap10min ON coinmarketcap10min (last_updated);
CREATE INDEX index_symbol_coinmarketcap10min ON coinmarketcap10min USING hash (symbol);

CREATE TABLE IF NOT EXISTS coinmarketcap30min (
	id SERIAL PRIMARY KEY,
	asset_id character varying(32) NOT NULL,  
	name character varying(32) NOT NULL, 
	symbol character varying(32) NOT NULL, 
	rank integer  NOT NULL, 
	price_usd_first real  NOT NULL, 
	price_usd_last real NOT NULL, 
	price_usd_low real  NOT NULL, 
	price_usd_high  real NOT NULL, 
	price_btc_first real  NOT NULL, 
	price_btc_last real  NOT NULL, 
	price_btc_low real  NOT NULL, 
	price_btc_high real  NOT NULL, 
	price_cny_first real  NOT NULL, 
	price_cny_last real  NOT NULL, 
	price_cny_low real  NOT NULL, 
	price_cny_high real  NOT NULL, 
	last_updated bigint NOT NULL,
	timestamp bigint NOT NULL,
       _group character varying(32) 
); 

CREATE INDEX index_timestamp_coinmarketcap30min ON coinmarketcap30min (timestamp);
CREATE INDEX index_last_updated_coinmarketcap30min ON coinmarketcap30min (last_updated);
CREATE INDEX index_symbol_coinmarketcap30min ON coinmarketcap30min USING hash (symbol);

CREATE TABLE IF NOT EXISTS coinmarketcap15min (
	id SERIAL PRIMARY KEY,
	asset_id character varying(32) NOT NULL,  
	name character varying(32) NOT NULL, 
	symbol character varying(32) NOT NULL, 
	rank integer  NOT NULL, 
	price_usd_first real  NOT NULL, 
	price_usd_last real NOT NULL, 
	price_usd_low real  NOT NULL, 
	price_usd_high  real NOT NULL, 
	price_btc_first real  NOT NULL, 
	price_btc_last real  NOT NULL, 
	price_btc_low real  NOT NULL, 
	price_btc_high real  NOT NULL, 
	price_cny_first real  NOT NULL, 
	price_cny_last real  NOT NULL, 
	price_cny_low real  NOT NULL, 
	price_cny_high real  NOT NULL, 
	last_updated bigint NOT NULL,
	timestamp bigint NOT NULL,
       _group character varying(32) 
); 

CREATE INDEX index_timestamp_coinmarketcap15min ON coinmarketcap15min (timestamp);
CREATE INDEX index_last_updated_coinmarketcap15min ON coinmarketcap15min (last_updated);
CREATE INDEX index_symbol_coinmarketcap15min ON coinmarketcap15min USING hash (symbol);

CREATE TABLE IF NOT EXISTS coinmarketcapday (
	id SERIAL PRIMARY KEY,
	asset_id character varying(32) NOT NULL,  
	name character varying(32) NOT NULL, 
	symbol character varying(32) NOT NULL, 
	rank integer  NOT NULL, 
	price_usd_first real  NOT NULL, 
	price_usd_last real NOT NULL, 
	price_usd_low real  NOT NULL, 
	price_usd_high  real NOT NULL, 
	price_btc_first real  NOT NULL, 
	price_btc_last real  NOT NULL, 
	price_btc_low real  NOT NULL, 
	price_btc_high real  NOT NULL, 
	price_cny_first real  NOT NULL, 
	price_cny_last real  NOT NULL, 
	price_cny_low real  NOT NULL, 
	price_cny_high real  NOT NULL, 
	last_updated bigint NOT NULL,
	timestamp bigint NOT NULL,
       _group character varying(32) 
); 

CREATE INDEX index_timestamp_coinmarketcapday ON coinmarketcapday (timestamp);
CREATE INDEX index_last_updated_coinmarketcapday ON coinmarketcapday (last_updated);
CREATE INDEX index_symbol_coinmarketcapday ON coinmarketcapday USING hash (symbol);

CREATE TABLE IF NOT EXISTS coinmarketcaphour (
	id SERIAL PRIMARY KEY,
	asset_id character varying(32) NOT NULL,  
	name character varying(32) NOT NULL, 
	symbol character varying(32) NOT NULL, 
	rank integer  NOT NULL, 
	price_usd_first real  NOT NULL, 
	price_usd_last real NOT NULL, 
	price_usd_low real  NOT NULL, 
	price_usd_high  real NOT NULL, 
	price_btc_first real  NOT NULL, 
	price_btc_last real  NOT NULL, 
	price_btc_low real  NOT NULL, 
	price_btc_high real  NOT NULL, 
	price_cny_first real  NOT NULL, 
	price_cny_last real  NOT NULL, 
	price_cny_low real  NOT NULL, 
	price_cny_high real  NOT NULL, 
	last_updated bigint NOT NULL,
	timestamp bigint NOT NULL,
       _group character varying(32) 
); 

CREATE INDEX index_timestamp_coinmarketcaphour ON coinmarketcaphour (timestamp);
CREATE INDEX index_last_updated_coinmarketcaphour ON coinmarketcaphour (last_updated);
CREATE INDEX index_symbol_coinmarketcaphour ON coinmarketcaphour USING hash (symbol);

CREATE TABLE IF NOT EXISTS coinmarketcapweek (
	id SERIAL PRIMARY KEY,
	asset_id character varying(32) NOT NULL,  
	name character varying(32) NOT NULL, 
	symbol character varying(32) NOT NULL, 
	rank integer  NOT NULL, 
	price_usd_first real  NOT NULL, 
	price_usd_last real NOT NULL, 
	price_usd_low real  NOT NULL, 
	price_usd_high  real NOT NULL, 
	price_btc_first real  NOT NULL, 
	price_btc_last real  NOT NULL, 
	price_btc_low real  NOT NULL, 
	price_btc_high real  NOT NULL, 
	price_cny_first real  NOT NULL, 
	price_cny_last real  NOT NULL, 
	price_cny_low real  NOT NULL, 
	price_cny_high real  NOT NULL, 
	last_updated bigint NOT NULL,
	timestamp bigint NOT NULL,
       _group character varying(32) 
); 

CREATE INDEX index_timestamp_coinmarketcapweek ON coinmarketcapweek (timestamp);
CREATE INDEX index_last_updated_coinmarketcapweek ON coinmarketcapweek (last_updated);
CREATE INDEX index_symbol_coinmarketcapweek ON coinmarketcapweek USING hash (symbol);


CREATE TABLE IF NOT EXISTS coinmarketcapcurrent (
	id SERIAL PRIMARY KEY,
	asset_id character varying(32) NOT NULL,  
	name character varying(32) NOT NULL, 
	symbol character varying(32) NOT NULL, 
	rank integer  NOT NULL, 
	price_usd_first real  NOT NULL, 
	price_usd_last real NOT NULL, 
	price_usd_low real  NOT NULL, 
	price_usd_high  real NOT NULL, 
	price_btc_first real  NOT NULL, 
	price_btc_last real  NOT NULL, 
	price_btc_low real  NOT NULL, 
	price_btc_high real  NOT NULL, 
	price_cny_first real  NOT NULL, 
	price_cny_last real  NOT NULL, 
	price_cny_low real  NOT NULL, 
	price_cny_high real  NOT NULL, 
	last_updated bigint NOT NULL,
	timestamp bigint NOT NULL,
       _group character varying(32) 
); 

CREATE INDEX index_timestamp_coinmarketcapcurrent ON coinmarketcapcurrent (timestamp);
CREATE INDEX index_last_updated_coinmarketcapcurrent ON coinmarketcapcurrent (last_updated);
CREATE INDEX index_symbol_coinmarketcapcurrent ON coinmarketcapcurrent USING hash (symbol);
CREATE INDEX index_group_coinmarketcapcurrent ON coinmarketcapcurrent USING hash (_group);
