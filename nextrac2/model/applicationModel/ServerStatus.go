package applicationModel

import "nexsoft.co.id/nexcommon/util"

type ServerStatus struct {
	Status string `json:"status"`
	Redis string `json:"redis"`
	Database string `json:"database"`
	//ElasticSearch string `json:"elasticsearch"`
}

func (object ServerStatus) String() string {
	return util.StructToJSON(object)
}
