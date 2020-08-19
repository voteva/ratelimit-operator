package controller

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllertest"
)

//StubClient contains fake.fakeClient class as wrapped instance and adds event pushing to informers
type StubClient struct {
	wrappedClient client.Client
	scheme        *runtime.Scheme
	cache         cache.Cache
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
		s.pushToInformer(obj, func(fakeInformer *controllertest.FakeInformer, metaobj metav1.Object) {
			fakeInformer.Add(metaobj)
		})
		return nil
	}
}

func (s StubClient) Delete(ctx context.Context, obj runtime.Object, opts ...client.DeleteOption) error {
	if err := s.wrappedClient.Delete(ctx, obj, opts...); err != nil {
		return err
	} else {
		s.pushToInformer(obj, func(fakeInformer *controllertest.FakeInformer, metaobj metav1.Object) {
			fakeInformer.Delete(metaobj)
		})
		return nil
	}
}

func (s StubClient) Update(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
	return s.updateOrPatch(obj, ctx, func() error {
		return s.wrappedClient.Update(ctx, obj, opts...)
	})
}

func (s StubClient) Patch(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error {
	return s.updateOrPatch(obj, ctx, func() error {
		return s.wrappedClient.Patch(ctx, obj, patch, opts...)
	})
}

func (s StubClient) DeleteAllOf(ctx context.Context, obj runtime.Object, opts ...client.DeleteAllOfOption) error {
	panic("Not implemented")
}

func (s StubClient) Status() client.StatusWriter {
	return s.wrappedClient.Status()
}

func convertRuntimeToMeta(Object runtime.Object) metav1.Object {
	if converted, ok := Object.(metav1.Object); ok {
		return converted
	} else {
		return nil
	}
}

func (s StubClient) pushToInformer(obj runtime.Object, Action func(fakeInformer *controllertest.FakeInformer, metaobj metav1.Object)) {
	informer, _ := s.cache.GetInformer(context.TODO(), obj)
	if fInformer, ok := informer.(*controllertest.FakeInformer); ok {
		Action(fInformer, convertRuntimeToMeta(obj))
	}
}

func (s StubClient) pushUpdateToInformer(obj runtime.Object, oldObj runtime.Object, Action func(fakeInformer *controllertest.FakeInformer, metaNewObj metav1.Object, metaOldObj metav1.Object)) {
	informer, _ := s.cache.GetInformer(context.TODO(), obj)
	if fInformer, ok := informer.(*controllertest.FakeInformer); ok {
		Action(fInformer, convertRuntimeToMeta(obj), convertRuntimeToMeta(oldObj))
	}
}

func (s StubClient) updateOrPatch(obj runtime.Object, ctx context.Context, f func() error) error {
	var oldObject runtime.Object
	meta := convertRuntimeToMeta(obj)
	if oldObject, err := s.scheme.New(obj.GetObjectKind().GroupVersionKind()); err != nil {
		return err
	} else {
		if err := s.wrappedClient.Get(ctx, types.NamespacedName{Namespace: meta.GetNamespace(), Name: meta.GetName()}, oldObject); err != nil {
			return err
		}
	}

	if err := f(); err != nil {
		return err
	} else {
		s.pushUpdateToInformer(obj, oldObject, func(fakeInformer *controllertest.FakeInformer, metaNewObj metav1.Object, metaOldObj metav1.Object) {
			fakeInformer.Update(metaOldObj, metaNewObj)
		})
		return nil
	}
}

func NewStubClient(scheme *runtime.Scheme, cache *cache.Cache, initObjs ...runtime.Object) client.Client {
	if initObjs == nil {
		initObjs = []runtime.Object{}
	}
	return &StubClient{
		wrappedClient: fake.NewFakeClientWithScheme(scheme, initObjs...),
		scheme:        scheme,
		cache: *cache,
	}
}
