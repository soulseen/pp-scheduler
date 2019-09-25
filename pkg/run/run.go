package run

import (
	"flag"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/soulseen/pp-scheduler/pkg/controller"
	"github.com/soulseen/pp-scheduler/pkg/predicate"
	"github.com/soulseen/pp-scheduler/pkg/prioritize"
	"github.com/soulseen/pp-scheduler/pkg/routes"
	_ "github.com/soulseen/pp-scheduler/pkg/sqlite"

	log "github.com/golang/glog"
)

var (
	PipelinePriority = prioritize.Prioritize{
		Name: "pipeline",
		Func: prioritize.Pipeline,
	}

	TruePredicate = predicate.Predicate{
		Name: "alwaystrue",
		Func: predicate.AlwaysTrue,
	}
)

func Run() {

	flag.Parse()

	log.Info("Start controller ....")
	go controller.RunController()

	router := httprouter.New()
	routes.AddVersion(router)

	predicates := []predicate.Predicate{TruePredicate}
	for _, p := range predicates {
		routes.AddPredicate(router, p)
	}

	priorities := []prioritize.Prioritize{PipelinePriority}
	for _, p := range priorities {
		routes.AddPrioritize(router, p)
	}

	log.Info("info: server starting on the port :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
