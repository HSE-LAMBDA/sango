package src

var GetIOBalancer func() *IOBalancer
var GetControllers func() []*Controller
var GetJbodControllers func() []*SANJBODController

func GetIOBalancerFactory(iob *IOBalancer) func() *IOBalancer {
	return func() *IOBalancer {
		return iob
	}
}


func GetControllersFactory(tvc []*Controller) func() []*Controller {
	return func() []*Controller {
		return tvc
	}
}


func GetJbodControllersFactory(tjc []*SANJBODController) func() []*SANJBODController {
	return func() []*SANJBODController {
		return tjc
	}
}
