package diff

// ChangeSet ...
type ChangeSet struct {
	line       int
	oldLine    string
	newLine    string
	changeType string // "new", "update", "delete"
}

// FilesInfo ...
type FilesInfo struct {
	oldFilePath string
	newFilePath string

	oldFileLineChecksums []string
	newFileLineChecksums []string

	oldFileLines map[string]string
	newFileLines map[string]string

	oldFileChecksum string
	newFileChecksum string
}

// Instantiate create diff instance with two files to compares
// 'of' stand for old file and 'nf' for new file
func Instantiate(oldFilePath, newFilePath string) *FilesInfo {
	ofLines, ofLinesChecksums, ofChecksum := readLines(oldFilePath)
	nfLines, nfLinesChecksums, nfChecksum := readLines(newFilePath)

	return &FilesInfo{
		oldFilePath: oldFilePath,
		newFilePath: newFilePath,

		oldFileChecksum: ofChecksum,
		newFileChecksum: nfChecksum,

		oldFileLineChecksums: ofLinesChecksums,
		newFileLineChecksums: nfLinesChecksums,

		oldFileLines: ofLines,
		newFileLines: nfLines,
	}
}

// CompareLines ...
func (df FilesInfo) CompareLines() (isSame bool, changeSets []ChangeSet) {
	isSame = false

	if df.oldFileChecksum == df.newFileChecksum {
		isSame = true
		return
	}

	oldFileLineChecksums, newFileLineChecksums := harmonizeSlicesSize(df.oldFileLineChecksums, df.newFileLineChecksums)
	for index, ofLine := range oldFileLineChecksums {
		if nfLine := newFileLineChecksums[index]; ofLine != nfLine {
			changeSets = append(changeSets, ChangeSet{
				line:       index + 1,
				oldLine:    df.oldFileLines[ofLine],
				newLine:    df.newFileLines[nfLine],
				changeType: "update",
			})
		}
	}

	return
}

// FindNewLines lines in new file not in old file
func (df FilesInfo) FindNewLines() (isSame bool, changeSets []ChangeSet) {
	isSame = false

	if df.oldFileChecksum == df.newFileChecksum {
		isSame = true
		return
	}

	for index, nfLine := range df.newFileLineChecksums {
		added := true
		for _, ofLine := range df.oldFileLineChecksums {
			if nfLine == ofLine {
				added = false
				break
			}
		}

		if added {
			changeSets = append(changeSets, ChangeSet{
				line:       index + 1,
				changeType: "new",
				newLine:    df.newFileLines[nfLine],
			})
		}
	}
	return
}

// FindRemovedLines find lines in old file not in new lines
func (df FilesInfo) FindRemovedLines() (isSame bool, changeSets []ChangeSet) {
	isSame = false

	if df.oldFileChecksum == df.newFileChecksum {
		isSame = true
		return
	}

	for index, ofLine := range df.oldFileLineChecksums {
		removed := true
		for _, nfLine := range df.newFileLineChecksums {
			if nfLine == ofLine {
				removed = false
				break
			}
		}

		if removed {
			changeSets = append(changeSets, ChangeSet{
				line:       index + 1,
				changeType: "delete",
				newLine:    df.oldFileLines[ofLine],
			})
		}
	}

	return
}
