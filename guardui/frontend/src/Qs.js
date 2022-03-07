import React from "react";
import TreeItem from '@mui/lab/TreeItem';
import { KeyVal } from './KeyVal';


function Qs(props) {
  var { data, onDataChange,nodeId,  name } = props;
  if (!data.kv) data.kv = {}

  //let value = data
  function handleKvChange(d) {
    //if (d) {
      //data.kv = d
      //console.log("handleKvChange", data)  
      onDataChange(data)
      //setData(value);
    //}
  };

return (
    <TreeItem nodeId={nodeId} label={name} sx={{ textAlign: "start"}}> 
        <KeyVal data={data.kv} nodeId={nodeId+">KeyVal"} name="KeyVal"  onDataChange={handleKvChange}></KeyVal>
    </TreeItem>
    )
}
export {Qs}