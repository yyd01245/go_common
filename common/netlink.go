package common

import (
	"fmt"
	"net"
	"errors"
	log "github.com/Sirupsen/logrus"	
	NT "github.com/vishvananda/netlink"
	// "github.com/vishvananda/netlink/nl"
	"golang.org/x/sys/unix"
)

const ADDROUTE = "add"
const DELROUTE = "del"

func LinkAddMacVlan(ifName string,parentIfname string) error{
	// list link 
	links, err := NT.LinkList()
	if err != nil {
		log.Errorf("Link list error: %v",err)
		return err
	}
	log.Infof("---links: %v",links)
	var parent NT.Link
	for _, l := range links {
		if l.Attrs().Name == parentIfname {
			// get parent link
			parent = l
		}
		if l.Attrs().Name == ifName {
			txt := fmt.Sprintf("ifname link:%v exsit",ifName)
			log.Infof(txt)
			return errors.New(txt)
		}
	}
	link := &NT.Macvlan{
		LinkAttrs: NT.LinkAttrs{Name: ifName, ParentIndex: parent.Attrs().Index},
		Mode:      NT.MACVLAN_MODE_BRIDGE,
	}
	log.Infof("---link:%v",link)

	if err := NT.LinkAdd(link); err != nil {
		log.Errorf("Link Add error: %v",err)
		return err
	}

	base := link.Attrs()

	result, err := NT.LinkByName(base.Name)
	if err != nil {
		log.Errorf("Link byname error: %v",err)
		return err
	}

	rBase := result.Attrs()

	if base.Index != 0 {
		if base.Index != rBase.Index {
			txt := fmt.Sprintf("index is %d, should be %d", rBase.Index, base.Index)
			log.Errorf(txt)
			return errors.New(txt)
		}
	}

	links, err = NT.LinkList()
	if err != nil {
		log.Errorf("Link list error: %v",err)
		return err
	}
	flag := false
	for _, l := range links {
		if l.Attrs().Name == link.Attrs().Name {
			log.Infof("Link macvlan properly:%v",l)
			flag = true
			break;
		}
	}
	if flag == false {
		log.Errorf("link macvlan add failed!!!")
		return errors.New("link macvlan add failed!!!")
	}
	// up 
	NT.LinkSetUp(link)
	return nil
}

func LinkDelMacVlan(ifName string) error{
	// list link 
	links, err := NT.LinkList()
	if err != nil {
		log.Errorf("Link list error: %v",err)
		return err
	}
	log.Infof("---links: %v",links)
	var link NT.Link
	flag := false
	for _, l := range links {

		if l.Attrs().Name == ifName {
			txt := fmt.Sprintf("ifname link:%v exsit",ifName)
			log.Infof(txt)
			link = l
			flag = true
			break;
		}
	}
	if !flag {
		txt := fmt.Sprintf("virtural macvlan:%v not exsit!!!",ifName)
		log.Infof(txt)
		return errors.New(txt)
	}
	if err := NT.LinkDel(link); err != nil {
		log.Errorf("Link Add error: %v",err)
		return err
	}

	links, err = NT.LinkList()
	if err != nil {
		log.Errorf("Link list error: %v",err)
		return err
	}
	flag = false
	for _, l := range links {
		if l.Attrs().Name == ifName {
			log.Infof("Link macvlan properly:%v",l)
			flag = true
			break;
		}
	}
	if flag {
		log.Errorf("link macvlan del failed!!!")
		return errors.New("link macvlan del failed!!!")
	}
	return nil
}

func GetLinkDevice(ifName string) (NT.Link,error) {
	var link NT.Link

	link, err := NT.LinkByName(ifName)
	if err != nil {
		log.Errorf("find link device:%v error: %v",ifName,err)
		return link,err
	}
	// bring the interface up
	if err = NT.LinkSetUp(link); err != nil {
		log.Errorf("setup link device:%v error: %v",ifName,err)
		return link,err
	}
	return link,nil
}

