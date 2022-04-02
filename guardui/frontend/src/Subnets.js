import React, {useState} from 'react';
import TreeItem from '@mui/lab/TreeItem';
import {IconButton, Divider} from '@mui/material';
import AddRoundedIcon from '@mui/icons-material/AddRounded';
import RemoveRoundedIcon from '@mui/icons-material/RemoveRounded';
import {SelectKeyDialog} from "./SelectKeyDialog";
import {AddKeyDialog} from './AddKeyDialog';
import { styled } from '@mui/material/styles';
import Chip from '@mui/material/Chip';
import Paper from '@mui/material/Paper';
import {FormHelperText} from '@mui/material';

const ListItem = styled('li')(({ theme }) => ({
  margin: theme.spacing(0.5),
}));


function Subnets(props) {
  var { data, onDataChange, nodeId, name, ipv6 } = props; 
  var value = {}
  
  //const [version, setVersion] = useState(0);
  const [dataVal, setData] = useState(data);
  const [newkeyOpen, setNewkey] = useState(false);
  const [delkeyOpen, delKey] = useState(false);
  console.log("Set dataVal", dataVal)
    
    function handleValDel() {
        delKey(true)
    }
    function handleValAdd() {
        setNewkey(true)
    }
    function onSelectKey(k) {
        console.log("handleDelKey",k)
        if (k !== "") {
          delete(value[k])
          console.log("handleDelKey", value)  
          onDataChange(value)
          setData(value)
          //setVersion(version+1);
        }
        delKey(false)
    }
    function onNewKey(k) {
        if (k !== "") {
            value[k] = true
            console.log("onNewKey", value)  
            setData(value)
        }
        setNewkey(false)
        onDataChange(value)
        //setVersion(version+1)
    };
    function handleDelete(k) {
        return () => {
            delete(value[k])
            console.log("handleDelete", k, value)  
            setData(value)
            onDataChange(value)
        }     
    };
   
    let keys = []
    for (let key in dataVal) {  
        value[key] = true
        keys.push(key)
    }
    console.log("Set value", value, "keys", keys)
    let ipregex = /^((25[0-5]|(2[0-4]|1[0-9]|[1-9]|)[0-9])(\.(?!$)|$)){4}$/;
    let cidr = /^((25[0-5]|(2[0-4]|1[0-9]|[1-9]|)[0-9])\.){3}(25[0-5]|(2[0-4]|1[0-9]|[1-9]|)[0-9])\/(3[0-2]|([12]|)[0-9])$/;
    if (ipv6) {
        cidr = new RegExp('^(' +
                                  '([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}' +
                                '|([0-9a-fA-F]{1,4}:){1,7}:' +
                                '|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}' +
                                '|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}' +
                                '|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}' +
                                '|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}' +
                                '|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}' +
                                '|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})' +
                                '|:((:[0-9a-fA-F]{1,4}){1,7}|:)' +
                            ')/(12[0-8]|((1[01]|[1-9]|)[0-9]))$') 
       
    }
    var subnets 
    if (keys.length>0) { 
        subnets = keys.map((key) => {
            let icon;
            return (
            <ListItem key={key}>
                <Chip
                icon={icon}
                label={key}
                onDelete={handleDelete(key)}
                />
            </ListItem>
            );
        })
    } else {
        subnets = "No Subnets"
    }
    return (
       <TreeItem nodeId={nodeId} label={name}>
            <SelectKeyDialog open={delkeyOpen} name="Key to delete" data ={Object.keys(value)} onClose={onSelectKey} ></SelectKeyDialog>
            <AddKeyDialog open={newkeyOpen} data ={Object.keys(value)} regex={cidr} title={"Add Subnet"} onClose={onNewKey} ></AddKeyDialog>
            <FormHelperText>Allowed subnets:</FormHelperText>
            <Paper sx={{
                    display: 'flex',
                    justifyContent: 'center',
                    flexWrap: 'wrap',
                    listStyle: 'none',
                    p: 0.5,
                    m: 0,
                }}
                component="ul"
            >
                {subnets}
                </Paper>
            <Divider>
            <IconButton color="primary" aria-label="Add Key" onClick={handleValAdd}>
                <AddRoundedIcon />
            </IconButton>
            <IconButton color="error" aria-label="Del Key" onClick={handleValDel}>
                <RemoveRoundedIcon />
            </IconButton>
        </Divider>
        
       </TreeItem>
    )

}
export {Subnets};                   