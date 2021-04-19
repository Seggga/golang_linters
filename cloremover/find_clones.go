package cloremover

import (
	"fmt"
	"io/fs"
	"sort"

	"github.com/sirupsen/logrus"
)

// FindClones looks for clone-files in the given directory and it's subdirectories
func FindClones(conf *ConfigType, log logrus.FieldLogger, fileSystem fs.FS) ([]fileData, error) {

	//enumerate all subdirectories in the given dirPath
	log.Debug("try to enumerate subdirectories in the given directory")
	dirSlice, err := enumDirs(conf, fileSystem)
	if err != nil {
		return nil, err
	}

	var (
		// pool of workers for each subdirectory
		pool = make(chan struct{}, len(dirSlice))
		// ch transfers data from workers to dataCollector
		ch = make(chan fileData, 1)
		// slice to hold files's data
		fileSlice []fileData
	)

	// start pool manager
	go manager(pool, ch)
	// start ch-writers - one goroutine for each subdirectory
	for _, someDir := range dirSlice {
		log.Debugf("try to read all files in the directory %s", someDir)
		go func(someDir string) {
			enumFiles(someDir, ch, fileSystem)
			pool <- struct{}{}
		}(someDir)
	}

	// start ch-reader
	for someData := range ch {
		fileSlice = append(fileSlice, someData)
	}
	log.Debug("list of files has been constructed")
	// obtain slice of clone-files only
	log.Debug("try to filter out unique files")
	return filterUnique(fileSlice), nil
}

// manager waits for all the enumFiles functions to end their work
func manager(pool <-chan struct{}, ch chan fileData) {
	for i := 0; i < cap(pool); i++ {
		<-pool
	}
	close(ch)
}

// enumFiles enumerates files in the given subdirectory and sends fileData structure about
// each file via the channel
func enumFiles(dirPath string, ch chan<- fileData, fileSystem fs.FS) {
	//files, _ := ioutil.ReadDir(dirPath)
	files, _ := fs.ReadDir(fileSystem, dirPath)
	for _, someFile := range files {
		if !someFile.IsDir() {

			someFileData := new(fileData)
			someFileData.dir = dirPath
			someFileData.fileName = someFile.Name()
			fileInfo, _ := someFile.Info()
			someFileData.sizeInBytes = uint64(fileInfo.Size())

			ch <- *someFileData
		}
	}
}

// enumDirs enumerates subdirectories in the given folder.
func enumDirs(conf *ConfigType, fileSystem fs.FS) ([]string, error) {
	/*

		// get absolute directory path
		dirPath, err := filepath.Abs(conf.DirPath)
		if err != nil {
			return nil, err
		}

		// check for dirPath existance
		if _, err := os.Stat(dirPath); err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("Entered directory was not found ( %s )\n", dirPath)
			} else {
				return nil, err
			}
		}
	*/
	//conf.DirPath = dirPath

	// get a slice of subdirectories
	var dirSlice []string
	/*
		err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				dirSlice = append(dirSlice, path)
			}
			return nil
		})
	*/
	//root := filepath.Base(filepath.Dir(conf.DirPath))
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			dirSlice = append(dirSlice, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	//dirSlice = append(dirSlice, conf.DirPath)
	//conf.DirPath = dirPath
	return dirSlice, nil
}

// filterUnique looks throw all the given []fileData and sets "id" field of each element.
// Clone-files are identical if their name and size are equal.
// filterUnique produces a slice with only data about clones and a map with a number of each ID in the given slice.
func filterUnique(fileSlice []fileData) []fileData /*, map[uint16]uint16*/ {
	// search for clones and mark the fileData structures with id that corresponds to the pair Name-Size
	index := uint32(1)
	idMap := make(map[string]uint32)
	cloneMap := make(map[uint32]uint32)
	for i := 0; i < len(fileSlice); i += 1 {
		cloneID := fileSlice[i].fileName + fmt.Sprint(fileSlice[i].sizeInBytes)
		id, ok := idMap[cloneID]
		if !ok {
			//this file is unique
			idMap[cloneID] = index
			fileSlice[i].id = index
			cloneMap[index] += 1
			index += 1
			continue
		}
		fileSlice[i].id = id
		cloneMap[id] += 1
	}
	// count capacity of slice to store only clone's data
	var capacity uint32
	for _, num := range cloneMap {
		if num > 1 {
			capacity += num
		}
	}
	// fill the slice of clones with the data
	fileSliceClones := make([]fileData, capacity)
	i := 0
	for _, someFileData := range fileSlice {
		if cloneMap[someFileData.id] > 1 {
			fileSliceClones[i] = someFileData
			i += 1
		}
	}

	sortData(fileSliceClones)

	return fileSliceClones
}

// sort according to flags user has set
func sortData(fileSlice []fileData) {
	sort.Slice(fileSlice, func(i, j int) bool {
		return fileSlice[i].id < fileSlice[j].id
	})
}

/*
func DirDigger(conf *ConfigType) ([]string, error) {

	// get absolute directory path
	dirPath, err := filepath.Abs(conf.DirPath)
	if err != nil {
		return nil, err
	}

	// check for dirPath existance
	if _, err := os.Stat(dirPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("Entered directory was not found ( %s )\n", dirPath)
		} else {
			return nil, err
		}
	}

	conf.DirPath = dirPath

	// get a slice of subdirectories
	var dirSlice []string

	currentFS := os.DirFS(dirPath)

	err = fs.WalkDir(currentFS, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && path != "." {
			dirSlice = append(dirSlice, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return dirSlice, nil
}

*/
