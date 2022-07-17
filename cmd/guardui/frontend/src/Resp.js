import React, {useEffect} from "react";
import TreeItem from '@mui/lab/TreeItem';
import { Headers } from './Headers';
import {Toggle} from './Guardian'


function Resp(props) {
  var { data, onDataChange, nodeId, name } = props;
  if (!data.url) data.url = {}
  if (!data.qs) data.qs = {}
  if (!data.headers) data.headers = {}
  
  useEffect(() => {
    Toggle([nodeId+">Headers"])
  }, [nodeId]);

  function handleHeadersChange(d) {
    console.log("handleHeadersChange", data, d)  
    onDataChange(data)
  };

return (
    <TreeItem nodeId={nodeId} label={name} sx={{ textAlign: "start"}}>
        <Headers data={data.headers} nodeId={nodeId+">Headers"} name="Headers" onDataChange={handleHeadersChange}></Headers>
    </TreeItem>
    )
}
export {Resp}
  