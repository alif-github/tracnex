-- some timing in PHP for different hashes.
-- the timings are from internet and not verified but offers relative cost
-- postgres:
--      digest(data bytea, type text) returns bytea
--          Standard algorithms are md5, sha1, sha224, sha256, sha384 and sha512
--
-- hash('crc32', 'The quick brown fox jumped over the lazy dog.');#  750ms   8 chars
-- hash('crc32b','The quick brown fox jumped over the lazy dog.');#  700ms   8 chars
-- hash('md5',   'The quick brown fox jumped over the lazy dog.');#  770ms  32 chars
-- hash('sha1',  'The quick brown fox jumped over the lazy dog.');#  880ms  40 chars
-- hash('sha256','The quick brown fox jumped over the lazy dog.');# 1490ms  64 chars
-- hash('sha384','The quick brown fox jumped over the lazy dog.');# 1830ms  96 chars
-- hash('sha512','The quick brown fox jumped over the lazy dog.');# 1870ms 128 chars
-- ============================================================================
-- use sha512 when security is also a concern


--	GET STACKED DIAGNOSTICS v_error_stack = PG_EXCEPTION_CONTEXT;
--	raise notice 'insert failed state % %', SQLSTATE, v_error_stack;
--	raise notice 'insert failed with the state SQLSTATE = % sqlerrm = %', SQLSTATE, sqlerrm;
--  RAISE;   after catching an exception this re-throws the exception up one level


CREATE TYPE record_status  AS ENUM ('active', 'locked', 'deleted', 'expired');
CREATE TYPE error_level    AS ENUM ('ERROR', 'WARNING', 'INFO', 'DEBUG');
CREATE TYPE error_severity AS ENUM ('CRITICAL', 'MAJOR', 'ALERT', 'MINOR', 'NOTICE');
CREATE TYPE language_code  AS ENUM ('id', 'en');

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



DROP FUNCTION IF EXISTS create_rc_json_type (varchar, varchar);
-- ----------------------------------------------------------
-- return json format in the form of:
-- {
--      "code"    : "",
--      "message" : ""
-- }
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION create_rc_json_type(IN    rc_code_ varchar,
                                               IN rc_message_ varchar)
RETURNS json
AS $BODY$
    DECLARE   _rc varchar;
BEGIN
	RETURN (SELECT
		    (SELECT row_to_json(r)
				FROM (SELECT rc_code_ AS code, rc_message_ AS message) r
			)
	);
END;
$BODY$
LANGUAGE PLPGSQL;



DROP FUNCTION IF EXISTS create_query_rc_json_type (varchar, varchar, json);
-- ----------------------------------------------------------
-- return json format in the form of:
-- {
--      "rc":
--      {
--          "code"    : "",
--          "message" : ""
--      },
--      "data":
--      {
--          p_data
--      }
-- }
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION create_query_rc_json_type(IN    p_rc_code varchar,
                                                     IN p_rc_message varchar,
                                                     IN       p_data json)
RETURNS json
AS $BODY$
    DECLARE   _rc json;
BEGIN
	_rc = create_rc_json_type (p_rc_code, p_rc_message);
	RETURN (SELECT row_to_json(t) as data
	          FROM (SELECT _rc AS rc, p_data::json as data) t );
END;
$BODY$
LANGUAGE PLPGSQL;


DROP FUNCTION IF EXISTS create_query_error_constant_json_type (varchar, language_code, json);
-- ----------------------------------------------------------
-- return json format in the form of:
-- {
--      "rc":
--      {
--          "code"    : "",
--          "message" : ""
--      },
--      "data":
--      {
--          null
--      }
-- }
-- but rc is made up of code and message taken
-- from table nex_error_code for the index error_constant
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION create_query_error_constant_json_type(IN    p_error_constant varchar,
                                                                 IN         p_lang_code language_code,
                                                                 IN              p_data json)
RETURNS json
AS $BODY$
    DECLARE     _rc             json;
                _error_code     varchar;
                _error_message  varchar;
BEGIN
    SELECT ec.code, em.message INTO _error_code, _error_message
          FROM nex_error_code ec, nex_error_message em
         WHERE ec.id = em.error_id
           AND ec.constant = p_error_constant
           AND em.language = p_lang_code;

	_rc = create_rc_json_type (_error_code, _error_message);
	RETURN (SELECT row_to_json(t) as data
	          FROM (SELECT _rc AS rc, p_data::json as data) t );
