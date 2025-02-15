package xiami

import (
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/mxget/pkg/concurrency"
	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

const (
	APISearch            = "https://acs.m.xiami.com/h5/mtop.alimusic.search.searchservice.searchsongs/1.0/?appKey=23649156"
	APIGetSongDetail     = "https://acs.m.xiami.com/h5/mtop.alimusic.music.songservice.getsongdetail/1.0/?appKey=23649156"
	APIGetSongLyric      = "https://acs.m.xiami.com/h5/mtop.alimusic.music.lyricservice.getsonglyrics/1.0/?appKey=23649156"
	APIGetSongs          = "https://acs.m.xiami.com/h5/mtop.alimusic.music.songservice.getsongs/1.0/?appKey=23649156"
	APIGetArtistInfo     = "https://acs.m.xiami.com/h5/mtop.alimusic.music.artistservice.getartistdetail/1.0/?appKey=23649156"
	APIGetArtistSongs    = "https://acs.m.xiami.com/h5/mtop.alimusic.music.songservice.getartistsongs/1.0/?appKey=23649156"
	APIGetAlbum          = "https://acs.m.xiami.com/h5/mtop.alimusic.music.albumservice.getalbumdetail/1.0/?appKey=23649156"
	APIGetPlaylistDetail = "https://h5api.m.xiami.com/h5/mtop.alimusic.music.list.collectservice.getcollectdetail/1.0/?appKey=23649156"
	APIGetPlaylistSongs  = "https://h5api.m.xiami.com/h5/mtop.alimusic.music.list.collectservice.getcollectsongs/1.0/?appKey=23649156"
	APILogin             = "https://h5api.m.xiami.com/h5/mtop.alimusic.xuser.facade.xiamiuserservice.login/1.0/?appKey=23649156"

	SongRequestLimit = 200
)

var (
	std = New(provider.Client())
)

type (
	CommonResponse struct {
		API string   `json:"api"`
		Ret []string `json:"ret"`
	}

	PagingVO struct {
		Count    string `json:"count"`
		Page     string `json:"page"`
		PageSize string `json:"pageSize"`
		Pages    string `json:"pages"`
	}

	ListenFile struct {
		Expire     string `json:"expire,omitempty"`
		FileSize   string `json:"fileSize"`
		Format     string `json:"format"`
		ListenFile string `json:"listenFile,omitempty"`
		Quality    string `json:"quality"`
		URL        string `json:"url,omitempty"`
	}

	Song struct {
		Album
		SongId       string       `json:"songId"`
		SongStringId string       `json:"songStringId"`
		SongName     string       `json:"songName"`
		Singers      string       `json:"singers"`
		SingerVOs    []Artist     `json:"singerVOs"`
		ListenFile   string       `json:"listenFile,omitempty"`
		ListenFiles  []ListenFile `json:"listenFiles"`
		Lyric        string       `json:"-"`
	}

	SearchSongsResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				PagingVO PagingVO `json:"pagingVO"`
				Songs    []*Song  `json:"songs"`
			} `json:"data"`
		} `json:"data"`
	}

	SongDetailResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				SongDetail Song `json:"songDetail"`
			} `json:"data"`
		} `json:"data"`
	}

	SongLyricResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				Lyrics []struct {
					Content      string `json:"content"`
					FlagOfficial string `json:"flagOfficial"`
					LyricURL     string `json:"lyricUrl"`
					Type         string `json:"type"`
				} `json:"lyrics"`
			} `json:"data"`
		} `json:"data"`
	}

	SongsResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				Songs []*Song `json:"songs"`
			} `json:"data"`
		} `json:"data"`
	}

	Artist struct {
		ArtistId       string `json:"artistId"`
		ArtistStringId string `json:"artistStringId"`
		ArtistName     string `json:"artistName"`
		ArtistLogo     string `json:"artistLogo"`
	}

	ArtistInfoResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				ArtistDetailVO Artist `json:"artistDetailVO"`
			} `json:"data"`
		} `json:"data"`
	}

	ArtistSongsResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				PagingVO PagingVO `json:"pagingVO"`
				Songs    []*Song  `json:"songs"`
			} `json:"data"`
		} `json:"data"`
	}

	Album struct {
		AlbumId       string `json:"albumId"`
		AlbumStringId string `json:"albumStringId"`
		AlbumName     string `json:"albumName"`
		AlbumLogo     string `json:"albumLogo"`
	}

	AlbumResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				AlbumDetail struct {
					Album
					Songs []*Song `json:"songs"`
				} `json:"albumDetail"`
			} `json:"data"`
		} `json:"data"`
	}

	PlaylistDetailResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				CollectDetail struct {
					ListId      string   `json:"listId"`
					CollectName string   `json:"collectName"`
					CollectLogo string   `json:"collectLogo"`
					SongCount   string   `json:"songCount"`
					AllSongs    []string `json:"allSongs"`
					Songs       []*Song  `json:"songs"`
					PagingVO    PagingVO `json:"pagingVO"`
				} `json:"collectDetail"`
			} `json:"data"`
		} `json:"data"`
	}

	PlaylistSongsResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				Songs    []*Song  `json:"songs"`
				PagingVO PagingVO `json:"pagingVO"`
			} `json:"data"`
		} `json:"data"`
	}

	LoginResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				AccessToken    string `json:"accessToken"`
				Expires        string `json:"expires"`
				NickName       string `json:"nickName"`
				RefreshExpires string `json:"refreshExpires"`
				RefreshToken   string `json:"refreshToken"`
				UserId         string `json:"userId"`
			} `json:"data"`
		} `json:"data"`
	}

	API struct {
		Client *sreq.Client
	}
)

