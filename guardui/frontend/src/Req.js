import React, {useEffect} from "react";
import TreeItem from '@mui/lab/TreeItem';

import { Set } from './Set';
import { Subnets } from './Subnets';
import { U8MinmaxSlice } from './U8MinmaxSlice';
import { Url } from './Url';
import { Qs } from './Qs';
import { Selection } from './Selection';

import { Headers } from './Headers';
import {Toggle} from './Guardian'


function Req(props) {
  var { data, onDataChange, nodeId, name } = props;
  if (!data.url) data.url = {}
  if (!data.qs) data.qs = {}
  if (!data.headers) data.headers = {}
  if (!data.contentlength) data.contentlength = []  
  if (!data.method) data.method = []
  if (!data.proto) data.proto = []
  if (!data.cip) data.cip = []

  useEffect(() => {
    Toggle([nodeId+">Url", nodeId+">Qs", nodeId+">Headers"])
  }, [nodeId]);
  
  function handleUrlChange(d) {
    console.log("handleUrlChange", data, d)  
    onDataChange(data)
  };
  function handleQsChange(d) {
    console.log("handleSegmentsChange", data, d)  
    onDataChange(data)
  };
  function handleHeadersChange(d) {
    console.log("handleHeadersChange", data, d)  
    onDataChange(data)
  };
  function handleContentLengthChange(d) {
    data.contentlength = d
    console.log("handleContentLengthChange", data)  
    onDataChange(data)
  };
  function handleMethodChange(d) {
    data.method = d
    console.log("handleMethodChange", data)  
    onDataChange(data)
  };
  function handleProtoChange(d) {
    data.proto = d
    console.log("handleProtoChange", data)  
    onDataChange(data)
  };
  function handleCipChange(d) {
    data.cip = d
    console.log("handleCipChange", data)  
    onDataChange(data)
  };
  var methods = ["GET", "POST", "PUT", "HEAD", "PATCH", "DELETE", "CONNECT", "OPTIONS", "TRACE" ]
  var protocols = ["HTTP/1.0", "HTTP/1.1",  "HTTP/2"]

  console.log("data.method", data.method)
return (
    <TreeItem nodeId={nodeId} label={name} sx={{ textAlign: "start"}}>
        <Selection data={data.method} nodeId={nodeId+">Method"} name="Method" selection={methods} onDataChange={handleMethodChange}></Selection>
        <Selection data={data.proto} nodeId={nodeId+">Proto"} name="Proto" selection={protocols} onDataChange={handleProtoChange}></Selection>
        <Subnets data={data.cip} nodeId={nodeId+">Cip"} name="Client Ip" onDataChange={handleCipChange}></Subnets>    
        <U8MinmaxSlice data={data.contentlength} nodeId={nodeId+">ContentLength"} name="Content Length" description="Range of the allowed content in powers of 2 (range is 2^min - 2^max)" onDataChange={handleContentLengthChange}></U8MinmaxSlice>    
        <Url data={data.url} nodeId={nodeId+">Url"} name="Url" onDataChange={handleUrlChange}></Url>
        <Qs data={data.qs} nodeId={nodeId+">Qs"} name="Qs" onDataChange={handleQsChange}></Qs>
        <Headers data={data.headers} nodeId={nodeId+">Headers"} name="Headers" onDataChange={handleHeadersChange}></Headers>
    </TreeItem>
    )
}
export {Req}
//<Set data={data.proto} nodeId={nodeId+">Proto"} name="Proto" onDataChange={handleProtoChange}></Set>    
          