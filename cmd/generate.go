package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"code-gen-go/utils"
	"github.com/dave/jennifer/jen"
)

const (
	inSrcFileName  = "../data.json"
	outPackageName = "entity"
)

const entityTypeField = "type"

//go:generate go run generate.go
func main() {
	filesToCreate := map[string]*jen.File{}

	incomingData, err := getFileContent(inSrcFileName)
	if err != nil {
		log.Fatal(err)
	}

	entityProps := getEntityPropsMap(incomingData)
	for entityName, entityPropsMap := range entityProps {
		goFile := jen.NewFile(outPackageName)
		goFile.Comment("Code generated by go generate; DO NOT EDIT.")

		/* Add entity type constant, as example:
		const TYPE_SPORT = "Sport"
		*/
		typeConstant := jen.Const().Defs(
			jen.Id("TYPE_" + strings.ToUpper(entityName)).Op("=").Lit(entityName).Line(),
		)
		goFile.Add(typeConstant)

		/* Add the main structure, as example:
		type Sport struct {
			Id       int    `db:"id" json:"id"`
			Title    string `db:"title" json:"title"`
			Short    string `db:"short" json:"short"`
			IsActive bool   `db:"isActive" json:"isActive"`
		}
		*/
		goStruct := jen.Type().Id(entityName).StructFunc(
			func(g *jen.Group) {
				for propName, propVal := range entityPropsMap {
					if propName == "type" {
						continue
					}
					goField := jen.Id(strings.Title(propName))

					switch propVal.(type) {
					case int:
						goField.Int()
					case float64:
						goField.Int()
					case string:
						goField.String()
					case bool:
						goField.Bool()
					default:
						continue
					}

					goField.Tag(map[string]string{"json": propName, "db": utils.ToSnakeCase(propName)})

					g.Add(goField)
				}
			},
		).Line()
		goFile.Add(goStruct)

		// Add simple functions depended on the property type
		for propName, propVal := range entityPropsMap {
			if propName == "type" {
				continue
			}

			switch propVal.(type) {
			case bool:
				/*
					func (e *League) CheckIsActive() bool {
						return e.IsActive == true
					}
				*/
				goFieldFunc := jen.Func().Params(
					jen.Id("e").Id("*" + strings.Title(entityName)),
				).Id("Check" + strings.Title(propName)).Params().Bool().Block(
					jen.Return(jen.Id("e." + strings.Title(propName)).Op("==").True()),
				).Line()

				goFile.Add(goFieldFunc)
			default:
				continue
			}
		}

		filesToCreate[entityName] = goFile
	}

	if err := writeFiles(filesToCreate); err != nil {
		log.Fatal(err)
	}
}

func getFileContent(fileName string) ([]map[string]interface{}, error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var entities []map[string]interface{}
	if err := json.Unmarshal(b, &entities); err != nil {
		return nil, err
	}

	return entities, err
}

func writeFiles(filesToCreate map[string]*jen.File) error {
	for fileName, f := range filesToCreate {
		if err := f.Save("../" + outPackageName + "/" + fileName + ".go"); err != nil {
			return err
		}
	}

	return nil
}

func getEntityPropsMap(entities []map[string]interface{}) map[string]map[string]interface{} {
	props := map[string]map[string]interface{}{} // map[entityName]map[propertyName]propertyValue

	for _, entity := range entities {
		// prepare entities types
		var entityName string
		for propName, propVal := range entity {
			if propName == entityTypeField {
				entityName = strings.Title(fmt.Sprintf("%v", propVal))
			}
		}

		if entityName == "" { // don't continue with no entity name
			continue
		}

		// prepare properties map - gather all properties even if a property appears in the one entity only
		propsMap, ok := props[entityName]
		if !ok {
			propsMap = map[string]interface{}{}
		}
		for propName, propVal := range entity {
			propsMap[propName] = propVal
		}
		props[entityName] = propsMap
	}

	return props
}
