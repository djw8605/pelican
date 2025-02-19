/***************************************************************
 *
 * Copyright (C) 2023, Pelican Project, Morgridge Institute for Research
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you
 * may not use this file except in compliance with the License.  You may
 * obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 ***************************************************************/

package director

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func hasServerAdWithName(serverAds []ServerAd, name string) bool {
	for _, serverAd := range serverAds {
		if serverAd.Name == name {
			return true
		}
	}
	return false
}

// Test getAdsForPath to make sure various nuanced cases work. Under the hood
// this really tests matchesPrefix, but we test this higher level function to
// avoid having to mess with the cache.
func TestGetAdsForPath(t *testing.T) {
	/*
		FLOW:
			- Set up a few dummy namespaces, origin, and cache ads
			- Record the ads
			- Query for a few paths and make sure the correct ads are returned
	*/
	nsAd1 := NamespaceAd{
		RequireToken: true,
		Path:         "/chtc",
		Issuer: url.URL{
			Scheme: "https",
			Host:   "wisc.edu",
		},
	}

	nsAd2 := NamespaceAd{
		RequireToken: false,
		Path:         "/chtc/PUBLIC",
		Issuer: url.URL{
			Scheme: "https",
			Host:   "wisc.edu",
		},
	}

	nsAd3 := NamespaceAd{
		RequireToken: false,
		Path:         "/chtc/PUBLIC2/",
		Issuer: url.URL{
			Scheme: "https",
			Host:   "wisc.edu",
		},
	}

	cacheAd1 := ServerAd{
		Name: "cache1",
		AuthURL: url.URL{
			Scheme: "https",
			Host:   "wisc.edu",
		},
		URL: url.URL{
			Scheme: "https",
			Host:   "wisc.edu",
		},
		Type: CacheType,
	}

	cacheAd2 := ServerAd{
		Name: "cache2",
		AuthURL: url.URL{
			Scheme: "https",
			Host:   "wisc.edu",
		},
		URL: url.URL{
			Scheme: "https",
			Host:   "wisc.edu",
		},
		Type: CacheType,
	}

	originAd1 := ServerAd{
		Name: "origin1",
		AuthURL: url.URL{
			Scheme: "https",
			Host:   "wisc.edu",
		},
		URL: url.URL{
			Scheme: "https",
			Host:   "wisc.edu",
		},
		Type: OriginType,
	}

	originAd2 := ServerAd{
		Name: "origin2",
		AuthURL: url.URL{
			Scheme: "https",
			Host:   "wisc.edu",
		},
		URL: url.URL{
			Scheme: "https",
			Host:   "wisc.edu",
		},
		Type: OriginType,
	}

	o1Slice := []NamespaceAd{nsAd1}
	o2Slice := []NamespaceAd{nsAd2, nsAd3}
	c1Slice := []NamespaceAd{nsAd1, nsAd2}
	RecordAd(originAd2, &o2Slice)
	RecordAd(originAd1, &o1Slice)
	RecordAd(cacheAd1, &c1Slice)
	RecordAd(cacheAd2, &o1Slice)

	nsAd, oAds, cAds := GetAdsForPath("/chtc")
	assert.Equal(t, nsAd.Path, "/chtc")
	assert.Equal(t, len(oAds), 1)
	assert.Equal(t, len(cAds), 2)
	assert.True(t, hasServerAdWithName(oAds, "origin1"))
	assert.True(t, hasServerAdWithName(cAds, "cache1"))
	assert.True(t, hasServerAdWithName(cAds, "cache2"))

	nsAd, oAds, cAds = GetAdsForPath("/chtc/")
	assert.Equal(t, nsAd.Path, "/chtc")
	assert.Equal(t, len(oAds), 1)
	assert.Equal(t, len(cAds), 2)
	assert.True(t, hasServerAdWithName(oAds, "origin1"))
	assert.True(t, hasServerAdWithName(cAds, "cache1"))
	assert.True(t, hasServerAdWithName(cAds, "cache2"))

	nsAd, oAds, cAds = GetAdsForPath("/chtc/PUBLI")
	assert.Equal(t, nsAd.Path, "/chtc")
	assert.Equal(t, len(oAds), 1)
	assert.Equal(t, len(cAds), 2)
	assert.True(t, hasServerAdWithName(oAds, "origin1"))
	assert.True(t, hasServerAdWithName(cAds, "cache1"))
	assert.True(t, hasServerAdWithName(cAds, "cache2"))

	nsAd, oAds, cAds = GetAdsForPath("/chtc/PUBLIC")
	assert.Equal(t, nsAd.Path, "/chtc/PUBLIC")
	assert.Equal(t, len(oAds), 1)
	assert.Equal(t, len(cAds), 1)
	assert.True(t, hasServerAdWithName(oAds, "origin2"))
	assert.True(t, hasServerAdWithName(cAds, "cache1"))

	nsAd, oAds, cAds = GetAdsForPath("/chtc/PUBLIC2")
	// since the stored path is actually /chtc/PUBLIC2/, the extra / is returned
	assert.Equal(t, nsAd.Path, "/chtc/PUBLIC2/")
	assert.Equal(t, len(oAds), 1)
	assert.Equal(t, len(cAds), 0)
	assert.True(t, hasServerAdWithName(oAds, "origin2"))
}
