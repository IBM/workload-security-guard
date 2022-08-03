import React, {useState} from "react";
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import Select from '@mui/material/Select';
import Button from "@mui/material/Button";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogTitle from "@mui/material/DialogTitle";



function SelectKeyDialog(props) {

  var { data, onClose, open, name} = props;
  const [keyname, setKeyname] = useState("");
  //const [opened, setOpen] = useState(open);
  var value = keyname
  
  function handleSelectKey(event) {
    let k = event.target.value
    setKeyname(k)
    console.log("handleSelectKey",k)
  }
  function handleCancel() {
    console.log("handleCancel")
    //setOpen(false)
    onClose("")
  }
  function handleOk() {
    //setOpen(false)
    onClose(keyname)
  }

  let i=0
  let listres = data.map((key) => {
    i = i+1 
    return (
    <MenuItem value={key} key={key}>{key}</MenuItem>
    )
  });
  console.log("listres", listres)
  return (
    <Dialog open={open} onClose={handleCancel} >
          <DialogTitle>{name}</DialogTitle>
          <DialogContent >
              <InputLabel id="delKey-label">Select key</InputLabel>
              <Select
                labelId="delKey-label"
                id="demo-simple-select-standard"
                fullWidth
                value={value}
                onChange={handleSelectKey}
                label="Key"
              >
                <MenuItem value="" key="">
                  
                </MenuItem>
                {listres}
              </Select>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleCancel}>Cancel</Button>
            <Button onClick={handleOk}>Ok</Button>
          </DialogActions>
    </Dialog>
  );
}

export {SelectKeyDialog};