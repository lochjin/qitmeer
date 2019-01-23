package blockchain

import (
	"github.com/noxproject/nox/common/anticone"
	"fmt"
	"container/list"
	"sort"
	"github.com/noxproject/nox/common/hash"
)

const (
	BlockDelay=15
	BlockRate=0.02
	SecurityLevel=0.01
)

type Phantom struct {
	// The general foundation framework of DAG
	bd *BlockDAG

	// This is a set that only include honest block and it is the common part of each
	// tips in the block dag, so it is a blue set too.
	commonBlueSet    *BlockSet

	// This is a set that only include honest block exclude from "commonBlueSet",
	// but it's not very stable.
	tempBlueSet      *BlockSet

	lastCommonBlocks *BlockSet

	// Well understood,this orderly array is the sorting of common set.
	commonOrder      []*hash.Hash

	// This orderly array is the sorting of the end of dag set.
	tempOrder        []*hash.Hash

	// If it happens that during two propagation delays only one block is created, this block is called hourglass block.
	// This means it reference all the tips and is reference by all following blocks.
	// Hourglass block is a strong signal of finality because its blue set is stable.
	hourglassBlocks *BlockSet

	// The block anticone size is all in the DAG which did not reference it and
	// were not referenced by it.
	anticoneSize int

	// The past set of block
	pastSetNum map[hash.Hash]uint
}

func (ph *Phantom) GetName() string {
	return phantom
}

func (ph *Phantom) Init(bd *BlockDAG) bool {
	ph.bd=bd

	ph.anticoneSize = anticone.GetSize(BlockDelay,BlockRate,SecurityLevel)

	ph.pastSetNum=map[hash.Hash]uint{}

	log.Info(fmt.Sprintf("anticone size:%d",ph.anticoneSize))

	return true
}

// This is an entry for update the block dag,you need pass in a block parameter,
// If add block have failure,it will return false.
func (ph *Phantom) AddBlock(b *Block) *list.List {
	if b == nil {
		return nil
	}
	ph.tempBlueSet=nil
	log.Trace(fmt.Sprintf("Add block:%v",b.GetHash().String()))

	ph.calculatePastBlockSetNum(b)
	ph.updateCommonBlueSet(b.GetHash())
	ph.updateHourglass()

	return ph.updateOrder(b)
}

// The past set is all its its ancestors.Because the past cannot be
// changed, so its number is fixed.
func (ph *Phantom) addPastSetNum(b *Block, num uint) {
	ph.pastSetNum[*b.GetHash()]=num
}

func (ph *Phantom) GetPastSetNum(b *Block) uint {
	return ph.pastSetNum[*b.GetHash()]
}

func isVirtualTip(b *Block, futureSet *BlockSet, anticone *BlockSet, children *BlockSet) bool {
	for k, _ := range children.GetMap() {
		if k.IsEqual(b.GetHash()) {
			return false
		}
		if !futureSet.Has(&k) && !anticone.Has(&k) {
			return false
		}
	}
	return true
}

// This function is used to GetAnticone recursion
func (ph *Phantom) recAnticone(b *Block, futureSet *BlockSet, anticone *BlockSet, h *hash.Hash) {
	if h.IsEqual(b.GetHash()) {
		return
	}
	node:=ph.bd.GetBlock(h)
	children := node.GetChildren()
	needRecursion := false
	if children == nil || children.Len() == 0 {
		needRecursion = true
	} else {
		needRecursion = isVirtualTip(b, futureSet, anticone, children)
	}
	if needRecursion {
		if !futureSet.Has(h) {
			anticone.Add(h)
		}
		parents := node.GetParents()

		//Because parents can not be empty, so there is no need to judge.
		for k, _ := range parents.GetMap() {
			ph.recAnticone(b, futureSet, anticone, &k)
		}
	}
}

