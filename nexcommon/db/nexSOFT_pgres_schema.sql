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

CREATE SCHEMA IF NOT EXISTS "nexSOFT" AUTHORIZATION "nexSOFT.admin";
SET search_path to 'nexSOFT';
SET ROLE "nexSOFT.admin";


-- ----------------------------------------------------
-- dropping all the tables
-- ----------------------------------------------------
DO $BODY$
DECLARE
    DROP_TABLES bool = true;
   
BEGIN
    IF DROP_TABLES = true THEN
        DROP TABLE    IF EXISTS nx_contact;
        DROP SEQUENCE IF EXISTS nx_contact_pkey_seq;
        
        DROP TABLE    IF EXISTS nx_user_role;
        DROP SEQUENCE IF EXISTS nx_user_role_pkey_seq;
        
        DROP TABLE    IF EXISTS nx_role;
        DROP SEQUENCE IF EXISTS nx_role_pkey_seq;
        
        DROP TABLE    IF EXISTS nx_user;
        DROP SEQUENCE IF EXISTS nx_user_pkey_seq;

        DROP TABLE    IF EXISTS nx_data_scope_list;
        DROP SEQUENCE IF EXISTS nx_data_scope_list_pkey_seq;

        DROP TABLE    IF EXISTS nx_data_scope;
        DROP SEQUENCE IF EXISTS nx_data_scope_pkey_seq;

        DROP TABLE    IF EXISTS nx_company_license;
        DROP SEQUENCE IF EXISTS nx_company_license_pkey_seq;
        
        DROP TABLE    IF EXISTS nx_company;
        DROP SEQUENCE IF EXISTS nx_company_pkey_seq;
        
        DROP TABLE    IF EXISTS nx_license;
        DROP SEQUENCE IF EXISTS nx_license_pkey_seq;

        DROP TABLE    IF EXISTS nx_address;
        DROP SEQUENCE IF EXISTS nx_address_pkey_seq;
        
        DROP TABLE    IF EXISTS nx_country;
        DROP SEQUENCE IF EXISTS nx_country_pkey_seq;
        
        DROP TABLE    IF EXISTS nx_reference;
        DROP SEQUENCE IF EXISTS nx_reference_pkey_seq;
        
    END IF;
END;
$BODY$;

-- --------------------------------------------------------


