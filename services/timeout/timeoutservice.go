package timeout

import (
	"fmt"
	"time"

	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

type timeoutservice struct{}

func (timeoutservice) GetBlob(objects.HCID) (objects.Blob, error) {
	time.Sleep(time.Second)
	return objects.Blob{}, fmt.Errorf("GetBlob Timeout")
}
func (timeoutservice) GetCommit(objects.HKID) (objects.Commit, error) {
	time.Sleep(time.Second)
	return objects.Commit{}, fmt.Errorf("GetCommit Timeout")
}
func (timeoutservice) GetTag(h objects.HKID, namesegment string) (objects.Tag, error) {
	time.Sleep(time.Second)
	return objects.Tag{}, fmt.Errorf("GetTag Timeout")
}

func (timeoutservice) GetTags(h objects.HKID) ([]objects.Tag, error) {
	time.Sleep(time.Second)
	return []objects.Tag{}, fmt.Errorf("GetTags Timeout")
}

func (timeoutservice) GetKey(objects.HKID) (objects.Blob, error) {
	time.Sleep(time.Second)
	return objects.Blob{}, fmt.Errorf("GetKey Timeout")
}

var Instance timeoutservice

func init() {
	services.Registerblobgeter(Instance)
	services.Registercommitgeter(Instance)
	services.Registertaggeter(Instance)
	services.Registerkeygeter(Instance)
}
