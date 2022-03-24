import React, {useState} from "react";
import Button from "@mui/material/Button";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogTitle from "@mui/material/DialogTitle";
import TextField from "@mui/material/TextField";
import {Alert, AlertTitle} from '@mui/material';


function AddKeyDialog(props) {
  var { data, onClose, open, regex, title} = props;
  const [keyname, setKeyname] = useState("");
  if (!title) {
    title = "Add Key"
  }
  
  function handleNewKey(event) {
    let k = event.target.value
    setKeyname(k)
    console.log("handleNewKey",k)
  }

  function handleCancel() {
    console.log("handleCancel")  
    onClose("")
  }
  function handleOk() {
    console.log("handleOk", keyname)  
    if (data.includes(keyname))  {
        console.log("handleOk key already exists") 
        onClose("")
        return
    }

    if (regex && !regex.test(keyname)) {
      console.log("handleOk ilegal key") 
      //onClose("")
      return
    }

    console.log("handleOk after closing")  
    onClose(keyname)
  }

  let alert
  if (regex && keyname.length) { 
    if (!regex.test(keyname)) { 
      alert = <Alert severity="warning">Ilegal Value</Alert>
    } else {
      alert =  <Alert severity="success">Legal Value</Alert>
    }  

  }

  return (       
    <Dialog open={open} onClose={handleCancel} >
    <DialogTitle>{title}</DialogTitle>
    <DialogContent >
      <TextField
        autoFocus
        margin="dense"
        id="name"
        label="Key Name"
        type="keyname"
        variant="standard"
        onChange={handleNewKey}
      />
    {alert}      
    </DialogContent>
    <DialogActions>
      <Button onClick={handleCancel}>Cancel</Button>
      <Button onClick={handleOk}>Ok</Button>
    </DialogActions>
  </Dialog>
  );
}
export {AddKeyDialog};