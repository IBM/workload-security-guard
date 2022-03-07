import React from "react";
import Box from '@mui/material/Box';
import Slider from '@mui/material/Slider';



function U8MinMax(props) {
  var { data, onDataChange, keyId } = props;
  
  //var value = data
  function handleChange(event, d) {
    console.log("U8MinMax",  d)
    onDataChange(keyId, d)
  };
  
  return (
    <Box sx={{ width: 600,  marginLeft: "10px" }}>
      <Slider
        value={data}
        onChange={handleChange}
        valueLabelDisplay="auto"
        max={255}
        size="small"
      />
    </Box>
  );
}
export {U8MinMax}