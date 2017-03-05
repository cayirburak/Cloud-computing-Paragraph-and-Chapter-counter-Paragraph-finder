package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	"strconv"
)

//global typeler tanımlanıyor.
type Paragraph []string
type Book struct {
	paragraphs []Paragraph //string dizisi tipinde paragraphs tanımlandı
}

//Book tipinde bir pointer tanımlandı.
var book *Book

//count fonksiyonu için global değişkenler
var paragraf_sayisi =0
var chapter=0

func check(e error){ //hata kontrolü yapan fonksiyon
	if e != nil {
           panic(e)
	}
}
// Bu kodu çalıştırmak için command line'dan `go run task.go` çalıştırman yeterli.

//readBook reads the book at filePath. Keep the at a glabal variable at access it at 'count' and 'query' functions 
func readBook(filePath string) (book *Book) {
	//	YOUR CODE HERE. Read the book and save it to a global variable, something like `var Book [][]string`
   
   dat , err := ioutil.ReadFile(filePath) //dosya adıyla okunuyor.
   check(err)
   kit := string(dat) //stringe dönüştürülüyor
   paragraf_sayisi = strings.Count(kit,"\r\n\r\n") //paragrafları sayıyor (chapterlarla beraber)
   chap1 := strings.Count(kit,"\r\n\r\n\r\n") //chapterların boşluklarını alıyor.
   paragraf_sayisi = paragraf_sayisi-chap1 //chapter boşluklarını sildik.
   chapter = (strings.Count(kit, "Chapter")) //chapter sayısını aldık.
   //fmt.Printf("\n%q\n", dat)  // karakter dizgesi analizi için kullanılacak.

	parts := strings.Split(string(dat), "\r\n\r\n\r\n\r\n" )  //her chapter gördüğünde ayırıyor
	paragraphs := make([]Paragraph, 0, len(parts)) //slice oluşturuluyor

	for _, part := range parts {
		paragraph := strings.Split(part, "\r\n\r\n") //paragraflar ayrıştırılıyor.

		paragraphs = append(paragraphs, paragraph) // her paragraf sırayla diziye aktarılıyor
	}

	return &Book{paragraphs} //işlemler sonucu oluşturulan değer döndürülüyor.

}

func query(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	result := "mevcut değil" //varsayılan olarak chapter sayısı bulunamazsa gönderilecek hata mesajı

	var paragraph *Paragraph
	c, _ := strconv.Atoi(q.Get("c")) // kullanıcıdan alınan sorgu string den int e çevriliyor.

	if c > 0 && len(book.paragraphs) >= c {
		paragraph = &book.paragraphs[c-1]
		result = strings.Join(*paragraph, "\n\n") //split ederken kaybedilen boşluklar geri kazandırılıyor.
	}

	p, _ := strconv.Atoi(q.Get("p")) // kullanıcıdan alınan sorgu string den int e çevriliyor.

	if paragraph != nil && p > 0 && len(*paragraph) > p {
		result = (*paragraph)[p]
	}else if paragraph != nil && len(*paragraph) <= p{ //kullanıcı geçersiz paragraf veya chapter sayısı girerse hata mesajı gösteriliyor
		result = "Mevcut değil."
	}

	fmt.Fprintf(w, result)
}
func count(w http.ResponseWriter, r *http.Request) {
	chapCount := 0
	paraCount := 0
	 
	 //global değişkenleri atıyoruz.
    chapCount = chapter  
    paraCount = paragraf_sayisi 
   	
	fmt.Fprintf(w, "chapter: %d\nparagraph: %d\n", chapCount, paraCount)
}
func otherwise(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}
func main() {
	
	book = readBook("book.txt")	
	http.HandleFunc("/count", count)
	http.HandleFunc("/query", query)
	http.HandleFunc("/", otherwise)
	http.ListenAndServe(":8080", nil)
}
