pragma solidity ^0.8.20;

contract OracleConsumer {
    struct PriceEntry {
        int256 price;
        uint256 timestamp;
        address submitter;
    }

    mapping(string => PriceEntry) public lastPrices;
    mapping(string => PriceEntry[]) public priceHistory;
    mapping(address => bool) public submitters;

    address public owner;
    address public pendingOwner;
    bool public paused;

    event PriceSubmitted(string indexed pair, int256 price, uint256 timestamp, address indexed submitter);
    event OwnershipTransferStarted(address indexed from, address indexed to);
    event OwnershipTransferred(address indexed from, address indexed to);
    event Paused(address indexed by);
    event Unpaused(address indexed by);

    error NotOwner();
    error NotAuthorizedSubmitter();
    error InvalidPrice();
    error EmptyPair();
    error NotPendingOwner();
    error ContractPaused();
    error ZeroAddress();
    error NoChange();

    constructor() {
        owner = msg.sender;
        submitters[msg.sender] = true; // Owner is an authorized submitter
        emit OwnershipTransferred(address(0), msg.sender);
    }

    modifier onlyOwner() {
        if (msg.sender != owner) {
            revert NotOwner();
        }
        _;
    }

    modifier onlySubmitter() {
        if (!submitters[msg.sender]) {
            revert NotAuthorizedSubmitter();
        }
        _;
    }

    modifier whenNotPaused() {
        if (paused) revert ContractPaused();
        _;
    }

    function addSubmitter(address _submitter) external onlyOwner {
        submitters[_submitter] = true;
    }

    function removeSubmitter(address _submitter) external onlyOwner {
        submitters[_submitter] = false;
    }

    /// @notice Nominate a new owner. Nothing changes until they accept.
    /// @dev Two-step so a mistyped address can't lock the contract forever.
    function transferOwnership(address newOwner) external onlyOwner {
        if (newOwner == address(0)) revert ZeroAddress();
        pendingOwner = newOwner;
        emit OwnershipTransferStarted(owner, newOwner);
    }

    /// @notice Claim ownership. Only the nominated address can call this.
    function acceptOwnership() external {
        if (msg.sender != pendingOwner) revert NotPendingOwner();
        address previous = owner;
        owner = pendingOwner;
        pendingOwner = address(0);
        emit OwnershipTransferred(previous, owner);
    }

    /// @notice Halt new submissions. Reads keep working so consumers can still
    ///         see the last price and judge its age themselves.
    function pause() external onlyOwner {
        if (paused) revert NoChange();
        paused = true;
        emit Paused(msg.sender);
    }

    function unpause() external onlyOwner {
        if (!paused) revert NoChange();
        paused = false;
        emit Unpaused(msg.sender);
    }

    function submitPrice(string calldata pair, int256 price) external onlySubmitter whenNotPaused {
        if (bytes(pair).length == 0) {
            revert EmptyPair();
        }
        if (price <= 0) {
            revert InvalidPrice();
        }

        PriceEntry memory newEntry = PriceEntry({price: price, timestamp: block.timestamp, submitter: msg.sender});

        lastPrices[pair] = newEntry;
        priceHistory[pair].push(newEntry);

        emit PriceSubmitted(pair, price, block.timestamp, msg.sender);
    }

    function getLastPrice(string calldata pair) external view returns (int256, uint256) {
        PriceEntry memory entry = lastPrices[pair];
        return (entry.price, entry.timestamp);
    }

    function getPriceHistory(string calldata pair) external view returns (uint256) {
        return priceHistory[pair].length;
    }
}
