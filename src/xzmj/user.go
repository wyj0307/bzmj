package xzmj

import (
	"common"
	"container/heap"
	"container/list"
	"errorValue"
	"logger"
	"rpc"
)

type MjLogicUser struct {
	mUid                     uint64
	mUserState               rpc.MjUserState
	mUserAction              uint32
	mSeatId                  uint32
	mLastAddCard             uint32
	mLastDealCard            uint32
	mLastPengCard            uint32
	mLastMingGangCard        uint32
	mLastAnGangCard          uint32
	mLastBuGangCard          uint32
	mRejectSuit              uint32
	mHuCard                  uint32
	mIsOffLine               bool
	mIsLeaveTable            bool
	mStartMatchCurrencyValue int64
	mRoomId                  uint32
	mOperateRecord           *list.List
	mSettlementInfo          *list.List
	mLastGameInfo            *MjLastUserGameInfo

	mOperateMask   *list.List
	mDealCards     IntHeap //玩家的出牌
	mHandCards     IntHeap //玩家的手牌
	mMingGangCards IntHeap //玩家的明杠
	mBuGangCards   IntHeap //玩家的补杠
	mAnGangCards   IntHeap //玩家的补杠
	mPengCards     IntHeap //玩家的碰牌
	mHuanCards     IntHeap //玩家的换牌
	mHuanInCards   IntHeap //玩家的换进的牌
	mInitHandCards IntHeap //初始手牌
}

var XzmjUser *MjLogicUser

func NewUser() *MjLogicUser {
	llXzmjUser := &MjLogicUser{}
	XzmjUser = llXzmjUser
	return XzmjUser
}

func (self *MjLogicUser) Init(seatId uint32, roomId uint32, uid uint64) {
	logger.Info("MjLogicUser:Init<ENTER>, seatId:%d,roomId:%d, uid:%d", seatId, roomId, uid)
	self.mUid = uid
	self.mUserState = rpc.MjUserState_Mj_User_State_Init
	self.mUserAction = Mj_User_Action_Init

	self.mSeatId = seatId

	self.mLastAddCard = MJ_BACK
	self.mLastDealCard = MJ_BACK
	self.mLastPengCard = MJ_BACK
	self.mLastMingGangCard = MJ_BACK
	self.mLastAnGangCard = MJ_BACK
	self.mLastBuGangCard = MJ_BACK
	self.mRejectSuit = MJ_BACK

	self.mHuCard = MJ_BACK
	self.mIsOffLine = false
	self.mIsLeaveTable = false
	self.mStartMatchCurrencyValue = int64(0)

	self.mRoomId = roomId
	self.mOperateMask.Init()
	heap.Init(&(self.mDealCards))
	heap.Init(&(self.mHandCards))
	heap.Init(&(self.mMingGangCards))
	heap.Init(&(self.mAnGangCards))
	heap.Init(&(self.mBuGangCards))
	heap.Init(&(self.mPengCards))
	self.mOperateRecord = list.New()
	self.mOperateRecord.Init()
	self.mSettlementInfo = list.New()
	self.mSettlementInfo.Init()
	self.mLastGameInfo = NewMjLastUserGameInfo()
	logger.Info("MjLogicUser:Init<LEAVE>, seatId:%d,roomId:%d, uid:%d", seatId, roomId, uid)
	return
}

func (self *MjLogicUser) ClearMatchInfo(IsAllMatchOver bool) {
	return
}

func (self *MjLogicUser) SetOperateMask(operateMask rpc.MJ_OPERATE_MASK) {
	self.ClearList(self.mOperateMask)
	self.mOperateMask.PushBack(operateMask)
	return
}

func (self *MjLogicUser) SetOperateMaskList(operateMasks *list.List) {
	self.ClearList(self.mOperateMask)
	for e := operateMasks.Front(); e != nil; e = e.Next() {
		self.mOperateMask.PushBack(e.Value)
	}
	return
}

func (self *MjLogicUser) ClearList(p *list.List) {
	e := p.Front()
	for e != nil {
		lTmp := e
		e = e.Next()
		p.Remove(lTmp)
	}
}

func (self *MjLogicUser) FindValueInList(l *list.List, element interface{}) *list.Element {
	var e *list.Element = nil
	for e = l.Front(); e != nil; e = e.Next() {
		if e.Value == element {
			return e
		}
	}
	return e
}

func (self *MjLogicUser) SetRejectSuit(rejectSuit uint32) {
	if rejectSuit >= MJ_WANG_SUIT && rejectSuit <= MJ_TIAO_SUIT {
		self.mRejectSuit = rejectSuit
	}
}

func (self *MjLogicUser) setUserLeaveTable() uint32 {
	self.mIsLeaveTable = true
	return errorValue.ERET_OK
}

func (self *MjLogicUser) setUserOffLine() uint32 {
	self.mIsOffLine = true
	return errorValue.ERET_OK
}

func (self *MjLogicUser) SetUserState(status rpc.MjUserState) uint32 {
	if self.mUserState == status {
		return errorValue.ERET_OK
	}
	self.mUserState = status
	return errorValue.ERET_OK
}

