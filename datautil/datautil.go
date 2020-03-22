// Package datautil contains functions that can be useful in data processing.
package datautil

// FindStrArr return the index of element `elm` in the string slice `slice`.
func FindStrArr(slice []string, elm string) (index int) {
	index = -1
	for i, item := range slice {
		if item == elm {
			index = i
		}
	}
	return
}

// FindIntArr return the index of element `elm` in the int slice `slice`.
func FindIntArr(slice []string, elm string) (index int) {
	index = -1
	for i, item := range slice {
		if item == elm {
			index = i
		}
	}
	return
}
