//import logo from './logo.svg';
import './App.css';
import {TextField, Box} from '@mui/material';
import React, { useState } from "react";

function Identity(props) {
    const { ns, sid, handleChange } = props;

  const [nsDataVal, setNs] = useState(ns);
  const [sidDataVal, setSid] = useState(sid);
    
  handleChange(nsDataVal, sidDataVal)
 
  function handleNsChange(event) {
    console.log("handleNsChange", event.target.value)
    setNs(event.target.value)
  }
  function handleSidChange(event) {
    console.log("handleSidChange", event.target.value)
    setSid(event.target.value)
  }
  return (
    <Box sx={{ display: "flex", justifyContent: "start", margin: "0.2em"}}>
        <TextField id="ns" label="Namespace" variant="standard" value={nsDataVal} onChange={handleNsChange} />
        <TextField id="sid" label="Service Name" variant="standard"  value={sidDataVal} onChange={handleSidChange} />
    </Box>
  );
}
export {Identity}