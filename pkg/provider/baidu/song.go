package baidu

import (
	"fmt"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func GetSong(songId string) (*provider.Song, error) {
	return std.GetSong(songId)
}

func (a *API) GetSong(songId string) (*provider.Song, error) {
	resp, err := a.GetSongRaw(songId)
	if err != nil {
		return nil, err
	}

	resp.SongInfo.URL = songURL(resp.SongURL.URL)
	a.patchSongLyric(&resp.SongInfo)
	songs := resolve(&resp.SongInfo)
	return songs[0], nil
}

func GetSongRaw(songId string) (*SongResponse, error) {
	return std.GetSongRaw(songId)
}

// 获取歌曲
func (a *API) GetSongRaw(songId string) (*SongResponse, error) {
	resp := new(SongResponse)
	err := a.Request(sreq.MethodGet, APIGetSong, sreq.WithQuery(aesCBCEncrypt(songId))).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrorCode != 22000 {
		return nil, fmt.Errorf("get song: %s", resp.ErrorMessage)
	}

	return resp, nil
}

func GetSongsRaw(songIds ...string) (*SongsResponse, error) {
	return std.GetSongsRaw(songIds...)
}

// 批量获取歌曲，遗留接口，不推荐使用
func (a *API) GetSongsRaw(songIds ...string) (*SongsResponse, error) {
	params := sreq.Params{
		"songIds": strings.Join(songIds, ","),
	}
	resp := new(SongsResponse)
	err := a.Request(sreq.MethodGet, APIGetSongs, sreq.WithQuery(params)).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrorCode != 22000 {
		return nil, fmt.Errorf("get songs: %d", resp.ErrorCode)
	}

	return resp, nil
}

func GetSongLyric(songId string) (string, error) {
	return std.GetSongLyric(songId)
}

func (a *API) GetSongLyric(songId string) (string, error) {
	resp, err := a.GetSongLyricRaw(songId)
	if err != nil {
		return "", err
	}

	return resp.LrcContent, nil
}

func GetSongLyricRaw(songId string) (*SongLyricResponse, error) {
	return std.GetSongLyricRaw(songId)
}

// 获取歌词
func (a *API) GetSongLyricRaw(songId string) (*SongLyricResponse, error) {
	params := sreq.Params{
		"songid": songId,
	}

	resp := new(SongLyricResponse)
	err := a.Request(sreq.MethodGet, APIGetSongLyric, sreq.WithQuery(params)).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrorCode != 0 && resp.ErrorCode != 22000 {
		return nil, fmt.Errorf("get lyric: %d", resp.ErrorCode)
	}

	return resp, nil
}
