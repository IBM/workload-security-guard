import React, {useEffect} from "react";
import TreeItem from '@mui/lab/TreeItem';
import { U8MinmaxSlice } from './U8MinmaxSlice';
import { Set } from './Set';
import {Toggle} from './Guardian'


function Process(props) {
  var { data, onDataChange, nodeId, name } = props;
  if (!data.responsetime) data.responsetime = []
  if (!data.completiontime) data.completiontime = []
  if (!data.tcp4peers) data.tcp4peers = {}
	if (!data.udp4peers) data.udp4peers = {}
	if (!data.udplite4peers) data.udplite4peers = {}
	if (!data.tcp6peers) data.tcp6peers = {}
	if (!data.udp6peers) data.udp6peers = {}
	if (!data.udplite6peers) data.udplite6peers = {}

  useEffect(() => {
    Toggle([nodeId+">responsetime", nodeId+">completiontime"])
  }, [nodeId]);

  function handleResponseTimeChange(key, d) {
    data.responsetime = d
    console.log("handleResponseTimeChange", data)  
    onDataChange(data)
  };

  function handleCompletionTimeChange(d) {
    data.completiontime = d
    console.log("handleCompletionTimeChange", data)  
    onDataChange(data)
  };

  function handleTcp4peersChange(d) {
    data.tcp4peers = d
    console.log("handleTcp4peersChange", data)  
    onDataChange(data)
  }

  function handleUdp4peersChange(d) {
    data.udp4peers = d
    console.log("handleUdp4peersChange", data)  
    onDataChange(data)
  }

  function handleUdpLite4peersChange(d) {
    data.udplite4peers = d
    console.log("handleUdpLite4peersChange", data)  
    onDataChange(data)
  }
  
  function handleTcp6peersChange(d) {
    data.tcp6peers = d
    console.log("handleTcp6peersChange", data)  
    onDataChange(data)
  }

  function handleUdp6peersChange(d) {
    data.udp6peers = d
    console.log("handleUdp6peersChange", data)  
    onDataChange(data)
  }
  
  function handleUdpLite6peersChange(d) {
    data.udplite6peers = d
    console.log("handleUdpLite6peersChange", data)  
    onDataChange(data)
  }

return (
    <TreeItem nodeId={nodeId} label={name}> 
        <U8MinmaxSlice data={data.responsetime} nodeId={nodeId+">responsetime"} name="Response Time" onDataChange={handleResponseTimeChange}></U8MinmaxSlice>
        <U8MinmaxSlice data={data.completiontime} nodeId={nodeId+">completiontime"} name="Completion Time" onDataChange={handleCompletionTimeChange}></U8MinmaxSlice>
        <Set data={data.tcp4peers} nodeId={nodeId+">tcp4peers"} name="TCP Ipv4 peers" onDataChange={handleTcp4peersChange}></Set>
        <Set data={data.udp4peers} nodeId={nodeId+">udp4peers"} name="UDP Ipv4 peers" onDataChange={handleUdp4peersChange}></Set>
        <Set data={data.udplite4peers} nodeId={nodeId+">udplite4peers"} name="UDP Lite Ipv4 peers" onDataChange={handleUdpLite4peersChange}></Set>
        <Set data={data.tcp6peers} nodeId={nodeId+">tcp6peers"} name="TCP Ipv6 peers" onDataChange={handleTcp6peersChange}></Set>
        <Set data={data.udp6peers} nodeId={nodeId+">udp6peers"} name="UDP Ipv6 peers" onDataChange={handleUdp6peersChange}></Set>
        <Set data={data.udplite6peers} nodeId={nodeId+">udplite6peers"} name="UDP Lite Ipv6 peers" onDataChange={handleUdpLite6peersChange}></Set>
    </TreeItem>
    )
}
export {Process}
  