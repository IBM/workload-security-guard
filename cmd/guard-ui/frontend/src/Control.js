import React, {useEffect, useState} from "react";
import {TreeItem} from '@mui/lab';
import {Stack, MenuItem, FormGroup} from '@mui/material';
import TextField from '@mui/material/TextField';
import MuiAlert from '@mui/material/Alert';

import {Box} from '@mui/material';
const alertMenuTips = {
                  "disable": "Never alert!", 
                  "manual": "Alert based on manual configuration", 
                  "auto": "Alert basd on automation", 
}

const learnMenuTips = {
                  "disable": "Never learn", 
                  "enable": "Learn from occurrences where no alert issued", 
                  "force": "Learn from occurrences that where not blocked (even during alert)", 
}

const blockMenuTips = {
                  "disable": "Never block", 
                  "enable": "Block on any alert", 
}

function Control(props) {
  const { data, nodeId, name } = props;
  if (!data.auto) data.auto = false
  if (!data.learn) data.learn = false
  if (!data.force) data.force = false
  if (!data.alert) data.alert = false
  if (!data.consult) data.consult = false
  if (!data.rpm) data.rpm = 60
  if (!data.block) data.block = false

  
  const [blockVal, setBlock] = useState(data.block);
  const [learnVal, setLearn] = useState(data.learn);
  const [forceVal, setForce] = useState(data.force);
  const [autoVal, setAuto] = useState(data.auto);
  const [alertVal, setAlert] = useState(data.alert);
  
  useEffect(() => {
    setAlert(data.alert)
    setLearn(data.learn)
    setForce(data.force)
    setBlock(data.block)
    setAuto(data.auto)
    console.log("euse effect of all", data)  
    
  }, [data]);

  
  /*
  function handleBlock(event) {
      console.log("handleBlock event.target.checked",event.target.checked, !data.block )
      data.block = !data.block
      setBlock(data.block);
  }
  function handleLearn(event) {
    console.log("handleLearn event.target.checked",event.target.checked, !data.learn )
    data.learn = !data.learn
    setBlock(data.learn);
  }
  */

  function handleAlertSelect(event) {
    console.log("handleAutoSelect",event.target.value )
    switch (event.target.value) {
        case "manual":
            data.auto = false
            data.alert = true
            break
        case "auto":
            data.auto = true
            data.alert = true
            break
        default:  // disable
            data.auto = false
            data.alert = false 
    }
    setAlert(data.alert)
    setAuto(data.auto)
  }


  function handleLearnSelect(event) {
    console.log("handleLearnSelect",event.target.value )
    switch (event.target.value) {
        case "enable":
            data.learn = true
            data.force = false
            break
        case "force":
            data.learn = true
            data.force = true
            break
        default:  // disable
            data.learn = false
            data.force = false
    }
    setLearn(data.learn)
    setForce(data.force)
  }
  function handleBlockSelect(event) {
    console.log("handleBlockSelect",event.target.value )
    switch (event.target.value) {
        case "enable":
            data.block = true
            break
        default:  // disable
            data.block = false
    }
    setBlock(data.block)
  }

  console.log("Control data", data)
  var alertMenuVal, learnMenuVal, blockMenuVal
  if (!alertVal) alertMenuVal = "disable"
  else alertMenuVal = autoVal ? "auto" : "manual"
  learnMenuVal = learnVal? (forceVal? "force" : "enable") : "disable"
  blockMenuVal = blockVal? "enable" : "disable"

  return (
    
        <TreeItem nodeId={nodeId} label={name} sx={{ textAlign: "start"}}>
          <FormGroup sx={{ margin: 2 }}>
            <Box sx={{ display:"flex",  justifyContent: "start", alignItems: "center"}}>
            
              <TextField
                id="alerting-select"
                select
                value={alertMenuVal}
                label="Select"
                onChange={handleAlertSelect}
                sx={{ margin: "2em"}}
                helperText="Alert control"
              >
                <MenuItem value={"disable"} title={alertMenuTips["disable"]}>No alerts</MenuItem>
                <MenuItem value={"manual"} title={alertMenuTips["manual"]}>Manual alert</MenuItem>
                <MenuItem value={"auto"} title={alertMenuTips["auto"]}>Auto alert</MenuItem>
              </TextField>
              <TextField
                id="learning-select"
                select
                value={learnMenuVal}
                label="Select"
                onChange={handleLearnSelect}
                sx={{ margin: "2em"}}
                helperText="Learning control"
              >
                <MenuItem value={"disable"} title={learnMenuTips["disable"]}>Disable</MenuItem>
                <MenuItem value={"enable"} title={learnMenuTips["enable"]}>Enable</MenuItem>
                <MenuItem value={"force"} title={learnMenuTips["force"]}>Force</MenuItem>
              </TextField>
              <TextField
                id="blocking-select"
                select
                value={blockMenuVal}
                label="Select"
                onChange={handleBlockSelect}
                sx={{ margin: "2em"}}
                helperText="Block control"
              >
                <MenuItem value={"disable"} title={blockMenuTips["disable"]}>Disable</MenuItem>
                <MenuItem value={"enable"} title={blockMenuTips["enable"]}>Enable</MenuItem>
              </TextField>
              
            </Box>
            <Stack spacing={2} sx={{ width: '100%' }}>
            <MuiAlert severity="info" elevation={6} >{alertMenuTips[alertMenuVal]}</MuiAlert>
            <MuiAlert severity="info" elevation={6}  >{learnMenuTips[learnMenuVal]} </MuiAlert>
            <MuiAlert severity="info" elevation={6}  >{blockMenuTips[blockMenuVal]} </MuiAlert>
            </Stack>
            
            
          </FormGroup>
        </TreeItem>
  );
}

export {Control};
// <FormControlLabel sx={{ margin: "2em"}} control={<Checkbox checked={blockVal} onChange={handleBlock}/>} label="Block on Alert" helperText="Alert control"/>
// <FormControlLabel sx={{ margin: "2em"}} control={<Checkbox checked={learnVal} onChange={handleLearn}/>} label="Enable Learning"/>
// {alertMenuTips[menuVal]} 
             
//<Select
//                labelId="alerting-select-label"
//                id="alerting-select"
//                value={menuVal}
//                label="Alert"
//                onChange={handleAutoSelect}
//                sx={{ margin: "2em"}}
//                helperText="Alert control"
//              ></Select>