func (self *MjLogicUser) SetHuanSanZhang(tagMj uint32, tagMj1 uint32, tagMj2 uint32) uint32 {
	self.mHuanCards.Clear()
	heap.Push(&(self.mHuanCards), tagMj)
	heap.Push(&(self.mHuanCards), tagMj1)
	heap.Push(&(self.mHuanCards), tagMj2)
	return errorValue.ERET_OK
}

func (self *MjLogicUser) clearOperateMask() uint32 {
	self.ClearList(self.mOperateMask)
	self.mOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_NONE)
	return errorValue.ERET_OK
}

func (self *MjLogicUser) ReconnectEnterGame() uint32 {
	self.mIsOffLine = false
	return errorValue.ERET_OK
}

func (self *MjLogicUser) GetOperateMask() *list.List {
	return self.mOperateMask
}

func (self *MjLogicUser) handleDianPao(seatId uint32, tagMj uint32) uint32 {
	if seatId == self.mSeatId {
		return errorValue.ERET_HU
	}
	return errorValue.ERET_OK
}

func (self *MjLogicUser) handleZiMo(tagMj uint32) uint32 {
	lIt := self.mHandCards.Find(tagMj)
	if lIt == self.mHandCards.Len() {
		return errorValue.ERET_OK
	}
	heap.Remove(&(self.mHandCards), lIt)
	return errorValue.ERET_OK
}

func (self *MjLogicUser) CalculateOperateMask(MjCountInPool uint32) *list.List {

	var lOperateMask *list.List
	lOperateMask.Init()
	//heap.Init(&lOperateMask)

	lHandCardSize := self.mHandCards.Len()
	lMingGangCardsSize := self.mMingGangCards.Len()
	lAnGangCardsSize := self.mAnGangCards.Len()
	lBuGangCardsSize := self.mBuGangCards.Len()
	lPengCardsSize := self.mPengCards.Len()

	lMjToalCount := lHandCardSize + lMingGangCardsSize*3 + lAnGangCardsSize*3 + lBuGangCardsSize*3 + lPengCardsSize*3
	if lMjToalCount == HAND_MJ_COUNT && self.mUserState != rpc.MjUserState_Mj_User_State_Hu && self.mUserState != rpc.MjUserState_Mj_User_State_GiveUp {

		lGang := false
		lBuGang := false
		lHu := false

		for lIts := range self.mHandCards {
			var lIt uint32 = uint32(lIts)
			var lResult uint32 = self.GetTagMjCountInHand(&lIt)
			if lResult == GANG_MJ_COUNT && lGang == false && (MJ_GET_SUIT(lIt) != self.mRejectSuit) && (self.mUserAction == Mj_User_Action_MoPai || self.mUserAction == Mj_User_Action_Init) && MjCountInPool > 0 {
				lGang = true //有暗杠
				lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_AN_GANG)
			}
			if (lBuGang == false) && (self.mUserAction == Mj_User_Action_MoPai) && (self.mPengCards.Len() > 0) && self.IsCanBuGangTagMj(lIt) && (MjCountInPool > 0) {
				lBuGang = true
				lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_BU_GANG) //补杠
			}
		}

		if lBuGang == false && self.IsCanBuGangTagMj(self.mLastAddCard) && MJ_GET_SUIT(self.mLastAddCard) != self.mRejectSuit && self.mUserAction == Mj_User_Action_MoPai && MjCountInPool > 0 {
			lBuGang = true
			lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_BU_GANG) //补杠
		}
		if self.IsHaveRejectSuit() == false && (self.mUserAction == Mj_User_Action_MoPai || self.mUserAction == Mj_User_Action_Init) {
			if self.IsCanHu(&(self.mHandCards), uint32(self.mHandCards.Len()), uint32(0), false) == true {
				lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_HU)
				lHu = true
			}
		}

		if lGang == false && lBuGang == false && lHu == false {
			lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_DEAL) //出牌
		} else {
			lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_GUO)
		}

	} else {
		lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_NONE)
	}
	return lOperateMask
}

func (self *MjLogicUser) CalculateOperateMaskForDealCard(seatId uint32, tagMj uint32, MjCountInPool uint32, RoomConfig *common.RoomConfig) *list.List {
	var lOperateMask *list.List
	//TODO
	lGuo := false
	lPeng := false
	lGang := false
	lHu := false

	self.ClearList(lOperateMask)
	if seatId < 0 || seatId > SIT_NUM_MAX || tagMj <= MJ_BACK || tagMj > MJ_9_TIAO || seatId == self.mSeatId || self.mUserState == rpc.MjUserState_Mj_User_State_Hu || self.mUserState == rpc.MjUserState_Mj_User_State_GiveUp || MjCountInPool == 0 {
		lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_NONE)
		return lOperateMask
	}

	if self.IsCanPengTagMj(tagMj) == true {
		lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_PENG) //碰牌
		if lGuo == false {
			lGuo = true
			lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_GUO)
		}
		lPeng = true
	}
	if self.IsCanMingGangTagMj(tagMj) == true && MjCountInPool > 0 {
		lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_MING_GANG) //明杠
		if lGuo == false {
			lGuo = true
			lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_GUO)
		}
		lGang = true
	}

	if self.IsHaveRejectSuit() == false {
		if self.IsCanHu(&(self.mHandCards), uint32(self.mHandCards.Len()), uint32(0), false) == true {
			lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_HU) //胡牌
			if lGuo == false {
				lGuo = true
				lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_GUO)
			}
			lHu = true
		}
	}
	if lPeng == false && lGang == false && lHu == false {
		lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_NONE) //空
	}

	return lOperateMask
}

