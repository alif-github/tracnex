-- PGOPTIONS="-c client_min_messages=error" PGPASSWORD=postgres psql --single-transaction -U postgres -h localhost -d nexSOFT -f nex_auth_pgres_schema.sql 

-- DROP DATABASE "nexSOFT";

-- CREATE DATABASE "nexSOFT"
--     WITH 
--     OWNER = "nexSOFT.admin"
--     ENCODING = 'UTF8'
--     LC_COLLATE = 'en_AU.UTF-8'
--     LC_CTYPE = 'en_AU.UTF-8'
--     TABLESPACE = pg_default
--     CONNECTION LIMIT = -1;
-- 

\c nexSOFT

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

CREATE SCHEMA IF NOT EXISTS "nex_auth" AUTHORIZATION "nexSOFT.admin";
SET search_path to 'nex_auth';
SET ROLE "nexSOFT.admin";


\i nex_postgres.include.sql
--	RAISE NOTICE '4:returned id from insert = %', _id;
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nex_user_pkey_seq;
CREATE TABLE IF NOT EXISTS nex_user
(
    id              bigint UNIQUE NOT NULL DEFAULT nextval('nex_user_pkey_seq'::regclass),
    status          varchar (32)  NOT NULL default 'active' references nex_reference (code), 
    username        varchar (256) NOT NULL UNIQUE,
    password        varchar (128) NOT NULL,
    first_name      varchar (64) NOT NULL,
    last_name       varchar (64) ,
    email           varchar (128) NOT NULL,
    phone           varchar (32) NOT NULL,
    company_id      bigint ,

    created_by      bigint NOT NULL,
    created_at      TIMESTAMP DEFAULT now (),
    last_updated_by bigint NOT NULL,
    last_updated_at TIMESTAMP DEFAULT now (),
    record_status   varchar (32) NOT NULL default 'active' references nex_reference (code)
                    
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;
INSERT INTO nex_user (username, password, first_name, last_name, email,
                     phone, company_id, created_by, last_updated_by)
       VALUES ('nexSOFT_admin', 'n3x50ftg2', 'nexSOFT', 'admin', 'postgres.admin@nexSOFT.co.id',
               '+62622122220139', 1, 1, 1);
COMMENT ON TABLE nex_user
    IS 'user information details';
CREATE INDEX nex_user_first_name_index  ON nex_user (first_name);
CREATE INDEX nex_user_username_index    ON nex_user (username);
CREATE INDEX nex_user_company_id_index  ON nex_user (company_id);
CREATE INDEX nex_user_email_index       ON nex_user (email);
CREATE INDEX nex_user_phone_index       ON nex_user (phone);

    
    
    
-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nex_contact_pkey_seq;
CREATE TABLE IF NOT EXISTS nex_contact
(
    id              bigint UNIQUE NOT NULL DEFAULT nextval('nex_contact_pkey_seq'::regclass),
    type            varchar (32) NOT NULL references nex_reference (code),
    details         varchar (256) NOT NULL,

    created_by      bigint NOT NULL,
    created_at      TIMESTAMP DEFAULT now (),
    last_updated_by bigint NOT NULL,
    last_updated_at TIMESTAMP DEFAULT now (),
    record_status   varchar (32) NOT NULL default 'active' references nex_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;
COMMENT ON TABLE nex_contact IS 'contacts of type email and phone';
INSERT INTO nex_contact (type, details, created_by, last_updated_by) 
       VALUES          ('phone', '+62622122220139', 1, 1),
                       ('email', 'sales@nexsoft.co.id', 1, 1);
CREATE INDEX nex_contact_type_index      ON nex_contact (type);
CREATE INDEX nex_contact_details_index           ON nex_contact (details);


-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nex_role_pkey_seq;
CREATE TABLE IF NOT EXISTS nex_role
(
    id              bigint UNIQUE NOT NULL DEFAULT nextval('nex_role_pkey_seq'::regclass),
    code            varchar (32)  NOT NULL,
    description     varchar (256) NOT NULL,

    created_by      bigint NOT NULL,
    created_at      TIMESTAMP DEFAULT now (),
    last_updated_by bigint NOT NULL,
    last_updated_at TIMESTAMP DEFAULT now (),
    record_status   varchar (32) NOT NULL default 'active' references nex_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;
COMMENT ON TABLE nex_role IS 'roles within the authorisation';
CREATE INDEX nex_role_code_index           ON nex_role (code);


-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nex_user_role_pkey_seq;
CREATE TABLE IF NOT EXISTS nex_user_role
(
    id              bigint UNIQUE NOT NULL DEFAULT nextval('nex_user_role_pkey_seq'::regclass),
    user_id         bigint NOT NULL references nex_user (id),
    role_id         bigint NOT NULL references nex_role (id),

    created_by      bigint NOT NULL,
    created_at      TIMESTAMP DEFAULT now (),
    last_updated_by bigint NOT NULL,
    last_updated_at TIMESTAMP DEFAULT now (),
    record_status   varchar (32) NOT NULL default 'active' references nex_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;
COMMENT ON TABLE nex_user_role IS 'links a user to role(s)';
CREATE INDEX nex_user_role_user_id_index           ON nex_user_role (user_id);
CREATE INDEX nex_user_role_role_id_index           ON nex_user_role (role_id);


-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nex_address_pkey_seq;
CREATE TABLE IF NOT EXISTS nex_address
(
    id              bigint UNIQUE NOT NULL DEFAULT nextval('nex_address_pkey_seq'::regclass),
    type            varchar (32) references nex_reference (code),
    RT              varchar (64),
    RW              varchar (64),
    kelurahan       varchar (64),
    kecamatan       varchar (64),
    kota            varchar (64),
    kabupaten       varchar (64),
    provinsi        varchar (64),
    country         varchar (64) NOT NULL DEFAULT 'Indonesia',
    country_code    int NOT NULL DEFAULT 62,
    line_address    varchar (512),
    unit_no         varchar (16),
    house_no        varchar (16),
    street_name     varchar (32),
    suburb          varchar (32),
    city            varchar (32),
    state           varchar (64),
    county          varchar (32),
    location        point,
    
    created_at      TIMESTAMP DEFAULT now (),
    last_updated_by bigint NOT NULL,
    last_updated_at TIMESTAMP DEFAULT now (),
    record_status   varchar (32) NOT NULL default 'active' references nex_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;
COMMENT ON TABLE nex_address IS 'address details';
COMMENT ON COLUMN nex_address.line_address IS 'this is the single line address where we do not have broken down address';








-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nex_data_scope_pkey_seq;
CREATE TABLE IF NOT EXISTS    nex_data_scope
(
    id                  bigint UNIQUE NOT NULL DEFAULT nextval('nex_data_scope_pkey_seq'::regclass),
    name                varchar (128) NOT NULL,
    type                varchar (32) NOT NULL,

    created_at          TIMESTAMP DEFAULT now (),
    last_updated_by     bigint NOT NULL,
    last_updated_at     TIMESTAMP DEFAULT now (),
    record_status       varchar (32) NOT NULL default 'active' references nex_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;
CREATE INDEX nex_data_scope_code_index ON nex_data_scope (name);
CREATE INDEX nex_data_scope_type_index ON nex_data_scope (type);




-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nex_data_scope_list_pkey_seq;
CREATE TABLE IF NOT EXISTS    nex_data_scope_list
(
    id                  bigint UNIQUE NOT NULL DEFAULT nextval('nex_data_scope_list_pkey_seq'::regclass),
    data_scope_id       bigint NOT NULL references nex_data_scope (id),
-- ??    data                varchar (32) NOT NULL,

    created_at          TIMESTAMP DEFAULT now (),
    last_updated_by     bigint NOT NULL,
    last_updated_at     TIMESTAMP DEFAULT now (),
    record_status       varchar (32) NOT NULL default 'active' references nex_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;




-- -----------------------------------------------------------------------------
-- -----------------------------------------------------------------------------
-- -----------------------------------------------------------------------------
-- -----------------------------------------------------------------------------
-- -----------------------------------------------------------------------------
-- -----------------------------------------------------------------------------
-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nex_company_pkey_seq;
CREATE TABLE IF NOT EXISTS nex_company
(
    id              bigint UNIQUE NOT NULL DEFAULT nextval('nex_company_pkey_seq'::regclass),
    type            varchar (32) references nex_reference (code),
    name            varchar (256) NOT NULL,
    legal_name      varchar (256),
    company_reg_id  varchar (256),
    website         varchar (256),
    contact_person  varchar (128),  /* name of the person to contact */
    contact_title   varchar (128),  /* title (position) of the contact */
    contact_details varchar (256),  /* email or phone no */
    parent_id       bigint, 

    created_at      TIMESTAMP DEFAULT now (),
    last_updated_by bigint NOT NULL,
    last_updated_at TIMESTAMP DEFAULT now (),
    record_status   varchar (32) NOT NULL default 'active' references nex_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;
COMMENT ON TABLE nex_company IS 'company details';
COMMENT ON COLUMN nex_company.type IS 'type and parent_id help describe the relationship between parent, division, branch, ...';
INSERT INTO nex_company (type, name, legal_name, company_reg_id, website, contact_person,
                        contact_title, contact_details, parent_id, last_updated_by)
          VALUES       ('parent', 'nexSOFT', '', '', 'nexsoft.co.id', 'Alex H Wreksoremboko',
                        'CEO', 'awreksoremboko@nexsoft.co.id', null, '1');



-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nex_license_pkey_seq;
CREATE TABLE IF NOT EXISTS    nex_license
(
    id                  bigint UNIQUE NOT NULL DEFAULT nextval('nex_license_pkey_seq'::regclass),
    type                varchar (32) references nex_reference (code),
    issuing_authority   varchar (64) NOT NULL,
                        
    created_at          TIMESTAMP DEFAULT now (),
    last_updated_by     bigint NOT NULL,
    last_updated_at     TIMESTAMP DEFAULT now (),
    record_status       varchar (32) NOT NULL default 'active' references nex_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;
COMMENT ON TABLE nex_license IS 'license details';
CREATE INDEX nex_license_type_index ON nex_license (type);



-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nex_company_license_pkey_seq;
CREATE TABLE IF NOT EXISTS    nex_company_license
(
    id                  bigint UNIQUE NOT NULL DEFAULT nextval('nex_company_license_pkey_seq'::regclass),
    company_id          bigint NOT NULL references nex_company (id),
    license_id          bigint NOT NULL references nex_license (id),
    code                varchar (32) NOT NULL,
    valid_from          date NOT NULL,
    expiry_date         date NOT NULL,

    created_at          TIMESTAMP DEFAULT now (),
    last_updated_by     bigint NOT NULL,
    last_updated_at     TIMESTAMP DEFAULT now (),
    record_status       varchar (32) NOT NULL default 'active' references nex_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;
COMMENT ON TABLE nex_company_license IS 'links license to company or division or branch';
COMMENT ON COLUMN nex_company_license.code IS 'This can be used as a license code that the app sends to be validated from the server side';
CREATE INDEX nex_company_license_code_index ON nex_company_license (code);
CREATE INDEX nex_company_license_company_id_index   ON nex_company_license (company_id);



-- -----------------------------------------------------------------------------
-- -----------------------------------------------------------------------------
-- -----------------------------------------------------------------------------




-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- Stored Procedures and functions
-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- ----------------------------------------------------------

-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- user functions
-- ----------------------------------------------------------
-- ----------------------------------------------------------

-- ----------------------------------------------------------
-- return id for a username. if the user is not found then 
-- it returns -1
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION get_user_id (username_   varchar)
RETURNS bigint 
AS $BODY$
DECLARE   _id    bigint;
BEGIN        
    _id := -1;
    SELECT id INTO _id FROM nex_user WHERE username = username_;
    IF NOT FOUND THEN 
        _id = -1;
    ELSEIF (_id IS NULL) THEN
		_id := -1;
	END IF;
    RETURN _id;
END; $BODY$
LANGUAGE PLPGSQL;

-- ----------------------------------------------------------
-- validate user information
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION validate_nex_user_info (IN   username_   varchar,  
                                                   IN   password_   varchar,   
                                                   IN first_name_   varchar, 
                                                   IN  last_name_   varchar, 
                                                   IN      email_   varchar, 
                                                   IN      phone_   varchar, 
                                                   IN company_id_   bigint, 
                                                   IN created_by_   bigint)
RETURNS rc_type AS $BODY$
DECLARE   _rc   rc_type;
BEGIN
	_rc.code    = 0; 
	_rc.message = '';
    IF     (username_   = null OR username_   = '') THEN _rc.code := '-1'; _rc.message := 'username can not be null or blank';
    ELSEIF (password_   = null OR password_   = '') THEN _rc.code := '-2'; _rc.message := 'password can not be null or blank';
    ELSEIF (first_name_ = null OR first_name_ = '') THEN _rc.code := '-3'; _rc.message := 'first_name_ can not be null or blank';
    ELSEIF (email_      = null OR email_      = '') THEN _rc.code := '-4'; _rc.message := 'email can not be null or blank';
    ELSEIF (phone_      = null OR phone_      = '') THEN _rc.code := '-5'; _rc.message := 'phone_ can not be null or blank';
    ELSEIF (created_by_ < 1)                        THEN _rc.code := '-6'; _rc.message := 'Invalid user for created_by column';
    ELSEIF (company_id_ < 1)                        THEN _rc.code := '-7'; _rc.message := 'Invalid company ID';
	END IF;
	RETURN _rc;
END; $BODY$
LANGUAGE PLPGSQL;


-- ----------------------------------------------------------
-- insert into user table 
-- insert_nex_user ('nexSOFT_admin', 'password', 'foad', 'momtazi', 'foad.momtazi@nexsoft.co.id', '+62 666 666 666 666', 1,1)
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION insert_nex_user (IN    username_   varchar,  
                                            IN    password_   varchar,   
                                            IN  first_name_   varchar, 
                                            IN   last_name_   varchar, 
                                            IN       email_   varchar, 
                                            IN       phone_   varchar, 
                                            IN  company_id_   bigint, 
                                            IN  created_by_   bigint)
RETURNS bigint AS $BODY$
DECLARE 
    _id 	bigint;
BEGIN
	INSERT INTO nex_user (username, password, first_name, last_name, email, phone, company_id, created_by, last_updated_by)
		 VALUES (username_, password_, first_name_, last_name_, email_, phone_, company_id_, created_by_, created_by_)
	RETURNING id INTO _id;

	RETURN _id;
END; $BODY$
LANGUAGE PLPGSQL;


-- ----------------------------------------------------------
-- add user 
-- add_nex_user ('nexSOFT_admin', 'password', 'foad', 'momtazi', 'foad.momtazi@nexsoft.co.id', '+62 666 666 666 666', 1,1)
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION add_nex_user (IN     username_   varchar,  
                                         IN     password_   varchar,   
                                         IN   first_name_   varchar, 
                                         IN    last_name_   varchar, 
                                         IN        email_   varchar, 
                                         IN        phone_   varchar, 
                                         IN   company_id_   bigint, 
                                         IN   created_by_   bigint)
RETURNS text AS $BODY$
DECLARE 
    _validate_response   rc_type;
    _id                  bigint;
    _result             json;
BEGIN
    _validate_response := validate_nex_user_info (username_, password_, first_name_, last_name_, email_, phone_, company_id_, created_by_) ;
    IF _validate_response.rc_code != 0 THEN
        RETURN create_query_rc_json_type (tochar (_validate_response.rc_code, '999'), _validate_response.rc_message, NULL);
    END IF;
    
  	_id = get_user_id (username_);
    IF (_id > 0) THEN
        RETURN create_query_rc_json_type ('-9', 'username is already in use', null);
    ELSE 
        _id := insert_nex_user (username_, password_, first_name_, 
                                  last_name_, email_, phone_, company_id_, created_by_);
       _result = (SELECT (row_to_json(root_node))
                  FROM (SELECT _id AS id) root_node );
       RETURN create_query_rc_json_type ('0', '', _result);
    END IF;
    
END; $BODY$
LANGUAGE PLPGSQL;



-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- return id for a role_code or role_name. if the role is not found then 
-- it returns -1
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION get_role_id (role_code_   varchar)
    RETURNS bigint 
AS $BODY$
DECLARE   _id    bigint;
BEGIN        
    _id := -1;
    SELECT id INTO _id FROM nex_role WHERE code = role_code_;
	IF (_id IS NULL) THEN
		_id := -1;
	END IF;
    RETURN _id;
END; $BODY$
LANGUAGE PLPGSQL;



DROP FUNCTION IF EXISTS insert_nex_role ();
-- ----------------------------------------------------------
-- create a new role
-- select create_role ()
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION insert_nex_role (IN    role_code_   varchar,
                                           IN  description_   varchar,
                                           IN   created_by_   bigint)
RETURNS bigint AS $BODY$
DECLARE 
    _id                  bigint;
BEGIN
    INSERT INTO nex_role (code, description, created_by, last_updated_by)
            VALUES (role_code_, description_, created_by_, created_by_)
            RETURNING id INTO _id;
    RETURN _id;
END; $BODY$
LANGUAGE plpgSQL;




DROP FUNCTION IF EXISTS create_role (varchar, varchar, bigint);
-- ----------------------------------------------------------
-- create a new role
-- select create_role ()
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION create_role (IN    role_code_   varchar,
                                        IN  description_   varchar,
                                        IN   created_by_   bigint)
RETURNS json AS $BODY$
DECLARE 
    v_id                bigint;
    v_result            json;
BEGIN
    v_id = get_role_id (role_code_);
    IF (v_id > 0) THEN
        RETURN create_query_rc_json_type ('-19', 'role code already exists', null);
    ELSE
        v_id = insert_nex_role (role_code_, description_, created_by_);
        v_result = (SELECT (row_to_json(root_node))
                    FROM (SELECT v_id AS id ) root_node );
        RETURN create_query_rc_json_type ('0', '', v_result);
    END IF;

END; $BODY$
LANGUAGE plpgSQL;




DROP FUNCTION IF EXISTS get_user_info ();
-- ----------------------------------------------------------
-- get all the references, sorted by type and code 
-- select * from get_all_references ()
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION get_user_info (IN username_  varchar)
RETURNS JSON AS $$
DECLARE 
    _id         bigint;
BEGIN
    _id := get_user_id (username_);
    IF _id < 1 THEN
        RETURN rc_as_json (-9, 'user not found');
    ELSE
        RETURN (SELECT array_to_json(array_agg(t ))
        from (
                SELECT id, type, code FROM nex_reference ORDER BY type, code
             )  t ) ;
    END IF;

END; $$
LANGUAGE plpgSQL;



DROP FUNCTION IF EXISTS get_roles_for_all_users ();
-- ----------------------------------------------------------
-- function RETURNs all users' role(s)
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION get_roles_for_all_users ()
    RETURNS json  AS 
$BODY$
DECLARE
    _result     json;
BEGIN

    _result = (SELECT array_to_json(array_agg(row_to_json(root_node))) AS user_roles 
        FROM (
                SELECT u.username, u.id AS user_id,
                (SELECT array_to_json(array_agg(row_to_json(role))) AS roles 
                FROM (
                        SELECT  ur.role_id, r.code
                          FROM  nex_role r
                         INNER  JOIN nex_user_role ur  ON (u.id = ur.user_id)
                         WHERE  ur.role_id = r.id
                              
                     ) role
                 )  FROM nex_user u
             ) root_node 
        );
        RETURN create_query_rc_json_type ('0', '', _result);
END; 
$BODY$
LANGUAGE 'plpgsql';



DROP FUNCTION IF EXISTS get_roles_from_username(varchar);
-- ----------------------------------------------------------
-- function RETURNs all roles for a username
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION get_roles_from_username (IN username_ varchar)
    RETURNS json  AS 
$BODY$

DECLARE     _id    bigint;
BEGIN

    _id = get_user_id (username_);
    IF _id = -1 THEN
        RETURN create_query_rc_json_type ('E87764', 'username does not exists', null);
    END IF;

    RETURN get_roles_from_userid (_id);

END; 
$BODY$
LANGUAGE 'plpgsql';
    



DROP FUNCTION IF EXISTS get_roles_from_userid(bigint);
-- ----------------------------------------------------------
-- function RETURNs all roles for a username
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION get_roles_from_userid (IN user_id_ bigint)
    RETURNS json  AS 
$BODY$
DECLARE     _result	    json;
            v_user_id    bigint = -1;
            v_username   varchar;
BEGIN

    SELECT u.id, u.username INTO v_user_id, v_username FROM nex_user u WHERE id=user_id_;
    IF NOT FOUND THEN
        RETURN create_query_rc_json_type ('E87764', 'username does not exists', null);
    END IF;
    _result = (SELECT (row_to_json(root_node)) AS user_roles 
    FROM (
--            SELECT u.username AS username, u.id AS user_id,
            SELECT v_username AS username, v_user_id AS user_id,
            (SELECT array_to_json(array_agg(row_to_json(role))) AS roles 
            FROM (
                    SELECT  ur.role_id, r.code
                      FROM  nex_role r, nex_user u
                     INNER  JOIN nex_user_role ur  ON (u.id = ur.user_id)
                     WHERE  ur.role_id = r.id
                          
                 ) role
             )  -- FROM nex_user u WHERE u.id = user_id_
         ) root_node 
    );
    RETURN create_query_rc_json_type ('0', '', _result);
END; 
$BODY$
LANGUAGE 'plpgsql';
    


-- ----------------------------------------------------------
-- test data
-- ----------------------------------------------------------
INSERT INTO nex_user (username, password, first_name, last_name, email, phone, created_by, last_updated_by) VALUES ('test_user_1', 'password', 'test', 'user 1', 'test@user1.com', '666 6666 6666', 1,1); 
INSERT INTO nex_user (username, password, first_name, last_name, email, phone, created_by, last_updated_by) VALUES ('test_user_2', 'password', 'test', 'user 2', 'test@user1.com', '666 6666 6666', 1,1); 

INSERT INTO nex_role (code,description,created_by,last_updated_by) VALUES ('read_only', 'simple user role', 1,1);
INSERT INTO nex_role (code,description,created_by,last_updated_by) VALUES ('manager',   'manager role', 1,1);
INSERT INTO nex_role (code,description,created_by,last_updated_by) VALUES ('admin',     'system admin', 1,1);

INSERT INTO nex_user_role (user_id,role_id,created_by,last_updated_by) VALUES (1,3,1,1);
INSERT INTO nex_user_role (user_id,role_id,created_by,last_updated_by) VALUES (2,1,1,1);
INSERT INTO nex_user_role (user_id,role_id,created_by,last_updated_by) VALUES (2,2,1,1);







