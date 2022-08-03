import React from "react";
import {Box, Button} from '@mui/material';



function OnOff(props) {
  var { data, keyId, onDataChange, name } = props;
  //const [value, setValue] = React.useState(data);
  //var value = data

  function handleChange(event) {
    console.log("handleChange data", keyId, !data)
    onDataChange(keyId, !data)
    //console.log("handleChange", keyId, !value)
    //onDataChange(keyId, !value)
    //value = !value;
    //color = value ? "success":"error"
  };
  var color = data ? "success":"error"
  var variant = data ? "contained":"outlined"
  //var color = value ? "success":"error"
  //var variant = value ? "contained":"outlined"
  console.log("onoff data", data)
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
