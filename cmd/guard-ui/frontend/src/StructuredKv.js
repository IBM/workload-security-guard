import React, {useState} from "react";
import TreeItem from '@mui/lab/TreeItem';

import { Structured } from "./Structured";

import {SelectKeyDialog} from "./SelectKeyDialog";
import {AddKeyDialog} from './AddKeyDialog';
import {IconButton, Divider, FormHelperText} from '@mui/material';
import AddRoundedIcon from '@mui/icons-material/AddRounded';
import RemoveRoundedIcon from '@mui/icons-material/RemoveRounded';
import {Box} from '@mui/material';


function StructuredKv(props) {
  var { data, onDataChange, nodeId, name } = props;
  console.log("keyval data start", data)
  
  console.log("keyval data", data)
  const [version, setVersion] = useState(0);
  const [newkeyOpen, setNewkey] = useState(false);
  const [delkeyOpen, delKey] = useState(false);

  function onValsChange(key, d) {
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
      delete(data[k])
      console.log("handleDelKey", data)  
      onDataChange(data)
    }
    delKey(false)
  }
  function onNewKey(k) {
    if (k !== "") {
      data[k] = {}
      console.log("onNewKey", data)  
    }
    setNewkey(false)
    onDataChange(data)
    setVersion(version+1)
  };
  
  let res = []
  for (let key in data) {  
    res.push(
        <Structured data={data[key]} nodeId={nodeId+key} key={nodeId+key} keyId={key} name={key} onDataChange={onValsChange}/>
    )
  }
  return (
      <TreeItem nodeId={nodeId} label={name}>
        <Box sx={{borderLeft: "1px solid orange", paddingLeft: "1em"}}>

        <SelectKeyDialog open={delkeyOpen} name="Key to delete" data ={Object.keys(data)} onClose={onSelectKey} ></SelectKeyDialog>
        <AddKeyDialog open={newkeyOpen} data ={Object.keys(data)} onClose={onNewKey} ></AddKeyDialog>
        <FormHelperText>Named Keys (per key specificcation of the allowed value)</FormHelperText>
        {res}
        <Divider>
           <IconButton color="primary" aria-label="Add Key" onClick={handleValAdd}>
                <AddRoundedIcon />
            </IconButton>
            <IconButton color="error" aria-label="Del Key" onClick={handleValDel}>
                <RemoveRoundedIcon />
            </IconButton>
        </Divider>
        </Box>
      </TreeItem>
        );
}

export {StructuredKv};
