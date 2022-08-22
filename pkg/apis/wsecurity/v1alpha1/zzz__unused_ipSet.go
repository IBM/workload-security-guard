package v1alpha1

/*
func (in CidrSet) DeepCopyInto(out *CidrSet) {
	copy((*out), in)
	cpy := make([]net.IPNet, len(in))
	for i, v := range in {
		copy(cpy[i].IP, v.IP)
		copy(cpy[i].Mask, v.Mask)
	}
	(*out) = cpy
}
*/

/*
// reconcile staff
func GetCidrsFromList(list []string) CidrSet {
	cidrs := make([]net.IPNet, 0, len(list))
	for _, v := range list {
		_, ipNet, err := net.ParseCIDR(v)
		if err == nil {
			var cidr net.IPNet
			//fmt.Printf("CIDRS found CIDR %v in GetCidrsFromList\n", ipNet)
			//cidrs = append(cidrs, *ipNet)
			cidr.IP = make(net.IP, len(ipNet.IP))
			cidr.Mask = make(net.IPMask, len(ipNet.Mask))
			copy(cidr.IP, ipNet.IP)
			copy(cidr.Mask, ipNet.Mask)
			cidrs = append(cidrs, cidr)
			continue
		}
		fmt.Printf("Ilegal cidr %s is skipped during GetCidrsFromList - %s\n", v, err.Error())
	}
	//if len(cidrs) > 0 {
	//fmt.Printf("GetCidrsFromList CIDRS %v\n", cidrs)
	//}
	return cidrs

}
*/
