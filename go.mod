module github.com/openconfig/lemming

go 1.23.4

toolchain go1.23.8

require (
	cloud.google.com/go/cloudbuild v1.19.2
	cloud.google.com/go/logging v1.13.0
	cloud.google.com/go/monitoring v1.22.1
	cloud.google.com/go/trace v1.11.3
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.49.0
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.25.0
	github.com/fatih/color v1.15.0
	github.com/go-logr/logr v1.4.2
	github.com/golang/glog v1.2.4
	github.com/google/go-cmp v0.7.0
	github.com/google/gopacket v1.1.19
	github.com/google/uuid v1.6.0
	github.com/googleapis/gax-go/v2 v2.14.1
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.0.0
	github.com/kentik/patricia v1.2.1
	github.com/mdlayher/genetlink v1.3.2
	github.com/open-traffic-generator/snappi/gosnappi v1.5.1
	github.com/openconfig/gnmi v0.13.0
	github.com/openconfig/gnoi v0.4.1
	github.com/openconfig/gnoigo v0.0.0-20240320202954-ebd033e3542c
	github.com/openconfig/gnsi v1.6.0
	github.com/openconfig/goyang v1.6.0
	github.com/openconfig/gribi v1.9.0
	github.com/openconfig/gribigo v0.1.1-0.20250417173317-8f33380bbbf9
	github.com/openconfig/kne v0.1.18
	github.com/openconfig/magna v0.0.0-20240326180454-518e16696c84
	github.com/openconfig/ondatra v0.5.8
	github.com/openconfig/testt v0.0.0-20220311054427-efbb1a32ec07
	github.com/openconfig/ygnmi v0.11.1
	github.com/openconfig/ygot v0.31.0
	github.com/osrg/gobgp/v3 v3.35.0
	github.com/p4lang/p4runtime v1.4.0-rc.5.0.20220728214547-13f0d02a521e
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.9.1
	github.com/spf13/pflag v1.0.6
	github.com/spf13/viper v1.19.0
	github.com/stoewer/go-strcase v1.3.0
	github.com/vishvananda/netlink v1.2.1
	go.opentelemetry.io/contrib/detectors/gcp v1.32.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.58.0
	go.opentelemetry.io/otel v1.33.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.6.0
	go.opentelemetry.io/otel/log v0.6.0
	go.opentelemetry.io/otel/sdk v1.33.0
	go.opentelemetry.io/otel/sdk/log v0.6.0
	go.opentelemetry.io/otel/sdk/metric v1.33.0
	go.opentelemetry.io/otel/trace v1.33.0
	go.uber.org/mock v0.2.0
	golang.org/x/oauth2 v0.25.0
	golang.org/x/sys v0.31.0
	google.golang.org/api v0.216.0
	google.golang.org/genproto/googleapis/api v0.0.0-20250106144421-5f5ef82da422
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.5
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/apimachinery v0.29.3
	k8s.io/client-go v0.29.3
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.120.1
	modernc.org/cc/v4 v4.1.3
)

require (
	bitbucket.org/creachadair/stringset v0.0.14 // indirect
	cloud.google.com/go v0.118.0 // indirect
	cloud.google.com/go/auth v0.14.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.7 // indirect
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	cloud.google.com/go/iam v1.3.1 // indirect
	cloud.google.com/go/longrunning v0.6.4 // indirect
	cloud.google.com/go/pubsub v1.45.3 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.25.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.49.0 // indirect
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/aristanetworks/arista-ceoslab-operator/v2 v2.1.2 // indirect
	github.com/carlmontanari/difflibgo v0.0.0-20210718194309-31b9e131c298 // indirect
	github.com/creack/pty v1.1.18 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13 // indirect
	github.com/drivenets/cdnos-controller v1.7.4 // indirect
	github.com/eapache/channels v1.1.0 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/emicklei/go-restful/v3 v3.12.0 // indirect
	github.com/evanphx/json-patch/v5 v5.9.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.4 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/josharian/native v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/jstemmer/go-junit-report/v2 v2.1.0 // indirect
	github.com/k-sone/critbitgo v1.4.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/magiconair/properties v1.8.9 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mdlayher/netlink v1.7.2 // indirect
	github.com/mdlayher/socket v0.5.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f // indirect
	github.com/networkop/meshnet-cni v0.3.1-0.20230525201116-d7c306c635cf // indirect
	github.com/open-traffic-generator/keng-operator v0.3.28 // indirect
	github.com/openconfig/attestz v0.2.0 // indirect
	github.com/openconfig/gocloser v0.0.0-20220310182203-c6c950ed3b0b // indirect
	github.com/openconfig/lemming/operator v0.2.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/scrapli/scrapligo v1.1.11 // indirect
	github.com/scrapli/scrapligocfg v1.0.0 // indirect
	github.com/sirikothe/gotextfsm v1.0.1-0.20200816110946-6aa2cfd355e4 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/srl-labs/srl-controller v0.6.1 // indirect
	github.com/srl-labs/srlinux-scrapli v0.6.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/vishvananda/netns v0.0.4 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.55.0 // indirect
	go.opentelemetry.io/otel/metric v1.33.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/exp v0.0.0-20250218142911-aa4b98e5adaa // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/time v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20250106144421-5f5ef82da422 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.29.3 // indirect
	k8s.io/apiextensions-apiserver v0.29.3 // indirect
	k8s.io/kube-openapi v0.0.0-20240228011516-70dd3763d340 // indirect
	k8s.io/utils v0.0.0-20240310230437-4693a0247e57 // indirect
	lukechampine.com/uint128 v1.3.0 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/opt v0.1.3 // indirect
	modernc.org/strutil v1.1.3 // indirect
	modernc.org/token v1.1.0 // indirect
	sigs.k8s.io/controller-runtime v0.17.2 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
