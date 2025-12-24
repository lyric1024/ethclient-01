// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

contract Counter{

    uint256 private count;
    address public owner;

    event Added(address indexed by, uint256 newCount);

    modifier onlyOwner() {
        require(msg.sender == owner, "caller is not owner");
        _;
    }

    constructor() {
        count = 0;
        owner = msg.sender;
    }

    function getCount() public view returns (uint256) {
        return count;
    }

    function add() public {
        count++;
        emit Added(msg.sender, count);
    }

}