package go_common

import (
	"fmt"
	cm "./common"
	"testing"
)


func TestGetIPRange(t *testing.T) {
	// testIP := "192.168.0.0/16"
	testIP := "1.0.8.0/21"
	
	startIP,endIP,number := cm.GetCidrIpRange(testIP)
	fmt.Printf("get %s, startIP:%s, endIP:%s, number=%d\n",testIP,startIP,endIP,number)
	start_id := cm.InetAtoN(startIP)
	end_id := cm.InetAtoN(endIP)
	fmt.Printf("get start_id:%d, end_id:%d, number=%d\n",start_id,end_id,number)

}