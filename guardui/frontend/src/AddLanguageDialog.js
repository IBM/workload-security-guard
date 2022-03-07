import React, {useState} from "react";
import Button from "@mui/material/Button";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogTitle from "@mui/material/DialogTitle";
import TextField from "@mui/material/TextField";


function AddLanguageDialog(props) {
  var { data, onClose, open} = props;
  const [keyname, setKeyname] = useState("");

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
        console.log("handleOk after key already exists")
        return
    }
    console.log("handleOk after closing")  
    onClose(keyname)
  }

  return (       
    <Dialog open={open} onClose={handleCancel} >
    <DialogTitle>Add Key</DialogTitle>
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
    </DialogContent>
    <DialogActions>
      <Button onClick={handleCancel}>Cancel</Button>
      <Button onClick={handleOk}>Ok</Button>
    </DialogActions>
  </Dialog>
  );
}
export {AddLanguageDialog};