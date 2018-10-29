# Hyperledger Fabric Lab 01

In this lab we have developed a basic chaincode for the Hyperledger Fabric Blockchain network using Go programming language, and executed it on the `chaincode-docker-devmode`.

## Preparation
Download and install following pre-requisites:
* Docker Community Edition (for Mac, Linux and Windows Professional Editions)
* Docker Toolbox (for Windows Home Edition)
* docker-compose (separate installation required only for Linux, otherwise included with Docker)
* Node.js (version 8.x.x only)
* Golang 
* Python
* cURL
* Git

Download Hyperledger Fabric Docker images, samples and binaries:
`curl -sSL http://bit.ly/2ysbOFE | bash -s 1.3.0`

## Setup
You will see a folder called `fabric-samples` downloaded by above cURL command

* Open a terminal or command prompt, and `cd` to the `fabric-samples/chaincode`
* Clone this repository into the chaincode folder
* Now cd to `fabric-samples/chaincode-docker-devmode`
* To create and start Docker containers, run:
`docker-compose -f docker-compose-simple.yaml up`
* Open two terminal or command prompt windows.
* In the first window, open `bash` terminal of `chaincode` docker container by executing following command:
`docker exec -it chaincode bash`
* In the second window, open `bash` terminal of `cli` docker container by executing following command:
`docker exec -it cli bash`

## Execution
In this step, we will build the `hyperledger_lab_01` chaincode and execute it in the `devmode`

**Chaincode Container**
* In the `chaincode` container's terminal, cd to `hyperledger_lab_01` folder, and run `go build`
* Now launch the chaincode with the following command:
`CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=hyperledger_lab_01:0 ./hyperledger_lab_01`

**CLI Container**
* In the `cli` container's terminal, run the following command to install `hyperledger_lab_01` chaincode
`peer chaincode install -p chaincodedev/chaincode/hyperledger_lab_01 -n hyperledger_lab_01 -v 0`
* Run the following command to instantiate the chaincode
`peer chaincode instantiate -n hyperledger_lab_01 -v 0 -c '{"Args":[]}' -C myc`
* Now you can use following set of commands to perform operations in the chaincode

```
peer chaincode invoke -n hyperledger_lab_01 -C myc -c '{"Args":["set", "key1", "Rahul Bairathi"]}' -C myc

peer chaincode query -n hyperledger_lab_01 -C myc -c '{"Args":["get", "key1"]}' -C myc

peer chaincode query -n hyperledger_lab_01 -C myc -c '{"Args":["history", "key1"]}' -C myc

peer chaincode invoke -n hyperledger_lab_01 -C myc -c '{"Args":["delete", "key1"]}' -C myc
```

## Updating and re-deploying chaincode
If you have made any change in the chaincode, then: 
* terminate the chaincode execution in the `chaincode` container's terminal by `ctrl+c`. 
* build the chaincode again with `go build`
* run the chaincode with `./hyperledger_lab01`

You don't have to re-instantiate and update the chaincode in CLI container.
