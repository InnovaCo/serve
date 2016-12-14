package utils

import "math/rand"

func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func Contains(elm string, list []string) bool {
	for _, v := range list {
		if v == elm {
			return true
		}
	}
	return false
}

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func MergeMaps(maps ...map[string]string) map[string]string {
	out := make(map[string]string, 0)
	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

func MapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if w, ok := b[k]; !ok || v != w {
			return false
		}
	}

	return true
}

type BySortIndex []map[string]string

func (a BySortIndex) Len() int           { return len(a) }
func (a BySortIndex) Less(i, j int) bool { return a[i]["sortIndex"] < a[j]["sortIndex"] }
func (a BySortIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }


var allLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var RandomString = func(length int) string {
    b := make([]rune, length)
    for i := range b {
        b[i] = allLetters[rand.Intn(len(allLetters))]
    }
    return string(b)
}

