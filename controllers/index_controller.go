package controllers

import (
	"digitalLibrary/utils"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

var commonRecover = utils.CommonRecover

type fn func(bson.M) (int64, error)

func GetAllCount() (map[string]int64, []error) {
	var wg sync.WaitGroup
	var errs []error

	allCounts := map[string]int64{}
	totalKeys := []string{"totalBooks", "totalBookInstances", "totalBookInstancesAvailable", "totalAuthors", "totalGenres"}
	fns := []func(primitive.M) (int64, error){GetTotalBooksCount, GetTotalBookInstancesCount, GetTotalBookInstancesCount, GetTotalAuthorsCount, GetTotalGenresCount}
	args := []bson.M{bson.M{}, bson.M{}, bson.M{"status": "Available"}, bson.M{}, bson.M{}, bson.M{}}

	for i, fn := range fns {
		wg.Add(1)
		go func(index int, f func(primitive.M) (int64, error)) {
			defer wg.Done()
			defer commonRecover("GetAllCount")
			count, err := f(args[index])
			if err != nil {
				errs = append(errs, err)
			}
			// fmt.Println(totalKeys[index], index)
			allCounts[totalKeys[index]] = count
		}(i, fn)
	}

	wg.Wait()
	return allCounts, errs
}
