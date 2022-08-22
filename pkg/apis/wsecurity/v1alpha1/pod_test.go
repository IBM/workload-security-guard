package v1alpha1

import (
	"os"
	"testing"
)

func updateAll(data []byte) {
	procNet = "/tmp/proc/net/"
	os.WriteFile(procNet+"tcp", data, 0644)
	os.WriteFile(procNet+"udp", data, 0644)
	os.WriteFile(procNet+"udplite", data, 0644)
	os.WriteFile(procNet+"tcp6", data, 0644)
	os.WriteFile(procNet+"udp6", data, 0644)
	os.WriteFile(procNet+"udplite6", data, 0644)
}

func TestPod_V1(t *testing.T) {

	os.MkdirAll(procNet, os.ModePerm)
	data1 := []byte(`  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
0: 02001102:E12E 0217D902:0050 06 00000000:00000000 03:00001599 00000000     0        0 0 3 0000000000000000
1: 02001102:C6F0 02B9FA02:0050 06 00000000:00000000 03:00001569 00000000     0        0 0 3 0000000000000000`)

	data2 := []byte(`  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
0: 02001104:E12E 0417D902:0050 06 00000000:00000000 03:00001599 00000000     0        0 0 3 0000000000000000
1: 02001104:C6F0 04B9FA02:0050 06 00000000:00000000 03:00001569 00000000     0        0 0 3 0000000000000000`)

	data3 := []byte(`  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
0: 02001103:E12E 0317D902:0050 06 00000000:00000000 03:00001599 00000000     0        0 0 3 0000000000000000
1: 02001103:C6F0 03B9FA02:0050 06 00000000:00000000 03:00001569 00000000     0        0 0 3 0000000000000000`)

	t.Run("CIDR management", func(t *testing.T) {
		var profile1 PodProfile
		var profile1a PodProfile
		var profile2 PodProfile
		var profile2a PodProfile
		var profile3 PodProfile
		var pile1 PodPile
		var pile2 PodPile
		var config1 PodConfig
		var config1a PodConfig
		var config2 PodConfig
		var config2a PodConfig

		updateAll(data1)
		profile1.Profile()
		profile1a.Profile()
		updateAll(data2)
		profile2.Profile()
		profile2a.Profile()
		updateAll(data3)
		profile3.Profile()

		pile1.Add(&profile1)
		config1.Learn(&pile1)
		if ret := config1.Decide(&profile1); ret != "" {
			t.Errorf("Decide return error %s when expected to ok", ret)
		}
		if ret := config1.Decide(&profile2); ret == "" {
			t.Errorf("Decide return ok when expected an error")
		}
		if ret := config1.Decide(&profile3); ret == "" {
			t.Errorf("Decide return ok when expected an error")
		}
		pile2.Add(&profile2)
		pile2.Merge(&pile1)
		config2.Learn(&pile2)
		if ret := config2.Decide(&profile1); ret != "" {
			t.Errorf("Decide return error %s when expected to ok", ret)
		}
		if ret := config2.Decide(&profile2); ret != "" {
			t.Errorf("Decide return error %s when expected to ok", ret)
		}
		if ret := config2.Decide(&profile3); ret != "" {
			t.Errorf("Decide return error %s when expected to ok", ret)
		}

		pile1.Clear()
		pile2.Clear()
		pile1.Add(&profile1a)
		pile2.Add(&profile2a)
		config1a.Learn(&pile1)
		config2a.Learn(&pile2)
		config1a.Fuse(&config2a)
		if ret := config1a.Decide(&profile1); ret != "" {
			t.Errorf("Decide return error %s when expected to ok", ret)
		}
		if ret := config1a.Decide(&profile2); ret != "" {
			t.Errorf("Decide return error %s when expected to ok", ret)
		}

		if ret := config1a.Decide(&profile3); ret != "" {
			t.Errorf("Decide return error %s when expected to ok", ret)
		}

	})
}
