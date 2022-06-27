package main

/*func TestGetAllBook(t *testing.T) {
	s := []struct {
		description string
		input       string
		output      []Book
	}{
		{"get all", "/books", []Book{
			{1, "hail", "eliot", "penguin", "02/02/2019"},
			{2, "harr potter", "jk", "penguin", "02/02/1996"},
		}},
	}

	for _, v := range s {

		req, err2 := http.NewRequest("GET", "v.input", nil)
		if err2 != nil {
			t.Fatalf("could not create request %v", err2)
		}
		w := httptest.NewRecorder()
		res := w.Result()
		getPosts(w, req)

		data, err := io.ReadAll(res.Body)

		var b Book
		errr := json.Unmarshal(data, &b)
		if errr != nil {
			fmt.Println("decode error")
		}
		assert.Equal(t, v.output, b)
		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status OK;got %v", res.StatusCode)
		}
		defer res.Body.Close()

	}

}*/
import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBook(t *testing.T) {
	testcases := []struct {
		desc   string
		req    string
		expRes []Book
	}{
		{"success", "/books", []Book{
			{120, "east of eden", Author{2, "theodar", "seuss", "03/10/1977", "doctor"}, "penguin", "02/02/1981"},
			{121, "the star", Author{1, "charlies", "rice", "02/02/1984", "howard"}, "epical", "09/10/1998"},
			{122, "brave", Author{3, "samuel", "clemens", "09/09/1980", "mark"}, "arihant", "08/04/1984"},
			{123, "east", Author{4, "charlotte", "bronte", "06/01/1981", "currer"}, "scholastic", "07/04/1979"},
		}},
	}

	for _, test := range testcases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "localhost:8000"+test.req, nil)

		getBooks(w, req)
		resp := w.Result()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var allBooks []Book

		err = json.Unmarshal(data, &allBooks)
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, test.expRes, allBooks)

		err = resp.Body.Close()

	}
}

func TestGetByID(t *testing.T) {
	testcases := []struct {
		desc   string
		input  string
		output Book
	}{
		{"The details for book with id 120 ", "/books/120",
			Book{120, "east of eden", Author{2, "theodar", "seuss", "03/10/1977", "doctor"}, "penguin", "02/02/1981"}},
	}

	for _, test := range testcases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", test.input, nil)

		getBookById(w, req)
		resp := w.Result()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var Book Book

		err = json.Unmarshal(data, &Book)
		if err != nil {
			return
		}

		assert.Equal(t, test.output, Book)

		err = resp.Body.Close()
	}
}

