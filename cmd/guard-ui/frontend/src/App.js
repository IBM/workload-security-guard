//import logo from './logo.svg';
import './App.css';
import {Guardian} from './Guardian';
import {Button, Collapse, Alert, AlertTitle, Stack} from '@mui/material';
import React, { useState, useRef } from "react";
import GetIcon from '@mui/icons-material/GetApp';
import SendIcon from '@mui/icons-material/Send';
import ClearIcon from '@mui/icons-material/Clear';
import LearnIcon from '@mui/icons-material/LightbulbOutlined';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell, { tableCellClasses } from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import {TextField} from '@mui/material';
import Radio from '@mui/material/Radio';
import RadioGroup from '@mui/material/RadioGroup';
import FormControlLabel from '@mui/material/FormControlLabel';
import { Typography } from '@mui/material';

import { styled } from '@mui/material/styles';

const yaml = require('js-yaml');

const StyledTableCell = styled(TableCell)(({ theme }) => ({
  [`&.${tableCellClasses.head}`]: {
    //backgroundColor: theme.palette.grey[300],
    
    fontSize: 18,
    padding: "2px",
    borderTop: "1px solid black"
   // borderLeft: "1px solid grey"
  },
  [`&.${tableCellClasses.body}`]: {
    fontSize: 14,
    padding: 0,
    borderSpacing: 0,
    margin: 0,
    border: "hidden"
  },
}));

const StyledButton = styled(Button)({
  //color: 'darkslategray',
  //backgroundColor: 'aliceblue',
  //padding: 8,
  borderRadius: 10,
  borderColor: "#d2e6fa",
  //color: "Blue",
  //borderColor: "red",
  width: "9em",
  padding: 1,
  margin: "0em"
});


const testData = {}