func GetRouteObject(link NT.Link,ipaddr string,outAddr string,tableID int) (*NT.Route,error) {
	_, dst, err := net.ParseCIDR(ipaddr)
	if err != nil {
		txt := fmt.Sprintf("check IP valid err:%v",err)
		log.Errorf(txt)
		return nil,errors.New(txt)
	}
	log.Debugf("dst cidr=%v",dst)
	src := net.ParseIP(outAddr)

	log.Debugf("src ip=%v",src)

	route := &NT.Route{
		LinkIndex: link.Attrs().Index,
		Dst:       dst,
		// Src:       src,
		Gw:       src,
		// Scope:     unix.RT_SCOPE_LINK,
		// Priority:  13,
		Table:     tableID,
		// Type:      unix.RTN_UNICAST,
		// Tos:       14,
	}
	log.Debugf("--- route: %v",route)
	return route,nil 
}

func GetDefaultRouteObject(link NT.Link,outAddr string,tableID int) (*NT.Route,error) {

	src := net.ParseIP(outAddr)

	log.Infof("src ip=%v",src)

	route := &NT.Route{
		LinkIndex: link.Attrs().Index,
		Dst:       nil,
		Src:       src,
		Gw:       nil,
		Scope:     unix.RT_SCOPE_LINK,
		// Priority:  13,
		Table:     tableID,
		// Type:      unix.RTN_UNICAST,
		// Tos:       14,
	}
	log.Infof("--- route: %v",route)
	return route,nil 
}

func NetAddOrDelRouteByLink(action string,route *NT.Route) error{
	
	// // add a gateway route
	if route == nil {
		txt := fmt.Sprintf("add route error: route is nil")
		log.Errorf(txt)
		return errors.New(txt)
	}
	if action == ADDROUTE {
		log.Debugf("---add route: %v",route)
		if err := NT.RouteAdd(route); err != nil {
			txt := fmt.Sprintf("add route error:%v",err)
			log.Warnf(txt)
			return errors.New(txt)
		}
	}else if action == DELROUTE {
		// list 
		// routes, err := NT.RouteListFiltered(NT.FAMILY_V4, route, NT.RT_FILTER_DST|NT.RT_FILTER_GW|NT.RT_FILTER_TABLE)
		// if err != nil {
		// 	txt := fmt.Sprintf("RouteListFiltered error:%v",err)
		// 	log.Errorf(txt)
		// 	return errors.New(txt)
		// }	
		// log.Infof("route list: %v",routes)
	
		// if len(routes) != 1 {
		// 	txt := fmt.Sprintf("Route not added properly:%v",routes)
		// 	log.Errorf(txt)
		// 	return errors.New(txt)
		// }
		log.Debugf("---delete route: %v",route)
		if err := NT.RouteDel(route); err != nil {
			txt := fmt.Sprintf("add route error:%v",err)
			log.Errorf(txt)
			return errors.New(txt)
		}
	}else {
		txt := "route unkown action type"
		log.Errorf(txt)
		return errors.New(txt)
	}
	return nil
}

// 批量添加路由，不进行匹配添加
func NetAddRoutePatch(link NT.Link,dstRoutes []string,outAddr string,tableID int) error{
	total := 0
	for _,ipAddr := range dstRoutes {
		// 
		route,err := GetRouteObject(link,ipAddr,outAddr,tableID)
		if err != nil {
			log.Errorf("get route from ipAddr:%s error!!",ipAddr)
			continue
		}
		err = NetAddOrDelRouteByLink("add",route)
		if err != nil {
			log.Errorf("add route:%v failed:%v!!",route,err)
			continue
		}
		total++
	}
	log.Infof("add route total:%v!",total)
	return nil
}

