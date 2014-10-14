//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

//
package interfaces

import (
	"log"
	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

//GetLocalSeed provides the root HKID for the local enviernment
func GetLocalSeed() string {
	return "c09b2765c6fd4b999d47c82f9cdf7f4b659bf7c29487cc0b357b8fc92ac8ad02"
}

func KeyLocalSeed() {
	HKIDstring := GetLocalSeed()
	h, er := objects.HkidFromHex(HKIDstring)
	if (er!=nil){
		log.Printf("local seed not valid hex /n")
	}
	_, err := services.GetKey(h)
	if (err!=nil){
		 objects.HkidFromDString("65232373562705602286177837897283294165955126"+
		"49112249373497830592072241416893611216069423804730437860475300564272"+
		"976762085068519188612732562106886379081213385", 10)
		}
	return
}