// This function can get anticone set for an block that you offered in the block dag,If
// the exclude set is not empty,the final result will exclude set that you passed in.
func (ph *Phantom) GetAnticone(b *Block, exclude *BlockSet) *BlockSet {
	futureSet := NewBlockSet()
	ph.GetFutureSet(futureSet, b)
	anticone := NewBlockSet()
	for k, _ := range ph.bd.tips.GetMap() {
		ph.recAnticone(b, futureSet, anticone, &k)
	}
	if exclude != nil {
		anticone.Exclude(exclude)
	}
	return anticone
}

func (ph *Phantom) GetFutureSet(fs *BlockSet, b *Block) {
	children := b.GetChildren()
	if children == nil || children.IsEmpty() {
		return
	}
	for k, _ := range children.GetMap() {
		if !fs.Has(&k) {
			fs.Add(&k)
			ph.GetFutureSet(fs, ph.bd.GetBlock(&k))
		}
	}
}

// Calculate the size of the past block set.Because the past block set of block
// is stable,we can calculate and save.
func (ph *Phantom) calculatePastBlockSetNum(b *Block) {

	if b.GetHash().IsEqual(ph.bd.GetGenesisHash()) {
		ph.addPastSetNum(b, 0)
		return
	}
	parents:=b.GetParents()
	if parents == nil || parents.IsEmpty() {
		return
	}
	parentsList:=[]*Block{}
	for k,_:=range parents.GetMap(){
		parentsList=append(parentsList,ph.bd.GetBlock(&k))
	}

	if len(parentsList) == 1 {
		ph.addPastSetNum(b, ph.GetPastSetNum(parentsList[0])+1)
		return
	}
	anticone := ph.GetAnticone(b, nil)

	anOther := ph.GetAnticone(parentsList[0], anticone)

	ph.addPastSetNum(b, ph.GetPastSetNum(parentsList[0])+uint(anOther.Len())+1)
}

func (ph *Phantom) sortBlockSet(set *BlockSet, bs *BlockSet) SortBlocks {
	sb0 := SortBlocks{}
	sb1 := SortBlocks{}

	for k, _ := range set.GetMap() {
		node:=ph.bd.GetBlock(&k)
		kv:=k
		if bs != nil && bs.Has(&k) {
			sb0 = append(sb0, SortBlock{&kv, ph.GetPastSetNum(node)})
		} else {
			sb1 = append(sb1, SortBlock{&kv, ph.GetPastSetNum(node)})
		}

	}
	sort.Sort(sb0)
	sort.Sort(sb1)
	sb0 = append(sb0, sb1...)
	return sb0
}

func (ph *Phantom) getPastSetByOrder(pastSet *BlockSet, exclude *BlockSet, h *hash.Hash) {
	if exclude.Has(h) || pastSet.Has(h) {
		return
	}

	if h.IsEqual(ph.bd.GetGenesisHash()) {
		return
	}

	parents := ph.bd.GetBlock(h).GetParents()
	parentsList := parents.List()
	if parents == nil || len(parentsList) == 0 {
		return
	}
	for _, v := range parentsList {

		pastSet.Add(v)
		ph.getPastSetByOrder(pastSet, exclude, v)
	}
}

