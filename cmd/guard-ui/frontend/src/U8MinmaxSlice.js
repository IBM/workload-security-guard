import React, {useState} from 'react';
import { U8MinMax } from './U8MinMax';
import TreeItem from '@mui/lab/TreeItem';
import {IconButton, Divider, FormHelperText} from '@mui/material';
import AddRoundedIcon from '@mui/icons-material/AddRounded';
import RemoveRoundedIcon from '@mui/icons-material/RemoveRounded';


function mergeData(data) {
    let newData = []
    newData.push(data[0])
    let j=0
    for (var i=1;i<data.length;i++) {
        if (data[i].min-1<=newData[j].max) {
            if (data[i].max>newData[j].max) {   
                newData[j].max = data[i].max
            } 
            continue
        }
        newData.push(data[i])
        j = j + 1
    }
    return newData    
}


function U8MinmaxSlice(props) {
    var { data, onDataChange, nodeId, name, description } = props;
    if (data.length===0)   data = [{min:0, max:0}] 
    data = mergeData(data)
    const [dataVal, setData] = useState(data);

    let value = [...dataVal]
    
    function handleU8MinmaxChange(key, d) {
      if (d) {
        value[key].min = d[0]
        value[key].max = d[1]
        value = mergeData(value)
        console.log("handleU8MinmaxChange", value)  
        onDataChange(value)
        setData(value);
      }
    };
    function handleU8MinmaxAdd() {
        let lastkey = value.length-1
        let newkey = value.length
        if (newkey>7) return
        if ((lastkey>=0) && (value[lastkey].max>253)) return
        if (lastkey>=0) {
            value.push({min:255, max:255})    
        }
        else {
            value.push({min:0, max:0})
        }
       
        console.log("handleU8MinmaxAdd", value)  
        onDataChange(value)
        setData(value);
    };
    function handleU8MinmaxDel() {
        if (value.length<=0) return
        value.pop()
        console.log("handleU8MinmaxDel", value)  
        onDataChange(value)
        setData(value);
    }
    var i = -1    
    var res = value.map((range) => {
            i = i+1
            return (
                <U8MinMax data={[range.min, range.max]} key={i} keyId={i} onDataChange={handleU8MinmaxChange}></U8MinMax>  
            )
        }
    
    );
    return (
       <TreeItem nodeId={nodeId} label={name}>
           <FormHelperText>{description} </FormHelperText>
           {res}
           <Divider>
           <IconButton color="primary" aria-label="add U8MinMax" onClick={handleU8MinmaxAdd}>
                <AddRoundedIcon />
            </IconButton>
            <IconButton color="error" aria-label="del U8MinMax" onClick={handleU8MinmaxDel}>
                <RemoveRoundedIcon />
            </IconButton>
            </Divider>
       </TreeItem>
    )

}
export {U8MinmaxSlice};
//<Box sx={{ display:"flex", alignItems: "center", justifyContent: "start", margin: "0.2em"}}>
//<Button sx={{ width: "12em", alignItems: "center", fontSize: "0.8em"}} onClick={handleChange} variant="contained">{name}</Button>
//<Box sx={{ display: "flex", justifyContent: "start", alignItems: "center" }}>
                                