END;
$BODY$
LANGUAGE PLPGSQL;


DROP FUNCTION IF EXISTS create_success_json_type (json);
-- ----------------------------------------------------------
-- return json format in the form of:
-- {
--      "rc":
--      {
--          "code"    : "OK",
--          "message" : ""
--      },
--      "data":
--      {
--          p_data
--      }
-- it assumes there are no errors and therefore code and
-- message are blank
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION create_success_json_type(IN p_data json)
RETURNS json
AS $BODY$
    DECLARE     _rc             json;
                _error_code     varchar;
                _error_message  varchar;
BEGIN
	_rc = create_rc_json_type ('OK', 'OK');
	RETURN (SELECT row_to_json(t) as data
	          FROM (SELECT _rc AS rc, p_data::json as data) t );
END;
$BODY$
LANGUAGE PLPGSQL;


DROP FUNCTION IF EXISTS create_insert_rc_ok_json_type (bigint);
-- ----------------------------------------------------------
-- return json format in the form of:
-- {
--      "rc":
--      {
--          "code"    : "OK",
--          "message" : ""
--      },
--      "data":
--      {
--          "id" :     bigint/int64
--      }
-- it assumes there are no errors and therefore code and
-- message are blank
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION create_insert_rc_ok_json_type(IN p_id  bigint)
RETURNS json
AS $BODY$
    DECLARE     _rc             json;
BEGIN

	_rc = create_rc_json_type ('OK', '');
	RETURN (SELECT row_to_json(t) as data
	          FROM (SELECT _rc AS rc, p_id AS id) t );
END;
$BODY$
LANGUAGE PLPGSQL;



DROP FUNCTION IF EXISTS create_insert_rc_json_type (varchar, language_code, bigint);
-- ----------------------------------------------------------
-- return json format in the form of:
-- {
--      "rc":
--      {
--          "code"    : "",
--          "message" : ""
--      },
--      "data":
--      {
--          "id" :     bigint/int64
--      }
-- it assumes there are no errors and therefore code and
-- message are blank
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION create_insert_rc_json_type(IN    p_error_constant  varchar,
                                                      IN         p_lang_code  language_code,
                                                      IN                p_id  bigint)
RETURNS json
AS $BODY$
    DECLARE     _rc             json;
                _error_code     varchar;
                _error_message  varchar;
BEGIN
    SELECT ec.code, em.message INTO _error_code, _error_message
          FROM nex_error_code ec, nex_error_message em
         WHERE ec.id = em.error_id
           AND ec.constant = p_error_constant
           AND em.language = p_lang_code;

	_rc = create_rc_json_type (_error_code, _error_message);
	RETURN (SELECT row_to_json(t) as data
	          FROM (SELECT _rc AS rc, p_id AS id) t );
END;
$BODY$
LANGUAGE PLPGSQL;


DROP FUNCTION IF EXISTS create_rc_from_constant_custom_code (varchar, varchar, json);
-- ----------------------------------------------------------
-- return json format in the form of:
-- {
--      "rc":
--      {
--          "code"    : "",
--          "message" : ""
--      },
--      "data":
--      {
--          "id" :     bigint/int64
--      }
-- it assumes there are no errors and therefore code and
-- message are blank
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION create_rc_from_constant_custom_code(
                                IN    p_error_constant  varchar,
                                IN           p_message  varchar,
                                IN              p_json  json)
RETURNS json
AS $BODY$
    DECLARE     _rc             json;
                _error_code     varchar;
                _error_message  varchar;
BEGIN
    SELECT ec.code INTO _error_code
          FROM nex_error_code ec, nex_error_message em
         WHERE ec.id = em.error_id
           AND ec.constant = p_error_constant
           AND em.language = 'en';
    IF NOT FOUND THEN
        _error_code = p_error_constant;
    END IF;
	_rc = create_rc_json_type (_error_code, p_message);
	RETURN (SELECT row_to_json(t) as data
	          FROM (SELECT _rc AS rc, p_json AS data) t );
END;
$BODY$
LANGUAGE PLPGSQL;

-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nex_reference_pkey_seq;
CREATE TABLE IF NOT EXISTS nex_reference
(
    id              bigint UNIQUE NOT NULL DEFAULT nextval('nex_reference_pkey_seq'::regclass),
    type            varchar(32) NOT NULL,
    code            varchar(32) NOT NULL UNIQUE
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;
CREATE INDEX nex_reference_type_index ON nex_reference (type);
CREATE INDEX nex_reference_code_index ON nex_reference (code);
CREATE INDEX nex_reference_combined_index  ON nex_reference (type, code);

INSERT INTO nex_reference (type, code) VALUES ('status',  'active');
INSERT INTO nex_reference (type, code) VALUES ('status',  'deleted');
INSERT INTO nex_reference (type, code) VALUES ('status',  'locked');
INSERT INTO nex_reference (type, code) VALUES ('status',  'change_password');
INSERT INTO nex_reference (type, code) VALUES ('contact', 'email');
INSERT INTO nex_reference (type, code) VALUES ('contact', 'phone');
INSERT INTO nex_reference (type, code) VALUES ('company', 'parent');
INSERT INTO nex_reference (type, code) VALUES ('company', 'division');
INSERT INTO nex_reference (type, code) VALUES ('company', 'branch');
INSERT INTO nex_reference (type, code) VALUES ('company', 'subsidiary');
INSERT INTO nex_reference (type, code) VALUES ('address', 'physical');
INSERT INTO nex_reference (type, code) VALUES ('address', 'postal');
INSERT INTO nex_reference (type, code) VALUES ('address', 'delivery');
INSERT INTO nex_reference (type, code) VALUES ('address', 'POBox');


-- --------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nex_country_pkey_seq;
CREATE TABLE IF NOT EXISTS nex_country
(
    id              bigint UNIQUE NOT NULL DEFAULT nextval('nex_country_pkey_seq'::regclass),
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
CREATE INDEX nex_country_name_index     ON nex_country (name);
CREATE INDEX nex_country_A2_code_index  ON nex_country (A2_code);
CREATE INDEX nex_country_A3_code_index  ON nex_country (A3_code);
CREATE INDEX nex_country_un_code_index  ON nex_country (un_code);
CREATE INDEX nex_country_cc_code_index  ON nex_country (cc_code);

-- ----------------------------------------------------------
-- insert reference data
-- ----------------------------------------------------------
-- country
INSERT INTO nex_country (name, A2_code, A3_code, un_code)
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
-- ('Côte d'Ivoire',  'CI', 'CIV', '384'),
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

-- -----------------------------------------------------------------------------
-- nex_error_code table
-- -----------------------------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nex_error_code_pkey_seq;
CREATE TABLE IF NOT EXISTS nex_error_code
(
    id                  bigint UNIQUE NOT NULL DEFAULT nextval('nex_error_code_pkey_seq'::regclass),
    code                varchar (16),
    constant            varchar (64),
    level               error_level,
    severity            error_severity,

    created_at          TIMESTAMP DEFAULT now ()
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;
CREATE INDEX nex_error_code_code_index     ON nex_error_code (code);
CREATE INDEX nex_error_code_constant_index ON nex_error_code (constant);


-- -----------------------------------------------------------------------------
-- nex_error_message table
-- -----------------------------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nex_error_message_pkey_seq;
CREATE TABLE IF NOT EXISTS nex_error_message
(
    id                  bigint UNIQUE NOT NULL DEFAULT nextval('nex_error_message_pkey_seq'::regclass),
    error_id            bigint NOT NULL references nex_error_code (id),
    message             varchar (512),        /* this is the message text in Indonesian hence _id */
    language            language_code,        /* this is the message text in englilsh hence _en */

    created_at          TIMESTAMP DEFAULT now ()
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;
CREATE INDEX nex_error_message_error_id_index ON nex_error_message(error_id);

DROP SEQUENCE IF EXISTS nex_locale_text_pkey_seq;
-- -----------------------------------------------------------------------------
-- nex_locale_text table
-- used to store all the language texts
-- -----------------------------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS nex_locale_text_pkey_seq;
CREATE TABLE IF NOT EXISTS nex_locale_text
(
    id                  bigint UNIQUE NOT NULL DEFAULT nextval('nex_locale_text_pkey_seq'::regclass),
    text_code           varchar (32) UNIQUE NOT NULL,
    module              varchar (64) NOT NULL,
    text                varchar (512),        /* this is the message text in Indonesian hence _id */
    language            language_code,        /* this is the message text in englilsh hence _en */

    created_at          TIMESTAMP DEFAULT now ()
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;
CREATE INDEX nex_locale_text_module_index ON nex_locale_text (module);

INSERT INTO nex_locale_text (text_code, module, text, language)
     VALUES ('T123476', 'test', 'Ini adalah pesan pertama dalam bahasa indonesia', 'id'),
            ('T123456', 'test', 'This is the first message in english', 'en');

-- -----------------------------------------------------------------------------
DROP FUNCTION IF EXISTS check_for_duplicate_error (varchar, varchar);
-- ----------------------------------------------------------
-- return id if key exists
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION check_for_duplicate_error (IN      p_error_code  varchar,
                                                      IN  p_error_constant  varchar)
RETURNS bigint AS $BODY$
    DECLARE     v_id         bigint;
BEGIN
    SELECT id INTO v_id
      FROM nex_error_code
     WHERE code = p_error_code OR constant = p_error_constant;
    IF NOT FOUND THEN
        v_id = -1;
    ELSE
        v_id = 0;
    END IF;
    RETURN v_id;

END;

$BODY$
LANGUAGE 'plpgsql';

DROP FUNCTION IF EXISTS get_error_id_from_code (varchar);
-- ----------------------------------------------------------
-- return id if key exists
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION get_error_id_from_code (IN  p_error_code  varchar)
RETURNS bigint AS $BODY$
    DECLARE     v_id         bigint;
BEGIN
    SELECT id INTO v_id
      FROM nex_error_code WHERE code = p_error_code;
    IF NOT FOUND THEN
        v_id = -1;
    END IF;
    RETURN v_id;

END;

$BODY$
LANGUAGE 'plpgsql';



DROP FUNCTION IF EXISTS get_error_id_from_constant (varchar);
-- ----------------------------------------------------------
-- return id if key exists
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION get_error_id_from_constant (IN  p_constant  varchar)
RETURNS bigint AS $BODY$
    DECLARE     v_id         bigint;
BEGIN
    SELECT id INTO v_id
      FROM nex_error_code WHERE constant = p_constant;
    IF NOT FOUND THEN
        v_id = -1;
    END IF;
    RETURN v_id;

END;

$BODY$
LANGUAGE 'plpgsql';


DROP FUNCTION IF EXISTS create_new_error(varchar, varchar, error_level, error_severity, varchar, language_code);
-- ----------------------------------------------------------
-- return status of a single service
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION create_new_error (IN      p_error_code  varchar,
                                             IN  p_error_constant  varchar,
                                             IN           p_level  error_level,
                                             IN        p_severity  error_severity,
                                             IN      p_error_text  varchar,
                                             IN        p_language  language_code)
RETURNS   json
AS $BODY$
    DECLARE     v_id    bigint = -1;
BEGIN
    v_id = check_for_duplicate_error (p_error_code, p_error_constant);
    IF v_id != -1 THEN
        RETURN create_insert_rc_json_type ('ROW_NOT_FOUND', 'en', null);
    ELSE
        INSERT INTO nex_error_code (code, constant, level, severity)
                           VALUES  (p_error_code, p_error_constant, p_level, p_severity)
            RETURNING id INTO v_id;
        INSERT INTO nex_error_message (error_id, message, language)
                           VALUES  (v_id, p_error_text, p_language);
        RETURN create_insert_rc_json_type (null, null, v_id);
    END IF;

END; $BODY$
LANGUAGE PLPGSQL;

DROP FUNCTION IF EXISTS get_error_message_from_code(varchar, language_code);
-- ----------------------------------------------------------
-- return status of a single service
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION get_error_message_from_code (IN  p_error_code  varchar,
                                                        IN   p_lang_code  language_code)
RETURNS   record
AS $BODY$
    DECLARE     _data    record;
BEGIN
    SELECT ec.id, ec.code, ec.constant, em.message INTO _data
          FROM nex_error_code ec, nex_error_message em
         WHERE ec.id = em.error_id
           AND ec.code = p_error_code
           AND em.language = p_lang_code;

    RETURN _data;
END; $BODY$
LANGUAGE PLPGSQL;


-- ----------------------------------------------------------
-- return status of a single service
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION get_error_info_from_constant(p_error_const character varying,
                                                        p_lang_code language_code)
    RETURNS json
    LANGUAGE 'plpgsql'
AS
$BODY$
DECLARE
    v_result  JSON;
    v_code    VARCHAR;
    v_message VARCHAR;
BEGIN
	SELECT ec.code, em.message INTO v_code, v_message
              FROM nex_error_code ec, nex_error_message em
             WHERE ec.id = em.error_id
               AND ec.constant = p_error_const
               AND em.language = p_lang_code;
    IF NOT FOUND THEN
		SELECT ec.code, em.message INTO v_code, v_message
              FROM nex_error_code ec, nex_error_message em
             WHERE ec.id = em.error_id
               AND ec.code = p_error_const
               AND em.language = p_lang_code;
		IF NOT FOUND THEN
			RETURN create_success_json_type(json_build_object('code',' E-5-UNK-UNK-001','message','Unknown error message. Please call our customer service'));
		ELSE
			RETURN create_success_json_type(json_build_object('code',v_code,'message',v_message));
		END IF;
	ELSE
		RETURN create_success_json_type(json_build_object('code',v_code,'message',v_message));
	END IF;

END$BODY$;


DROP FUNCTION IF EXISTS read_all_text();
-- ----------------------------------------------------------
-- return status of a single service
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION read_all_text ()
RETURNS JSON
AS $BODY$
BEGIN
    RETURN (SELECT array_to_json(array_agg(t ))
	from (
             SELECT text_code, text, language
             FROM nex_locale_text
         )  t ) ;

END; $BODY$
LANGUAGE PLPGSQL;

DROP FUNCTION IF EXISTS read_text (language_code);
-- ----------------------------------------------------------
-- return status of a single service
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION read_text (lang_  language_code)
RETURNS JSON
AS $BODY$
BEGIN
    RETURN (SELECT array_to_json(array_agg(t ))
	from (
		    SELECT text_code, text FROM nex_locale_text WHERE language=lang_
		 )  t ) ;

END; $BODY$
LANGUAGE PLPGSQL;

DROP FUNCTION IF EXISTS read_text_id ();
-- ----------------------------------------------------------
-- return status of a single service
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION read_text_id ()
RETURNS JSON
AS $BODY$
BEGIN
    RETURN read_text ('id');

END; $BODY$
LANGUAGE PLPGSQL;


DROP FUNCTION IF EXISTS read_text_en ();
-- ----------------------------------------------------------
-- return status of a single service
-- ----------------------------------------------------------
CREATE OR REPLACE FUNCTION read_text_en ()
RETURNS JSON
AS $BODY$
BEGIN
    RETURN read_text ('en');

END; $BODY$
LANGUAGE PLPGSQL;

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
		    SELECT id, type, code FROM nex_reference ORDER BY type, code
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
		    SELECT id, type, code FROM nex_reference WHERE type=type_ ORDER BY code
		 )  t ) ;

END; $$
LANGUAGE plpgSQL;

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
		    SELECT id, name, A2_code, A3_code, un_code, cc_code FROM nex_country ORDER BY cc_code
		 )  t ) ;

END; $$
LANGUAGE plpgSQL;

CREATE OR REPLACE FUNCTION create_rc_from_error_code(p_code character varying,
                                                     p_message character varying,
                                                     p_data json)
    RETURNS json
    LANGUAGE 'plpgsql'
AS
$BODY$
DECLARE
    _rc            json;
    _error_code    varchar;
    _error_message varchar;
BEGIN
	_rc = create_rc_json_type (p_code, p_message);
	RETURN (SELECT row_to_json(t) as data
	          FROM (SELECT _rc AS rc, p_data::json as data) t );
END;
$BODY$;

-- =============================================================================
CREATE OR REPLACE FUNCTION check_table_exists(p_table_name character varying,
                                              p_schema_name character varying)
    RETURNS boolean
    LANGUAGE 'plpgsql'
AS
$BODY$
BEGIN
    RETURN (SELECT EXISTS(
                           SELECT 1
                           FROM information_schema.tables
                           WHERE table_schema = p_schema_name
   		AND    table_name = p_table_name
   ));
END$BODY$;