func (ph *Phantom) GetTempOrder(tempOrder *[]*hash.Hash, tempOrderM *BlockSet, bs *BlockSet, h *hash.Hash, exclude *BlockSet) {

	//1.If h that has already appeared must be excluded.
	if exclude != nil && exclude.Has(h) {
		return
	}
	node:=ph.bd.GetBlock(h)
	parents := node.GetParents()

	//2.If its father hasn't sorted,the function must return.
	if parents != nil && parents.Len() > 0 {
		for k, _ := range parents.GetMap() {
			if exclude != nil && exclude.Has(&k) {
				continue
			}
			if !tempOrderM.Has(&k) {
				return
			}
		}
	}
	var anticone *BlockSet

	//3.Search some uncle block that it is in front of me, then
	//make sure they are sorted.
	if !tempOrderM.Has(h) {
		if !ph.bd.GetGenesisHash().IsEqual(h) && !ph.lastCommonBlocks.Has(h) {
			anticone = ph.GetAnticone(node, exclude)
			//
			if !anticone.IsEmpty() {
				ansb := ph.sortBlockSet(anticone, bs)
				if bs.Has(h) {
					for _, av := range ansb {
						avNode:=ph.bd.GetBlock(av.h)
						if bs.Has(av.h) && ph.GetPastSetNum(avNode) < ph.GetPastSetNum(node) && !tempOrderM.Has(av.h) {
							ph.GetTempOrder(tempOrder, tempOrderM, bs, av.h, exclude)
						}
					}
				} else {
					for _, av := range ansb {
						if bs.Has(av.h) && !tempOrderM.Has(av.h) {
							ph.GetTempOrder(tempOrder, tempOrderM, bs, av.h, exclude)
						}
					}
				}

			}
		}

	}

	//4.Add myself to the array
	if !tempOrderM.Has(h) {
		(*tempOrder) = append(*tempOrder, h)
		tempOrderM.Add(h)
	}

	//5.Sort all my children
	childrenSrc := node.GetChildren()
	children := childrenSrc.Clone()
	if exclude != nil {
		children.Exclude(exclude)
	}
	if children == nil || children.Len() == 0 {
		return
	}
	pastSet := NewBlockSet()
	redSet := NewBlockSet()
	sb := ph.sortBlockSet(children, bs)

	for _, v := range sb {

		if bs.Has(v.h) {
			if !tempOrderM.Has(v.h) {
				pastSet.Clear()
				redSet.Clear()
				var excludeT *BlockSet
				if exclude != nil {
					excludeT = tempOrderM.Clone()
					excludeT.AddSet(exclude)
				} else {
					excludeT = tempOrderM
				}

				ph.getPastSetByOrder(pastSet, excludeT, v.h)

				inbs := pastSet.Intersection(anticone)
				if inbs != nil && inbs.Len() > 0 {
					insb := ph.sortBlockSet(inbs, bs)

					for _, v0 := range insb {
						if bs.Has(v0.h) {
							if !tempOrderM.Has(v0.h) {
								ph.GetTempOrder(tempOrder, tempOrderM, bs, v0.h, exclude)
							}
						} else {
							redSet.Add(v0.h)
						}
					}
					if !redSet.IsEmpty() {
						pastSet.Exclude(redSet)
						isAllOrder := true
						for k, _ := range pastSet.GetMap() {
							if !tempOrderM.Has(&k) {
								isAllOrder = false
								break
							}
						}
						if isAllOrder {
							redsb := ph.sortBlockSet(redSet, bs)
							for _, v1 := range redsb {
								ph.GetTempOrder(tempOrder, tempOrderM, bs, v1.h, exclude)
							}
						}
					}

				}
			}
			ph.GetTempOrder(tempOrder, tempOrderM, bs, v.h, exclude)
		}
	}
	for _, v := range sb {
		if !bs.Has(v.h) {
			ph.GetTempOrder(tempOrder, tempOrderM, bs, v.h, exclude)
		}
	}
}

