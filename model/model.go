package model

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/Sirupsen/logrus"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/remohammadi/delphi/common"
)

var (
	db *gorm.DB
)

type Question struct {
	ID              int64      `json:"id"              gorm:"primary_key"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	DeletedAt       *time.Time `json:"-"               sql:"index"`
	Title           string     `json:"title"           gorm:"type:varchar(255)"`
	CompareQuestion string     `json:"compareQuestion" gorm:"type:varchar(255)"`
	Items           []Item     `json:"items"`
	Votes           []Vote     `json:"votes"`
}

func (q *Question) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", rand.Int())
	return nil
}

type Item struct {
	ID         int64      `json:"id"          gorm:"primary_key"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `json:"-"           sql:"index"`
	QuestionID int64      `json:"questionId"  gorm:"index"`
	Title      string     `json:"title"       gorm:"type:varchar(140)"`
}

func (i *Item) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", rand.Int())
	return nil
}

type Vote struct {
	ID           int64      `json:"id"            gorm:"primary_key"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"-"             sql:"index"`
	QuestionID   int64      `json:"questionId"  gorm:"index"`
	OrderedItems string     `json:"orderedItems"  gorm:"type:varchar(500)"`
}

func (v *Vote) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", rand.Int())
	return nil
}

func init() {
	var err error
	engine := common.ConfigString("DB_ENGINE")
	if engine == "" {
		logrus.Panic("No database engine is specified")
	}
	params := common.ConfigString("DB_PARAMS")
	db, err = gorm.Open(engine, params)
	if err != nil {
		logrus.WithFields(logrus.Fields{"engine": engine, "params": params}).WithError(err).Panic("Failed to connect database")
	}
	db.SetLogger(&common.GormLogger{})
	db.LogMode(common.ConfigString("LOG_LEVEL") == "DEBUG")
	db.AutoMigrate(&Question{}, &Item{}, &Vote{})
}

type RestFuncs struct{}

func (rf *RestFuncs) GetQuestion(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	q := Question{}
	if e := db.Preload("Items").First(&q, id).Error; e != nil {
		logrus.WithError(e).Debugf("GetQuestion :: id=%s", id)
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&q)
}

func (rf *RestFuncs) PostQuestion(w rest.ResponseWriter, r *rest.Request) {
	q := Question{}
	if err := r.DecodeJsonPayload(&q); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := db.Save(&q).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&q)
}

func (rf *RestFuncs) Vote(w rest.ResponseWriter, r *rest.Request) {
	vote := Vote{}
	if err := r.DecodeJsonPayload(&vote); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := Question{}
	if db.First(&q, vote.QuestionID).Error != nil {
		rest.NotFound(w, r)
		return
	}

	if err := db.Save(&vote).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&vote)
}

func Api() *rest.Api {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	rf := &RestFuncs{}
	router, err := rest.MakeRouter(
		rest.Post("/questions", rf.PostQuestion),
		rest.Get("/questions/:id", rf.GetQuestion),
		rest.Put("/vote", rf.Vote),
	)
	if err != nil {
		logrus.WithError(err).Fatal("While initializing api router")
	}
	api.SetApp(router)

	return api
}
