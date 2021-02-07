package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"sync"
)

var (
	mEstate         sync.RWMutex
	esCountByWidth  [][]Estate
	esCountByHeight [][]Estate
	esCountByRent   [][]Estate
	estateByID      = make(map[int64]Estate)
)

func appendCountCache(es Estate) {
	switch {
	case es.DoorWidth < 80:
		esCountByWidth[0] = append(esCountByWidth[0], es)
	case es.DoorWidth < 110:
		esCountByWidth[1] = append(esCountByWidth[1], es)
	case es.DoorWidth < 150:
		esCountByWidth[2] = append(esCountByWidth[2], es)
	default:
		esCountByWidth[3] = append(esCountByWidth[3], es)
	}

	switch {
	case es.DoorHeight < 80:
		esCountByHeight[0] = append(esCountByHeight[0], es)
	case es.DoorHeight < 110:
		esCountByHeight[1] = append(esCountByHeight[1], es)
	case es.DoorHeight < 150:
		esCountByHeight[2] = append(esCountByHeight[2], es)
	default:
		esCountByHeight[3] = append(esCountByHeight[3], es)
	}

	switch {
	case es.Rent < 50000:
		esCountByRent[0] = append(esCountByRent[0], es)
	case es.Rent < 100000:
		esCountByRent[1] = append(esCountByRent[1], es)
	case es.Rent < 150000:
		esCountByRent[2] = append(esCountByRent[2], es)
	default:
		esCountByRent[3] = append(esCountByRent[3], es)
	}
}

// 降順そーと
type byPopularity []Estate

func (a byPopularity) Len() int { return len(a) }
func (a byPopularity) Less(i, j int) bool {
	if a[i].Popularity == a[j].Popularity {
		return a[i].ID < a[j].ID
	}
	return a[i].Popularity > a[j].Popularity
}
func (a byPopularity) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func sortCountCache() {
	sort.Sort(byPopularity(esCountByWidth[0]))
	sort.Sort(byPopularity(esCountByWidth[1]))
	sort.Sort(byPopularity(esCountByWidth[2]))
	sort.Sort(byPopularity(esCountByWidth[3]))

	sort.Sort(byPopularity(esCountByHeight[0]))
	sort.Sort(byPopularity(esCountByHeight[1]))
	sort.Sort(byPopularity(esCountByHeight[2]))
	sort.Sort(byPopularity(esCountByHeight[3]))

	sort.Sort(byPopularity(esCountByRent[0]))
	sort.Sort(byPopularity(esCountByRent[1]))
	sort.Sort(byPopularity(esCountByRent[2]))
	sort.Sort(byPopularity(esCountByRent[3]))

	//for i := 0; i < 4; i++ {
	//	fmt.Printf("width[%v]=%v\n", i, len(esCountByWidth[i]))
	//}
	//for i := 0; i < 4; i++ {
	//	fmt.Printf("height[%v]=%v\n", i, len(esCountByHeight[i]))
	//}
	//for i := 0; i < 4; i++ {
	//	fmt.Printf("rent[%v]=%v\n", i, len(esCountByRent[i]))
	//}
}

func initEstateCache() {
	mEstate.Lock()
	defer mEstate.Unlock()

	estateByID = make(map[int64]Estate)

	var estateCache []Estate
	err := db.Select(&estateCache, "SELECT * FROM estate")
	if err != nil {
		log.Printf("ERROR!!! failed to load estate: %v", err)
		panic(err)
	}

	fmt.Printf("estate total=%v\n", len(estateCache))

	esCountByWidth = make([][]Estate, 4)
	esCountByHeight = make([][]Estate, 4)
	esCountByRent = make([][]Estate, 4)

	for _, es := range estateCache {
		appendCountCache(es)
		estateByID[es.ID] = es
	}
	sortCountCache()

	//fmt.Println("countbywidth")
	//for i, s := range esCountByWidth {
	//	for j, t := range s {
	//		fmt.Printf("%v,%v, w=%v p=%v\n", i, j, t.DoorWidth, t.Popularity)
	//	}
	//}
	//fmt.Println("countbyheight")
	//for i, s := range esCountByHeight {
	//	for j, t := range s {
	//		fmt.Printf("%v,%v, h=%v p=%v\n", i, j, t.DoorHeight, t.Popularity)
	//	}
	//}
	//fmt.Println("countbyrent")
	//for i, s := range esCountByRent {
	//	for j, t := range s {
	//		fmt.Printf("%v,%v, r=%v p=%v\n", i, j, t.Rent, t.Popularity)
	//	}
	//}
}

func searchEstateByWidth(widthID int) []Estate {
	mEstate.RLock()
	defer mEstate.RUnlock()

	if widthID >= len(esCountByWidth) {
		panic("bad width id")
	}

	res := make([]Estate, len(esCountByWidth[widthID]))
	copy(res, esCountByWidth[widthID])
	return res
}

func searchEstateByHeight(heightID int) []Estate {
	mEstate.RLock()
	defer mEstate.RUnlock()

	if heightID >= len(esCountByHeight) {
		panic("bad height id")
	}

	res := make([]Estate, len(esCountByHeight[heightID]))
	copy(res, esCountByHeight[heightID])
	return res
}

func searchEstateByRent(id int) []Estate {
	mEstate.RLock()
	defer mEstate.RUnlock()

	if id >= len(esCountByRent) {
		panic("bad rent id")
	}

	res := make([]Estate, len(esCountByRent[id]))
	copy(res, esCountByRent[id])
	return res
}

func getEstateByID(id int64) (Estate, error) {
	mEstate.RLock()
	defer mEstate.RUnlock()

	es, ok := estateByID[id]
	if !ok {
		return es, sql.ErrNoRows
	}
	return es, nil
}