func (ph *Phantom) updateCommonOrder(tip *hash.Hash, blueSet *BlockSet, isRollBack bool, exclude *BlockSet, curLastCommonBS *BlockSet, startIndex int) {

	if tip.IsEqual(ph.bd.GetGenesisHash()) {
		ph.commonOrder = []*hash.Hash{}
		return
	}
	node:=ph.bd.GetBlock(tip)
	parents := node.GetParents()

	if parents.HasOnly(ph.bd.GetGenesisHash()) {
		if len(ph.commonOrder) == 0 {
			ph.commonOrder = append(ph.commonOrder,ph.bd.GetGenesisHash())
		}
	}

	if !isRollBack {
		if blueSet == nil {
			return
		}
		tempOrder := []*hash.Hash{}
		tempOrderM := NewBlockSet()

		lpsb := ph.sortBlockSet(ph.lastCommonBlocks, blueSet)

		for _, v := range lpsb {
			ph.GetTempOrder(&tempOrder, tempOrderM, blueSet, v.h, exclude)
		}
		toLen := len(tempOrder)
		var poLen int = 0
		for i := 0; i < toLen; i++ {
			if ph.lastCommonBlocks.Has(tempOrder[i]) {
				continue
			}
			index := startIndex + i
			poLen = len(ph.commonOrder)
			if index < poLen {
				ph.commonOrder[index] = tempOrder[i]
			} else {
				ph.commonOrder = append(ph.commonOrder, tempOrder[i])
			}
		}
		poLen = len(ph.commonOrder)
		for i := poLen - 1; i >= 0; i-- {
			if ph.commonOrder[i]!=nil {
				if !curLastCommonBS.Has(ph.commonOrder[i]) {
					log.Error("order errer:end block is not new common block")
				}
				break
			}

		}

	} else {
		poLen := len(ph.commonOrder)
		rNum := 0
		for i := poLen - 1; i >= 0; i-- {
			if curLastCommonBS.Has(ph.commonOrder[i]) {
				break
			}
			ph.commonOrder[i] = nil
			rNum++
		}
		if (poLen - rNum) != startIndex {
			log.Error("order errer:number")
		}
	}
}

func (ph *Phantom) recPastBlockSet(genealogy *BlockSet, tipsAncestors *map[hash.Hash]*BlockSet, tipsGenealogy *map[hash.Hash]*BlockSet) {

	var maxPastHash *hash.Hash = nil
	var maxPastNum uint = 0
	var tipsHash *hash.Hash = nil

	for tk, v := range *tipsAncestors {
		tkv:=tk
		if v.Len() == 1 && v.Has(ph.bd.GetGenesisHash()) {
			continue
		}

		for k, _ := range v.GetMap() {
			kv:=k
			node:=ph.bd.GetBlock(&kv)
			pastNum := ph.GetPastSetNum(node)
			if maxPastHash == nil || maxPastNum < pastNum {
				maxPastHash = &kv
				maxPastNum = pastNum
				tipsHash = &tkv
			}
		}

	}
	if maxPastHash == nil {
		return
	}
	parents := ph.bd.GetBlock(maxPastHash).GetParents()
	if parents == nil || parents.Len() == 0 {
		return
	}
	(*tipsAncestors)[*tipsHash].Remove(maxPastHash)
	for k, _ := range parents.GetMap() {

		if !(*tipsGenealogy)[*tipsHash].Has(&k) {
			(*tipsAncestors)[*tipsHash].Add(&k)
			(*tipsGenealogy)[*tipsHash].Add(&k)
			if genealogy != nil {
				genealogy.Add(&k)
			}
		}

	}
}

func (ph *Phantom) calLastCommonBlocks(tip *hash.Hash) *BlockSet {
	tips := ph.bd.GetTips()
	if tips == nil {
		return nil
	}
	tipsList := tips.List()
	if len(tipsList) <= 1 {
		return nil
	}
	tipsGenealogy:=make(map[hash.Hash]*BlockSet)
	tipsAncestors := make(map[hash.Hash]*BlockSet)
	for _, v := range tipsList {
		tipsAncestors[*v] = NewBlockSet()
		tipsAncestors[*v].Add(v)

		tipsGenealogy[*v]=NewBlockSet()
		tipsGenealogy[*v].Add(v)
	}

	//
	for {
		hasDifferent := false
		for k, v := range tipsAncestors {
			if k.IsEqual(tip) {
				continue
			}
			if !tipsAncestors[*tip].IsEqual(v) {
				hasDifferent = true
				break
			}
		}
		if !hasDifferent {
			break
		}
		ph.recPastBlockSet(nil, &tipsAncestors, &tipsGenealogy)
	}
	return tipsAncestors[*tip]
}

