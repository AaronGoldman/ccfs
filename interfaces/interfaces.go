package interfaces

import ()

const LocalSeed = "c09b2765c6fd4b999d47c82f9cdf7f4b659bf7c29487cc0b357b8fc92ac8ad02"

//GetLocalSeed provides the root HKID for the local environment
func GetLocalSeed() string {
	return LocalSeed
}
