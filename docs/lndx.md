# Cross-Chain PoC Demo

This is a PoC for the cross-chain hub.

## Protocol

- Carol creates an invoice and sends it to Alice.
- Alice creates a hodl invoice with the same payment hash and sends the hodl invoice to Carol.
- Carol shows the hodl invoice to Bob.
- Bob pays the hodl invoice.
- Alice accepts the htlc.
- Alice pays the original invoice created by Carol.
- Carol accepts the htlc and sends back the preimage to settle the payment.
- Alice gets the preimage from Carol and uses it to settle the hodl invoice.

## Demo

Download [Polar](https://github.com/jamaljsr/polar) to setup a lighting cluster for demo.

### Build Docker Image

Clone this repository and checkout the branch `lndx`.

```
git clone git@github.com:doitian/lnd.git
cd lnd
git checkout lndx
```

Build a docker image tagged `lndx`

```
docker build -t lndx -f polar/Dockerfile .
```

### Add Custom Image

Add a custom image in Polar via the menu item "Manage Images" in the top right corner.

- Click the button "Add a Custom Node"
- Name: LNDX
- Implementation: LND
- Docker Image: lndx:latest
- Command: Keep it unchanged

### Create a Network

Create a network with 1 LNDX node, 2 LND nodes, and 1 Bitcoin Core node.

Polar will name the LNDX node as alice, and the 2 LND nodes as bob and carol.

Start the network.

### Setup Channels

- Deposit funds to Alice and Bob.
- Create a channel from Bob to Alice.
- Create a channel from Alice to Carol.

### Complete a Payment

1. Click the node Carol and create an invoice from the Actions tab. Copy the
payment request string.
2. Click the node Alice and launch a terminal from the Actions tab. In the
terminal window, create a hodl invice using the following command:

    ```
    lncli addinvoiceproxy --pay_req PAY_REQ_COPIED_IN_STEP_1
    ```

    Copy the value of the field `payment_request` in the JSON response.

3. Click the node Bob and select "Pay Invoice" from the Actions tab. Paste the
   payment request copied in step 2.

Check that both channels have moved 50,000 sats to the destination.