func NetSyncSopeLinkRouteTable(link NT.Link,srcTableID int,dstTableID int) error{
	total := 0
	// 不检查 default route,仅校验 scope link
	route := NT.Route{
		LinkIndex: link.Attrs().Index,
		Scope: 		unix.RT_SCOPE_LINK,
		Table:     srcTableID,
	}

	routes, err := NT.RouteListFiltered(NT.FAMILY_V4, &route, 
			NT.RT_FILTER_TABLE|NT.RT_FILTER_SCOPE)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered error:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}	
	
	dstRoute := NT.Route{
		LinkIndex: link.Attrs().Index,
		Scope: 		unix.RT_SCOPE_LINK,
		Table:     dstTableID,
	}

	dstRoutes, err := NT.RouteListFiltered(NT.FAMILY_V4, &dstRoute, 
			NT.RT_FILTER_TABLE|NT.RT_FILTER_SCOPE)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered error:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}	

	for _,R := range routes {
		flag := false
		if R.Dst == nil {
			// default
			continue
		}
		for index,r := range dstRoutes {
			if r.Dst == nil {
				continue
			}
			// 找到则 60.191.85.65/32 格式为带掩码
			if R.Dst.String() == r.Dst.String() {
				flag = true
				log.Infof("find route:%v, route:%v",r,R)
				dstRoutes = append(dstRoutes[:index],dstRoutes[index+1:]...)
				break
			}
		}
		if flag {
			continue
		}
		// log.Infof("delete route: %v!",R)
		R.Table = dstTableID 
		err = NetAddOrDelRouteByLink("add",&R)
		if err != nil {
			log.Errorf("add route:%v failed:%v!!",R,err)
			continue
		}
		total++
	}
	log.Infof("NetSyncSopeLinkRouteTable del route total:%v!",total)
	total = 0
	// 添加
	// log.Infof("last need to delete scope link route:%v,len=%v!",dstRoutes,len(dstRoutes))
	for _,route := range dstRoutes {
		if route.Dst == nil {
			continue
		}
		log.Infof("need to delete scope link  ip:%v",route)
		err = NetAddOrDelRouteByLink("del",&route)
		if err != nil {
			log.Errorf("add route:%v failed:%v!!",route,err)
			continue
		}
		total++
	}
	log.Infof("delete scope link route total:%v!",total)

	return nil
}

func RemoveRepByMap(slc []string) []string {
	result := []string{}
	tempMap := map[string]byte{}  // 存放不重复主键
	for _, e := range slc{
			l := len(tempMap)
			tempMap[e] = 0
			if len(tempMap) != l{  // 加入map后，map长度变化，则元素不重复
					result = append(result, e)
			}
	}
	return result
}

func NetSyncRoutePatch(link NT.Link,data []string,tableID int,scope int,gwIP string) error{
	total := 0
	gw := net.ParseIP(gwIP)

	dstRoutes := RemoveRepByMap(data)
	log.Infof("remoce repeat data: len=%d",len(dstRoutes))
	// 不能删除 default 和 scope link 以及保留的地址
	route := NT.Route{
		LinkIndex: link.Attrs().Index,
		Scope: 		NT.Scope(scope),
		Table:     tableID,
		Gw:			gw,	
	}

	routes, err := NT.RouteListFiltered(NT.FAMILY_V4, &route, 
			NT.RT_FILTER_TABLE|NT.RT_FILTER_SCOPE|NT.RT_FILTER_GW)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered error:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}	
	for _,R := range routes {
		flag := false
		for index,ipAddr := range dstRoutes {
			// 找到则 60.191.85.65/32 格式为带掩码
			if R.Dst == nil {
				continue
			}
			if R.Dst.String() == ipAddr {
				flag = true
				log.Debugf("find route:%v, route:%v",ipAddr,R)
				dstRoutes = append(dstRoutes[:index],dstRoutes[index+1:]...)
				break
			}
		}
		if flag {
			continue
		}
		// log.Infof("delete route: %v!",R)
		err = NetAddOrDelRouteByLink("del",&R)
		if err != nil {
			log.Errorf("add route:%v failed:%v!!",R,err)
			continue
		}
		total++
	}
	log.Infof("NetSyncRoutePatch del route total:%v!",total)
	total = 0
	// 添加
	log.Infof("last need to add route len=%v!",len(dstRoutes))
	for _,ipAddr := range dstRoutes {
		// 
		if ipAddr == "" {
			continue
		}
		log.Debugf("begin add ip:%v",ipAddr)
		route,err := GetRouteObject(link,ipAddr,gwIP,tableID)
		if err != nil {
			log.Errorf("get route from ipAddr:%s error!!",ipAddr)
			continue
		}
		err = NetAddOrDelRouteByLink("add",route)
		if err != nil {
			log.Warnf("add route:%v failed:%v!!",route,err)
			continue
		}
		total++
	}
	log.Infof("NetSyncRoutePatch add route total:%v!",total)

	return nil
}

