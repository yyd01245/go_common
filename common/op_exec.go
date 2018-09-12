package common

import (
	// "fmt"
	"net"
	"fmt"
	// "errors"
	"strconv"
	"strings"
	"os/exec"
	"bytes"
	"math/big"
	log "github.com/Sirupsen/logrus"

)


func ExecScripts(binPath string,args []string) (string, error) {
	stderr := bytes.NewBuffer(nil)
	stdout := bytes.NewBuffer(nil)

	cmd := exec.Command(binPath, args...)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		log.Debugf("exec command: %v  args: %v err=%v",binPath,
				args,err)
	}
	// log.Infof("exec command: ",binPath," args: %v",
	// 	args," success")	
	outputErr := string(stderr.Bytes())
	if len(outputErr) > 0 {
		log.Debugf("exec command: stderr: %v",outputErr)
	}
	output := string(stdout.Bytes())
	if len(output) > 0 {
		log.Debugf("exec command: stdout: ",output)

	}	
	return output,err
}

func ExecScriptsPipe(binPath string,args []string) (string, error) {
	stderr := bytes.NewBuffer(nil)
	stdout := bytes.NewBuffer(nil)
	strArgs := ""
	for _,value := range args{
		strArgs += (" " + value)
	}
	input := fmt.Sprintf(`%s%s`,binPath,strArgs)
	cmd := exec.Command("/bin/sh", "-c", input)
	// log.Infof("pipe exec input:%s",input)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		log.Warnf("exec command: /bin/sh -c %v  args: %v err=%v",binPath,
			input,err)
	}
	// log.Infof("exec command: /bin/sh -c  ",binPath," args: %v",
	// 	input," success")	
	outputErr := string(stderr.Bytes())
	if len(outputErr) > 0 {
		log.Debugf("exec command: stderr: %v",outputErr)
	}
	output := string(stdout.Bytes())
	if len(output) > 0 {
		log.Debugf("exec command: stdout: ",output)

	}	
	return output,err
}

func GetUciValue(binPath string,args []string) (string, error) {

	outString,err := ExecScripts(binPath,args)
	if err != nil {
		log.Warnf("get uci %v ,error",args)
		return "",err
	}
	outString = strings.Replace(outString, "\n", "", -1) 
	// log.Debugf("uci show data: %v",outString)
	// outString = strings.Replace(outString, " ", "", -1) 
	// log.Debugf("uci show data: %v",outString)
	index := strings.Index(outString,"=")
	// output := strings.Split(outString[index+1:],"'")
	output := strings.Replace(outString[index+1:], "'", "", -1) 

	// log.Debugf("get data %v,len=%v",output,len(output))
	return output,nil

	// log.Debugf("output[0]=%v",output[0]);
	// return output[0],nil

}
func GetUciValueList(binPath string,args []string) ([]string, error) {

	outString,err := ExecScripts(binPath,args)
	if err != nil {
		log.Warnf("get uci %v ,error",args)
		return []string{},err
	}
	outString = strings.Replace(outString, "\n", "", -1) 
	index := strings.Index(outString,"=")
	outputstring := strings.Replace(outString[index+1:], "'", "", -1) 
	output := strings.Split(outputstring," ")
	log.Debugf("get data %v,len=%v",output,len(output))
	// if len(output) < 3 {
	// 	txt := fmt.Sprintf("get uci %v, result:%s, parse error",args,outString)
	// 	return "",errors.New(txt)
	// }
	return output,nil

}

func SetUciValue(value string) error{
	// txt := fmt.Sprintf("%s=%s",key,value)
	args := []string{"set",value}
	_,err := ExecScripts("uci",args)
	if err != nil {
		log.Errorf("set uci %v ,error",args)
		return err
	}
	return nil
}

func CommitUciValue(value string) error{
	// txt := fmt.Sprintf("%s=%s",key,value)
	args := []string{"commit",value}
	_,err := ExecScripts("uci",args)
	if err != nil {
		log.Errorf("set uci %v ,error",args)
		return err
	}
	return nil
}

func GetMacAddrs() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
			log.Errorf("fail to get net interfaces: %v", err)
			return ""
	}
	macAddr := ""
	for _, netInterface := range netInterfaces {
			if netInterface.Name == "eth0" {
				macAddr = netInterface.HardwareAddr.String()
				break
			}
	}
	return macAddr
}

func CheckIPValid(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		log.Errorf("wrong ipAddr format")
		return false
	}
	ip = ip.To4()
	if ip == nil {
		log.Errorf("wrong ipAddr to To4 format")
		return false
	}
	// return binary.BigEndian.Uint32(ip), nil
	return true
}

// CheckPrivateIPValid 检查 ip 是否是有效的，私有地址段
// ipaddr 带掩码 192.168.0.0/24
func CheckPrivateIPValid(ipaddr string) bool {
	ret := false
	log.Debugf("check ip valid: %v",ipaddr)
	ip, _, err := net.ParseCIDR(ipaddr)
	if err != nil {
		log.Errorf("check IP valid err:%v",err)
		return ret
	}
	ipv4Value := ip.To4()
	if ipv4Value == nil {
		return ret
	}
	ip0 := net.ParseIP("0.0.0.0")
	if ip.Equal(ip0) == true {
		return ret
	}
	return true
}

