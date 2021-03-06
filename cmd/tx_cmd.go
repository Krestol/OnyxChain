/*
 * Copyright (C) 2019 The onyxchain Authors
 * This file is part of The onyxchain library.
 *
 * The onyxchain is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The onyxchain is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The onyxchain.  If not, see <http://www.gnu.org/licenses/>.
 */
package cmd

import (
	"encoding/hex"
	"fmt"
	cmdcom "github.com/OnyxPay/OnyxChain/cmd/common"
	"github.com/OnyxPay/OnyxChain/cmd/utils"
	"github.com/OnyxPay/OnyxChain/common"
	nutils "github.com/OnyxPay/OnyxChain/smartcontract/service/native/utils"
	"github.com/urfave/cli"
	"strconv"
	"strings"
)

var SendTxCommand = cli.Command{
	Name:        "sendtx",
	Usage:       "Send raw transaction to OnyxChain",
	Description: "Send raw transaction to OnyxChain.",
	ArgsUsage:   "<rawtx>",
	Action:      sendTx,
	Flags: []cli.Flag{
		utils.RPCPortFlag,
		utils.PrepareExecTransactionFlag,
	},
}

func sendTx(ctx *cli.Context) error {
	SetRpcPort(ctx)
	if ctx.NArg() < 1 {
		PrintErrorMsg("Missing raw tx argument.")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}
	rawTx := ctx.Args().First()

	isPre := ctx.IsSet(utils.GetFlagName(utils.PrepareExecTransactionFlag))
	if isPre {
		preResult, err := utils.PrepareSendRawTransaction(rawTx)
		if err != nil {
			return err
		}
		if preResult.State == 0 {
			return fmt.Errorf("prepare execute transaction failed. %v", preResult)
		}
		PrintInfoMsg("Prepare execute transaction success.")
		PrintInfoMsg("Gas limit:%d", preResult.Gas)
		PrintInfoMsg("Result:%v", preResult.Result)
		return nil
	}
	txHash, err := utils.SendRawTransactionData(rawTx)
	if err != nil {
		return err
	}
	PrintInfoMsg("Send transaction success.")
	PrintInfoMsg("  TxHash:%s", txHash)
	PrintInfoMsg("\nTip:")
	PrintInfoMsg("  Using './onyxchain info status %s' to query transaction status.", txHash)
	return nil
}

var TxCommond = cli.Command{
	Name:  "buildtx",
	Usage: "Build transaction",
	Subcommands: []cli.Command{
		TransferTxCommond,
		ApproveTxCommond,
		TransferFromTxCommond,
		WithdrawOXGTxCommond,
	},
	Description: "Build transaction",
}

var TransferTxCommond = cli.Command{
	Name:        "transfer",
	Usage:       "Build transfer transaction",
	Description: "Build transfer transaction.",
	Action:      transferTx,
	Flags: []cli.Flag{
		utils.WalletFileFlag,
		utils.TransactionGasPriceFlag,
		utils.TransactionGasLimitFlag,
		utils.TransactionPayerFlag,
		utils.TransactionAssetFlag,
		utils.TransactionFromFlag,
		utils.TransactionToFlag,
		utils.TransactionAmountFlag,
	},
}

var ApproveTxCommond = cli.Command{
	Name:        "approve",
	Usage:       "Build approve transaction",
	Description: "Build approve transaction.",
	Action:      approveTx,
	Flags: []cli.Flag{
		utils.WalletFileFlag,
		utils.TransactionGasPriceFlag,
		utils.TransactionGasLimitFlag,
		utils.TransactionPayerFlag,
		utils.ApproveAssetFlag,
		utils.ApproveAssetFromFlag,
		utils.ApproveAssetToFlag,
		utils.ApproveAmountFlag,
	},
}

var TransferFromTxCommond = cli.Command{
	Name:        "transferfrom",
	Usage:       "Build transfer from transaction",
	Description: "Build transfer from transaction.",
	Action:      transferFromTx,
	Flags: []cli.Flag{
		utils.WalletFileFlag,
		utils.TransactionGasPriceFlag,
		utils.TransactionGasLimitFlag,
		utils.ApproveAssetFlag,
		utils.TransactionPayerFlag,
		utils.TransferFromSenderFlag,
		utils.ApproveAssetFromFlag,
		utils.ApproveAssetToFlag,
		utils.TransferFromAmountFlag,
	},
}