func (self *MjLogicUser) CalculateOperateMaskForBuGang(seatId uint32, tagMj uint32, MjCountInPool uint32) *list.List {
	var lOperateMask *list.List
	self.ClearList(lOperateMask)
	if seatId < 0 || seatId > SIT_NUM_MAX || tagMj <= MJ_BACK || tagMj > MJ_9_TIAO || seatId == self.mSeatId || self.mUserState == rpc.MjUserState_Mj_User_State_Hu || self.mUserState == rpc.MjUserState_Mj_User_State_GiveUp || MjCountInPool == 0 {
		lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_NONE)
		return lOperateMask
	}

	if self.IsHaveRejectSuit() == false {
		if self.IsCanHu(&(self.mHandCards), uint32(self.mHandCards.Len()), 0, false) == true {
			lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_HU)  //胡牌(补杠时抢杠)
			lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_GUO) //胡牌(补杠时附带过权限)
		} else {
			lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_NONE)
		}
	} else {
		lOperateMask.PushBack(rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_NONE)
	}

	return lOperateMask
}

/*
func (self *MjLogicUser) CopyMjCards(sourceMjCard *IntHeap, destMjCards *IntHeap) {
	copy(*destMjCards, *sourceMjCard)
}

func (self *MjLogicUser) SortMjCards(sourceMjCard *IntHeap) {
	return
}
*/

func (self *MjLogicUser) TestHu7Dui(pPai *IntHeap, mHandCardSize uint32) bool {
	if mHandCardSize != HAND_MJ_COUNT {
		return false
	}
	var i uint32
	for i = 0; i < 7; i++ {
		if (*pPai)[2*i] != (*pPai)[2*i+1] {
			return false
		}
	}
	return true
}

func (self *MjLogicUser) TestNormalCanHu(mHandCardTmp *IntHeap, mHandCardSize uint32, nProcessedCount uint32, bFoundJiang bool) bool {
	var lMjIndex0 int = -1
	var lMjIndex1 int = -1
	var lMjIndex2 int = -1
	var lMjCard uint32 = MJ_BACK

	if nProcessedCount == mHandCardSize {
		return true
	}

	lMjIndex0 = self.GetNextValidMJ(mHandCardTmp, mHandCardSize, 0, &lMjCard)
	if lMjIndex0 < 0 || lMjCard == MJ_BACK {
		return false
	}

	lTemp := false

	//	首先检测坎
	lMjIndex1 = self.FindMjIndex(mHandCardTmp, mHandCardSize, uint32(lMjIndex0)+1, &lMjCard)
	if lMjIndex1 >= 0 {
		lMjIndex2 = self.FindMjIndex(mHandCardTmp, mHandCardSize, uint32(lMjIndex1)+1, &lMjCard)
		if lMjIndex2 >= 0 {
			// 暂时删除已经选出的牌
			(*mHandCardTmp)[uint32(lMjIndex0)] = MJ_BACK
			(*mHandCardTmp)[uint32(lMjIndex1)] = MJ_BACK
			(*mHandCardTmp)[uint32(lMjIndex2)] = MJ_BACK
			// 递归调用测试剩余部分情况
			lTemp = self.TestNormalCanHu(mHandCardTmp, mHandCardSize, nProcessedCount+PENG_MJ_COUNT, bFoundJiang)

			// 重新设置牌
			(*mHandCardTmp)[uint32(lMjIndex0)] = lMjCard
			(*mHandCardTmp)[uint32(lMjIndex1)] = lMjCard
			(*mHandCardTmp)[uint32(lMjIndex2)] = lMjCard

			if lTemp == true {
				return true
			}
		}
	}
	if bFoundJiang == false {
		//找对子
		lMjIndex1 = self.FindMjIndex(mHandCardTmp, mHandCardSize, uint32(lMjIndex0)+1, &lMjCard)
		if lMjIndex1 >= 0 {
			//	暂时删除已经选出的牌
			(*mHandCardTmp)[uint32(lMjIndex0)] = MJ_BACK
			(*mHandCardTmp)[uint32(lMjIndex1)] = MJ_BACK

			//	递归调用测试剩余部分情况
			lTemp = self.TestNormalCanHu(mHandCardTmp, mHandCardSize, nProcessedCount+DUIZI_MJ_COUNT, true)

			(*mHandCardTmp)[uint32(lMjIndex0)] = lMjCard
			(*mHandCardTmp)[uint32(lMjIndex1)] = lMjCard

			if lTemp == true {
				return true
			}
		}
	}

	//	寻找顺子
	var lTmCard1 uint32 = lMjCard + 1
	lMjIndex1 = self.FindMjIndex(mHandCardTmp, mHandCardSize, uint32(lMjIndex0)+1, &(lTmCard1))
	if lMjIndex1 >= 0 {
		var lTmCard2 uint32 = lMjCard + 2
		lMjIndex2 = self.FindMjIndex(mHandCardTmp, mHandCardSize, uint32(lMjIndex1)+1, &(lTmCard2))
		if lMjIndex2 >= 0 {
			// 暂时删除已经选出的牌
			(*mHandCardTmp)[uint32(lMjIndex0)] = MJ_BACK
			(*mHandCardTmp)[uint32(lMjIndex1)] = MJ_BACK
			(*mHandCardTmp)[uint32(lMjIndex2)] = MJ_BACK

			//	递归调用测试剩余部分情况
			lTemp = self.TestNormalCanHu(mHandCardTmp, mHandCardSize, nProcessedCount+HUA_MJ_COUNT, bFoundJiang)
			(*mHandCardTmp)[uint32(lMjIndex0)] = lMjCard
			(*mHandCardTmp)[uint32(lMjIndex1)] = lMjCard + 1
			(*mHandCardTmp)[uint32(lMjIndex2)] = lMjCard + 2

			if lTemp == true {
				return true
			}
		}
	}
	return false
}

