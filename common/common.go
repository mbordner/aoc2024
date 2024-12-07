package common

func PopulateStringCombinationsAtLength(results map[string]bool, pickChars string, prefix string, length int) {
	if length == 0 {
		results[prefix] = true
		return
	}

	for i := 0; i < len(pickChars); i++ {
		nextPrefix := prefix + string(pickChars[i])
		PopulateStringCombinationsAtLength(results, pickChars, nextPrefix, length-1)
	}
}
