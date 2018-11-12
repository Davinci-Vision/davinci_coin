var DACToken = artifacts.require("DACToken");
var DACTimelock = artifacts.require("DACTimelock");

var DACTokenVesting = artifacts.require("DACTokenVesting");

module.exports = async function(deployer) {

  // deployer.deploy(DACToken).then(()=>{
  //   var instance = DACToken;
  //   console.log('DACToken is ', instance.address);
  //   return deployer.deploy(DACTimelock, instance.address);
  // });

  deployer.deploy(DACToken);
  // deployer.deploy(DACTokenVesting, account_three, Math.floor(new Date().getTime()/1000), 0, 1000, 100, true);

 };
