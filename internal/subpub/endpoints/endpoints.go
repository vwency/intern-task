package endpoints

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/intern-task/internal/subpub/service"
)

type Endpoints struct {
	Subscribe   endpoint.Endpoint
	Publish     endpoint.Endpoint
	Unsubscribe endpoint.Endpoint
}

func MakeEndpoints(s *service.SubPubService) Endpoints {
	return Endpoints{
		Subscribe:   makeSubscribeEndpoint(s),
		Publish:     makePublishEndpoint(s),
		Unsubscribe: makeUnsubscribeEndpoint(s),
	}
}
