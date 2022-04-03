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
  //if (!data.runes) data.runes = []
  if (!data.schars) data.schars = []
  if (!data.spaces) data.spaces = []
  if (!data.unicodes) data.unicodes = []
  if (!data.nonreadables) data.nonreadables = []
  if (!data.sequences) data.sequences = []
  

  //if (!data.numbers) data.numbers = []
  //if (!data.words) data.words = []
  if (!data.flags) data.flags = 0
  if (!data.unicodeFlags) data.unicodeFlags = []

  useEffect(() => {
    Toggle([nodeId+">Digits", nodeId+">Letters", nodeId+">Spaces", nodeId+">Special Chars", 
            nodeId+">Sequences",  nodeId+">Unicodes",  nodeId+">NonReadables", nodeId+">Flags", nodeId+">UnicodeTypes"])
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
  //function onRunesChange(d) {
  //  data.runes = d 
  //  console.log("onRunesChange", d, data)
  //  onDataChange(keyId,data)
  //}
  function onSpecialCharsChange(d) {
    data.schars = d 
    console.log("onSCharsChange", d, data)
    onDataChange(keyId,data)
  }
  function onSequencesChange(d) {
    data.sequences = d 
    console.log("onSequencesChange", d, data)
    onDataChange(keyId,data)
  }
  //function onNumbersChange(d) {
  //  data.numbers = d 
  //  console.log("onNumbersChange", d, data)
  //  onDataChange(keyId,data)
  //}
  //function onWordsChange(d) {
  //  data.words = d 
  //  console.log("onWordsChange", d, data)
  //  onDataChange(keyId,data)
  //}
  function onFlagsChange(d) {
    data.flags = d 
    console.log("onFlagsChange", d, data)
    onDataChange(keyId,data)
  }
  function onUnicodeTypesChange(d) {
    data.unicodeFlags = d 
    console.log("onUnicodesChange", d, data)
    onDataChange(keyId,data)
  }
  function onSpacesChange(d) {
    data.spaces = d 
    console.log("onSpacesChange", d, data)
    onDataChange(keyId,data)
  }
  function onNonReadablesChange(d) {
    data.nonreadables = d 
    console.log("onNonReadablesChange", d, data)
    onDataChange(keyId,data)
  }
  function onUnicodesChange(d) {
    data.unicodes = d 
    console.log("onUnicodesChange", d, data)
    onDataChange(keyId,data)
  }
  
  console.log("SimpleVal", data)

  return (
  
    
      <TreeItem nodeId={nodeId} label={name}>
        <TreeItem nodeId={nodeId+">Counters"} label={"Rune Counters"}>
          <U8MinmaxSlice data={data.digits} nodeId={nodeId+">Digits"} name="Digits 0-9" description="Number of digits (runes 0-9)" onDataChange={onDigitsChange}></U8MinmaxSlice>
          <U8MinmaxSlice data={data.letters} nodeId={nodeId+">Letters"} name="Letters A-Z a-z" description="Number of letters (runes A-Z a-z)" onDataChange={onLettersChange}></U8MinmaxSlice>
          < U8MinmaxSlice data={data.runes} nodeId={nodeId+">Spaces"} name="Spaces" description="Number of space runes" onDataChange={onSpacesChange}></U8MinmaxSlice>
          <U8MinmaxSlice data={data.schars} nodeId={nodeId+">SpecialChars"} name="Special Chars" description="Number of Special Charcter runes" onDataChange={onSpecialCharsChange}></U8MinmaxSlice>
          <U8MinmaxSlice data={data.unicodes} nodeId={nodeId+">Unicodes"} name="Unicodes" description="Number of Unicode runes" onDataChange={onUnicodesChange}></U8MinmaxSlice>
          <U8MinmaxSlice data={data.nonreadables} nodeId={nodeId+">NonReadables"} name="Non Readables" description="Number of Non Readable runes" onDataChange={onNonReadablesChange}></U8MinmaxSlice>
        </TreeItem>
        <U8MinmaxSlice data={data.sequences} nodeId={nodeId+">Sequences"} name="Sequences" description="Number of sequances of runes of the same type (e.g. letters: aabc digits: 987 special chars: !?!)" onDataChange={onSequencesChange}></U8MinmaxSlice>
        <FlagsSlice data={data.flags} nodeId={nodeId+">Flags"} name="Special Char Types" onDataChange={onFlagsChange}></FlagsSlice>
        <UnicodeSlice data={data.unicodeFlags} nodeId={nodeId+">UnicodeTypes"} name="Unicode Types" onDataChange={onUnicodeTypesChange}></UnicodeSlice>
      </TreeItem>
 
        );
}

export {SimpleVal};
// </div><Box sx={{ display: "flex", justifyContent: "start", margin: "0.2em"}}>
// <Button sx={{ width: "12em", alignItems: "center", fontSize: "0.8em"}} onClick={handleChange} variant="contained">{name}</Button>
// </Box><Box sx={{ display: showVal ? "flex" : "none", flexDirection:  "column", justifyContent: "start"}}>
