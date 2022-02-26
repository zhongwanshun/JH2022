//SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./Time.sol";

library limitTime {
  //时间 把分钟转换为秒,
    function calculation(uint _minute) internal pure returns (uint) {
        uint second = _minute * 60;
        return second;
    }
}   

contract shareTlm {
 
     //合约主持人EOA
     address private host_tlm;
     //押金最低额度
     uint private _MinimumAmount = 10 ** 15;  //1 finney
    //超时归还赔偿的租金倍数
     uint overTimeMultiple = 2;
    //不归还物品赔偿的租金倍数
     uint noReturnMultiple = 10;
    //最大借物数
    uint allowBorrowNum = 3;
    //当前时间  但是每次查询时间前都需要重新调用getTime()

    Time time = new Time(); 
    uint public _time ;
     

    using limitTime for uint;
    uint internal _availableTime;  //下架时间

 //物品信息
 struct Goods_tlm {
   address owner;	//出借人EOA
   address borrower;//借用人EOA
   uint ethPledge;	//押金
   bool available;  //是否已上架
   bool isBorrow;   //是否已借出
   bool exist;      //物品是否存在
   uint availableTime; //物品所有者设置的下架时间
 }

//押金信息
struct yaJin{
    address addr; //地址
    uint money;   //押金金额
    uint id;    //此合约调用者交过押金的笔数
    bool exist; //是否符合最低借物押金
    uint borrowNum;   //此人手上借的物品数量，0则为没有借东西
}

//通过地址 映射 押金结构体
 mapping(address => yaJin) private yajinData_Tlm;

// //通过地址 映射 自己物品状况的结构体
//  mapping(address => myGoods) private myGoods_Tlm;

 //存储所有贴纸(分类)信息
 mapping(string => mapping(uint => Goods_tlm))  goodsData_tlm;
 
 //所有贴纸(分类)物品笔数
 mapping(string => uint)  goodsInx_tlm;
 
 //贴纸(分类)是否存在的记录
 mapping(string => bool)  goodsChk_tlm;

 //记录合约主持人
 constructor ()  {
   host_tlm = msg.sender;
 }
 	
 //只有主持人才可执行 
 modifier onlyHost() {
   require(msg.sender == host_tlm,
   "only host can do this");
    _;
 }

 //查询贴纸(分类)是否已经存在 与其有多少个物品
 function stickExistNum(string memory stickName_tlm) public view returns(bool,uint) {
   return (goodsChk_tlm[stickName_tlm],goodsInx_tlm[stickName_tlm]);
 }


 //添加一种贴纸(分类)
 function addSticker(string memory stickName_tlm) public onlyHost {
  //贴纸(分类)不存在，才可以添加
	 bool _exist;
   (_exist,) = stickExistNum(stickName_tlm);
   require(!_exist, 
		 "stick already exist");

   //设置可以使用此类贴纸
  goodsChk_tlm[stickName_tlm] = true; 
  
  //触发添加贴纸的事件
  emit addStickerEvnt("addSticker", stickName_tlm);
 }

 //添加贴纸(分类)事件
 event addStickerEvnt(string indexed eventType, string stickName_tlm);
 
 //添加物品
 function addGoods(string memory stickName_tlm, uint ethPledge, bool available,uint surplusTime) public returns(uint){
  //贴纸(分类)必须存在 
  bool  _exist;
   (_exist,) = stickExistNum(stickName_tlm);
    require(_exist, "stick not exist");
  if(available){  
    //预定下架时间 = 用户设置时间 + 现在的时间 block.timestamp
	 _availableTime = limitTime.calculation(surplusTime) + _time;	
  }
   //物品序号加1   
   goodsInx_tlm[stickName_tlm] +=1;
   uint inx = goodsInx_tlm[stickName_tlm];
   
   //新的物品信息
   Goods_tlm memory goodsTlm = Goods_tlm({
	 owner: msg.sender,		//出借人EOA
	 borrower: address(0),			//借用人EOA
	 ethPledge: ethPledge,	//租金
	 available: available,	//是否已上架
	 isBorrow: false,		//是否已借出 
	 exist: true,			//确认信息存在
   availableTime: _availableTime  //下架时间
   });
   
   //数据存储至映射结构
   goodsData_tlm[stickName_tlm][inx] = goodsTlm;
   
   //触发添加物品事件
   emit addGoodsEvnt("addGoods", stickName_tlm, inx);
   
   //返回数据索引
   return inx;
 }
 
 //添加物品事件
 event addGoodsEvnt(string indexed eventType, string stickName_tlm, uint inx);
 
 //判断物品是否存在,与如果存在它的所有者
 function isGoodExist(string memory stickName_tlm, uint inx) public view returns(bool,address){
   //贴纸(分类)必须存在 
   bool _exist;
   (_exist,) = stickExistNum(stickName_tlm);
   require(_exist, 
		 "stick not exist");
   
   return (goodsData_tlm[stickName_tlm][inx].exist,goodsData_tlm[stickName_tlm][inx].owner);   
 }
 
 //设置物品上下架,或者修改上架的时间
 function setGoodsStatus(string memory stickName_tlm, uint inx, bool available,uint surplusTime) public { 
   //物品必须存在 
    bool  _exist;
   (_exist,) = isGoodExist(stickName_tlm, inx);
   require(_exist, 
		 "goods not exist");
	
   //必须是出借人才可以改变状态	
   require(goodsData_tlm[stickName_tlm][inx].owner == msg.sender,
           "not goods owner");
		   
   //物品必须没被借出	
   require(!goodsData_tlm[stickName_tlm][inx].isBorrow,
           "goods already lend");
	
   //改变上下架状态	
   goodsData_tlm[stickName_tlm][inx].available = available;
   //改变物品上下架时间	
   if (available){ //如果是上架，则就要有一个下架时间。下架则为0
     goodsData_tlm[stickName_tlm][inx].availableTime = surplusTime + _time;
   }else{
    goodsData_tlm[stickName_tlm][inx].availableTime = 0;
   }

 }
 
 //查询物品是否上下架及剩余时间
 function isGoodsAvailable(string  memory stickName_tlm, uint inx ) public view returns(bool,uint) { 
   //物品必须存在 
  bool  _exist;
   (_exist,) = isGoodExist(stickName_tlm, inx);
   require(_exist, 
		 "goods not exist");

  uint sy_surplusTime; //剩余时间
  if (goodsData_tlm[stickName_tlm][inx].available) {
    if(_availableTime > _time){
      sy_surplusTime = _availableTime - _time;  //剩余时间 = 预定下架时间 - 现在时间
    }
  }
   //返回上下架状态	
   return (goodsData_tlm[stickName_tlm][inx].available,sy_surplusTime);
 }
 
 //查询物品借出状态
 function isGoodsLend(string memory stickName_tlm, uint inx) public view returns(bool,uint) { 
   //物品必须存在 
   bool  _exist;
   (_exist,) = isGoodExist(stickName_tlm, inx);
   require(_exist, 
		 "goods not exist");
     //剩余时间
     uint jc_surplusTime;
     //如果物品状态为上架，剩余时间为预定时间 - 现在
  if(goodsData_tlm[stickName_tlm][inx].isBorrow){ 
    if(goodsData_tlm[stickName_tlm][inx].availableTime > _time){
        jc_surplusTime = goodsData_tlm[stickName_tlm][inx].availableTime - _time;
    }
  }
   //返回借出状态与剩余时间
   return (goodsData_tlm[stickName_tlm][inx].isBorrow,jc_surplusTime);
 }
 
 // 借入物品
 function borrowGoods(string memory stickName_tlm, uint inx) public payable { 
   //物品必须存在 
   bool  _exist;
   (_exist,) = isGoodExist(stickName_tlm, inx);
   require(_exist, 
		 "goods not exist");
   
   //物品必须是可用状态
   require(goodsData_tlm[stickName_tlm][inx].available,
           "goods not available");
		   
   //物品必须没被借出	
   require(!goodsData_tlm[stickName_tlm][inx].isBorrow,
           "goods already lend");
	
   //租金必要符合设置
   require(goodsData_tlm[stickName_tlm][inx].ethPledge == msg.value,
           "eth pledge not match");
   //押金必要符合设置  
   require(yajinData_Tlm[msg.sender].exist,
           "yaJin less than Minimum Amount");
    //一个人不能一直借
    require(yajinData_Tlm[msg.sender].borrowNum < allowBorrowNum,
    "You have borrowed a lot of things. Please return them in time. " );
		   
   //设置借用人EOA
   goodsData_tlm[stickName_tlm][inx].borrower = msg.sender;
   
   //设置为已借出
   goodsData_tlm[stickName_tlm][inx].isBorrow = true;

   //借用人 借用物品数 +1
  yajinData_Tlm[msg.sender].borrowNum += 1;
  
  //借用人在合约中的钱要加上租金，即他在合约中的钱等于预先付的押金与接物品时的租金
  yajinData_Tlm[msg.sender].money += goodsData_tlm[stickName_tlm][inx].ethPledge;
  
   //触发借出事件
   emit borrowGoodsEvnt("borrowEvn", stickName_tlm, inx, msg.sender);
 }
 
 //物品借出事件
 event borrowGoodsEvnt(string indexed eventType, string stickName_tlm, uint inx, address borrower);
 
 //查询物品借出人
 function queryBorrower(string memory stickName_tlm, uint inx) public view returns(address) { 
   //物品必须存在 
   bool  _exist;
   (_exist,) = isGoodExist(stickName_tlm, inx);
   require(_exist, 
		 "goods not exist");
	
   //物品必须已被借出	
   require(goodsData_tlm[stickName_tlm][inx].isBorrow,
           "goods not lend");
		   
   //返回借出人
   return goodsData_tlm[stickName_tlm][inx].borrower;
 }
 
 //设置物品归还 ==> 未超时 || 已超时 ==> 都已归还
 function doGoodsReturn(string memory stickName_tlm, uint inx) public { 
   //物品必须存在 
  bool  _exist;
   (_exist,) = isGoodExist(stickName_tlm, inx);
   require(_exist, 
		 "goods not exist");
	
   //必须是出借人才可以改变状态	
   require(goodsData_tlm[stickName_tlm][inx].owner == msg.sender,
           "not goods owner");
		   
   //物品必须已被借出	
   require(goodsData_tlm[stickName_tlm][inx].isBorrow,
           "goods not lend");
  
  	// overTimeMultiple代表的是 超时归还赔偿的租金倍数，全局
   //将租金或补偿返还所有者 
   uint pledge = goodsData_tlm[stickName_tlm][inx].ethPledge ;
   if(goodsData_tlm[stickName_tlm][inx].availableTime < _time){  //如果超时走这里
          uint borrowerMaxMoney = yajinData_Tlm[goodsData_tlm[stickName_tlm][inx].borrower].money;  //该借物人全部的押金

        if(pledge * overTimeMultiple > borrowerMaxMoney){
          pledge = borrowerMaxMoney;
        }else{
           pledge *= overTimeMultiple;    // 原租金*超时多扣倍数
        }
   } 
   
  payable(goodsData_tlm[stickName_tlm][inx].owner).transfer(pledge);

  //借用人 借用物品数 -1  goodsData_tlm[stickName_tlm][inx].borrower ==>物品借用者
  yajinData_Tlm[goodsData_tlm[stickName_tlm][inx].borrower].borrowNum -= 1;
  // 押金 -pledge
  yajinData_Tlm[goodsData_tlm[stickName_tlm][inx].borrower].money -= pledge;

   //触发归还事件
   emit returnGoodsEvnt("returnEvn", stickName_tlm, inx, goodsData_tlm[stickName_tlm][inx].borrower);
   
   //设置借用人EOA
   goodsData_tlm[stickName_tlm][inx].borrower = address(0);
   
   //设置为未借出
   goodsData_tlm[stickName_tlm][inx].isBorrow = false;
 }
  //物品归还事件
 event returnGoodsEvnt(string indexed eventType, string stickName_tlm, uint inx, address borrower);
 

 //设置物品未归还 ==>  已超时 ==> 未归还。（就当他还了，东西送给他，获取数倍押金）
 function noGoodsReturn(string memory stickName_tlm, uint inx) public { 
   //物品必须存在 
  bool  _exist;
   (_exist,) = isGoodExist(stickName_tlm, inx);
   require(_exist, 
		 "goods not exist");
	
   //必须是出借人才可以改变状态	
   require(goodsData_tlm[stickName_tlm][inx].owner == msg.sender,
           "not goods owner");
		   
   //物品必须已被借出	
   require(goodsData_tlm[stickName_tlm][inx].isBorrow,
           "goods not lend");

  // 物品归还必须已超时
   require(goodsData_tlm[stickName_tlm][inx].availableTime < _time,
           "goods must return overtime");
  
  	// overTimeMultiple代表的是 超时归还赔偿的租金倍数，全局
   //补偿返还所有者，而且借用人原本付的租金也别想赖账，所以应该一起给过去，即共 ethPledge*(noReturnMultiple +1 )
  uint pledge = goodsData_tlm[stickName_tlm][inx].ethPledge *= (noReturnMultiple + 1 );

  //物品借用者的全部押金
  uint borrowerMaxMoney = yajinData_Tlm[goodsData_tlm[stickName_tlm][inx].borrower].money;
   
   if (pledge > borrowerMaxMoney){ 
        pledge = borrowerMaxMoney;
   }
  payable(goodsData_tlm[stickName_tlm][inx].owner).transfer(pledge);
  
  //借用人 借用物品数 -1 ,押金 -pledge。
  yajinData_Tlm[goodsData_tlm[stickName_tlm][inx].borrower].borrowNum -= 1;
  yajinData_Tlm[goodsData_tlm[stickName_tlm][inx].borrower].money -= pledge;
  

   //触发物品所有者更换事件，可以记录他的不诚信
   emit noReturnGoodsEvnt("noreturnEvn", stickName_tlm, inx, goodsData_tlm[stickName_tlm][inx].borrower);
     //更改这个物品信息
    Goods_tlm memory newGoodsTlm = Goods_tlm({
          owner:goodsData_tlm[stickName_tlm][inx].borrower,		//出借人EOA,先从之前的物品信息中拿到借物人地址
          borrower: address(0),			//借用人EOA
          ethPledge: 0,	//押金
          available: false,	//是否已上架
          isBorrow: false,		//是否已借出 
          exist: true,			//确认信息存在
          availableTime: 0  //下架时间
   });
    //删除这个物品的原有信息  delete mapping(string => mapping(uint => Goods))
    delete goodsData_tlm[stickName_tlm][inx];
    goodsData_tlm[stickName_tlm][inx] = newGoodsTlm;
 }
//声明物品所有者更换事件
event noReturnGoodsEvnt(string indexed eventType,string stickName_tlm,uint inx, address borrower);


  //查询合约余额
  function queryBalance() public view returns (uint) {
      require(msg.sender == host_tlm,"You are not host");
	    return address(this).balance;
  }
  
//声明付押金事件
event fyj(string indexed a, address add);

//充值押金  
function deposit() public payable{
    require( msg.value >= 10 ** 15,"Cannot lost 0.001 ether,finney"); //  每次付的押金大于等于 10 ** 15 Wei
    
    address addr = msg.sender;

    yajinData_Tlm[addr].addr = addr;
    yajinData_Tlm[addr].money += msg.value;
    yajinData_Tlm[addr].id++;

    //如果给的钱不够押金的最低额度 
    if (msg.value >= _MinimumAmount){   
      yajinData_Tlm[addr].exist = true;
    } else{
      yajinData_Tlm[addr].exist = false;
    }   
  //触发付押金事件
    emit fyj("add",msg.sender);
}

//查看当前合约调用者的押金
function chaYaJin(address addr) public view returns(uint, address, uint,uint, bool){
     return (yajinData_Tlm[addr].id,yajinData_Tlm[addr].addr,yajinData_Tlm[addr].money,yajinData_Tlm[addr].borrowNum,yajinData_Tlm[addr].exist);
}

//将押金返还借用人（提现）
 function withdraw(uint _amount) public {   //你想提现的金额
  
  require(yajinData_Tlm[ msg.sender].borrowNum == 0,
      "Please return the goods first");  //请先归还物品,否则锁定账户押金
      
  require(yajinData_Tlm[msg.sender].money >= _amount,
      "Your deposit is less than the money you want to withdraw"
    ); //   <=合约里的钱

   //将租金返还所有者
  uint pledge = yajinData_Tlm[msg.sender].money;
  payable(msg.sender).transfer(_amount);
  
   //触发退押金事件
   emit returnMoneyEvnt("returnMoney",msg.sender,_amount,pledge);
   
   //押金还剩多少
   yajinData_Tlm[msg.sender].money = pledge - _amount; 
   
   //如果押金小于最低押金，则转为false
   if (yajinData_Tlm[msg.sender].money < _MinimumAmount){
        yajinData_Tlm[msg.sender].exist = false;
   }
 }
 
 //退押金事件
 event returnMoneyEvnt(string indexed eventType,address addr,uint amonut , uint pledge);
 
 //修改时间
      function getTime() public {
            _time = time.callTime();
        }
}