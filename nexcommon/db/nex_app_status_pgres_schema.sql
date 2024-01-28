-- PGOPTIONS="-c client_min_messages=error" PGPASSWORD=postgres psql --single-transaction -U postgres -h localhost -d postgres -f nex_app_status_pgres_schema.sql 

SET ROLE "postgres";

DO $BODY$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'nexSOFT.admin') THEN
        CREATE ROLE "nexSOFT.admin" WITH
              LOGIN
              SUPERUSER
              INHERIT
              CREATEDB
              CREATEROLE
              NOREPLICATION;
         GRANT postgres TO "nexSOFT.admin";
     END IF;
END;
$BODY$;

CREATE SCHEMA IF NOT EXISTS "nex_app_status" AUTHORIZATION "nexSOFT.admin";
SET search_path to 'nex_app_status';
SET ROLE "nexSOFT.admin";

DROP TYPE IF EXISTS protocol CASCADE;
CREATE TYPE protocol AS ENUM ('get', 'post');

DROP TYPE IF EXISTS health_status CASCADE;
CREATE TYPE health_status AS ENUM ('green', 'amber', 'red');


DROP TABLE IF EXISTS nex_service CASCADE;
DROP SEQUENCE IF EXISTS nex_service_pkey_seq;
-- -----------------------------------------------------------------------------
-- nex_service table
-- this describes a service details for the purpose of health check
-- -----------------------------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nex_service_pkey_seq;
CREATE TABLE IF NOT EXISTS nex_service
(
    id                  bigint UNIQUE NOT NULL DEFAULT nextval('nex_service_pkey_seq'::regclass),
    name                varchar (64)  NOT NULL,  /* could be product group or module */
    category            varchar (32),
    description         varchar (256),
    parent              bigint,             /* describes the service dependency */
    health_check_mode   protocol,           /* get or post */
    health_check_url    varchar (512),      /* the url to call the get the health status back */
    health_ckeck_data   varchar (2048),     /* data for http post */
    
    created_by          bigint NOT NULL,
    created_at          TIMESTAMP DEFAULT now (),
    last_updated_by     bigint NOT NULL,
    last_updated_at     TIMESTAMP DEFAULT now (),
    record_status       varchar (32) NOT NULL default 'active'
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;
INSERT INTO nex_service (name, category, description, parent, health_check_mode, health_check_url, health_ckeck_data, created_by, last_updated_by)
       VALUES ('root', 'root', 'the root of the tree. it is not a service', 0, 'post', 'https://nexSOFT.co.id', '', 1, 1);


DROP TABLE IF EXISTS nex_service_status CASCADE;
DROP SEQUENCE IF EXISTS nex_service_status_pkey_seq;
-- -----------------------------------------------------------------------------
-- nex_service_status table
-- this describes a service health check response 
-- -----------------------------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nex_service_status_pkey_seq;
CREATE TABLE IF NOT EXISTS nex_service_status
(
    id                  bigint UNIQUE NOT NULL DEFAULT nextval('nex_service_status_pkey_seq'::regclass),
    service_id          bigint NOT NULL references nex_service (id),
    status              health_status,
    response_message    varchar (2048),
    
    created_at          TIMESTAMP DEFAULT now ()
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;



DROP TABLE IF EXISTS nex_txn_log CASCADE;
DROP SEQUENCE IF EXISTS nex_txn_log_pkey_seq;
-- -----------------------------------------------------------------------------
-- nex_txn_log table
-- for request response times and SLA tracking 
-- -----------------------------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nex_txn_log_pkey_seq;
CREATE TABLE IF NOT EXISTS nex_txn_log
(
    id                  bigint UNIQUE NOT NULL DEFAULT nextval('nex_txn_log_pkey_seq'::regclass),
    txn_id              varchar (64),       /* can be UUID for the request or similar unique identifier */
    type                varchar (64),
    client              varchar (64),
    source_system       varchar (64),
    destination         varchar (64),
    response_code       varchar (64),       /* success or error code */
    start_timestamp     TIMESTAMP,
    end_timestamp       TIMESTAMP,
    
    created_at          TIMESTAMP DEFAULT now ()
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;


DROP TABLE IF EXISTS nex_txn_audit CASCADE;
DROP SEQUENCE IF EXISTS nex_txn_audit_pkey_seq;
-- -----------------------------------------------------------------------------
-- nex_txn_audit table
-- for request response audit logging 
-- -----------------------------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nex_txn_audit_pkey_seq;
CREATE TABLE IF NOT EXISTS nex_txn_audit
(
    id                  bigint UNIQUE NOT NULL DEFAULT nextval('nex_txn_audit_pkey_seq'::regclass),
    txn_id              varchar (64),       /* can be UUID for the request or similar unique identifier */
    type                varchar (64),
    user_id             bigint,
    request             varchar (512),
    response            varchar (2048),
    
    created_at          TIMESTAMP DEFAULT now ()
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;



-- -----------------------------------------------------------------------------
-- -----------------------------------------------------------------------------
-- -----------------------------------------------------------------------------
-- -----------------------------------------------------------------------------
-- -----------------------------------------------------------------------------
-- -----------------------------------------------------------------------------
-- -----------------------------------------------------------------------------
DROP FUNCTION IF EXISTS save_service_status;
-- ----------------------------------------------------------
-- return status of a single service 
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION save_service_status (IN        service_id_  bigint, 
                                                IN            status_  health_status,
                                                IN  response_message_  varchar)
RETURNS void 
AS $BODY$
BEGIN
    INSERT INTO nex_service_status (service_id, status, response_message)
        VALUES   (service_id_, status_, response_message_);
END; $BODY$
LANGUAGE PLPGSQL;


