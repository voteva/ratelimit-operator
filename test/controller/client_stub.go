package controller

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type StubClient struct {
	wrappedClient client.Client
	handler StubEventHandler
}

func (s StubClient) Get(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
	return s.wrappedClient.Get(ctx, key, obj)
}

func (s StubClient) List(ctx context.Context, list runtime.Object, opts ...client.ListOption) error {
	return s.wrappedClient.List(ctx, list, opts...)
}

func (s StubClient) Create(ctx context.Context, obj runtime.Object, opts ...client.CreateOption) error {
	if err := s.wrappedClient.Create(ctx, obj, opts...); err != nil {
		return err
	} else {

		s.handler.OnAdd(obj)
		return nil
	}
}

func (s StubClient) Delete(ctx context.Context, obj runtime.Object, opts ...client.DeleteOption) error {
	return s.wrappedClient.Delete(ctx, obj, opts...)
}

func (s StubClient) Update(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
	return s.wrappedClient.Update(ctx, obj, opts...)
}

func (s StubClient) Patch(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error {
	return s.wrappedClient.Patch(ctx, obj, patch, opts...)
}

func (s StubClient) DeleteAllOf(ctx context.Context, obj runtime.Object, opts ...client.DeleteAllOfOption) error {
	return s.wrappedClient.DeleteAllOf(ctx, obj, opts...)
}

func (s StubClient) Status() client.StatusWriter {
	return s.wrappedClient.Status()
}

func NewStubClient(clientScheme *runtime.Scheme, initObjs ...runtime.Object) client.Client {
	return &StubClient{wrappedClient: fake.NewFakeClientWithScheme(clientScheme, initObjs...)}
}
