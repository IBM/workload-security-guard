import React, {useState, useEffect} from "react";
import TreeItem from '@mui/lab/TreeItem';
import { SimpleVal } from "./SimpleVal";
import {SelectKeyDialog} from "./SelectKeyDialog";
import {AddKeyDialog} from './AddKeyDialog';
import {IconButton, Divider} from '@mui/material';
import AddRoundedIcon from '@mui/icons-material/AddRounded';
import RemoveRoundedIcon from '@mui/icons-material/RemoveRounded';
import {Toggle} from './Guardian'

function KeyVal(props) {
  var { data, onDataChange, nodeId, name } = props;
  console.log("keyval data start", data)
  if (!data["vals"]) data["vals"] = {}
  if (!data["otherVals"]) data["otherVals"] = {}
  if (!data["otherKeynames"]) data["otherKeynames"] = {}
  //if (!data["minimalSet"]) data["minimalSet"] = []

  useEffect(() => {
    //Toggle([nodeId+">minimalSet", nodeId+">OtherVals", nodeId+">OtherKeynames"])
    Toggle([nodeId+">OtherVals", nodeId+">OtherKeynames"])
  }, [nodeId]);
  
  console.log("keyval data", data)
  const [version, setVersion] = useState(0);
  const [newkeyOpen, setNewkey] = useState(false);
  const [delkeyOpen, delKey] = useState(false);

  function onValsChange(key, d) {
    onDataChange(data)
  }
  function onOtherValsChange(key, d) {
    onDataChange(data)  
  }
  function onOtherKeynamesChange(key, d) {
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
    }
    delKey(false)
  }
  function onNewKey(k) {
    if (k !== "") {
      data.vals[k] = {}
      console.log("onNewKey", data)  
    }
    setNewkey(false)
    onDataChange(data)
    setVersion(version+1)
  };
  
  let res = []
  for (let key in data.vals) {  
    let v =  data.vals[key]
    res.push(
      <SimpleVal data={v} nodeId={nodeId+key} key={nodeId+key} keyId={key} name={key} onDataChange={onValsChange}/>
    )
  }
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
        <SimpleVal data={data.otherVals} nodeId={nodeId+">OtherVals"} keyId="otherVals" name="Other Vals" onDataChange={onOtherValsChange}></SimpleVal>
        <SimpleVal data={data.otherKeynames} nodeId={nodeId+">OtherKeynames"} keyId="otherKeynames" name="Other Keynames" onDataChange={onOtherKeynamesChange}></SimpleVal>
      </TreeItem>
 
        );
}

export {KeyVal};
// </div><Box sx={{ display: "flex", justifyContent: "start", margin: "0.2em"}}>
// <Button sx={{ width: "12em", alignItems: "center", fontSize: "0.8em"}} onClick={handleChange} variant="contained">{name}</Button>
// </Box><Box sx={{ display: showVal ? "flex" : "none", flexDirection:  "column", justifyContent: "start"}}>
//<CheckedKeySlice data={minimalSet} nodeId={nodeId+">minimalSet"} keyId="minimalSet" name="Minimal Set" onDataChange={onMinimalSetChange} />
        