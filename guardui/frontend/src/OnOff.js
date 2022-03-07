import React from "react";
import {Box, Button} from '@mui/material';



function OnOff(props) {
  var { data, keyId, onDataChange, name } = props;
    
  //const [value, setValue] = React.useState(data);
  var value = data

  function handleChange(event) {
    console.log("handleChange", keyId, !value)
    onDataChange(keyId, !value)
    value = !value;
    color = value ? "success":"error"
  };
  var color = value ? "success":"error"
  var variant = value ? "contained":"outlined"
  console.log("color", color)
  return (
    <Box >
      <Button
        onClick={handleChange}
        color={color}
        variant={variant}
        
        style={{maxWidth: '9.2em', minWidth: '1.2em', height: '1.4em', textTransform: "none", padding:5, margin:1}}
      >{name}</Button>
    </Box>
  );
}
export {OnOff}
// 