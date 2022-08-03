import React from "react";
import {Box} from '@mui/material';
import { OnOff } from './OnOff';
import TreeItem from '@mui/lab/TreeItem';




function UnicodeBlock(props) {
  var { data, onDataChange, keyId, start } = props;
    
  //const [dataVal, setData] = useState(data);
  
  let value = data

  function handleOnOffChange(key, d) {
    if (d) {
      value |= 0x1 << key 
    } else {
      value &= ~(0x1 << key) 
    }
    
    console.log("UnicodeBlock", "key", key, "data", d, "value", value)
    onDataChange(keyId, value)
    //setData(value)
  };
  
  var current = start
  let buffer = []
  for (var i=0;i<32;i++) {
    buffer.push(<OnOff data={value & (0x1 << i)} key={i} keyId={i} name={(current).toString(16)} onDataChange={handleOnOffChange}></OnOff>)  
    current = current + 128
  }
  

return (
  <TreeItem nodeId={(start).toString(16)} label={(start).toString(16)}>
    <Box sx={{ display:"flex",  justifyContent: "start"}}>
               {buffer}
    </Box>
  </TreeItem> 
    )
}
export {UnicodeBlock}
//<Box sx={{ display:"flex",  justifyContent: "start"}}>
//</Box><Box sx={{ width: "8em",  fontSize: "0.8em"}}>{start}:</Box>
    