var WithdrawOXGTxCommond = cli.Command{
	Action:      withdrawOXGTx,
	Name:        "withdrawoxg",
	Usage:       "Build Withdraw OXG transaction",
	Description: "Build Withdraw OXG transaction",
	ArgsUsage:   "<address|label|index>",
	Flags: []cli.Flag{
		utils.RPCPortFlag,
		utils.WalletFileFlag,
		utils.TransactionGasPriceFlag,
		utils.TransactionGasLimitFlag,
		utils.TransactionPayerFlag,
		utils.WithdrawOXGAmountFlag,
		utils.WithdrawOXGReceiveAccountFlag,
	},
}

func transferTx(ctx *cli.Context) error {
	if !ctx.IsSet(utils.GetFlagName(utils.TransactionToFlag)) ||
		!ctx.IsSet(utils.GetFlagName(utils.TransactionFromFlag)) ||
		!ctx.IsSet(utils.GetFlagName(utils.TransactionAmountFlag)) {
		PrintErrorMsg("Missing %s %s or %s argument.", utils.TransactionToFlag.Name, utils.TransactionFromFlag.Name, utils.TransactionAmountFlag.Name)
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	gasPrice := ctx.Uint64(utils.TransactionGasPriceFlag.Name)
	gasLimit := ctx.Uint64(utils.TransactionGasLimitFlag.Name)

	asset := ctx.String(utils.GetFlagName(utils.TransactionAssetFlag))
	if asset == "" {
		asset = utils.ASSET_ONX
	}
	from := ctx.String(utils.GetFlagName(utils.TransactionFromFlag))
	fromAddr, err := cmdcom.ParseAddress(from, ctx)
	if err != nil {
		return err
	}
	to := ctx.String(utils.GetFlagName(utils.TransactionToFlag))
	toAddr, err := cmdcom.ParseAddress(to, ctx)
	if err != nil {
		return err
	}

	var payer common.Address
	payerAddr := ctx.String(utils.GetFlagName(utils.TransactionPayerFlag))
	if payerAddr != "" {
		payerAddr, err = cmdcom.ParseAddress(payerAddr, ctx)
		if err != nil {
			return err
		}
	} else {
		payerAddr = fromAddr
	}

	payer, err = common.AddressFromBase58(payerAddr)
	if err != nil {
		return fmt.Errorf("invalid payer address:%s", err)
	}

	var amount uint64
	amountStr := ctx.String(utils.TransactionAmountFlag.Name)
	switch strings.ToLower(asset) {
	case "onyx":
		amount = utils.ParseOnx(amountStr)
		amountStr = utils.FormatOnx(amount)
	case "oxg":
		amount = utils.ParseOxg(amountStr)
		amountStr = utils.FormatOxg(amount)
	default:
		return fmt.Errorf("unsupport asset:%s", asset)
	}

	err = utils.CheckAssetAmount(asset, amount)
	if err != nil {
		return err
	}

	mutTx, err := utils.TransferTx(gasPrice, gasLimit, asset, fromAddr, toAddr, amount)
	if err != nil {
		return err
	}
	mutTx.Payer = payer

	tx, err := mutTx.IntoImmutable()
	if err != nil {
		return fmt.Errorf("IntoImmutable error:%s", err)
	}
	sink := common.ZeroCopySink{}
	err = tx.Serialization(&sink)
	if err != nil {
		return fmt.Errorf("tx serialization error:%s", err)
	}
	PrintInfoMsg("Transfer raw tx:")
	PrintInfoMsg(hex.EncodeToString(sink.Bytes()))
	return nil
}

func approveTx(ctx *cli.Context) error {
	asset := ctx.String(utils.GetFlagName(utils.ApproveAssetFlag))
	from := ctx.String(utils.GetFlagName(utils.ApproveAssetFromFlag))
	to := ctx.String(utils.GetFlagName(utils.ApproveAssetToFlag))
	amountStr := ctx.String(utils.GetFlagName(utils.ApproveAmountFlag))
	if asset == "" ||
		from == "" ||
		to == "" ||
		amountStr == "" {
		PrintErrorMsg("Missing %s %s %s or %s argument.", utils.ApproveAssetFlag.Name, utils.ApproveAssetFromFlag.Name, utils.ApproveAssetToFlag.Name, utils.ApproveAmountFlag.Name)
		cli.ShowSubcommandHelp(ctx)
		return nil
	}
	fromAddr, err := cmdcom.ParseAddress(from, ctx)
	if err != nil {
		return err
	}
	toAddr, err := cmdcom.ParseAddress(to, ctx)
	if err != nil {
		return err
	}

	var payer common.Address
	payerAddr := ctx.String(utils.GetFlagName(utils.TransactionPayerFlag))
	if payerAddr != "" {
		payerAddr, err = cmdcom.ParseAddress(payerAddr, ctx)
		if err != nil {
			return err
		}
	} else {
		payerAddr = fromAddr
	}

	payer, err = common.AddressFromBase58(payerAddr)
	if err != nil {
		return fmt.Errorf("invalid payer address:%s", err)
	}

	var amount uint64
	switch strings.ToLower(asset) {
	case "onyx":
		amount = utils.ParseOnx(amountStr)
		amountStr = utils.FormatOnx(amount)
	case "oxg":
		amount = utils.ParseOxg(amountStr)
		amountStr = utils.FormatOxg(amount)
	default:
		return fmt.Errorf("unsupport asset:%s", asset)
	}

	err = utils.CheckAssetAmount(asset, amount)
	if err != nil {
		return err
	}

	gasPrice := ctx.Uint64(utils.TransactionGasPriceFlag.Name)
	gasLimit := ctx.Uint64(utils.TransactionGasLimitFlag.Name)

	mutTx, err := utils.ApproveTx(gasPrice, gasLimit, asset, fromAddr, toAddr, amount)
	if err != nil {
		return err
	}
	mutTx.Payer = payer

	tx, err := mutTx.IntoImmutable()
	if err != nil {
		return fmt.Errorf("IntoImmutable error:%s", err)
	}
	sink := common.ZeroCopySink{}
	err = tx.Serialization(&sink)
	if err != nil {
		return fmt.Errorf("tx serialization error:%s", err)
	}
	PrintInfoMsg("Approve raw tx:")
	PrintInfoMsg(hex.EncodeToString(sink.Bytes()))
	return nil
}

func transferFromTx(ctx *cli.Context) error {
	asset := ctx.String(utils.GetFlagName(utils.ApproveAssetFlag))
	from := ctx.String(utils.GetFlagName(utils.ApproveAssetFromFlag))
	to := ctx.String(utils.GetFlagName(utils.ApproveAssetToFlag))
	amountStr := ctx.String(utils.GetFlagName(utils.TransferFromAmountFlag))
	if asset == "" ||
		from == "" ||
		to == "" ||
		amountStr == "" {
		PrintErrorMsg("Missing %s %s %s or %s argument.", utils.ApproveAssetFlag.Name, utils.ApproveAssetFromFlag.Name, utils.ApproveAssetToFlag.Name, utils.TransferFromAmountFlag.Name)
		cli.ShowSubcommandHelp(ctx)
		return nil
	}
	fromAddr, err := cmdcom.ParseAddress(from, ctx)
	if err != nil {
		return err
	}
	toAddr, err := cmdcom.ParseAddress(to, ctx)
	if err != nil {
		return err
	}

	var sendAddr string
	sender := ctx.String(utils.GetFlagName(utils.TransferFromSenderFlag))
	if sender == "" {
		sendAddr = toAddr
	} else {
		sendAddr, err = cmdcom.ParseAddress(sender, ctx)
		if err != nil {
			return err
		}
	}

	var payer common.Address
	payerAddr := ctx.String(utils.GetFlagName(utils.TransactionPayerFlag))
	if payerAddr != "" {
		payerAddr, err = cmdcom.ParseAddress(payerAddr, ctx)
		if err != nil {
			return err
		}
	} else {
		payerAddr = sendAddr
	}

	payer, err = common.AddressFromBase58(payerAddr)
	if err != nil {
		return fmt.Errorf("invalid payer address:%s", err)
	}

	var amount uint64
	switch strings.ToLower(asset) {
	case "onyx":
		amount = utils.ParseOnx(amountStr)
		amountStr = utils.FormatOnx(amount)
	case "oxg":
		amount = utils.ParseOxg(amountStr)
		amountStr = utils.FormatOxg(amount)
	default:
		return fmt.Errorf("unsupport asset:%s", asset)
	}

	err = utils.CheckAssetAmount(asset, amount)
	if err != nil {
		return err
	}

	gasPrice := ctx.Uint64(utils.TransactionGasPriceFlag.Name)
	gasLimit := ctx.Uint64(utils.TransactionGasLimitFlag.Name)

	mutTx, err := utils.TransferFromTx(gasPrice, gasLimit, asset, sendAddr, fromAddr, toAddr, amount)
	if err != nil {
		return err
	}
	mutTx.Payer = payer

	tx, err := mutTx.IntoImmutable()
	if err != nil {
		return fmt.Errorf("IntoImmutable error:%s", err)
	}
	sink := common.ZeroCopySink{}
	err = tx.Serialization(&sink)
	if err != nil {
		return fmt.Errorf("tx serialization error:%s", err)
	}
	PrintInfoMsg("TransferFrom raw tx:")
	PrintInfoMsg(hex.EncodeToString(sink.Bytes()))
	return nil
}

func withdrawOXGTx(ctx *cli.Context) error {
	SetRpcPort(ctx)
	if ctx.NArg() < 1 {
		PrintErrorMsg("Missing account argument.")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}
	addrArg := ctx.Args().First()
	accAddr, err := cmdcom.ParseAddress(addrArg, ctx)
	if err != nil {
		return err
	}

	fromAddr := nutils.OnxContractAddress.ToBase58()

	var amount uint64
	amountStr := ctx.String(utils.GetFlagName(utils.TransferFromAmountFlag))
	if amountStr == "" {
		balance, err := utils.GetAllowance("oxg", fromAddr, accAddr)
		if err != nil {
			return err
		}
		amount, err = strconv.ParseUint(balance, 10, 64)
		if err != nil {
			return err
		}
		if amount <= 0 {
			return fmt.Errorf("haven't unbound oxg")
		}
		amountStr = utils.FormatOxg(amount)
	} else {
		amount = utils.ParseOxg(amountStr)
		if amount <= 0 {
			return fmt.Errorf("haven't unbound oxg")
		}
		amountStr = utils.FormatOxg(amount)
	}

	var payer common.Address
	payerAddr := ctx.String(utils.GetFlagName(utils.TransactionPayerFlag))
	if payerAddr != "" {
		payerAddr, err = cmdcom.ParseAddress(payerAddr, ctx)
		if err != nil {
			return err
		}
	} else {
		payerAddr = accAddr
	}
	payer, err = common.AddressFromBase58(payerAddr)
	if err != nil {
		return fmt.Errorf("invalid payer address:%s", err)
	}

	var receiveAddr string
	receive := ctx.String(utils.GetFlagName(utils.WithdrawOXGReceiveAccountFlag))
	if receive == "" {
		receiveAddr = accAddr
	} else {
		receiveAddr, err = cmdcom.ParseAddress(receive, ctx)
		if err != nil {
			return err
		}
	}

	gasPrice := ctx.Uint64(utils.TransactionGasPriceFlag.Name)
	gasLimit := ctx.Uint64(utils.TransactionGasLimitFlag.Name)

	PrintInfoMsg("Withdraw account:%s", accAddr)
	PrintInfoMsg("Receive account:%s", receiveAddr)
	PrintInfoMsg("Withdraw OXG amount:%v", amount)
	mutTx, err := utils.TransferFromTx(gasPrice, gasLimit, "oxg", accAddr, fromAddr, receiveAddr, amount)
	if err != nil {
		return err
	}

	mutTx.Payer = payer
	tx, err := mutTx.IntoImmutable()
	if err != nil {
		return fmt.Errorf("IntoImmutable error:%s", err)
	}
	sink := common.ZeroCopySink{}
	err = tx.Serialization(&sink)
	if err != nil {
		return fmt.Errorf("tx serialization error:%s", err)
	}
	PrintInfoMsg("Withdraw raw tx:")
	PrintInfoMsg(hex.EncodeToString(sink.Bytes()))
	return nil
}
