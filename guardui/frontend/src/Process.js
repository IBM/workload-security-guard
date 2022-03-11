import React, {useEffect} from "react";
import TreeItem from '@mui/lab/TreeItem';
import { U8MinmaxSlice } from './U8MinmaxSlice';
import {Toggle} from './Guardian'


function Process(props) {
  var { data, onDataChange, nodeId, name } = props;
  if (!data.responsetime) data.responsetime = []
  if (!data.completiontime) data.completiontime = []
  
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

return (
    <TreeItem nodeId={nodeId} label={name}> 
        <U8MinmaxSlice data={data.responsetime} nodeId={nodeId+">responsetime"} name="Response Time" onDataChange={handleResponseTimeChange}></U8MinmaxSlice>
        <U8MinmaxSlice data={data.completiontime} nodeId={nodeId+">completiontime"} name="Completion Time" onDataChange={handleCompletionTimeChange}></U8MinmaxSlice>
    </TreeItem>
    )
}
export {Process}
  