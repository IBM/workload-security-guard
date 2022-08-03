import TreeItem from '@mui/lab/TreeItem';
import React, {useState} from 'react';
import { SimpleVal } from "./SimpleVal";
import { StructuredKv } from "./StructuredKv";

import Radio from '@mui/material/Radio';
import RadioGroup from '@mui/material/RadioGroup';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormControl from '@mui/material/FormControl';

function Structured(props) {
  var { data, onDataChange, nodeId, name } = props;
  console.log("Structured data is",data)
  if (!data.val) data.val = {}
  if (!data.kv) data.kv = {}
  console.log("Structured data is",data)
  
  if (!data.kind) data.kind = "skip"
  
  const [kind, setKind] = useState(data.kind);
  const [showVal, setShowVal] = useState(data.kind === "array" || data.kind === "string" || data.kind === "number" || data.kind === "boolean");
  const [showKeyVal, setShowKeyVal] = useState(data.kind === "object");

  function handleKindSelect(event) {
    console.log("v",event.target.value )
   
    data.kind = event.target.value
    setKind(event.target.value)
    setShowVal(data.kind === "array" || data.kind === "string" || data.kind === "number" || data.kind === "boolean");
    setShowKeyVal(data.kind === "object");
    onDataChange(data)
  }

  function handleKvChange(d) {
    console.log("handleKvChange", data, d)  
    onDataChange(data)
  };

  function onValChange(keyId, d) {
    console.log("onValChange", data, keyId, d)  
    onDataChange(data)  
  };
  
return (
    <TreeItem nodeId={nodeId} label={name} sx={{ textAlign: "start"}}>
      <FormControl sx={{ margin: "0px"}}>
        <RadioGroup
        row
        name="kind-group"
        value={kind}
        onChange={handleKindSelect}
      >
            <FormControlLabel value="object" control={<Radio />} label="Object" />
            <FormControlLabel value="array" control={<Radio />} label="List" />
            <FormControlLabel value="string" control={<Radio />} label="String" />
            <FormControlLabel value="number" control={<Radio />} label="Number" />
            <FormControlLabel value="boolean" control={<Radio />} label="Boolean" />
            <FormControlLabel value="skip" control={<Radio />} label="Ignore" />
        </RadioGroup>
        {showVal?<SimpleVal data={data.val} nodeId={nodeId+">Val"} keyId="val" name="Val" onDataChange={onValChange}></SimpleVal> : null }
        {showKeyVal?<StructuredKv data={data.kv} nodeId={nodeId+">KeyVal"} name="KeyVal" onDataChange={handleKvChange}></StructuredKv>: null }
        
      </FormControl>
    </TreeItem>
    )
}
export {Structured}
//<Set data={data.proto} nodeId={nodeId+">Proto"} name="Proto" onDataChange={handleProtoChange}></Set>    
          