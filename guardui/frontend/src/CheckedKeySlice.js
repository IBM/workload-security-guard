import React from "react";
import {Box} from '@mui/material';
import { OnOff } from './OnOff';
import TreeItem from '@mui/lab/TreeItem';


function CheckedKeySlice(props) { 
  var { data, onDataChange, nodeId, keyId, name } = props;
  const [version, setVersion] = React.useState(0); 
 
  let value = data

  function handleOnOffChange(key, d) {
      value[key] = d 
      console.log("key", key, "data", d, "value", value)
      onDataChange(keyId,   value)
      //setData(value)
      setVersion(version+1)
  };
  
  let buffer = []
  for (var key in value) {  
    buffer.push(<OnOff data={value[key]} key={key} keyId={key} name={key} onDataChange={handleOnOffChange}/>)  
  }

return (
    <TreeItem nodeId={nodeId} label={name}>
     <Box sx={{ display: "flex", justifyContent: "start", alignItems: "center"}}>
           {buffer}
        </Box>
    </TreeItem> 
    )
}
export {CheckedKeySlice}
//<Box sx={{ display:"flex", alignItems: "center", justifyContent: "start", margin: "0.2em"}}>
//<Button sx={{ width: "12em", alignItems: "center", fontSize: "0.8em"}} onClick={handleChange} variant="contained">{name}</Button>
               