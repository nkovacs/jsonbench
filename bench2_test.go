package main

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/nkovacs/jsonx"
)

type CBAvatar struct {
	Url string `json:"url"`
}

type CBGravatar struct {
	Avatars []*CBAvatar `json:"avatars"`
}

type CBGithub struct {
	Followers int `json:"followers"`
}

type CBName struct {
	FullName string `json:"fullName"`
}

type CBPerson struct {
	Name     *CBName     `json:"name"`
	Github   *CBGithub   `json:"github"`
	Gravatar *CBGravatar `json:"gravatar"`
}

type MediumPayload struct {
	Person  *CBPerson `json:"person"`
	Company string    `json:"compnay"`
}

var mediumFixture = []byte(`{
	"person": {
	  "id": "d50887ca-a6ce-4e59-b89f-14f0b5d03b03",
	  "name": {
		"fullName": "Leonid Bugaev",
		"givenName": "Leonid",
		"familyName": "Bugaev"
	  },
	  "email": "leonsbox@gmail.com",
	  "gender": "male",
	  "location": "Saint Petersburg, Saint Petersburg, RU",
	  "geo": {
		"city": "Saint Petersburg",
		"state": "Saint Petersburg",
		"country": "Russia",
		"lat": 59.9342802,
		"lng": 30.3350986
	  },
	  "bio": "Senior engineer at Granify.com",
	  "site": "http://flickfaver.com",
	  "avatar": "https://d1ts43dypk8bqh.cloudfront.net/v1/avatars/d50887ca-a6ce-4e59-b89f-14f0b5d03b03",
	  "employment": {
		"name": "www.latera.ru",
		"title": "Software Engineer",
		"domain": "gmail.com"
	  },
	  "facebook": {
		"handle": "leonid.bugaev"
	  },
	  "github": {
		"handle": "buger",
		"id": 14009,
		"avatar": "https://avatars.githubusercontent.com/u/14009?v=3",
		"company": "Granify",
		"blog": "http://leonsbox.com",
		"followers": 95,
		"following": 10
	  },
	  "twitter": {
		"handle": "flickfaver",
		"id": 77004410,
		"bio": null,
		"followers": 2,
		"following": 1,
		"statuses": 5,
		"favorites": 0,
		"location": "",
		"site": "http://flickfaver.com",
		"avatar": null
	  },
	  "linkedin": {
		"handle": "in/leonidbugaev"
	  },
	  "googleplus": {
		"handle": null
	  },
	  "angellist": {
		"handle": "leonid-bugaev",
		"id": 61541,
		"bio": "Senior engineer at Granify.com",
		"blog": "http://buger.github.com",
		"site": "http://buger.github.com",
		"followers": 41,
		"avatar": "https://d1qb2nb5cznatu.cloudfront.net/users/61541-medium_jpg?1405474390"
	  },
	  "klout": {
		"handle": null,
		"score": null
	  },
	  "foursquare": {
		"handle": null
	  },
	  "aboutme": {
		"handle": "leonid.bugaev",
		"bio": null,
		"avatar": null
	  },
	  "gravatar": {
		"handle": "buger",
		"urls": [
		],
		"avatar": "http://1.gravatar.com/avatar/f7c8edd577d13b8930d5522f28123510",
		"avatars": [
		  {
			"url": "http://1.gravatar.com/avatar/f7c8edd577d13b8930d5522f28123510",
			"type": "thumbnail"
		  }
		]
	  },
	  "fuzzy": false
	},
	"company": null
}`)

func BenchmarkJsoniterMediumPayloadDecode(b *testing.B) {
	var data MediumPayload
	b.Run("standard", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err := json.Unmarshal(mediumFixture, &data)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("jsonx", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err := jsonx.Unmarshal(mediumFixture, &data)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("jsoniter", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err := jsoniterStd.Unmarshal(mediumFixture, &data)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkJsoniterMediumPayloadEncode(b *testing.B) {
	var data MediumPayload
	err := json.Unmarshal(mediumFixture, &data)
	if err != nil {
		b.Fatal(err)
	}
	b.Run("standard", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			bts, err := json.Marshal(data)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(len(bts)))
		}
	})
	b.Run("jsonx", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			bts, err := jsonx.Marshal(data)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(len(bts)))
		}
	})
	b.Run("jsonx encoder", func(b *testing.B) {
		var buf bytes.Buffer
		enc := jsonx.NewEncoder(&buf)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err := enc.Encode(data)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(buf.Len()))
			buf.Reset()
		}
	})
	b.Run("jsoniter", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			bts, err := jsoniterStd.Marshal(data)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(len(bts)))
		}
	})
}