func (self *MjLogicUser) FindRejectSuitMj() int {
	var lResult int = -1
	for lMj := range self.mHandCards {
		if MJ_GET_SUIT(uint32(lMj)) == self.mRejectSuit {
			lResult = int(lMj)
		}
	}

	return lResult
}

func (self *MjLogicUser) FindMjIndex(mjCard *IntHeap, nCount uint32, nStartIndex uint32, card *uint32) int {
	if MJ_BACK == *card || nStartIndex < 0 || nStartIndex >= nCount {
		return -1
	}

	var i uint32 = nStartIndex
	for i = nStartIndex; i < nCount; i++ {
		if *card == (*mjCard)[i] {
			return int(i)
		}
	}

	return -1
}

//出牌
func (self *MjLogicUser) DealCard(tagMj uint32) uint32 {
	for i := 0; i < (self.mHandCards).Len(); i++ {
		if (self.mHandCards)[i] == tagMj {
			//self.mHandCards.Remove(i)
			heap.Remove(&(self.mHandCards), i)
			heap.Push(&(self.mDealCards), tagMj)
			self.mLastDealCard = tagMj
			self.mUserAction = Mj_User_Action_DealPai
			return errorValue.ERET_OK
		}
	}
	return errorValue.ERET_DEAL
}

//摸牌
func (self *MjLogicUser) AddCard(tagMj uint32) uint32 {
	heap.Push(&self.mHandCards, tagMj)
	self.mLastAddCard = tagMj
	self.mUserAction = Mj_User_Action_MoPai
	return errorValue.ERET_OK
}

//碰牌
//func (self *MjLogicUser) PengCard(tagMj uint32, roomCfg *RoomConfig) uint32 {
func (self *MjLogicUser) PengCard(tagMj uint32) uint32 {
	//self.mLastGuoDetailInfo.Init(); //操作过手，清楚过操作信息。
	for i := 0; i < 2; i++ {
		lIndex := (self.mHandCards).Find(tagMj)
		if int(lIndex) == (self.mHandCards).Len() {
			return errorValue.ERET_DEAL
		}
		//self.mHandCards.Remove(i)
		heap.Remove(&(self.mHandCards), i)
	}
	heap.Push(&(self.mPengCards), tagMj)
	self.mLastPengCard = tagMj
	self.mUserAction = Mj_User_Action_PengPai

	/*
			//添加过手胡信息
			uint32 multiple = 0;
		    HuCardAnalysis analyze;
			analyze.initAllCard(mHandCards, mPengCards, mMingGangCards, mAnGangCards, mBuGangCards);
		    auto huResult = analyze.Analyze();
		    if(huResult.Type != FINAL_CARD_TYPE_NONE)
		    {
				FinalCardResultByRoomConfig(roomCfg, &huResult);
				huResult.Type &= ~FINAL_CARD_TYPE_JIN_GOU_DIAO;
		        multiple = GetMultipleByCardType(huResult.Type) * GetMultipleByGenCount(huResult.GenCount);
				mLastGuoDetailInfo.mGuoHu = true;
				mLastGuoDetailInfo.mGuoHuFan = multiple;
				mLastGuoDetailInfo.mTagMj = tagMj;

		    }
	*/
	return errorValue.ERET_OK
}