func (ph *Phantom) calLastCommonBlocksPBS(pastBlueSet *map[hash.Hash]*BlockSet) {
	/////
	lastPFuture := NewBlockSet()
	for k, _ := range ph.lastCommonBlocks.GetMap() {
		ph.GetFutureSet(lastPFuture, ph.bd.GetBlock(&k))
	}

	if ph.lastCommonBlocks.Len() == 1 {
		lpbHash := ph.lastCommonBlocks.List()[0]
		if pastBlueSet != nil {
			(*pastBlueSet)[*lpbHash] = NewBlockSet()
		}

		//pastBlueSet[lpbHash].Add(lpbHash)

	} else {
		lastTempBlueSet := NewBlockSet()
		lpbAnti := make(map[hash.Hash]*BlockSet)

		for k, _ := range ph.lastCommonBlocks.GetMap() {
			lpbAnti[k] = ph.GetAnticone(ph.bd.GetBlock(&k), lastPFuture)
			lastTempBlueSet.AddSet(lpbAnti[k])
		}
		if pastBlueSet != nil {
			for k, _ := range lastTempBlueSet.GetMap() {
				if !ph.commonBlueSet.Has(&k) {
					lastTempBlueSet.Remove(&k)
				}
			}
			for k, _ := range ph.lastCommonBlocks.GetMap() {
				(*pastBlueSet)[k] = lastTempBlueSet.Clone()
				(*pastBlueSet)[k].Exclude(lpbAnti[k])
				(*pastBlueSet)[k].Remove(&k)
			}
		}

	}
}

func (ph *Phantom) calculateBlueSet(parents *BlockSet, exclude *BlockSet, pastBlueSet *map[hash.Hash]*BlockSet, useCommon bool) *BlockSet {

	parentsPBSS := make(map[hash.Hash]*BlockSet)
	for k, _ := range parents.GetMap() {
		if _, ok := (*pastBlueSet)[k]; ok {
			parentsPBSS[k] = (*pastBlueSet)[k]
		} else {
			parentsPBSS[k] = NewBlockSet()
		}

	}

	maxBluePBSHash := GetMaxLenBlockSet(parentsPBSS)
	if maxBluePBSHash == nil {
		return nil
	}
	//
	result := NewBlockSet()
	result.AddSet(parentsPBSS[*maxBluePBSHash])
	result.Add(maxBluePBSHash)

	if parents.Len() == 1 {
		return result
	}

	maxBlueAnBS := ph.GetAnticone(ph.bd.GetBlock(maxBluePBSHash), exclude)

	//

	if maxBlueAnBS != nil && maxBlueAnBS.Len() > 0 {

		for k, _ := range maxBlueAnBS.GetMap() {
			bAnBS := ph.GetAnticone(ph.bd.GetBlock(&k), exclude)
			if bAnBS == nil || bAnBS.Len() == 0 {
				continue
			}
			inBS := result.Intersection(bAnBS)
			if useCommon {
				inPBS := ph.commonBlueSet.Intersection(bAnBS)
				inBS.AddSet(inPBS)
			}

			if inBS == nil || inBS.Len() <= ph.anticoneSize {
				result.Add(&k)
			}
		}
	}
	return result
}

func (ph *Phantom) calculatePastBlueSet(h *hash.Hash, pastBlueSet *map[hash.Hash]*BlockSet, useCommon bool) {

	_, ok := (*pastBlueSet)[*h]
	if ok {
		return
	}
	if h.IsEqual(ph.bd.GetGenesisHash()) {
		(*pastBlueSet)[*h] = NewBlockSet()
		return
	}
	//
	parents := ph.bd.GetBlock(h).GetParents()
	if parents == nil || parents.IsEmpty() {
		return
	} else if parents.HasOnly(ph.bd.GetGenesisHash()) {
		(*pastBlueSet)[*h] = NewBlockSet()
		(*pastBlueSet)[*h].Add(ph.bd.GetGenesisHash())
		return
	}

	for k, _ := range parents.GetMap() {
		ph.calculatePastBlueSet(&k, pastBlueSet, useCommon)
	}
	//
	anticone := ph.GetAnticone(ph.bd.GetBlock(h), nil)
	(*pastBlueSet)[*h] = ph.calculateBlueSet(parents, anticone, pastBlueSet, useCommon)
}

