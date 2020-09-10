module ratelimit-operator

go 1.13

require (
	cloud.google.com/go v0.50.0 // indirect
	github.com/Azure/go-autorest/autorest v0.9.4 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.8.3 // indirect
	github.com/champly/lib4go v0.0.0-20200508051201-2cb0f4ccb079
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.3.1
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/operator-framework/operator-sdk v0.18.2
	github.com/solo-io/protoc-gen-ext v0.0.9
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.5.1
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/sys v0.0.0-20200810151505-1b9f1253b3ed // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	istio.io/api v0.0.0-20200521171657-32375f234cc1
	istio.io/client-go v0.0.0-20200518164621-ef682e2929e5
	istio.io/gogo-genproto v0.0.0-20200422223746-8166b73efbae // indirect
	k8s.io/api v0.18.2
	k8s.io/apiextensions-apiserver v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.6.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	k8s.io/client-go => k8s.io/client-go v0.18.2 // Required by prometheus-operator
)
