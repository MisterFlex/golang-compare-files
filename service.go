package diff

// ChangeSet ...
type ChangeSet struct {
	line       int
	oldLine    string
	newLine    string
	changeType string // "new", "update", "delete"
}

// CompareLines compare lines of two files
func Compare(currentFilePath, newFilePath string) (isTheSame bool, differentLines []ChangeSet) {
	cfLines, cfChecksum := readLines(currentFilePath)
	nfLines, nfChecksum := readLines(newFilePath)
	isTheSame = true

	// there is some difference between those two files
	if cfChecksum != nfChecksum {
		var removedLines, newLines, updatedLines []ChangeSet

		sizeOfCfLines := len(cfLines)
		sizeOfNfLines := len(nfLines)

		if sizeOfCfLines > sizeOfNfLines {
			// removed lines
			for i, val := range cfLines[sizeOfNfLines:] {
				removedLines = append(removedLines, ChangeSet{
					line:       sizeOfNfLines + 1 + i,
					changeType: "delete",
					oldLine:    val,
					newLine:    "",
				})
			}
		} else {
			// new lines
			for i, val := range nfLines[sizeOfCfLines:] {
				newLines = append(newLines, ChangeSet{
					line:       sizeOfCfLines + 1 + i,
					changeType: "new",
					oldLine:    "",
					newLine:    val,
				})
			}
		}

		// changesets
		for index, line := range cfLines {
			if index == len(nfLines) {
				break
			}

			if nLine := nfLines[index]; line != nLine {
				updatedLines = append(updatedLines, ChangeSet{
					line:       index + 1,
					changeType: "update",
					oldLine:    line,
					newLine:    nLine,
				})
			}
		}

		differentLines = append(differentLines, append(updatedLines, newLines...)...)
		isTheSame = false
	}

	return
}

// CompareLines ...
func CompareLines(oldFile, newFile string) (isSame bool, changeSets []ChangeSet) {
	oldFileLines, ofChecksum := readLines(oldFile)
	newFileLines, nfChecksum := readLines(newFile)
	isSame = false

	if ofChecksum == nfChecksum {
		isSame = true
		return
	}

	oldFileLines, newFileLines = harmonizeSlicesSize(oldFileLines, newFileLines)
	for index, ofLine := range oldFileLines {
		if nfLine := newFileLines[index]; ofLine != nfLine {
			changeSets = append(changeSets, ChangeSet{
				line:       index + 1,
				oldLine:    ofLine,
				newLine:    nfLine,
				changeType: "update",
			})
		}
	}

	return
}

// FindNewLines lines in new file not in old file
func FindNewLines(oldFile, newFile string) (isSame bool, changeSets []ChangeSet) {
	oldFileLines, ofChecksum := readLines(oldFile)
	newFileLines, nfChecksum := readLines(newFile)
	isSame = false

	if ofChecksum == nfChecksum {
		isSame = true
		return
	}

	for index, nfLine := range newFileLines {
		added := true
		for _, ofLine := range oldFileLines {
			if nfLine == ofLine {
				added = false
				break
			}
		}

		if added {
			changeSets = append(changeSets, ChangeSet{
				line:       index + 1,
				changeType: "new",
				newLine:    nfLine,
			})
		}
	}
	return
}

// FindRemovedLines find lines in old file not in new lines
func FindRemovedLines(oldFile, newFile string) (isSame bool, changeSets []ChangeSet) {
	oldFileLines, ofChecksum := readLines(oldFile)
	newFileLines, nfChecksum := readLines(newFile)
	isSame = false

	if ofChecksum == nfChecksum {
		isSame = true
		return
	}

	for index, ofLine := range oldFileLines {
		removed := true
		for _, nfLine := range newFileLines {
			if nfLine == ofLine {
				removed = false
				break
			}
		}

		if removed {
			changeSets = append(changeSets, ChangeSet{
				line:       index + 1,
				changeType: "delete",
				newLine:    ofLine,
			})
		}
	}

	return
}
