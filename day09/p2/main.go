package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"strconv"
)

const (
	null = int64(-1)
)

type Disk []int64
type FileID int64

type File struct {
	id      FileID
	length  int64
	address int64
}

type Free struct {
	length  int64
	address int64
}

func main() {
	data, _ := file.GetContent("../data.txt")
	nextFileID := FileID(0)

	files := make(map[FileID]*File)
	free := make([]Free, 0, 1000)

	address := int64(0)
	for i := 0; i < len(data); i++ {
		val := getIntVal(data[i])
		if i%2 == 0 {
			f := File{id: nextFileID, length: val, address: address}
			files[nextFileID] = &f
			nextFileID++
		} else {
			f := Free{length: val, address: address}
			free = append(free, f)
		}
		address += val
	}

	for id := nextFileID - 1; id > 0; id-- {
		f := files[id]
		for j := 0; j < len(free); j++ {
			// if this free space is after the file, we wouldn't want to move it
			if free[j].address > f.address {
				break
			}
			// if we can move this file back
			if free[j].length >= f.length {
				last := f.address
				// move file
				f.address = free[j].address
				// adjust free space
				remaining := free[j].length - f.length
				if remaining == 0 {
					if j == 0 {
						free = free[1:]
					} else {
						free = append(free[0:j], free[j+1:]...)
					}
				} else {
					free[j].length = remaining
					free[j].address += f.length
				}
				for k := j; k < len(free); k++ {
					if free[k].address+free[k].length == last {
						free[k].length += f.length
						if k < len(free)-1 {
							if free[k].address+free[k].length == free[k+1].address {
								free[k].length += free[k+1].length
								free = append(free[0:k+1], free[k+2:]...)
							}
						}
						break
					} else if free[k].address > last {
						if last+f.length == free[k].address {
							free[k].address -= f.length
						} else {
							insert := Free{address: last, length: f.length}
							if k == 0 {
								free = append([]Free{insert}, free...)
							} else {
								end := free[k:]
								free = append(append(free[0:k], insert), end...)
							}
						}
						break
					}
				}
				break
			}
		}
	}

	var maxFile *File
	for id := FileID(0); id < nextFileID; id++ {
		if maxFile == nil {
			maxFile = files[id]
		} else {
			if files[id].address > maxFile.address {
				maxFile = files[id]
			}
		}
	}
	disk := make(Disk, maxFile.address+maxFile.length)
	for i := 0; i < len(disk); i++ {
		disk[i] = null
	}

	for id := FileID(0); id < nextFileID; id++ {
		for i := files[id].address; i < files[id].address+files[id].length; i++ {
			disk[i] = int64(id)
		}
	}

	sum := int64(0)
	for i := 0; i < len(disk); i++ {
		if disk[i] != null {
			sum += int64(i) * disk[i]
		}
	}

	fmt.Println(sum)
}

func getIntVal(c byte) int64 {
	val, _ := strconv.ParseInt(string(c), 10, 64)
	return val
}
