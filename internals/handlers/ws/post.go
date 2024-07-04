package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	conf "gateway/config"
	"gateway/internals/tools"
	md "gateway/model"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

func createPost(w http.ResponseWriter, data string, userId int) int {
	var post md.Post
	post.UserId = userId
	if err := json.Unmarshal([]byte(data), &post); err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0
	}

	postPath := "./static/images/posts/"
	filename, err := decodeBase64Image(post.Img, postPath)
	post.Img = filename
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0
	}

	bodyData := md.RequestBody{
		Action: "createPost",
		Body:   post,
	}

	resp, err := tools.SendRequest(w, bodyData, "POST", conf.URLPost)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func getOnePost(w http.ResponseWriter, data string) (int, string) {
	var post md.Post
	if err := json.Unmarshal([]byte(data), &post); err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0, ""
	}

	bodyData := md.RequestBody{
		Action: "getOne",
		Body:   post,
	}

	resp, err := tools.SendRequest(w, bodyData, "POST", conf.URLPost)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0, ""
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0, ""
	}

	return resp.StatusCode, string(responseBody)
}

func getAllPost(w http.ResponseWriter) (int, string) {
	bodyData := md.RequestBody{
		Action: "getAll",
	}
	resp, err := tools.SendRequest(w, bodyData, "POST", conf.URLPost)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0, ""
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0, ""
	}

	var posts []md.Post
	if err := json.Unmarshal([]byte(responseBody), &posts); err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0, ""
	}

	for i := 0; i < len(posts); i++ {
		base64Code, err := encodeImageToBase64(posts[i].Img)
		if err != nil {
			fmt.Println("Internal server error: " + err.Error())
			return 0, ""
		}

		posts[i].ImgBase64 = base64Code
	}

	response, err := json.Marshal(posts)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0, ""
	}

	return resp.StatusCode, string(response)
}

func deletePost(w http.ResponseWriter, data string) int {
	var post md.Post
	if err := json.Unmarshal([]byte(data), &post); err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0
	}

	bodyData := md.RequestBody{
		Action: "delete",
		Body:   post,
	}

	resp, err := tools.SendRequest(w, bodyData, "POST", conf.URLPost)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func getLastPost(w http.ResponseWriter) (int, string) {
	status, posts := getAllPost(w)
	if status != http.StatusOK {
		return 0, ""
	}

	var tabPost []md.Post
	if err := json.Unmarshal([]byte(posts), &tabPost); err != nil {
		return 0, ""
	}

	var tabToSend = []md.Post{tabPost[0]}
	toSend, err := json.Marshal(tabToSend)
	if err != nil {
		return 0, ""
	}

	return 200, string(toSend)
}

func decodeBase64Image(base64Image string, outputFilePath string) (string, error) {

	fmt.Println("Début du décodage de l'image")

	// Vérifier si la chaîne commence par "data:"
	if strings.HasPrefix(base64Image, "data:") {
		// Extraire la partie Base64 après la virgule
		parts := strings.SplitN(base64Image, ",", 2)
		if len(parts) != 2 {
			return "", errors.New("format de chaîne data: invalide")
		}
		base64Image = parts[1]
	}

	fmt.Println("Base64Image après extraction:", base64Image[:30]) // Afficher les premiers caractères pour vérifier

	// Décoder la chaîne Base64 en un tableau d'octets
	data, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		fmt.Println("Erreur lors du décodage Base64:", err)
		return "", err
	}

	fmt.Println("Décodage Base64 réussi. Longueur des données:", len(data))

	// Vérifier le contenu décodé (uniquement les premiers octets pour éviter une sortie massive)

	fileName := outputFilePath + uuid.New().String() + ".jpeg"

	// Écrire les octets décodés dans un fichier
	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		fmt.Println("Erreur lors de l'écriture du fichier:", err)
		return "", err
	}

	fmt.Println("Écriture du fichier réussie:", outputFilePath)
	return fileName, nil
}

func encodeImageToBase64(imagePath string) (string, error) {
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la lecture du fichier image: %w", err)
	}

	base64Encoded := base64.StdEncoding.EncodeToString(imageBytes)

	return base64Encoded, nil
}
