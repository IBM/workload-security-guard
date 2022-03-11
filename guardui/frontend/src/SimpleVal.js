import React, {useEffect} from "react";
import { U8MinmaxSlice } from './U8MinmaxSlice';
import { FlagsSlice } from './FlagsSlice';
import { UnicodeSlice } from './UnicodeSlice';
import {Toggle} from './Guardian'


import TreeItem from '@mui/lab/TreeItem';

function SimpleVal(props) {
  var { data, onDataChange, nodeId, keyId, name } = props;
  console.log("SimpleVal data", data, name)
  
  if (!data.digits) data.digits = []
  if (!data.letters) data.letters = []
  if (!data.runes) data.runes = []
  if (!data.schars) data.schars = []
  if (!data.numbers) data.numbers = []
  if (!data.words) data.words = []
  if (!data.flags) data.flags = 0
  if (!data.unicodeFlags) data.unicodeFlags = []

  useEffect(() => {
    Toggle([nodeId+">Digits", nodeId+">Letters", nodeId+">Runes", nodeId+">Special Chars", 
            nodeId+">Numbers", nodeId+">Words", nodeId+">Flags", nodeId+">Unicodes"])
  }, [nodeId]);

  function onDigitsChange(d) {
    data.digits = d 
    console.log("onDigitsChange", d, data)
    onDataChange(keyId,data)
  }
  function onLettersChange(d) {
    data.letters = d 
    console.log("onLettersChange", d, data)
    onDataChange(keyId,data)
  }
  function onRunesChange(d) {
    data.runes = d 
    console.log("onRunesChange", d, data)
    onDataChange(keyId,data)
  }
  function onSCharsChange(d) {
    data.schars = d 
    console.log("onSCharsChange", d, data)
    onDataChange(keyId,data)
  }
  function onNumbersChange(d) {
    data.numbers = d 
    console.log("onNumbersChange", d, data)
    onDataChange(keyId,data)
  }
  function onWordsChange(d) {
    data.words = d 
    console.log("onWordsChange", d, data)
    onDataChange(keyId,data)
  }
  function onFlagsChange(d) {
    data.flags = d 
    console.log("onFlagsChange", d, data)
    onDataChange(keyId,data)
  }
  function onUnicodesChange(d) {
    data.unicodeFlags = d 
    console.log("onUnicodesChange", d, data)
    onDataChange(keyId,data)
  }
  
  console.log("SimpleVal", data)

  return (
  
    
      <TreeItem nodeId={nodeId} label={name}>
        <U8MinmaxSlice data={data.digits} nodeId={nodeId+">Digits"} name="Digits" onDataChange={onDigitsChange}></U8MinmaxSlice>
        <U8MinmaxSlice data={data.letters} nodeId={nodeId+">Letters"} name="Letters" onDataChange={onLettersChange}></U8MinmaxSlice>
        <U8MinmaxSlice data={data.runes} nodeId={nodeId+">Runes"} name="Runes" onDataChange={onRunesChange}></U8MinmaxSlice>
        <U8MinmaxSlice data={data.schars} nodeId={nodeId+">Special Chars"} name="Special Chars" onDataChange={onSCharsChange}></U8MinmaxSlice>
        <U8MinmaxSlice data={data.numbers} nodeId={nodeId+">Numbers"} name="Numbers" onDataChange={onNumbersChange}></U8MinmaxSlice>
        <U8MinmaxSlice data={data.words} nodeId={nodeId+">Words"} name="Words" onDataChange={onWordsChange}></U8MinmaxSlice>
        <FlagsSlice data={data.flags} nodeId={nodeId+">Flags"} name="Flags" onDataChange={onFlagsChange}></FlagsSlice>
        <UnicodeSlice data={data.unicodeFlags} nodeId={nodeId+">Unicodes"} name="Unicodes" onDataChange={onUnicodesChange}></UnicodeSlice>
      </TreeItem>
 
        );
}

export {SimpleVal};
// </div><Box sx={{ display: "flex", justifyContent: "start", margin: "0.2em"}}>
// <Button sx={{ width: "12em", alignItems: "center", fontSize: "0.8em"}} onClick={handleChange} variant="contained">{name}</Button>
// </Box><Box sx={{ display: showVal ? "flex" : "none", flexDirection:  "column", justifyContent: "start"}}>
