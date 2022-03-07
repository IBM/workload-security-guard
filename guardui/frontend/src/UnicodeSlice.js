import React, {useState} from "react";
import {Box} from '@mui/material';
import { UnicodeBlock } from './UnicodeBlock';
import TreeItem from '@mui/lab/TreeItem';
import {IconButton, Divider} from '@mui/material';
import AddRoundedIcon from '@mui/icons-material/AddRounded';
import RemoveRoundedIcon from '@mui/icons-material/RemoveRounded';
import LanguageRoundedIcon from '@mui/icons-material/LanguageRounded';
import {UnicodeTable} from './UnicodeTable';
import {SelectKeyDialog} from "./SelectKeyDialog";


function mergeData(data) {
  for (var i=data.length-1;i>=0;i--) {
      if (data[i]!==0) break
      data.pop()
  }
  return data    
}

function UnicodeSlice(props) {
  var { data, onDataChange, nodeId, name } = props;
  const [languageOpen, setLanguage] = useState(false);
  if (!data) data = []
  console.log("UnicodeSlice before merge ", data) 
  data = mergeData(data)
  console.log("UnicodeSlice after merge ", data) 
  const [dataVal, setData] = useState(data);
  console.log("UnicodeSlice dataVal ", dataVal) 
  
  let value = [...dataVal]
  
  console.log("UnicodeSlice", value) 
  function handleUnicodeBlockChange(key, d) {
    value[key] = d
    console.log("handleUnicodeBlockChange", value)  
    onDataChange(value)
    setData(value);
  };
  function handleUnicodeBlockAdd() {
    value.push(0)
    console.log("handleUnicodeBlockAdd", value)  
    onDataChange(value)
    setData([...value]);
  };
  function handleUnicodeBlockDel() {
    value.pop()
    console.log("handleUnicodeBlockDel", value)  
    onDataChange(value)
    setData([...value]);
  }

  function handleAddLanguage() {
    setLanguage(true)
  }
  
  function onSelectKey(k) {
    console.log("onSelectKey",k)
    if (k !== "") {
      let languages = UnicodeTable[k]
      console.log("onSelectKey languages",languages)
      for (var i=0; i<languages.length;i++) {
        let l = languages[i] 
        if (l===0) continue
        l -= 1
        let l_block = Math.floor(l/0x20)
        let l_index = l%0x20
        let length = value.length
        console.log("onSelectKey", i, l, l_block, l_index, length, value) 
        for (var j=length; j< l_block+1; j++) {
          value.push(0)
        }
        console.log("onSelectKey", i, value.length, value) 
        value[l_block] |= 0x1 << l_index 
      }
      console.log("onSelectKey languages", languages, value)  
      onDataChange(value)
      setData([...value]);
      //setVersion(version+1);
    }
    setLanguage(false)
  }
  var current = 0x80

  console.log("UnicodeSlice name", name, "value", value)
  let buffer = []
  if (value) {
    for (var i=0;i<value.length;i++) {
        buffer.push(<UnicodeBlock key={i} keyId={i} data={value[i]} start={current} onDataChange={handleUnicodeBlockChange}></UnicodeBlock>)  
        current += 0x1000
      }
  }

 
  

return (
    <TreeItem nodeId={nodeId} label={name}>
          <Box sx={{ display:  "flex" , flexDirection:  "column", justifyContent: "start"}}>
           {buffer}
          </Box>
          <Divider>
          <SelectKeyDialog open={languageOpen} name="Language to add" data ={Object.keys(UnicodeTable)} onClose={onSelectKey} ></SelectKeyDialog>
 
           <IconButton color="primary" aria-label="add UnicodeBlock" onClick={handleUnicodeBlockAdd}>
                <AddRoundedIcon />
            </IconButton>
            <IconButton color="primary" aria-label="add Language" onClick={handleAddLanguage}>
                <LanguageRoundedIcon />
            </IconButton>
            <IconButton color="error" aria-label="del UnicodeBlock" onClick={handleUnicodeBlockDel}>
                <RemoveRoundedIcon />
            </IconButton>
            
            </Divider>
    </TreeItem>
    )
}
export {UnicodeSlice}
  