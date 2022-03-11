//import logo from './logo.svg';
import './App.css';
import {Guardian} from './Guardian';
import {Identity} from './Identity'
import {Button, Box, Collapse, Alert, AlertTitle, Stack} from '@mui/material';
import React, { useState } from "react";
import GetIcon from '@mui/icons-material/GetApp';
import SendIcon from '@mui/icons-material/Send';
import ClearIcon from '@mui/icons-material/Clear';
import LearnIcon from '@mui/icons-material/LightbulbOutlined';
import ExpandMoreIcon from '@mui/icons-material/ExpandMoreRounded';
import ExpandLessIcon from '@mui/icons-material/ExpandLessRounded';

const testData = {}
const testData1 = {
  "req": {
    "url": {
      "val": {
        "flags": 32768,
        "runes": [
          {
            "min": 0,
            "max": 64
          }
        ],
        "digits": [
          {
            "min": 0,
            "max": 64
          },
          {
            "min": 20,
            "max": 26
          }
        ],
        "letters": [
          {
            "min": 0,
            "max": 64
          },
          {
            "min": 70,
            "max": 80
          }
        ],
        "schars": null,
        "words": [
          {
            "min": 0,
            "max": 16
          }
        ],
        "numbers": [
          {
            "min": 0,
            "max": 16
          }
        ],
        "unicodeFlags": [32124, 32124, 32124]
      },
      "segments": [
        {
          "min": 0,
          "max": 8
        }
      ]
    },
    "qs": {
      "kv": {
        "vals": null,
        "minimalSet": null,
        "otherVals": {
          "flags": 0,
          "runes": [
            {
              "min": 0,
              "max": 32
            }
          ],
          "digits": [
            {
              "min": 0,
              "max": 32
            }
          ],
          "letters": [
            {
              "min": 0,
              "max": 32
            }
          ],
          "schars": null,
          "words": [
            {
              "min": 0,
              "max": 16
            }
          ],
          "numbers": [
            {
              "min": 0,
              "max": 16
            }
          ],
          "unicodeFlags": null
        },
        "otherKeynames": {
          "flags": 0,
          "runes": [
            {
              "min": 0,
              "max": 16
            }
          ],
          "digits": [
            {
              "min": 0,
              "max": 16
            }
          ],
          "letters": [
            {
              "min": 0,
              "max": 16
            }
          ],
          "schars": null,
          "words": [
            {
              "min": 0,
              "max": 4
            }
          ],
          "numbers": [
            {
              "min": 0,
              "max": 4
            }
          ],
          "unicodeFlags": null
        }
      }
    },
    "headers": {
      "kv": {
        "vals": {
          "Accept": {
            "flags": 34359796736,
          "runes": [
            {
              "min": 0,
              "max": 32
            }
          ],
          "digits": [
            {
              "min": 0,
              "max": 32
            }
          ],
          "letters": [
            {
              "min": 0,
              "max": 32
            }
          ],
          "schars": [
            {
              "min": 0,
              "max": 8
            }
          ],
          "words": [
            {
              "min": 0,
              "max": 16
            }
          ],
          "numbers": [
            {
              "min": 0,
              "max": 16
            }
          ],
          "unicodeFlags": null
          }
        },
        "minimalSet": null,
        "otherVals": {
          "flags": 34359796736,
          "runes": [
            {
              "min": 0,
              "max": 32
            }
          ],
          "digits": [
            {
              "min": 0,
              "max": 32
            }
          ],
          "letters": [
            {
              "min": 0,
              "max": 32
            }
          ],
          "schars": [
            {
              "min": 0,
              "max": 8
            }
          ],
          "words": [
            {
              "min": 0,
              "max": 16
            }
          ],
          "numbers": [
            {
              "min": 0,
              "max": 16
            }
          ],
          "unicodeFlags": null
        },
        "otherKeynames": {
          "flags": 0,
          "runes": [
            {
              "min": 0,
              "max": 16
            }
          ],
          "digits": [
            {
              "min": 0,
              "max": 16
            },
            {
              "min": 20,
              "max": 26
            }
          ],
          "letters": [
            {
              "min": 0,
              "max": 16
            }
          ],
          "schars": [
            {
              "min": 0,
              "max": 2
            }
          ],
          "words": [
            {
              "min": 0,
              "max": 4
            }
          ],
          "numbers": [
            {
              "min": 0,
              "max": 4
            }
          ],
          "unicodeFlags": [32124, 32124, 32124]
        }
      }
    }
  },
  "consult": {
    "active": false,
    "rpm": 0
  },
  "forceAllow": false
}
function App() {
  const [dataVal, setData] = useState(testData);
  const [successVal, setSuccess] = useState("");
  const [errorVal, setError] = useState("");
  const [infoVal, setInfo] = useState("");
    
  let ns = ""
  let sid = ""
  let collapse, expand

  function onGetClick() {
    setSuccess("")
    setInfo("")
    setError("")

    let url = "guardian/"+ns+"/"+sid
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
    /*
    let success = fetchData();
    if (success) {
      setSuccess("Succeasfuly read Guardian from cluster");
    } else {
      setError("Error while reading Guardian from cluster")
    }
    */
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
    setInfo("Approved Learned Critiria to be used as Configured Critiria");
    if (collapse) {
      collapse()
    }
  };
  function onSetClick() {
    setSuccess("")
    setInfo("")
    setError("")
    var d =  {}
    d.control = dataVal.control
    d.configured = dataVal.configured
    console.log("onSetClick", d)
    const requestOptions = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(d)
    };
    let url = "guardian/"+ns+"/"+sid
    fetch(url, requestOptions)
      .then(response => {
        console.log("fetch post response ",response.ok)
        if (response.ok) {
          setSuccess("Succeasfuly updated Guardian in cluster");
          if (collapse) {
            collapse()
          }
        } else {
          setError("Error while trying to update Guardian in cluster")
        }
      })
      .catch( err => {
        console.log("mmm... fetch post error...",err)
        setError("Error while trying to update Guardian in cluster")
      })
/*
    let success = postData(dataVal)
    if (success) {
      
    } else {
      
    }
    if (collapse) {
      collapse()
    }
    */
  }

  function handleIdentityChange(nsVal, sidVal) {
    ns = nsVal
    sid = sidVal
  }
 
  function setCollapse(c) {
    console.log("setCollapse at APP")
    collapse = c
  }
  function setExpand(e) {
    console.log("setExpand at APP")
    expand = e
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
        <Box sx={{ display: "flex", justifyContent: "start", margin: "0.2em"}}>
            <Button sx={{ width: "12em"}} variant="outlined" endIcon={<GetIcon />} onClick={onGetClick}>Get</Button>
            <Button sx={{ width: "12em"}} variant="outlined" startIcon={<ClearIcon />} onClick={onClearClick}>Default</Button>
            <Button sx={{ width: "12em"}} variant="outlined" startIcon={<LearnIcon />} onClick={onLearnedClick}>Learned</Button>
            <Button sx={{ width: "12em"}} variant="outlined" endIcon={<SendIcon />} onClick={onSetClick}>Set</Button>
        </Box>
        <Identity ns="" sid="" handleChange={handleIdentityChange} ></Identity>
        
        
      <Guardian data={dataVal} setCollapse={setCollapse} setExpand={setExpand}>Guard</Guardian>
      </form>
    </div>
    
  );
}

export default App;