//明杠
func (self *MjLogicUser) MingGangCard(seatId uint32, tagMj uint32) uint32 {
	//self.mLastGuoDetailInfo.Init(); //操作过手，清楚过操作信息。
	var i uint32 = 0
	for ; i < MAX_PLAYER_NUM_PER_TABLE-1; i++ {
		Postion := self.mHandCards.Find(tagMj)
		if Postion == self.mHandCards.Len() {
			return errorValue.ERET_DEAL
		}
		//self.Remove(Postion)
		heap.Remove(&(self.mHandCards), Postion)
	}
	heap.Push(&(self.mMingGangCards), tagMj)
	self.mLastMingGangCard = tagMj
	self.mUserAction = Mj_User_Action_MingGang
	return errorValue.ERET_OK
}

//暗杠
func (self *MjLogicUser) AnGangCard(tagMj uint32) uint32 {
	//mLastGuoDetailInfo.Init(); //操作过手，清楚过操作信息。
	for i := 0; i < 4; i++ {
		lIt := self.mHandCards.Find(tagMj)
		if lIt == self.mHandCards.Len() {
			return errorValue.ERET_DEAL
		}
		heap.Remove(&(self.mHandCards), lIt)
	}
	heap.Push(&(self.mAnGangCards), tagMj)
	self.mLastAnGangCard = tagMj
	self.mUserAction = Mj_User_Action_AnGang

	return errorValue.ERET_OK
}

//补杠
func (self *MjLogicUser) BuGangCard(tagMj uint32) uint32 {
	//mLastGuoDetailInfo.Init(); //操作过手，清楚过操作信息。
	lIt := self.mHandCards.Find(tagMj)
	if lIt == self.mHandCards.Len() {
		return errorValue.ERET_DEAL
	}

	heap.Remove(&(self.mHandCards), lIt)

	lIt = self.mPengCards.Find(tagMj)
	if lIt == self.mPengCards.Len() {
		return errorValue.ERET_DEAL
	}

	heap.Remove(&(self.mPengCards), lIt)
	heap.Push(&(self.mBuGangCards), tagMj)
	self.mLastBuGangCard = tagMj
	self.mUserAction = Mj_User_Action_BuGang
	return errorValue.ERET_OK
}

/*
   胡了
   参数1:出了这张胡牌的玩家的座位id。(有可能是自己即自摸，别人就是点炮)
   参数2:胡的牌
*/
func (self *MjLogicUser) HuCard(seatId uint32, tagMj uint32) uint32 {
	//self.mLastGuoDetailInfo.Init(); //操作过手，清楚过操作信息。
	if seatId < 0 || seatId > SIT_NUM_MAX || tagMj <= MJ_BACK || tagMj > MJ_9_TIAO {
		return errorValue.ERET_HU
	}

	if seatId == self.mSeatId {
		//自摸
		lResult := self.handleZiMo(tagMj)
		if lResult == errorValue.ERET_HU {
			return errorValue.ERET_HU
		}
	} else {
		//点炮
		lResult := self.handleDianPao(seatId, tagMj)
		if lResult == errorValue.ERET_HU {
			return errorValue.ERET_HU
		}
	}

	self.SetUserState(rpc.MjUserState_Mj_User_State_Hu)
	self.mHuCard = tagMj
	self.mUserAction = Mj_User_Action_HuPai
	return errorValue.ERET_OK
}

func (self *MjLogicUser) DeleteLastDealCard(tagMj uint32) uint32 {
	lIt := self.mDealCards.Find(tagMj)
	if lIt == self.mDealCards.Len() {
		return errorValue.ERET_DEAL
	}

	heap.Remove(&(self.mDealCards), lIt)
	return errorValue.ERET_OK
}

func (self *MjLogicUser) DeleteLastBuGang(tagMj uint32) uint32 {
	lIt := self.mBuGangCards.Find(tagMj)
	if lIt == self.mBuGangCards.Len() {
		return errorValue.ERET_DEAL
	}

	heap.Remove(&(self.mBuGangCards), lIt)
	heap.Push(&(self.mPengCards), tagMj)
	return errorValue.ERET_OK
}

func (self *MjLogicUser) GetNextValidMJ(MjCard *IntHeap, nCount uint32, nStartIndex uint32, card *uint32) int {
	var lResult int = -1
	var i uint32 = 0
	*card = MJ_BACK

	if MjCard.Len() == 0 {
		return lResult
	}

	if nStartIndex < 0 || nStartIndex >= nCount {
		return lResult
	}

	for i = nStartIndex; i < nCount; i++ {
		if MJ_BACK != (*MjCard)[i] {
			*card = (*MjCard)[i]
			lResult = int(i)
			break
		}
	}

	return lResult
}

func (self *MjLogicUser) GetTagMjCountInHand(tagMj *uint32) uint32 {
	var lMjCount uint32 = 0
	handMjSize := self.mHandCards.Len()
	if handMjSize > 0 {
		for lit := range self.mHandCards {
			if uint32(lit) == *tagMj {
				lMjCount += 1
			}
		}
	}
	return lMjCount
}

func (self *MjLogicUser) GetUserState() rpc.MjUserState {
	return self.mUserState
}

func (self *MjLogicUser) GetRejectSuit() uint32 {
	return self.mRejectSuit
}

