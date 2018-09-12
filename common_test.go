package go_common

import (
	"fmt"
	cm "./common"
	"testing"
)


func TestGetIPRange(t *testing.T) {
	// testIP := "192.168.0.0/16"
	testIP := "1.0.1.0/22"
	
	startIP,endIP := cm.GetCidrIpRange(testIP)
	fmt.Printf("get %s, startIP:%s, endIP:%s\n",testIP,startIP,endIP)

}