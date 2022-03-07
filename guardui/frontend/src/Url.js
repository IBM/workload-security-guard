import React from "react";
import TreeItem from '@mui/lab/TreeItem';
import { U8MinmaxSlice } from './U8MinmaxSlice';
import { SimpleVal } from './SimpleVal';


function Url(props) {
  var { data, onDataChange, nodeId, name } = props;
  if (!data.val) data.val = {}
  if (!data.segments) data.segments = []
  //let value = data
  function handleValChange(key, d) {
    //if (d) {
      //data.val = d
      //console.log("handleValChange", data)  
      onDataChange(data)
    //}
  };
  function handleSegmentsChange(d) {
    //if (d) {
      data.segments = d
      console.log("handleSegmentsChange", data)  
      onDataChange(data)
    //}
  };

return (
    <TreeItem nodeId={nodeId} label={name}> 
        <U8MinmaxSlice data={data.segments} nodeId={nodeId+">Segements"} name="Url Segments" onDataChange={handleSegmentsChange}></U8MinmaxSlice>
        <SimpleVal data={data.val} keyId={"url"} nodeId={nodeId+">Val"} name="Url Characters" onDataChange={handleValChange}></SimpleVal>
    </TreeItem>
    )
}
export {Url}
  