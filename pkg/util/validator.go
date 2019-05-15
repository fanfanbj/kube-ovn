package util

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func ValidateLogicalSwitch(annotations map[string]string) error {
	cidrStr := annotations[CidrAnnotation]
	if cidrStr == "" {
		return fmt.Errorf("cidr is required for logical switch")
	}
	_, cidr, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return fmt.Errorf("%s is a invalid cidr %v", cidrStr, err)
	}

	gatewayStr := annotations[GatewayAnnotation]
	if gatewayStr == "" {
		return fmt.Errorf("gateway is required for logical switch")
	}
	gateway := net.ParseIP(gatewayStr)
	if gateway == nil {
		return fmt.Errorf("%s  is not a valid gateway", gatewayStr)
	}
	if !cidr.Contains(gateway) {
		return fmt.Errorf("gateway address %s not in cidr range", gatewayStr)
	}

	excludeIps := annotations[ExcludeIpsAnnotation]
	if excludeIps != "" {
		ipRanges := strings.Split(excludeIps, " ")
		ips := []string{}
		for _, ipr := range ipRanges {
			ips = append(ips, strings.Split(ipr, "..")...)
		}
		for _, ip := range ips {
			if net.ParseIP(ip) == nil {
				return fmt.Errorf("ip %s in exclude_ips is not a valid address", ip)
			}
		}
	}

	private := annotations[PrivateSwitchAnnotation]
	if private != "" && private != "true" && private != "false" {
		return fmt.Errorf("%s can only be \"true\" or \"false\"", PrivateSwitchAnnotation)
	}

	allow := annotations[AllowAccessAnnotation]
	if allow != "" {
		for _, cidr := range strings.Split(allow, ",") {
			if _, _, err := net.ParseCIDR(cidr); err != nil {
				return fmt.Errorf("%s in %s is not a valid address", cidr, AllowAccessAnnotation)
			}
		}
	}

	gwType := annotations[GWTypeAnnotation]
	if gwType != "" {
		if gwType != GWDistributedMode && gwType != GWCentralizedMode {
			return fmt.Errorf("%s is not a valid %s", gwType, GWTypeAnnotation)
		}
	}

	return nil
}

func ValidatePodNetwork(annotations map[string]string) error {
	if ip := annotations[IpAddressAnnotation]; ip != "" {
		_, _, err := net.ParseCIDR(ip)
		if err != nil {
			return fmt.Errorf("%s is not a valid %s", ip, IpAddressAnnotation)
		}
	}

	mac := annotations[MacAddressAnnotation]
	if mac != "" {
		if _, err := net.ParseMAC(mac); err != nil {
			return fmt.Errorf("%s is not a valid %s", mac, MacAddressAnnotation)
		}
	}

	ipPool := annotations[IpPoolAnnotation]
	if ipPool != "" {
		for _, ip := range strings.Split(ipPool, ",") {
			if net.ParseIP(ip) == nil {
				return fmt.Errorf("%s in %s is not a valid address", ip, IpPoolAnnotation)
			}
		}
	}

	ingress := annotations[IngressRateAnnotation]
	if ingress != "" {
		if _, err := strconv.Atoi(ingress); err != nil {
			return fmt.Errorf("%s is not a valid %s", ingress, IngressRateAnnotation)
		}
	}

	egress := annotations[EgressRateAnnotation]
	if egress != "" {
		if _, err := strconv.Atoi(ingress); err != nil {
			return fmt.Errorf("%s is not a valid %s", ingress, EgressRateAnnotation)
		}
	}

	return nil
}