func NetVerifyRoutePatch(link NT.Link,dstRoutes []string,tableID int,scope int,gwIP string) error{
	total := 0
	gw := net.ParseIP(gwIP)

	// 不能删除 default 和 scope link 以及保留的地址
	route := NT.Route{
		LinkIndex: link.Attrs().Index,
		Scope: 		NT.Scope(scope),
		Table:     tableID,
		Gw:			gw,	
	}

	routes, err := NT.RouteListFiltered(NT.FAMILY_V4, &route, 
			NT.RT_FILTER_TABLE|NT.RT_FILTER_SCOPE|NT.RT_FILTER_GW)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered error:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}	

	for _,R := range routes {
		flag := false
		for index,ipAddr := range dstRoutes {
			if R.Dst == nil {
				continue
			}
			// 找到则 60.191.85.65/32 格式为带掩码
			if R.Dst.String() == ipAddr {
				flag = true
				log.Infof("find route:%v, route:%v",ipAddr,R)
				dstRoutes = append(dstRoutes[:index],dstRoutes[index+1:]...)
				break
			}
		}
		if flag {
			continue
		}
	}
	total = 0
	// 添加
	log.Infof("NetVerifyRoutePatch last need to add route:%v,len=%v!",dstRoutes,len(dstRoutes))
	for _,ipAddr := range dstRoutes {
		// 
		if ipAddr == "" {
			continue
		}
		log.Infof("begin add ip:%v",ipAddr)
		route,err := GetRouteObject(link,ipAddr,gwIP,tableID)
		if err != nil {
			log.Errorf("get route from ipAddr:%s error!!",ipAddr)
			continue
		}
		err = NetAddOrDelRouteByLink("add",route)
		if err != nil {
			log.Errorf("add route:%v failed:%v!!",route,err)
			continue
		}
		total++
	}
	log.Infof("NetVerifyRoutePatch add route total:%v!",total)

	return nil
}

// 批量删除路由，进行匹配需要删除的数据，排查 scope
func NetDelRoutePatch(link NT.Link,dstRoutes []string,tableID int,scope int,gwIP string) error{
	total := 0
	gw := net.ParseIP(gwIP)

	// 不能删除 default 和 scope link 以及保留的地址
	route := NT.Route{
		LinkIndex: link.Attrs().Index,
		Scope: 		NT.Scope(scope),
		Table:     tableID,
		Gw:			gw,	
	}

	routes, err := NT.RouteListFiltered(NT.FAMILY_V4, &route, 
			NT.RT_FILTER_TABLE|NT.RT_FILTER_SCOPE|NT.RT_FILTER_GW)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered error:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}	
	for _,R := range routes {
		flag := false
		for index,ipAddr := range dstRoutes {
			// 找到则 60.191.85.65/32 格式为带掩码
			if R.Dst == nil {
				continue
			}
			if R.Dst.String() == ipAddr {
				flag = true
				log.Infof("find route:%v, route:%v",ipAddr,R)
				dstRoutes = append(dstRoutes[:index],dstRoutes[index+1:]...)
				break
			}
		}
		if !flag {
			continue
		}
		// log.Infof("delete route: %v!",R)
		err = NetAddOrDelRouteByLink("del",&R)
		if err != nil {
			log.Errorf("add route:%v failed:%v!!",R,err)
			continue
		}
		total++
	}


	log.Infof("NetDelRoutePatch del route total:%v!",total)
	return nil
}

func NetFindRouteByLink(link NT.Link,tableID int,route *NT.Route) error{


	routes, err := NT.RouteListFiltered(NT.FAMILY_V4, route, 
			NT.RT_FILTER_TABLE|NT.RT_FILTER_SCOPE|NT.RT_FILTER_DST|NT.RT_FILTER_SRC|NT.RT_FILTER_GW)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered error:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}	
	log.Infof("route list: %v",routes)
	log.Infof("route list len: %v",len(routes))
	if len(routes) != 1 {
		txt := fmt.Sprintf("RouteListFiltered failed:%v",route)
		log.Errorf(txt)
		return errors.New(txt)
	}
	return nil
}


func NetListRouteByLink(link NT.Link,tableID int,scope int,gwIP string) error{

	// add a gateway route
	gw := net.ParseIP(gwIP)
	route := NT.Route{
		LinkIndex: link.Attrs().Index,
		Scope: 		NT.Scope(scope),
		Table:     tableID,
		Gw:				gw,
	}

	routes, err := NT.RouteListFiltered(NT.FAMILY_V4, &route, 
			NT.RT_FILTER_TABLE|NT.RT_FILTER_SCOPE|NT.RT_FILTER_GW)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered error:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}	
	log.Infof("route list: %v",routes)
	log.Infof("route list len: %v",len(routes))
	
	return nil
}
