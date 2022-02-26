// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

contract Time{

 function getDate() private view returns(uint){
        uint _time = block.timestamp;
        return _time;
    }
    
    function callTime() external view returns(uint){
      uint tim = getDate();
      return(tim);
    }
}