function App() {
  const [dataVal, setData] = useState(testData);
  const [successVal, setSuccess] = useState("");
  const [errorVal, setError] = useState("");
  const [infoVal, setInfo] = useState("");
  const [sidVal, setSid] = useState("");
  const [nsVal, setNs] = useState("");
  const [cmVal, setCm] = useState("crd");
  
  let collapse

  function onGetClick() {
    setSuccess("")
    setInfo("")
    setError("")
          
    let url = "guardian/"+cmVal+"/"+nsVal+"/"+sidVal
    console.log("fetching", url)
    fetch(url)
    .then( response => {
      const isJson = response.headers.get('content-type')?.includes('application/json');
      if (!response.ok) { 
        console.log("mmm... not ok!"); 
        throw response 
      }
      console.log("content-type", response.headers.get('content-type'))
      if (!isJson) { 
        console.log("mmm... not json!"); 
        throw new Error("Response is not a json...") 
      }
      return response.json()  //we only get here if there is no error
    })
    .then( json => {
      setData(json); 
      setSuccess("Succeasfuly read Guardian from cluster");
      if (collapse) {
        collapse()
      }
    })
    .catch( err => {
      console.log("mmm... fetch error...",err)
      setError("Error while reading Guardian from cluster")
    })
  };

  function onClearClick() {
    setSuccess("")
    setInfo("")
    setError("")
    setData(testData);
    setInfo("Resumed Guardian to default values");
    if (collapse) {
      collapse()
    }
  };
  function onLearnedClick() {
    setSuccess("")
    setInfo("")
    setError("")
    var d =  {}
    d.control = dataVal.control
    d.configured = JSON.parse(JSON.stringify(dataVal.learned));
    d.learned = dataVal.learned 
    setData(d);
    setInfo("Approved Learned Criteria to be used as Configured Criteria");
    if (collapse) {
      collapse()
    }
  };
  function onSetClick() {
    setSuccess("")
    setInfo("")
    setError("")
    console.log("onSetClick", dataVal)
    const requestOptions = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(dataVal)
    };
    let url = "guardian/"+cmVal+"/"+nsVal+"/"+sidVal
    fetch(url, requestOptions)
      .then(response => {
        console.log("fetch post response ",response.ok)
        if (response.ok) {
          setSuccess("Succeasfuly updated Guardian in cluster");
          if (collapse) {
            collapse()
          }
        } else {
          response.json().then(backendResponse => {
            if (backendResponse["message"]) {
              setError(backendResponse["message"])
            } else {
              setError("Error while trying to update Guardian in cluster")
            }
          })
        }
      })
      .catch( err => {
        console.log("mmm... fetch post error...",err)
        setError("Error while trying to update Guardian in cluster")
      })
  }
 
  function setCollapse(c) {
    console.log("setCollapse at APP")
    collapse = c
  }


	function onLoadHandler(event) {
		var reader = new FileReader();
    reader.onload = function(){
      	var bytes = reader.result;
        try {
          var d = yaml.load(bytes);
			  
          var ns = d.metadata.namespace
          var sid = d.metadata.name
        } catch (error) {
          console.error(error);
          alert("Failed to parse file")
          refLoad.current.value = "";
          return
        }
        

        if (ns && sid && d.apiVersion === "wsecurity.ibmresearch.com/v1" 
            && d.kind === "Guardian" && d.spec !== undefined) {
              setCm("crd")
              setNs(ns)
              setSid(sid)
              setData(d.spec)
              if (collapse) {
                collapse()
              }
              refLoad.current.value = "";
              return
        }
        if (ns && sid && d.apiVersion === "v1" && d.kind === "ConfigMap"
            && d.data !== undefined && d.data.Guardian !== undefined) {
              setCm("cm")
              setNs(ns)
              setSid(sid)
              var json = JSON.parse(d.data.Guardian)
              setData(json)
              
              if (collapse) {
                collapse()
              }
              refLoad.current.value = "";
              return
        }
        alert("File structure not recognised.")
        refLoad.current.value = "";
    };
    reader.readAsText(event.target.files[0]);
	};
	

    const refLoad = useRef()
    function onLoadClick() {
        refLoad.current.click()
    }

	function onSaveClick() {
    if (!sidVal || !nsVal || !cmVal) {
      alert("Identity not set")
      return
    }
    var json
    if (cmVal === "cm") {
      json = {
        apiVersion: "v1",
        kind: "ConfigMap",
        metadata: {name: sidVal, namespace:nsVal},
        data: {Guardian: JSON.stringify(dataVal)}
      }
    } else {
      json = {
        apiVersion: "wsecurity.ibmresearch.com/v1",
        kind: "Guardian",
        metadata: {name: sidVal, namespace:nsVal},
        spec: dataVal
      }   
    }
    var bytes  = yaml.dump(json);     
    
		const downloadLink = document.createElement('a');
		document.body.appendChild(downloadLink);

		downloadLink.href = "data:text/html;charset=utf-8,"+encodeURIComponent(bytes);
		downloadLink.target = '_self';
		downloadLink.download = cmVal+"."+nsVal+"."+sidVal+".yml";
		downloadLink.click(); 
	}
  function handleNsChange(event) {
    console.log("handleNsChange", event.target.value)
    setNs(event.target.value)
  }
  function handleSidChange(event) {
    console.log("handleSidChange", event.target.value)
    setSid(event.target.value)
  }
  function handleCmChange(event) {
    console.log("handleCmChange", event.target.value)
    setCm(event.target.value)
  }

  return (
    <div className="App">
      <Stack sx={{ width: '100%' }} spacing={2}>
      <Collapse in={errorVal !== ""}>
        <Alert severity="error" onClose={() => {setError("")}}>
          <AlertTitle>Error</AlertTitle>
          {errorVal}
        </Alert>
      </Collapse>
      <Collapse in={infoVal !== ""}>
        <Alert severity="info" onClose={() => {setInfo("")}}>
          <AlertTitle>Info</AlertTitle>
          {infoVal}
        </Alert>
      </Collapse>
      <Collapse in={successVal !== ""}>
        <Alert severity="success" onClose={() => {setSuccess("")}}>
          <AlertTitle>Success</AlertTitle>
          {successVal}
        </Alert>
      </Collapse>
      </Stack>
      <form>
      <Table  aria-label="customized table" sx={{ width: 'auto',  marginLeft: "2em", border:"hidden"}}>
        <TableHead  >
          <TableRow >
            <StyledTableCell sx={{color: "#d21976"}} colSpan={2} align="center">File System</StyledTableCell>
            <StyledTableCell sx={{color: "#19d275"}} colSpan={2} align="center"></StyledTableCell>
            <StyledTableCell sx={{color: "#d27519"}} colSpan={2} align="center">Kube Cluster</StyledTableCell>
          </TableRow>
        </TableHead>
        <TableBody>
            <TableRow sx={{  margin: "0em", padding: "0em", borderSpacing: "0"}}>
              <StyledTableCell>
                <StyledButton sx={{color: "#d21976"}} variant="outlined" endIcon={<GetIcon />} onClick={onLoadClick}>Load</StyledButton>
              </StyledTableCell>
              <StyledTableCell>
                <StyledButton sx={{color: "#d21976"}} variant="outlined" endIcon={<SendIcon />} onClick={onSaveClick}>Save</StyledButton>
              </StyledTableCell>
              <StyledTableCell>
                <StyledButton sx={{color: "#19d275"}}  variant="outlined" startIcon={<ClearIcon />} onClick={onClearClick}>Default</StyledButton>
              </StyledTableCell>
              <StyledTableCell>
                <StyledButton sx={{color: "#19d275"}}  variant="outlined" startIcon={<LearnIcon />} onClick={onLearnedClick}>Learned</StyledButton>
              </StyledTableCell>
              <StyledTableCell>
                <StyledButton sx={{color: "#d27519"}} variant="outlined" endIcon={<GetIcon />} onClick={onGetClick}>Get</StyledButton>
              </StyledTableCell>
              <StyledTableCell>
                <StyledButton sx={{color: "#d27519"}} variant="outlined" endIcon={<SendIcon />} onClick={onSetClick}>Set</StyledButton>
              </StyledTableCell>
            </TableRow>
            <TableRow>
              <TableCell align="center" colSpan={2}>
                <RadioGroup sx={{padding:0, paddingLeft:"3em" }} column  value={cmVal}  onChange={handleCmChange} name="crd-cm-group">
                  <FormControlLabel value="crd" control={<Radio sx={{padding:"5px" }} />} label={<Typography sx={{ fontSize: 12 }}>CRD</Typography>} />
                  <FormControlLabel value="cm" control={<Radio  sx={{padding:"5px" }} />} label={<Typography sx={{ fontSize: 12 }}>ConfigMap</Typography>} />
                </RadioGroup>
              </TableCell>
              <TableCell align="center" colSpan={2}>
                <TextField sx={{ margin: "1em"}} id="ns" label="Namespace" variant="standard" value={nsVal} onChange={handleNsChange} />
              </TableCell>
              <TableCell align="center" colSpan={2}>
                <TextField sx={{ margin: "1em"}} id="sid" label="Service Name" variant="standard"  value={sidVal} onChange={handleSidChange} />
              </TableCell>
              </TableRow>
        </TableBody>
      </Table>
      <input style={{ "display": "none" }}  accept=".yml,.yaml" ref={refLoad} type="file" onChange={onLoadHandler} />
	  		    
        
        
      <Guardian data={dataVal} setCollapse={setCollapse} >Guard</Guardian>
      </form>
    </div>
    
  );
}

//<Box sx={{ display: "flex", justifyContent: "start", margin: "0.2em"}}>
//            <Button sx={{ width: "12em"}} variant="outlined" endIcon={<GetIcon />} onClick={onLoadClick}>Load</Button>
//            <input style={{ "display": "none" }}  accept=".yml,.yaml" ref={refLoad} type="file" onChange={onLoadHandler} />
//	  		    <Button sx={{ width: "12em"}} variant="outlined" endIcon={<SendIcon />} onClick={onSaveClick}>Save</Button>
//            <Button sx={{ width: "12em"}} variant="outlined" startIcon={<ClearIcon />} onClick={onClearClick}>Default</Button>
//            <Button sx={{ width: "12em"}} variant="outlined" startIcon={<LearnIcon />} onClick={onLearnedClick}>Learned</Button>
//            <Button sx={{ width: "12em"}} variant="outlined" endIcon={<GetIcon />} onClick={onGetClick}>Get</Button>
//            <Button sx={{ width: "12em"}} variant="outlined" endIcon={<SendIcon />} onClick={onSetClick}>Set</Button>
//        </Box>
  
export default App;
