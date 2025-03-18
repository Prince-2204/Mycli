package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	_ "github.com/go-sql-driver/mysql"
	"go_cli.princeaman.net/internal/models"
)

type application struct {
	datas *models.DataModel
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


const (
    apiEndpoint = "https://api.openai.com/v1/chat/completions"
)

func fetchAIResponse(search string) {
    client := resty.New()
	apiKey :="YOUR_API_KEY"
    response, err := client.R().
        SetAuthToken(apiKey).
        SetHeader("Content-Type", "application/json").
        SetBody(map[string]interface{}{
            "model":      "gpt-3.5-turbo",
            "messages":   []interface{}{map[string]interface{}{"role": "system", "content": search}},
            "max_tokens": 50,
        }).
        Post(apiEndpoint)

    if err != nil {
        log.Fatalf("Error while sending the request: %v", err)
    }
	body := response.Body()

    var data map[string]interface{}
    err = json.Unmarshal(body, &data)
    if err != nil {
        fmt.Println("Error while decoding JSON response:", err)
        return
    }

   
    content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
    fmt.Println(content)
 
}

func dirName(){
	path,err := os.Getwd()
	if err !=nil{
		log.Fatal(err)
	}

	fmt.Println("Directory is: ",path)
}

func getFileNames(){
	var filenames []string
	dirpath,err := os.Getwd()
	if err!= nil{
		log.Fatal(err)
	}

	err1:= filepath.WalkDir(dirpath,func(path string , d fs.DirEntry,err error)error{
		if err!=nil{
			return nil
		}
		if !d.IsDir(){
			filenames = append(filenames,path )
		}
		return nil
	})
	if err1!= nil{
		log.Fatal(err1)
	}
	fmt.Print("Files in your Directory is:\n")
	for _, fileName := range filenames {
		fmt.Println(fileName)
	}

}

func main() {
	dsn := flag.String("dsn", "Your sql detail", "MySQL data source name")


	addTask := flag.String("add", "", "Add a new task")
	delTask := flag.Int("del", -1, "Delete a task by ID")
	listTasks := flag.Bool("list", false, "List all tasks")
	completeTask := flag.Int("complete", -1, "Mark a task as completed")
	search := flag.String("search","","Search anything with gpt")
	directory := flag.Bool("dir",false,"Get the current directory")
	listfiles := flag.Bool("ls",false,"Get the list of files in your directory")
	
	flag.Parse()

	
	db, err := openDB(*dsn)
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer db.Close()

	app := &application{
		datas: &models.DataModel{DB: db},
	}

	
	if *addTask != "" {
		err := app.datas.AddTask(*addTask)
		if err == nil {
			fmt.Println("Task added:", *addTask)
		}
	} else if *delTask != -1 {
		err := app.datas.DeleteTask(*delTask)
		if err == nil {
			fmt.Println("Task deleted:", *delTask)
		}
	} else if *completeTask != -1 {
		err := app.datas.MarkTaskCompleted(*completeTask)
		if err == nil {
			fmt.Println("Task marked as completed:", *completeTask)
		}
	} else if *listTasks {
		tasks, err := app.datas.ListTasks()
		if err == nil {
			fmt.Println("Tasks:")
			for _, task := range tasks {
				fmt.Println("-", task)
			}
		}
	}else if *search != ""{
		fetchAIResponse(*search)
	}else if *directory{
		dirName()
	}else if *listfiles{
		getFileNames()
	} else {
		fmt.Println("No valid command provided.")
	}
}
