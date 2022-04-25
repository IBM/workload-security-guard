import React, {useEffect, useState} from "react";
import { Req } from './Req';
import { Resp } from './Resp';
import { Process } from './Process';
import { OnOff } from './OnOff';
import {TreeItem} from '@mui/lab';
import {Toggle} from './Guardian'

function Critiria(props) {
  const { data, onDataChange, nodeId, name} = props;
  if (!data.req) data.req = {}
  if (!data.resp) data.resp = {}
  if (!data.process) data.process = {}

  const [activeVal, setActive] = useState(data.active);
  
  useEffect(() => {
    Toggle([nodeId+">Req", nodeId+">Resp", nodeId+">Process"])
  }, [nodeId]);

  function handleReqChange(d) {
    console.log("handleReqChange", name, data, d)
    onDataChange(data)
  }
  
  function handleRespChange(d) {
    console.log("handleRespChange", name, data, d)
    onDataChange(data)
  }
  
  function handleProcessChange(d) {
    console.log("handleProcessChange", name, data, d)
    onDataChange(data)
  }

  function handleOnOffChange(key, d) {
    setActive(d)
    data.active = d
    onDataChange(data)
};

  console.log("Critiria data",name,  data)
 
  return (   
        <TreeItem nodeId={name} label={name} sx={{ textAlign: "start"}}>
          <OnOff data={activeVal} key={0} keyId={0} name={["Active"]} onDataChange={handleOnOffChange}/>
          <Req data={data.req} nodeId={nodeId+">Req"} name="Request" onDataChange={handleReqChange}></Req>   
          <Resp data={data.resp} nodeId={nodeId+">Resp"} name="Response" onDataChange={handleRespChange}></Resp>   
          <Process data={data.process} nodeId={nodeId+">Process"} name="Process" onDataChange={handleProcessChange}></Process>   
        </TreeItem>
  );
}

export {Critiria};
