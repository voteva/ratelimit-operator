module ratelimit-operator

go 1.13

require (
	cloud.google.com/go v0.50.0 // indirect
	github.com/Azure/go-autorest/autorest v0.9.4 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.8.3 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/operator-framework/operator-sdk v0.18.2
	github.com/spf13/pflag v1.0.5
	istio.io/api v0.0.0-20200521171657-32375f234cc1
	istio.io/client-go v0.0.0-20200518164621-ef682e2929e5
	istio.io/gogo-genproto v0.0.0-20200422223746-8166b73efbae // indirect
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.6.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	k8s.io/client-go => k8s.io/client-go v0.18.2 // Required by prometheus-operator
)