func (self *MjLogicUser) GetHandCards() IntHeap {
	return self.mHandCards
}

func (self *MjLogicUser) GetDealCards() *IntHeap {
	return &(self.mDealCards)
}

func (self *MjLogicUser) GetMingGangCards() *IntHeap {
	return &(self.mMingGangCards)
}

func (self *MjLogicUser) GetAnGangCards() *IntHeap {
	return &(self.mAnGangCards)
}

func (self *MjLogicUser) GetBuGangCards() *IntHeap {
	return &(self.mBuGangCards)
}

func (self *MjLogicUser) GetPengCards() *IntHeap {
	return &(self.mPengCards)
}

func (self *MjLogicUser) GetHuanCards() *IntHeap {
	return &(self.mHuanCards)
}

func (self *MjLogicUser) GetHuanInCards() *IntHeap {
	return &(self.mHuanInCards)
}

func (self *MjLogicUser) GetLastOperateRecord() interface{} {
	/*
		if self.mOperateRecord.Len() <= 0 {
			return errorValue.ERET_SYS_ERR
		}
	*/
	return self.mOperateRecord.Back().Value
}

func (self *MjLogicUser) GetCardCount(card uint32) uint32 {
	var lCount uint32 = 0
	for _, p := range self.mHandCards {
		if p == card {
			lCount += 1
		}
	}

	for _, p := range self.mPengCards {
		if p == card {
			lCount += 3
		}
	}
	for _, p := range self.mAnGangCards {
		if p == card {
			lCount += 4
		}
	}
	for _, p := range self.mBuGangCards {
		if p == card {
			lCount += 4
		}
	}
	for _, p := range self.mMingGangCards {
		if p == card {
			lCount += 4
		}
	}

	if self.mHuCard == card {
		lCount += 1
	}

	return lCount
}

func (self *MjLogicUser) GetGameCurCards(cards *IntHeap) uint32 {
	for _, p := range self.mHandCards {
		heap.Push(cards, p)
	}

	for _, p := range self.mPengCards {
		heap.Push(cards, p)
	}
	for _, p := range self.mAnGangCards {
		heap.Push(cards, p)
	}
	for _, p := range self.mBuGangCards {
		heap.Push(cards, p)
	}
	for _, p := range self.mMingGangCards {
		heap.Push(cards, p)
	}

	if self.mHuCard == MJ_BACK {
		heap.Push(cards, self.mHuCard)
	}
	return errorValue.ERET_OK
}

func (self *MjLogicUser) CanOperater() bool {
	for p := self.mOperateMask.Front(); p != nil; p = p.Next() {
		if (p.Value).(rpc.MJ_OPERATE_MASK) != rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_NONE {
			return true
		}
	}

	return false
}

func (self *MjLogicUser) CanDeal() bool {
	lIt := self.FindValueInList(self.mOperateMask, rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_DEAL)
	if lIt != nil {
		return true
	}
	return false
}

func (self *MjLogicUser) CanDealByTag(TagMj uint32) bool {
	lIt := self.FindValueInList(self.mOperateMask, rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_DEAL)
	if lIt == nil {
		return false
	}
	if self.IsHaveRejectSuit() && MJ_GET_SUIT(TagMj) != self.mRejectSuit {
		return false
	}
	return true
}

func (self *MjLogicUser) CanHuanSanZhang() bool {
	lIt := self.FindValueInList(self.mOperateMask, rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_HUAN_SAN_ZHANG)
	if lIt != nil {
		return true
	}
	return false
}

func (self *MjLogicUser) CanReject() bool {
	lIt := self.FindValueInList(self.mOperateMask, rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_REJECTSUIT)
	if lIt != nil {
		return true
	}
	return false
}

func (self *MjLogicUser) CanPeng() bool {
	lIt := self.FindValueInList(self.mOperateMask, rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_PENG)
	if lIt != nil {
		return true
	}
	return false
}

func (self *MjLogicUser) CanPengByTag(TagMj uint32) bool {
	lIt := self.FindValueInList(self.mOperateMask, rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_PENG)
	if lIt == nil {
		return false
	}

	if self.IsCanPengTagMj(TagMj) == false {
		return false
	}
	return true
}

func (self *MjLogicUser) CanMingGangByTag(TagMj uint32) bool {
	lPostion := self.FindValueInList(self.mOperateMask, rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_MING_GANG)
	if lPostion == nil {
		return false
	}

	if self.IsCanMingGangTagMj(TagMj) == false {
		return false
	}

	return true
}

func (self *MjLogicUser) CanMingGang() bool {
	lPostion := self.FindValueInList(self.mOperateMask, rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_MING_GANG)
	if lPostion != nil {
		return true
	}
	return false
}

func (self *MjLogicUser) CanAnGangByTag(TagMj uint32) bool {

	if TagMj <= MJ_BACK || TagMj > MJ_9_TIAO || (MJ_GET_SUIT(TagMj) == self.mRejectSuit) {
		return false
	}

	lPostion := self.FindValueInList(self.mOperateMask, rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_AN_GANG)
	if lPostion == nil {
		return false
	}

	lResult := self.GetTagMjCountInHand(&TagMj)
	if lResult != GANG_MJ_COUNT {
		return false
	}
	return true
}

