package main

import (
	"fmt"
	"os"
	"time"

	"gioui.org/widget"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

type Player struct {
	CurrentAlbum        string
	CurrentTrack        int
	CurrentSongElapsed  float64
	CurrentSongProgress *widget.Float
	Lib                 Library

	otoCtx            *oto.Context
	readyChan         chan struct{}
	currentPlayer     *oto.Player
	currentFile       *os.File
	currentDecodedMp3 *mp3.Decoder
	time              time.Time
}

func InitPlayer(library Library) (Player, error) {
	op := &oto.NewContextOptions{}

	// Usually 44100 or 48000. Other values might cause distortions in Oto
	op.SampleRate = 44100

	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	op.ChannelCount = 2

	// Format of the source. go-mp3's format is signed 16bit integers.
	op.Format = oto.FormatSignedInt16LE

	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		return Player{}, err
	}

	return Player{
		CurrentAlbum:        "",
		CurrentTrack:        -1,
		CurrentSongProgress: new(widget.Float),
		Lib:                 library,

		otoCtx:    otoCtx,
		readyChan: readyChan,
		time:      time.Now(),
	}, nil
}

func (p *Player) closeCurrentFile() {
	if p.currentFile != nil {
		err := p.currentPlayer.Close()
		if err != nil {
			panic("player.Close failed: " + err.Error())
		}

		p.currentFile.Close()
	}
}

func (p *Player) Update() {
	if p.IsPlaying() {
		dt := time.Since(p.time)
		p.time = time.Now()
		p.CurrentSongElapsed += dt.Seconds()
		fmt.Println(p.CurrentSongElapsed)

		p.CurrentSongProgress.Value = float32(p.CurrentSongElapsed / p.Lib.Songs[p.CurrentAlbum][p.CurrentTrack].Duration.Seconds())
	}
}

func (p *Player) PlaySong(album string, track int) error {
	if p.currentFile != nil {
		p.Pause()
		p.closeCurrentFile()
	}

	if songs, ok := p.Lib.Songs[album]; ok {
		if track >= len(songs) {
			return fmt.Errorf("track number %d out of bounds in player.library.songs[\"%s\"]", track, album)
		}

		song := songs[track]

		var err error

		p.currentFile, err = os.Open(song.FilePath)
		if err != nil {
			return err
		}

		p.currentDecodedMp3, err = mp3.NewDecoder(p.currentFile)
		if err != nil {
			return err
		}
		<-p.readyChan

		p.currentPlayer = p.otoCtx.NewPlayer(p.currentDecodedMp3)
		p.currentPlayer.Play()

		p.CurrentAlbum = album
		p.CurrentTrack = track
		p.CurrentSongElapsed = 0

		return nil
	}

	return fmt.Errorf("%s is not an entry in player.library.songs", album)
}

func (p *Player) IsPlaying() bool {
	if p.currentPlayer == nil || !p.currentPlayer.IsPlaying() {
		return false
	}
	return true
}

func (p *Player) TogglePlayBack() {
	if p.IsPlaying() {
		p.Pause()
		return
	}
	p.Resume()
}

func (p *Player) Resume() {
	if p.currentPlayer == nil {
		return
	}
	p.currentPlayer.Play()
}

func (p *Player) Pause() {
	if p.currentPlayer == nil {
		return
	}
	p.currentPlayer.Pause()
}