func (ph *Phantom) updateCommonBlueSet(tip *hash.Hash){

	if tip.IsEqual(ph.bd.GetGenesisHash()) {
		//needOrderBS.Add(tip)
		ph.commonBlueSet = NewBlockSet()
		ph.lastCommonBlocks = NewBlockSet()
		ph.updateCommonOrder(tip, nil, false, nil, nil, 0)

		return
	}
	parents := ph.bd.GetBlock(tip).GetParents()

	if parents.HasOnly(ph.bd.GetGenesisHash()) {
		//needOrderBS.AddList(bd.tempOrder)
		ph.commonBlueSet.Clear()
		ph.commonBlueSet.Add(ph.bd.GetGenesisHash())
		ph.lastCommonBlocks.Clear()
		ph.lastCommonBlocks.Add(ph.bd.GetGenesisHash())
		ph.updateCommonOrder(tip, nil, false, nil, nil, 0)

	} else {
		tips := ph.bd.GetTips()
		if tips.Len() <= 1 {
			//needOrderBS.Add(tip)
			return
		}
		curLastCommonBS := ph.calLastCommonBlocks(tip)
		if curLastCommonBS.IsEqual(ph.lastCommonBlocks) {
			return
		}
		curLPFuture := NewBlockSet()
		for k, _ := range curLastCommonBS.GetMap() {
			ph.GetFutureSet(curLPFuture, ph.bd.GetBlock(&k))
		}

		lastPFuture := NewBlockSet()
		for k, _ := range ph.lastCommonBlocks.GetMap() {
			ph.GetFutureSet(lastPFuture, ph.bd.GetBlock(&k))
		}
		//
		pastBlueSet := make(map[hash.Hash]*BlockSet)

		if lastPFuture.Contain(curLPFuture) {
			//needOrderBS.AddSet(lastPFuture)
			//
			oExclude := NewBlockSet()
			oExclude.AddSet(curLPFuture)
			for k, _ := range ph.lastCommonBlocks.GetMap() {
				oExclude.AddSet(ph.bd.GetBlock(&k).GetParents())
			}

			ph.calLastCommonBlocksPBS(&pastBlueSet)

			for k, _ := range curLastCommonBS.GetMap() {
				ph.calculatePastBlueSet(&k, &pastBlueSet, false)
			}
			commonBlueSet := ph.calculateBlueSet(curLastCommonBS, curLPFuture, &pastBlueSet, false)
			//
			ph.updateCommonOrder(tip, commonBlueSet, false, oExclude, curLastCommonBS, int(ph.bd.GetBlockTotal())-lastPFuture.Len())
			//
			ph.commonBlueSet.AddSet(commonBlueSet)
			ph.lastCommonBlocks = curLastCommonBS
		} else if curLPFuture.Contain(lastPFuture) {
			//needOrderBS.AddSet(curLPFuture)

			ph.updateCommonOrder(tip, nil, true, nil, curLastCommonBS, int(ph.bd.GetBlockTotal())-curLPFuture.Len())
			ph.commonBlueSet.Exclude(curLPFuture)
			ph.lastCommonBlocks = curLastCommonBS
		} else {
			log.Error("error:common set")
		}

	}

}

func (ph *Phantom) GetTempBlueSet() *BlockSet {
	//
	tips := ph.bd.GetTips()
	//

	result := NewBlockSet()
	if tips.HasOnly(ph.bd.GetGenesisHash()) {
		result = NewBlockSet()
		result.Add(ph.bd.GetGenesisHash())
	} else {
		pastBlueSet := make(map[hash.Hash]*BlockSet)

		ph.calLastCommonBlocksPBS(&pastBlueSet)

		for k, _ := range tips.GetMap() {
			ph.calculatePastBlueSet(&k, &pastBlueSet, false)
		}

		result = ph.calculateBlueSet(tips, nil, &pastBlueSet, false)
	}
	return result
}