func (self *MjLogicUser) CanAnGang() bool {
	lPostion := self.FindValueInList(self.mOperateMask, rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_AN_GANG)
	if lPostion != nil {
		return true
	}
	return false
}

func (self *MjLogicUser) CanBuGangByTag(TagMj uint32) bool {

	if TagMj <= MJ_BACK || TagMj > MJ_9_TIAO || (MJ_GET_SUIT(TagMj) == self.mRejectSuit) {
		return false
	}

	lPostion := self.FindValueInList(self.mOperateMask, rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_BU_GANG)
	if lPostion == nil {
		return false
	}

	lResult := self.IsCanBuGangTagMj(TagMj)
	if lResult != false {
		return false
	}
	return true
}

func (self *MjLogicUser) CanBuGang() bool {
	lPostion := self.FindValueInList(self.mOperateMask, rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_BU_GANG)
	if lPostion != nil {
		return true
	}
	return false
}

func (self *MjLogicUser) CanHuByTag(HuMj uint32) bool {
	lIt := self.FindValueInList(self.mOperateMask, rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_HU)
	if lIt == nil {
		return false
	}

	if HuMj <= MJ_BACK || HuMj > MJ_9_TIAO || (MJ_GET_SUIT(HuMj) == self.mRejectSuit || self.IsHaveRejectSuit() == true) {
		return false
	}

	if self.IsHandCards3NAnd2() {
		if self.mDealCards.Len() != 0 {
			if self.mLastAddCard != HuMj {
				return false
			}
		} else {
			if self.mLastAddCard != 0 {
				if self.mLastAddCard != HuMj {
					return false
				}
			} else {
				lTrue := false
				for _, it := range self.mHandCards {
					if it == HuMj {
						lTrue = true
						break
					}
				}
				if lTrue == false {
					return false
				}
			}
		}
		if self.IsCanHu(&(self.mHandCards), uint32(self.mHandCards.Len()), uint32(0), false) == false {
			return false
		}
	} else {
		if self.IsCanHu(&(self.mHandCards), uint32(self.mHandCards.Len()), uint32(0), false) == false {
			return false
		}
	}

	return true
}

func (self *MjLogicUser) CanHu() bool {
	lPostion := self.FindValueInList(self.mOperateMask, rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_HU)
	if lPostion != nil {
		return true
	}
	return false
}

func (self *MjLogicUser) CanGuo() bool {
	lPostion := self.FindValueInList(self.mOperateMask, rpc.MJ_OPERATE_MASK_MJ_OPERATE_MASK_GUO)
	if lPostion != nil {
		return true
	}
	return false
}

func (self *MjLogicUser) IsOnActivestate() bool {
	if self.IsFree() || self.mIsOffLine || self.mIsLeaveTable {
		return false
	}

	return true
}

func (self *MjLogicUser) IsHaveOperateRecord() bool {
	if self.mOperateRecord.Len() > 0 {
		return true
	}
	return false
}

func (self *MjLogicUser) IsHandCards3NAnd2() bool {

	lResult := self.mHandCards.Len() % 3
	if lResult == 2 {
		return true
	}

	return false
}

func (self *MjLogicUser) IsCanOperateForDealCard(seatId uint32, tagMj uint32, roomCfg *common.RoomConfig) bool {
	//TODO
	return false
}

func (self *MjLogicUser) IsCanOperateForBuGang(seatId uint32, tagMj uint32, roomCfg *common.RoomConfig) bool {
	//TODO
	return false
}
func (self *MjLogicUser) IsCanHu(mHandCardTmp *IntHeap, mHandCardSize uint32, nProcessedCount uint32, bFoundJiang bool) bool {
	lRet := false

	lRet = self.TestHu7Dui(mHandCardTmp, mHandCardSize)
	if lRet == true {
		return true
	}
	lRet = self.TestNormalCanHu(mHandCardTmp, mHandCardSize, nProcessedCount, bFoundJiang)
	if lRet == true {
		return true
	}

	return false
}

func (self *MjLogicUser) IsCanMingGangTagMj(tagMj uint32) bool {
	lResult := false

	if tagMj <= MJ_BACK || tagMj > MJ_9_TIAO || (MJ_GET_SUIT(tagMj) == self.mRejectSuit) {
		return lResult
	}

	lTagMjCount := self.GetTagMjCountInHand(&tagMj)
	if (lTagMjCount + 1) == GANG_MJ_COUNT {
		lResult = true
	}

	return lResult
}

func (self *MjLogicUser) IsCanBuGangTagMj(tagMj uint32) bool {
	lResult := false

	if tagMj <= MJ_BACK || tagMj > MJ_9_TIAO || (MJ_GET_SUIT(tagMj) == self.mRejectSuit) {
		return lResult
	}

	for _, lIt := range self.mPengCards {
		if lIt == tagMj {
			lResult = true
			break
		}
	}

	return lResult
}

