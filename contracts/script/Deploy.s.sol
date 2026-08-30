// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.20;

import {Script, console2} from "forge-std/Script.sol";
import {OracleConsumer} from "../src/OracleConsumer.sol";

contract DeployScript is Script {
    function run() external {
        vm.startBroadcast();
        OracleConsumer oracle = new OracleConsumer();
        vm.stopBroadcast();

        console2.log("OracleConsumer deployed at:", address(oracle));
    }
}
