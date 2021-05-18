package boolean

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	irmodels "github.com/mikkelstb/ir_models"
)

type TermDictionary struct {
	Terms []Term
	punctuations string
	newline_regex *regexp.Regexp
}

type Term struct {
	ID                 string
	Document_frequency int
	Postings_list      []int
}

func (this *TermDictionary) Sort() {
	sort.Slice(this.Terms, func(i, j int) bool { return (this.Terms)[i].ID < (this.Terms)[j].ID })
}

func (this *TermDictionary) say() {
	fmt.Println(this.Terms)
}

func (this *TermDictionary) Init() {
	this.punctuations = "/()\"~.,;:!?-–*#%'»«"
	this.newline_regex = regexp.MustCompile(`\r?\n`)
}

func (this *TermDictionary) AddDocument(document irmodels.Article) {

	types := make(map[string]int)

	//Make a slice with all words
	//document.Text = strings.
	document.Text = this.newline_regex.ReplaceAllString(document.Text, "")
	tokens := strings.Split(document.Text, " ")

	for _, token := range tokens {

		//Remove all punctuations, and convert to lower case
		token = strings.Trim(token, this.punctuations)
		token = strings.ToLower(token)

		if len([]rune(token)) > 2 {
			types[token]++
		}
	}
	//Make slice for new terms, max number is the number of types
	newTerms := make([]Term, 0, len(types))
	for wtype := range types {
		found := this.isPresent(wtype)
		if found > -1 {
			this.updateTerm(wtype, found, document.Doc_id)
		} else {
			newTerms = append(newTerms, Term{ID: wtype, Document_frequency: 1, Postings_list: []int{document.Doc_id}})
		}
	}
	this.addTerms(newTerms, document.Doc_id)
	this.Sort()
}

func (this *TermDictionary) updateTerm(term string, index int, docID int) {
	this.Terms[index].Postings_list = append(this.Terms[index].Postings_list, docID)
	this.Terms[index].Document_frequency++
}

func (this *TermDictionary) addTerms(terms []Term, doc_id int) {
	this.Terms = append(this.Terms, terms...)
}

func (this *TermDictionary) isPresent(term string) int {

	low := 0
	high := len(this.Terms) - 1

	for low <= high {
		median := (low + high) / 2
		if this.Terms[median].ID < term {
			low = median + 1
		} else {
			high = median - 1
		}
	}
	if low == len(this.Terms) || this.Terms[low].ID != term {
		return -1
	} else {
		//fmt.Println("Found " + term + " on index: " + strconv.Itoa(low))
		return low
	}
}

func (this *TermDictionary) intersect(post_list_a []int, post_list_b []int) []int {
	a := 0
	b := 0
	var result []int
	for a < len(post_list_a) && b < len(post_list_b) {
		if post_list_a[a] == post_list_b[b] {
			result = append(result, post_list_a[a])
			a++
			b++
		} else if post_list_a[a] < post_list_b[b] {
			a++
		} else {
			b++
		}
	}
	return result
}

func (this *TermDictionary) intersectMultiple(list ...[]int) []int {

	//fmt.Println("List length: " + strconv.Itoa(len(list)))

	if len(list) == 0 {
		return nil
	} else if len(list) == 1 {
		return list[0]
	}

	//Sorting slices according to size
	sort.Slice(list, func(i, j int) bool {
		return len(list[i]) < len(list[j])
	})

	result := list[0]
	list = list[1:]

	for (len(list) > 0) && (len(result) > 0) {
		result = this.intersect(result, list[0])
		list = list[1:]
	}

	//fmt.Println(result)

	return result
}

func (this *TermDictionary) Search (searchstring string) []int {
	queries := this.parseQuery(searchstring)
	var results [][]int

	for _, query := range queries {

		index := this.isPresent(query)
		if index > -1 {
			results = append(results, this.Terms[index].Postings_list)
		} else { return nil }	
	}

	return this.intersectMultiple(results...)
}


func (this *TermDictionary) parseQuery (query string) []string {
	return strings.Split(query, " ")
}