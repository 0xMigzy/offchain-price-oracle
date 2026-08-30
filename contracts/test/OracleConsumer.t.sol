pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {OracleConsumer} from "../src/OracleConsumer.sol";

contract OracleConsumerTest is Test {
    OracleConsumer oracle;

    address owner = address(0x1);
    address submitter = address(0x2);
    address attacker = address(0x3);
    address newOwner = address(0x4);

    string constant PAIR = "ETH/USD";
    string constant PAIR2 = "BTC/USD";
    int256 constant PRICE = 2500e8;

    // Events must be redeclared here for vm.expectEmit to match against.
    event PriceSubmitted(string indexed pair, int256 price, uint256 timestamp, address indexed submitter);
    event OwnershipTransferStarted(address indexed from, address indexed to);
    event OwnershipTransferred(address indexed from, address indexed to);
    event Paused(address indexed by);
    event Unpaused(address indexed by);

    function setUp() public {
        vm.startPrank(owner);
        oracle = new OracleConsumer();
        vm.stopPrank();
    }

    /// @dev Shortcut for tests that don't care how the submitter got added.
    function _addSubmitter(address who) internal {
        vm.prank(owner);
        oracle.addSubmitter(who);
    }

    
    function test_OwnerCanAddSubmitter() public {
        vm.prank(owner);
        oracle.addSubmitter(submitter);
        assertTrue(oracle.submitters(submitter));
    }

    function test_SubmitterCanSubmitPrice() public {
        vm.prank(owner);
        oracle.addSubmitter(submitter);

        vm.prank(submitter);
        oracle.submitPrice("ETH/USD", 2500e8);

        (int256 price, uint256 ts, ) = oracle.lastPrices("ETH/USD");

        assertEq(price, 2500e8);
        assertEq(ts, block.timestamp);
    }

       function test_unauthorizedSubmitterReverts() public {
        vm.prank(attacker);
        vm.expectRevert(OracleConsumer.NotAuthorizedSubmitter.selector);
        oracle.submitPrice("ETH/USD", 2500e8);
    }

        function test_Constructor_SetsOwner() public view {
        assertEq(oracle.owner(), owner);
    }

    function test_Constructor_OwnerIsSubmitter() public view {
        assertTrue(oracle.submitters(owner));
    }

    function test_Constructor_StartsUnpaused() public view {
        assertFalse(oracle.paused());
    }

    function test_Constructor_NoPendingOwner() public view {
        assertEq(oracle.pendingOwner(), address(0));
    }

    function test_Constructor_EmitsOwnershipTransferred() public {
        vm.expectEmit(true, true, false, false);
        emit OwnershipTransferred(address(0), owner);
        vm.prank(owner);
        new OracleConsumer();
    }
    
        function test_AddSubmitter_RevertsWhenNotOwner() public {
        vm.prank(attacker);
        vm.expectRevert(OracleConsumer.NotOwner.selector);
        oracle.addSubmitter(attacker);
    }

    function test_RemoveSubmitter_OwnerCanRemove() public {
        _addSubmitter(submitter);
        vm.prank(owner);
        oracle.removeSubmitter(submitter);
        assertFalse(oracle.submitters(submitter));
    }

    function test_RemoveSubmitter_RevokedAddressCannotSubmit() public {
        _addSubmitter(submitter);
        vm.prank(owner);
        oracle.removeSubmitter(submitter);

        vm.prank(submitter);
        vm.expectRevert(OracleConsumer.NotAuthorizedSubmitter.selector);
        oracle.submitPrice(PAIR, PRICE);
    }

    function test_RemoveSubmitter_RevertsWhenNotOwner() public {
        _addSubmitter(submitter);
        vm.prank(attacker);
        vm.expectRevert(OracleConsumer.NotOwner.selector);
        oracle.removeSubmitter(submitter);
    }

    function test_RemoveSubmitter_OwnerCanRemoveThemselves() public {
        vm.prank(owner);
        oracle.removeSubmitter(owner);
        assertFalse(oracle.submitters(owner));
        assertEq(oracle.owner(), owner); // still admin, just can't submit
    }

        function test_TransferOwnership_SetsPendingOwner() public {
        vm.prank(owner);
        oracle.transferOwnership(newOwner);
        assertEq(oracle.pendingOwner(), newOwner);
    }

    function test_TransferOwnership_DoesNotChangeOwnerYet() public {
        vm.prank(owner);
        oracle.transferOwnership(newOwner);
        assertEq(oracle.owner(), owner);
    }

    function test_TransferOwnership_OldOwnerStillHasPowerUntilAccepted() public {
        vm.prank(owner);
        oracle.transferOwnership(newOwner);

        vm.prank(owner);
        oracle.pause();
        assertTrue(oracle.paused());
    }

    function test_TransferOwnership_RevertsOnZeroAddress() public {
        vm.prank(owner);
        vm.expectRevert(OracleConsumer.ZeroAddress.selector);
        oracle.transferOwnership(address(0));
    }

    function test_TransferOwnership_RevertsWhenNotOwner() public {
        vm.prank(attacker);
        vm.expectRevert(OracleConsumer.NotOwner.selector);
        oracle.transferOwnership(attacker);
    }

    function test_TransferOwnership_EmitsEvent() public {
        vm.expectEmit(true, true, false, false);
        emit OwnershipTransferStarted(owner, newOwner);
        vm.prank(owner);
        oracle.transferOwnership(newOwner);
    }

    function test_AcceptOwnership_TransfersOwner() public {
        vm.prank(owner);
        oracle.transferOwnership(newOwner);
        vm.prank(newOwner);
        oracle.acceptOwnership();
        assertEq(oracle.owner(), newOwner);
    }

    function test_AcceptOwnership_ClearsPendingOwner() public {
        vm.prank(owner);
        oracle.transferOwnership(newOwner);
        vm.prank(newOwner);
        oracle.acceptOwnership();
        assertEq(oracle.pendingOwner(), address(0));
    }

    function test_AcceptOwnership_OldOwnerLosesPower() public {
        vm.prank(owner);
        oracle.transferOwnership(newOwner);
        vm.prank(newOwner);
        oracle.acceptOwnership();

        vm.prank(owner);
        vm.expectRevert(OracleConsumer.NotOwner.selector);
        oracle.pause();
    }

    function test_AcceptOwnership_RevertsForWrongCaller() public {
        vm.prank(owner);
        oracle.transferOwnership(newOwner);

        vm.prank(attacker);
        vm.expectRevert(OracleConsumer.NotPendingOwner.selector);
        oracle.acceptOwnership();
    }

    function test_AcceptOwnership_RevertsWhenNoTransferPending() public {
        vm.prank(newOwner);
        vm.expectRevert(OracleConsumer.NotPendingOwner.selector);
        oracle.acceptOwnership();
    }

    function test_AcceptOwnership_CannotBeReplayed() public {
        vm.prank(owner);
        oracle.transferOwnership(newOwner);
        vm.prank(newOwner);
        oracle.acceptOwnership();

        vm.prank(newOwner);
        vm.expectRevert(OracleConsumer.NotPendingOwner.selector);
        oracle.acceptOwnership();
    }

    function test_AcceptOwnership_EmitsEvent() public {
        vm.prank(owner);
        oracle.transferOwnership(newOwner);

        vm.expectEmit(true, true, false, false);
        emit OwnershipTransferred(owner, newOwner);
        vm.prank(newOwner);
        oracle.acceptOwnership();
    }

        function test_Pause_SetsPaused() public {
        vm.prank(owner);
        oracle.pause();
        assertTrue(oracle.paused());
    }

    function test_Pause_BlocksSubmissions() public {
        _addSubmitter(submitter);
        vm.prank(owner);
        oracle.pause();

        vm.prank(submitter);
        vm.expectRevert(OracleConsumer.ContractPaused.selector);
        oracle.submitPrice(PAIR, PRICE);
    }

    function test_Pause_ReadsStillWork() public {
        _addSubmitter(submitter);
        vm.prank(submitter);
        oracle.submitPrice(PAIR, PRICE);

        vm.prank(owner);
        oracle.pause();

        (int256 price,) = oracle.getLastPrice(PAIR);
        assertEq(price, PRICE);
    }

    function test_Pause_RevertsWhenAlreadyPaused() public {
        vm.startPrank(owner);
        oracle.pause();
        vm.expectRevert(OracleConsumer.NoChange.selector);
        oracle.pause();
        vm.stopPrank();
    }

    function test_Pause_RevertsWhenNotOwner() public {
        vm.prank(attacker);
        vm.expectRevert(OracleConsumer.NotOwner.selector);
        oracle.pause();
    }

    function test_Pause_EmitsEvent() public {
        vm.expectEmit(true, false, false, false);
        emit Paused(owner);
        vm.prank(owner);
        oracle.pause();
    }

    function test_Unpause_ResumesSubmissions() public {
        _addSubmitter(submitter);
        vm.startPrank(owner);
        oracle.pause();
        oracle.unpause();
        vm.stopPrank();

        vm.prank(submitter);
        oracle.submitPrice(PAIR, PRICE);

        (int256 price,) = oracle.getLastPrice(PAIR);
        assertEq(price, PRICE);
    }

    function test_Unpause_RevertsWhenNotPaused() public {
        vm.prank(owner);
        vm.expectRevert(OracleConsumer.NoChange.selector);
        oracle.unpause();
    }

    function test_Unpause_RevertsWhenNotOwner() public {
        vm.prank(owner);
        oracle.pause();

        vm.prank(attacker);
        vm.expectRevert(OracleConsumer.NotOwner.selector);
        oracle.unpause();
    }

    function test_Unpause_EmitsEvent() public {
        vm.prank(owner);
        oracle.pause();

        vm.expectEmit(true, false, false, false);
        emit Unpaused(owner);
        vm.prank(owner);
        oracle.unpause();
    }

        function test_SubmitPrice_OwnerCanSubmit() public {
        vm.prank(owner);
        oracle.submitPrice(PAIR, PRICE);
        (int256 price,) = oracle.getLastPrice(PAIR);
        assertEq(price, PRICE);
    }

    function test_SubmitPrice_RevertsOnEmptyPair() public {
        _addSubmitter(submitter);
        vm.prank(submitter);
        vm.expectRevert(OracleConsumer.EmptyPair.selector);
        oracle.submitPrice("", PRICE);
    }

    function test_SubmitPrice_RevertsOnZeroPrice() public {
        _addSubmitter(submitter);
        vm.prank(submitter);
        vm.expectRevert(OracleConsumer.InvalidPrice.selector);
        oracle.submitPrice(PAIR, 0);
    }

    function test_SubmitPrice_RevertsOnNegativePrice() public {
        _addSubmitter(submitter);
        vm.prank(submitter);
        vm.expectRevert(OracleConsumer.InvalidPrice.selector);
        oracle.submitPrice(PAIR, -1);
    }

    function test_SubmitPrice_EmitsEvent() public {
        _addSubmitter(submitter);

        vm.expectEmit(true, true, false, true);
        emit PriceSubmitted(PAIR, PRICE, block.timestamp, submitter);

        vm.prank(submitter);
        oracle.submitPrice(PAIR, PRICE);
    }

    function test_SubmitPrice_RecordsTimestampCorrectly() public {
        _addSubmitter(submitter);
        vm.warp(1_700_000_000);

        vm.prank(submitter);
        oracle.submitPrice(PAIR, PRICE);

        (, uint256 ts) = oracle.getLastPrice(PAIR);
        assertEq(ts, 1_700_000_000);
    }

        function test_SubmitPrice_SecondSubmissionOverwritesLast() public {
        _addSubmitter(submitter);
        vm.startPrank(submitter);
        oracle.submitPrice(PAIR, PRICE);
        oracle.submitPrice(PAIR, 2600e8);
        vm.stopPrank();

        (int256 price,) = oracle.getLastPrice(PAIR);
        assertEq(price, 2600e8);
    }

    function test_SubmitPrice_HistoryGrows() public {
        _addSubmitter(submitter);
        assertEq(oracle.getPriceHistory(PAIR), 0);

        vm.startPrank(submitter);
        oracle.submitPrice(PAIR, PRICE);
        assertEq(oracle.getPriceHistory(PAIR), 1);

        oracle.submitPrice(PAIR, 2600e8);
        assertEq(oracle.getPriceHistory(PAIR), 2);
        vm.stopPrank();
    }

    function test_SubmitPrice_HistoryKeepsOldEntries() public {
        _addSubmitter(submitter);
        vm.startPrank(submitter);
        oracle.submitPrice(PAIR, PRICE);
        oracle.submitPrice(PAIR, 2600e8);
        vm.stopPrank();

        (int256 first,,) = oracle.priceHistory(PAIR, 0);
        (int256 second,,) = oracle.priceHistory(PAIR, 1);
        assertEq(first, PRICE);
        assertEq(second, 2600e8);
    }

    function test_SubmitPrice_PairsAreIsolated() public {
        _addSubmitter(submitter);
        vm.startPrank(submitter);
        oracle.submitPrice(PAIR, PRICE);
        oracle.submitPrice(PAIR2, 60000e8);
        vm.stopPrank();

        (int256 eth,) = oracle.getLastPrice(PAIR);
        (int256 btc,) = oracle.getLastPrice(PAIR2);
        assertEq(eth, PRICE);
        assertEq(btc, 60000e8);
        assertEq(oracle.getPriceHistory(PAIR), 1);
    }

    function test_SubmitPrice_MultipleSubmittersRecordedSeparately() public {
        _addSubmitter(submitter);

        vm.prank(submitter);
        oracle.submitPrice(PAIR, PRICE);
        vm.prank(owner);
        oracle.submitPrice(PAIR, 2600e8);

        (,, address latest) = oracle.lastPrices(PAIR);
        assertEq(latest, owner);

        (,, address earlier) = oracle.priceHistory(PAIR, 0);
        assertEq(earlier, submitter);
    }

        function test_GetLastPrice_ReturnsSubmittedValues() public {
        _addSubmitter(submitter);
        vm.warp(1_700_000_000);
        vm.prank(submitter);
        oracle.submitPrice(PAIR, PRICE);

        (int256 price, uint256 ts) = oracle.getLastPrice(PAIR);
        assertEq(price, PRICE);
        assertEq(ts, 1_700_000_000);
    }

    /// @dev Documents current behaviour, not desired behaviour: an unknown pair
    ///      is indistinguishable from a real price of zero. If you later add a
    ///      NoPriceAvailable revert, flip this test to expect it.
    function test_GetLastPrice_UnknownPairReturnsZeroes() public view {
        (int256 price, uint256 ts) = oracle.getLastPrice("DOGE/MARS");
        assertEq(price, 0);
        assertEq(ts, 0);
    }

    function test_GetPriceHistory_UnknownPairReturnsZero() public view {
        assertEq(oracle.getPriceHistory("DOGE/MARS"), 0);
    }

        function testFuzz_SubmitPrice_RejectsNonPositive(int256 price) public {
        vm.assume(price <= 0);
        _addSubmitter(submitter);

        vm.prank(submitter);
        vm.expectRevert(OracleConsumer.InvalidPrice.selector);
        oracle.submitPrice(PAIR, price);
    }

    function testFuzz_SubmitPrice_AcceptsAnyPositive(int256 price) public {
        vm.assume(price > 0);
        _addSubmitter(submitter);

        vm.prank(submitter);
        oracle.submitPrice(PAIR, price);

        (int256 stored,) = oracle.getLastPrice(PAIR);
        assertEq(stored, price);
    }

    function testFuzz_OnlyOwnerCanAddSubmitter(address caller) public {
        vm.assume(caller != owner);

        vm.prank(caller);
        vm.expectRevert(OracleConsumer.NotOwner.selector);
        oracle.addSubmitter(attacker);
    }

    function testFuzz_OnlySubmittersCanSubmit(address caller) public {
        vm.assume(caller != owner);

        vm.prank(caller);
        vm.expectRevert(OracleConsumer.NotAuthorizedSubmitter.selector);
        oracle.submitPrice(PAIR, PRICE);
    }

    function testFuzz_AnyPairNameWorks(string calldata pair) public {
        vm.assume(bytes(pair).length > 0);
        _addSubmitter(submitter);

        vm.prank(submitter);
        oracle.submitPrice(pair, PRICE);

        (int256 stored,) = oracle.getLastPrice(pair);
        assertEq(stored, PRICE);
    }

    

}