import TreeItem from '@mui/lab/TreeItem';

import { SimpleVal } from "./SimpleVal";
import { Structured } from "./Structured";


function RBody(props) {
  var { data, onDataChange, nodeId, name } = props;
  if (!data.unstructured) data.unstructured = {}
  if (!data.structured) data.structured = {}
  
  function handleUnstructuredChange(d) {
    console.log("handleUnstructuredChange", data, d)  
    onDataChange(data)
  };
  function handleStructuredChange(d) {
    console.log("handleStructuredChange", data, d)  
    onDataChange(data)
  };
  
return (
    <TreeItem nodeId={nodeId} label={name} sx={{ textAlign: "start"}}>
        <SimpleVal data={data.unstructured} nodeId={nodeId+">Unstructured"} keyId="unstructured" name="Unstructured" onDataChange={handleUnstructuredChange}></SimpleVal>
        <Structured data={data.structured} nodeId={nodeId+">Structured"} name="Structured" onDataChange={handleStructuredChange}></Structured>
    </TreeItem>
    )
}
export {RBody}
//<Set data={data.proto} nodeId={nodeId+">Proto"} name="Proto" onDataChange={handleProtoChange}></Set>    
          