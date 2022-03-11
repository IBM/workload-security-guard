import React, {useEffect} from "react";
import TreeItem from '@mui/lab/TreeItem';

import { Url } from './Url';
import { Qs } from './Qs';
import { Headers } from './Headers';
import {Toggle} from './Guardian'


function Req(props) {
  var { data, onDataChange, nodeId, name } = props;
  if (!data.url) data.url = {}
  if (!data.qs) data.qs = {}
  if (!data.headers) data.headers = {}

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

return (
    <TreeItem nodeId={nodeId} label={name} sx={{ textAlign: "start"}}>
        <Url data={data.url} nodeId={nodeId+">Url"} name="Url" onDataChange={handleUrlChange}></Url>
        <Qs data={data.qs} nodeId={nodeId+">Qs"} name="Qs" onDataChange={handleQsChange}></Qs>
        <Headers data={data.headers} nodeId={nodeId+">Headers"} name="Headers" onDataChange={handleHeadersChange}></Headers>
    </TreeItem>
    )
}
export {Req}
  