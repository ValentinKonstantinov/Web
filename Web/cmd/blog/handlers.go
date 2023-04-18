package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type indexPageData struct {
	Logo            string
	PageTitle       string
	PageSubtitle    string
	PageButton      string
	Home            string
	Categories      string
	About           string
	Contact         string
	MenuNature      string
	MenuPhotography string
	MenuRelaxation  string
	MenuVacation    string
	MenuTravel      string
	MenuAdventure   string
	MainTitle       string
	Title           string
	BottomHeader    string
	BottomButton    string
	Featured        []featuredPostData
	PostsData       []postsDatas
}

type featuredPostData struct {
	PostHeader    string `db:"title"`
	PostText      string `db:"subtitle"`
	PublishDate   string `db:"publish_date"`
	PostAuthor    string `db:"author"`
	PostAuthorUrl string `db:"author_url"`
	PostImageUrl  string `db:"image_url"`
	PostCategory  string `db:"category"`
}

type postsDatas struct {
	PostHeader    string `db:"title"`
	PostText      string `db:"subtitle"`
	PublishDate   string `db:"publish_date"`
	PostAuthor    string `db:"author"`
	PostAuthorUrl string `db:"author_url"`
	PostImageUrl  string `db:"image_url"`
}

func index(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		featuredPosts, err := featured(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500) // В случае ошибки парсинга - возвращаем 500
			log.Println(err)
			return // Не забываем завершить выполнение ф-ии
		}

		posts, err := postsData(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500) // В случае ошибки парсинга - возвращаем 500
			log.Println(err)
			return // Не забываем завершить выполнение ф-ии
		}

		ts, err := template.ParseFiles("pages/index.html") // Главная страница блога
		if err != nil {
			http.Error(w, "Internal Server Error", 500) // В случае ошибки парсинга - возвращаем 500
			log.Println(err)
			return // Не забываем завершить выполнение ф-ии
		}

		data := indexPageData{
			Logo:            "Escape",
			PageTitle:       "Let's do it together",
			PageSubtitle:    "We travel the world in search of stories. Come along for the ride.",
			PageButton:      "View Latest Posts",
			Home:            "HOME",
			Categories:      "CATEGORIES",
			About:           "ABOUT",
			Contact:         "CONTACT",
			MenuNature:      "Nature",
			MenuPhotography: "Photography",
			MenuRelaxation:  "Relaxation",
			MenuVacation:    "Vacation",
			MenuTravel:      "Travel",
			MenuAdventure:   "Adventure",
			MainTitle:       "Featured Posts",
			Title:           "Most Recent",
			BottomHeader:    "Stay in Touch",
			BottomButton:    "Submit",
			Featured:        featuredPosts,
			PostsData:       posts,
		}

		err = ts.Execute(w, data) // Заставляем шаблонизатор вывести шаблон в тело ответа
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		log.Println("Request completed successfully")
	}
}

func featured(db *sqlx.DB) ([]featuredPostData, error) {
	const query = `
		SELECT
			title,
			subtitle,
			publish_date,
			author,
			author_url,
			image_url,
			category
		FROM
			post
		WHERE featured = 1
	` // Составляем SQL-запрос для получения записей для секции featured-posts

	var posts []featuredPostData // Заранее объявляем массив с результирующей информацией

	err := db.Select(&posts, query) // Делаем запрос в базу данных
	if err != nil {                 // Проверяем, что запрос в базу данных не завершился с ошибкой
		return nil, err
	}

	return posts, nil
}

func postsData(db *sqlx.DB) ([]postsDatas, error) {
	const query = `
		SELECT
			title,
			subtitle,
			publish_date,
			author,
			author_url,
			image_url
		FROM
			post
		WHERE featured = 0
	` // Составляем SQL-запрос для получения записей для секции featured-posts

	var posts []postsDatas // Заранее объявляем массив с результирующей информацией

	err := db.Select(&posts, query) // Делаем запрос в базу данных
	if err != nil {                 // Проверяем, что запрос в базу данных не завершился с ошибкой
		return nil, err
	}

	return posts, nil
}
