package lemmatizer

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

type Lemmatizer interface {
	Lemmatize(word string) []string
}

type ItalianLemmatizer struct {
	suffixMap        map[string][]string
	suffixList       []string
	lemmatizationMap map[string][]string
}

var Lemmatizers map[string]Lemmatizer

var itLemmatizer *ItalianLemmatizer

func (lemmatizer ItalianLemmatizer) Lemmatize(word string) []string {
	if lems, ok := itLemmatizer.lemmatizationMap[word]; ok {
		return lems
	}
	var lemmas []string
	exist := make(map[string]bool)
	for _, suffix := range lemmatizer.suffixList {
		if strings.HasSuffix(word, suffix) {
			values := lemmatizer.suffixMap[suffix]
			for _, v := range values {
				lemma := word[:len(word)-len(suffix)] + v
				if !exist[lemma] {
					lemmas = append(lemmas, lemma)
					exist[lemma] = true
				}
			}
		}
	}
	return lemmas
}

func init() {
	Lemmatizers = make(map[string]Lemmatizer)
	itSuffixMap := make(map[string][]string)
	itSuffixMap["o"] = []string{"are", "ere", "ire"}
	itSuffixMap["i"] = []string{"o", "e", "a", "are", "ere", "ire"}
	itSuffixMap["a"] = []string{"o", "are", "ere", "ire"}
	itSuffixMap["e"] = []string{"o", "a", "ere", "ire"}
	itSuffixMap["ò"] = []string{"are"}
	itSuffixMap["ate"] = []string{"are"}
	itSuffixMap["ano"] = []string{"are", "ere", "ire"}
	itSuffixMap["ete"] = []string{"ere"}
	itSuffixMap["ono"] = []string{"ere", "ire"}
	itSuffixMap["ino"] = []string{"are"}
	itSuffixMap["iamo"] = []string{"are", "ire", "ere"}
	itSuffixMap["isca"] = []string{"ire"}
	itSuffixMap["isco"] = []string{"ire"}
	itSuffixMap["isci"] = []string{"ire"}
	itSuffixMap["isce"] = []string{"ire"}
	itSuffixMap["iscono"] = []string{"ire"}
	itSuffixMap["iscano"] = []string{"ire"}
	itSuffixMap["ite"] = []string{"ire"}
	itSuffixMap["ato"] = []string{"are"}
	itSuffixMap["ati"] = []string{"are"}
	itSuffixMap["ata"] = []string{"are"}
	itSuffixMap["ate"] = []string{"are"}
	itSuffixMap["uto"] = []string{"ere"}
	itSuffixMap["uti"] = []string{"ere"}
	itSuffixMap["ute"] = []string{"ere"}
	itSuffixMap["uta"] = []string{"ere"}
	itSuffixMap["ito"] = []string{"ire"}
	itSuffixMap["iti"] = []string{"ire"}
	itSuffixMap["ita"] = []string{"ire"}
	itSuffixMap["avo"] = []string{"are"}
	itSuffixMap["avi"] = []string{"are"}
	itSuffixMap["ava"] = []string{"are"}
	itSuffixMap["avamo"] = []string{"are"}
	itSuffixMap["avate"] = []string{"are"}
	itSuffixMap["avano"] = []string{"are"}
	itSuffixMap["evo"] = []string{"ere"}
	itSuffixMap["evi"] = []string{"ere"}
	itSuffixMap["eva"] = []string{"ere"}
	itSuffixMap["evamo"] = []string{"ere"}
	itSuffixMap["evate"] = []string{"ere"}
	itSuffixMap["evano"] = []string{"ere"}
	itSuffixMap["ivo"] = []string{"ire"}
	itSuffixMap["ivi"] = []string{"ire"}
	itSuffixMap["iva"] = []string{"ire"}
	itSuffixMap["ivamo"] = []string{"ire"}
	itSuffixMap["ivate"] = []string{"ire"}
	itSuffixMap["ivano"] = []string{"ire"}
	itSuffixMap["ai"] = []string{"are"}
	itSuffixMap["asti"] = []string{"are"}
	itSuffixMap["ammo"] = []string{"are"}
	itSuffixMap["aste"] = []string{"are"}
	itSuffixMap["arono"] = []string{"are"}
	itSuffixMap["ei"] = []string{"ere"}
	itSuffixMap["etti"] = []string{"ere"}
	itSuffixMap["esti"] = []string{"ere"}
	itSuffixMap["é"] = []string{"ere"}
	itSuffixMap["ette"] = []string{"ere"}
	itSuffixMap["emmo"] = []string{"ere"}
	itSuffixMap["este"] = []string{"ere"}
	itSuffixMap["erono"] = []string{"ere"}
	itSuffixMap["ettero"] = []string{"ere"}
	itSuffixMap["ii"] = []string{"ire"}
	itSuffixMap["isti"] = []string{"ire"}
	itSuffixMap["ì"] = []string{"ire"}
	itSuffixMap["immo"] = []string{"ire"}
	itSuffixMap["iste"] = []string{"ire"}
	itSuffixMap["irono"] = []string{"ire"}
	itSuffixMap["erò"] = []string{"are", "ere"}
	itSuffixMap["erai"] = []string{"are", "ere"}
	itSuffixMap["erà"] = []string{"are", "ere"}
	itSuffixMap["eremo"] = []string{"are", "ere"}
	itSuffixMap["erete"] = []string{"are", "ere"}
	itSuffixMap["eranno"] = []string{"are", "ere"}
	itSuffixMap["irò"] = []string{"ire"}
	itSuffixMap["irai"] = []string{"ire"}
	itSuffixMap["irà"] = []string{"ire"}
	itSuffixMap["iremo"] = []string{"ire"}
	itSuffixMap["irete"] = []string{"ire"}
	itSuffixMap["iranno"] = []string{"ire"}
	itSuffixMap["erei"] = []string{"are", "ere"}
	itSuffixMap["eresti"] = []string{"are", "ere"}
	itSuffixMap["erebbe"] = []string{"are", "ere"}
	itSuffixMap["eremmo"] = []string{"are", "ere"}
	itSuffixMap["ereste"] = []string{"are", "ere"}
	itSuffixMap["erebbero"] = []string{"are", "ere"}
	itSuffixMap["irei"] = []string{"ire"}
	itSuffixMap["iresti"] = []string{"ire"}
	itSuffixMap["irebbe"] = []string{"ire"}
	itSuffixMap["iremmo"] = []string{"ire"}
	itSuffixMap["ireste"] = []string{"ire"}
	itSuffixMap["irebbero"] = []string{"ire"}
	itSuffixMap["iate"] = []string{"are", "ire", "ere"}
	itSuffixMap["assi"] = []string{"are"}
	itSuffixMap["asse"] = []string{"are"}
	itSuffixMap["assimo"] = []string{"are"}
	itSuffixMap["assero"] = []string{"are"}
	itSuffixMap["essi"] = []string{"ere"}
	itSuffixMap["esse"] = []string{"ere"}
	itSuffixMap["essimo"] = []string{"ere"}
	itSuffixMap["essero"] = []string{"ere"}
	itSuffixMap["issi"] = []string{"ire"}
	itSuffixMap["isse"] = []string{"ire"}
	itSuffixMap["issimo"] = []string{"ire"}
	itSuffixMap["issero"] = []string{"ire"}
	itSuffixMap["ante"] = []string{"are"}
	itSuffixMap["anti"] = []string{"are"}
	itSuffixMap["ente"] = []string{"ere", "ire"}
	itSuffixMap["enti"] = []string{"ere", "ire"}
	itSuffixMap["ando"] = []string{"are"}
	itSuffixMap["endo"] = []string{"ere", "ire"}
	itSuffixList := make([]string, 0, len(itSuffixMap))
	for key := range itSuffixMap {
		itSuffixList = append(itSuffixList, key)
	}
	sort.Slice(itSuffixList, func(i, j int) bool {
		return len(itSuffixList[i]) > len(itSuffixList[j])
	})
	itLemmatizer = &ItalianLemmatizer{suffixMap: itSuffixMap, suffixList: itSuffixList}
	itLemmatizer.lemmatizationMap = make(map[string][]string)
	exist := make(map[string]bool)
	f, _ := os.Open("../../assets/lemmatization/it.txt")
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		lem := scanner.Text()
		scanner.Scan()
		conj := scanner.Text()
		if !exist[conj] {
			itLemmatizer.lemmatizationMap[conj] = []string{lem}
			exist[conj] = true
		} else {
			itLemmatizer.lemmatizationMap[conj] = append(itLemmatizer.lemmatizationMap[conj], lem)
		}
	}
	Lemmatizers["it"] = itLemmatizer
}
