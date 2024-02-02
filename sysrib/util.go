package sysrib

import (
	"fmt"
	"net/netip"
)

func canonicalPrefix(prefix string) (netip.Prefix, error) {
	pfx, err := netip.ParsePrefix(prefix)
	if err != nil {
		// TODO(wenbli): This should be a Reconciler.Validate checker function.
		return netip.Prefix{}, fmt.Errorf("BGP GUE Policy prefix cannot be parsed: %v", err)
	}
	return pfx.Masked(), nil // make canonical
}

// PrintRIB prints the IPv4 and IPv6 RIBs.
func (s *SysRIB) PrintRIB() {
	fmt.Println("RIB DUMP")
	s.mu.RLock()
	defer s.mu.RUnlock()
	for niName, ni := range s.NI {
		fmt.Printf("================== network-instance: %q ==================\n", niName)
		for it := ni.IPV4.Iterate(); it.Next(); {
			for _, route := range it.Tags() {
				fmt.Println(route)
			}
			fmt.Println("------------------")
		}
		for it := ni.IPV6.Iterate(); it.Next(); {
			for _, route := range it.Tags() {
				fmt.Println(route)
			}
			fmt.Println("------------------")
		}
		fmt.Printf("================== END network-instance: %q ==================\n", niName)
	}
}

// PrintProgrammedRoutes prints the programmed routes.
func (s *Server) PrintProgrammedRoutes() {
	fmt.Println("Programmed Routes DUMP (From Sysrib)")
	s.programmedRoutesMu.Lock()
	defer s.programmedRoutesMu.Unlock()
	for _, route := range s.programmedRoutes {
		fmt.Printf("%+v\n", route)
		fmt.Println("------------------")
	}
}
