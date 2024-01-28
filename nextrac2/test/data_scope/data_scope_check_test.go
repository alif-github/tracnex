package data_scope_test

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/service"
	"testing"
)

func getMap() map[string]interface{} {
	test := make(map[string]interface{})
	test2 := make(map[string][]string)
	var test3 []map[string][]string

	test2["nexsoft.asset_owner_id"] = []string{"2"}
	test2["nexsoft.country.region.province"] = []string{"2"}

	test3 = append(test3, test2)

	test["nexsoft.btl_account_id"] = []string{"3"}
	test["nexsoft.asset_owner_id"] = []string{"2"}
	test["nexsoft.asset_owner_id._and"] = test3
	return test
}

func TestCase1(testing *testing.T) {
	dataScope := getMap()

	needScope := []string{"nexsoft.btl_account_id"}
	expectedResult := "{\"nexsoft.btl_account_id\":[\"3\"]}"

	//fmt.Println(util.StructToJSON(service.CheckScope(dataScope, needScope)))
	if util.StructToJSON(service.CheckScope(dataScope, needScope)) != expectedResult {
		testing.Errorf("Not Expected Result At Test Case 1")
	}
}

func TestCase2(testing *testing.T) {
	dataScope := getMap()
	needScope := []string{"nexsoft.asset_owner_id"}
	expectedResult := "{\"nexsoft.asset_owner_id\":[\"2\"]}"

	//fmt.Println(util.StructToJSON(service.CheckScope(dataScope, needScope)))
	if util.StructToJSON(service.CheckScope(dataScope, needScope)) != expectedResult {
		testing.Errorf("Not Expected Result At Test Case 2")
	}
}

func TestCase3(testing *testing.T) {
	dataScope := getMap()
	needScope := []string{"nexsoft.btl_account_id", "nexsoft.asset_owner_id"}
	expectedResult := "{\"nexsoft.asset_owner_id\":[\"2\"],\"nexsoft.btl_account_id\":[\"3\"]}"

	//fmt.Println(util.StructToJSON(service.CheckScope(dataScope, needScope)))
	if util.StructToJSON(service.CheckScope(dataScope, needScope)) != expectedResult {
		testing.Errorf("Not Expected Result At Test Case 3")
	}
}

func TestCase4(testing *testing.T) {
	dataScope := getMap()
	needScope := []string{"nexsoft.btl_account_id", "nexsoft.asset_owner_id", "nexsoft.country.region.province"}
	expectedResult := "{\"_and\":[{\"nexsoft.asset_owner_id\":[\"2\"],\"nexsoft.country.region.province\":[\"2\"]}],\"nexsoft.btl_account_id\":[\"3\"]}"

	//fmt.Println(util.StructToJSON(service.CheckScope(dataScope, needScope)))
	if util.StructToJSON(service.CheckScope(dataScope, needScope)) != expectedResult {
		testing.Errorf("Not Expected Result At Test Case 4")
	}
}

func TestCase5(testing *testing.T) {
	dataScope := getMap()
	needScope := []string{"nexsoft.asset_owner_id", "nexsoft.country.region.province"}
	expectedResult := "{\"_and\":[{\"nexsoft.asset_owner_id\":[\"2\"],\"nexsoft.country.region.province\":[\"2\"]}]}"

	//fmt.Println(util.StructToJSON(service.CheckScope(dataScope, needScope)))
	if util.StructToJSON(service.CheckScope(dataScope, needScope)) != expectedResult {
		testing.Errorf("Not Expected Result At Test Case 5")
	}
}

func TestCase6(testing *testing.T) {
	dataScope := getMap()
	needScope := []string{"nexsoft.btl_account_id", "nexsoft.country.region.province"}
	expectedResult := "null"

	//fmt.Println(util.StructToJSON(service.CheckScope(dataScope, needScope)))
	if util.StructToJSON(service.CheckScope(dataScope, needScope)) != expectedResult {
		testing.Errorf("Not Expected Result At Test Case 6")
	}
}

func TestCase7(testing *testing.T) {
	dataScope := getMap()
	needScope := []string{"nexsoft.country.region.province"}
	expectedResult := "null"

	//fmt.Println(util.StructToJSON(service.CheckScope(dataScope, needScope)))
	if util.StructToJSON(service.CheckScope(dataScope, needScope)) != expectedResult {
		testing.Errorf("Not Expected Result At Test Case 7")
	}
}
