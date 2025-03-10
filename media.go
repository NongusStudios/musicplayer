package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"go.senan.xyz/taglib"
)

type Directory map[string][]string

type Album struct {
	name   string
	artist string
	date   string
	length time.Duration // seconds
	songs  []Song
}

type Song struct {
	trackNumber int
	title       string
	artists     []string
	date        string
	length      time.Duration // seconds
	filePath    string
}

type Library struct {
	albums []Album
}

var SupportedFileExtensions = []string{".mp3"}

func ReadFilesInDirectory(path string) (Directory, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	directory := make(Directory)

	for _, entry := range entries {
		fullFileName := fmt.Sprintf("%s/%s", path, entry.Name())
		if entry.IsDir() {
			subDirectory, e := ReadFilesInDirectory(fullFileName)
			if e != nil {
				return nil, e
			}

			for subPath, files := range subDirectory {
				directory[subPath] = append(directory[subPath], files...)
			}

			continue
		}
		if slices.Contains(SupportedFileExtensions, filepath.Ext(fullFileName)) {
			directory[path] = append(directory[path], entry.Name())
		}
	}

	return directory, nil
}

func CombineLibraries(libA Library, libB Library) Library {
	lib := Library{}
	lib.albums = append(lib.albums, libA.albums...)
	lib.albums = append(lib.albums, libB.albums...)
	return lib
}

func GetAudioFileDuration(path string) (time.Duration, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return 0, err
	}

	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		return 0, err
	}

	const sampleSize = int64(4)
	samples := decoder.Length() / sampleSize
	duration := time.Duration(samples/int64(decoder.SampleRate())) * time.Second

	return duration, nil
}

func ReadDirectoryAudioTags(dir Directory) (Library, error) {
	songsByAlbum := make(map[string][]Song)

	for path, files := range dir {
		for _, file := range files {
			filePath := fmt.Sprintf("%s/%s", path, file)

			tags, err := taglib.ReadTags(filePath)
			if err != nil {
				return Library{}, err
			}

			// Get album
			album := path
			if len(tags[taglib.Album]) > 0 {
				album = tags[taglib.Album][0]
			}

			// Get title
			title := file
			if len(tags[taglib.Title]) > 0 {
				title = tags[taglib.Title][0]
			}

			// Get date
			date := "n/a"
			if len(tags[taglib.Date]) > 0 {
				date = tags[taglib.Date][0]
			}

			// Get track number
			tn, err := strconv.Atoi(tags[taglib.TrackNumber][0])
			if err != nil {
				return Library{}, err
			}

			// Get duration
			duration, err := GetAudioFileDuration(filePath)
			if err != nil {
				return Library{}, err
			}

			songsByAlbum[album] = append(songsByAlbum[tags[taglib.Album][0]], Song{
				trackNumber: tn,
				title:       title,
				artists:     tags[taglib.Artist],
				date:        date,
				length:      duration,
				filePath:    filePath,
			})
		}
	}

	lib := Library{}

	for albumName, songs := range songsByAlbum {
		slices.SortFunc(songs, func(a Song, b Song) int {
			return a.trackNumber - b.trackNumber
		})

		var totalLength time.Duration
		for _, song := range songs {
			totalLength += song.length
		}

		artist := "unknown"
		if len(songs[0].artists) > 0 {
			artist = songs[0].artists[0]
		}

		album := Album{
			name:   albumName,
			artist: artist,
			date:   songs[0].date,
			length: totalLength,
			songs:  slices.Clone(songs),
		}

		lib.albums = append(lib.albums, album)
	}

	return lib, nil
}

func GetLibraryFromMediaDirectories(dirs []string) (Library, error) {
	library := Library{}

	// Read Media Directories
	for _, path := range dirs {
		dir, err := ReadFilesInDirectory(path)
		if err != nil {
			return Library{}, err
		}

		var lib Library
		lib, err = ReadDirectoryAudioTags(dir)
		if err != nil {
			return Library{}, err
		}

		library = CombineLibraries(library, lib)
	}

	return library, nil
}
