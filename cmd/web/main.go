package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"golangs.org/snippetbox/pkg/models/mysql" // Новый импорт

	_ "github.com/go-sql-driver/mysql"
)

// Добавляем поле snippets в структуру application. Это позволит
// сделать объект SnippetModel доступным для наших обработчиков.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel
}

func main() {
	// Создаем новый флаг командной строки, значение по умолчанию: ":4000".
	// Добавляем небольшую справку, объясняющая, что содержит данный флаг. 
	// Значение флага будет сохранено в переменной addr.
	addr := flag.String("addr", ":4000", "Сетевой адрес веб-сервера")
	dsn := flag.String("dsn", "zzz:000@/snippetbox?parseTime=true", "Название MySQL источника данных")
	// Мы вызываем функцию flag.Parse() для извлечения флага из командной строки.
	// Она считывает значение флага из командной строки и присваивает его содержимое
	// переменной. Вам нужно вызвать ее *до* использования переменной addr
	// иначе она всегда будет содержать значение по умолчанию ":4000". 
	// Если есть ошибки во время извлечения данных - приложение будет остановле
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	// Инициализируем экземпляр mysql.SnippetModel и добавляем его в зависимостях.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &mysql.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	
	// Значение, возвращаемое функцией flag.String(), является указателем на значение
	// из флага, а не самим значением. Нам нужно убрать ссылку на указатель
	// то есть перед использованием добавьте к нему префикс *. Обратите внимание, что мы используем
	// функцию log.Printf() для записи логов в журнал работы нашего приложения.
	infoLog.Printf("Запуск сервера на http://127.0.0.1%s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
