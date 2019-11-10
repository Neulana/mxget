package baidu

import (
	"strings"

	"github.com/winterssy/mxget/pkg/concurrency"
	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

const (
	APISearch       = "http://musicapi.qianqian.com/v1/restserver/ting?method=baidu.ting.search.merge&from=android&version=8.1.4.0&format=json&type=-1&isNew=1"
	APIGetSong      = "http://musicapi.qianqian.com/v1/restserver/ting?method=baidu.ting.song.getInfos&format=json&from=android&version=8.1.4.0"
	APIGetSongs     = "http://music.taihe.com/data/music/fmlink"
	APIGetSongLyric = "http://musicapi.qianqian.com/v1/restserver/ting?method=baidu.ting.song.lry&format=json&from=android&version=8.1.4.0"
	APIGetArtist    = "http://musicapi.qianqian.com/v1/restserver/ting?method=baidu.ting.artist.getSongList&from=android&version=8.1.4.0&format=json&order=2"
	APIGetAlbum     = "http://musicapi.qianqian.com/v1/restserver/ting?method=baidu.ting.album.getAlbumInfo&from=android&version=8.1.4.0&format=json"
	APIGetPlaylist  = "http://musicapi.qianqian.com/v1/restserver/ting?method=baidu.ting.ugcdiy.getBaseInfo&from=android&version=8.1.4.0"
)

var (
	std = New(provider.Client())
)

type (
	CommonResponse struct {
		ErrorCode    int    `json:"error_code,omitempty"`
		ErrorMessage string `json:"error_message,omitempty"`
	}

	Song struct {
		SongId     string `json:"song_id"`
		Title      string `json:"title"`
		TingUid    string `json:"ting_uid"`
		Author     string `json:"author"`
		AlbumId    string `json:"album_id"`
		AlbumTitle string `json:"album_title"`
		PicBig     string `json:"pic_big"`
		LrcLink    string `json:"lrclink"`
		CopyType   string `json:"copy_type"`
		URL        string `json:"-"`
		Lyric      string `json:"-"`
	}

	SearchSongsResponse struct {
		CommonResponse
		Result struct {
			SongInfo struct {
				SongList []*Song `json:"song_list"`
			} `json:"song_info"`
		} `json:"result"`
	}

	URL struct {
		ShowLink    string `json:"show_link"`
		FileFormat  string `json:"file_format"`
		FileBitrate int    `json:"file_bitrate"`
		FileLink    string `json:"file_link"`
	}

	SongResponse struct {
		CommonResponse
		SongInfo Song `json:"songinfo"`
		SongURL  struct {
			URL []URL `json:"url"`
		} `json:"songurl"`
	}

	SongsResponse struct {
		ErrorCode int `json:"errorCode"`
		Data      struct {
			SongList []*struct {
				SongId     int    `json:"songId"`
				SongName   string `json:"songName"`
				ArtistId   string `json:"artistId"`
				ArtistName string `json:"artistName"`
				AlbumId    int    `json:"albumId"`
				AlbumName  string `json:"albumName"`
				SongPicBig string `json:"songPicBig"`
				LrcLink    string `json:"lrcLink"`
				CopyType   int    `json:"copyType"`
				SongLink   string `json:"songLink"`
				ShowLink   string `json:"showLink"`
				Format     string `json:"format"`
				Rate       int    `json:"rate"`
			} `json:"songList"`
		} `json:"data"`
	}

	SongLyricResponse struct {
		CommonResponse
		Title      string `json:"title"`
		LrcContent string `json:"lrcContent"`
	}

	ArtistResponse struct {
		CommonResponse
		ArtistInfo struct {
			TingUid   string `json:"ting_uid"`
			Name      string `json:"name"`
			AvatarBig string `json:"avatar_big"`
		} `json:"artistinfo"`
		SongList []*Song `json:"songlist"`
	}

	AlbumResponse struct {
		CommonResponse
		AlbumInfo struct {
			AlbumId string `json:"album_id"`
			Title   string `json:"title"`
			PicBig  string `json:"pic_big"`
		} `json:"albuminfo"`
		SongList []*Song `json:"songlist"`
	}

	PlaylistResponse struct {
		CommonResponse
		Result struct {
			Info struct {
				ListId    string `json:"list_id"`
				ListTitle string `json:"list_title"`
				ListPic   string `json:"list_pic"`
			} `json:"info"`
			SongList []*Song `json:"songlist"`
		}
	}

	API struct {
		Client *sreq.Client
	}
)

func (s *SearchSongsResponse) String() string {
	return provider.ToJSON(s, false)
}

func (s *SongResponse) String() string {
	return provider.ToJSON(s, false)
}

func (s *SongsResponse) String() string {
	return provider.ToJSON(s, false)
}

func (s *SongLyricResponse) String() string {
	return provider.ToJSON(s, false)
}

func (a *ArtistResponse) String() string {
	return provider.ToJSON(a, false)
}

func (a *AlbumResponse) String() string {
	return provider.ToJSON(a, false)
}

func (p *PlaylistResponse) String() string {
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
	return provider.BaiDu
}

func (a *API) Request(method string, url string, opts ...sreq.RequestOption) *sreq.Response {
	defaultOpts := []sreq.RequestOption{
		sreq.WithHeaders(sreq.Headers{
			"Origin":  "http://music.taihe.com",
			"Referer": "http://music.taihe.com",
		}),
	}
	opts = append(opts, defaultOpts...)
	return a.Client.Request(method, url, opts...)
}

func resolve(src ...*Song) []*provider.Song {
	songs := make([]*provider.Song, 0, len(src))
	for _, s := range src {
		songs = append(songs, &provider.Song{
			Name:     strings.TrimSpace(s.Title),
			Artist:   strings.TrimSpace(strings.ReplaceAll(s.Author, ",", "/")),
			Album:    strings.TrimSpace(s.AlbumTitle),
			PicURL:   strings.Split(s.PicBig, "@")[0],
			Lyric:    s.Lyric,
			Playable: s.URL != "",
			URL:      s.URL,
		})
	}
	return songs
}

func songURL(urls []URL) string {
	for _, i := range urls {
		if i.FileFormat == "mp3" {
			return i.ShowLink
		}
	}
	return ""
}

func (a *API) patchSongURL(songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		c.Add(1)
		go func(s *Song) {
			resp, err := a.GetSongRaw(s.SongId)
			if err == nil {
				s.URL = songURL(resp.SongURL.URL)
				if s.LrcLink == "" && resp.SongInfo.LrcLink != "" {
					s.LrcLink = resp.SongInfo.LrcLink
				}
			}
			c.Done()
		}(s)
	}
	c.Wait()
}

func (a *API) patchSongLyric(songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		c.Add(1)
		go func(s *Song) {
			lyric, err := a.Request(sreq.MethodGet, s.LrcLink).Text()
			if err == nil {
				s.Lyric = lyric
			}
			c.Done()
		}(s)
	}
	c.Wait()
}