func TestPostBook(t *testing.T) {
	testcases := []struct {
		desc       string
		book       Book
		statusCode int
	}{
		{"Details posted.", Book{1, "ABC", Author{1, "anu", "07-10-1998", "cilios", "w"}, "Arihanth", "18-08-2018"},
			http.StatusCreated},
		{"Invalid publication.", Book{2, "", Author{2, "anu", "07-10-1998", "cilios", "x"}, "Oxford", "21-04-1985"},
			http.StatusBadRequest},
		{"Invalid author.", Book{3, "ABC", Author{9, "anu", "19-12-1972", "cilios", "c"}, "scholastic", "21-04-1985"},
			http.StatusBadRequest,
		},
		{"Invalid publication date.", Book{4, "ABC", Author{4, "anu", "19-12-1972", "cilios", "d"}, "penguin", "21-04-1879"},
			http.StatusBadRequest,
		},
		{"Invalid publication date.", Book{5, "ABC", Author{5, "anu", "19-12-1972", "cilios", ""}, "penguin", "21-04-2031"},
			http.StatusBadRequest,
		},
	}

	for _, test := range testcases {
		newData, _ := json.Marshal(test.book)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/books", bytes.NewBuffer(newData))
		createPostBook(w, req)
		resp := w.Result()

		assert.Equal(t, test.statusCode, resp.StatusCode)

	}
}
func TestPostAuthor(t *testing.T) {
	testcases := []struct {
		desc       string
		author     Author
		statusCode int
	}{
		{"Valid case.", Author{1, "charlies", "rice", "02/02/1984", "howard"},
			http.StatusCreated},
		{"Invalid case.", Author{0, "anu", "sharma", "07-10-1998", "a"},
			http.StatusBadRequest},
		{"Invalid case.", Author{2, "", "lian", "07-10-1998", "x"},
			http.StatusBadRequest},
		{"Invalid case.", Author{0, "", "toy", "07-10-1998", "y"},
			http.StatusBadRequest},
	}

	for _, test := range testcases {
		newData, _ := json.Marshal(test.author)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/postauthor", bytes.NewBuffer(newData))
		createPostAuthor(w, req)
		resp := w.Result()

		assert.Equal(t, test.statusCode, resp.StatusCode)

	}
}
func TestPutBook(t *testing.T) {
	testcases := []struct {
		desc      string
		reqId     string
		reqBody   Book
		expRes    Book
		expStatus int
	}{
		{"Valid Details", "1", Book{Title: "hello", Author: Author{1, "elior", "gold", "02/02/1967", "c"}, "penguin", "02/02/1960"}, Book{}, 200},
		{"Publication should be Scholastic/Pengiun/Arihanth", "1", Book{Title: "Jay", Author: nil, Publication: "Arvind", PublishedDate: 11 / 03 / 2002}, Book{}, http.StatusBadRequest},
		{"Published date should be between 1880 and 2022", "1", Book{Title: "", Author: nil, Publication: "", PublishedDate: "1/1/1870"}, Book{}, http.StatusBadRequest},
		{"Published date should be between 1880 and 2022", "1", Book{Title: "", Author: nil, Publication: "", PublishedDate: "1/1/2222"}, Book{}, http.StatusBadRequest},
		{"Author should exist", "1", Book{}, Book{}, http.StatusBadRequest},
		{"Title can't be empty", "1", Book{Title: "", Author: nil, Publication: "", PublishedDate: ""}, Book{}, http.StatusBadRequest},
	}
	for i, tc := range testcases {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.reqBody)
		req := httptest.NewRequest(http.MethodPost, "localhost:8000/book/"+tc.reqId, bytes.NewReader(body))
		putBook(w, req)
		defer w.Result().Body.Close()

		if w.Result().StatusCode != tc.expStatus {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
		res, _ := io.ReadAll(w.Result().Body)
		resBook := Book{}
		json.Unmarshal(res, &resBook)
		if resBook != tc.expRes {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}

func TestPutAuthor(t *testing.T) {
	testcases := []struct {
		desc      string
		reqBody   Author
		expRes    Author
		expStatus int
	}{
		{"Valid details", Author{FirstName: "RD", LastName: "Sharma", Dob: "2/11/1989", PenName: "Sharma"}, Author{1, "RD", "Sharma", "2/11/1989", "Sharma"}, http.StatusOK},
		{"InValid details", Author{FirstName: "", LastName: "Sharma", Dob: "2/11/1989", PenName: "Sharma"}, Author{}, http.StatusBadRequest},
	}
	for i, tc := range testcases {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.reqBody)
		req := httptest.NewRequest(http.MethodPost, "localhost:8000/author", bytes.NewReader(body))
		putAuthor(w, req)
		defer w.Result().Body.Close()

		if w.Result().StatusCode != tc.expStatus {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
		res, _ := io.ReadAll(w.Result().Body)
		resAuthor := Author{}
		json.Unmarshal(res, &resAuthor)
		if resAuthor != tc.expRes {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}
func TestDeleteBook(t *testing.T) {
	testcases := []struct {
		desc      string
		reqId     string
		expStatus int
	}{
		{"Valid Details", "120", http.StatusOK},
		{"Book does not exists", "100", http.StatusNotFound},
	}
	for i, tc := range testcases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "localhost:8000/deletebook/"+tc.reqId, nil)
		deleteBook(w, req)
		defer w.Result().Body.Close()

		if w.Result().StatusCode != tc.expStatus {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}

func TestDeleteAuthor(t *testing.T) {
	testcases := []struct {
		desc      string
		reqId     string
		expStatus int
	}{
		{"Valid Details", "1", http.StatusOK},
		{"Author does not exists", "100", http.StatusBadRequest},
	}
	for i, tc := range testcases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "localhost:8000/deleteauthor/"+tc.reqId, nil)
		deleteAuthor(w, req)
		defer w.Result().Body.Close()

		if w.Result().StatusCode != tc.expStatus {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}
