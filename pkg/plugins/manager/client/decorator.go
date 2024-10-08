package client

import (
	"context"
	"errors"

	"github.com/grafana/grafana-plugin-sdk-go/backend"

	"github.com/grafana/grafana/pkg/plugins"
)

// Decorator allows a plugins.Client to be decorated with middlewares.
type Decorator struct {
	client      plugins.Client
	middlewares []plugins.ClientMiddleware
}

var (
	_ = plugins.Client(&Decorator{})
)

// NewDecorator creates a new plugins.client decorator.
func NewDecorator(client plugins.Client, middlewares ...plugins.ClientMiddleware) (*Decorator, error) {
	if client == nil {
		return nil, errors.New("client cannot be nil")
	}

	return &Decorator{
		client:      client,
		middlewares: middlewares,
	}, nil
}

func (d *Decorator) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	if req == nil {
		return nil, errNilRequest
	}
	ctx = backend.WithEndpoint(ctx, backend.EndpointQueryData)
	ctx = backend.WithPluginContext(ctx, req.PluginContext)
	ctx = backend.WithUser(ctx, req.PluginContext.User)

	client := clientFromMiddlewares(d.middlewares, d.client)

	return client.QueryData(ctx, req)
}

func (d *Decorator) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	if req == nil {
		return errNilRequest
	}

	ctx = backend.WithEndpoint(ctx, backend.EndpointCallResource)
	ctx = backend.WithPluginContext(ctx, req.PluginContext)
	ctx = backend.WithUser(ctx, req.PluginContext.User)

	if sender == nil {
		return errors.New("sender cannot be nil")
	}

	client := clientFromMiddlewares(d.middlewares, d.client)
	return client.CallResource(ctx, req, sender)
}

func (d *Decorator) CollectMetrics(ctx context.Context, req *backend.CollectMetricsRequest) (*backend.CollectMetricsResult, error) {
	if req == nil {
		return nil, errNilRequest
	}

	ctx = backend.WithEndpoint(ctx, backend.EndpointCollectMetrics)
	ctx = backend.WithPluginContext(ctx, req.PluginContext)
	ctx = backend.WithUser(ctx, req.PluginContext.User)

	client := clientFromMiddlewares(d.middlewares, d.client)
	return client.CollectMetrics(ctx, req)
}

func (d *Decorator) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	if req == nil {
		return nil, errNilRequest
	}

	ctx = backend.WithEndpoint(ctx, backend.EndpointCheckHealth)
	ctx = backend.WithPluginContext(ctx, req.PluginContext)
	ctx = backend.WithUser(ctx, req.PluginContext.User)

	client := clientFromMiddlewares(d.middlewares, d.client)
	return client.CheckHealth(ctx, req)
}

func (d *Decorator) SubscribeStream(ctx context.Context, req *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
	if req == nil {
		return nil, errNilRequest
	}

	ctx = backend.WithEndpoint(ctx, backend.EndpointSubscribeStream)
	ctx = backend.WithPluginContext(ctx, req.PluginContext)
	ctx = backend.WithUser(ctx, req.PluginContext.User)

	client := clientFromMiddlewares(d.middlewares, d.client)
	return client.SubscribeStream(ctx, req)
}

func (d *Decorator) PublishStream(ctx context.Context, req *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
	if req == nil {
		return nil, errNilRequest
	}

	ctx = backend.WithEndpoint(ctx, backend.EndpointPublishStream)
	ctx = backend.WithPluginContext(ctx, req.PluginContext)
	ctx = backend.WithUser(ctx, req.PluginContext.User)

	client := clientFromMiddlewares(d.middlewares, d.client)
	return client.PublishStream(ctx, req)
}

func (d *Decorator) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
	if req == nil {
		return errNilRequest
	}

	ctx = backend.WithEndpoint(ctx, backend.EndpointRunStream)
	ctx = backend.WithPluginContext(ctx, req.PluginContext)
	ctx = backend.WithUser(ctx, req.PluginContext.User)

	if sender == nil {
		return errors.New("sender cannot be nil")
	}

	client := clientFromMiddlewares(d.middlewares, d.client)
	return client.RunStream(ctx, req, sender)
}

func (d *Decorator) ValidateAdmission(ctx context.Context, req *backend.AdmissionRequest) (*backend.ValidationResponse, error) {
	if req == nil {
		return nil, errNilRequest
	}

	ctx = backend.WithEndpoint(ctx, backend.EndpointValidateAdmission)
	ctx = backend.WithPluginContext(ctx, req.PluginContext)
	ctx = backend.WithUser(ctx, req.PluginContext.User)

	client := clientFromMiddlewares(d.middlewares, d.client)
	return client.ValidateAdmission(ctx, req)
}

func (d *Decorator) MutateAdmission(ctx context.Context, req *backend.AdmissionRequest) (*backend.MutationResponse, error) {
	if req == nil {
		return nil, errNilRequest
	}

	ctx = backend.WithEndpoint(ctx, backend.EndpointMutateAdmission)
	ctx = backend.WithPluginContext(ctx, req.PluginContext)
	ctx = backend.WithUser(ctx, req.PluginContext.User)

	client := clientFromMiddlewares(d.middlewares, d.client)
	return client.MutateAdmission(ctx, req)
}

func (d *Decorator) ConvertObjects(ctx context.Context, req *backend.ConversionRequest) (*backend.ConversionResponse, error) {
	if req == nil {
		return nil, errNilRequest
	}

	ctx = backend.WithEndpoint(ctx, backend.EndpointConvertObject)
	ctx = backend.WithPluginContext(ctx, req.PluginContext)
	ctx = backend.WithUser(ctx, req.PluginContext.User)

	client := clientFromMiddlewares(d.middlewares, d.client)
	return client.ConvertObjects(ctx, req)
}

func clientFromMiddlewares(middlewares []plugins.ClientMiddleware, finalClient plugins.Client) plugins.Client {
	if len(middlewares) == 0 {
		return finalClient
	}

	reversed := reverseMiddlewares(middlewares)
	next := finalClient

	for _, m := range reversed {
		next = m.CreateClientMiddleware(next)
	}

	return next
}

func reverseMiddlewares(middlewares []plugins.ClientMiddleware) []plugins.ClientMiddleware {
	reversed := make([]plugins.ClientMiddleware, len(middlewares))
	copy(reversed, middlewares)

	for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}

	return reversed
}