func (s *SearchSongsResponse) String() string {
	return provider.ToJSON(s, false)
}

func (s *SongDetailResponse) String() string {
	return provider.ToJSON(s, false)
}

func (s *SongsResponse) String() string {
	return provider.ToJSON(s, false)
}

func (a *ArtistInfoResponse) String() string {
	return provider.ToJSON(a, false)
}

func (a *ArtistSongsResponse) String() string {
	return provider.ToJSON(a, false)
}

func (a *AlbumResponse) String() string {
	return provider.ToJSON(a, false)
}

func (p *PlaylistDetailResponse) String() string {
	return provider.ToJSON(p, false)
}

func (p *PlaylistSongsResponse) String() string {
	return provider.ToJSON(p, false)
}

func (p *LoginResponse) String() string {
	return provider.ToJSON(p, false)
}

func New(client *sreq.Client) *API {
	if client == nil {
		client = sreq.New(nil)
		client.SetDefaultRequestOpts(
			sreq.WithHeaders(sreq.Headers{
				"User-Agent": provider.UserAgent,
			}),
		)
	}
	return &API{
		Client: client,
	}
}

func Client() provider.API {
	return std
}

func (a *API) Platform() int {
	return provider.XiaMi
}

func (a *API) Request(method string, url string, opts ...sreq.RequestOption) *sreq.Response {
	defaultOpts := []sreq.RequestOption{
		sreq.WithHeaders(sreq.Headers{
			"Origin":  "https://h.xiami.com",
			"Referer": "https://h.xiami.com",
		}),
	}
	opts = append(opts, defaultOpts...)
	return a.Client.Request(method, url, opts...)
}

func (a *API) getToken(url string) (string, error) {
	const XiaMiToken = "_m_h5_tk"
	token, err := a.Client.FilterCookie(url, XiaMiToken)
	if err != nil {
		// 如果在cookie jar中没有找到对应cookie，发送预请求获取
		token, err = a.Request(sreq.MethodGet, url).Cookie(XiaMiToken)
	}

	if err != nil {
		return "", fmt.Errorf("can't get token: %w", err)
	}

	return strings.Split(token.Value, "_")[0], nil
}

func (c *CommonResponse) check() error {
	for _, s := range c.Ret {
		if strings.HasPrefix(s, "FAIL") {
			return errors.New(s)
		}
	}
	return nil
}

func songURL(listenFiles []ListenFile) string {
	for _, i := range listenFiles {
		if i.Quality == "l" {
			return i.URL + i.ListenFile
		}
	}
	return ""
}

func (a *API) patchSongLyric(songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		c.Add(1)
		go func(s *Song) {
			lyric, err := a.GetSongLyric(s.SongId)
			if err == nil {
				s.Lyric = lyric
			}
			c.Done()
		}(s)
	}
	c.Wait()
}

func resolve(src ...*Song) []*provider.Song {
	songs := make([]*provider.Song, 0, len(src))
	for _, s := range src {
		url := songURL(s.ListenFiles)
		songs = append(songs, &provider.Song{
			Name:     strings.TrimSpace(s.SongName),
			Artist:   strings.TrimSpace(strings.ReplaceAll(s.Singers, " / ", "/")),
			Album:    strings.TrimSpace(s.AlbumName),
			PicURL:   s.AlbumLogo,
			Lyric:    s.Lyric,
			Playable: url != "",
			URL:      url,
		})
	}
	return songs
}
