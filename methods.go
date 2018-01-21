package portal

type method struct {
	// This method does not require user authorization
	NoAuth bool
}

var rpcMethods = map[string]method{
	// == Blockchain ==
	"callcontract":          {true},
	"getaccountinfo":        {true},
	"getbestblockhash":      {true},
	"getblock":              {true},
	"getblockchaininfo":     {true},
	"getblockcount":         {true},
	"getblockhash":          {true},
	"getblockheader":        {true},
	"getchaintips":          {true},
	"getdifficulty":         {true},
	"getmempoolancestors":   {true},
	"getmempooldescendants": {true},
	"getmempoolentry":       {true},
	"getmempoolinfo":        {true},
	"getrawmempool":         {true},
	"gettransactionreceipt": {true},
	"gettxout":              {true},
	"gettxoutproof":         {true},
	"gettxoutsetinfo":       {true},
	"listcontracts":         {true},
	// preciousblock "blockhash"
	// pruneblockchain
	// verifychain ( checklevel nblocks )
	"verifytxoutproof": {true},

	// == Control ==
	// getinfo
	"getinfo": {true},
	// getmemoryinfo
	// help ( "command" )
	// stop

	// == Generating ==
	// generate nblocks ( maxtries )
	// generatetoaddress nblocks address (maxtries)

	// == Mining ==
	// getblocktemplate ( TemplateRequest )
	// getmininginfo
	"getmininginfo": {true},
	// getnetworkhashps ( nblocks height )
	// getstakinginfo
	"getstakinginfo": {true},
	// getsubsidy [nTarget]
	// prioritisetransaction <txid> <priority delta> <fee delta>
	// submitblock "hexdata" ( "jsonparametersobject" )

	// == Network ==
	// addnode "node" "add|remove|onetry"
	// clearbanned
	// disconnectnode "node"
	// getaddednodeinfo ( "node" )
	// getconnectioncount
	// getnettotals
	"getnettotals": {true},
	// getnetworkinfo
	"getnetworkinfo": {true},
	// getpeerinfo
	"getpeerinfo": {true},
	// listbanned
	// ping
	// setban "subnet" "add|remove" (bantime) (absolute)
	// setnetworkactive true|false

	// == Rawtransactions ==
	// createrawtransaction [{"txid":"id","vout":n},...] {"address":amount,"data":"hex",...} ( locktime )
	// decoderawtransaction "hexstring"
	// decodescript "hexstring"
	// fromhexaddress "hexaddress"
	"fromhexaddress": {true},
	// fundrawtransaction "hexstring" ( options )
	// gethexaddress "address"
	"gethexaddress": {true},
	// getrawtransaction "txid" ( verbose )
	// sendrawtransaction "hexstring" ( allowhighfees )
	// signrawtransaction "hexstring" ( [{"txid":"id","vout":n,"scriptPubKey":"hex","redeemScript":"hex"},...] ["privatekey1",...] sighashtype )

	// == Util ==
	// createmultisig nrequired ["key",...]
	// estimatefee nblocks
	// estimatepriority nblocks
	// estimatesmartfee nblocks
	// estimatesmartpriority nblocks
	// signmessagewithprivkey "privkey" "message"
	// validateaddress "address"
	// verifymessage "address" "signature" "message"

	// == Wallet ==
	// abandontransaction "txid"
	// addmultisigaddress nrequired ["key",...] ( "account" )
	// addwitnessaddress "address"
	// backupwallet "destination"
	// bumpfee "txid" ( options )
	// createcontract "bytecode" (gaslimit gasprice "senderaddress" broadcast)
	"createcontract": {false},
	// dumpprivkey "address"
	// dumpwallet "filename"
	// encryptwallet "passphrase"
	// getaccount "address"
	"getaccount": {true},
	// DEPRECATED getaccountaddress "account"
	// DEPRECATED getaddressesbyaccount "account"
	// getbalance ( "account" minconf include_watchonly )
	"getbalance": {true},
	// getnewaddress ( "account" )
	"getnewaddress": {false},
	// getrawchangeaddress
	// DEPRECATED getreceivedbyaccount "account" ( minconf )
	// getreceivedbyaddress "address" ( minconf )
	"getreceivedbyaddress": {true},
	// gettransaction "txid" ( include_watchonly )
	"gettransaction": {true},
	// getunconfirmedbalance
	"getunconfirmedbalance": {true},
	// getwalletinfo
	"getwalletinfo": {true},
	// importaddress "address" ( "label" rescan p2sh )
	// importmulti "requests" "options"
	// importprivkey "qtum" ( "label" ) ( rescan )
	// importprunedfunds
	// importpubkey "pubkey" ( "label" rescan )
	// importwallet "filename"
	// keypoolrefill ( newsize )
	"keypoolrefill": {false},
	// DEPRECATED listaccounts ( minconf include_watchonly)
	// listaddressgroupings
	// listlockunspent
	// listreceivedbyaccount ( minconf include_empty include_watchonly)
	// listreceivedbyaddress ( minconf include_empty include_watchonly)
	"listreceivedbyaddress": {false},
	// listsinceblock ( "blockhash" target_confirmations include_watchonly)
	"listsinceblock": {true},
	// listtransactions ( "account" count skip include_watchonly)
	"listtransactions": {true},
	// listunspent ( minconf maxconf  ["addresses",...] [include_unsafe] )
	"listunspent": {false},
	// lockunspent unlock ([{"txid":"txid","vout":n},...])
	// DEPRECATED move "fromaccount" "toaccount" amount ( minconf "comment" )
	// removeprunedfunds "txid"
	// reservebalance [<reserve> [amount]]
	// DEPRECATED sendfrom "fromaccount" "toaddress" amount ( minconf "comment" "comment_to" )
	// sendmany "fromaccount" {"address":amount,...} ( minconf "comment" ["address",...] )
	"sendmany": {false},
	// sendmanywithdupes "fromaccount" {"address":amount,...} ( minconf "comment" ["address",...] )
	"sendmanywithdupes": {false},
	// sendtoaddress "address" amount ( "comment" "comment_to" subtractfeefromamount )
	"sendtoaddress": {false},
	// sendtocontract "contractaddress" "data" (amount gaslimit gasprice senderaddress broadcast)
	"sendtocontract": {false},
	// DEPRECATED. setaccount "address" "account"
	// settxfee amount
	// signmessage "address" "message"
	"waitforlogs": {true},

	"help": {true},
}
