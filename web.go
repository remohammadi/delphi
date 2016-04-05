package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/remohammadi/delphi/common"
)

func main() {
	// Static files
	fs := http.FileServer(http.Dir(common.ConfigString("STATIC_DIR")))
	http.Handle("/s/", http.StripPrefix("/s/", fs))
	http.Handle("/favicon.ico", fs)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		base := "base"
		name := r.URL.Path[1:]
		var data map[string]interface{}
		contentType := "text/html; charset=utf-8"
		if name == "" {
			name = "index"
		}
		if strings.HasPrefix(name, "email-") {
			err := r.ParseForm()
			if err != nil {
				logrus.WithField("name", name).WithError(err).Debug("Failed to ParseForm")
			} else {
				data = make(map[string]interface{})
				for k, v := range r.Form {
					data[k] = v[0]
				}
				logrus.WithField("name", name).WithField("data", data).Debug("Context for email template")
			}

			if strings.HasPrefix(name, "email-html/") {
				base = "email-base-html"
			}
			if strings.HasPrefix(name, "email-text/") {
				base = "email-base-text"
				contentType = "text/plain; charset=utf-8"
			}
			name = name[11:]
		}

		name = fmt.Sprintf("%s.tmpl", name)

		writeFunc, err := common.RenderTemplate(w, base, name, data)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("404 (%s)", err)))
		} else {
			w.Header().Set("Content-Type", contentType)
			writeFunc()
		}

		logrus.WithFields(logrus.Fields{
			logrus.ErrorKey: err,
			"method":        r.Method,
			"uri":           r.RequestURI,
			"referer":       r.Referer(),
			"user-agent":    r.UserAgent(),
		}).Infoln()
	})

	bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
	fmt.Printf("listening on %s...", bind)
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		panic(err)
	}
}
