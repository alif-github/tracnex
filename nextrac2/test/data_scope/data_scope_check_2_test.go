package data_scope_test

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/service"
	"testing"
)

func getMap_2() map[string]interface{} {
	test := make(map[string]interface{})
	test2 := make(map[string][]string)
	var test3 []map[string][]string

	test2["nexsoft.asset_owner_id"] = []string{"2"}
	test2["nexsoft.country.region.province"] = []string{"2"}

	test3 = append(test3, test2)

	test["nexsoft.btl_account_id"] = []string{"3"}
	test["nexsoft.asset_owner_id"] = []string{"2"}
	test["nexsoft.country.region.province"] = []string{"all"}

	test["nexsoft.asset_owner_id._and"] = test3
	return test
}

func TestCase1_2(testing *testing.T) {
	dataScope := getMap_2()
	needScope := []string{"nexsoft.btl_account_id"}
	expectedResult := "{\"nexsoft.btl_account_id\":[\"3\"]}"

	//fmt.Println(util.StructToJSON(service.CheckScope(dataScope, needScope)))
	if util.StructToJSON(service.CheckScope(dataScope, needScope)) != expectedResult {
		testing.Errorf("Not Expected Result At Test Case 1_2")
	}
}

func TestCase2_2(testing *testing.T) {
	dataScope := getMap_2()
	needScope := []string{"nexsoft.asset_owner_id"}
	expectedResult := "{\"nexsoft.asset_owner_id\":[\"2\"]}"

	//fmt.Println(util.StructToJSON(service.CheckScope(dataScope, needScope)))
	if util.StructToJSON(service.CheckScope(dataScope, needScope)) != expectedResult {
		testing.Errorf("Not Expected Result At Test Case 2_2")
	}
}

func TestCase3_2(testing *testing.T) {
	dataScope := getMap_2()
	needScope := []string{"nexsoft.btl_account_id", "nexsoft.asset_owner_id"}
	expectedResult := "{\"nexsoft.asset_owner_id\":[\"2\"],\"nexsoft.btl_account_id\":[\"3\"]}"

	//fmt.Println(util.StructToJSON(service.CheckScope(dataScope, needScope)))
	if util.StructToJSON(service.CheckScope(dataScope, needScope)) != expectedResult {
		testing.Errorf("Not Expected Result At Test Case 3_2")
	}
}

func TestCase4_2(testing *testing.T) {
	dataScope := getMap_2()
	needScope := []string{"nexsoft.btl_account_id", "nexsoft.asset_owner_id", "nexsoft.country.region.province"}
	expectedResult := "{\"_and\":[{\"nexsoft.asset_owner_id\":[\"2\"],\"nexsoft.country.region.province\":[\"2\"]}],\"nexsoft.btl_account_id\":[\"3\"]}"

	//fmt.Println(util.StructToJSON(service.CheckScope(dataScope, needScope)))
	if util.StructToJSON(service.CheckScope(dataScope, needScope)) != expectedResult {
		testing.Errorf("Not Expected Result At Test Case 4_2")
	}
}

func TestCase5_2(testing *testing.T) {
	dataScope := getMap_2()
	needScope := []string{"nexsoft.asset_owner_id", "nexsoft.country.region.province"}
	expectedResult := "{\"_and\":[{\"nexsoft.asset_owner_id\":[\"2\"],\"nexsoft.country.region.province\":[\"2\"]}]}"

	//fmt.Println(util.StructToJSON(service.CheckScope(dataScope, needScope)))
	if util.StructToJSON(service.CheckScope(dataScope, needScope)) != expectedResult {
		testing.Errorf("Not Expected Result At Test Case 5_2")
	}
}

func TestCase6_2(testing *testing.T) {
	dataScope := getMap_2()
	needScope := []string{"nexsoft.btl_account_id", "nexsoft.country.region.province"}
	expectedResult := "{\"nexsoft.btl_account_id\":[\"3\"],\"nexsoft.country.region.province\":[\"all\"]}"

	//fmt.Println(util.StructToJSON(service.CheckScope(dataScope, needScope)))
	if util.StructToJSON(service.CheckScope(dataScope, needScope)) != expectedResult {
		testing.Errorf("Not Expected Result At Test Case 6_2")
	}
}

func TestCase7_2(testing *testing.T) {
	dataScope := getMap_2()
	needScope := []string{"nexsoft.country.region.province"}
	expectedResult := "{\"nexsoft.country.region.province\":[\"all\"]}"

	//fmt.Println(util.StructToJSON(service.CheckScope(dataScope, needScope)))
	if util.StructToJSON(service.CheckScope(dataScope, needScope)) != expectedResult {
		testing.Errorf("Not Expected Result At Test Case 7_2")
	}
}
