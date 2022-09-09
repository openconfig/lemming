module github.com/openconfig/lemming

go 1.18

require (
	github.com/golang/glog v1.0.0
	github.com/golang/mock v1.6.0
	github.com/google/go-cmp v0.5.8
	github.com/google/gopacket v1.1.19
	github.com/h-fam/errdiff v1.0.2
	github.com/kentik/patricia v1.2.0
	github.com/openconfig/gnmi v0.0.0-20220617175856-41246b1b3507
	github.com/openconfig/gnoi v0.0.0-20220809151450-6bddacd72ef8
	github.com/openconfig/gnsi v0.0.0-20220906172358-1eda48d90de6
	github.com/openconfig/goyang v1.1.0
	github.com/openconfig/gribi v0.1.1-0.20220622162620-08d53dffce45
	github.com/openconfig/gribigo v0.0.0-20220802181317-805e943d8714
	github.com/openconfig/ondatra v0.0.0-20220902223518-87f933d5bfae
	github.com/openconfig/ygnmi v0.2.2
	github.com/openconfig/ygot v0.24.3
	github.com/osrg/gobgp/v3 v3.5.0
	github.com/p4lang/p4runtime v1.3.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.11.0
	github.com/vishvananda/netlink v1.1.1-0.20210330154013-f5de75959ad5
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad
	google.golang.org/grpc v1.49.0
	google.golang.org/protobuf v1.28.1
	k8s.io/klog/v2 v2.60.1
)

require (
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13 // indirect
	github.com/eapache/channels v1.1.0 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jstemmer/go-junit-report/v2 v2.0.1-0.20220823220451-7b10b4285462 // indirect
	github.com/k-sone/critbitgo v1.4.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/open-traffic-generator/snappi/gosnappi v0.8.5 // indirect
	github.com/openconfig/gocloser v0.0.0-20220310182203-c6c950ed3b0b // indirect
	github.com/openconfig/grpctunnel v0.0.0-20220524190229-125331eabdde // indirect
	github.com/openconfig/kne v0.1.4 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pelletier/go-toml/v2 v2.0.0-beta.8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/cobra v1.4.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/vishvananda/netns v0.0.0-20200728191858-db3c7e526aae // indirect
	golang.org/x/crypto v0.0.0-20220518034528-6f7dac969898 // indirect
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220608133413-ed9918b62aac // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	lukechampine.com/uint128 v1.1.1 // indirect
)

replace (
	github.com/openconfig/kne => ../kne
	github.com/openconfig/ondatra => ../ondatra
)