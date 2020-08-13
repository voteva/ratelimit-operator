package controller

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type StubClient struct {
	wrappedClient client.Client
}

func (s StubClient) Get(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
	panic("implement me")
}

func (s StubClient) List(ctx context.Context, list runtime.Object, opts ...client.ListOption) error {
	panic("implement me")
}

func (s StubClient) Create(ctx context.Context, obj runtime.Object, opts ...client.CreateOption) error {
	panic("implement me")
}

func (s StubClient) Delete(ctx context.Context, obj runtime.Object, opts ...client.DeleteOption) error {
	panic("implement me")
}

func (s StubClient) Update(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
	panic("implement me")
}

func (s StubClient) Patch(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error {
	panic("implement me")
}

func (s StubClient) DeleteAllOf(ctx context.Context, obj runtime.Object, opts ...client.DeleteAllOfOption) error {
	panic("implement me")
}

func (s StubClient) Status() client.StatusWriter {
	panic("implement me")
}

func NewStubClient(clientScheme *runtime.Scheme, initObjs ...runtime.Object) client.Client {
	return &StubClient{wrappedClient: fake.NewFakeClientWithScheme(clientScheme, initObjs...)}
}