func (self *MjLogicUser) IsHaveRejectSuit() bool {

	lResult := false
	for _, lMj := range self.mHandCards {
		if MJ_GET_SUIT(lMj) == self.mRejectSuit {
			lResult = true
			break
		}
	}

	return lResult
}

func (self *MjLogicUser) IsCanPengTagMj(tagMj uint32) bool {
	lResult := false

	if tagMj <= 0 || tagMj > MJ_9_TIAO || (MJ_GET_SUIT(tagMj) == self.mRejectSuit) {
		return lResult
	}

	lTagMjCount := self.GetTagMjCountInHand(&tagMj)
	if (lTagMjCount + 1) >= PENG_MJ_COUNT {
		lResult = true
	}

	return lResult
}

func (self *MjLogicUser) IsMatch() bool {
	if self.mUserState == rpc.MjUserState_Mj_User_State_Init || self.mUserState == rpc.MjUserState_Mj_User_State_Observer_Sit {
		return false
	}

	return true
}

func (self *MjLogicUser) IsWaitRecharge() bool {
	if self.mUserState == rpc.MjUserState_Mj_User_State_Wait_Recharge {
		return true
	}

	return false
}

func (self *MjLogicUser) IsGiveUp() bool {
	if self.mUserState == rpc.MjUserState_Mj_User_State_GiveUp {
		return true
	}

	return false
}

func (self *MjLogicUser) IsOffLine() bool {

	return self.mIsOffLine
}

func (self *MjLogicUser) IsFree() bool {
	if self.mUserState == rpc.MjUserState_Mj_User_State_Init {
		return true
	}
	return false
}

func (self *MjLogicUser) IsCanContinueGame() bool {
	if self.IsFree() || self.mIsLeaveTable || self.mUserState == rpc.MjUserState_Mj_User_State_GiveUp {
		return false
	}

	return true
}

func (self *MjLogicUser) IsHu() bool {
	if self.mUserState != rpc.MjUserState_Mj_User_State_Hu {
		return false
	}
	return true
}

func (self *MjLogicUser) IsReady() bool {
	if self.mUserState != rpc.MjUserState_Mj_User_State_Ready {
		return false
	}
	return true
}

func (self *MjLogicUser) IsLeaveTable() bool {
	return self.mIsLeaveTable
}

func (self *MjLogicUser) IsEnd() bool {
	if self.mUserState == rpc.MjUserState_Mj_User_State_Hu || self.mUserState == rpc.MjUserState_Mj_User_State_GiveUp {
		return true
	}

	return false
}

func (self *MjLogicUser) AiCalculateOperateInfo(mask rpc.MJ_OPERATE_MASK, tagMj uint32, tagMj1 uint32, tagMj2 uint32) uint32 {
	return errorValue.ERET_OK
}

func (self *MjLogicUser) GetStartMatchCurrencyValue() int64 {
	return self.mStartMatchCurrencyValue
}

func (self *MjLogicUser) InitStartMatchCurrencyValue(currencyType rpc.CURRENCY_TYPE) uint32 {
	//TODO:
	//访问数据库，获取金币数量
	var lRet uint32 = errorValue.ERET_OK
	if currencyType == rpc.CURRENCY_TYPE_CURRENCY_TYPE_GOLD {
		self.mStartMatchCurrencyValue = int64(0)
	} else if currencyType == rpc.CURRENCY_TYPE_CURRENCY_TYPE_POINT {
		self.mStartMatchCurrencyValue = int64(0)
	}
	return lRet
}

func (self *MjLogicUser) GetSettlementType() rpc.MJ_SETTLEMENT_TYPE {
	//TODO : 获取结算结果
	return rpc.MJ_SETTLEMENT_TYPE_MJ_SETTLEMENT_TYPE_NONE
}

func (self *MjLogicUser) GetSettlementInfo() *list.List {
	return self.mSettlementInfo
}

func (self *MjLogicUser) GetLastGameInfo() *MjLastUserGameInfo {
	return self.mLastGameInfo
}

func (self *MjLogicUser) CacheLastGameInfo() {
	self.mLastGameInfo.RejectSuit = self.mRejectSuit
	self.mLastGameInfo.HuCard = self.mHuCard
	self.mLastGameInfo.DealCards = self.mDealCards
	self.mLastGameInfo.HandCards = self.mHandCards
	self.mLastGameInfo.MingGangCards = self.mMingGangCards
	self.mLastGameInfo.AnGangCards = self.mAnGangCards
	self.mLastGameInfo.BuGangCards = self.mBuGangCards
	self.mLastGameInfo.PengCards = self.mPengCards
	self.mLastGameInfo.InitHandCards = self.mInitHandCards
	self.mLastGameInfo.HuanCards = self.mHuanCards
	self.mLastGameInfo.HuanInCards = self.mHuanInCards
	self.mLastGameInfo.SettlementInfo = self.mSettlementInfo
	self.mLastGameInfo.OperateRecord = self.mOperateRecord
}
