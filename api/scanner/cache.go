package scanner

import (
	"path"
	"sync"

	"github.com/pkg/errors"
)

type AlbumScannerCache struct {
	path_contains_photos map[string]bool
	photo_types          map[string]MediaType
	ignore_data          map[string][]string
	mutex                sync.Mutex
}

func MakeAlbumCache() *AlbumScannerCache {
	return &AlbumScannerCache{
		path_contains_photos: make(map[string]bool),
		photo_types:          make(map[string]MediaType),
		ignore_data:          make(map[string][]string),
	}
}

// Insert single album directory in cache
func (c *AlbumScannerCache) InsertAlbumPath(path string, contains_photo bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.path_contains_photos[path] = contains_photo
}

// Insert album path and all parent directories up to the given root directory in cache
func (c *AlbumScannerCache) InsertAlbumPaths(end_path string, root string, contains_photo bool) {
	curr_path := path.Clean(end_path)
	root_path := path.Clean(root)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for curr_path != root_path || curr_path == "." {

		c.path_contains_photos[curr_path] = contains_photo

		curr_path = path.Dir(curr_path)
	}
}

func (c *AlbumScannerCache) AlbumContainsPhotos(path string) *bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	contains_photo, found := c.path_contains_photos[path]
	if found {
		// log.Printf("Album cache hit: %s\n", path)
		return &contains_photo
	}

	return nil
}

// func (c *AlbumScannerCache) InsertPhotoType(path string, content_type MediaType) {
// 	c.mutex.Lock()
// 	defer c.mutex.Unlock()

// 	(c.photo_types)[path] = content_type
// }

func (c *AlbumScannerCache) GetMediaType(path string) (*MediaType, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	result, found := c.photo_types[path]
	if found {
		// log.Printf("Image cache hit: %s\n", path)
		return &result, nil
	}

	mediaType, err := getMediaType(path)
	if err != nil {
		return nil, errors.Wrapf(err, "get media type (%s)", path)
	}

	if mediaType != nil {
		(c.photo_types)[path] = *mediaType
	}

	return mediaType, nil
}

func (c *AlbumScannerCache) GetAlbumIgnore(path string) *[]string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ignore_data, found := c.ignore_data[path]
	if found {
		return &ignore_data
	}

	return nil
}

func (c *AlbumScannerCache) InsertAlbumIgnore(path string, ignore_data []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.ignore_data[path] = ignore_data
}
