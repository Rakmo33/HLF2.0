const {
  Gateway,
  Wallets,
  TxEventHandler,
  GatewayOptions,
  DefaultEventHandlerStrategies,
  TxEventHandlerFactory,
} = require("fabric-network");
const fs = require("fs");
const EventStrategies = require("fabric-network/lib/impl/event/defaulteventhandlerstrategies");
const path = require("path");
const log4js = require("log4js");
const logger = log4js.getLogger("BasicNetwork");
const util = require("util");

const helper = require("./helper");
const { blockListener, contractListener } = require("./Listeners");

const invokeTransaction = async (
  channelName,
  chaincodeName,
  fcn,
  args,
  username,
  org_name,
  transientData
) => {
  try {
    const ccp = await helper.getCCP(org_name);

    const walletPath = await helper.getWalletPath(org_name);
    const wallet = await Wallets.newFileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    let identity = await wallet.get(username);
    if (!identity) {
      console.log(
        `An identity for the user ${username} does not exist in the wallet, so registering user`
      );
      await helper.getRegisteredUser(username, org_name, true);
      identity = await wallet.get(username);
      console.log("Run the registerUser.js application before retrying");
      return;
    }

    const connectOptions = {
      wallet,
      identity: username,
      discovery: { enabled: true, asLocalhost: true },
      // eventHandlerOptions: EventStrategies.NONE
    };

    const gateway = new Gateway();
    await gateway.connect(ccp, connectOptions);

    const network = await gateway.getNetwork(channelName);
    const contract = network.getContract(chaincodeName);

    // Important: Please dont set listener here, I just showed how to set it. If we are doing here, it will set on every invoke call.
    // Instead create separate function and call it once server started, it will keep listening.
    // await contract.addContractListener(contractListener);
    // await network.addBlockListener(blockListener);

    // Multiple smartcontract in one chaincode
    let result;
    let message;

    switch (fcn) {
      case "CreateCar":
        result = await contract.submitTransaction(
          "SmartContract:" + fcn,
          args[0]
        );
        result = { txid: result.toString() };
        break;
      case "UpdateCarOwner":
        console.log("=============");
        result = await contract.submitTransaction(
          "SmartContract:" + fcn,
          args[0],
          args[1]
        );
        result = { txid: result.toString() };
        break;
      case "CreateDocument":
        result = await contract.submitTransaction(
          "DocumentContract:" + fcn,
          args[0]
        );
        console.log(result.toString());
        result = { txid: result.toString() };
        break;
      //custom functions by Omkar
      case "SetTransaction":
        result = await contract.submitTransaction(
          "SmartContract:" + fcn,
          args[0]
        );
        result = { txid: result.toString() };
        break;
      case "CreatePrivateCar":
        let carData = JSON.parse(transientData);
        let key = Object.keys(carData)[0];
        const transientDataBuffer = {};

        transientDataBuffer[key] = Buffer.from(JSON.stringify(carData.car));

        result = await contract
          .createTransaction("SmartContract:" + fcn)
          .setTransient(transientDataBuffer)
          .submit();
        // result = await contract.submitTransaction('SmartContract:'+fcn, [transientData]);

        console.log(`Result is =========================: ${result}`);
        result = { txid: result.toString() };
        break;
      case "SetPrivateTransaction":
        let transactionData = JSON.parse(transientData);
        let key2 = Object.keys(transactionData)[0];
        const transientDataBuffer2 = {};

        transientDataBuffer2[key2] = Buffer.from(JSON.stringify(transactionData.transaction));

        result = await contract
          .createTransaction("SmartContract:" + fcn)
          .setTransient(transientDataBuffer2)
          .submit(JSON.stringify(args));

        console.log(`Result is =========================: `, result);
        result = { txid: result.toString() };
        break;
      default:
        break;
    }

    await gateway.disconnect();

    // result = JSON.parse(result.toString());

    let response = {
      message: message,
      result,
    };

    return response;
  } catch (error) {
    console.log(`Getting error: ${error}`);
    return error.message;
  }
};

exports.invokeTransaction = invokeTransaction;
