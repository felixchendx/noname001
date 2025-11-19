package route

import (
	"fmt"
)

// temp experimental

// TODO: reassess these stuffs when designing web-wide links
// TODO: move this ?, to /web/module/stream/constant/...

const (
	navLink__blank              = ""
	navLink__streamGroupListing = "/stream/stream-group/listing"
	navLink__streamGroupDetail  = "/stream/stream-group/detail-02"

	actLink__streamGroupDoAdd    = "/stream/stream-group/detail-02/do/add"
	actLink__streamGroupDoEdit   = "/stream/stream-group/detail-02/do/edit"
	actLink__streamGroupDoDelete = "/stream/stream-group/detail-02/do/delete"
)

func navLink__streamGroupDetailWithID(id string) (string) {
	return navLink__streamGroupDetail + fmt.Sprintf("?id=%s", id)
}

var (
	navLinks map[string]string = map[string]string{
		"listing": navLink__streamGroupListing,
	}

	// navLinks navLinks = navLinks{
	// 	listing: navLink__streamGroupListing,
	// }
)

// type navLinks struct {
// 	listing string
// }
