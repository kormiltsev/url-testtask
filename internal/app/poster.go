package app

import (
	st "github.com/kormiltsev/url-testtask/internal/storage"
)

// find duplicates in catalog
func Post(data st.Request) (st.Request, string) {
	if data.Url != "" {
		if Conf.DBtype == "local" {
			// get short url by url requested from catalog
			surl, found := FindUrl(Catalog, data.Url)
			if found {
				data.Surl = surl
				return data, "done. already exists"
			} else {
				return PostUniq(data)
			}
		} else {
			// connect to postgres
			req, comment := st.PostPostgres(data.Url)
			data = req
			return data, comment
		}

	} else {
		return data, "requested url is empty"
	}
}

// create short url and check for duplicates
func PostUniq(data st.Request) (st.Request, string) {
	// get random string
	surl := GetRandomStringFaster(Conf.Settings.Qty, Conf.Settings.Letters)
	if _, found := FindSurl(Catalog, surl); found {
		// repeat till create uniq short url
		return PostUniq(data)
	} else {
		data.Surl = surl
		AddLink(data)
		return data, "done. new short_url"
	}
}

// add new url and short url in catalog
func AddLink(data st.Request) {
	Catalog.List = append(Catalog.List, data)
}
