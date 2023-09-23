package Interface

import (
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
)

func ReadOpenapiDocs(docsUrl string) (openapi3.Paths, error) {
	loader := openapi3.NewLoader()

	openApiData, err := getData(docsUrl)
	if err != nil {
		return nil, err
	}

	doc, err := loadFromData(openApiData, loader)
	if err != nil {
		return nil, err
	}

	if err = doc.Validate(loader.Context); err != nil {
		return nil, err
	}

	return doc.Paths, err
}

func getData(docsUrl string) ([]byte, error) {
	request, _ := CreateRequest(docsUrl, "GET", "")
	status, openApiData, _ := FullRequest(request)

	if status != 200 {
		errorMsg := fmt.Sprintf("openapi: getting openapi docs (HTTP STATUS not 200, status: %s)", status)
		err := errors.New(errorMsg)

		data := string(openApiData)
		fmt.Println(data)

		return nil, err
	}

	return openApiData, nil
}

func loadFromData(data []byte, loader *openapi3.Loader) (*openapi3.T, error) {
	doc, err := loader.LoadFromData(data)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
