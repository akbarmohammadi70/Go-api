package main
import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	
	"fmt"
	"io/ioutil"
	"net/http"
    "encoding/json"
    "github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

//go get -u github.com/go-sql-driver/mysql
// go get -u github.com/gorilla/mux
//go get -u gorm.io/gorm
//go get -u gorm.io/driver/mysql

func(article Article) TableName() string {
	return "article"
}

func ShowArticles(writer http.ResponseWriter,request *http.Request){
    //switch request.Method {
    //   case "GET":
    //        writer.WriteHeader(http.StatusOk)
    //        json.NewEncoder(writer).Encode("this is a Get method")
    //   case "POST":
    //        writer.WriteHeader(http.StatusBadRequest)
    //        json.NewEncoder(writer).Encode("this is a Post method") 
    //   default:
    //        json.NewEncoder(writer).Encode("other requests..")
    // }
    // fmt.Fprintf(writer, format: "this is a articles page")
    //writer.Header().Set("content-type","application/json")  // add heaser json
    //writer.WriteHeader(http.StatusOk)   // Or StatusBadRequest Or StatusCreated Or ...
    // json.NewEncoder(writer).Encode(Articles)
	//https://gorm.io/docs/connecting_to_the_database.html
	// dsn := "go:qazwsx@tcp(192.168.56.118:3306)/go?charset=utf8mb4&parseTime=True&loc=Local"
	input := mux.Vars(request)
	id := input["id"]
	db, err := getConnection()		
	
	if err != nil {
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(err.Error())
		return
	}

	oneArticle := Article{} 
	db.Find(&oneArticle,id)
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(oneArticle)


}
func showHomePage(writer http.ResponseWriter,request *http.Request){
    fmt.Fprintf(writer , "this is a Pist article")
}
func showArticlesPost(writer http.ResponseWriter,request *http.Request){
    fmt.Fprintf(writer , "this is a home page")
}

func showOneArticle(writer http.ResponseWriter,request *http.Request){
    input := mux.Vars(request)
    id := input["id"]
    for _ , article := range Articles {
        if article.Id == id {
            writer.Header().Set("content-type","application/json")
            json.NewEncoder(writer).Encode(article)
            
        }
        
    }
}

func getConnection() (*gorm.DB,error){
	// dbDriver := "mysql"
	// dbUser := "go"
	// dbPassword := "qazwsx"
	// dbName := "go"
	
	// db, err := gorn.Open(dbDriver, dbUser+":"+dbPassword+"@tcp(192.168.56.118:3306)"+"/"+dbName)
	dsn := "go:qazwsx@tcp(192.168.56.118:3306)/go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})		

	return db , err
	
}

func addNewArticle(writer http.ResponseWriter,request *http.Request) {
    reqBody , _ := ioutil.ReadAll(request.Body)
    oneArticle := Article{}
    err := json.Unmarshal(reqBody, &oneArticle)
    // if err != nil {
    //     // Info is a Logger with LogLevel INF
	// 	fmt.Printf("notok")
		
	// }
	// Articles = append(Articles , oneArticle)
	//github.com/go-sql-driver/mysql
	// dbDriver := "mysql"
	// dbUser := "go"
	// dbPassword := "qazwsx"
	// dbName := "go"
	
	// db, err := sql.Open(dbDriver, dbUser+":"+dbPassword+"@tcp(192.168.56.118:3306)"+"/"+dbName)
	

	// if err != nil {
	// 	writer.WriteHeader(http.StatusConflict)
	// 	json.NewEncoder(writer).Encode(err.Error())
	// 	return
	// }
	db , err := getConnection()
	//defer db.Close()   
	
	if err != nil {
	 	writer.WriteHeader(http.StatusConflict)
	 	json.NewEncoder(writer).Encode(err.Error())
	 	return
		}
	
	// insert , err := db.Prepare("insert into article (title , description , content) values (?,?,?)")
	db.Select("title","description","content").Create(oneArticle)
	if err != nil {
		writer.WriteHeader(http.StatusCreated)
		json.NewEncoder(writer).Encode("article append")
		return
	}

	// insert.Exec(oneArticle.Title,oneArticle.Description,oneArticle.Content)
	// defer insert.Close()



// 	writer.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(writer).Encode("article append ...")
}
func UpdateArticle(writer http.ResponseWriter,request *http.Request) {
	inputs := mux.Vars(request)
	id := inputs["id"]
	reqBody , _ := ioutil.ReadAll(request.Body)
    oneArticle := Article{}
    err := json.Unmarshal(reqBody, &oneArticle)
	db, err := getConnection()
	if err != nil {
		writer.WriteHeader(http.StatusConflict)
		json.NewEncoder(writer).Encode(err.Error())
		return
	}
	db.Model(&oneArticle).Where("id = ?",id).Updates(&oneArticle)
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode("article update..")

}



func deleteOneArticle(writer http.ResponseWriter,request *http.Request) {
     inputs := mux.Vars(request)
     id := inputs["id"]
     db , _ := getConnection()
	 oneArticle := Article{}
	 db.Delete(&oneArticle,id)
    //  for index, article := range Articles{
    //      if article.Id==id{
	// 		Articles=append(Articles[:index],Articles[index+1:]...)
	// 	}
    //  }
     writer.WriteHeader(http.StatusOK)
     json.NewEncoder(writer).Encode("deleted Successfully")
 }

func handleRequests() {
    my_router := mux.NewRouter()
    my_router.HandleFunc("/",showHomePage).Methods("GET")   // add method GET with gorilla mux
    my_router.HandleFunc("/articles/{id}",showOneArticle).Methods("GET")
	my_router.HandleFunc("/articles",ShowArticles).Methods("GET")
    my_router.HandleFunc("/articles",addNewArticle).Methods("POST") // add method POST gorilla mux
    my_router.HandleFunc("/articles/{id}",deleteOneArticle).Methods("DELETE")
	my_router.HandleFunc("/articles/{id}",UpdateArticle).Methods("PUT")

    http.ListenAndServe(":8484", my_router)
    
    
}



type Article struct {
    Id string
    Title string
    Description string
    Content string
}

var Articles [] Article

func main() {
    Articles = append(Articles , Article{Id:"1",Title:"title number 1" , Description:"description number 1" , Content: "content number 1"})
    Articles = append(Articles , Article{Id:"2",Title:"title number 2" , Description:"description number 2" , Content: "content number 2"})
    Articles = append(Articles , Article{Id:"3",Title:"title number 3" , Description:"description number 3" , Content: "content number 3"})

  
    handleRequests()
}