func (ph *Phantom) getTempBS() *BlockSet{
	if ph.tempBlueSet==nil {
		ph.tempBlueSet=ph.GetTempBlueSet()
	}
	return ph.tempBlueSet
}

func (ph *Phantom) recCalHourglass(genealogy *BlockSet, ancestors *BlockSet) {

	var maxPastHash *hash.Hash = nil
	var maxPastNum uint = 0

	for k, _ := range ancestors.GetMap() {
		pastNum := ph.GetPastSetNum(ph.bd.GetBlock(&k))
		if maxPastHash == nil || maxPastNum < pastNum {
			maxPastHash = &k
			maxPastNum = pastNum
		}
	}

	if maxPastHash == nil {
		return
	}
	parents := ph.bd.GetBlock(maxPastHash).GetParents()
	if parents == nil || parents.Len() == 0 {
		return
	}
	ancestors.Remove(maxPastHash)
	for k, _ := range parents.GetMap() {
		if !genealogy.Has(&k) {
			ancestors.Add(&k)
			genealogy.Add(&k)
		}
	}

}

func (ph *Phantom) updateHourglass(){
	tips := ph.bd.GetTips()
	if tips == nil||tips.Len()==0 {
		return
	}
	if ph.hourglassBlocks==nil {
		ph.hourglassBlocks=NewBlockSet()
	}
	if tips.HasOnly(ph.bd.GetGenesisHash()){

		ph.hourglassBlocks.Add(ph.bd.GetGenesisHash())
		return
	}
	tempNum:=0
	for k,_:=range tips.GetMap(){
		parents:=ph.bd.GetBlock(&k).GetParents()
		if parents!=nil&&parents.HasOnly(ph.bd.GetGenesisHash()) {
			tempNum++
		}
	}
	if tempNum==tips.Len() {
		return
	}
	//
	genealogy:=NewBlockSet()
	ancestors:=NewBlockSet()

	for k,_:=range tips.GetMap(){
		genealogy.Add(&k)
		ancestors.Add(&k)
	}
	tempBs:=ph.getTempBS()

	for  {
		ph.recCalHourglass(genealogy,ancestors)

		ne0:=tempBs.Intersection(ancestors)
		ne1:=ph.commonBlueSet.Intersection(ancestors)
		ne0.AddSet(ne1)

		ancestors=ne0


		//
		if ancestors.IsEmpty()||ancestors.HasOnly(ph.bd.GetGenesisHash()) {
			ph.hourglassBlocks.Clear()
			ph.hourglassBlocks.Add(ph.bd.GetGenesisHash())
			return
		}

		sb := ph.sortBlockSet(ancestors,nil)
		for _,v:=range sb{
			anti:=ph.GetAnticone(ph.bd.GetBlock(v.h),nil)
			if anti.Len()==0 {
				ph.hourglassBlocks.Exclude(genealogy)
				ph.hourglassBlocks.Add(v.h)
				return
			}else{
				banti0:=tempBs.Intersection(anti)
				banti1:=ph.commonBlueSet.Intersection(anti)
				banti0.AddSet(banti1)

				if banti0.Len()==0 {
					ph.hourglassBlocks.Exclude(genealogy)
					ph.hourglassBlocks.Add(v.h)
					return
				}
			}
		}
	}
}

