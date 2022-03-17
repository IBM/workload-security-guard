import React, {useState} from "react";
import { Control } from './Control';
import { Critiria } from './Critiria';
import {TreeView} from '@mui/lab';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';

let fullExpanded = {"Configured":true, "Control":true, "Learned":true}
let Toggle = (nodeIds) => {
  let n = {} 
  for (let k in fullExpanded) {
    n[k] = true
  }
  for (let i=0;i<nodeIds.length;i++) {
    n[nodeIds[i]] = true
  }
  fullExpanded = n
};

function Guardian(props) {
  const { data, setCollapse } = props;
  if (!data.configured) data.configured = {}
  if (!data.learned) data.learned = {}
  if (!data.control) data.control = {}
  const [expanded, setExpanded] = useState([]);
  
  //useEffect(() => {
  //  console.log("Guardian Exanding", fullExpanded)
  //  setExpanded(Object.keys(fullExpanded));
  //}, [fullExpanded]);

  const handleToggle = (event, nodeIds) => {
    setExpanded(nodeIds);
  };

  function collapse() {
    setExpanded([])
  }
  setCollapse(collapse)

  function handleConfiguredChange(d) {
    console.log("handleConfiguredChange", data, d)
  }

  function handleLearnedChange(d) {
    console.log("handleLearnedChange", data, d)
  }

  function handleControldChange(d) {
    console.log("handleConfiguredChange", data, d)
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
        <Control data={data.control}  nodeId="Control" name="Control" onDataChange={handleControldChange}></Control>    
        <Critiria data={data.configured}  nodeId="Configured" name="Configured" onDataChange={handleConfiguredChange}></Critiria>    
        <Critiria data={data.learned}  nodeId="Learned" name="Learned" onDataChange={handleLearnedChange}></Critiria>  
    </TreeView>
  );
}

export {Guardian, Toggle};
//  <div><pre>{JSON.stringify(data, null, 2) }</pre></div>
//<Box sx={{ display:"flex",  justifyContent: "start"}}>