package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/seppo0010/wikipedia-go"
)

var imageType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Image",
		Fields: graphql.Fields{
			"url": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(wikipedia.Image).Url, nil
				},
			},
			"title": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(wikipedia.Image).Title, nil
				},
			},
			"description_url": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(wikipedia.Image).DescriptionUrl, nil
				},
			},
		},
	},
)

var referenceType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Reference",
		Fields: graphql.Fields{
			"url": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(wikipedia.Reference).Url, nil
				},
			},
		},
	},
)

var pageType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Page",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(wikipedia.Page).Id()
				},
			},
			"title": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(wikipedia.Page).Title()
				},
			},
			"content": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(wikipedia.Page).Content()
				},
			},
			"html_content": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(wikipedia.Page).HtmlContent()
				},
			},
			"summary": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(wikipedia.Page).Summary()
				},
			},
			"images": &graphql.Field{
				Type: graphql.NewList(imageType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					images := make([]wikipedia.Image, 0, 100)
					for imageResult := range p.Source.(wikipedia.Page).Images() {
						if imageResult.Err != nil {
							return nil, imageResult.Err
						}
						images = append(images, imageResult.Image)
					}
					return images, nil
				},
			},
			"references": &graphql.Field{
				Type: graphql.NewList(referenceType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					references := make([]wikipedia.Reference, 0, 100)
					for referenceResult := range p.Source.(wikipedia.Page).Extlinks() {
						if referenceResult.Err != nil {
							return nil, referenceResult.Err
						}
						references = append(references, referenceResult.Reference)
					}
					return references, nil
				},
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"page": &graphql.Field{
				Type: pageType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"title": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					idQuery, isOK := p.Args["id"].(string)
					if isOK {
						return wiki.PageFromId(idQuery), nil
					}
					titleQuery, isOK := p.Args["title"].(string)
					if isOK {
						return wiki.Page(titleQuery), nil
					}
					return nil, nil
				},
			},
		},
	})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	return result
}

var wiki wikipedia.Wikipedia = wikipedia.NewWikipedia()

func main() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query()["query"][0], schema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Now server is running on port 8080")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8080/graphql?query={page(id:\"4138548\"){title}}'")
	http.ListenAndServe(":8080", nil)
}
