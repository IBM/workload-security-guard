import React, {useState} from "react";
import TreeItem from '@mui/lab/TreeItem';
import { SimpleVal } from "./SimpleVal";
import {SelectKeyDialog} from "./SelectKeyDialog";
import {AddKeyDialog} from './AddKeyDialog';
import {CheckedKeySlice} from "./CheckedKeySlice";
import {IconButton, Divider} from '@mui/material';
import AddRoundedIcon from '@mui/icons-material/AddRounded';
import RemoveRoundedIcon from '@mui/icons-material/RemoveRounded';

function KeyVal(props) {
  var { data, onDataChange, nodeId, name } = props;
  console.log("keyval data start", data)
  if (!data["vals"]) data["vals"] = {}
  if (!data["otherVals"]) data["otherVals"] = {}
  if (!data["otherKeynames"]) data["otherKeynames"] = {}
  if (!data["minimalSet"]) data["minimalSet"] = []

  //let value = data
  console.log("keyval data", data)
  let minSet = {}
  for (var key in data.vals) {  
    minSet[key] = data.minimalSet.includes(key)
  }
  const [version, setVersion] = useState(0);
  const [newkeyOpen, setNewkey] = useState(false);
  const [delkeyOpen, delKey] = useState(false);
  const [minimalSet, setMinimalSet] = useState(minSet);
  

 

  function onValsChange(key, d) {
    //data.vals[key] = d 
    //console.log("onValsChange", key, d, data)
    onDataChange(data)
  }
  function onOtherValsChange(key, d) {
    //data.otherVals = d 
    //console.log("onOtherValsChange", d, data)
    onDataChange(data)  
  }
  function onOtherKeynamesChange(key, d) {
    //data.otherKeynames = d 
    //console.log("onOtherKeynamesChange", d, data)
    onDataChange(data)
  }
  function onMinimalSetChange(key, d) {
    data.minimalSet = []
    for (var k in d) { 
      if (d[k]) data.minimalSet.push(k)
    }
    console.log("onMinimalSetChange", d, data)
    onDataChange(data)
  }
  function handleValDel() {
    delKey(true)
  }
  function handleValAdd() {
    setNewkey(true)
  }
  function onSelectKey(k) {
    console.log("handleDelKey",k)
    if (k !== "") {
      delete(data.vals[k])
      console.log("handleDelKey", data)  
      onDataChange(data)
      //setVersion(version+1);
    }
    delKey(false)
  }
  function onNewKey(k) {
    if (k !== "") {
      data.vals[k] = {}
      console.log("onNewKey", data)  
    }
    setNewkey(false)
    let minSet = {}
    for (key in data.vals) {  
      minSet[key] = data.minimalSet.includes(key)
    }
    console.log("onNewKey  minSet", minSet)
    onDataChange(data)
    setMinimalSet(minSet)
    setVersion(version+1)
    
    //setVersion(version+1);
  };
  
  let res = []
  for (key in data.vals) {  
    let v =  data.vals[key]
    res.push(
      <SimpleVal data={v} nodeId={nodeId+key} key={nodeId+key} keyId={key} name={key} onDataChange={onValsChange}/>
    )
  }
  console.log("******* KEy Val Refresh ******", minimalSet)
  return (
      <TreeItem nodeId={nodeId} label={name}>
        <SelectKeyDialog open={delkeyOpen} name="Key to delete" data ={Object.keys(data.vals)} onClose={onSelectKey} ></SelectKeyDialog>
        <AddKeyDialog open={newkeyOpen} data ={Object.keys(data.vals)} onClose={onNewKey} ></AddKeyDialog>
        {res}
        <Divider>
           <IconButton color="primary" aria-label="Add Key" onClick={handleValAdd}>
                <AddRoundedIcon />
            </IconButton>
            <IconButton color="error" aria-label="Del Key" onClick={handleValDel}>
                <RemoveRoundedIcon />
            </IconButton>
        </Divider>
        <CheckedKeySlice data={minimalSet} nodeId={nodeId+">minimalSet"} keyId="minimalSet" name="Minimal Set" onDataChange={onMinimalSetChange} />
        <SimpleVal data={data.otherVals} nodeId={nodeId+">OtherVals"} keyId="otherVals" name="Other Vals" onDataChange={onOtherValsChange}></SimpleVal>
        <SimpleVal data={data.otherKeynames} nodeId={nodeId+">OtherKeynames"} keyId="otherKeynames" name="Other Keynames" onDataChange={onOtherKeynamesChange}></SimpleVal>
      </TreeItem>
 
        );
}

export {KeyVal};
// </div><Box sx={{ display: "flex", justifyContent: "start", margin: "0.2em"}}>
// <Button sx={{ width: "12em", alignItems: "center", fontSize: "0.8em"}} onClick={handleChange} variant="contained">{name}</Button>
// </Box><Box sx={{ display: showVal ? "flex" : "none", flexDirection:  "column", justifyContent: "start"}}>
