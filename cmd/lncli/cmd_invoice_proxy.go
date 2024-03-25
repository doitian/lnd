package main

import (
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/urfave/cli"
)

var addInvoiceProxyCommand = cli.Command{
	Name:     "addinvoiceproxy",
	Category: "Invoices",
	Usage:    "Add a proxy invoice for the payment request.",
	Description: `
	Accepts a payment request and add a proxy invoice for it.

	When the proxy invoice is paid, fulfill the payment request
	to get the preimage.
	`,
	ArgsUsage: "--pay_req=R",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "pay_req",
			Usage: "a zpay32 encoded payment request to fulfill",
		},
	},
	Action: addInvoiceProxy,
}

func addInvoiceProxy(ctx *cli.Context) error {
	// Show command help if no arguments provided
	if (ctx.NArg() == 0 && ctx.NumFlags() == 0) || !ctx.IsSet("pay_req") {
		_ = cli.ShowCommandHelp(ctx, "addinvoiceproxy")
		return nil
	}

	req := &lnrpc.PayReqString{PayReq: stripPrefix(ctx.String("pay_req"))}

	ctxc := getContext()
	client, cleanUp := getClient(ctx)
	defer cleanUp()

	resp, err := client.AddInvoiceProxy(ctxc, req)
	if err != nil {
		return err
	}

	printRespJSON(resp)

	return nil
}
