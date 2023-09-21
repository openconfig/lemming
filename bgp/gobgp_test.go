package bgp

import (
	"testing"

	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/ygot/ygot"
)

func TestValidatePrefixSetMode(t *testing.T) {
	tests := []struct {
		desc     string
		inConfig *oc.Root
		wantErr  bool
	}{{
		desc:     "nil config",
		inConfig: nil,
	}, {
		desc: "IPv4-no-prefixes",
		inConfig: func() *oc.Root {
			root := &oc.Root{}
			ps := root.GetOrCreateRoutingPolicy().GetOrCreateDefinedSets().GetOrCreatePrefixSet("foo")
			ps.SetMode(oc.PrefixSet_Mode_IPV4)
			return root
		}(),
	}, {
		desc: "IPv4",
		inConfig: func() *oc.Root {
			root := &oc.Root{}
			ps := root.GetOrCreateRoutingPolicy().GetOrCreateDefinedSets().GetOrCreatePrefixSet("foo")
			ps.SetMode(oc.PrefixSet_Mode_IPV4)
			if err := ps.AppendPrefix(&oc.RoutingPolicy_DefinedSets_PrefixSet_Prefix{
				IpPrefix:        ygot.String("1.1.1.1/32"),
				MasklengthRange: ygot.String("exact"),
			}); err != nil {
				t.Error(err)
			}
			if err := ps.AppendPrefix(&oc.RoutingPolicy_DefinedSets_PrefixSet_Prefix{
				IpPrefix:        ygot.String("2.2.2.2/32"),
				MasklengthRange: ygot.String("exact"),
			}); err != nil {
				t.Error(err)
			}
			return root
		}(),
	}, {
		desc: "IPv4-mode-unset",
		inConfig: func() *oc.Root {
			root := &oc.Root{}
			ps := root.GetOrCreateRoutingPolicy().GetOrCreateDefinedSets().GetOrCreatePrefixSet("foo")
			if err := ps.AppendPrefix(&oc.RoutingPolicy_DefinedSets_PrefixSet_Prefix{
				IpPrefix:        ygot.String("1.1.1.1/32"),
				MasklengthRange: ygot.String("exact"),
			}); err != nil {
				t.Error(err)
			}
			if err := ps.AppendPrefix(&oc.RoutingPolicy_DefinedSets_PrefixSet_Prefix{
				IpPrefix:        ygot.String("2.2.2.2/32"),
				MasklengthRange: ygot.String("exact"),
			}); err != nil {
				t.Error(err)
			}
			return root
		}(),
		wantErr: true,
	}, {
		desc: "IPv6",
		inConfig: func() *oc.Root {
			root := &oc.Root{}
			ps := root.GetOrCreateRoutingPolicy().GetOrCreateDefinedSets().GetOrCreatePrefixSet("foo")
			ps.SetMode(oc.PrefixSet_Mode_IPV6)
			if err := ps.AppendPrefix(&oc.RoutingPolicy_DefinedSets_PrefixSet_Prefix{
				IpPrefix:        ygot.String("2001::/32"),
				MasklengthRange: ygot.String("exact"),
			}); err != nil {
				t.Error(err)
			}
			if err := ps.AppendPrefix(&oc.RoutingPolicy_DefinedSets_PrefixSet_Prefix{
				IpPrefix:        ygot.String("2002::/32"),
				MasklengthRange: ygot.String("exact"),
			}); err != nil {
				t.Error(err)
			}
			return root
		}(),
	}, {
		desc: "mixed",
		inConfig: func() *oc.Root {
			root := &oc.Root{}
			ps := root.GetOrCreateRoutingPolicy().GetOrCreateDefinedSets().GetOrCreatePrefixSet("foo")
			ps.SetMode(oc.PrefixSet_Mode_MIXED)
			if err := ps.AppendPrefix(&oc.RoutingPolicy_DefinedSets_PrefixSet_Prefix{
				IpPrefix:        ygot.String("2001::/32"),
				MasklengthRange: ygot.String("exact"),
			}); err != nil {
				t.Error(err)
			}
			if err := ps.AppendPrefix(&oc.RoutingPolicy_DefinedSets_PrefixSet_Prefix{
				IpPrefix:        ygot.String("2.2.2.2/32"),
				MasklengthRange: ygot.String("exact"),
			}); err != nil {
				t.Error(err)
			}
			return root
		}(),
	}, {
		desc: "mixed-with-only-ipv4",
		inConfig: func() *oc.Root {
			root := &oc.Root{}
			ps := root.GetOrCreateRoutingPolicy().GetOrCreateDefinedSets().GetOrCreatePrefixSet("foo")
			ps.SetMode(oc.PrefixSet_Mode_MIXED)
			if err := ps.AppendPrefix(&oc.RoutingPolicy_DefinedSets_PrefixSet_Prefix{
				IpPrefix:        ygot.String("1.1.1.1/32"),
				MasklengthRange: ygot.String("exact"),
			}); err != nil {
				t.Error(err)
			}
			if err := ps.AppendPrefix(&oc.RoutingPolicy_DefinedSets_PrefixSet_Prefix{
				IpPrefix:        ygot.String("2.2.2.2/32"),
				MasklengthRange: ygot.String("exact"),
			}); err != nil {
				t.Error(err)
			}
			return root
		}(),
	}, {
		desc: "mixed-with-only-ipv6",
		inConfig: func() *oc.Root {
			root := &oc.Root{}
			ps := root.GetOrCreateRoutingPolicy().GetOrCreateDefinedSets().GetOrCreatePrefixSet("foo")
			ps.SetMode(oc.PrefixSet_Mode_MIXED)
			if err := ps.AppendPrefix(&oc.RoutingPolicy_DefinedSets_PrefixSet_Prefix{
				IpPrefix:        ygot.String("2001::/32"),
				MasklengthRange: ygot.String("exact"),
			}); err != nil {
				t.Error(err)
			}
			if err := ps.AppendPrefix(&oc.RoutingPolicy_DefinedSets_PrefixSet_Prefix{
				IpPrefix:        ygot.String("2002::/32"),
				MasklengthRange: ygot.String("exact"),
			}); err != nil {
				t.Error(err)
			}
			return root
		}(),
	}, {
		desc: "IPv4-not-pure",
		inConfig: func() *oc.Root {
			root := &oc.Root{}
			ps := root.GetOrCreateRoutingPolicy().GetOrCreateDefinedSets().GetOrCreatePrefixSet("foo")
			ps.SetMode(oc.PrefixSet_Mode_IPV4)
			if err := ps.AppendPrefix(&oc.RoutingPolicy_DefinedSets_PrefixSet_Prefix{
				IpPrefix:        ygot.String("2001::/32"),
				MasklengthRange: ygot.String("exact"),
			}); err != nil {
				t.Error(err)
			}
			if err := ps.AppendPrefix(&oc.RoutingPolicy_DefinedSets_PrefixSet_Prefix{
				IpPrefix:        ygot.String("2.2.2.2/32"),
				MasklengthRange: ygot.String("exact"),
			}); err != nil {
				t.Error(err)
			}
			return root
		}(),
		wantErr: true,
	}, {
		desc: "IPv6-not-pure",
		inConfig: func() *oc.Root {
			root := &oc.Root{}
			ps := root.GetOrCreateRoutingPolicy().GetOrCreateDefinedSets().GetOrCreatePrefixSet("foo")
			ps.SetMode(oc.PrefixSet_Mode_IPV6)
			if err := ps.AppendPrefix(&oc.RoutingPolicy_DefinedSets_PrefixSet_Prefix{
				IpPrefix:        ygot.String("2001::/32"),
				MasklengthRange: ygot.String("exact"),
			}); err != nil {
				t.Error(err)
			}
			if err := ps.AppendPrefix(&oc.RoutingPolicy_DefinedSets_PrefixSet_Prefix{
				IpPrefix:        ygot.String("2.2.2.2/32"),
				MasklengthRange: ygot.String("exact"),
			}); err != nil {
				t.Error(err)
			}
			return root
		}(),
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := validatePrefixSetMode(tt.inConfig)
			if gotErr := (err != nil); gotErr != tt.wantErr {
				t.Errorf("gotErr %v, wantErr: %v", err, tt.wantErr)
			}
		})
	}
}