-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nx_reference_pkey_seq;
CREATE TABLE IF NOT EXISTS nx_reference
(
    id              bigint UNIQUE NOT NULL DEFAULT nextval('nx_reference_pkey_seq'::regclass),
    type            varchar(32) NOT NULL,
    code            varchar(32) NOT NULL UNIQUE
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;
CREATE INDEX nx_reference_type_index ON nx_reference (type);
CREATE INDEX nx_reference_code_index ON nx_reference (code);
CREATE INDEX nx_reference_combined_index  ON nx_reference (type, code);

INSERT INTO nx_reference (type, code) VALUES ('status',  'active');
INSERT INTO nx_reference (type, code) VALUES ('status',  'deleted');
INSERT INTO nx_reference (type, code) VALUES ('status',  'locked');
INSERT INTO nx_reference (type, code) VALUES ('status',  'change_password');
INSERT INTO nx_reference (type, code) VALUES ('contact', 'email');
INSERT INTO nx_reference (type, code) VALUES ('contact', 'phone');
INSERT INTO nx_reference (type, code) VALUES ('company', 'parent');
INSERT INTO nx_reference (type, code) VALUES ('company', 'division');
INSERT INTO nx_reference (type, code) VALUES ('company', 'branch');
INSERT INTO nx_reference (type, code) VALUES ('company', 'subsidiary');
INSERT INTO nx_reference (type, code) VALUES ('address', 'physical');
INSERT INTO nx_reference (type, code) VALUES ('address', 'postal');
INSERT INTO nx_reference (type, code) VALUES ('address', 'delivery');
INSERT INTO nx_reference (type, code) VALUES ('address', 'POBox');


-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nx_country_pkey_seq;
CREATE TABLE IF NOT EXISTS nx_country
(
    id              bigint UNIQUE NOT NULL DEFAULT nextval('nx_country_pkey_seq'::regclass),
    name            varchar(128) NOT NULL,
    A2_code         varchar (4) NOT NULL,
    A3_code         varchar (6) NOT NULL,
    un_code         varchar (6) NOT NULL,
    cc_code         varchar (6) /* cc=Country Calling */
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;
CREATE INDEX nx_country_name_index     ON nx_country (name);
CREATE INDEX nx_country_A2_code_index  ON nx_country (A2_code);
CREATE INDEX nx_country_A3_code_index  ON nx_country (A3_code);
CREATE INDEX nx_country_un_code_index  ON nx_country (un_code);
CREATE INDEX nx_country_cc_code_index  ON nx_country (cc_code);



-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nx_user_pkey_seq;
CREATE TABLE IF NOT EXISTS nx_user
(
    id              bigint UNIQUE NOT NULL DEFAULT nextval('nx_user_pkey_seq'::regclass),
    status          varchar (32)  NOT NULL default 'active' references nx_reference (code), 
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
    record_status   varchar (32) NOT NULL default 'active' references nx_reference (code)
                    
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;
INSERT INTO nx_user (username, password, first_name, last_name, email,
                     phone, company_id, created_by, last_updated_by)
       VALUES ('nexSOFT_admin', 'n3x50ftg2', 'nexSOFT', 'admin', 'postgres.admin@nexSOFT.co.id',
               '+62622122220139', 1, 1, 1);
COMMENT ON TABLE nx_user
    IS 'user information details';
CREATE INDEX nx_user_first_name_index  ON nx_user (first_name);
CREATE INDEX nx_user_username_index    ON nx_user (username);
CREATE INDEX nx_user_company_id_index  ON nx_user (company_id);
CREATE INDEX nx_user_email_index       ON nx_user (email);
CREATE INDEX nx_user_phone_index       ON nx_user (phone);

    
    
    
-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nx_contact_pkey_seq;
CREATE TABLE IF NOT EXISTS nx_contact
(
    id              bigint UNIQUE NOT NULL DEFAULT nextval('nx_contact_pkey_seq'::regclass),
    type            varchar (32) NOT NULL references nx_reference (code),
    details         varchar (256) NOT NULL,

    created_by      bigint NOT NULL,
    created_at      TIMESTAMP DEFAULT now (),
    last_updated_by bigint NOT NULL,
    last_updated_at TIMESTAMP DEFAULT now (),
    record_status   varchar (32) NOT NULL default 'active' references nx_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;
COMMENT ON TABLE nx_contact IS 'contacts of type email and phone';
INSERT INTO nx_contact (type, details, created_by, last_updated_by) 
       VALUES          ('phone', '+62622122220139', 1, 1),
                       ('email', 'sales@nexsoft.co.id', 1, 1);
CREATE INDEX nx_contact_type_index      ON nx_contact (type);
CREATE INDEX nx_contact_details_index           ON nx_contact (details);


-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nx_role_pkey_seq;
CREATE TABLE IF NOT EXISTS nx_role
(
    id              bigint UNIQUE NOT NULL DEFAULT nextval('nx_role_pkey_seq'::regclass),
    code            varchar (32)  NOT NULL references nx_reference (code),
    description     varchar (256) NOT NULL,

    created_by      bigint NOT NULL,
    created_at      TIMESTAMP DEFAULT now (),
    last_updated_by bigint NOT NULL,
    last_updated_at TIMESTAMP DEFAULT now (),
    record_status   varchar (32) NOT NULL default 'active' references nx_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;
COMMENT ON TABLE nx_role IS 'roles within the authorisation';
CREATE INDEX nx_role_code_index           ON nx_role (code);


-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nx_user_role_pkey_seq;
CREATE TABLE IF NOT EXISTS nx_user_role
(
    id              bigint UNIQUE NOT NULL DEFAULT nextval('nx_user_role_pkey_seq'::regclass),
    user_id         bigint NOT NULL references nx_user (id),
    role_id         bigint NOT NULL references nx_role (id),

    created_by      bigint NOT NULL,
    created_at      TIMESTAMP DEFAULT now (),
    last_updated_by bigint NOT NULL,
    last_updated_at TIMESTAMP DEFAULT now (),
    record_status   varchar (32) NOT NULL default 'active' references nx_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;
COMMENT ON TABLE nx_user_role IS 'links a user to role(s)';
CREATE INDEX nx_user_role_user_id_index           ON nx_user_role (user_id);
CREATE INDEX nx_user_role_role_id_index           ON nx_user_role (role_id);


-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nx_data_scope_pkey_seq;
CREATE TABLE IF NOT EXISTS    nx_data_scope
(
    id                  bigint UNIQUE NOT NULL DEFAULT nextval('nx_data_scope_pkey_seq'::regclass),
    name                varchar (128) NOT NULL,
    type                varchar (32) NOT NULL,

    created_at          TIMESTAMP DEFAULT now (),
    last_updated_by     bigint NOT NULL,
    last_updated_at     TIMESTAMP DEFAULT now (),
    record_status       varchar (32) NOT NULL default 'active' references nx_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;
CREATE INDEX nx_data_scope_code_index ON nx_data_scope (name);
CREATE INDEX nx_data_scope_type_index ON nx_data_scope (type);




-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nx_data_scope_list_pkey_seq;
CREATE TABLE IF NOT EXISTS    nx_data_scope_list
(
    id                  bigint UNIQUE NOT NULL DEFAULT nextval('nx_data_scope_list_pkey_seq'::regclass),
    data_scope_id       bigint NOT NULL references nx_data_scope (id),
-- ??    data                varchar (32) NOT NULL,

    created_at          TIMESTAMP DEFAULT now (),
    last_updated_by     bigint NOT NULL,
    last_updated_at     TIMESTAMP DEFAULT now (),
    record_status       varchar (32) NOT NULL default 'active' references nx_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;





-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nx_company_pkey_seq;
CREATE TABLE IF NOT EXISTS nx_company
(
    id              bigint UNIQUE NOT NULL DEFAULT nextval('nx_company_pkey_seq'::regclass),
    type            varchar (32) references nx_reference (code),
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
    record_status   varchar (32) NOT NULL default 'active' references nx_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;
COMMENT ON TABLE nx_company IS 'company details';
COMMENT ON COLUMN nx_company.type IS 'type and parent_id help describe the relationship between parent, division, branch, ...';
INSERT INTO nx_company (type, name, legal_name, company_reg_id, website, contact_person,
                        contact_title, contact_details, parent_id, last_updated_by)
          VALUES       ('parent', 'nexSOFT', '', '', 'nexsoft.co.id', 'Alex H Wreksoremboko',
                        'CEO', 'awreksoremboko@nexsoft.co.id', null, '1');



-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nx_address_pkey_seq;
CREATE TABLE IF NOT EXISTS nx_address
(
    id              bigint UNIQUE NOT NULL DEFAULT nextval('nx_address_pkey_seq'::regclass),
    type            varchar (32) references nx_reference (code),
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
    record_status   varchar (32) NOT NULL default 'active' references nx_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;
COMMENT ON TABLE nx_address IS 'address details';
COMMENT ON COLUMN nx_address.line_address IS 'this is the single line address where we do not have broken down address';


-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nx_license_pkey_seq;
CREATE TABLE IF NOT EXISTS    nx_license
(
    id                  bigint UNIQUE NOT NULL DEFAULT nextval('nx_license_pkey_seq'::regclass),
    type                varchar (32) references nx_reference (code),
    issuing_authority   varchar (64) NOT NULL,
                        
    created_at          TIMESTAMP DEFAULT now (),
    last_updated_by     bigint NOT NULL,
    last_updated_at     TIMESTAMP DEFAULT now (),
    record_status       varchar (32) NOT NULL default 'active' references nx_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;
COMMENT ON TABLE nx_license IS 'license details';
CREATE INDEX nx_license_type_index ON nx_license (type);



-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nx_company_license_pkey_seq;
CREATE TABLE IF NOT EXISTS    nx_company_license
(
    id                  bigint UNIQUE NOT NULL DEFAULT nextval('nx_company_license_pkey_seq'::regclass),
    company_id          bigint NOT NULL references nx_company (id),
    license_id          bigint NOT NULL references nx_license (id),
    code                varchar (32) NOT NULL,
    valid_from          date NOT NULL,
    expiry_date         date NOT NULL,

    created_at          TIMESTAMP DEFAULT now (),
    last_updated_by     bigint NOT NULL,
    last_updated_at     TIMESTAMP DEFAULT now (),
    record_status       varchar (32) NOT NULL default 'active' references nx_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;
COMMENT ON TABLE nx_company_license IS 'links license to company or division or branch';
COMMENT ON COLUMN nx_company_license.code IS 'This can be used as a license code that the app sends to be validated from the server side';
CREATE INDEX nx_company_license_code_index ON nx_company_license (code);
CREATE INDEX nx_company_license_company_id_index   ON nx_company_license (company_id);




-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- insert reference data
-- ----------------------------------------------------------
-- country
INSERT INTO nx_country (name, A2_code, A3_code, un_code)  
       VALUES 
-- ('Afghanistan',  'AF', 'AFG', '004'),
-- ('ALA Aland Islands',  'AX', 'ALA', '248'),
-- ('Albania',  'AL', 'ALB', '008'),
-- ('Algeria',  'DZ', 'DZA', '012'),
-- ('American Samoa',  'AS', 'ASM', '016'),
-- ('Andorra',  'AD', 'AND', '020'),
-- ('Angola',  'AO', 'AGO', '024'),
-- ('Anguilla',  'AI', 'AIA', '660'),
-- ('Antarctica',  'AQ', 'ATA', '010'),
-- ('Antigua and Barbuda',  'AG', 'ATG', '028'),
-- ('Argentina',  'AR', 'ARG', '032'),
-- ('Armenia',  'AM', 'ARM', '051'),
-- ('Aruba',  'AW', 'ABW', '533'),
('Australia',  'AU', 'AUS', '036'),
-- ('Austria',  'AT', 'AUT', '040'),
-- ('Azerbaijan',  'AZ', 'AZE', '031'),
-- ('Bahamas',  'BS', 'BHS', '044'),
-- ('Bahrain',  'BH', 'BHR', '048'),
-- ('Bangladesh',  'BD', 'BGD', '050'),
-- ('Barbados',  'BB', 'BRB', '052'),
-- ('Belarus',  'BY', 'BLR', '112'),
-- ('Belgium',  'BE', 'BEL', '056'),
-- ('Belize',  'BZ', 'BLZ', '084'),
-- ('Benin',  'BJ', 'BEN', '204'),
-- ('Bermuda',  'BM', 'BMU', '060'),
-- ('Bhutan',  'BT', 'BTN', '064'),
-- ('Bolivia',  'BO', 'BOL', '068'),
-- ('Bosnia and Herzegovina',  'BA', 'BIH', '070'),
-- ('Botswana',  'BW', 'BWA', '072'),
-- ('Bouvet Island',  'BV', 'BVT', '074'),
-- ('Brazil',  'BR', 'BRA', '076'),
-- ('British Virgin Islands',  'VG', 'VGB', '092'),
-- ('British Indian Ocean Territory',  'IO', 'IOT', '086'),
-- ('Brunei Darussalam',  'BN', 'BRN', '096'),
-- ('Bulgaria',  'BG', 'BGR', '100'),
-- ('Burkina Faso',  'BF', 'BFA', '854'),
-- ('Burundi',  'BI', 'BDI', '108'),
-- ('Cambodia',  'KH', 'KHM', '116'),
-- ('Cameroon',  'CM', 'CMR', '120'),
-- ('Canada',  'CA', 'CAN', '124'),
-- ('Cape Verde',  'CV', 'CPV', '132'),
-- ('Cayman Islands',  'KY', 'CYM', '136'),
-- ('Central African Republic',  'CF', 'CAF', '140'),
-- ('Chad',  'TD', 'TCD', '148'),
-- ('Chile',  'CL', 'CHL', '152'),
-- ('China',  'CN', 'CHN', '156'),
-- ('Hong Kong, SAR China',  'HK', 'HKG', '344'),
-- ('Macao, SAR China',  'MO', 'MAC', '446'),
-- ('Christmas Island',  'CX', 'CXR', '162'),
-- ('Cocos (Keeling) Islands',  'CC', 'CCK', '166'),
-- ('Colombia',  'CO', 'COL', '170'),
-- ('Comoros',  'KM', 'COM', '174'),
-- ('Congo (Brazzaville)',  'CG', 'COG', '178'),
-- ('Congo, (Kinshasa)',  'CD', 'COD', '180'),
-- ('Cook Islands',  'CK', 'COK', '184'),
-- ('Costa Rica',  'CR', 'CRI', '188'),
-- ('Côte d''Ivoire',  'CI', 'CIV', '384'),
-- ('Croatia',  'HR', 'HRV', '191'),
-- ('Cuba',  'CU', 'CUB', '192'),
-- ('Cyprus',  'CY', 'CYP', '196'),
-- ('Czech Republic',  'CZ', 'CZE', '203'),
-- ('Denmark',  'DK', 'DNK', '208'),
-- ('Djibouti',  'DJ', 'DJI', '262'),
-- ('Dominica',  'DM', 'DMA', '212'),
-- ('Dominican Republic',  'DO', 'DOM', '214'),
-- ('Ecuador',  'EC', 'ECU', '218'),
-- ('Egypt',  'EG', 'EGY', '818'),
-- ('El Salvador',  'SV', 'SLV', '222'),
-- ('Equatorial Guinea',  'GQ', 'GNQ', '226'),
-- ('Eritrea',  'ER', 'ERI', '232'),
-- ('Estonia',  'EE', 'EST', '233'),
-- ('Ethiopia',  'ET', 'ETH', '231'),
-- ('Falkland Islands (Malvinas)',  'FK', 'FLK', '238'),
-- ('Faroe Islands',  'FO', 'FRO', '234'),
-- ('Fiji',  'FJ', 'FJI', '242'),
-- ('Finland',  'FI', 'FIN', '246'),
-- ('France',  'FR', 'FRA', '250'),
-- ('French Guiana',  'GF', 'GUF', '254'),
-- ('French Polynesia',  'PF', 'PYF', '258'),
-- ('French Southern Territories',  'TF', 'ATF', '260'),
-- ('Gabon',  'GA', 'GAB', '266'),
-- ('Gambia',  'GM', 'GMB', '270'),
-- ('Georgia',  'GE', 'GEO', '268'),
-- ('Germany',  'DE', 'DEU', '276'),
-- ('Ghana',  'GH', 'GHA', '288'),
-- ('Gibraltar',  'GI', 'GIB', '292'),
-- ('Greece',  'GR', 'GRC', '300'),
-- ('Greenland',  'GL', 'GRL', '304'),
-- ('Grenada',  'GD', 'GRD', '308'),
-- ('Guadeloupe',  'GP', 'GLP', '312'),
-- ('Guam',  'GU', 'GUM', '316'),
-- ('Guatemala',  'GT', 'GTM', '320'),
-- ('Guernsey',  'GG', 'GGY', '831'),
-- ('Guinea',  'GN', 'GIN', '324'),
-- ('Guinea-Bissau',  'GW', 'GNB', '624'),
-- ('Guyana',  'GY', 'GUY', '328'),
-- ('Haiti',  'HT', 'HTI', '332'),
-- ('Heard and Mcdonald Islands',  'HM', 'HMD', '334'),
-- ('Holy See (Vatican City State)',  'VA', 'VAT', '336'),
-- ('Honduras',  'HN', 'HND', '340'),
-- ('Hungary',  'HU', 'HUN', '348'),
-- ('Iceland',  'IS', 'ISL', '352'),
-- ('India',  'IN', 'IND', '356'),
('Indonesia',  'ID', 'IDN', '360'),
-- ('Iran, Islamic Republic of',  'IR', 'IRN', '364'),
-- ('Iraq',  'IQ', 'IRQ', '368'),
-- ('Ireland',  'IE', 'IRL', '372'),
-- ('Isle of Man',  'IM', 'IMN', '833'),
-- ('Israel',  'IL', 'ISR', '376'),
-- ('Italy',  'IT', 'ITA', '380'),
-- ('Jamaica',  'JM', 'JAM', '388'),
-- ('Japan',  'JP', 'JPN', '392'),
-- ('Jersey',  'JE', 'JEY', '832'),
-- ('Jordan',  'JO', 'JOR', '400'),
-- ('Kazakhstan',  'KZ', 'KAZ', '398'),
-- ('Kenya',  'KE', 'KEN', '404'),
-- ('Kiribati',  'KI', 'KIR', '296'),
-- ('Korea (North)',  'KP', 'PRK', '408'),
-- ('Korea (South)',  'KR', 'KOR', '410'),
-- ('Kuwait',  'KW', 'KWT', '414'),
-- ('Kyrgyzstan',  'KG', 'KGZ', '417'),
-- ('Lao PDR',  'LA', 'LAO', '418'),
-- ('Latvia',  'LV', 'LVA', '428'),
-- ('Lebanon',  'LB', 'LBN', '422'),
-- ('Lesotho',  'LS', 'LSO', '426'),
-- ('Liberia',  'LR', 'LBR', '430'),
-- ('Libya',  'LY', 'LBY', '434'),
-- ('Liechtenstein',  'LI', 'LIE', '438'),
-- ('Lithuania',  'LT', 'LTU', '440'),
-- ('Luxembourg',  'LU', 'LUX', '442'),
-- ('Macedonia, Republic of',  'MK', 'MKD', '807'),
-- ('Madagascar',  'MG', 'MDG', '450'),
-- ('Malawi',  'MW', 'MWI', '454'),
('Malaysia',  'MY', 'MYS', '458'),
-- ('Maldives',  'MV', 'MDV', '462'),
-- ('Mali',  'ML', 'MLI', '466'),
-- ('Malta',  'MT', 'MLT', '470'),
-- ('Marshall Islands',  'MH', 'MHL', '584'),
-- ('Martinique',  'MQ', 'MTQ', '474'),
-- ('Mauritania',  'MR', 'MRT', '478'),
-- ('Mauritius',  'MU', 'MUS', '480'),
-- ('Mayotte',  'YT', 'MYT', '175'),
-- ('Mexico',  'MX', 'MEX', '484'),
-- ('Micronesia, Federated States of',  'FM', 'FSM', '583'),
-- ('Moldova',  'MD', 'MDA', '498'),
-- ('Monaco',  'MC', 'MCO', '492'),
-- ('Mongolia',  'MN', 'MNG', '496'),
-- ('Montenegro',  'ME', 'MNE', '499'),
-- ('Montserrat',  'MS', 'MSR', '500'),
-- ('Morocco',  'MA', 'MAR', '504'),
-- ('Mozambique',  'MZ', 'MOZ', '508'),
-- ('Myanmar',  'MM', 'MMR', '104'),
-- ('Namibia',  'NA', 'NAM', '516'),
-- ('Nauru',  'NR', 'NRU', '520'),
-- ('Nepal',  'NP', 'NPL', '524'),
-- ('Netherlands',  'NL', 'NLD', '528'),
-- ('Netherlands Antilles',  'AN', 'ANT', '530'),
-- ('New Caledonia',  'NC', 'NCL', '540'),
-- ('New Zealand',  'NZ', 'NZL', '554'),
-- ('Nicaragua',  'NI', 'NIC', '558'),
-- ('Niger',  'NE', 'NER', '562'),
-- ('Nigeria',  'NG', 'NGA', '566'),
-- ('Niue',  'NU', 'NIU', '570'),
-- ('Norfolk Island',  'NF', 'NFK', '574'),
-- ('Northern Mariana Islands',  'MP', 'MNP', '580'),
-- ('Norway',  'NO', 'NOR', '578'),
-- ('Oman',  'OM', 'OMN', '512'),
-- ('Pakistan',  'PK', 'PAK', '586'),
-- ('Palau',  'PW', 'PLW', '585'),
-- ('Palestinian Territory',  'PS', 'PSE', '275'),
-- ('Panama',  'PA', 'PAN', '591'),
-- ('Papua New Guinea',  'PG', 'PNG', '598'),
-- ('Paraguay',  'PY', 'PRY', '600'),
-- ('Peru',  'PE', 'PER', '604'),
-- ('Philippines',  'PH', 'PHL', '608'),
-- ('Pitcairn',  'PN', 'PCN', '612'),
-- ('Poland',  'PL', 'POL', '616'),
-- ('Portugal',  'PT', 'PRT', '620'),
-- ('Puerto Rico',  'PR', 'PRI', '630'),
-- ('Qatar',  'QA', 'QAT', '634'),
-- ('Réunion',  'RE', 'REU', '638'),
-- ('Romania',  'RO', 'ROU', '642'),
-- ('Russian Federation',  'RU', 'RUS', '643'),
-- ('Rwanda',  'RW', 'RWA', '646'),
-- ('Saint-Barthélemy',  'BL', 'BLM', '652'),
-- ('Saint Helena',  'SH', 'SHN', '654'),
-- ('Saint Kitts and Nevis',  'KN', 'KNA', '659'),
-- ('Saint Lucia',  'LC', 'LCA', '662'),
-- ('Saint-Martin (French part)',  'MF', 'MAF', '663'),
-- ('Saint Pierre and Miquelon',  'PM', 'SPM', '666'),
-- ('Saint Vincent and Grenadines',  'VC', 'VCT', '670'),
-- ('Samoa',  'WS', 'WSM', '882'),
-- ('San Marino',  'SM', 'SMR', '674'),
-- ('Sao Tome and Principe',  'ST', 'STP', '678'),
-- ('Saudi Arabia',  'SA', 'SAU', '682'),
-- ('Senegal',  'SN', 'SEN', '686'),
-- ('Serbia',  'RS', 'SRB', '688'),
-- ('Seychelles',  'SC', 'SYC', '690'),
-- ('Sierra Leone',  'SL', 'SLE', '694'),
('Singapore',  'SG', 'SGP', '702'),
-- ('Slovakia',  'SK', 'SVK', '703'),
-- ('Slovenia',  'SI', 'SVN', '705'),
-- ('Solomon Islands',  'SB', 'SLB', '090'),
-- ('Somalia',  'SO', 'SOM', '706'),
-- ('South Africa',  'ZA', 'ZAF', '710'),
-- ('South Georgia and the South Sandwich Islands',  'GS', 'SGS', '239'),
-- ('South Sudan',  'SS', 'SSD', '728'),
-- ('Spain',  'ES', 'ESP', '724'),
-- ('Sri Lanka',  'LK', 'LKA', '144'),
-- ('Sudan',  'SD', 'SDN', '736'),
-- ('Suriname',  'SR', 'SUR', '740'),
-- ('Svalbard and Jan Mayen Islands',  'SJ', 'SJM', '744'),
-- ('Swaziland',  'SZ', 'SWZ', '748'),
-- ('Sweden',  'SE', 'SWE', '752'),
-- ('Switzerland',  'CH', 'CHE', '756'),
-- ('Syrian Arab Republic (Syria)',  'SY', 'SYR', '760'),
-- ('Taiwan, Republic of China',  'TW', 'TWN', '158'),
-- ('Tajikistan',  'TJ', 'TJK', '762'),
-- ('Tanzania, United Republic of',  'TZ', 'TZA', '834'),
-- ('Thailand',  'TH', 'THA', '764'),
-- ('Timor-Leste',  'TL', 'TLS', '626'),
-- ('Togo',  'TG', 'TGO', '768'),
-- ('Tokelau',  'TK', 'TKL', '772'),
-- ('Tonga',  'TO', 'TON', '776'),
-- ('Trinidad and Tobago',  'TT', 'TTO', '780'),
-- ('Tunisia',  'TN', 'TUN', '788'),
-- ('Turkey',  'TR', 'TUR', '792'),
-- ('Turkmenistan',  'TM', 'TKM', '795'),
-- ('Turks and Caicos Islands',  'TC', 'TCA', '796'),
-- ('Tuvalu',  'TV', 'TUV', '798'),
-- ('Uganda',  'UG', 'UGA', '800'),
-- ('Ukraine',  'UA', 'UKR', '804'),
-- ('United Arab Emirates',  'AE', 'ARE', '784'),
('United Kingdom',  'GB', 'GBR', '826'),
('United States of America',  'US', 'USA', '840'),
-- ('US Minor Outlying Islands',  'UM', 'UMI', '581'),
-- ('Uruguay',  'UY', 'URY', '858'),
-- ('Uzbekistan',  'UZ', 'UZB', '860'),
-- ('Vanuatu',  'VU', 'VUT', '548'),
-- ('Venezuela (Bolivarian Republic)',  'VE', 'VEN', '862'),
-- ('Viet Nam',  'VN', 'VNM', '704'),
-- ('Virgin Islands, US',  'VI', 'VIR', '850'),
-- ('Wallis and Futuna Islands',  'WF', 'WLF', '876'),
-- ('Western Sahara',  'EH', 'ESH', '732'),
-- ('Yemen',  'YE', 'YEM', '887'),
-- ('Zambia',  'ZM', 'ZMB', '894'),
 ('Zimbabwe',  'ZW', 'ZWE', '716');

-- 



-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- Stored Procedures and functions
-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- drop the functions and types first
DROP FUNCTION IF EXISTS template;
DROP FUNCTION IF EXISTS get_user_id;
DROP FUNCTION IF EXISTS validate_nx_user_info;
DROP FUNCTION IF EXISTS insert_nx_user;
DROP FUNCTION IF EXISTS add_nx_user;
DROP FUNCTION IF EXISTS get_all_references ();
DROP FUNCTION IF EXISTS get_all_references_of_type;


DROP TYPE IF EXISTS insert_rc;
DROP TYPE IF EXISTS rc;

-- ----------------------------------------------------------
-- template for a typical function
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION template (username_   varchar)
RETURNS integer 
AS $BODY$
DECLARE _rc    integer;
BEGIN        
    _rc = 0;
    
    RETURN _rc;
END; $BODY$
LANGUAGE PLPGSQL;

-- ----------------------------------------------------------
-- user types
-- ----------------------------------------------------------
CREATE TYPE rc AS (
	code 	integer,
	message  varchar
);

-- ----------------------------------------------------------
-- user types
-- ----------------------------------------------------------
CREATE TYPE insert_rc AS (
    id          bigint,
	rc			rc
);




DROP FUNCTION IF EXISTS rc_as_json;
-- ----------------------------------------------------------
-- return rc in json format
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION rc_as_json (rc_  rc)
RETURNS text 
AS $BODY$
BEGIN
	RETURN (select row_to_json(t)
		from (select rc_.rc_code AS code, rc_.rc_message AS message) t );
END; $BODY$
LANGUAGE PLPGSQL;



DROP FUNCTION IF EXISTS ins_rc_as_json;
-- ----------------------------------------------------------
-- return insert_rc in json format
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION ins_rc_as_json (ins_resp_  insert_rc)
RETURNS text 
AS $BODY$
BEGIN
	RETURN (select row_to_json(t)
	from (select _ins_resp.id AS id, 
		    (select row_to_json(r)
				from (select ins_resp_.rc.rc_code AS code, ins_resp_.rc.rc_message AS message) r 
			) AS rc
	) t );
END; $BODY$
LANGUAGE PLPGSQL;

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
    SELECT id INTO _id FROM nx_user WHERE username = username_;
	IF (_id IS NULL) THEN
		_id := -1;
	END IF;
    RETURN _id;
END; $BODY$
LANGUAGE PLPGSQL;

-- ----------------------------------------------------------
-- validate user information
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION validate_nx_user_info (IN   username_   varchar,  
                                                  IN   password_   varchar,   
                                                  IN first_name_   varchar, 
                                                  IN  last_name_   varchar, 
                                                  IN      email_   varchar, 
                                                  IN      phone_   varchar, 
                                                  IN company_id_   bigint, 
                                                  IN created_by_   bigint)
RETURNS rc AS $BODY$
DECLARE   _rc   rc;
BEGIN
	_rc.code    = 0; 
	_rc.message = '';
    IF     (username_   = null OR username_   = '') THEN _rc.code := -1; _rc.message := 'username can not be null or blank';
    ELSEIF (password_   = null OR password_   = '') THEN _rc.code := -2; _rc.message := 'password can not be null or blank';
    ELSEIF (first_name_ = null OR first_name_ = '') THEN _rc.code := -3; _rc.message := 'first_name_ can not be null or blank';
    ELSEIF (email_      = null OR email_      = '') THEN _rc.code := -4; _rc.message := 'email can not be null or blank';
    ELSEIF (phone_      = null OR phone_      = '') THEN _rc.code := -5; _rc.message := 'phone_ can not be null or blank';
    ELSEIF (created_by_ < 1)                        THEN _rc.code := -6; _rc.message := 'Invalid user for created_by column';
    ELSEIF (company_id_ < 1)                        THEN _rc.code := -7; _rc.message := 'Invalid company ID';
	END IF;
	RETURN _rc;
END; $BODY$
LANGUAGE PLPGSQL;


-- ----------------------------------------------------------
-- insert into user table 
-- insert_nx_user ('nexSOFT_admin', 'password', 'foad', 'momtazi', 'foad.momtazi@nexsoft.co.id', '+62 666 666 666 666', 1,1)
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION insert_nx_user (IN    username_   varchar,  
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
	INSERT INTO nx_user (username, password, first_name, last_name, email, phone, company_id, created_by, last_updated_by)
		 VALUES (username_, password_, first_name_, last_name_, email_, phone_, company_id_, created_by_, created_by_)
	RETURNING id INTO _id;
--	RAISE NOTICE '4:returned id from insert = %', _id;
	RETURN _id;
END; $BODY$
LANGUAGE PLPGSQL;


-- ----------------------------------------------------------
-- add user 
-- add_nx_user ('nexSOFT_admin', 'password', 'foad', 'momtazi', 'foad.momtazi@nexsoft.co.id', '+62 666 666 666 666', 1,1)
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION add_nx_user (IN     username_   varchar,  
                                        IN     password_   varchar,   
                                        IN   first_name_   varchar, 
                                        IN    last_name_   varchar, 
                                        IN        email_   varchar, 
                                        IN        phone_   varchar, 
                                        IN   company_id_   bigint, 
                                        IN   created_by_   bigint)
RETURNS text AS $BODY$
DECLARE 
    _ins_resp            insert_response;
    _validate_response   response_record;
BEGIN
	_validate_response.rc_code = 
	_validate_response.rc_message = '';
	
    _validate_response := validate_nx_user_info (username_, password_, first_name_, last_name_, email_, phone_, company_id_, created_by_) ;
	_ins_resp.rc = _validate_response;
    IF (_validate_response.rc_code = 0) THEN
    	_ins_resp.id := get_user_id (username_);
        IF (_ins_resp.id > 0) THEN
            _validate_response.rc_code    := -9;
            _validate_response.rc_message := 'username is already in use';
			_ins_resp.rc = _validate_response;
			_ins_resp.id = -1;
        ELSE 
			_ins_resp.id := insert_nx_user (username_, password_, first_name_, 
		                              last_name_, email_, phone_, company_id_, created_by_);
        END IF;
    END IF;
	
	RETURN (select row_to_json(t)
	from (
	  select _ins_resp.id AS id, 
		    (select row_to_json(r)
				from (select _validate_response.rc_code AS code, _validate_response.rc_message AS message) r 
			) AS rc
	) t );
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
    SELECT id INTO _id FROM nx_role WHERE code = role_code_;
	IF (_id IS NULL) THEN
		_id := -1;
	END IF;
    RETURN _id;
END; $BODY$
LANGUAGE PLPGSQL;



DROP FUNCTION IF EXISTS insert_nx_role ();
-- ----------------------------------------------------------
-- create a new role
-- select create_role ()
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION insert_nx_role (IN    role_code_   varchar,
                                           IN  description_   varchar
                                           IN   created_by_   bigint)
RETURNS bigint AS $BODY$
DECLARE 
    _id                  bigint;
BEGIN
    INSERT INTO nx_role (code, description, created_by, last_updated_by)
            VALUES (role_code_, description_, created_by_, created_by_)
            RETURNING id INTO _id;
    RETURN _id;
END; $BODY$
LANGUAGE plpgSQL;




DROP FUNCTION IF EXISTS create_role ();
-- ----------------------------------------------------------
-- create a new role
-- select create_role ()
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION create_role (IN    role_code_   varchar,
                                        IN  description_   varchar
                                        IN   created_by_   bigint)
RETURNS text AS $BODY$
DECLARE 
    _ins_resp            insert_response;
    _validate_response   response_record;
    _id                  bigint;
BEGIN
    _id = select get_role_id (role_code_);
    if (_id > 0) THEN
        _validate_response.rc_code    := -19;
        _validate_response.rc_message := 'role is already in use';
        _ins_resp.rc = _validate_response;
        _ins_resp.id = -1;
    ELSE
        _ins_resp.id = insert_nx_role (role_code_, description_, created_by_);
    END IF;

    RETURN _ins_resp;

END; $BODY$
LANGUAGE plpgSQL;













-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nx_role_pkey_seq;
CREATE TABLE IF NOT EXISTS nx_role
(
    id              bigint UNIQUE NOT NULL DEFAULT nextval('nx_role_pkey_seq'::regclass),
    code            varchar (32)  NOT NULL references nx_reference (code),
    description     varchar (256) NOT NULL,

    created_by      bigint NOT NULL,
    created_at      TIMESTAMP DEFAULT now (),
    last_updated_by bigint NOT NULL,
    last_updated_at TIMESTAMP DEFAULT now (),
    record_status   varchar (32) NOT NULL default 'active' references nx_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;
COMMENT ON TABLE nx_role IS 'roles within the authorisation';
CREATE INDEX nx_role_code_index           ON nx_role (code);


-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nx_user_role_pkey_seq;
CREATE TABLE IF NOT EXISTS nx_user_role
(
    id              bigint UNIQUE NOT NULL DEFAULT nextval('nx_user_role_pkey_seq'::regclass),
    user_id         bigint NOT NULL references nx_user (id),
    role_id         bigint NOT NULL references nx_role (id),

    created_by      bigint NOT NULL,
    created_at      TIMESTAMP DEFAULT now (),
    last_updated_by bigint NOT NULL,
    last_updated_at TIMESTAMP DEFAULT now (),
    record_status   varchar (32) NOT NULL default 'active' references nx_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;
COMMENT ON TABLE nx_user_role IS 'links a user to role(s)';
CREATE INDEX nx_user_role_user_id_index           ON nx_user_role (user_id);
CREATE INDEX nx_user_role_role_id_index           ON nx_user_role (role_id);


-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nx_data_scope_pkey_seq;
CREATE TABLE IF NOT EXISTS    nx_data_scope
(
    id                  bigint UNIQUE NOT NULL DEFAULT nextval('nx_data_scope_pkey_seq'::regclass),
    name                varchar (128) NOT NULL,
    type                varchar (32) NOT NULL,

    created_at          TIMESTAMP DEFAULT now (),
    last_updated_by     bigint NOT NULL,
    last_updated_at     TIMESTAMP DEFAULT now (),
    record_status       varchar (32) NOT NULL default 'active' references nx_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;
CREATE INDEX nx_data_scope_code_index ON nx_data_scope (name);
CREATE INDEX nx_data_scope_type_index ON nx_data_scope (type);




-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nx_data_scope_list_pkey_seq;
CREATE TABLE IF NOT EXISTS    nx_data_scope_list
(
    id                  bigint UNIQUE NOT NULL DEFAULT nextval('nx_data_scope_list_pkey_seq'::regclass),
    data_scope_id       bigint NOT NULL references nx_data_scope (id),
-- ??    data                varchar (32) NOT NULL,

    created_at          TIMESTAMP DEFAULT now (),
    last_updated_by     bigint NOT NULL,
    last_updated_at     TIMESTAMP DEFAULT now (),
    record_status       varchar (32) NOT NULL default 'active' references nx_reference (code)
) 
WITH (OIDS = FALSE)
TABLESPACE pg_default;










-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- ----------------------------------------------------------

DROP FUNCTION IF EXISTS get_all_references ();
-- ----------------------------------------------------------
-- get all the references, sorted by type and code 
-- select * from get_all_references ()
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION get_all_references ()
RETURNS text AS $$
BEGIN
	RETURN (SELECT array_to_json(array_agg(t ))
	from (
		    SELECT id, type, code FROM nx_reference ORDER BY type, code
		 )  t ) ;

END; $$
LANGUAGE plpgSQL;


DROP FUNCTION IF EXISTS get_all_references (varchar);
-- ----------------------------------------------------------
-- get all the references of a type, sorted by code 
-- select * from get_all_references ('address')
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION get_all_references (type_ varchar)
RETURNS text AS $$
BEGIN
	RETURN (SELECT array_to_json(array_agg(t ))
	from (
		    SELECT id, type, code FROM nx_reference WHERE type=type_ ORDER BY code
		 )  t ) ;

END; $$
LANGUAGE plpgSQL;



-- ----------------------------------------------------------
-- ----------------------------------------------------------
-- ----------------------------------------------------------

DROP FUNCTION IF EXISTS get_all_countries ();
-- ----------------------------------------------------------
-- get all the countries sorted by cc_code 
-- select * from get_all_countries ()
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION get_all_countries ()
RETURNS text AS $$
BEGIN
	RETURN (SELECT array_to_json(array_agg(t ))
	from (
		    SELECT id, name, A2_code, A3_code, un_code, cc_code FROM nx_country ORDER BY cc_code
		 )  t ) ;

END; $$
LANGUAGE plpgSQL;


-- 
-- 
-- 
-- 
-- 
-- 





