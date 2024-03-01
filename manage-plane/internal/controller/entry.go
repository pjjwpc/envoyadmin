package controller

type ApiGroup struct {
	CdsApi          ClusterController
	EdsApi          EndPointController
	EnvoyClusterApi EnvoyClusterController
	EnvoyNodeApi    EnvoyNodeController
	SysManagerApi   SysManagerController
}

var ApiGroupApp = new(ApiGroup)
