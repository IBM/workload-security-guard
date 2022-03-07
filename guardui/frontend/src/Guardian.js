import React, {useState} from "react";
import { Req } from './Req';
import {TreeView, TreeItem} from '@mui/lab';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import {FormControlLabel, Select, MenuItem, FormGroup, Checkbox} from '@mui/material';

import {Box} from '@mui/material';
const autotips = {"disable": "Alert based on manual configuration", "alert": "Alert basd on automation", "learn": "Learning used for recomendations only"}

function Guardian(props) {
  const { data, setCollapse } = props;
  if (!data.auto) data.auto = "disable"
  if (!data.forceAllow) data.forceAllow = false
  if (!data.req) data.req = {}
  
  const [expanded, setExpanded] = useState([]);
  const [allowVal, setAllow] = useState(!data.forceAllow);
  const [autoVal, setAuto] = useState(data.auto);

  const handleToggle = (event, nodeIds) => {
    setExpanded(nodeIds);
  };

  function collapse() {
    setExpanded([])
  }
  setCollapse(collapse)

  function handleReqChange(d) {
    console.log("handleReqChange", d)
  }
  
  function handleAllow(event) {
      console.log("handleAllow event.target.checked",event.target.checked )
      data.forceAllow = !allowVal
      setAllow(!allowVal);
  }

  function handleAutoSelect(event) {
    console.log("handleAutoSelect",event.target.value )
    data.auto = event.target.value
    setAuto(event.target.value);
  }

  console.log("Guardian data", data)
  if (Object.keys(data).length === 0) {
    return (
      <div>
        Loading...
      </div>
    )
  }
      
  return (
    <TreeView
      aria-label="file system navigator"
      defaultCollapseIcon={<ExpandMoreIcon />}
      defaultExpandIcon={<ChevronRightIcon />}
      expanded={expanded}
      onNodeToggle={handleToggle}
      sx={{ height: 480, 
        flexGrow: 1, 
        maxWidth: 800, 
        overflowY: 'auto', 
        border: "1px solid", 
        borderRadius: "1em",
        margin: "1em",
        justifyContent: "start", 
        alignItems: "center"
      }}
    > 
    
        <TreeItem nodeId="Constrols" label="Controls" sx={{ textAlign: "start"}}>
          <FormGroup>
            <Box sx={{ display:"flex",  justifyContent: "start", alignItems: "center"}}>
            
              <Select
                labelId="selelct automation"
                id="automation-select"
                value={autoVal}
                label="Automation"
                onChange={handleAutoSelect}
                sx={{ margin: "2em"}}
              >
                <MenuItem value={"disable"} title={autotips["disable"]}>Automation disabled</MenuItem>
                <MenuItem value={"learn"} title={autotips["learn"]}>Learn only</MenuItem>
                <MenuItem value={"alert"} title={autotips["alert"]}>Auto alert</MenuItem>
              </Select>
              {autotips[autoVal]} 
              <FormControlLabel sx={{ margin: "2em"}} control={<Checkbox checked={!allowVal} onChange={handleAllow}/>} label="Block on Alert"/>
            </Box>
          </FormGroup>
        </TreeItem>
        <Req data={data.req} nodeId="Req" name="Request" onDataChange={handleReqChange}></Req>    
         
    </TreeView>
  );
}

export {Guardian};
//  <div><pre>{JSON.stringify(data, null, 2) }</pre></div>
//<Box sx={{ display:"flex",  justifyContent: "start"}}>