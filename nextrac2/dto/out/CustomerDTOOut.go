package out

import "time"

type GetListCustomerResponse struct {
	ID          int64     `json:"id"`
	CompanyID   string    `json:"company_id"`
	BranchID    string    `json:"branch_id"`
	CompanyName string    `json:"company_name"`
	Product     string    `json:"product"`
	UserAmount  int64     `json:"user_amount"`
	ExpDate     time.Time `json:"exp_date"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ViewCustomerResponse struct {
	CompanyID      string    `json:"company_id"`
	BranchID       string    `json:"branch_id"`
	CompanyName    string    `json:"company_name"`
	City           string    `json:"city"`
	Implementer    string    `json:"implementer"`
	Implementation time.Time `json:"implementation"`
	Product        string    `json:"product"`
	Version        string    `json:"version"`
	LicenseType    string    `json:"license_type"`
	UserAmount     int64     `json:"user_amount"`
	ExpDate        time.Time `json:"exp_date"`
	CreatedBy      int64     `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedClient  string    `json:"created_client"`
	UpdatedBy      int64     `json:"updated_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedClient  string    `json:"updated_client"`
}

type CustomerListResponse struct {
	ID           int64     `json:"id"`
	Npwp         string    `json:"npwp"`
	CustomerName string    `json:"customer_name"`
	Address      string    `json:"address"`
	ProvinceID   int64     `json:"province_id"`
	ProvinceName string    `json:"province_name"`
	DistrictID   int64     `json:"district_id"`
	DistrictName string    `json:"district_name"`
	Phone        string    `json:"phone"`
	Status       string    `json:"status"`
	CreatedBy    int64     `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedBy    int64     `json:"updated_by"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CustomerListByStatusResponse struct {
	ID                  int64     `json:"id"`
	MDBCompanyProfileID int64     `json:"mdb_company_profile_id"`
	Npwp                string    `json:"npwp"`
	CustomerName        string    `json:"customer_name"`
	Address             string    `json:"address"`
	ProvinceID          int64     `json:"province_id"`
	ProvinceName        string    `json:"province_name"`
	DistrictID          int64     `json:"district_id"`
	DistrictName        string    `json:"district_name"`
	Phone               string    `json:"phone"`
	Status              string    `json:"status"`
	CreatedBy           int64     `json:"created_by"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedBy           int64     `json:"updated_by"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type CustomerViewResponse struct {
	ID                      int64                         `json:"id"`
	IsPrincipal             bool                          `json:"is_principal"`
	IsParent                bool                          `json:"is_parent"`
	ParentCustomerID        int64                         `json:"parent_customer_id"`
	ParentCustomerName      string                        `json:"parent_customer_name"`
	MDBParentCustomerID     int64                         `json:"mdb_parent_customer_id"`
	MDBCompanyProfileID     int64                         `json:"mdb_company_profile_id"`
	Npwp                    string                        `json:"npwp"`
	MDBCompanyTitleID       int64                         `json:"mdb_company_title_id"`
	CompanyTitle            string                        `json:"company_title"`
	CustomerName            string                        `json:"customer_name"`
	Address                 string                        `json:"address"`
	Address2                string                        `json:"address_2"`
	Address3                string                        `json:"address_3"`
	Hamlet                  string                        `json:"hamlet"`
	Neighbourhood           string                        `json:"neighbourhood"`
	CountryID               int64                         `json:"country_id"`
	ProvinceID              int64                         `json:"province_id"`
	ProvinceName            string                        `json:"province_name"`
	DistrictID              int64                         `json:"district_id"`
	DistrictName            string                        `json:"district_name"`
	SubDistrictID           int64                         `json:"sub_district_id"`
	SubDistrictName         string                        `json:"sub_district_name"`
	UrbanVillageID          int64                         `json:"urban_village_id"`
	UrbanVillageName        string                        `json:"urban_village_name"`
	PostalCodeID            int64                         `json:"postal_code_id"`
	PostalCode              string                        `json:"postal_code_name"`
	Longitude               float64                       `json:"long"`
	Latitude                float64                       `json:"lat"`
	Phone                   string                        `json:"phone"`
	AlternativePhone        string                        `json:"alternative_phone"`
	Fax                     string                        `json:"fax"`
	CompanyEmail            string                        `json:"company_email"`
	AlternativeCompanyEmail string                        `json:"alternative_company_email"`
	CustomerSource          string                        `json:"customer_source"`
	TaxName                 string                        `json:"tax_name"`
	TaxAddress              string                        `json:"tax_address"`
	SalesmanID              int64                         `json:"salesman_id"`
	SalesmanName            string                        `json:"salesman_name"`
	RefCustomerID           int64                         `json:"ref_customer_id"`
	RefCustomerName         string                        `json:"ref_customer_name"`
	DistributorOF           string                        `json:"distributor_of"`
	CustomerGroupID         int64                         `json:"customer_group_id"`
	CustomerGroupName       string                        `json:"customer_group_name"`
	CustomerCategoryID      int64                         `json:"customer_category_id"`
	CustomerCategoryName    string                        `json:"customer_category_name"`
	Status                  string                        `json:"status"`
	IsUsed                  bool                          `json:"is_used"`
	CustomerContact         []CustomerContactViewResponse `json:"customer_contact"`
	CreatedBy               int64                         `json:"created_by"`
	CreatedAt               time.Time                     `json:"created_at"`
	CreatedName             string                        `json:"created_name"`
	UpdatedBy               int64                         `json:"updated_by"`
	UpdatedAt               time.Time                     `json:"updated_at"`
	UpdatedName             string                        `json:"updated_name"`
}

type CustomerErrorResponse struct {
	InsertCustomerResponse DetailCustomerInsertResponse `json:"customer_response"`
	PreviousRequest        PreviousPayload              `json:"previous_payload"`
}

type DetailCustomerInsertResponse struct {
	InsertCustomerDetail        CustomerErrorStatus        `json:"customer_detail"`
	InsertCustomerContactDetail CustomerContactErrorStatus `json:"customer_contact_detail"`
}

type CustomerErrorStatus struct {
	Npwp    string `json:"npwp"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type CustomerContactErrorStatus struct {
	Nik     string `json:"nik"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type PreviousPayload struct {
	ID                          int64                     `json:"id"`
	IsPrincipal                 bool                      `json:"is_principal"`
	IsParent                    bool                      `json:"is_parent"`
	ParentCustomerID            int64                     `json:"parent_customer_id"`
	MDBParentCustomerID         int64                     `json:"mdb_parent_customer_id"`
	MDBCompanyProfileID         int64                     `json:"mdb_company_profile_id"`
	Npwp                        string                    `json:"npwp"`
	MDBCompanyTitleID           int64                     `json:"mdb_company_title_id"`
	CompanyTitle                string                    `json:"company_title"`
	CustomerName                string                    `json:"customer_name"`
	Address                     string                    `json:"address"`
	Address2                    string                    `json:"address_2"`
	Address3                    string                    `json:"address_3"`
	Hamlet                      string                    `json:"hamlet"`
	Neighbourhood               string                    `json:"neighbourhood"`
	CountryID                   int64                     `json:"country_id"`
	ProvinceID                  int64                     `json:"province_id"`
	DistrictID                  int64                     `json:"district_id"`
	SubDistrictID               int64                     `json:"sub_district_id"`
	UrbanVillageID              int64                     `json:"urban_village_id"`
	PostalCodeID                int64                     `json:"postal_code_id"`
	Longitude                   float64                   `json:"long"`
	Latitude                    float64                   `json:"lat"`
	PhoneCountryCode            string                    `json:"phone_country_code"`
	Phone                       string                    `json:"phone"`
	AlternativePhoneCountryCode string                    `json:"alternative_phone_country_code"`
	AlternativePhone            string                    `json:"alternative_phone"`
	Fax                         string                    `json:"fax"`
	CompanyEmail                string                    `json:"company_email"`
	AlternativeCompanyEmail     string                    `json:"alternative_company_email"`
	CustomerSource              string                    `json:"customer_source"`
	TaxName                     string                    `json:"tax_name"`
	TaxAddress                  string                    `json:"tax_address"`
	SalesmanID                  int64                     `json:"salesman_id"`
	RefCustomerID               int64                     `json:"ref_customer_id"`
	DistributorOF               string                    `json:"distributor_of"`
	CustomerGroupID             int64                     `json:"customer_group_id"`
	CustomerCategoryID          int64                     `json:"customer_category_id"`
	Status                      string                    `json:"status"`
	CustomerContact             []PreviousCustomerContact `json:"customer_contact"`
	UpdatedAt                   time.Time                 `json:"updated_at"`
	IsSuccess                   bool                      `json:"is_success"`
}

type PreviousCustomerContact struct {
	ID                 int64     `json:"id"`
	CustomerID         int64     `json:"customer_id"`
	MdbPersonProfileID int64     `json:"mdb_person_profile_id"`
	Nik                string    `json:"nik"`
	MdbPersonTitleID   int64     `json:"mdb_person_title_id"`
	PersonTitle        string    `json:"person_title"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Sex                string    `json:"sex"`
	Address            string    `json:"address"`
	Hamlet             string    `json:"hamlet"`
	Neighbourhood      string    `json:"neighbourhood"`
	ProvinceID         int64     `json:"province_id"`
	DistrictID         int64     `json:"district_id"`
	Phone              string    `json:"phone"`
	Email              string    `json:"email"`
	MdbPositionID      int64     `json:"mdb_position_id"`
	PositionName       string    `json:"position_name"`
	Status             string    `json:"status"`
	Action             int64     `json:"action"`
	UpdatedAt          time.Time `json:"updated_at"`
	IsSuccess          bool      `json:"is_success"`
}

type InternalGetListCustomerResponse struct {
	ID                  int64  `json:"id"`
	MDBCompanyProfileId int64  `json:"mdb_company_profile_id"`
	NPWP                string `json:"npwp"`
	IsPrincipal         bool   `json:"is_principal"`
	IsParent            bool   `json:"is_parent"`
	CompanyTitle        string `json:"company_title"`
	CustomerName        string `json:"customer_name"`
	Address             string `json:"address"`
	Phone               string `json:"phone"`
	CompanyEmail        string `json:"company_email"`
}

type InternalGetListDistributorResponse struct {
	ClientID            string    `json:"client_id"`
	AuthUserID          int64     `json:"auth_user_id"`
	MDBCompanyProfileId int64     `json:"mdb_company_profile_id"`
	ClientTypeName      string    `json:"client_type_name"`
	DistributorID       int64     `json:"distributor_id"`
	PrincipalID         int64     `json:"principal_id"`
	DistTitle           string    `json:"dist_title"`
	DistName            string    `json:"dist_name"`
	DistNPWP            string    `json:"dist_npwp"`
	DistAddress         string    `json:"dist_address"`
	DistHamlet          string    `json:"dist_hamlet"`
	DistNeighbourhood   string    `json:"dist_neighbourhood"`
	DistCountry         int64     `json:"dist_country"`
	DistProvince        int64     `json:"dist_province"`
	DistDistrict        int64     `json:"dist_district"`
	DistSubdistrict     int64     `json:"dist_subdistrict"`
	DistUrbanvillage    int64     `json:"dist_urbanvillage"`
	DistPostalcode      int64     `json:"dist_postalcode"`
	DistLicenseVariant  string    `json:"dist_license_variant"`
	Longitude           float64   `json:"longitude"`
	Latitude            float64   `json:"latitude"`
	DistPhone           string    `json:"dist_phone"`
	DistFax             string    `json:"dist_fax"`
	DistEmail           string    `json:"dist_email"`
	DistJoindate        time.Time `json:"dist_joindate"`
	DistFromdate        time.Time `json:"dist_fromdate"`
	DistExpirydate      time.Time `json:"dist_expirydate"`
	CompanyID           string    `json:"company_id"`
	BranchID            string    `json:"branch_id"`
	ActivationDate      time.Time `json:"activation_date"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type DetailErrorCustomerResponse struct {
	ID           int64  `json:"id"`
	NPWP         string `json:"npwp"`
	CustomerName string `json:"customer_name"`
}
