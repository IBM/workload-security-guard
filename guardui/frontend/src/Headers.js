import React, {useEffect} from "react";
import TreeItem from '@mui/lab/TreeItem';
import { KeyVal } from './KeyVal';
import {Toggle} from './Guardian'
import {FormHelperText} from '@mui/material';



function Headers(props) {
  var { data, onDataChange, nodeId, name } = props;
  
  if (!data.kv) data.kv = {}
  useEffect(() => {
    Toggle([nodeId+">KeyVal"])
  }, [nodeId]);

  function handleKvChange(d) {
   onDataChange(data)
  };

return (
    <TreeItem nodeId={nodeId} label={name}> 
        <FormHelperText>Specification of the Headers</FormHelperText>
        <KeyVal data={data.kv} nodeId={nodeId+">KeyVal"} name="KeyVal" onDataChange={handleKvChange}></KeyVal>
    </TreeItem>
    )
}
export {Headers}