func convertQuardsToInt(splits []string) []int {
	quardsInt := []int{}

	for _, quard := range splits {
		j, err := strconv.Atoi(quard)
		if err != nil {
			panic(err)
		}
		quardsInt = append(quardsInt, j)
	}

	return quardsInt
}

func InetAtoN(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

func GetNumberIPAddresses(networkSize int) int {
	return 2 << uint(31-networkSize)
}
func GetCidrIpRange(cidr string) (string, string,int) {
	log.Debugf("check ip valid: %v",cidr)
	ipv4, ipv4Net, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Errorf("check IP valid err:%v",err)
		return "","",0
	}
	ipv4Value := ipv4.To4()
	if ipv4Value == nil {
		return "","",0
	}
	ipBegin := ipv4Value.String()

	networkSize,_ := ipv4Net.Mask.Size()

	// networkQuads := s.GetNetworkPortionQuards()
	networkQuads := convertQuardsToInt(strings.Split(ipBegin, "."))
	numberIPAddress := GetNumberIPAddresses(networkSize)
	networkRangeQuads := []string{}
	subnet_mask := 0xFFFFFFFF << uint(32-networkSize)
	networkRangeQuads = append(networkRangeQuads, fmt.Sprintf("%d", (networkQuads[0]&(subnet_mask>>24))+(((numberIPAddress-1)>>24)&0xFF)))
	networkRangeQuads = append(networkRangeQuads, fmt.Sprintf("%d", (networkQuads[1]&(subnet_mask>>16))+(((numberIPAddress-1)>>16)&0xFF)))
	networkRangeQuads = append(networkRangeQuads, fmt.Sprintf("%d", (networkQuads[2]&(subnet_mask>>8))+(((numberIPAddress-1)>>8)&0xFF)))
	networkRangeQuads = append(networkRangeQuads, fmt.Sprintf("%d", (networkQuads[3]&(subnet_mask>>0))+(((numberIPAddress-1)>>0)&0xFF)))

	return ipBegin,strings.Join(networkRangeQuads, "."),numberIPAddress

}

// GetCidrNetwork


func FindIfnameByAddresses(ipAddr string) (string,error) {
	result := ""
	ifaces, err := net.Interfaces()
	if err != nil {
			return result,err
	}

	for _, ifi := range ifaces {
		addrs, err := ifi.Addrs()
		if err != nil {
			log.Warnf("localAddresses: %v\n", err.Error())
			continue
		}

		ipInput := net.ParseIP(ipAddr)
		for _, a := range addrs {
			log.Debugf("%v -- %v\n", ifi.Name, a)
			ip, ipNet, err := net.ParseCIDR(a.String())
			if err != nil {
				// log.Errorf("check IP valid err:%v",err)
				continue
			}
			ipv4Value := ip.To4()
			if ipv4Value == nil {
				continue
			}
			if ipNet.Contains(ipInput) {
				result = ifi.Name
				log.Debugf("get ip:%v, ifname:%v",ipAddr,ifi.Name)
			}
		}
		// fmt.Printf("%v\n", ifi)
	}
	return result,err
}

func FindIfname(Ifname string) ([]string,error) {
	result := []string{}
	ifaces, err := net.Interfaces()
	if err != nil {
			return result,err
	}

	for _, ifi := range ifaces {
		addrs, err := ifi.Addrs()
		if err != nil {
			log.Warnf("localAddresses: %v\n", err.Error())
			continue
		}

		for _, a := range addrs {
			log.Debugf("%v -- %v\n", ifi.Name, a)
			ip, _, err := net.ParseCIDR(a.String())
			if err != nil {
				// log.Errorf("check IP valid err:%v",err)
				continue
			}
			ipv4Value := ip.To4()
			if ipv4Value == nil {
				continue
			}
			if Ifname == ifi.Name || ifi.Name == "lo" {
				continue
			}
			result = append(result,ifi.Name)
		}
		// fmt.Printf("%v\n", ifi)
	}
	return result,err
}

func GetMacAddressByIP(ip string,dhcpFile string) (string,error) {
	macAddr := ""
	// 获取 mac 地址, 先通过 dhcp client 获取到当前 mac 地址
	// 如果获取不到，则通过 uci show cascade 得到 mac 地址
	output,err := ReadFileAll(dhcpFile)
	if err != nil {
		return macAddr,err
	}
	log.Debugf("read dhcp file %v",output)
	outputLine := strings.Split(output,"\n")
	for _,value := range outputLine {
		if strings.Index(value,ip) >= 0 {
			// find
			log.Debugf("find mac addr in dhcp file :%v",value)
			data := strings.Split(value," ")
			macAddr = data[1]
			break;
		}
	}
	return macAddr,nil
}

func GetIPByMacAddress(macAddr string,dhcpFile string) (string,error) {
	ip := ""
	// 获取 mac 地址, 先通过 dhcp client 获取到当前 mac 地址
	// 如果获取不到，则通过 uci show cascade 得到 mac 地址
	output,err := ReadFileAll(dhcpFile)
	if err != nil {
		return ip,err
	}
	log.Debugf("read dhcp file %v",output)
	outputLine := strings.Split(output,"\n")
	for _,value := range outputLine {
		if strings.Index(value,macAddr) >= 0 {
			// find
			log.Debugf("find mac addr in dhcp file :%v",value)
			data := strings.Split(value," ")
			ip = data[2]
			break;
		}
	}
	return ip,nil
}