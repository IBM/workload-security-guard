import React, {useState} from "react";
import {Box} from '@mui/material';
import { OnOff } from './OnOff';
import TreeItem from '@mui/lab/TreeItem';



function Selection(props) {
  var { data, onDataChange, nodeId, name, selection } = props;
  
  const [dataVal, setData] = useState(data);
 
  let value = dataVal

  function handleOnOffChange(key, d) {
      if (d) {
        value |= 0x1 << key 
      } else {
        value &= ~(0x1 << key) 
      }
      
      console.log("key", key, "data", d, "value", value)
      let results = {}
      for (var i=0;i<selection.length;i++) {
            if (value & (0x1 << i)) {
                results[selection[i]] = true
            } 
      }
      onDataChange(results)
      setData(value)
  };
  
  let buffer = []
  for (var i=0;i<selection.length;i++) {
    buffer.push(<OnOff data={value & (0x1 << i)} key={i} keyId={i} name={selection[i]} onDataChange={handleOnOffChange}/>)  
  }


return (
    <TreeItem nodeId={nodeId} label={name}>
     <Box sx={{ display: "flex", justifyContent: "start", alignItems: "center"}}>
           {buffer}
        </Box>
    </TreeItem> 
    )
}
export {Selection}
//<Box sx={{ display:"flex", alignItems: "center", justifyContent: "start", margin: "0.2em"}}>
//<Button sx={{ width: "12em", alignItems: "center", fontSize: "0.8em"}} onClick={handleChange} variant="contained">{name}</Button>
               