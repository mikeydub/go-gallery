package generator

import "text/template"

var dataloadersTemplate = template.Must(template.New("generated").
	Funcs(template.FuncMap{
		"lcFirst": lcFirst,
	}).
	Parse(`
// Code generated by github.com/mikeydub/go-gallery/cmd/dataloaders, DO NOT EDIT.

package {{.Package}}

import (
    "sync"
    "time"
    "context"
	"github.com/jackc/pgx/v4"
	"github.com/mikeydub/go-gallery/cmd/dataloaders/generator"

    {{range .ImportPaths}}
    "{{.}}"
    {{end}}
)

type autoCacheWithKey[TKey any, TResult any] interface {
	getKeyForResult(TResult) TKey
}

type autoCacheWithKeys[TKey any, TResult any] interface {
	getKeysForResult(TResult) []TKey
}

type PreFetchHook func(context.Context, string) context.Context
type PostFetchHook func(context.Context, string)

type NotFound[TKey any] struct {
	Key TKey
}

func (e NotFound[TKey]) Error() string {
	return fmt.Sprintf("result not found with key: %v", e.Key)
}

{{range .Definitions}}
// {{.Name}} batches and caches requests          
type {{.Name}} struct {
	generator.Dataloader[{{.KeyType.String}}, {{.ResultType.String}}]
}

// new{{.Name}} creates a new {{.Name}} with the given settings, functions, and options
func new{{.Name}}(
	ctx context.Context,
	maxBatchSize int,
	batchTimeout time.Duration,
	cacheResults bool,
	publishResults bool,
	fetch func(context.Context, []{{.KeyType.String}}) ([]{{.ResultType.String}}, []error),
	preFetchHook PreFetchHook,
	postFetchHook PostFetchHook,
	) *{{.Name}} {
	fetchWithHooks := func(ctx context.Context, keys []{{.KeyType.String}}) ([]{{.ResultType.String}}, []error) {
		// Allow the preFetchHook to modify and return a new context
		if preFetchHook != nil {
			ctx = preFetchHook(ctx, "{{.Name}}")
		}

		results, errors := fetch(ctx, keys)

		if postFetchHook != nil {
			postFetchHook(ctx, "{{.Name}}")
		}
		
		return results, errors
	}

	{{ if .KeyIsComparable }}
	dataloader := generator.NewDataloader(ctx, maxBatchSize, batchTimeout, cacheResults, publishResults, fetchWithHooks)
	{{ else }}
	dataloader := generator.NewDataloaderWithNonComparableKey(ctx, maxBatchSize, batchTimeout, cacheResults, publishResults, fetchWithHooks)
	{{ end }}
	d := &{{.Name}}{
		Dataloader: *dataloader,
	}
	
	return d
}

{{ if .CanAutoCacheDBID }}
func (*{{.Name}}) getKeyForResult(result {{.ResultType.String}}) persist.DBID {
	return result.ID
}
{{- end }}

{{ if .IsCustomBatch }}
func load{{.Name}}(q *coredb.Queries) func(context.Context, []{{.KeyType.String}}) ([]{{.CustomBatching.LoaderResultType.String}}, []error) {
	return func(ctx context.Context, params []{{.KeyType.String}}) ([]{{.CustomBatching.LoaderResultType.String}}, []error) {
		queryResults, err := q.{{.Name}}(ctx, params)

		results := make([]{{.CustomBatching.LoaderResultType.String}}, len(params))
		errors := make([]error, len(params))

		if err != nil {
			for i := range errors {
				errors[i] = err
			}

			return results, errors
		}


		hasResults := make([]bool, len(params))

		for _, result := range queryResults {
			results[result.{{.CustomBatching.BatchKeyIndexFieldName}}-1] = result{{.CustomBatching.ResultFieldName}}
			hasResults[result.{{.CustomBatching.BatchKeyIndexFieldName}}-1] = true
		}

		for i, hasResult := range hasResults {
			if !hasResult {
				errors[i] = NotFound[{{.KeyType.String}}]{Key: params[i]}
			}
		}

		return results, errors
	}
}
{{ else }}
func load{{.Name}}(q *coredb.Queries) func(context.Context, []{{.KeyType.String}}) ([]{{.ResultType.String}}, []error) {
	return func(ctx context.Context, params []{{.KeyType.String}}) ([]{{.ResultType.String}}, []error) {
		results := make([]{{.ResultType.String}}, len(params))
		errors := make([]error, len(params))

		b := q.{{.Name}}(ctx, params)
		defer b.Close()

		{{ if .ResultType.IsSlice }}
		b.Query(func(i int, r {{.ResultType.String}}, err error) {
		{{- else }}
		b.QueryRow(func(i int, r {{.ResultType.String}}, err error) {
		{{- end }}
			results[i], errors[i] = r, err
			if errors[i] == pgx.ErrNoRows {
				errors[i] = NotFound[{{.KeyType.String}}]{Key: params[i]}
			}
		})

		return results, errors
	}
}
{{end}}
{{end}}
`))

var apiTemplate = template.Must(template.New("api").
	Funcs(template.FuncMap{
		"lcFirst": lcFirst,
	}).
	Parse(`
// Code generated by github.com/mikeydub/go-gallery/cmd/dataloaders, DO NOT EDIT.

package {{.Package}}

import (
    "time"
    "context"

	// TODO: Add to ImportPaths
	//"github.com/mikeydub/go-gallery/db/gen/coredb"
    
	{{range .ImportPaths}}
    "{{.}}"
    {{end}}
)

type Loaders struct {
{{range .Definitions}}
	{{.Name}} *{{.Name}}
{{- end}}
}

func NewLoaders(ctx context.Context, q *coredb.Queries, disableCaching bool, preFetchHook PreFetchHook, postFetchHook PostFetchHook) *Loaders {
	loaders := &Loaders{}

	{{range .Definitions}}
	loaders.{{.Name}} = new{{.Name}}(ctx, {{.MaxBatchSize}}, time.Duration({{.BatchTimeout}}), !disableCaching, {{.PublishResults}}, load{{.Name}}(q), preFetchHook, postFetchHook)
	{{- end}}

	{{range .Subscriptions}}
	{{- if .SingleKey }}
	loaders.{{.Target}}.RegisterResultSubscriber(func(result {{.Result}}) {
	{{- if .ResultIsSlice }}
	for _, entry := range result {
		loaders.{{.Subscriber}}.Prime(loaders.{{.Subscriber}}.getKeyForResult(entry{{.Field}}), entry{{.Field}})
	}
	{{- else }}
	loaders.{{.Subscriber}}.Prime(loaders.{{.Subscriber}}.getKeyForResult(result{{.Field}}), result{{.Field}})
	{{- end }}
	})
	{{- end }}
	{{- if .ManyKeys }}
	loaders.{{.Target}}.RegisterResultSubscriber(func(result {{.Result}}) {
	{{- if .ResultIsSlice }}
	for _, entry := range result {
		for _, key := range loaders.{{.Subscriber}}.getKeysForResult(entry{{.Field}}) {
			loaders.{{.Subscriber}}.Prime(key, entry{{.Field}})
		}
	}
	{{- else }}
	for _, key := range loaders.{{.Subscriber}}.getKeysForResult(result{{.Field}}) {
		loaders.{{.Subscriber}}.Prime(key, result{{.Field}})
	}
	{{- end }}
	})
	{{- end }}
	{{- end }}

	return loaders
}
`))
