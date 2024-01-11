package controller

type ApiGroup struct {
	CdsApi          ClusterController
	EdsApi          EndPointController
	EnvoyClusterApi EnvoyClusterController
	EnvoyNodeApi    EnvoyNodeController
}

var ApiGroupApp = new(ApiGroup)
