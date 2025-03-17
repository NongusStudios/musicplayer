package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/dhowden/tag"
	"github.com/hajimehoshi/go-mp3"
)

type Directory map[string][]string

type Album struct {
	artist   string
	year     int
	coverArt []byte
	duration time.Duration // seconds
}

type Song struct {
	title       string
	trackNumber int
	duration    time.Duration // seconds
	filePath    string
}

type Library struct {
	// key = album name
	albums map[string]Album
	songs  map[string][]Song
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

func IndexMediaDirectory(lib *Library, dir Directory) error {
	for path, files := range dir {
		for _, file := range files {
			filePath := fmt.Sprintf("%s/%s", path, file)

			f, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer f.Close()

			metaData, err := tag.ReadFrom(f)
			if err != nil {
				return err
			}

			// Check if albums has an entry for the album of the current song
			// If not add it
			if _, ok := lib.albums[metaData.Album()]; !ok {
				lib.albums[metaData.Album()] = Album{
					artist:   metaData.Artist(),
					year:     metaData.Year(),
					coverArt: slices.Clone(metaData.Picture().Data),
				}
			}

			duration, err := GetAudioFileDuration(filePath)
			if err != nil {
				return err
			}
			album := lib.albums[metaData.Album()]
			album.duration += duration
			lib.albums[metaData.Album()] = album

			trackNumber, _ := metaData.Track()

			lib.songs[metaData.Album()] = append(lib.songs[metaData.Album()], Song{
				title:       metaData.Title(),
				trackNumber: trackNumber,
				duration:    duration,
				filePath:    filePath,
			})
		}
	}

	// Sort Songs by track number
	for _, songs := range lib.songs {
		slices.SortFunc(songs, func(a Song, b Song) int {
			return a.trackNumber - b.trackNumber
		})
	}

	return nil
}

func GetLibraryFromMediaDirectories(dirs []string) (Library, error) {
	library := Library{}
	library.albums = make(map[string]Album)
	library.songs = make(map[string][]Song)

	// Read Media Directories
	for _, path := range dirs {
		dir, err := ReadFilesInDirectory(path)
		if err != nil {
			return Library{}, err
		}

		err = IndexMediaDirectory(&library, dir)
		if err != nil {
			return Library{}, err
		}
	}

	return library, nil
}
