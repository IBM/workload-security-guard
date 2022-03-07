import React from "react";
import TreeItem from '@mui/lab/TreeItem';

import { Url } from './Url';
import { Qs } from './Qs';
import { Headers } from './Headers';


function Req(props) {
  var { data, onDataChange, nodeId, name } = props;
  if (!data.url) data.url = {}
  if (!data.qs) data.qs = {}
  if (!data.headers) data.headers = {}

  let value = data
  function handleUrlChange(d) {
    //if (d) {
      //value.val = d
      //console.log("handleUrlChange", value)  
      onDataChange(value)
      //setData(value);
    //}
  };
  function handleQsChange(d) {
    //if (d) {
      //value.qs = d
      //console.log("handleSegmentsChange", value)  
      onDataChange(value)
      //setData(value);
    //}
  };
  function handleHeadersChange(d) {
    //if (d) {
      //value.headers = d
      //console.log("handleSegmentsChange", value)  
      onDataChange(value)
      //setData(value);
    //}
  };

return (
    <TreeItem nodeId={nodeId} label={name} sx={{ textAlign: "start"}}>
        <Url data={data.url} nodeId="ReqUrl" name="Url" onDataChange={handleUrlChange}></Url>
        <Qs data={data.qs} nodeId="ReqQs" name="Qs" onDataChange={handleQsChange}></Qs>
        <Headers data={data.headers} nodeId="ReqHeaders" name="Headers" onDataChange={handleHeadersChange}></Headers>
    </TreeItem>
    )
}
export {Req}
  