func (ph *Phantom) updateOrder(b *Block) *list.List{
	ph.tempOrder=[]*hash.Hash{}
	refNodes:=list.New()
	if ph.bd.GetBlockTotal() == 1 {
		ph.tempOrder=append(ph.tempOrder, ph.bd.GetGenesisHash())
		refNodes.PushBack(*ph.bd.GetGenesisHash())
		b.order=0
		return refNodes
	}
	tempOrder := []*hash.Hash{}
	tempOrderM := NewBlockSet()
	//
	blueSet := ph.getTempBS()
	lpsb := ph.sortBlockSet(ph.lastCommonBlocks, nil)
	exclude := NewBlockSet()
	for k, _ := range ph.lastCommonBlocks.GetMap() {
		exclude.AddSet(ph.bd.GetBlock(&k).GetParents())
	}
	for _, v := range lpsb {
		ph.GetTempOrder(&tempOrder, tempOrderM, blueSet, v.h, exclude)
	}
	tLen := len(tempOrder)
	//
	pNum:=ph.GetCommonOrderNum()
	tIndex:=0
	for i := 0; i < tLen; i++ {
		if !ph.lastCommonBlocks.Has(tempOrder[i]) {
			ph.tempOrder = append(ph.tempOrder, tempOrder[i])
			//
			node:=ph.bd.GetBlock(tempOrder[i])

			node.order=uint(pNum+tIndex)
			tIndex++
			if node.order==0 {
				log.Error(fmt.Sprintf("Order error:%v",*node.GetHash()))
			}
		}
	}
	checkOrder:=ph.GetCommonOrderNum()+len(ph.tempOrder)
	if uint(checkOrder)!=ph.bd.GetBlockTotal() {
		log.Error(fmt.Sprintf("Order error:The number is a problem"))
	}
	//////
	tips:=ph.bd.GetTips()
	if tips.HasOnly(b.GetHash())||ph.tempOrder[len(ph.tempOrder)-1].IsEqual(b.GetHash()) {
		b.order=uint(ph.bd.GetBlockTotal()-1)
		refNodes.PushBack(b.GetHash())
		return refNodes
	}
	////
	tLen = len(ph.tempOrder)
	for i:=tLen-1;i>=0;i-- {
		refNodes.PushFront(ph.tempOrder[i])
		if ph.tempOrder[i].IsEqual(b.GetHash()) {
			break
		}
	}
	return refNodes
}

func (ph *Phantom) GetCommonOrderNum() int{
	pLen:=len(ph.commonOrder)

	if pLen>0 {
		var i int
		for i=pLen-1;i>=0 ;i--  {
			if ph.commonOrder[i]!=nil {
				break
			}
		}
		return i+1
	}
	return 0
}
/*
func (ph *Phantom) GetBlockOrder(h *hash.Hash) int32{
	var result int32=-1
	if ph.tempOrder==nil {
		return result
	}
	result=int32(ph.bd.GetBlockTotal())
	tLen:=len(ph.tempOrder)
	if tLen>0 {
		for i:=tLen-1;i>=0 ;i--  {
			if ph.tempOrder[i]!=nil {
				result--
				if h.IsEqual(ph.tempOrder[i]) {
					return result
				}
			}
		}
	}
	pLen:=len(ph.commonOrder)
	if pLen>0 {
		for i:=pLen-1;i>=0 ;i--  {
			if ph.commonOrder[i]!=nil {
				result--
				if h.IsEqual(ph.commonOrder[i]) {
					return result
				}
			}
		}
	}

	return -1
}
*/


func (ph *Phantom) GetBlockByOrder(order uint) *hash.Hash{
	if ph.tempOrder==nil||order<0 {
		return nil
	}
	pNum:=uint(ph.GetCommonOrderNum())
	if order<pNum {
		return ph.commonOrder[order]
	}
	rIndex:=order-pNum
	tLen:=len(ph.tempOrder)
	if rIndex<uint(tLen) {
		return ph.tempOrder[rIndex]
	}
	return nil
}

func (ph *Phantom) GetTipsList() []*Block {
	return nil
}
///////
type SortBlock struct {
	h          *hash.Hash
	pastSetNum uint
}

type SortBlocks []SortBlock

func (a SortBlocks) Len() int {
	return len(a)
}

func (a SortBlocks) Less(i, j int) bool {
	if a[i].pastSetNum == a[j].pastSetNum {
		return a[i].h.String() < a[j].h.String()
	}
	return a[i].pastSetNum < a[j].pastSetNum
}

func (a SortBlocks) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
/////////
