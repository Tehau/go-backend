# GO Backend 

Backend in Goland with Ricky and Morty data.

### Docker build

````
docker scan go-backend-rick
````

````
docker run --rm -p 5000:5000 go-backend-rick
````

### Local run

Change in file character.data.go at line 54, if run in a Docker
````
filePath := "data/rickandmortycharacter.json"
````

and if you run locally
````
filePath := "src/data/rickandmortycharacter.json"
````