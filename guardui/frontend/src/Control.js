import React, {useState} from "react";
import {TreeItem} from '@mui/lab';
import {FormControlLabel, Select, MenuItem, FormGroup, Checkbox} from '@mui/material';

import {Box} from '@mui/material';
const autotips = {"disable": "Alert based on manual configuration", "alert": "Alert basd on automation", "learn": "Learning used for recomendations only"}

function Control(props) {
  const { data, nodeId, name } = props;
  if (!data.auto) data.auto = false
  if (!data.learn) data.learn = false
  if (!data.consult) data.consult = false
  if (!data.rpm) data.rpm = 60
  if (!data.block) data.block = false

  const [blockVal, setBlock] = useState(data.block);
  const [learnVal, setLearn] = useState(data.learn);
  const [autoVal, setAuto] = useState(data.auto);
  
  
  function handleBlock(event) {
      console.log("handleBlock event.target.checked",event.target.checked, !data.block )
      data.block = !data.block
      setBlock(data.block);
  }

  function handleAutoSelect(event) {
    console.log("handleAutoSelect",event.target.value )
    switch (event.target.value) {
        case "learn":
            data.learn = true
            data.auto = false
            data.consult = false
            setLearn(true)
            setAuto(false)
            break
        case "alert":
            data.learn = true
            data.auto = true
            data.consult = false
            setLearn(true)
            setAuto(true)
            break
        default:
            data.learn = false
            data.auto = false
            data.consult = false
            setLearn(false)
            setAuto(false)
           
        
    }
  }

  console.log("Control data", data)
  var menuVal
  if (!learnVal) menuVal = "disable"
  else menuVal = autoVal ? "alert" : "learn"
  return (
    
        <TreeItem nodeId={nodeId} label={name} sx={{ textAlign: "start"}}>
          <FormGroup>
            <Box sx={{ display:"flex",  justifyContent: "start", alignItems: "center"}}>
            
              <Select
                labelId="selelct automation"
                id="automation-select"
                value={menuVal}
                label="Automation"
                onChange={handleAutoSelect}
                sx={{ margin: "2em"}}
              >
                <MenuItem value={"disable"} title={autotips["disable"]}>Automation disabled</MenuItem>
                <MenuItem value={"learn"} title={autotips["learn"]}>Learn only</MenuItem>
                <MenuItem value={"alert"} title={autotips["alert"]}>Auto alert</MenuItem>
              </Select>
              {autotips[menuVal]} 
              <FormControlLabel sx={{ margin: "2em"}} control={<Checkbox checked={blockVal} onChange={handleBlock}/>} label="Block on Alert"/>
            </Box>
          </FormGroup>
        </TreeItem>
  );
}

export {Control};
