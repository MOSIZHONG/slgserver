package entity

import (
	"go.uber.org/zap"
	"slgserver/db"
	"slgserver/log"
	"slgserver/model"
	"sync"
)

type ArmyMgr struct {
	mutex     sync.RWMutex
	armyById  map[int]*model.Army
	armByCityId map[int][]*model.Army

}

var AMgr = &ArmyMgr{
	armyById: make(map[int]*model.Army),
	armByCityId: make(map[int][]*model.Army),
}

func (this* ArmyMgr) Get(aid int) (*model.Army, error){
	this.mutex.RLock()
	a, ok := this.armyById[aid]
	this.mutex.RUnlock()
	if ok {
		return a, nil
	}else{
		army := &model.Army{}
		ok, err := db.MasterDB.Table(model.Army{}).Where("id=?", aid).Get(army)
		if ok {
			this.mutex.Lock()
			this.armyById[aid] = army
			if _, r:= this.armByCityId[army.CityId]; r == false{
				this.armByCityId[army.CityId] = make([]*model.Army, 0)
			}
			this.armByCityId[army.CityId] = append(this.armByCityId[army.CityId], army)
			this.mutex.Unlock()
			return army, nil
		}else{
			return nil, err
		}
	}
}

func (this* ArmyMgr) GetOrCreate(rid int, cid int, order int8) (*model.Army, error){

	this.mutex.RLock()
	armys, ok := this.armByCityId[cid]
	if ok {
		for _, v := range armys {
			if v.Order == order{
				return v, nil
			}
		}
	}

	//需要创建
	a := &model.Army{RId: rid, Order: order, CityId: cid,
		FirstId: -1, SecondId: -1, ThirdId: -1,
		FirstSoldierCnt: 0, SecondSoldierCnt: 0, ThirdSoldierCnt: 0}
	_, err := db.MasterDB.Insert(a)
	if err == nil{
		return a, nil
	}else{
		log.DefaultLog.Warn("db error", zap.Error(err))
		return nil, err
	}
}

