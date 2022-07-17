import React, {useEffect, useState} from "react";
import { Req } from './Req';
import { Resp } from './Resp';
import { RBody } from './RBody';
import { Process } from './Process';
import { OnOff } from './OnOff';
import {TreeItem} from '@mui/lab';
import {Toggle} from './Guardian'

function Critiria(props) {
  const { data, onDataChange, nodeId, name} = props;
  //console.log("data.active in Critiria", data.active)  
  
  if (!data.active) data.active = false
  if (!data.req) data.req = {}
  if (!data.resp) data.resp = {}
  if (!data.reqbody) data.reqbody = {}
  if (!data.respbody) data.respbody = {}
  if (!data.process) data.process = {}

  const [active, setActive] = useState(data.active);
  console.log("data.active in Critiria", data.active)  
  //const [activeVal, setActive] = useState(data.active);
  //console.log("val in Critiria", val)  
  console.log("val.active in Critiria", active)  
  
  useEffect(() => {
    Toggle([nodeId+">Req", nodeId+">Resp", nodeId+">Process"])
  }, [nodeId]);

  useEffect(() => {
    setActive(data.active)
    console.log("val.active in Critiria", data.active)  
    
  }, [data.active]);

  function handleReqChange(d) {
    console.log("handleReqChange", name, data, d)
    onDataChange(data)
  }
  
  function handleRespChange(d) {
    console.log("handleRespChange", name, data, d)
    onDataChange(data)
  }
  
  function handleReqBodyChange(d) {
    console.log("handleReqBodyChange", name, data, d)
    onDataChange(data)
  }
  
  function handleRespBodyChange(d) {
    console.log("handleRespBodyChange", name, data, d)
    onDataChange(data)
  }
  
  function handleProcessChange(d) {
    console.log("handleProcessChange", name, data, d)
    onDataChange(data)
  }

  function handleOnOffChange(key, d) {
    //setActive(d)
    data.active = d
    console.log("handleOnOffChange active", d, "data", data)
    onDataChange(data)
    //setActive(d)
    setActive(data.active)
};

  console.log("Critiria data",name,  data)
 
  return (   
        <TreeItem nodeId={name} label={name} sx={{ textAlign: "start"}}>
          <OnOff data={data.active} key={0} keyId={0} name={["Active"]} onDataChange={handleOnOffChange}/>
          <Req data={data.req} nodeId={nodeId+">Req"} name="Request" onDataChange={handleReqChange}></Req>   
          <Resp data={data.resp} nodeId={nodeId+">Resp"} name="Response" onDataChange={handleRespChange}></Resp>   
          <RBody data={data.reqbody} nodeId={nodeId+">ReqBody"} name="ReqBody" onDataChange={handleReqBodyChange}></RBody>   
          <RBody data={data.respbody} nodeId={nodeId+">RespBody"} name="RespBody" onDataChange={handleRespBodyChange}></RBody>   
          <Process data={data.process} nodeId={nodeId+">Process"} name="Process" onDataChange={handleProcessChange}></Process>   
        </TreeItem>
  );
}

export {Critiria};
