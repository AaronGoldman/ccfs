package main

import ()

type entry struct {
	Hash        []byte
	TypeString  string
	nameSegment string
}

type list []entry
