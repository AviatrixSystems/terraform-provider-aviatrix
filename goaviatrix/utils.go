package goaviatrix

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var ErrNotFound = fmt.Errorf("ErrNotFound")

func ExpandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}

// difference returns the elements in a that aren't in b
func Difference(a, b []string) []string {
	mb := map[string]bool{}
	for _, x := range b {
		mb[x] = true
	}
	ab := make([]string, 0)
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			ab = append(ab, x)
		}
	}
	return ab
}

// DifferenceSlice returns the one-dimension elements in two-dimension slice a that aren't in two-dimension b
func DifferenceSlice(a, b [][]string) [][]string {
	if a == nil || len(a) == 0 || b == nil || len(b) == 0 {
		return a
	}

	aa := make([]string, 0)
	for i := range a {
		temp := ""
		for j := range a[i] {
			temp += a[i][j]
		}
		aa = append(aa, temp)
	}

	bb := make([]string, 0)
	for t := range b {
		temp := ""
		for m := range b[t] {
			temp += b[t][m]
		}
		bb = append(bb, temp)
	}

	mb := map[string]bool{}
	for x := range bb {
		mb[bb[x]] = true
	}
	ab := make([][]string, 0)
	for x := range aa {
		if _, ok := mb[aa[x]]; !ok {
			ab = append(ab, a[x])
		}
	}
	return ab
}

func ReadFile(local_filepath string) (string, string, error) {
	// File being read must be .json
	// Returns filename, contents of json file, error string
	if filepath.Ext(local_filepath) != ".json" {
		return "", "", errors.New("local filepath doesn't lead to a json file")
	}
	filename := filepath.Base(local_filepath)
	jsonFile, err := os.Open(local_filepath)
	if err != nil {
		return "", "", errors.New("Failed to open local json file: " + err.Error())
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return "", "", errors.New("Failed to read local json file: " + err.Error())
	}
	contents := string(byteValue[:])
	return filename, contents, nil
}

func ReadPemFile(local_filepath string) (string, string, error) {
	// File being read must be .pem
	// Returns filename, contents of pem file, error string
	if filepath.Ext(local_filepath) != ".pem" {
		return "", "", errors.New("local filepath doesn't lead to a pem file")
	}
	filename := filepath.Base(local_filepath)
	pemFile, err := os.Open(local_filepath)
	if err != nil {
		return "", "", errors.New("Failed to open local pem file: " + err.Error())
	}
	defer pemFile.Close()
	byteValue, err := ioutil.ReadAll(pemFile)
	if err != nil {
		return "", "", errors.New("Failed to read local pem file: " + err.Error())
	}
	contents := string(byteValue[:])
	return filename, contents, nil
}

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func TagListStrColon(tagListStr []string) []string {
	if tagListStr != nil {
		for i := range tagListStr {
			tagListStr[i] = strings.ReplaceAll(tagListStr[i], ":", "\\\\:")
			tagListStr[i] = strings.Replace(tagListStr[i], "\\\\:", ":", 1)
		}
		return tagListStr
	}
	return nil
}

func CompareMapOfInterface(map1 map[string]interface{}, map2 map[string]interface{}) bool {
	if map1 == nil && map2 == nil {
		return true
	}
	if map1 == nil || map2 == nil || len(map1) != len(map2) {
		return false
	}

	for key := range map1 {
		if val, ok := map2[key]; ok {
			if map1[key] != val {
				return false
			}
			continue
		}
		return false
	}
	return true
}
