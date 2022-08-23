import './App.css';
import {TextField, Box} from '@mui/material';
import React from "react";
import Radio from '@mui/material/Radio';
import RadioGroup from '@mui/material/RadioGroup';
import FormControlLabel from '@mui/material/FormControlLabel';

function Identity(props) {
  const { setSid, sidVal, setNs, nsVal, setCm, cmVal } = props;

    
  function handleNsChange(event) {
    console.log("handleNsChange", event.target.value)
    setNs(event.target.value)
  }
  function handleSidChange(event) {
    console.log("handleSidChange", event.target.value)
    setSid(event.target.value)
  }
  function handleCmChange(event) {
    console.log("handleCmChange", event.target.value)
    setCm(event.target.value)
  }
  return (
    <Box sx={{ display: "flex", width: "100%", height: "auto", justifyContent: "center", margin: "0.2em", marginLeft: "1em"}}>
        <RadioGroup row value={cmVal}  onChange={handleCmChange} name="crd-cm-group">
          <FormControlLabel value="crd" control={<Radio />} label="CRD" />
          <FormControlLabel value="cm" control={<Radio />} label="ConfigMap" />
        </RadioGroup>
        <TextField sx={{ margin: "1em"}} id="ns" label="Namespace" variant="standard" value={nsVal} onChange={handleNsChange} />
        <TextField sx={{ margin: "1em"}} id="sid" label="Service Name" variant="standard"  value={sidVal} onChange={handleSidChange} />
    </Box>
  );
}
export {Identity}